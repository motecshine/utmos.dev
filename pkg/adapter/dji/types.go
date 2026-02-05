// Package dji provides the DJI protocol adapter implementation.
package dji

import (
	"encoding/json"

	"github.com/utmos/utmos/pkg/rabbitmq"
)

// Re-export constants from rabbitmq package for convenience.
const (
	VendorDJI            = rabbitmq.VendorDJI
	ActionPropertyReport = rabbitmq.ActionPropertyReport
	ActionPropertySet    = rabbitmq.ActionPropertySet
	ActionEventReport    = rabbitmq.ActionEventReport
	ActionServiceCall    = rabbitmq.ActionServiceCall
	ActionServiceReply   = rabbitmq.ActionServiceReply
	ActionDeviceOnline   = rabbitmq.ActionDeviceOnline
	ActionDeviceOffline  = rabbitmq.ActionDeviceOffline
)

// DJI-specific actions not in rabbitmq package.
const (
	ActionStatusReply        = "status.reply"
	ActionEventReply         = "event.reply"
	ActionDeviceRequest      = "device.request"
	ActionDeviceRequestReply = "device.request.reply"
	ActionDRCCommand         = "drc.command"
	ActionDRCEvent           = "drc.event"
)

// TopicType represents the type of DJI MQTT topic.
type TopicType string

// Topic type constants.
const (
	TopicTypeOSD           TopicType = "osd"
	TopicTypeState         TopicType = "state"
	TopicTypeServices      TopicType = "services"
	TopicTypeServicesReply TopicType = "services_reply"
	TopicTypeEvents        TopicType = "events"
	TopicTypeEventsReply   TopicType = "events_reply"
	TopicTypeStatus        TopicType = "status"
	TopicTypeStatusReply   TopicType = "status_reply"
	TopicTypeRequests      TopicType = "requests"
	TopicTypeRequestsReply TopicType = "requests_reply"
	TopicTypeDRCUp         TopicType = "drc/up"
	TopicTypeDRCDown       TopicType = "drc/down"
)

// Direction represents the message direction.
type Direction string

// Direction constants.
const (
	DirectionUplink   Direction = "uplink"
	DirectionDownlink Direction = "downlink"
)

// Message represents a DJI protocol message.
type Message struct {
	TID       string          `json:"tid"`                  // Transaction ID
	BID       string          `json:"bid"`                  // Business ID
	Timestamp int64           `json:"timestamp,omitempty"`  // Message timestamp (milliseconds)
	Method    string          `json:"method,omitempty"`     // Protocol method
	NeedReply *int            `json:"need_reply,omitempty"` // 0 or 1 to indicate if reply is needed
	Data      json.RawMessage `json:"data,omitempty"`       // Business data
	GatewaySN string          `json:"gateway,omitempty"`    // Gateway serial number
}

// Validate validates the DJI message.
func (m *Message) Validate() error {
	if m.TID == "" {
		return ErrMissingTID
	}
	if m.BID == "" {
		return ErrMissingBID
	}
	return nil
}

// NeedReplyBool returns the need_reply field as a boolean.
func (m *Message) NeedReplyBool() bool {
	if m.NeedReply == nil {
		return false
	}
	return *m.NeedReply == 1
}

// ToJSON serializes the message to JSON bytes.
func (m *Message) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}
