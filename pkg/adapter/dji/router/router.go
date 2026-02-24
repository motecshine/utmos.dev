// Package router provides service and event routing for the DJI adapter.
package router

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
)

// HandlerResponse is the unified response type for both service and event handlers.
type HandlerResponse struct {
	Result int             `json:"result"`
	Output json.RawMessage `json:"output,omitempty"`
}

// HandlerFunc is the unified function signature for service and event handlers.
type HandlerFunc func(ctx context.Context, data json.RawMessage) (*HandlerResponse, error)

// ErrMethodNotFound is returned when a method is not registered.
var ErrMethodNotFound = fmt.Errorf("method not found")

// ErrMethodAlreadyRegistered is returned when trying to register a duplicate method.
var ErrMethodAlreadyRegistered = fmt.Errorf("method already registered")

// handlerRegistry is a generic, thread-safe registry for method handlers.
// It is embedded by ServiceRouter and EventRouter to avoid duplicating
// List, Has, register, and get logic.
type handlerRegistry[T any] struct {
	mu       sync.RWMutex
	handlers map[string]T
}

// newHandlerRegistry creates a new handlerRegistry.
func newHandlerRegistry[T any]() handlerRegistry[T] {
	return handlerRegistry[T]{
		handlers: make(map[string]T),
	}
}

// register adds a handler for the given method. Returns an error wrapping
// sentinelErr if the method is already registered.
func (r *handlerRegistry[T]) register(method string, handler T, sentinelErr error) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.handlers[method]; exists {
		return fmt.Errorf("%w: %s", sentinelErr, method)
	}

	r.handlers[method] = handler
	return nil
}

// get retrieves the handler for the given method.
func (r *handlerRegistry[T]) get(method string) (T, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	handler, exists := r.handlers[method]
	if !exists {
		var zero T
		return zero, fmt.Errorf("%w: %s", ErrMethodNotFound, method)
	}

	return handler, nil
}

// list returns all registered method names.
func (r *handlerRegistry[T]) list() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	methods := make([]string, 0, len(r.handlers))
	for method := range r.handlers {
		methods = append(methods, method)
	}
	return methods
}

// has checks if a method is registered.
func (r *handlerRegistry[T]) has(method string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.handlers[method]
	return exists
}
