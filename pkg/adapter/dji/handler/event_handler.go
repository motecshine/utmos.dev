package handler

import (
	"context"
	"fmt"

	dji "github.com/utmos/utmos/pkg/adapter/dji"
	"github.com/utmos/utmos/pkg/adapter/dji/router"
	"github.com/utmos/utmos/pkg/rabbitmq"
)

// EventHandler handles Event messages.
//nolint:dupl // Structure similar to RequestHandler but types differ
type EventHandler struct {
	router *router.EventRouter
	cfg    MessageConfig
}

// NewEventHandler creates a new Event handler.
func NewEventHandler(r *router.EventRouter) *EventHandler {
	return &EventHandler{
		router: r,
		cfg: MessageConfig{
			ReplyTopicType: dji.TopicTypeEventsReply,
			RequestAction:  dji.ActionEventReport,
			ReplyAction:    dji.ActionEventReply,
			MessageType:    "event_request",
			ReplyType:      "event_reply",
		},
	}
}

// Handle processes an Event message and returns a StandardMessage.
func (h *EventHandler) Handle(_ context.Context, msg *dji.Message, topic *dji.TopicInfo) (*rabbitmq.StandardMessage, error) {
	sm, err := HandleMessage(msg, topic, h.cfg, DefaultDataBuilder)
	if err != nil {
		return nil, fmt.Errorf("failed to build event data: %w", err)
	}
	return sm, nil
}

// GetTopicType returns the topic type this handler processes.
func (h *EventHandler) GetTopicType() dji.TopicType {
	return dji.TopicTypeEvents
}

// GetRouter returns the event router.
func (h *EventHandler) GetRouter() *router.EventRouter {
	return h.router
}

// Ensure EventHandler implements Handler interface.
var _ Handler = (*EventHandler)(nil)
