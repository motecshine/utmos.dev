package rabbitmq

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

// StandardMessage represents the standard message format for RabbitMQ.
type StandardMessage struct {
	Data      json.RawMessage `json:"data"`      // Business data
	TID       string          `json:"tid"`       // Transaction ID (UUID)
	BID       string          `json:"bid"`       // Business ID (UUID)
	Service   string          `json:"service"`   // Sending service name
	Action    string          `json:"action"`    // Action identifier
	DeviceSN  string          `json:"device_sn"` // Device serial number
	Timestamp int64           `json:"timestamp"` // Millisecond Unix timestamp
}

// MessageHeader represents RabbitMQ message headers.
type MessageHeader struct {
	Traceparent string `json:"traceparent"`      // W3C Trace Context
	Tracestate  string `json:"tracestate"`       // W3C Trace State
	MessageType string `json:"message_type"`     // property, event, service
	Vendor      string `json:"vendor,omitempty"` // Vendor identifier (optional)
}

// NewStandardMessage creates a new standard message with auto-generated IDs and timestamp.
func NewStandardMessage(service, action, deviceSN string, data interface{}) (*StandardMessage, error) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return &StandardMessage{
		TID:       uuid.New().String(),
		BID:       uuid.New().String(),
		Timestamp: time.Now().UnixMilli(),
		Service:   service,
		Action:    action,
		DeviceSN:  deviceSN,
		Data:      dataBytes,
	}, nil
}

// NewStandardMessageWithIDs creates a new standard message with specified IDs.
func NewStandardMessageWithIDs(tid, bid, service, action, deviceSN string, data interface{}) (*StandardMessage, error) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return &StandardMessage{
		TID:       tid,
		BID:       bid,
		Timestamp: time.Now().UnixMilli(),
		Service:   service,
		Action:    action,
		DeviceSN:  deviceSN,
		Data:      dataBytes,
	}, nil
}

// Validate validates the message format.
func (m *StandardMessage) Validate() error {
	if m.TID == "" {
		return errors.New("TID is required")
	}
	if m.BID == "" {
		return errors.New("BID is required")
	}
	if m.Timestamp == 0 {
		return errors.New("Timestamp is required")
	}
	if m.Service == "" {
		return errors.New("Service is required")
	}
	if m.Action == "" {
		return errors.New("Action is required")
	}
	if m.DeviceSN == "" {
		return errors.New("DeviceSN is required")
	}
	return nil
}

// GetData unmarshals the Data field into the provided interface.
func (m *StandardMessage) GetData(v interface{}) error {
	return json.Unmarshal(m.Data, v)
}

// SetData marshals the provided interface into the Data field.
func (m *StandardMessage) SetData(v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	m.Data = data
	return nil
}

// ToBytes serializes the message to JSON bytes.
func (m *StandardMessage) ToBytes() ([]byte, error) {
	return json.Marshal(m)
}

// FromBytes deserializes JSON bytes into a StandardMessage.
func FromBytes(data []byte) (*StandardMessage, error) {
	var msg StandardMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}
