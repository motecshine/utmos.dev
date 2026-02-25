// Package adapter provides the protocol adapter framework for IoT message conversion.
package adapter

import (
	"context"

	"github.com/utmos/utmos/pkg/rabbitmq"
)

// ProcessedMessage represents a message after processing by an uplink processor.
// This is the public interface that adapters should use.
type ProcessedMessage struct {
	Original    *rabbitmq.StandardMessage `json:"-"`
	MessageType MessageType               `json:"message_type"`
	DeviceSN    string                    `json:"device_sn"`
	Vendor      string                    `json:"vendor"`
	Properties  map[string]any            `json:"properties,omitempty"`
	Events      []Event                   `json:"events,omitempty"`
	Timestamp   int64                     `json:"timestamp"`
}

// Event represents a device event in a processed message.
type Event struct {
	Name   string         `json:"name"`
	Params map[string]any `json:"params"`
	Output map[string]any `json:"output,omitempty"`
}

// UplinkProcessor defines the interface for uplink message processors.
// Vendor-specific adapters (DJI, Tuya, etc.) implement this interface.
type UplinkProcessor interface {
	// GetVendor returns the vendor identifier (e.g., "dji", "tuya").
	GetVendor() string

	// CanProcess checks if this processor can handle the given message.
	CanProcess(msg *rabbitmq.StandardMessage) bool

	// Process processes a standard message and returns processed data.
	Process(ctx context.Context, msg *rabbitmq.StandardMessage) (*ProcessedMessage, error)
}

// DownlinkDispatcher defines the interface for downlink message dispatchers.
// Vendor-specific adapters implement this interface to handle service calls.
type DownlinkDispatcher interface {
	// GetVendor returns the vendor identifier.
	GetVendor() string

	// CanDispatch checks if this dispatcher can handle the given vendor.
	CanDispatch(vendor string) bool

	// Dispatch dispatches a service call to the device.
	Dispatch(ctx context.Context, call *ServiceCall) (*DispatchResult, error)
}

// ServiceCall represents a service call request to a device.
type ServiceCall struct {
	ID       string         `json:"id"`
	DeviceSN string         `json:"device_sn"`
	Vendor   string         `json:"vendor"`
	Method   string         `json:"method"`
	Params   map[string]any `json:"params,omitempty"`
	Timeout  int64          `json:"timeout,omitempty"` // Timeout in milliseconds
}

// DispatchResult represents the result of a service call dispatch.
type DispatchResult struct {
	CallID     string         `json:"call_id"`
	DeviceSN   string         `json:"device_sn"`
	Vendor     string         `json:"vendor"`
	Method     string         `json:"method"`
	Success    bool           `json:"success"`
	Error      string         `json:"error,omitempty"`
	Response   map[string]any `json:"response,omitempty"`
	RoutingKey string         `json:"routing_key,omitempty"`
}
