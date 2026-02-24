// Package uplink provides the IoT uplink message processing service
package uplink

import (
	"context"
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/utmos/utmos/internal/uplink/processor"
	"github.com/utmos/utmos/internal/uplink/router"
	"github.com/utmos/utmos/internal/uplink/storage"
	"github.com/utmos/utmos/pkg/adapter"
	"github.com/utmos/utmos/pkg/rabbitmq"
)

// ServiceConfig holds configuration for the uplink service
type ServiceConfig struct {
	// Queue configuration
	QueueName   string
	RoutingKeys []string

	// Storage configuration
	Influx *storage.Config

	// Router configuration
	Router *router.Config

	// Processing configuration
	EnableStorage bool
	EnableRouting bool
}

// DefaultServiceConfig returns default service configuration
// Note: RoutingKeys should be configured externally to support multiple vendors
// The default uses wildcard patterns that can be overridden via configuration
func DefaultServiceConfig() *ServiceConfig {
	return &ServiceConfig{
		QueueName: "iot.uplink.messages",
		RoutingKeys: []string{
			"iot.*.#", // Wildcard pattern for all vendors - configure specific vendors externally
		},
		Influx:        storage.DefaultConfig(),
		Router:        router.DefaultConfig(),
		EnableStorage: true,
		EnableRouting: true,
	}
}

// Service is the main uplink service
type Service struct {
	config *ServiceConfig
	logger *logrus.Entry

	// Components
	registry   *processor.Registry
	handler    *processor.MessageHandler
	storage    *storage.Storage
	router     *router.Router
	subscriber *rabbitmq.Subscriber
	publisher  *rabbitmq.Publisher

	// State
	running bool
	mu      sync.RWMutex
	cancel  context.CancelFunc
}

// NewService creates a new uplink service
func NewService(
	config *ServiceConfig,
	subscriber *rabbitmq.Subscriber,
	publisher *rabbitmq.Publisher,
	logger *logrus.Entry,
) *Service {
	if config == nil {
		config = DefaultServiceConfig()
	}
	if logger == nil {
		logger = logrus.NewEntry(logrus.StandardLogger())
	}

	svcLogger := logger.WithField("component", "uplink-service")

	// Create processor registry and handler
	registry := processor.NewRegistry(svcLogger)
	handler := processor.NewMessageHandler(registry, svcLogger)

	// Note: Processors should be registered by the caller using RegisterProcessor()
	// This allows for vendor-agnostic service initialization

	// Create storage
	var influxStorage *storage.Storage
	if config.EnableStorage {
		influxStorage = storage.NewStorage(config.Influx, svcLogger)
	}

	// Create router
	var msgRouter *router.Router
	if config.EnableRouting {
		msgRouter = router.NewRouter(publisher, config.Router, svcLogger)
	}

	svc := &Service{
		config:     config,
		logger:     svcLogger,
		registry:   registry,
		handler:    handler,
		storage:    influxStorage,
		router:     msgRouter,
		subscriber: subscriber,
		publisher:  publisher,
	}

	// Set up the processing pipeline
	handler.SetOnProcessed(svc.onMessageProcessed)

	return svc
}

// Start starts the uplink service
func (s *Service) Start(ctx context.Context) error {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return fmt.Errorf("service already running")
	}
	s.running = true
	s.mu.Unlock()

	_, s.cancel = context.WithCancel(ctx)

	s.logger.Info("Starting uplink service")

	// Start router if enabled
	if s.router != nil {
		if err := s.router.Start(); err != nil {
			s.logger.WithError(err).Warn("Failed to start router")
		}
	}

	// Subscribe to message queues
	if s.subscriber != nil {
		if err := s.subscribeToQueues(); err != nil {
			s.mu.Lock()
			s.running = false
			s.mu.Unlock()
			return fmt.Errorf("failed to subscribe to queues: %w", err)
		}
	}

	s.logger.Info("Uplink service started")
	return nil
}

// Stop stops the uplink service
func (s *Service) Stop() error {
	s.mu.Lock()
	if !s.running {
		s.mu.Unlock()
		return nil
	}
	s.running = false
	s.mu.Unlock()

	s.logger.Info("Stopping uplink service")

	// Cancel context
	if s.cancel != nil {
		s.cancel()
	}

	// Stop router
	if s.router != nil {
		if err := s.router.Stop(); err != nil {
			s.logger.WithError(err).Warn("Failed to stop router")
		}
	}

	// Close storage
	if s.storage != nil {
		if err := s.storage.Close(); err != nil {
			s.logger.WithError(err).Warn("Failed to close storage")
		}
	}

	// Unsubscribe from queues
	if s.subscriber != nil {
		if err := s.subscriber.Unsubscribe(s.config.QueueName); err != nil {
			s.logger.WithError(err).Warn("Failed to unsubscribe from queue")
		}
	}

	s.logger.Info("Uplink service stopped")
	return nil
}

// subscribeToQueues subscribes to the configured message queues
func (s *Service) subscribeToQueues() error {
	return s.subscriber.Subscribe(s.config.QueueName, func(ctx context.Context, msg *rabbitmq.StandardMessage) error {
		return s.handler.Handle(ctx, msg)
	})
}

// tryOperation attempts an operation and logs/wraps errors with a descriptive name
func (s *Service) tryOperation(name string, deviceSN string, fn func() error) error {
	if err := fn(); err != nil {
		s.logger.WithError(err).WithField("device_sn", deviceSN).Errorf("Failed to %s message", name)
		return fmt.Errorf("%s error: %w", name, err)
	}
	return nil
}

// processStep defines a conditional post-processing operation.
type processStep struct {
	name    string
	enabled bool
	run     func() error
}

func (s *Service) onMessageProcessed(ctx context.Context, processed *adapter.ProcessedMessage) error {
	steps := []processStep{
		{"store", s.storage != nil && s.config.EnableStorage, func() error {
			return s.storage.WriteProcessedMessage(ctx, processed)
		}},
		{"route", s.router != nil && s.config.EnableRouting, func() error {
			return s.router.Route(ctx, processed)
		}},
	}

	var errs []error
	for _, step := range steps {
		if step.enabled {
			if err := s.tryOperation(step.name, processed.DeviceSN, step.run); err != nil {
				errs = append(errs, err)
			}
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("processing errors: %v", errs)
	}

	return nil
}

// IsRunning returns whether the service is running
func (s *Service) IsRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.running
}

// GetRegistry returns the processor registry
func (s *Service) GetRegistry() *processor.Registry {
	return s.registry
}

// GetHandler returns the message handler
func (s *Service) GetHandler() *processor.MessageHandler {
	return s.handler
}

// GetStorage returns the InfluxDB storage
func (s *Service) GetStorage() *storage.Storage {
	return s.storage
}

// GetRouter returns the message router
func (s *Service) GetRouter() *router.Router {
	return s.router
}

// RegisterProcessor registers a new processor
func (s *Service) RegisterProcessor(p adapter.UplinkProcessor) {
	s.registry.Register(p)
}

// ProcessMessage manually processes a message (for testing)
func (s *Service) ProcessMessage(ctx context.Context, msg *rabbitmq.StandardMessage) error {
	return s.handler.Handle(ctx, msg)
}

// Stats holds service statistics
type Stats struct {
	Running           bool
	RegisteredVendors []string
	StorageEnabled    bool
	RoutingEnabled    bool
}

// GetStats returns service statistics
func (s *Service) GetStats() *Stats {
	return &Stats{
		Running:           s.IsRunning(),
		RegisteredVendors: s.registry.ListVendors(),
		StorageEnabled:    s.config.EnableStorage,
		RoutingEnabled:    s.config.EnableRouting,
	}
}
