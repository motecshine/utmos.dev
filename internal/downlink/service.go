// Package downlink provides the iot-downlink service implementation
package downlink

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/utmos/utmos/internal/downlink/dispatcher"
	"github.com/utmos/utmos/internal/downlink/retry"
	"github.com/utmos/utmos/internal/downlink/router"
	"github.com/utmos/utmos/pkg/adapter"
	"github.com/utmos/utmos/pkg/rabbitmq"
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
	config     *Config
	logger     *logrus.Entry
	registry   *dispatcher.DispatcherRegistry
	handler    *dispatcher.DispatchHandler
	retryHandler *retry.Handler
	router     *router.Router
	publisher  *rabbitmq.Publisher
	subscriber *rabbitmq.Subscriber

	mu       sync.RWMutex
	running  bool
	cancelFn context.CancelFunc

	// Metrics
	processedCount int64
	failedCount    int64
}

// NewService creates a new downlink service
func NewService(config *Config, publisher *rabbitmq.Publisher, logger *logrus.Entry) *Service {
	if config == nil {
		config = DefaultConfig()
	}
	if logger == nil {
		logger = logrus.NewEntry(logrus.StandardLogger())
	}

	serviceLogger := logger.WithField("service", "iot-downlink")

	// Create dispatcher registry
	registry := dispatcher.NewDispatcherRegistry(serviceLogger)

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

	svc := &Service{
		config:       config,
		logger:       serviceLogger,
		registry:     registry,
		handler:      handler,
		retryHandler: retryHandler,
		router:       routerInstance,
		publisher:    publisher,
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

// consumeMessages consumes messages from RabbitMQ
func (s *Service) consumeMessages(ctx context.Context) {
	s.logger.Info("Starting message consumer")

	// This would be implemented based on the actual consumer interface
	// For now, we'll just wait for context cancellation
	<-ctx.Done()
	s.logger.Info("Message consumer stopped")
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
}

// incrementProcessed increments the processed counter
func (s *Service) incrementProcessed() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.processedCount++
}

// incrementFailed increments the failed counter
func (s *Service) incrementFailed() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.failedCount++
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
