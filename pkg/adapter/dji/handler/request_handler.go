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

// RequestHandler handles device-initiated request messages.
type RequestHandler struct {
	router *router.ServiceRouter
}

// NewRequestHandler creates a new Request handler.
func NewRequestHandler(r *router.ServiceRouter) *RequestHandler {
	return &RequestHandler{
		router: r,
	}
}

// Handle processes a Request message and returns a StandardMessage.
func (h *RequestHandler) Handle(ctx context.Context, msg *dji.Message, topic *dji.TopicInfo) (*rabbitmq.StandardMessage, error) {
	if msg == nil {
		return nil, fmt.Errorf("nil message")
	}
	if topic == nil {
		return nil, fmt.Errorf("nil topic info")
	}

	// Determine if this is a request or reply
	isReply := topic.Type == dji.TopicTypeRequestsReply

	// Build StandardMessage
	action := dji.ActionDeviceRequest
	if isReply {
		action = dji.ActionDeviceRequestReply
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

	// Build request data
	var data json.RawMessage
	var err error

	if isReply {
		data, err = h.buildRequestReplyData(msg, topic)
	} else {
		data, err = h.buildRequestData(msg, topic)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to build request data: %w", err)
	}
	sm.Data = data

	return sm, nil
}

// GetTopicType returns the topic type this handler processes.
func (h *RequestHandler) GetTopicType() dji.TopicType {
	return dji.TopicTypeRequests
}

// buildRequestData builds data for device request.
func (h *RequestHandler) buildRequestData(msg *dji.Message, topic *dji.TopicInfo) (json.RawMessage, error) {
	result := make(map[string]interface{})

	result["device_sn"] = topic.DeviceSN
	result["gateway_sn"] = topic.GatewaySN
	result["message_type"] = "device_request"
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

// buildRequestReplyData builds data for device request reply.
func (h *RequestHandler) buildRequestReplyData(msg *dji.Message, topic *dji.TopicInfo) (json.RawMessage, error) {
	result := make(map[string]interface{})

	result["device_sn"] = topic.DeviceSN
	result["gateway_sn"] = topic.GatewaySN
	result["message_type"] = "device_request_reply"
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

// GetRouter returns the service router.
func (h *RequestHandler) GetRouter() *router.ServiceRouter {
	return h.router
}

// Ensure RequestHandler implements Handler interface.
var _ Handler = (*RequestHandler)(nil)
