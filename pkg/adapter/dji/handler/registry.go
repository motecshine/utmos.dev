package handler

import (
	"fmt"
	"sync"

	dji "github.com/utmos/utmos/pkg/adapter/dji"
)

// ErrHandlerNotFound is returned when no handler is registered for a topic type.
var ErrHandlerNotFound = fmt.Errorf("handler not found")

// ErrHandlerAlreadyRegistered is returned when a handler is already registered for a topic type.
var ErrHandlerAlreadyRegistered = fmt.Errorf("handler already registered")

// Registry manages handlers for different topic types.
type Registry struct {
	mu       sync.RWMutex
	handlers map[dji.TopicType]Handler
}

// NewRegistry creates a new handler registry.
func NewRegistry() *Registry {
	return &Registry{
		handlers: make(map[dji.TopicType]Handler),
	}
}

// Register registers a handler for a topic type.
func (r *Registry) Register(handler Handler) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	topicType := handler.GetTopicType()
	if _, exists := r.handlers[topicType]; exists {
		return fmt.Errorf("%w: %s", ErrHandlerAlreadyRegistered, topicType)
	}

	r.handlers[topicType] = handler
	return nil
}

// Get returns the handler for a topic type.
func (r *Registry) Get(topicType dji.TopicType) (Handler, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	handler, exists := r.handlers[topicType]
	if !exists {
		return nil, fmt.Errorf("%w: %s", ErrHandlerNotFound, topicType)
	}

	return handler, nil
}

// Has checks if a handler is registered for a topic type.
func (r *Registry) Has(topicType dji.TopicType) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.handlers[topicType]
	return exists
}

// List returns all registered topic types.
func (r *Registry) List() []dji.TopicType {
	r.mu.RLock()
	defer r.mu.RUnlock()

	types := make([]dji.TopicType, 0, len(r.handlers))
	for t := range r.handlers {
		types = append(types, t)
	}
	return types
}

// MustRegister registers a handler and panics on error.
func (r *Registry) MustRegister(handler Handler) {
	if err := r.Register(handler); err != nil {
		panic(err)
	}
}
