// Package router provides service and event routing for the DJI adapter.
package router

import (
	"context"
	"encoding/json"
	"fmt"
)

// Result codes for service responses.
const (
	ResultSuccess    = 0      // Success
	ResultParamError = 314000 // Parameter error
)

// validateAndBuildOutput validates optional JSON data against type T and builds the response output.
// Returns (output, resultCode) where resultCode is ResultParamError on invalid JSON.
func validateAndBuildOutput[T any](data json.RawMessage, method, status string) (json.RawMessage, int) {
	if len(data) > 0 {
		var req T
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, ResultParamError
		}
	}
	return json.RawMessage(fmt.Sprintf(`{"method": %q, "status": %q}`, method, status)), ResultSuccess
}

// simpleHandlerFunc creates a handler that validates data against type T.
func simpleHandlerFunc[T any](method, status string) HandlerFunc {
	return func(_ context.Context, data json.RawMessage) (*HandlerResponse, error) {
		output, result := validateAndBuildOutput[T](data, method, status)
		return &HandlerResponse{Result: result, Output: output}, nil
	}
}

// noDataHandlerFunc creates a handler for commands/events without data payload.
func noDataHandlerFunc(method, status string) HandlerFunc {
	return func(_ context.Context, _ json.RawMessage) (*HandlerResponse, error) {
		return &HandlerResponse{
			Result: ResultSuccess,
			Output: json.RawMessage(fmt.Sprintf(`{"method": %q, "status": %q}`, method, status)),
		}, nil
	}
}

// SimpleCommandHandler creates a handler for commands that just need to be accepted.
// It validates the JSON data against the provided type T and returns an accepted response.
func SimpleCommandHandler[T any](method string) ServiceHandlerFunc {
	return simpleHandlerFunc[T](method, "accepted")
}

// NoDataCommandHandler creates a handler for commands without data payload.
func NoDataCommandHandler(method string) ServiceHandlerFunc {
	return noDataHandlerFunc(method, "accepted")
}

// SimpleEventHandler creates a handler for events that just need to be acknowledged.
// It validates the JSON data against the provided type T and returns a received response.
func SimpleEventHandler[T any](method string) EventHandlerFunc {
	return simpleHandlerFunc[T](method, "received")
}

// NoDataEventHandler creates a handler for events without data payload.
func NoDataEventHandler(method string) EventHandlerFunc {
	return noDataHandlerFunc(method, "received")
}

// RegistrationError represents an error during handler registration.
type RegistrationError struct {
	Method string
	Err    error
}

func (e *RegistrationError) Error() string {
	return fmt.Sprintf("failed to register handler for method %q: %v", e.Method, e.Err)
}

func (e *RegistrationError) Unwrap() error {
	return e.Err
}

// registerAll registers multiple handlers into a handlerRegistry and collects any errors.
// Returns nil if all registrations succeed, otherwise returns the first error.
func registerAll[T any](reg *handlerRegistry[T], handlers map[string]T) error {
	for method, handler := range handlers {
		if err := reg.register(method, handler, ErrMethodAlreadyRegistered); err != nil {
			return &RegistrationError{Method: method, Err: err}
		}
	}
	return nil
}

// RegisterHandlers registers multiple service handlers and collects any errors.
// Returns nil if all registrations succeed, otherwise returns the first error.
func RegisterHandlers(r *ServiceRouter, handlers map[string]ServiceHandlerFunc) error {
	return registerAll(&r.registry, handlers)
}

// RegisterEventHandlers registers multiple event handlers and collects any errors.
// Returns nil if all registrations succeed, otherwise returns the first error.
func RegisterEventHandlers(r *EventRouter, handlers map[string]EventHandlerFunc) error {
	return registerAll(&r.registry, handlers)
}
