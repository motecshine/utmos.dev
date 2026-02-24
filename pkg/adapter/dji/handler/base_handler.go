package handler

import (
	"encoding/json"
	"fmt"
	"time"

	dji "github.com/utmos/utmos/pkg/adapter/dji"
	"github.com/utmos/utmos/pkg/rabbitmq"
)

// MessageConfig holds configuration for building StandardMessage.
type MessageConfig struct {
	ReplyTopicType dji.TopicType
	RequestAction  string
	ReplyAction    string
	MessageType    string
	ReplyType      string
}

// DataBuilder is a function type for building message data.
type DataBuilder func(msg *dji.Message, topic *dji.TopicInfo, isReply bool, cfg MessageConfig) (json.RawMessage, error)

// HandleMessage is a generic handler function that processes DJI messages.
func HandleMessage(msg *dji.Message, topic *dji.TopicInfo, cfg MessageConfig, builder DataBuilder) (*rabbitmq.StandardMessage, error) {
	if err := ValidateInputs(msg, topic); err != nil {
		return nil, err
	}

	sm := BuildStandardMessage(msg, topic, cfg)
	isReply := topic.Type == cfg.ReplyTopicType

	data, err := builder(msg, topic, isReply, cfg)
	if err != nil {
		return nil, err
	}
	sm.Data = data

	return sm, nil
}

// DefaultDataBuilder is the default data builder for request/reply messages.
func DefaultDataBuilder(msg *dji.Message, topic *dji.TopicInfo, isReply bool, cfg MessageConfig) (json.RawMessage, error) {
	if isReply {
		return BuildReplyData(msg, topic, cfg.ReplyType)
	}
	return BuildRequestData(msg, topic, cfg.MessageType, nil)
}

// BuildStandardMessage creates a StandardMessage from DJI message and topic.
func BuildStandardMessage(msg *dji.Message, topic *dji.TopicInfo, cfg MessageConfig) *rabbitmq.StandardMessage {
	isReply := topic.Type == cfg.ReplyTopicType

	action := cfg.RequestAction
	if isReply {
		action = cfg.ReplyAction
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

	if sm.Timestamp == 0 {
		sm.Timestamp = time.Now().UnixMilli()
	}

	return sm
}

// BuildRequestData builds data for request messages.
func BuildRequestData(msg *dji.Message, topic *dji.TopicInfo, messageType string, extraFields map[string]any) (json.RawMessage, error) {
	result := map[string]any{
		"device_sn":    topic.DeviceSN,
		"gateway_sn":   topic.GatewaySN,
		"message_type": messageType,
		"method":       msg.Method,
	}

	for k, v := range extraFields {
		result[k] = v
	}

	if len(msg.Data) > 0 {
		var data any
		if err := json.Unmarshal(msg.Data, &data); err == nil {
			result["data"] = data
		} else {
			result["raw_data"] = string(msg.Data)
		}
	}

	if msg.NeedReply != nil {
		result["need_reply"] = msg.NeedReplyBool()
	}

	return json.Marshal(result)
}

// BuildReplyData builds data for reply messages.
func BuildReplyData(msg *dji.Message, topic *dji.TopicInfo, messageType string) (json.RawMessage, error) {
	result := map[string]any{
		"device_sn":    topic.DeviceSN,
		"gateway_sn":   topic.GatewaySN,
		"message_type": messageType,
		"method":       msg.Method,
	}

	if len(msg.Data) > 0 {
		var replyData map[string]any
		if err := json.Unmarshal(msg.Data, &replyData); err == nil {
			if resultCode, ok := replyData["result"].(float64); ok {
				result["result"] = int(resultCode)
			}
			if output, ok := replyData["output"]; ok {
				result["output"] = output
			}
			result["data"] = replyData
		} else {
			result["raw_data"] = string(msg.Data)
		}
	}

	return json.Marshal(result)
}

// ValidateInputs validates message and topic are not nil.
func ValidateInputs(msg *dji.Message, topic *dji.TopicInfo) error {
	if msg == nil {
		return fmt.Errorf("nil message")
	}
	if topic == nil {
		return fmt.Errorf("nil topic info")
	}
	return nil
}
