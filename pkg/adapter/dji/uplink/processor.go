// Package uplink provides DJI uplink message processing functionality
package uplink

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"

	dji "github.com/utmos/utmos/pkg/adapter/dji"
	"github.com/utmos/utmos/pkg/rabbitmq"
)

// MessageType represents the type of message being processed
type MessageType string

const (
	// MessageTypeProperty is the property message type.
	MessageTypeProperty MessageType = "property"
	// MessageTypeEvent is the event message type.
	MessageTypeEvent MessageType = "event"
	// MessageTypeService is the service message type.
	MessageTypeService MessageType = "service"
	// MessageTypeStatus is the status message type.
	MessageTypeStatus MessageType = "status"
)

// ProcessedMessage represents a message after processing
type ProcessedMessage struct {
	Original    *rabbitmq.StandardMessage
	MessageType MessageType
	DeviceSN    string
	Vendor      string
	Properties  map[string]any
	Events      []Event
	Timestamp   int64
}

// Event represents a device event
type Event struct {
	Name   string         `json:"name"`
	Params map[string]any `json:"params"`
	Output map[string]any `json:"output,omitempty"`
}

// Processor processes DJI protocol messages
type Processor struct {
	vendor string
	logger *logrus.Entry
}

// NewProcessor creates a new DJI processor
func NewProcessor(logger *logrus.Entry) *Processor {
	if logger == nil {
		logger = logrus.NewEntry(logrus.StandardLogger())
	}
	return &Processor{
		vendor: dji.VendorDJI,
		logger: logger.WithField("processor", dji.VendorDJI),
	}
}

// GetVendor returns the vendor name
func (p *Processor) GetVendor() string {
	return p.vendor
}

// CanProcess checks if this processor can handle the given message
func (p *Processor) CanProcess(msg *rabbitmq.StandardMessage) bool {
	if msg.ProtocolMeta != nil && msg.ProtocolMeta.Vendor == dji.VendorDJI {
		return true
	}
	// Check action prefix for DJI-specific actions
	return strings.HasPrefix(msg.Action, "dji.") ||
		strings.HasPrefix(msg.Action, "property.") ||
		strings.HasPrefix(msg.Action, "event.") ||
		strings.HasPrefix(msg.Action, "device.")
}

// Process processes a DJI message
func (p *Processor) Process(ctx context.Context, msg *rabbitmq.StandardMessage) (*ProcessedMessage, error) {
	if msg == nil {
		return nil, fmt.Errorf("message is nil")
	}

	processed := &ProcessedMessage{
		Original:   msg,
		DeviceSN:   msg.DeviceSN,
		Vendor:     dji.VendorDJI,
		Properties: make(map[string]any),
		Events:     make([]Event, 0),
		Timestamp:  msg.Timestamp,
	}

	// Determine message type from action
	processed.MessageType = p.determineMessageType(msg.Action)

	// Parse the data based on message type
	switch processed.MessageType {
	case MessageTypeProperty:
		if err := p.processPropertyMessage(msg, processed); err != nil {
			return nil, fmt.Errorf("failed to process property message: %w", err)
		}
	case MessageTypeEvent:
		if err := p.processEventMessage(msg, processed); err != nil {
			return nil, fmt.Errorf("failed to process event message: %w", err)
		}
	case MessageTypeService:
		if err := p.processServiceMessage(msg, processed); err != nil {
			return nil, fmt.Errorf("failed to process service message: %w", err)
		}
	case MessageTypeStatus:
		if err := p.processStatusMessage(msg, processed); err != nil {
			return nil, fmt.Errorf("failed to process status message: %w", err)
		}
	default:
		// For unknown types, just extract raw data
		if err := p.processRawMessage(msg, processed); err != nil {
			return nil, fmt.Errorf("failed to process raw message: %w", err)
		}
	}

	p.logger.WithFields(logrus.Fields{
		"device_sn":    processed.DeviceSN,
		"message_type": processed.MessageType,
		"tid":          msg.TID,
	}).Debug("Processed DJI message")

	return processed, nil
}

// determineMessageType determines the message type from the action
func (p *Processor) determineMessageType(action string) MessageType {
	switch {
	case strings.HasPrefix(action, "property."):
		return MessageTypeProperty
	case strings.HasPrefix(action, "event."):
		return MessageTypeEvent
	case strings.HasPrefix(action, "service."):
		return MessageTypeService
	case strings.HasPrefix(action, "device.online"), strings.HasPrefix(action, "device.offline"):
		return MessageTypeStatus
	default:
		return MessageTypeProperty
	}
}

// processPropertyMessage processes property report messages
func (p *Processor) processPropertyMessage(msg *rabbitmq.StandardMessage, processed *ProcessedMessage) error {
	if len(msg.Data) == 0 {
		return nil
	}

	// Try to parse as DJI OSD/State format
	var data map[string]any
	if err := json.Unmarshal(msg.Data, &data); err != nil {
		return fmt.Errorf("failed to unmarshal property data: %w", err)
	}

	// Extract properties from DJI format
	// DJI messages may have nested structure like {"data": {...}}
	if dataField, ok := data["data"].(map[string]any); ok {
		processed.Properties = p.flattenProperties(dataField)
	} else {
		processed.Properties = p.flattenProperties(data)
	}

	return nil
}

// processEventMessage processes event messages
func (p *Processor) processEventMessage(msg *rabbitmq.StandardMessage, processed *ProcessedMessage) error {
	if len(msg.Data) == 0 {
		return nil
	}

	var data map[string]any
	if err := json.Unmarshal(msg.Data, &data); err != nil {
		return fmt.Errorf("failed to unmarshal event data: %w", err)
	}

	// Extract event information
	event := Event{
		Params: make(map[string]any),
	}

	// Get method/event name from protocol meta or data
	if msg.ProtocolMeta != nil && msg.ProtocolMeta.Method != "" {
		event.Name = msg.ProtocolMeta.Method
	} else if method, ok := data["method"].(string); ok {
		event.Name = method
	} else {
		event.Name = "unknown"
	}

	// Extract params
	if params, ok := data["data"].(map[string]any); ok {
		event.Params = params
	} else {
		event.Params = data
	}

	// Extract output if present
	if output, ok := data["output"].(map[string]any); ok {
		event.Output = output
	}

	processed.Events = append(processed.Events, event)

	return nil
}

// processServiceMessage processes service call messages
func (p *Processor) processServiceMessage(msg *rabbitmq.StandardMessage, processed *ProcessedMessage) error {
	if len(msg.Data) == 0 {
		return nil
	}

	var data map[string]any
	if err := json.Unmarshal(msg.Data, &data); err != nil {
		return fmt.Errorf("failed to unmarshal service data: %w", err)
	}

	// Store service call data as properties
	processed.Properties = data

	return nil
}

// processStatusMessage processes device status messages
func (p *Processor) processStatusMessage(msg *rabbitmq.StandardMessage, processed *ProcessedMessage) error {
	if len(msg.Data) == 0 {
		return nil
	}

	var data map[string]any
	if err := json.Unmarshal(msg.Data, &data); err != nil {
		return fmt.Errorf("failed to unmarshal status data: %w", err)
	}

	// Extract online status
	if online, ok := data["online"].(bool); ok {
		processed.Properties["online"] = online
	}

	// Determine online status from action
	if strings.HasSuffix(msg.Action, ".online") {
		processed.Properties["online"] = true
	} else if strings.HasSuffix(msg.Action, ".offline") {
		processed.Properties["online"] = false
	}

	return nil
}

// processRawMessage processes unknown message types
func (p *Processor) processRawMessage(msg *rabbitmq.StandardMessage, processed *ProcessedMessage) error {
	if len(msg.Data) == 0 {
		return nil
	}

	var data map[string]any
	if err := json.Unmarshal(msg.Data, &data); err != nil {
		// If not JSON object, store as raw
		processed.Properties["raw"] = string(msg.Data)
		return nil
	}

	processed.Properties = data
	return nil
}

// flattenProperties flattens nested properties for easier storage
func (p *Processor) flattenProperties(data map[string]any) map[string]any {
	result := make(map[string]any)
	p.flattenPropertiesRecursive("", data, result)
	return result
}

// flattenPropertiesRecursive recursively flattens nested maps
func (p *Processor) flattenPropertiesRecursive(prefix string, data map[string]any, result map[string]any) {
	for key, value := range data {
		fullKey := key
		if prefix != "" {
			fullKey = prefix + "." + key
		}

		switch v := value.(type) {
		case map[string]any:
			// Recursively flatten nested maps
			p.flattenPropertiesRecursive(fullKey, v, result)
		default:
			result[fullKey] = value
		}
	}
}
