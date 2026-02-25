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
	cfg     MessageConfig
}

// NewServiceHandler creates a new Service handler.
func NewServiceHandler(r *router.ServiceRouter) *ServiceHandler {
	return &ServiceHandler{
		router:  r,
		timeout: config.ServiceCallTimeout,
		cfg: MessageConfig{
			ReplyTopicType: dji.TopicTypeServicesReply,
			RequestAction:  dji.ActionServiceCall,
			ReplyAction:    dji.ActionServiceReply,
			MessageType:    "service_request",
			ReplyType:      "service_reply",
		},
	}
}

// Handle processes a Service message and returns a StandardMessage.
func (h *ServiceHandler) Handle(_ context.Context, msg *dji.Message, topic *dji.TopicInfo) (*rabbitmq.StandardMessage, error) {
	builder := func(msg *dji.Message, topic *dji.TopicInfo, isReply bool, cfg MessageConfig) (json.RawMessage, error) {
		if isReply {
			return BuildReplyData(msg, topic, cfg.ReplyType)
		}
		extraFields := map[string]any{
			"timeout_ms": h.timeout.Milliseconds(),
		}
		return BuildRequestData(msg, topic, cfg.MessageType, extraFields)
	}

	sm, err := HandleMessage(msg, topic, h.cfg, builder)
	if err != nil {
		return nil, fmt.Errorf("failed to build service data: %w", err)
	}
	return sm, nil
}

// GetTopicType returns the topic type this handler processes.
func (h *ServiceHandler) GetTopicType() dji.TopicType {
	return dji.TopicTypeServices
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
