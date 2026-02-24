// Package gateway provides the IoT gateway service
package gateway

import (
	"context"
	"fmt"
	"sync"
	"time"

	pahomqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"

	"github.com/utmos/utmos/internal/gateway/bridge"
	"github.com/utmos/utmos/internal/gateway/connection"
	"github.com/utmos/utmos/internal/gateway/mqtt"
	"github.com/utmos/utmos/pkg/rabbitmq"
)

// Gateway service configuration defaults
const (
	// DefaultCleanupInterval is the default interval for cleaning up stale connections
	DefaultCleanupInterval = 5 * time.Minute
	// DefaultMaxStaleAge is the default maximum age for stale connections
	DefaultMaxStaleAge = 24 * time.Hour
	// DefaultDisconnectQuiesce is the default quiesce time for MQTT disconnect (ms)
	DefaultDisconnectQuiesce = 1000
)

// ServiceConfig holds configuration for the gateway service
type ServiceConfig struct {
	// MQTT configuration
	MQTT *mqtt.Config

	// Bridge configuration
	UplinkBridge   *bridge.UplinkBridgeConfig
	DownlinkBridge *bridge.DownlinkBridgeConfig

	// Connection cleanup configuration
	CleanupInterval time.Duration
	MaxStaleAge     time.Duration

	// Topic subscriptions
	SubscribeTopics []string
}

// DefaultServiceConfig returns default service configuration
func DefaultServiceConfig() *ServiceConfig {
	return &ServiceConfig{
		MQTT:            mqtt.DefaultConfig(),
		UplinkBridge:    bridge.DefaultUplinkBridgeConfig(),
		DownlinkBridge:  bridge.DefaultDownlinkBridgeConfig(),
		CleanupInterval: DefaultCleanupInterval,
		MaxStaleAge:     DefaultMaxStaleAge,
		SubscribeTopics: []string{
			"thing/product/+/+/#",
			"sys/product/+/#",
		},
	}
}

// Service is the main gateway service
type Service struct {
	config *ServiceConfig
	logger *logrus.Entry

	// Components
	mqttClient    *mqtt.Client
	mqttHandler   *mqtt.Handler
	uplinkBridge  *bridge.UplinkBridge
	downlinkBridge *bridge.DownlinkBridge
	connManager   *connection.Manager

	// RabbitMQ
	publisher  *rabbitmq.Publisher
	subscriber *rabbitmq.Subscriber

	// State
	running bool
	mu      sync.RWMutex
	cancel  context.CancelFunc
}

// NewService creates a new gateway service
func NewService(
	config *ServiceConfig,
	publisher *rabbitmq.Publisher,
	subscriber *rabbitmq.Subscriber,
	logger *logrus.Entry,
) *Service {
	if config == nil {
		config = DefaultServiceConfig()
	}
	if logger == nil {
		logger = logrus.NewEntry(logrus.StandardLogger())
	}

	svcLogger := logger.WithField("component", "gateway-service")

	// Create MQTT client
	mqttClient := mqtt.NewClient(config.MQTT, svcLogger)

	// Create MQTT handler
	mqttHandler := mqtt.NewHandler(svcLogger)

	// Create connection manager
	connManager := connection.NewManager(svcLogger)

	// Create bridges
	uplinkBridge := bridge.NewUplinkBridge(publisher, config.UplinkBridge, svcLogger)
	downlinkBridge := bridge.NewDownlinkBridge(mqttClient, subscriber, config.DownlinkBridge, svcLogger)

	return &Service{
		config:         config,
		logger:         svcLogger,
		mqttClient:     mqttClient,
		mqttHandler:    mqttHandler,
		uplinkBridge:   uplinkBridge,
		downlinkBridge: downlinkBridge,
		connManager:    connManager,
		publisher:      publisher,
		subscriber:     subscriber,
	}
}

// Start starts the gateway service
func (s *Service) Start(ctx context.Context) error {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return fmt.Errorf("service already running")
	}
	s.running = true
	s.mu.Unlock()

	ctx, s.cancel = context.WithCancel(ctx)

	s.logger.Info("Starting gateway service")

	// Setup MQTT client handlers
	s.setupMQTTHandlers()

	// Connect to MQTT broker
	if err := s.mqttClient.Connect(ctx); err != nil {
		s.mu.Lock()
		s.running = false
		s.mu.Unlock()
		return fmt.Errorf("failed to connect to MQTT broker: %w", err)
	}

	// Subscribe to topics
	if err := s.subscribeToTopics(); err != nil {
		s.mqttClient.Disconnect(DefaultDisconnectQuiesce)
		s.mu.Lock()
		s.running = false
		s.mu.Unlock()
		return fmt.Errorf("failed to subscribe to topics: %w", err)
	}

	// Start downlink bridge
	if err := s.downlinkBridge.Start(ctx); err != nil {
		s.logger.WithError(err).Warn("Failed to start downlink bridge")
	}

	// Start connection cleanup routine
	s.connManager.StartCleanupRoutine(ctx, s.config.CleanupInterval, s.config.MaxStaleAge)

	s.logger.Info("Gateway service started")
	return nil
}

// Stop stops the gateway service
func (s *Service) Stop() error {
	s.mu.Lock()
	if !s.running {
		s.mu.Unlock()
		return nil
	}
	s.running = false
	s.mu.Unlock()

	s.logger.Info("Stopping gateway service")

	// Cancel context
	if s.cancel != nil {
		s.cancel()
	}

	// Stop downlink bridge
	s.downlinkBridge.Stop()

	// Disconnect MQTT client
	s.mqttClient.Disconnect(1000)

	s.logger.Info("Gateway service stopped")
	return nil
}

// setupMQTTHandlers configures MQTT client handlers
func (s *Service) setupMQTTHandlers() {
	// Set message handler to route through mqtt.Handler
	s.mqttClient.SetMessageHandler(func(client *mqtt.Client, msg pahomqtt.Message) {
		s.mqttHandler.Handle(nil, msg)
	})

	// Set connect handler
	s.mqttClient.SetConnectHandler(func(client *mqtt.Client) {
		s.logger.Info("MQTT connection established")
	})

	// Set connection lost handler
	s.mqttClient.SetConnectionLostHandler(func(client *mqtt.Client, err error) {
		s.logger.WithError(err).Warn("MQTT connection lost")
	})

	// Register uplink bridge processor
	uplinkProcessor := s.uplinkBridge.CreateProcessor("#")
	s.mqttHandler.RegisterProcessor(uplinkProcessor)

	// Setup connection tracking callbacks
	s.connManager.SetOnConnect(func(state *connection.DeviceState) {
		s.logDeviceStateChange(state, "Device connected")
	})

	s.connManager.SetOnDisconnect(func(state *connection.DeviceState) {
		s.logDeviceStateChange(state, "Device disconnected")
	})
}

// logDeviceStateChange logs a device connection state change event
func (s *Service) logDeviceStateChange(state *connection.DeviceState, event string) {
	s.logger.WithFields(logrus.Fields{
		"device_sn": state.DeviceSN,
		"client_id": state.ClientID,
	}).Info(event)
}

// subscribeToTopics subscribes to configured MQTT topics
func (s *Service) subscribeToTopics() error {
	for _, topic := range s.config.SubscribeTopics {
		if err := s.mqttClient.Subscribe(topic, s.config.MQTT.QoS, nil); err != nil {
			return fmt.Errorf("failed to subscribe to %s: %w", topic, err)
		}
		s.logger.WithField("topic", topic).Debug("Subscribed to topic")
	}
	return nil
}

// IsRunning returns whether the service is running
func (s *Service) IsRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.running
}

// IsMQTTConnected returns whether the MQTT client is connected
func (s *Service) IsMQTTConnected() bool {
	return s.mqttClient.IsConnected()
}

// GetOnlineDeviceCount returns the number of online devices
func (s *Service) GetOnlineDeviceCount() int {
	return s.connManager.GetOnlineCount()
}

// GetConnectionManager returns the connection manager
func (s *Service) GetConnectionManager() *connection.Manager {
	return s.connManager
}

// GetMQTTClient returns the MQTT client
func (s *Service) GetMQTTClient() *mqtt.Client {
	return s.mqttClient
}

// GetMQTTHandler returns the MQTT handler
func (s *Service) GetMQTTHandler() *mqtt.Handler {
	return s.mqttHandler
}

// RegisterDevice registers a device connection
func (s *Service) RegisterDevice(deviceSN, clientID, ipAddress string) *connection.DeviceState {
	return s.connManager.Connect(deviceSN, clientID, ipAddress)
}

// UnregisterDevice unregisters a device connection
func (s *Service) UnregisterDevice(deviceSN string) *connection.DeviceState {
	return s.connManager.Disconnect(deviceSN)
}

// IsDeviceOnline checks if a device is online
func (s *Service) IsDeviceOnline(deviceSN string) bool {
	return s.connManager.IsOnline(deviceSN)
}
