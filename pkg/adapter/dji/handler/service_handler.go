package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	dji "github.com/utmos/utmos/pkg/adapter/dji"
	"github.com/utmos/utmos/pkg/adapter/dji/config"
	"github.com/utmos/utmos/pkg/adapter/dji/router"
	"github.com/utmos/utmos/pkg/rabbitmq"
)

// ServiceHandler handles Service call and reply messages.
type ServiceHandler struct {
	router  *router.ServiceRouter
	timeout time.Duration
}

// NewServiceHandler creates a new Service handler.
func NewServiceHandler(r *router.ServiceRouter) *ServiceHandler {
	return &ServiceHandler{
		router:  r,
		timeout: config.ServiceCallTimeout,
	}
}

// Handle processes a Service message and returns a StandardMessage.
func (h *ServiceHandler) Handle(ctx context.Context, msg *dji.Message, topic *dji.TopicInfo) (*rabbitmq.StandardMessage, error) {
	if msg == nil {
		return nil, fmt.Errorf("nil message")
	}
	if topic == nil {
		return nil, fmt.Errorf("nil topic info")
	}

	// Determine if this is a request or reply
	isReply := topic.Type == dji.TopicTypeServicesReply

	// Build StandardMessage
	action := dji.ActionServiceCall
	if isReply {
		action = dji.ActionServiceReply
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

	// Build service data
	var data json.RawMessage
	var err error

	if isReply {
		data, err = h.buildServiceReplyData(msg, topic)
	} else {
		data, err = h.buildServiceRequestData(msg, topic)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to build service data: %w", err)
	}
	sm.Data = data

	return sm, nil
}

// GetTopicType returns the topic type this handler processes.
// Note: This handler processes both services and services_reply.
func (h *ServiceHandler) GetTopicType() dji.TopicType {
	return dji.TopicTypeServices
}

// buildServiceRequestData builds data for service request.
func (h *ServiceHandler) buildServiceRequestData(msg *dji.Message, topic *dji.TopicInfo) (json.RawMessage, error) {
	result := make(map[string]interface{})

	result["device_sn"] = topic.DeviceSN
	result["gateway_sn"] = topic.GatewaySN
	result["message_type"] = "service_request"
	result["method"] = msg.Method
	result["timeout_ms"] = h.timeout.Milliseconds()

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

// buildServiceReplyData builds data for service reply.
func (h *ServiceHandler) buildServiceReplyData(msg *dji.Message, topic *dji.TopicInfo) (json.RawMessage, error) {
	result := make(map[string]interface{})

	result["device_sn"] = topic.DeviceSN
	result["gateway_sn"] = topic.GatewaySN
	result["message_type"] = "service_reply"
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

// SetTimeout sets the service call timeout.
func (h *ServiceHandler) SetTimeout(timeout time.Duration) {
	h.timeout = timeout
}

// GetRouter returns the service router.
func (h *ServiceHandler) GetRouter() *router.ServiceRouter {
	return h.router
}

// Ensure ServiceHandler implements Handler interface.
var _ Handler = (*ServiceHandler)(nil)
