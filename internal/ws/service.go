// Package ws provides WebSocket service functionality
package ws

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/utmos/utmos/internal/ws/hub"
	"github.com/utmos/utmos/internal/ws/push"
	"github.com/utmos/utmos/internal/ws/subscription"
	"github.com/utmos/utmos/pkg/metrics"
	"github.com/utmos/utmos/pkg/rabbitmq"
)

// WebSocket upgrader configuration defaults
const (
	// DefaultReadBufferSize is the default read buffer size for WebSocket connections
	DefaultReadBufferSize = 1024
	// DefaultWriteBufferSize is the default write buffer size for WebSocket connections
	DefaultWriteBufferSize = 1024
)

// ServiceConfig holds WebSocket service configuration
type ServiceConfig struct {
	// Hub configuration
	HubConfig *hub.Config
	// Client configuration
	ClientConfig *hub.ClientConfig
	// Pusher configuration
	PusherConfig *push.Config
	// AllowedOrigins for WebSocket upgrade
	AllowedOrigins []string
}

// DefaultServiceConfig returns default service configuration
func DefaultServiceConfig() *ServiceConfig {
	return &ServiceConfig{
		HubConfig:      hub.DefaultConfig(),
		ClientConfig:   hub.DefaultClientConfig(),
		PusherConfig:   push.DefaultConfig(),
		AllowedOrigins: []string{"*"},
	}
}

// Service is the WebSocket service
type Service struct {
	config     *ServiceConfig
	hub        *hub.Hub
	subManager *subscription.Manager
	pusher     *push.Pusher
	subscriber *rabbitmq.Subscriber
	upgrader   websocket.Upgrader
	logger     *logrus.Entry
	msgMetrics *metrics.MessageMetrics

	running   bool
	runningMu sync.RWMutex
	done      chan struct{}
}

// NewService creates a new WebSocket service
func NewService(config *ServiceConfig, metricsCollector *metrics.Collector, logger *logrus.Entry) *Service {
	if config == nil {
		config = DefaultServiceConfig()
	}
	if logger == nil {
		logger = logrus.NewEntry(logrus.StandardLogger())
	}

	var msgMetrics *metrics.MessageMetrics
	if metricsCollector != nil {
		msgMetrics = metrics.NewMessageMetrics(metricsCollector)
	}

	// Create hub
	h := hub.NewHub(config.HubConfig, logger)

	// Create subscription manager
	subManager := subscription.NewManager(logger)

	// Create pusher
	pusher := push.NewPusher(config.PusherConfig, h, subManager, logger)

	// Create upgrader
	upgrader := websocket.Upgrader{
		ReadBufferSize:  DefaultReadBufferSize,
		WriteBufferSize: DefaultWriteBufferSize,
		CheckOrigin: func(r *http.Request) bool {
			if len(config.AllowedOrigins) == 0 {
				return true
			}
			origin := r.Header.Get("Origin")
			for _, allowed := range config.AllowedOrigins {
				if allowed == "*" || allowed == origin {
					return true
				}
			}
			return false
		},
	}

	svc := &Service{
		config:     config,
		hub:        h,
		subManager: subManager,
		pusher:     pusher,
		upgrader:   upgrader,
		logger:     logger.WithField("component", "ws-service"),
		msgMetrics: msgMetrics,
		done:       make(chan struct{}),
	}

	// Set hub callbacks
	h.SetOnConnect(svc.onClientConnect)
	h.SetOnDisconnect(svc.onClientDisconnect)
	h.SetOnMessage(svc.onClientMessage)

	return svc
}

// SetSubscriber sets the RabbitMQ subscriber for receiving messages
func (s *Service) SetSubscriber(subscriber *rabbitmq.Subscriber) {
	s.subscriber = subscriber
}

// Start starts the WebSocket service
func (s *Service) Start(ctx context.Context) error {
	s.runningMu.Lock()
	if s.running {
		s.runningMu.Unlock()
		return nil
	}
	s.running = true
	s.done = make(chan struct{})
	s.runningMu.Unlock()

	s.logger.Info("Starting WebSocket service")

	// Start hub
	s.hub.Start()

	// Start pusher
	s.pusher.Start()

	// Start RabbitMQ consumer if configured
	if s.subscriber != nil {
		go s.consumeMessages(ctx)
	}

	return nil
}

// Stop stops the WebSocket service
func (s *Service) Stop() error {
	s.runningMu.Lock()
	if !s.running {
		s.runningMu.Unlock()
		return nil
	}
	s.running = false
	s.runningMu.Unlock()

	close(s.done)

	// Stop pusher
	s.pusher.Stop()

	// Stop hub
	s.hub.Stop()

	s.logger.Info("WebSocket service stopped")
	return nil
}

// IsRunning returns whether the service is running
func (s *Service) IsRunning() bool {
	s.runningMu.RLock()
	defer s.runningMu.RUnlock()
	return s.running
}

// HandleWebSocket handles WebSocket upgrade requests
func (s *Service) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	if !s.IsRunning() {
		http.Error(w, "Service not running", http.StatusServiceUnavailable)
		return
	}

	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.logger.WithError(err).Warn("Failed to upgrade WebSocket connection")
		return
	}

	// Generate client ID
	clientID := uuid.New().String()

	// Extract metadata from request
	deviceSN := r.URL.Query().Get("device_sn")
	userID := r.URL.Query().Get("user_id")

	// Create client
	client := hub.NewClient(clientID, conn, s.hub, s.config.ClientConfig, s.logger)
	client.DeviceSN = deviceSN
	client.UserID = userID

	// Register with hub
	s.hub.Register(client)

	// Start client pumps
	client.Start()

	s.logger.WithFields(logrus.Fields{
		"client_id": clientID,
		"device_sn": deviceSN,
		"user_id":   userID,
	}).Info("WebSocket client connected")
}

// onClientConnect is called when a client connects
func (s *Service) onClientConnect(client *hub.Client) {
	s.logger.WithField("client_id", client.ID).Debug("Client connected callback")
}

// onClientDisconnect is called when a client disconnects
func (s *Service) onClientDisconnect(client *hub.Client) {
	// Unsubscribe from all topics
	s.subManager.UnsubscribeAll(client.ID)

	s.logger.WithField("client_id", client.ID).Debug("Client disconnected callback")
}

// onClientMessage is called when a message is received from a client
func (s *Service) onClientMessage(client *hub.Client, msg *hub.Message) {
	s.logger.WithFields(logrus.Fields{
		"client_id": client.ID,
		"msg_type":  msg.Type,
		"event":     msg.Event,
	}).Debug("Client message received")

	// Handle subscription messages
	switch msg.Type {
	case hub.MessageTypeSubscribe:
		if msg.Event != "" {
			s.subManager.Subscribe(client.ID, msg.Event)
		}
	case hub.MessageTypeUnsubscribe:
		if msg.Event != "" {
			s.subManager.Unsubscribe(client.ID, msg.Event)
		}
	}
}

// consumeMessages consumes messages from RabbitMQ and pushes to clients
func (s *Service) consumeMessages(ctx context.Context) {
	if s.subscriber == nil {
		return
	}

	// Subscribe to the WebSocket queue
	err := s.subscriber.Subscribe("iot.ws.queue", func(ctx context.Context, msg *rabbitmq.StandardMessage) error {
		s.handleRabbitMQMessage(ctx, msg)
		return nil
	})
	if err != nil {
		s.logger.WithError(err).Error("Failed to subscribe to RabbitMQ")
		return
	}

	// Wait for done signal
	select {
	case <-s.done:
	case <-ctx.Done():
	}

	// Unsubscribe when done
	if err := s.subscriber.Unsubscribe("iot.ws.queue"); err != nil {
		return
	}
}

// handleRabbitMQMessage handles a message from RabbitMQ
func (s *Service) handleRabbitMQMessage(ctx context.Context, msg *rabbitmq.StandardMessage) {
	if msg == nil {
		return
	}

	tr := otel.Tracer("iot-ws")
	_, span := tr.Start(ctx, "ws.message.push",
		trace.WithAttributes(
			attribute.String("device_sn", msg.DeviceSN),
			attribute.String("action", msg.Action),
		),
	)
	defer span.End()

	// Determine topic from action
	topic := s.actionToTopic(msg.Action)

	// Create WebSocket message
	wsMsg := &hub.Message{
		Type:    hub.MessageTypeEvent,
		Event:   topic,
		TraceID: msg.TID,
	}

	// Parse data
	if msg.Data != nil {
		var data any
		if err := json.Unmarshal(msg.Data, &data); err == nil {
			wsMsg.Data = data
		} else {
			wsMsg.Data = string(msg.Data)
		}
	}

	// Push to subscribers
	s.pusher.PushToTopic(topic, wsMsg)

	// Also push to device-specific topic if device_sn is present
	if msg.DeviceSN != "" {
		deviceTopic := "device." + msg.DeviceSN + "." + topic
		s.pusher.PushToTopic(deviceTopic, wsMsg)
	}

	if s.msgMetrics != nil {
		s.msgMetrics.ProcessedTotal.WithLabelValues("iot-ws", "", msg.Action, "success").Inc()
	}
}

// actionToTopic converts an action to a WebSocket topic
func (s *Service) actionToTopic(action string) string {
	return action
}

// Hub returns the WebSocket hub
func (s *Service) Hub() *hub.Hub {
	return s.hub
}

// SubscriptionManager returns the subscription manager
func (s *Service) SubscriptionManager() *subscription.Manager {
	return s.subManager
}

// Pusher returns the message pusher
func (s *Service) Pusher() *push.Pusher {
	return s.pusher
}

// GetStats returns service statistics
func (s *Service) GetStats() ServiceStats {
	pushed, dropped := s.pusher.GetMetrics()
	return ServiceStats{
		ConnectedClients: s.hub.GetClientCount(),
		ActiveTopics:     s.subManager.GetTopicCount(),
		MessagesPushed:   pushed,
		MessagesDropped:  dropped,
		QueueLength:      s.pusher.QueueLength(),
	}
}

// ServiceStats holds service statistics
type ServiceStats struct {
	ConnectedClients int   `json:"connected_clients"`
	ActiveTopics     int   `json:"active_topics"`
	MessagesPushed   int64 `json:"messages_pushed"`
	MessagesDropped  int64 `json:"messages_dropped"`
	QueueLength      int   `json:"queue_length"`
}
