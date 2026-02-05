package router

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
)

// EventHandlerFunc is the function signature for event handlers.
type EventHandlerFunc func(ctx context.Context, data json.RawMessage) (*EventResponse, error)

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

// EventResponse represents an event response.
type EventResponse struct {
	Result int             `json:"result"`
	Output json.RawMessage `json:"output,omitempty"`
}

// EventRouter routes events to their handlers.
type EventRouter struct {
	mu       sync.RWMutex
	handlers map[string]EventHandlerFunc
}

// NewEventRouter creates a new EventRouter.
func NewEventRouter() *EventRouter {
	return &EventRouter{
		handlers: make(map[string]EventHandlerFunc),
	}
}

// RegisterEventHandler registers a handler for an event method.
func (r *EventRouter) RegisterEventHandler(method string, handler EventHandlerFunc) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.handlers[method]; exists {
		return fmt.Errorf("event handler for method %s already registered", method)
	}

	r.handlers[method] = handler
	return nil
}

// RouteEvent routes an event request to the appropriate handler.
func (r *EventRouter) RouteEvent(ctx context.Context, req *EventRequest) (*EventResponse, error) {
	r.mu.RLock()
	handler, exists := r.handlers[req.Method]
	r.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("unknown event method: %s", req.Method)
	}

	return handler(ctx, req.Data)
}

// List returns all registered event methods.
func (r *EventRouter) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	methods := make([]string, 0, len(r.handlers))
	for method := range r.handlers {
		methods = append(methods, method)
	}
	return methods
}

// Has returns true if a handler is registered for the given method.
func (r *EventRouter) Has(method string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.handlers[method]
	return exists
}
