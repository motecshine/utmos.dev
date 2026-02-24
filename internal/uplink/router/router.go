// Package router provides message routing functionality for iot-uplink
package router

import (
	"context"
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/utmos/utmos/pkg/adapter"
	"github.com/utmos/utmos/pkg/rabbitmq"
)

// RoutingKey constants for different destinations
const (
	RoutingKeyWSProperty = "iot.ws.property"
	RoutingKeyWSEvent    = "iot.ws.event"
	RoutingKeyWSStatus   = "iot.ws.status"
	RoutingKeyAPIProperty = "iot.api.property"
	RoutingKeyAPIEvent    = "iot.api.event"
)

// Config holds router configuration
type Config struct {
	Exchange         string
	EnableWSRouting  bool
	EnableAPIRouting bool
}

// DefaultConfig returns default router configuration
func DefaultConfig() *Config {
	return &Config{
		Exchange:         "iot.topic",
		EnableWSRouting:  true,
		EnableAPIRouting: true,
	}
}

// Router routes processed messages to other services
type Router struct {
	publisher *rabbitmq.Publisher
	config    *Config
	logger    *logrus.Entry
	mu        sync.RWMutex
	running   bool
}

// NewRouter creates a new message router
func NewRouter(publisher *rabbitmq.Publisher, config *Config, logger *logrus.Entry) *Router {
	if config == nil {
		config = DefaultConfig()
	}
	if logger == nil {
		logger = logrus.NewEntry(logrus.StandardLogger())
	}

	return &Router{
		publisher: publisher,
		config:    config,
		logger:    logger.WithField("component", "message-router"),
	}
}

// Start starts the router
func (r *Router) Start() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.running {
		return fmt.Errorf("router already running")
	}

	r.running = true
	r.logger.Info("Message router started")
	return nil
}

// Stop stops the router
func (r *Router) Stop() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.running {
		return nil
	}

	r.running = false
	r.logger.Info("Message router stopped")
	return nil
}

// Route routes a processed message to appropriate destinations
func (r *Router) Route(ctx context.Context, msg *adapter.ProcessedMessage) error {
	if msg == nil {
		return fmt.Errorf("message is nil")
	}

	if r.publisher == nil {
		return fmt.Errorf("publisher not initialized")
	}

	var errs []error

	// Route to WebSocket service
	if r.config.EnableWSRouting {
		if err := r.routeToWS(ctx, msg); err != nil {
			errs = append(errs, fmt.Errorf("ws routing failed: %w", err))
		}
	}

	// Route to API service
	if r.config.EnableAPIRouting {
		if err := r.routeToAPI(ctx, msg); err != nil {
			errs = append(errs, fmt.Errorf("api routing failed: %w", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("routing errors: %v", errs)
	}

	return nil
}

// routeToWS routes message to WebSocket service
func (r *Router) routeToWS(ctx context.Context, msg *adapter.ProcessedMessage) error {
	routingKey := r.getWSRoutingKey(msg.MessageType)

	stdMsg, err := r.createStandardMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to create standard message: %w", err)
	}

	if err := r.publisher.Publish(ctx, routingKey, stdMsg); err != nil {
		return fmt.Errorf("failed to publish to WS: %w", err)
	}

	r.logger.WithFields(logrus.Fields{
		"device_sn":   msg.DeviceSN,
		"routing_key": routingKey,
	}).Debug("Routed message to WS")

	return nil
}

// routeToAPI routes message to API service
func (r *Router) routeToAPI(ctx context.Context, msg *adapter.ProcessedMessage) error {
	// Only route property and event messages to API
	if msg.MessageType != adapter.MessageTypeProperty && msg.MessageType != adapter.MessageTypeEvent {
		return nil
	}

	routingKey := r.getAPIRoutingKey(msg.MessageType)

	stdMsg, err := r.createStandardMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to create standard message: %w", err)
	}

	if err := r.publisher.Publish(ctx, routingKey, stdMsg); err != nil {
		return fmt.Errorf("failed to publish to API: %w", err)
	}

	r.logger.WithFields(logrus.Fields{
		"device_sn":   msg.DeviceSN,
		"routing_key": routingKey,
	}).Debug("Routed message to API")

	return nil
}

// getWSRoutingKey returns the routing key for WebSocket service
func (r *Router) getWSRoutingKey(msgType adapter.MessageType) string {
	switch msgType {
	case adapter.MessageTypeProperty:
		return RoutingKeyWSProperty
	case adapter.MessageTypeEvent:
		return RoutingKeyWSEvent
	case adapter.MessageTypeStatus:
		return RoutingKeyWSStatus
	default:
		return RoutingKeyWSProperty
	}
}

// getAPIRoutingKey returns the routing key for API service
func (r *Router) getAPIRoutingKey(msgType adapter.MessageType) string {
	switch msgType {
	case adapter.MessageTypeProperty:
		return RoutingKeyAPIProperty
	case adapter.MessageTypeEvent:
		return RoutingKeyAPIEvent
	default:
		return RoutingKeyAPIProperty
	}
}

// createStandardMessage creates a StandardMessage from ProcessedMessage
func (r *Router) createStandardMessage(msg *adapter.ProcessedMessage) (*rabbitmq.StandardMessage, error) {
	// Prepare data payload
	data := map[string]any{
		"properties": msg.Properties,
		"events":     msg.Events,
	}

	action := r.getAction(msg.MessageType)

	stdMsg, err := rabbitmq.NewStandardMessage("iot-uplink", action, msg.DeviceSN, data)
	if err != nil {
		return nil, err
	}

	// Preserve original TID/BID if available
	if msg.Original != nil {
		stdMsg.TID = msg.Original.TID
		stdMsg.BID = msg.Original.BID
	}

	// Set protocol meta
	stdMsg.ProtocolMeta = &rabbitmq.ProtocolMeta{
		Vendor: msg.Vendor,
	}

	return stdMsg, nil
}

// getAction returns the action string for a message type
func (r *Router) getAction(msgType adapter.MessageType) string {
	switch msgType {
	case adapter.MessageTypeProperty:
		return "property.processed"
	case adapter.MessageTypeEvent:
		return "event.processed"
	case adapter.MessageTypeService:
		return "service.processed"
	case adapter.MessageTypeStatus:
		return "status.processed"
	default:
		return "message.processed"
	}
}

// IsRunning returns whether the router is running
func (r *Router) IsRunning() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.running
}

// RouteCallback is a callback function for routing
type RouteCallback func(ctx context.Context, msg *adapter.ProcessedMessage) error

// MultiRouter routes messages to multiple destinations
type MultiRouter struct {
	routers []RouteCallback
	logger  *logrus.Entry
	mu      sync.RWMutex
}

// NewMultiRouter creates a new multi-router
func NewMultiRouter(logger *logrus.Entry) *MultiRouter {
	if logger == nil {
		logger = logrus.NewEntry(logrus.StandardLogger())
	}
	return &MultiRouter{
		routers: make([]RouteCallback, 0),
		logger:  logger.WithField("component", "multi-router"),
	}
}

// AddRouter adds a router callback
func (m *MultiRouter) AddRouter(router RouteCallback) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.routers = append(m.routers, router)
}

// Route routes a message to all registered routers
func (m *MultiRouter) Route(ctx context.Context, msg *adapter.ProcessedMessage) error {
	m.mu.RLock()
	routers := make([]RouteCallback, len(m.routers))
	copy(routers, m.routers)
	m.mu.RUnlock()

	var errs []error
	for _, router := range routers {
		if err := router(ctx, msg); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("routing errors: %v", errs)
	}

	return nil
}
