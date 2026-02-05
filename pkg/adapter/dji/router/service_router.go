package router

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
)

// ServiceHandlerFunc is the function signature for service handlers.
type ServiceHandlerFunc func(ctx context.Context, data json.RawMessage) (*ServiceResponse, error)

// ServiceRouter routes service method calls to appropriate handlers.
type ServiceRouter struct {
	mu       sync.RWMutex
	handlers map[string]ServiceHandlerFunc
}

// NewServiceRouter creates a new service router.
func NewServiceRouter() *ServiceRouter {
	return &ServiceRouter{
		handlers: make(map[string]ServiceHandlerFunc),
	}
}

// ServiceRequest represents a service call request.
type ServiceRequest struct {
	Method string          `json:"method"`
	Data   json.RawMessage `json:"data,omitempty"`
}

// ServiceResponse represents a service call response.
type ServiceResponse struct {
	Result int             `json:"result"`
	Output json.RawMessage `json:"output,omitempty"`
}

// RegisterServiceHandler registers a service handler.
func (r *ServiceRouter) RegisterServiceHandler(method string, handler ServiceHandlerFunc) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.handlers[method]; exists {
		return fmt.Errorf("%w: %s", ErrMethodAlreadyRegistered, method)
	}

	r.handlers[method] = handler
	return nil
}

// RouteService routes a service request and returns a response.
func (r *ServiceRouter) RouteService(ctx context.Context, req *ServiceRequest) (*ServiceResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("nil service request")
	}

	r.mu.RLock()
	handler, exists := r.handlers[req.Method]
	r.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("%w: %s", ErrMethodNotFound, req.Method)
	}

	return handler(ctx, req.Data)
}

// List returns all registered methods.
func (r *ServiceRouter) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	methods := make([]string, 0, len(r.handlers))
	for method := range r.handlers {
		methods = append(methods, method)
	}
	return methods
}

// Has checks if a method is registered.
func (r *ServiceRouter) Has(method string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.handlers[method]
	return exists
}
