package handler

import (
	dji "github.com/utmos/utmos/pkg/adapter/dji"
	"github.com/utmos/utmos/pkg/adapter/dji/router"
)

// EventHandler handles Event messages.
type EventHandler struct {
	baseHandler
	router *router.EventRouter
}

// NewEventHandler creates a new Event handler.
func NewEventHandler(r *router.EventRouter) *EventHandler {
	return &EventHandler{
		baseHandler: newBaseHandler(dji.TopicTypeEvents, eventHandlerConfig, "event"),
		router:      r,
	}
}

// GetRouter returns the event router.
func (h *EventHandler) GetRouter() *router.EventRouter {
	return h.router
}

// Ensure EventHandler implements Handler interface.
var _ Handler = (*EventHandler)(nil)

// RequestHandler handles device-initiated request messages.
type RequestHandler struct {
	baseHandler
	router *router.ServiceRouter
}

// NewRequestHandler creates a new Request handler.
func NewRequestHandler(r *router.ServiceRouter) *RequestHandler {
	return &RequestHandler{
		baseHandler: newBaseHandler(dji.TopicTypeRequests, requestHandlerConfig, "request"),
		router:      r,
	}
}

// GetRouter returns the service router.
func (h *RequestHandler) GetRouter() *router.ServiceRouter {
	return h.router
}

// Ensure RequestHandler implements Handler interface.
var _ Handler = (*RequestHandler)(nil)
