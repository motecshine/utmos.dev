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

// DRCHandler handles DRC (Drone Remote Control) messages.
type DRCHandler struct {
	serviceRouter    *router.ServiceRouter
	eventRouter      *router.EventRouter
	heartbeatTimeout time.Duration
}

// NewDRCHandler creates a new DRC handler.
func NewDRCHandler(sr *router.ServiceRouter, er *router.EventRouter) *DRCHandler {
	return &DRCHandler{
		serviceRouter:    sr,
		eventRouter:      er,
		heartbeatTimeout: config.DRCHeartbeatTimeout,
	}
}

// Handle processes a DRC message and returns a StandardMessage.
func (h *DRCHandler) Handle(ctx context.Context, msg *dji.Message, topic *dji.TopicInfo) (*rabbitmq.StandardMessage, error) {
	if msg == nil {
		return nil, fmt.Errorf("nil message")
	}
	if topic == nil {
		return nil, fmt.Errorf("nil topic info")
	}

	// Determine if this is up (command) or down (event)
	isDown := topic.Type == dji.TopicTypeDRCDown

	// Build StandardMessage
	action := dji.ActionDRCCommand
	if isDown {
		action = dji.ActionDRCEvent
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

	// Build DRC data
	var data json.RawMessage
	var err error

	if isDown {
		data, err = h.buildDRCEventData(msg, topic)
	} else {
		data, err = h.buildDRCCommandData(msg, topic)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to build DRC data: %w", err)
	}
	sm.Data = data

	return sm, nil
}

// GetTopicType returns the topic type this handler processes.
func (h *DRCHandler) GetTopicType() dji.TopicType {
	return dji.TopicTypeDRCUp
}

// buildDRCCommandData builds data for DRC command (up).
func (h *DRCHandler) buildDRCCommandData(msg *dji.Message, topic *dji.TopicInfo) (json.RawMessage, error) {
	result := make(map[string]interface{})

	result["device_sn"] = topic.DeviceSN
	result["gateway_sn"] = topic.GatewaySN
	result["message_type"] = "drc_command"
	result["method"] = msg.Method
	result["heartbeat_timeout_ms"] = h.heartbeatTimeout.Milliseconds()

	// Include raw data
	if len(msg.Data) > 0 {
		var data interface{}
		if err := json.Unmarshal(msg.Data, &data); err == nil {
			result["data"] = data
		} else {
			result["raw_data"] = string(msg.Data)
		}
	}

	return json.Marshal(result)
}

// buildDRCEventData builds data for DRC event (down).
func (h *DRCHandler) buildDRCEventData(msg *dji.Message, topic *dji.TopicInfo) (json.RawMessage, error) {
	result := make(map[string]interface{})

	result["device_sn"] = topic.DeviceSN
	result["gateway_sn"] = topic.GatewaySN
	result["message_type"] = "drc_event"
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

	return json.Marshal(result)
}

// SetHeartbeatTimeout sets the DRC heartbeat timeout.
func (h *DRCHandler) SetHeartbeatTimeout(timeout time.Duration) {
	h.heartbeatTimeout = timeout
}

// GetServiceRouter returns the service router.
func (h *DRCHandler) GetServiceRouter() *router.ServiceRouter {
	return h.serviceRouter
}

// GetEventRouter returns the event router.
func (h *DRCHandler) GetEventRouter() *router.EventRouter {
	return h.eventRouter
}

// Ensure DRCHandler implements Handler interface.
var _ Handler = (*DRCHandler)(nil)
