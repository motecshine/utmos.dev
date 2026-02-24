package router

import (
	"context"
	"encoding/json"
	"fmt"
)

// ServiceHandlerFunc is an alias for HandlerFunc used in service handler registration.
type ServiceHandlerFunc = HandlerFunc

// ServiceRouter routes service method calls to appropriate handlers.
type ServiceRouter struct {
	registry handlerRegistry[HandlerFunc]
}

// NewServiceRouter creates a new service router.
func NewServiceRouter() *ServiceRouter {
	return &ServiceRouter{
		registry: newHandlerRegistry[HandlerFunc](),
	}
}

// ServiceRequest represents a service call request.
type ServiceRequest struct {
	Method string          `json:"method"`
	Data   json.RawMessage `json:"data,omitempty"`
}

// ServiceResponse is an alias for HandlerResponse used in service handlers.
type ServiceResponse = HandlerResponse

// RegisterServiceHandler registers a service handler.
func (r *ServiceRouter) RegisterServiceHandler(method string, handler ServiceHandlerFunc) error {
	return r.registry.register(method, handler, ErrMethodAlreadyRegistered)
}

// RouteService routes a service request and returns a response.
func (r *ServiceRouter) RouteService(ctx context.Context, req *ServiceRequest) (*HandlerResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("nil service request")
	}

	handler, err := r.registry.get(req.Method)
	if err != nil {
		return nil, err
	}

	return handler(ctx, req.Data)
}

// List returns all registered methods.
func (r *ServiceRouter) List() []string {
	return r.registry.list()
}

// Has checks if a method is registered.
func (r *ServiceRouter) Has(method string) bool {
	return r.registry.has(method)
}
