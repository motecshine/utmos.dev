// Package adapter provides the protocol adapter framework for IoT message conversion.
package adapter

import (
	"encoding/json"
	"errors"

	"github.com/utmos/utmos/pkg/rabbitmq"
)

// MessageType represents the type of protocol message.
type MessageType string

// Message type constants.
const (
	MessageTypeProperty MessageType = "property"
	MessageTypeEvent    MessageType = "event"
	MessageTypeService  MessageType = "service"
	MessageTypeStatus   MessageType = "status"
)

// String returns the string representation of the message type.
func (mt MessageType) String() string {
	return string(mt)
}

// ProtocolMessage represents a protocol-specific message parsed from raw bytes.
type ProtocolMessage struct {
	Vendor      string          `json:"vendor"`               // Vendor identifier (dji, tuya, generic)
	Topic       string          `json:"topic"`                // Original MQTT topic
	DeviceSN    string          `json:"device_sn"`            // Device serial number
	GatewaySN   string          `json:"gateway_sn,omitempty"` // Gateway serial number (if applicable)
	MessageType MessageType     `json:"message_type"`         // Message type (property, event, service, status)
	Method      string          `json:"method,omitempty"`     // Protocol method (e.g., thing.property.post)
	TID         string          `json:"tid,omitempty"`        // Transaction ID
	BID         string          `json:"bid,omitempty"`        // Business ID
	Timestamp   int64           `json:"timestamp,omitempty"`  // Message timestamp (milliseconds)
	NeedReply   bool            `json:"need_reply,omitempty"` // Whether response is required
	Data        json.RawMessage `json:"data,omitempty"`       // Business data payload
	Extra       map[string]any  `json:"extra,omitempty"`      // Additional protocol-specific fields
}

// Validate validates the protocol message.
func (pm *ProtocolMessage) Validate() error {
	if pm.Vendor == "" {
		return errors.New("vendor is required")
	}
	if pm.Topic == "" {
		return errors.New("topic is required")
	}
	if pm.DeviceSN == "" {
		return errors.New("device_sn is required")
	}
	if pm.MessageType == "" {
		return errors.New("message_type is required")
	}
	return nil
}

// ProtocolAdapter defines the interface for protocol adapters.
// Each vendor (DJI, Tuya, etc.) implements this interface to handle
// protocol-specific message parsing and conversion.
type ProtocolAdapter interface {
	// GetVendor returns the vendor identifier (e.g., "dji", "tuya").
	GetVendor() string

	// ParseRawMessage parses raw bytes into a protocol-specific message.
	// The topic parameter is the original MQTT topic.
	ParseRawMessage(topic string, payload []byte) (*ProtocolMessage, error)

	// ToStandardMessage converts a protocol message to a standard message.
	ToStandardMessage(pm *ProtocolMessage) (*rabbitmq.StandardMessage, error)

	// FromStandardMessage converts a standard message to a protocol message.
	// Used for downlink message conversion.
	FromStandardMessage(sm *rabbitmq.StandardMessage) (*ProtocolMessage, error)

	// GetRawPayload returns the raw payload bytes for sending to device.
	// Used for downlink message serialization.
	GetRawPayload(pm *ProtocolMessage) ([]byte, error)
}
