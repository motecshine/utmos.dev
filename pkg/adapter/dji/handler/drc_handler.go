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
	cfg              MessageConfig
}

// NewDRCHandler creates a new DRC handler.
func NewDRCHandler(sr *router.ServiceRouter, er *router.EventRouter) *DRCHandler {
	return &DRCHandler{
		serviceRouter:    sr,
		eventRouter:      er,
		heartbeatTimeout: config.DRCHeartbeatTimeout,
		cfg: MessageConfig{
			ReplyTopicType: dji.TopicTypeDRCDown,
			RequestAction:  dji.ActionDRCCommand,
			ReplyAction:    dji.ActionDRCEvent,
			MessageType:    "drc_command",
			ReplyType:      "drc_event",
		},
	}
}

// Handle processes a DRC message and returns a StandardMessage.
func (h *DRCHandler) Handle(_ context.Context, msg *dji.Message, topic *dji.TopicInfo) (*rabbitmq.StandardMessage, error) {
	builder := func(msg *dji.Message, topic *dji.TopicInfo, isReply bool, cfg MessageConfig) (json.RawMessage, error) {
		if isReply {
			return h.buildDRCEventData(msg, topic)
		}
		extraFields := map[string]any{
			"heartbeat_timeout_ms": h.heartbeatTimeout.Milliseconds(),
		}
		return BuildRequestData(msg, topic, cfg.MessageType, extraFields)
	}

	sm, err := HandleMessage(msg, topic, h.cfg, builder)
	if err != nil {
		return nil, fmt.Errorf("failed to build DRC data: %w", err)
	}
	return sm, nil
}

// GetTopicType returns the topic type this handler processes.
func (h *DRCHandler) GetTopicType() dji.TopicType {
	return dji.TopicTypeDRCUp
}

// buildDRCEventData builds data for DRC event (down).
// DRC events don't have need_reply field, so we use a custom builder.
func (h *DRCHandler) buildDRCEventData(msg *dji.Message, topic *dji.TopicInfo) (json.RawMessage, error) {
	result := map[string]any{
		"device_sn":    topic.DeviceSN,
		"gateway_sn":   topic.GatewaySN,
		"message_type": "drc_event",
		"method":       msg.Method,
	}

	tryUnmarshalData(msg.Data, result)

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
