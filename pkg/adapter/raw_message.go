package adapter

import (
	"errors"
	"time"
)

// RawMessage represents a raw message received from the MQTT broker
// via iot-gateway before protocol-specific parsing.
type RawMessage struct {
	Vendor    string            `json:"vendor"`    // Vendor identifier (dji, tuya, generic)
	Topic     string            `json:"topic"`     // Original MQTT topic
	Payload   []byte            `json:"payload"`   // Raw message payload
	QoS       int               `json:"qos"`       // MQTT QoS level
	Timestamp int64             `json:"timestamp"` // Receive timestamp (milliseconds)
	Headers   map[string]string `json:"headers"`   // Additional headers (traceparent, etc.)
}

// NewRawMessage creates a new RawMessage with the given parameters.
func NewRawMessage(vendor, topic string, payload []byte, qos int) *RawMessage {
	return &RawMessage{
		Vendor:    vendor,
		Topic:     topic,
		Payload:   payload,
		QoS:       qos,
		Timestamp: time.Now().UnixMilli(),
		Headers:   make(map[string]string),
	}
}

// Validate validates the raw message.
func (rm *RawMessage) Validate() error {
	if rm.Vendor == "" {
		return errors.New("vendor is required")
	}
	if rm.Topic == "" {
		return errors.New("topic is required")
	}
	if len(rm.Payload) == 0 {
		return errors.New("payload is required")
	}
	return nil
}

// WithHeader adds a header to the raw message and returns the message for chaining.
func (rm *RawMessage) WithHeader(key, value string) *RawMessage {
	if rm.Headers == nil {
		rm.Headers = make(map[string]string)
	}
	rm.Headers[key] = value
	return rm
}

// GetHeader returns the value of a header.
func (rm *RawMessage) GetHeader(key string) string {
	if rm.Headers == nil {
		return ""
	}
	return rm.Headers[key]
}
