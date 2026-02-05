package dji

import (
	"context"

	"github.com/utmos/utmos/pkg/adapter"
	"github.com/utmos/utmos/pkg/rabbitmq"
)

// MessageHandler defines the interface for message handlers.
type MessageHandler interface {
	Handle(ctx context.Context, msg *Message, topic *TopicInfo) (*rabbitmq.StandardMessage, error)
	GetTopicType() TopicType
}

// HandlerRegistry defines the interface for handler registry.
type HandlerRegistry interface {
	Get(topicType TopicType) (MessageHandler, error)
}

// Adapter implements the ProtocolAdapter interface for DJI devices.
type Adapter struct {
	converter *Converter
	registry  HandlerRegistry
}

// NewAdapter creates a new DJI adapter.
func NewAdapter() *Adapter {
	return &Adapter{
		converter: NewConverter(),
	}
}

// SetHandlerRegistry sets the handler registry for the adapter.
func (a *Adapter) SetHandlerRegistry(registry HandlerRegistry) {
	a.registry = registry
}

// HandleMessage processes a raw message using the appropriate handler.
func (a *Adapter) HandleMessage(ctx context.Context, topic string, payload []byte) (*rabbitmq.StandardMessage, error) {
	// Parse the topic
	topicInfo, err := ParseTopic(topic)
	if err != nil {
		return nil, err
	}

	// Parse the message
	msg, err := ParseMessage(payload)
	if err != nil {
		return nil, err
	}

	// If registry is set, try to get handler
	if a.registry != nil {
		h, err := a.registry.Get(topicInfo.Type)
		if err == nil {
			return h.Handle(ctx, msg, topicInfo)
		}
	}

	// Fall back to converter for unsupported topic types
	return a.converter.ToStandardMessage(msg, topicInfo)
}

// GetVendor returns the vendor identifier.
func (a *Adapter) GetVendor() string {
	return VendorDJI
}

// ParseRawMessage parses raw bytes into a protocol message.
func (a *Adapter) ParseRawMessage(topic string, payload []byte) (*adapter.ProtocolMessage, error) {
	// Parse the topic
	topicInfo, err := ParseTopic(topic)
	if err != nil {
		return nil, err
	}

	// Parse the message
	djiMsg, err := ParseMessage(payload)
	if err != nil {
		return nil, err
	}

	// Determine message type
	msgType := mapTopicTypeToMessageType(topicInfo.Type)

	return &adapter.ProtocolMessage{
		Vendor:      VendorDJI,
		Topic:       topic,
		DeviceSN:    topicInfo.DeviceSN,
		GatewaySN:   topicInfo.GatewaySN,
		MessageType: msgType,
		Method:      djiMsg.Method,
		TID:         djiMsg.TID,
		BID:         djiMsg.BID,
		Timestamp:   djiMsg.Timestamp,
		NeedReply:   djiMsg.NeedReplyBool(),
		Data:        djiMsg.Data,
	}, nil
}

// ToStandardMessage converts a protocol message to a standard message.
func (a *Adapter) ToStandardMessage(pm *adapter.ProtocolMessage) (*rabbitmq.StandardMessage, error) {
	// Reconstruct DJI message from protocol message
	var needReply *int
	if pm.NeedReply {
		one := 1
		needReply = &one
	}

	djiMsg := &Message{
		TID:       pm.TID,
		BID:       pm.BID,
		Timestamp: pm.Timestamp,
		Method:    pm.Method,
		NeedReply: needReply,
		Data:      pm.Data,
	}

	// Parse topic info
	topicInfo, err := ParseTopic(pm.Topic)
	if err != nil {
		// If topic parsing fails, create minimal topic info
		topicInfo = &TopicInfo{
			Type:     mapMessageTypeToTopicType(pm.MessageType),
			DeviceSN: pm.DeviceSN,
			Raw:      pm.Topic,
		}
	}

	return a.converter.ToStandardMessage(djiMsg, topicInfo)
}

// FromStandardMessage converts a standard message to a protocol message.
func (a *Adapter) FromStandardMessage(sm *rabbitmq.StandardMessage) (*adapter.ProtocolMessage, error) {
	djiMsg, err := a.converter.FromStandardMessage(sm)
	if err != nil {
		return nil, err
	}

	// Determine message type from action
	msgType := mapActionToMessageType(sm.Action)

	// Build topic
	topicType := MapActionToTopicType(sm.Action)
	topic := BuildTopic(topicType, sm.DeviceSN)

	return &adapter.ProtocolMessage{
		Vendor:      VendorDJI,
		Topic:       topic,
		DeviceSN:    sm.DeviceSN,
		MessageType: msgType,
		Method:      djiMsg.Method,
		TID:         djiMsg.TID,
		BID:         djiMsg.BID,
		Timestamp:   djiMsg.Timestamp,
		Data:        djiMsg.Data,
	}, nil
}

// GetRawPayload returns the raw payload for sending to device.
func (a *Adapter) GetRawPayload(pm *adapter.ProtocolMessage) ([]byte, error) {
	var needReply *int
	if pm.NeedReply {
		one := 1
		needReply = &one
	}

	djiMsg := &Message{
		TID:       pm.TID,
		BID:       pm.BID,
		Timestamp: pm.Timestamp,
		Method:    pm.Method,
		NeedReply: needReply,
		Data:      pm.Data,
	}

	return djiMsg.ToJSON()
}

// mapTopicTypeToMessageType maps DJI topic type to adapter message type.
func mapTopicTypeToMessageType(tt TopicType) adapter.MessageType {
	switch tt {
	case TopicTypeOSD, TopicTypeState:
		return adapter.MessageTypeProperty
	case TopicTypeEvents:
		return adapter.MessageTypeEvent
	case TopicTypeServices, TopicTypeServicesReply:
		return adapter.MessageTypeService
	case TopicTypeStatus, TopicTypeStatusReply:
		return adapter.MessageTypeStatus
	default:
		return adapter.MessageTypeProperty
	}
}

// mapMessageTypeToTopicType maps adapter message type to DJI topic type.
func mapMessageTypeToTopicType(mt adapter.MessageType) TopicType {
	switch mt {
	case adapter.MessageTypeProperty:
		return TopicTypeOSD
	case adapter.MessageTypeEvent:
		return TopicTypeEvents
	case adapter.MessageTypeService:
		return TopicTypeServices
	case adapter.MessageTypeStatus:
		return TopicTypeStatus
	default:
		return TopicTypeOSD
	}
}

// mapActionToMessageType maps standard action to adapter message type.
func mapActionToMessageType(action string) adapter.MessageType {
	switch action {
	case ActionPropertyReport, ActionPropertySet:
		return adapter.MessageTypeProperty
	case ActionEventReport:
		return adapter.MessageTypeEvent
	case ActionServiceCall, ActionServiceReply:
		return adapter.MessageTypeService
	case ActionDeviceOnline, ActionDeviceOffline:
		return adapter.MessageTypeStatus
	default:
		return adapter.MessageTypeProperty
	}
}

// Register registers the DJI adapter with the global registry.
func Register() {
	adapter.Register(NewAdapter())
}

// Ensure Adapter implements ProtocolAdapter interface.
var _ adapter.ProtocolAdapter = (*Adapter)(nil)
