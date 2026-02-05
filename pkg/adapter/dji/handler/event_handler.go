package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	dji "github.com/utmos/utmos/pkg/adapter/dji"
	"github.com/utmos/utmos/pkg/adapter/dji/router"
	"github.com/utmos/utmos/pkg/rabbitmq"
)

// EventHandler handles Event messages.
type EventHandler struct {
	router *router.EventRouter
}

// NewEventHandler creates a new Event handler.
func NewEventHandler(r *router.EventRouter) *EventHandler {
	return &EventHandler{
		router: r,
	}
}

// Handle processes an Event message and returns a StandardMessage.
func (h *EventHandler) Handle(ctx context.Context, msg *dji.Message, topic *dji.TopicInfo) (*rabbitmq.StandardMessage, error) {
	if msg == nil {
		return nil, fmt.Errorf("nil message")
	}
	if topic == nil {
		return nil, fmt.Errorf("nil topic info")
	}

	// Determine if this is a request or reply
	isReply := topic.Type == dji.TopicTypeEventsReply

	// Build StandardMessage
	action := dji.ActionEventReport
	if isReply {
		action = dji.ActionEventReply
	}

	sm := &rabbitmq.StandardMessage{
		TID:       msg.TID,
		BID:       msg.BID,
		Timestamp: msg.Timestamp,
		DeviceSN:  topic.DeviceSN,
		Service:   dji.VendorDJI,
		Action:    action,
		ProtocolMeta: &rabbitmq.ProtocolMeta{
			Vendor:        dji.VendorDJI,
			OriginalTopic: topic.Raw,
			Method:        msg.Method,
		},
	}

	// Set timestamp if not provided
	if sm.Timestamp == 0 {
		sm.Timestamp = time.Now().UnixMilli()
	}

	// Build event data
	var data json.RawMessage
	var err error

	if isReply {
		data, err = h.buildEventReplyData(msg, topic)
	} else {
		data, err = h.buildEventRequestData(msg, topic)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to build event data: %w", err)
	}
	sm.Data = data

	return sm, nil
}

// GetTopicType returns the topic type this handler processes.
func (h *EventHandler) GetTopicType() dji.TopicType {
	return dji.TopicTypeEvents
}

// buildEventRequestData builds data for event request.
func (h *EventHandler) buildEventRequestData(msg *dji.Message, topic *dji.TopicInfo) (json.RawMessage, error) {
	result := make(map[string]interface{})

	result["device_sn"] = topic.DeviceSN
	result["gateway_sn"] = topic.GatewaySN
	result["message_type"] = "event_request"
	result["method"] = msg.Method

	// Include raw data
	if len(msg.Data) > 0 {
		var data interface{}
		if err := json.Unmarshal(msg.Data, &data); err == nil {
			result["data"] = data
		} else {
			result["raw_data"] = string(msg.Data)
		}
	}

	// Include need_reply flag
	result["need_reply"] = msg.NeedReplyBool()

	return json.Marshal(result)
}

// buildEventReplyData builds data for event reply.
func (h *EventHandler) buildEventReplyData(msg *dji.Message, topic *dji.TopicInfo) (json.RawMessage, error) {
	result := make(map[string]interface{})

	result["device_sn"] = topic.DeviceSN
	result["gateway_sn"] = topic.GatewaySN
	result["message_type"] = "event_reply"
	result["method"] = msg.Method

	// Parse reply data to extract result code
	if len(msg.Data) > 0 {
		var replyData map[string]interface{}
		if err := json.Unmarshal(msg.Data, &replyData); err == nil {
			// Extract result code
			if resultCode, ok := replyData["result"].(float64); ok {
				result["result"] = int(resultCode)
			}
			// Include output data
			if output, ok := replyData["output"]; ok {
				result["output"] = output
			}
			// Include full reply data
			result["data"] = replyData
		} else {
			result["raw_data"] = string(msg.Data)
		}
	}

	return json.Marshal(result)
}

// GetRouter returns the event router.
func (h *EventHandler) GetRouter() *router.EventRouter {
	return h.router
}

// Ensure EventHandler implements Handler interface.
var _ Handler = (*EventHandler)(nil)
