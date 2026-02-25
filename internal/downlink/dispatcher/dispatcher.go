// Package dispatcher provides message dispatching functionality for iot-downlink
package dispatcher

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/utmos/utmos/pkg/adapter"
	"github.com/utmos/utmos/pkg/rabbitmq"
	"github.com/utmos/utmos/pkg/registry"
)

// ServiceCallType represents the type of service call
type ServiceCallType string

const (
	// ServiceCallTypeCommand is the command service call type.
	ServiceCallTypeCommand ServiceCallType = "command"
	// ServiceCallTypeProperty is the property service call type.
	ServiceCallTypeProperty ServiceCallType = "property"
	// ServiceCallTypeConfig is the config service call type.
	ServiceCallTypeConfig ServiceCallType = "config"
)

// ServiceCallStatus represents the status of a service call
type ServiceCallStatus string

const (
	// ServiceCallStatusPending indicates the service call is pending.
	ServiceCallStatusPending ServiceCallStatus = "pending"
	// ServiceCallStatusSent indicates the service call has been sent.
	ServiceCallStatusSent ServiceCallStatus = "sent"
	// ServiceCallStatusSuccess indicates the service call succeeded.
	ServiceCallStatusSuccess ServiceCallStatus = "success"
	// ServiceCallStatusFailed indicates the service call failed.
	ServiceCallStatusFailed ServiceCallStatus = "failed"
	// ServiceCallStatusTimeout indicates the service call timed out.
	ServiceCallStatusTimeout ServiceCallStatus = "timeout"
	// ServiceCallStatusRetrying indicates the service call is being retried.
	ServiceCallStatusRetrying ServiceCallStatus = "retrying"
)

// ServiceCallParams represents service call parameters as raw JSON.
// Using json.RawMessage allows for type-safe handling while deferring parsing.
type ServiceCallParams = json.RawMessage

// ServiceCall represents a service call request
type ServiceCall struct {
	ID          string            `json:"id"`
	DeviceSN    string            `json:"device_sn"`
	Vendor      string            `json:"vendor"`
	Method      string            `json:"method"`
	Params      ServiceCallParams `json:"params"`
	CallType    ServiceCallType   `json:"call_type"`
	Status      ServiceCallStatus `json:"status"`
	TID         string            `json:"tid"`
	BID         string            `json:"bid"`
	CreatedAt   time.Time         `json:"created_at"`
	SentAt      *time.Time        `json:"sent_at,omitempty"`
	CompletedAt *time.Time        `json:"completed_at,omitempty"`
	RetryCount  int               `json:"retry_count"`
	MaxRetries  int               `json:"max_retries"`
	Error       string            `json:"error,omitempty"`
}

// DispatchResult represents the result of a dispatch operation
type DispatchResult struct {
	Success   bool
	MessageID string
	Error     error
	SentAt    time.Time
}

// Dispatcher defines the interface for message dispatchers
type Dispatcher interface {
	// Dispatch dispatches a service call to the appropriate destination
	Dispatch(ctx context.Context, call *ServiceCall) (*DispatchResult, error)

	// GetVendor returns the vendor this dispatcher handles
	GetVendor() string

	// CanDispatch checks if this dispatcher can handle the given call
	CanDispatch(call *ServiceCall) bool
}

// AdapterDispatcher wraps an adapter.DownlinkDispatcher to implement the internal Dispatcher interface
type AdapterDispatcher struct {
	adapter adapter.DownlinkDispatcher
}

// NewAdapterDispatcher creates a new adapter dispatcher wrapper
func NewAdapterDispatcher(a adapter.DownlinkDispatcher) *AdapterDispatcher {
	return &AdapterDispatcher{adapter: a}
}

// GetVendor returns the vendor name
func (d *AdapterDispatcher) GetVendor() string {
	return d.adapter.GetVendor()
}

// CanDispatch checks if this dispatcher can handle the given call
func (d *AdapterDispatcher) CanDispatch(call *ServiceCall) bool {
	return d.adapter.CanDispatch(call.Vendor)
}

// Dispatch dispatches a service call
func (d *AdapterDispatcher) Dispatch(ctx context.Context, call *ServiceCall) (*DispatchResult, error) {
	// Parse params from JSON
	var params map[string]any
	if len(call.Params) > 0 {
		if err := json.Unmarshal(call.Params, &params); err != nil {
			return &DispatchResult{
				Success: false,
				Error:   fmt.Errorf("failed to parse params: %w", err),
				SentAt:  time.Now(),
			}, err
		}
	}

	// Convert to adapter format
	adapterCall := &adapter.ServiceCall{
		ID:       call.ID,
		DeviceSN: call.DeviceSN,
		Vendor:   call.Vendor,
		Method:   call.Method,
		Params:   params,
	}

	result, err := d.adapter.Dispatch(ctx, adapterCall)
	if err != nil {
		return &DispatchResult{
			Success: false,
			Error:   err,
			SentAt:  time.Now(),
		}, err
	}

	// Update call with result
	if result.Success {
		now := time.Now()
		call.SentAt = &now
		call.Status = ServiceCallStatusSent
	}

	return &DispatchResult{
		Success:   result.Success,
		MessageID: result.CallID,
		SentAt:    time.Now(),
	}, nil
}

// Registry manages registered dispatchers using the generic registry
type Registry struct {
	*registry.Registry[Dispatcher]
}

// NewRegistry creates a new dispatcher registry
func NewRegistry(logger *logrus.Entry) *Registry {
	return &Registry{
		Registry: registry.New[Dispatcher]("dispatcher-registry", logger),
	}
}

// GetForCall returns a dispatcher that can handle the given call
func (r *Registry) GetForCall(call *ServiceCall) (Dispatcher, bool) {
	// Use GetOrFind: first try vendor lookup, then predicate
	return r.GetOrFind(call.Vendor, func(d Dispatcher) bool {
		return d.CanDispatch(call)
	})
}

// DispatchHandler handles service call dispatching
type DispatchHandler struct {
	registry     *Registry
	logger       *logrus.Entry
	onDispatched func(ctx context.Context, call *ServiceCall, result *DispatchResult) error
}

// NewDispatchHandler creates a new dispatch handler
func NewDispatchHandler(registry *Registry, logger *logrus.Entry) *DispatchHandler {
	if logger == nil {
		logger = logrus.NewEntry(logrus.StandardLogger())
	}
	return &DispatchHandler{
		registry: registry,
		logger:   logger.WithField("component", "dispatch-handler"),
	}
}

// SetOnDispatched sets the callback for dispatched calls
func (h *DispatchHandler) SetOnDispatched(callback func(ctx context.Context, call *ServiceCall, result *DispatchResult) error) {
	h.onDispatched = callback
}

// Handle dispatches a service call
func (h *DispatchHandler) Handle(ctx context.Context, call *ServiceCall) (*DispatchResult, error) {
	if call == nil {
		return nil, fmt.Errorf("service call is nil")
	}

	tr := otel.Tracer("iot-downlink")
	ctx, span := tr.Start(ctx, "downlink.dispatch",
		trace.WithAttributes(
			attribute.String("device_sn", call.DeviceSN),
			attribute.String("vendor", call.Vendor),
			attribute.String("method", call.Method),
		),
	)
	defer span.End()

	// Find appropriate dispatcher
	dispatcher, found := h.registry.GetForCall(call)
	if !found {
		h.logger.WithFields(logrus.Fields{
			"device_sn": call.DeviceSN,
			"vendor":    call.Vendor,
			"method":    call.Method,
		}).Warn("No dispatcher found for service call")
		err := fmt.Errorf("no dispatcher found for vendor: %s", call.Vendor)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// Dispatch the call
	result, err := dispatcher.Dispatch(ctx, call)
	if err != nil {
		h.logger.WithError(err).WithFields(logrus.Fields{
			"device_sn": call.DeviceSN,
			"vendor":    call.Vendor,
			"method":    call.Method,
		}).Error("Failed to dispatch service call")
		wrappedErr := fmt.Errorf("failed to dispatch: %w", err)
		span.RecordError(wrappedErr)
		span.SetStatus(codes.Error, wrappedErr.Error())
		return result, wrappedErr
	}

	h.logger.WithFields(logrus.Fields{
		"device_sn":  call.DeviceSN,
		"vendor":     call.Vendor,
		"method":     call.Method,
		"message_id": result.MessageID,
	}).Debug("Service call dispatched")

	// Call the dispatched callback if set
	if h.onDispatched != nil {
		if err := h.onDispatched(ctx, call, result); err != nil {
			h.logger.WithError(err).WithField("call_id", call.ID).Error("Failed to handle dispatched callback")
		}
	}

	return result, nil
}

// RegisterDispatcher registers a dispatcher
func (h *DispatchHandler) RegisterDispatcher(dispatcher Dispatcher) {
	h.registry.Register(dispatcher)
}

// RegisterAdapterDispatcher registers an adapter.DownlinkDispatcher by wrapping it
func (h *DispatchHandler) RegisterAdapterDispatcher(a adapter.DownlinkDispatcher) {
	h.registry.Register(NewAdapterDispatcher(a))
}

// UnregisterDispatcher unregisters a dispatcher
func (h *DispatchHandler) UnregisterDispatcher(vendor string) {
	h.registry.Unregister(vendor)
}

// BaseDispatcher provides common functionality for dispatchers
type BaseDispatcher struct {
	vendor    string
	publisher *rabbitmq.Publisher
	logger    *logrus.Entry
}

// NewBaseDispatcher creates a new base dispatcher
func NewBaseDispatcher(vendor string, publisher *rabbitmq.Publisher, logger *logrus.Entry) *BaseDispatcher {
	if logger == nil {
		logger = logrus.NewEntry(logrus.StandardLogger())
	}
	return &BaseDispatcher{
		vendor:    vendor,
		publisher: publisher,
		logger:    logger.WithField("dispatcher", vendor),
	}
}

// GetVendor returns the vendor name
func (d *BaseDispatcher) GetVendor() string {
	return d.vendor
}

// Publisher returns the RabbitMQ publisher
func (d *BaseDispatcher) Publisher() *rabbitmq.Publisher {
	return d.publisher
}

// Logger returns the dispatcher's logger
func (d *BaseDispatcher) Logger() *logrus.Entry {
	return d.logger
}

// NewServiceCall creates a new service call
func NewServiceCall(deviceSN, vendor, method string, params json.RawMessage) *ServiceCall {
	return &ServiceCall{
		DeviceSN:   deviceSN,
		Vendor:     vendor,
		Method:     method,
		Params:     params,
		CallType:   ServiceCallTypeCommand,
		Status:     ServiceCallStatusPending,
		CreatedAt:  time.Now(),
		MaxRetries: 3,
	}
}
