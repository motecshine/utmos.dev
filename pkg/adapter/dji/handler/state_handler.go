package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	dji "github.com/utmos/utmos/pkg/adapter/dji"
	"github.com/utmos/utmos/pkg/rabbitmq"
)

// StateHandler handles State (property change) messages.
type StateHandler struct{}

// NewStateHandler creates a new State handler.
func NewStateHandler() *StateHandler {
	return &StateHandler{}
}

// Handle processes a State message and returns a StandardMessage.
func (h *StateHandler) Handle(ctx context.Context, msg *dji.Message, topic *dji.TopicInfo) (*rabbitmq.StandardMessage, error) {
	if msg == nil {
		return nil, fmt.Errorf("nil message")
	}
	if topic == nil {
		return nil, fmt.Errorf("nil topic info")
	}

	// Build StandardMessage
	sm := &rabbitmq.StandardMessage{
		TID:       msg.TID,
		BID:       msg.BID,
		Timestamp: msg.Timestamp,
		DeviceSN:  topic.DeviceSN,
		Service:   dji.VendorDJI,
		Action:    dji.ActionPropertyReport,
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

	// Build state data
	data, err := h.buildStateData(msg.Data, topic)
	if err != nil {
		return nil, fmt.Errorf("failed to build state data: %w", err)
	}
	sm.Data = data

	return sm, nil
}

// GetTopicType returns the topic type this handler processes.
func (h *StateHandler) GetTopicType() dji.TopicType {
	return dji.TopicTypeState
}

// buildStateData converts state data to a data map for StandardMessage.
func (h *StateHandler) buildStateData(data json.RawMessage, topic *dji.TopicInfo) (json.RawMessage, error) {
	result := make(map[string]interface{})

	result["device_sn"] = topic.DeviceSN
	result["gateway_sn"] = topic.GatewaySN
	result["message_type"] = "state"

	// Parse the raw state data
	if len(data) > 0 {
		var stateData map[string]interface{}
		if err := json.Unmarshal(data, &stateData); err != nil {
			return nil, fmt.Errorf("failed to parse state data: %w", err)
		}
		result["properties"] = stateData

		// Extract changed property names for quick access
		changedProps := make([]string, 0, len(stateData))
		for key := range stateData {
			changedProps = append(changedProps, key)
		}
		result["changed_properties"] = changedProps
	}

	return json.Marshal(result)
}

// Ensure StateHandler implements Handler interface.
var _ Handler = (*StateHandler)(nil)
