package dji

import (
	"time"

	"github.com/utmos/utmos/pkg/rabbitmq"
)

// Converter handles conversion between DJI messages and standard messages.
type Converter struct{}

// NewConverter creates a new Converter.
func NewConverter() *Converter {
	return &Converter{}
}

// ToStandardMessage converts a DJI message to a standard message.
func (c *Converter) ToStandardMessage(djiMsg *Message, topicInfo *TopicInfo) (*rabbitmq.StandardMessage, error) {
	action := MapTopicTypeToAction(topicInfo.Type)

	// Build protocol metadata
	needReply := djiMsg.NeedReplyBool()
	protocolMeta := &rabbitmq.ProtocolMeta{
		Vendor:        VendorDJI,
		OriginalTopic: topicInfo.Raw,
		Method:        djiMsg.Method,
		NeedReply:     &needReply,
	}

	timestamp := djiMsg.Timestamp
	if timestamp == 0 {
		timestamp = time.Now().UnixMilli()
	}

	return &rabbitmq.StandardMessage{
		TID:          djiMsg.TID,
		BID:          djiMsg.BID,
		Timestamp:    timestamp,
		Service:      "dji-adapter",
		Action:       action,
		DeviceSN:     topicInfo.DeviceSN,
		Data:         djiMsg.Data,
		ProtocolMeta: protocolMeta,
	}, nil
}

// FromStandardMessage converts a standard message to a DJI message.
func (c *Converter) FromStandardMessage(stdMsg *rabbitmq.StandardMessage) (*Message, error) {
	var method string
	if stdMsg.ProtocolMeta != nil {
		method = stdMsg.ProtocolMeta.Method
	}

	return &Message{
		TID:       stdMsg.TID,
		BID:       stdMsg.BID,
		Timestamp: stdMsg.Timestamp,
		Method:    method,
		Data:      stdMsg.Data,
	}, nil
}

// MapTopicTypeToAction maps a DJI topic type to a standard action.
func MapTopicTypeToAction(topicType TopicType) string {
	switch topicType {
	case TopicTypeOSD, TopicTypeState:
		return ActionPropertyReport
	case TopicTypeEvents:
		return ActionEventReport
	case TopicTypeStatus:
		return ActionDeviceOnline
	case TopicTypeServicesReply:
		return ActionServiceReply
	case TopicTypeServices:
		return ActionServiceCall
	case TopicTypeStatusReply:
		return ActionStatusReply
	default:
		return "unknown"
	}
}

// MapActionToTopicType maps a standard action to a DJI topic type.
func MapActionToTopicType(action string) TopicType {
	switch action {
	case ActionPropertyReport:
		return TopicTypeOSD
	case ActionEventReport:
		return TopicTypeEvents
	case ActionDeviceOnline, ActionDeviceOffline:
		return TopicTypeStatus
	case ActionServiceCall:
		return TopicTypeServices
	case ActionServiceReply:
		return TopicTypeServicesReply
	default:
		return TopicTypeServices
	}
}
