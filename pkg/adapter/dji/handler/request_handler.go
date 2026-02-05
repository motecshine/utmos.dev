package handler

import (
	"context"
	"fmt"

	dji "github.com/utmos/utmos/pkg/adapter/dji"
	"github.com/utmos/utmos/pkg/adapter/dji/router"
	"github.com/utmos/utmos/pkg/rabbitmq"
)

// RequestHandler handles device-initiated request messages.
type RequestHandler struct {
	router *router.ServiceRouter
	cfg    MessageConfig
}

// NewRequestHandler creates a new Request handler.
func NewRequestHandler(r *router.ServiceRouter) *RequestHandler {
	return &RequestHandler{
		router: r,
		cfg: MessageConfig{
			ReplyTopicType: dji.TopicTypeRequestsReply,
			RequestAction:  dji.ActionDeviceRequest,
			ReplyAction:    dji.ActionDeviceRequestReply,
			MessageType:    "device_request",
			ReplyType:      "device_request_reply",
		},
	}
}

// Handle processes a Request message and returns a StandardMessage.
func (h *RequestHandler) Handle(ctx context.Context, msg *dji.Message, topic *dji.TopicInfo) (*rabbitmq.StandardMessage, error) {
	sm, err := HandleMessage(msg, topic, h.cfg, DefaultDataBuilder)
	if err != nil {
		return nil, fmt.Errorf("failed to build request data: %w", err)
	}
	return sm, nil
}

// GetTopicType returns the topic type this handler processes.
func (h *RequestHandler) GetTopicType() dji.TopicType {
	return dji.TopicTypeRequests
}

// GetRouter returns the service router.
func (h *RequestHandler) GetRouter() *router.ServiceRouter {
	return h.router
}

// Ensure RequestHandler implements Handler interface.
var _ Handler = (*RequestHandler)(nil)
