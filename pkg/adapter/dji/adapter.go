package dji

import (
	"github.com/utmos/utmos/pkg/adapter"
	"github.com/utmos/utmos/pkg/rabbitmq"
)

// Adapter implements the ProtocolAdapter interface for DJI devices.
type Adapter struct {
	converter *Converter
}

// NewAdapter creates a new DJI adapter.
func NewAdapter() *Adapter {
	return &Adapter{
		converter: NewConverter(),
	}
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
