// Package downlink provides the iot-downlink service implementation
package downlink

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/utmos/utmos/internal/downlink/dispatcher"
	"github.com/utmos/utmos/internal/downlink/model"
	"github.com/utmos/utmos/internal/downlink/retry"
	"github.com/utmos/utmos/internal/downlink/router"
	"github.com/utmos/utmos/pkg/adapter"
	"github.com/utmos/utmos/pkg/metrics"
	"github.com/utmos/utmos/pkg/models"
	"github.com/utmos/utmos/pkg/rabbitmq"
	"github.com/utmos/utmos/pkg/repository"
)

// Config holds downlink service configuration
type Config struct {
	// Retry configuration
	RetryConfig *retry.Config

	// Router configuration
	RouterConfig *router.Config

	// EnableRetry enables retry mechanism
	EnableRetry bool

	// EnableRouting enables routing to gateway
	EnableRouting bool

	// RetryWorkerInterval is the interval for retry worker
	RetryWorkerInterval time.Duration
}

// DefaultConfig returns default service configuration
func DefaultConfig() *Config {
	return &Config{
		RetryConfig:         retry.DefaultConfig(),
		RouterConfig:        router.DefaultConfig(),
		EnableRetry:         true,
		EnableRouting:       true,
		RetryWorkerInterval: 5 * time.Second,
	}
}

// Service is the main downlink service
type Service struct {
	config       *Config
	logger       *logrus.Entry
	db           *gorm.DB
	registry     *dispatcher.Registry
	handler      *dispatcher.DispatchHandler
	retryHandler *retry.Handler
	router       *router.Router
	publisher    *rabbitmq.Publisher
	subscriber   *rabbitmq.Subscriber
	repository   *model.ServiceCallRepository
	msgLogRepo   *repository.MessageLogRepository

	mu       sync.RWMutex
	running  bool
	cancelFn context.CancelFunc

	// Metrics
	processedCount int64
	failedCount    int64
	msgMetrics     *metrics.MessageMetrics
}

// NewService creates a new downlink service
func NewService(config *Config, publisher *rabbitmq.Publisher, metricsCollector *metrics.Collector, logger *logrus.Entry, db *gorm.DB) *Service {
	if config == nil {
		config = DefaultConfig()
	}
	if logger == nil {
		logger = logrus.NewEntry(logrus.StandardLogger())
	}

	serviceLogger := logger.WithField("service", "iot-downlink")

	// Create dispatcher registry
	registry := dispatcher.NewRegistry(serviceLogger)

	// Create dispatch handler
	handler := dispatcher.NewDispatchHandler(registry, serviceLogger)

	// Create retry handler
	var retryHandler *retry.Handler
	if config.EnableRetry {
		retryHandler = retry.NewHandler(config.RetryConfig, serviceLogger)
	}

	// Create router
	var routerInstance *router.Router
	if config.EnableRouting {
		routerInstance = router.NewRouter(publisher, config.RouterConfig, serviceLogger)
	}

	// Create message metrics
	var msgMetrics *metrics.MessageMetrics
	if metricsCollector != nil {
		msgMetrics = metrics.NewMessageMetrics(metricsCollector)
	}

	// Create repository if DB is available
	var repo *model.ServiceCallRepository
	var msgLogRepo *repository.MessageLogRepository
	if db != nil {
		repo = model.NewServiceCallRepository(db)
		msgLogRepo = repository.NewMessageLogRepository(db)
	}

	svc := &Service{
		config:       config,
		logger:       serviceLogger,
		db:           db,
		registry:     registry,
		handler:      handler,
		retryHandler: retryHandler,
		router:       routerInstance,
		publisher:    publisher,
		msgMetrics:   msgMetrics,
		repository:   repo,
		msgLogRepo:   msgLogRepo,
	}

	// Set up callbacks
	handler.SetOnDispatched(svc.onDispatched)

	if retryHandler != nil {
		retryHandler.SetOnRetry(svc.onRetry)
		retryHandler.SetOnDeadLetter(svc.onDeadLetter)
	}

	return svc
}

// RegisterDispatcher registers a vendor dispatcher
func (s *Service) RegisterDispatcher(d dispatcher.Dispatcher) {
	s.registry.Register(d)
}

// RegisterAdapterDispatcher registers an adapter.DownlinkDispatcher by wrapping it
func (s *Service) RegisterAdapterDispatcher(a adapter.DownlinkDispatcher) {
	s.registry.Register(dispatcher.NewAdapterDispatcher(a))
}

// SetSubscriber sets the RabbitMQ subscriber
func (s *Service) SetSubscriber(subscriber *rabbitmq.Subscriber) {
	s.subscriber = subscriber
}

// Start starts the downlink service
func (s *Service) Start(ctx context.Context) error {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return fmt.Errorf("service already running")
	}

	ctx, cancel := context.WithCancel(ctx)
	s.cancelFn = cancel
	s.running = true
	s.mu.Unlock()

	s.logger.Info("Starting downlink service")

	// Start retry worker if enabled
	if s.config.EnableRetry && s.retryHandler != nil {
		s.retryHandler.StartRetryWorker(ctx, s.config.RetryWorkerInterval)
	}

	// Start consuming messages if subscriber is set
	if s.subscriber != nil {
		go s.consumeMessages(ctx)
	}

	s.logger.Info("Downlink service started")
	return nil
}

// Stop stops the downlink service
func (s *Service) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return nil
	}

	s.logger.Info("Stopping downlink service")

	if s.cancelFn != nil {
		s.cancelFn()
	}

	s.running = false
	s.logger.Info("Downlink service stopped")
	return nil
}

// IsRunning returns whether the service is running
func (s *Service) IsRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.running
}

// Dispatch dispatches a service call
func (s *Service) Dispatch(ctx context.Context, call *dispatcher.ServiceCall) (*dispatcher.DispatchResult, error) {
	if call == nil {
		return nil, fmt.Errorf("service call is nil")
	}

	s.logger.WithFields(logrus.Fields{
		"device_sn": call.DeviceSN,
		"vendor":    call.Vendor,
		"method":    call.Method,
	}).Debug("Dispatching service call")

	result, err := s.handler.Handle(ctx, call)
	if err != nil {
		s.incrementFailed()

		// Schedule retry if enabled
		if s.config.EnableRetry && s.retryHandler != nil {
			s.retryHandler.ScheduleRetry(call, err.Error())
		}

		return result, err
	}

	s.incrementProcessed()
	return result, nil
}

// consumeMessages consumes service reply messages from RabbitMQ
func (s *Service) consumeMessages(ctx context.Context) {
	s.logger.Info("Starting service reply consumer")

	// Set up the service reply queue
	queueName := "iot.downlink.service.reply"
	routingKeyPattern := "iot.*.device.service.reply"

	// Setup queue with binding if subscriber has access to client
	if s.subscriber != nil && s.publisher != nil {
		// Get the underlying client from subscriber to setup queue
		// The subscriber holds a reference to the client
		if err := s.setupServiceReplyQueue(queueName, routingKeyPattern); err != nil {
			s.logger.WithError(err).Warn("Failed to setup service reply queue, will rely on existing queue")
		}
	}

	// Subscribe to consume service replies
	if err := s.subscriber.Subscribe(queueName, s.handleServiceReply); err != nil {
		s.logger.WithError(err).Error("Failed to subscribe to service reply queue")
		return
	}

	s.logger.WithFields(logrus.Fields{
		"queue":       queueName,
		"routing_key": routingKeyPattern,
	}).Info("Subscribed to service reply queue")

	// Wait for context cancellation
	<-ctx.Done()

	// Unsubscribe when done
	if err := s.subscriber.Unsubscribe(queueName); err != nil {
		s.logger.WithError(err).Warn("Failed to unsubscribe from service reply queue")
	}

	s.logger.Info("Service reply consumer stopped")
}

// setupServiceReplyQueue declares and binds the service reply queue
func (s *Service) setupServiceReplyQueue(queueName, routingKeyPattern string) error {
	if s.subscriber == nil {
		return fmt.Errorf("subscriber not configured")
	}

	client := s.subscriber.Client()
	if client == nil {
		return fmt.Errorf("rabbitmq client not available")
	}

	return client.SetupQueueWithBinding(queueName, routingKeyPattern)
}

// handleServiceReply handles incoming service reply messages
func (s *Service) handleServiceReply(ctx context.Context, msg *rabbitmq.StandardMessage) error {
	if msg == nil {
		return nil
	}

	s.logger.WithFields(logrus.Fields{
		"tid":       msg.TID,
		"bid":       msg.BID,
		"device_sn": msg.DeviceSN,
		"service":   msg.Service,
		"action":    msg.Action,
	}).Debug("Received service reply")

	// If we have a repository, look up the ServiceCall by TID and update status
	if s.repository != nil {
		return s.handleServiceReplyWithDB(ctx, msg)
	}

	// If no repository, just log and acknowledge
	s.logger.WithField("tid", msg.TID).Debug("No repository available, skipping correlation")
	return nil
}

// handleServiceReplyWithDB handles service reply with database correlation
func (s *Service) handleServiceReplyWithDB(ctx context.Context, msg *rabbitmq.StandardMessage) error {
	if msg.TID == "" {
		s.logger.Warn("Service reply has no TID, skipping")
		return nil
	}

	// Find the service call by TID
	call, err := s.repository.FindByTID(msg.TID)
	if err != nil {
		s.logger.WithError(err).WithField("tid", msg.TID).Warn("Failed to find service call by TID")
		// Don't return error to avoid Nack - the message might be for a different correlation ID
		return nil
	}

	if call == nil {
		s.logger.WithField("tid", msg.TID).Warn("No service call found for TID")
		return nil
	}

	s.logger.WithFields(logrus.Fields{
		"call_id": call.ID,
		"tid":     msg.TID,
		"status":  call.Status,
	}).Debug("Found service call for reply")

	// Check if the call is already in a terminal state
	if call.IsCompleted() {
		// Late response - log for audit but don't reopen terminal state
		s.logger.WithFields(logrus.Fields{
			"call_id":  call.ID,
			"tid":      msg.TID,
			"status":   call.Status,
			"response": string(msg.Data),
		}).Info("Late service reply received for completed call - audited but not processed")

		// Log late response to MessageLog for audit (FR-019)
		if s.msgLogRepo != nil {
			errMsg := "late response after terminal state"
			_ = s.msgLogRepo.CreateLateResponseRecord(ctx, call.DeviceSN, msg.TID, call.BID,
				"iot-downlink", "service.reply", models.MessageDirectionDownlink, errMsg)
		}
		return nil
	}

	// Parse response data
	var responseData map[string]any
	if msg.Data != nil {
		if err := json.Unmarshal(msg.Data, &responseData); err != nil {
			s.logger.WithError(err).Warn("Failed to parse response data")
		}
	}

	// Determine success/failure based on response content
	// A response with error field indicates failure
	var terminalStatus models.MessageStatus
	if errStr, hasError := responseData["error"].(string); hasError && errStr != "" {
		call.MarkFailed(errStr)
		terminalStatus = models.MessageStatusFailed
	} else {
		if err := call.MarkSuccess(responseData); err != nil {
			s.logger.WithError(err).WithField("call_id", call.ID).Error("Failed to mark service call as success")
			return err
		}
		terminalStatus = models.MessageStatusSuccess
	}

	// Update in database
	if err := s.repository.Update(call); err != nil {
		s.logger.WithError(err).WithField("call_id", call.ID).Error("Failed to update service call")
		return err
	}

	// Log terminal state to MessageLog (FR-017)
	if s.msgLogRepo != nil {
		var errMsg string
		if terminalStatus == models.MessageStatusFailed {
			if errStr, ok := responseData["error"].(string); ok {
				errMsg = errStr
			}
		}
		_ = s.msgLogRepo.CreateTerminalStateRecord(ctx, call.DeviceSN, msg.TID, call.BID,
			"iot-downlink", "service.reply", models.MessageDirectionDownlink, terminalStatus, errMsg)
	}

	s.logger.WithFields(logrus.Fields{
		"call_id": call.ID,
		"tid":     msg.TID,
		"status":  call.Status,
	}).Info("Service call status updated from reply")

	return nil
}

// onDispatched is called when a service call is dispatched
func (s *Service) onDispatched(ctx context.Context, call *dispatcher.ServiceCall, result *dispatcher.DispatchResult) error {
	// Route to gateway if enabled
	if s.config.EnableRouting && s.router != nil {
		_, err := s.router.Route(ctx, call, result)
		if err != nil {
			s.logger.WithError(err).WithField("call_id", call.ID).Error("Failed to route to gateway")
			return err
		}
	}
	return nil
}

// onRetry is called when a retry is attempted
func (s *Service) onRetry(ctx context.Context, call *dispatcher.ServiceCall) error {
	s.logger.WithFields(logrus.Fields{
		"call_id":     call.ID,
		"device_sn":   call.DeviceSN,
		"retry_count": call.RetryCount,
	}).Debug("Retrying service call")

	// Log retry attempt to MessageLog
	if s.msgLogRepo != nil {
		errMsg := fmt.Sprintf("retry attempt %d", call.RetryCount)
		_ = s.msgLogRepo.CreateRetryRecord(ctx, call.DeviceSN, call.TID, call.BID,
			"iot-downlink", "service.request", models.MessageDirectionDownlink, call.RetryCount, errMsg)
	}

	_, err := s.handler.Handle(ctx, call)
	return err
}

// onDeadLetter is called when a call is moved to dead letter
func (s *Service) onDeadLetter(entry *retry.DeadLetterEntry) {
	s.logger.WithFields(logrus.Fields{
		"call_id":   entry.Call.ID,
		"device_sn": entry.Call.DeviceSN,
		"error":     entry.Error,
		"retries":   entry.Retries,
	}).Warn("Service call moved to dead letter queue")

	// Log to MessageLog for audit
	if s.msgLogRepo != nil {
		ctx := context.Background()
		_ = s.msgLogRepo.CreateTerminalStateRecord(ctx, entry.Call.DeviceSN, entry.Call.TID, entry.Call.BID,
			"iot-downlink", "service.request", models.MessageDirectionDownlink, models.MessageStatusFailed, entry.Error)
	}
}

// incrementCounter increments the given counter and records a Prometheus metric with the given status.
func (s *Service) incrementCounter(counter *int64, status string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	*counter++

	if s.msgMetrics != nil {
		s.msgMetrics.ProcessedTotal.WithLabelValues("iot-downlink", "", "command", status).Inc()
	}
}

// incrementProcessed increments the processed counter
func (s *Service) incrementProcessed() {
	s.incrementCounter(&s.processedCount, "success")
}

// incrementFailed increments the failed counter
func (s *Service) incrementFailed() {
	s.incrementCounter(&s.failedCount, "error")
}

// GetMetrics returns service metrics
func (s *Service) GetMetrics() (processed, failed int64) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.processedCount, s.failedCount
}

// GetRetryMetrics returns retry metrics
func (s *Service) GetRetryMetrics() (pending, deadLetter int) {
	if s.retryHandler == nil {
		return 0, 0
	}
	return s.retryHandler.GetPendingRetries(), s.retryHandler.GetDeadLetterCount()
}

// GetRouterMetrics returns router metrics
func (s *Service) GetRouterMetrics() (routed, failed int64) {
	if s.router == nil {
		return 0, 0
	}
	return s.router.GetMetrics()
}

// GetRegisteredVendors returns list of registered vendors
func (s *Service) GetRegisteredVendors() []string {
	return s.registry.ListVendors()
}
