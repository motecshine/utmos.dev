// Package router provides message routing functionality for iot-downlink
package router

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/utmos/utmos/internal/downlink/dispatcher"
	"github.com/utmos/utmos/pkg/rabbitmq"
)

const (
	// RoutingKeyGatewayDownlink is the routing key for downlink messages to gateway
	RoutingKeyGatewayDownlink = "iot.gateway.downlink"

	// RoutingKeyGatewayCommand is the routing key for command messages
	RoutingKeyGatewayCommand = "iot.gateway.command"

	// RoutingKeyGatewayProperty is the routing key for property set messages
	RoutingKeyGatewayProperty = "iot.gateway.property"
)

// Config holds router configuration
type Config struct {
	// DefaultRoutingKey is the default routing key for messages
	DefaultRoutingKey string
	// EnableMetrics enables routing metrics
	EnableMetrics bool
}

// DefaultConfig returns default router configuration
func DefaultConfig() *Config {
	return &Config{
		DefaultRoutingKey: RoutingKeyGatewayDownlink,
		EnableMetrics:     true,
	}
}

// Router routes dispatched messages to the gateway service
type Router struct {
	publisher *rabbitmq.Publisher
	config    *Config
	logger    *logrus.Entry
	mu        sync.RWMutex

	// Metrics
	routedCount  int64
	failedCount  int64
}

// NewRouter creates a new router
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
		logger:    logger.WithField("component", "downlink-router"),
	}
}

// RouteResult represents the result of a routing operation
type RouteResult struct {
	Success    bool
	RoutingKey string
	Error      error
}

// Route routes a service call result to the gateway
func (r *Router) Route(ctx context.Context, call *dispatcher.ServiceCall, result *dispatcher.DispatchResult) (*RouteResult, error) {
	if call == nil {
		return nil, fmt.Errorf("service call is nil")
	}

	if r.publisher == nil {
		return nil, fmt.Errorf("publisher not initialized")
	}

	// Determine routing key based on call type
	routingKey := r.getRoutingKey(call)

	// Create gateway message
	msg, err := r.createGatewayMessage(call, result)
	if err != nil {
		r.incrementFailed()
		return &RouteResult{
			Success:    false,
			RoutingKey: routingKey,
			Error:      err,
		}, fmt.Errorf("failed to create gateway message: %w", err)
	}

	// Publish to gateway
	if err := r.publisher.Publish(ctx, routingKey, msg); err != nil {
		r.incrementFailed()
		return &RouteResult{
			Success:    false,
			RoutingKey: routingKey,
			Error:      err,
		}, fmt.Errorf("failed to publish to gateway: %w", err)
	}

	r.incrementRouted()

	r.logger.WithFields(logrus.Fields{
		"device_sn":   call.DeviceSN,
		"method":      call.Method,
		"routing_key": routingKey,
		"tid":         call.TID,
	}).Debug("Routed message to gateway")

	return &RouteResult{
		Success:    true,
		RoutingKey: routingKey,
	}, nil
}

// getRoutingKey determines the routing key for a service call
func (r *Router) getRoutingKey(call *dispatcher.ServiceCall) string {
	switch call.CallType {
	case dispatcher.ServiceCallTypeCommand:
		return RoutingKeyGatewayCommand
	case dispatcher.ServiceCallTypeProperty:
		return RoutingKeyGatewayProperty
	default:
		return r.config.DefaultRoutingKey
	}
}

// gatewayPayload is the typed payload sent in gateway messages.
type gatewayPayload struct {
	Method         string                 `json:"method"`
	Params         any                    `json:"params"`
	DispatchResult *gatewayDispatchResult `json:"dispatch_result,omitempty"`
}

// gatewayDispatchResult is the typed dispatch result embedded in gateway payloads.
type gatewayDispatchResult struct {
	Success   bool      `json:"success"`
	MessageID string    `json:"message_id"`
	SentAt    time.Time `json:"sent_at"`
	Error     string    `json:"error,omitempty"`
}

// createGatewayMessage creates a StandardMessage for the gateway
func (r *Router) createGatewayMessage(call *dispatcher.ServiceCall, result *dispatcher.DispatchResult) (*rabbitmq.StandardMessage, error) {
	// Create typed payload
	payload := &gatewayPayload{
		Method: call.Method,
		Params: call.Params,
	}

	if result != nil {
		dr := &gatewayDispatchResult{
			Success:   result.Success,
			MessageID: result.MessageID,
			SentAt:    result.SentAt,
		}
		if result.Error != nil {
			dr.Error = result.Error.Error()
		}
		payload.DispatchResult = dr
	}

	// Create standard message
	msg, err := rabbitmq.NewStandardMessageWithIDs(
		call.TID,
		call.BID,
		"iot-downlink",
		r.getAction(call.CallType),
		call.DeviceSN,
		payload,
	)
	if err != nil {
		return nil, err
	}

	// Set protocol meta
	msg.ProtocolMeta = &rabbitmq.ProtocolMeta{
		Vendor: call.Vendor,
		Method: call.Method,
	}

	return msg, nil
}

// getAction returns the action string for a call type
func (r *Router) getAction(callType dispatcher.ServiceCallType) string {
	switch callType {
	case dispatcher.ServiceCallTypeCommand:
		return "command.send"
	case dispatcher.ServiceCallTypeProperty:
		return "property.set"
	case dispatcher.ServiceCallTypeConfig:
		return "config.update"
	default:
		return "downlink.send"
	}
}

// incrementRouted increments the routed counter
func (r *Router) incrementRouted() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.routedCount++
}

// incrementFailed increments the failed counter
func (r *Router) incrementFailed() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.failedCount++
}

// GetMetrics returns routing metrics
func (r *Router) GetMetrics() (routed, failed int64) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.routedCount, r.failedCount
}

// ResetMetrics resets routing metrics
func (r *Router) ResetMetrics() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.routedCount = 0
	r.failedCount = 0
}

// BatchRouter handles batch routing of multiple messages
type BatchRouter struct {
	router *Router
	logger *logrus.Entry
}

// NewBatchRouter creates a new batch router
func NewBatchRouter(router *Router, logger *logrus.Entry) *BatchRouter {
	if logger == nil {
		logger = logrus.NewEntry(logrus.StandardLogger())
	}
	return &BatchRouter{
		router: router,
		logger: logger.WithField("component", "batch-router"),
	}
}

// BatchRouteResult represents the result of a batch routing operation
type BatchRouteResult struct {
	Total     int
	Succeeded int
	Failed    int
	Results   []*RouteResult
}

// RouteBatch routes multiple service calls
func (br *BatchRouter) RouteBatch(ctx context.Context, calls []*dispatcher.ServiceCall, results map[string]*dispatcher.DispatchResult) *BatchRouteResult {
	batchResult := &BatchRouteResult{
		Total:   len(calls),
		Results: make([]*RouteResult, 0, len(calls)),
	}

	for _, call := range calls {
		var dispatchResult *dispatcher.DispatchResult
		if results != nil {
			dispatchResult = results[call.ID]
		}

		result, err := br.router.Route(ctx, call, dispatchResult)
		if err != nil {
			br.logger.WithError(err).WithField("call_id", call.ID).Error("Failed to route call")
			batchResult.Failed++
			if result == nil {
				result = &RouteResult{
					Success: false,
					Error:   err,
				}
			}
		} else {
			batchResult.Succeeded++
		}
		batchResult.Results = append(batchResult.Results, result)
	}

	br.logger.WithFields(logrus.Fields{
		"total":     batchResult.Total,
		"succeeded": batchResult.Succeeded,
		"failed":    batchResult.Failed,
	}).Debug("Batch routing completed")

	return batchResult
}
