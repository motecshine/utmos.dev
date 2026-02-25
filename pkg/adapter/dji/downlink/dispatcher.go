// Package downlink provides DJI message dispatching functionality
package downlink

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	dji "github.com/utmos/utmos/pkg/adapter/dji"
	"github.com/utmos/utmos/pkg/rabbitmq"
)

const (
	// DJI routing key pattern: iot.{vendor}.{action}
	// These are generated from the vendor constant to ensure consistency
	routingKeyPattern = "iot.%s.%s"
)

// Routing key actions
const (
	actionServiceCall  = "service.call"
	actionPropertySet  = "property.set"
	actionConfigUpdate = "config.update"
)

// GetRoutingKey generates a routing key for the given action
func GetRoutingKey(action string) string {
	return fmt.Sprintf(routingKeyPattern, dji.VendorDJI, action)
}

// Routing keys (generated from pattern for backward compatibility)
var (
	RoutingKeyServiceCall  = GetRoutingKey(actionServiceCall)
	RoutingKeyPropertySet  = GetRoutingKey(actionPropertySet)
	RoutingKeyConfigUpdate = GetRoutingKey(actionConfigUpdate)
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

// ServiceCall represents a service call request
type ServiceCall struct {
	ID          string            `json:"id"`
	DeviceSN    string            `json:"device_sn"`
	Vendor      string            `json:"vendor"`
	Method      string            `json:"method"`
	Params      map[string]any    `json:"params"`
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
	CallID     string
	DeviceSN   string
	Vendor     string
	Method     string
	Success    bool
	MessageID  string
	ErrorMsg   string
	RoutingKey string
	SentAt     time.Time
}

// Dispatcher dispatches service calls to DJI devices
type Dispatcher struct {
	vendor    string
	publisher *rabbitmq.Publisher
	logger    *logrus.Entry
}

// NewDispatcher creates a new DJI dispatcher
func NewDispatcher(publisher *rabbitmq.Publisher, logger *logrus.Entry) *Dispatcher {
	if logger == nil {
		logger = logrus.NewEntry(logrus.StandardLogger())
	}
	return &Dispatcher{
		vendor:    dji.VendorDJI,
		publisher: publisher,
		logger:    logger.WithField("dispatcher", dji.VendorDJI),
	}
}

// GetVendor returns the vendor name
func (d *Dispatcher) GetVendor() string {
	return d.vendor
}

// CanDispatch checks if this dispatcher can handle the given call
func (d *Dispatcher) CanDispatch(call *ServiceCall) bool {
	return call.Vendor == dji.VendorDJI
}

// Dispatch dispatches a service call to a DJI device
func (d *Dispatcher) Dispatch(ctx context.Context, call *ServiceCall) (*DispatchResult, error) {
	if call == nil {
		return nil, fmt.Errorf("service call is nil")
	}

	if d.publisher == nil {
		return nil, fmt.Errorf("publisher not initialized")
	}

	// Generate IDs if not set
	if call.TID == "" {
		call.TID = uuid.New().String()
	}
	if call.BID == "" {
		call.BID = uuid.New().String()
	}
	if call.ID == "" {
		call.ID = uuid.New().String()
	}

	// Determine routing key based on call type
	routingKey := d.getRoutingKey(call.CallType)

	// Create standard message
	stdMsg, err := d.createStandardMessage(call)
	if err != nil {
		return d.failResult(call, err, "failed to create standard message")
	}

	// Publish to RabbitMQ
	if err := d.publisher.Publish(ctx, routingKey, stdMsg); err != nil {
		return d.failResult(call, err, "failed to publish message")
	}

	sentAt := time.Now()
	call.SentAt = &sentAt
	call.Status = ServiceCallStatusSent

	d.logger.WithFields(logrus.Fields{
		"device_sn":   call.DeviceSN,
		"method":      call.Method,
		"routing_key": routingKey,
		"tid":         call.TID,
	}).Debug("Dispatched DJI service call")

	return &DispatchResult{
		CallID:     call.ID,
		DeviceSN:   call.DeviceSN,
		Vendor:     call.Vendor,
		Method:     call.Method,
		Success:    true,
		MessageID:  call.TID,
		RoutingKey: routingKey,
		SentAt:     sentAt,
	}, nil
}

// failResult creates a failed DispatchResult for a service call.
func (d *Dispatcher) failResult(call *ServiceCall, err error, msg string) (*DispatchResult, error) {
	return &DispatchResult{
		CallID:   call.ID,
		DeviceSN: call.DeviceSN,
		Vendor:   call.Vendor,
		Method:   call.Method,
		Success:  false,
		ErrorMsg: err.Error(),
		SentAt:   time.Now(),
	}, fmt.Errorf("%s: %w", msg, err)
}

// getRoutingKey returns the routing key for a call type
func (d *Dispatcher) getRoutingKey(callType ServiceCallType) string {
	switch callType {
	case ServiceCallTypeProperty:
		return RoutingKeyPropertySet
	case ServiceCallTypeConfig:
		return RoutingKeyConfigUpdate
	default:
		return RoutingKeyServiceCall
	}
}

// createStandardMessage creates a StandardMessage from a ServiceCall
func (d *Dispatcher) createStandardMessage(call *ServiceCall) (*rabbitmq.StandardMessage, error) {
	// Create data payload in DJI format
	data := map[string]any{
		"method": call.Method,
		"data":   call.Params,
	}

	action := d.getAction(call.CallType)

	stdMsg, err := rabbitmq.NewStandardMessageWithIDs(
		call.TID,
		call.BID,
		"iot-downlink",
		action,
		call.DeviceSN,
		data,
	)
	if err != nil {
		return nil, err
	}

	// Set protocol meta
	stdMsg.ProtocolMeta = &rabbitmq.ProtocolMeta{
		Vendor: dji.VendorDJI,
		Method: call.Method,
	}

	return stdMsg, nil
}

// getAction returns the action string for a call type
func (d *Dispatcher) getAction(callType ServiceCallType) string {
	switch callType {
	case ServiceCallTypeProperty:
		return "property.set"
	case ServiceCallTypeConfig:
		return "config.update"
	default:
		return "service.call"
	}
}

// NewServiceCall creates a new service call
func NewServiceCall(deviceSN, method string, params map[string]any) *ServiceCall {
	return &ServiceCall{
		DeviceSN:   deviceSN,
		Vendor:     dji.VendorDJI,
		Method:     method,
		Params:     params,
		CallType:   ServiceCallTypeCommand,
		Status:     ServiceCallStatusPending,
		CreatedAt:  time.Now(),
		MaxRetries: 3,
	}
}
