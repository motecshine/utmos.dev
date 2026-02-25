package router

import (
	"context"
	"encoding/json"
	"fmt"
)

// EventHandlerFunc is an alias for HandlerFunc used in event handler registration.
type EventHandlerFunc = HandlerFunc

// EventRequest represents an event request.
type EventRequest struct {
	Method    string          `json:"method"`
	NeedReply *int            `json:"need_reply,omitempty"`
	Data      json.RawMessage `json:"data,omitempty"`
}

// NeedReplyBool returns true if need_reply is set and non-zero.
func (r *EventRequest) NeedReplyBool() bool {
	return r.NeedReply != nil && *r.NeedReply != 0
}

// EventResponse is an alias for HandlerResponse used in event handlers.
type EventResponse = HandlerResponse

// EventRouter routes events to their handlers.
type EventRouter struct {
	registry handlerRegistry[HandlerFunc]
}

// NewEventRouter creates a new EventRouter.
func NewEventRouter() *EventRouter {
	return &EventRouter{
		registry: newHandlerRegistry[HandlerFunc](),
	}
}

// RegisterEventHandler registers a handler for an event method.
func (r *EventRouter) RegisterEventHandler(method string, handler EventHandlerFunc) error {
	return r.registry.register(method, handler, ErrMethodAlreadyRegistered)
}

// RouteEvent routes an event request to the appropriate handler.
func (r *EventRouter) RouteEvent(ctx context.Context, req *EventRequest) (*HandlerResponse, error) {
	handler, err := r.registry.get(req.Method)
	if err != nil {
		return nil, fmt.Errorf("unknown event method: %s", req.Method)
	}

	return handler(ctx, req.Data)
}

// List returns all registered event methods.
func (r *EventRouter) List() []string {
	return r.registry.list()
}

// Has returns true if a handler is registered for the given method.
func (r *EventRouter) Has(method string) bool {
	return r.registry.has(method)
}
