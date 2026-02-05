// Package router provides service and event routing for the DJI adapter.
package router

import (
	"context"
	"encoding/json"
	"fmt"
)

// SimpleCommandHandler creates a handler for commands that just need to be accepted.
// It validates the JSON data against the provided type T and returns an accepted response.
func SimpleCommandHandler[T any](method string) ServiceHandlerFunc {
	return func(_ context.Context, data json.RawMessage) (*ServiceResponse, error) {
		if len(data) > 0 {
			var req T
			if err := json.Unmarshal(data, &req); err != nil {
				return &ServiceResponse{Result: ResultParamError}, nil
			}
		}

		return &ServiceResponse{
			Result: ResultSuccess,
			Output: json.RawMessage(fmt.Sprintf(`{"method": %q, "status": "accepted"}`, method)),
		}, nil
	}
}

// NoDataCommandHandler creates a handler for commands without data payload.
func NoDataCommandHandler(method string) ServiceHandlerFunc {
	return func(_ context.Context, _ json.RawMessage) (*ServiceResponse, error) {
		return &ServiceResponse{
			Result: ResultSuccess,
			Output: json.RawMessage(fmt.Sprintf(`{"method": %q, "status": "accepted"}`, method)),
		}, nil
	}
}

// Result codes for service responses.
const (
	ResultSuccess    = 0      // Success
	ResultParamError = 314000 // Parameter error
)

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

// RegisterHandlers registers multiple handlers and collects any errors.
// Returns nil if all registrations succeed, otherwise returns the first error.
func RegisterHandlers(r *ServiceRouter, handlers map[string]ServiceHandlerFunc) error {
	for method, handler := range handlers {
		if err := r.RegisterServiceHandler(method, handler); err != nil {
			return &RegistrationError{Method: method, Err: err}
		}
	}
	return nil
}

// SimpleEventHandler creates a handler for events that just need to be acknowledged.
// It validates the JSON data against the provided type T and returns a received response.
func SimpleEventHandler[T any](method string) EventHandlerFunc {
	return func(_ context.Context, data json.RawMessage) (*EventResponse, error) {
		if len(data) > 0 {
			var req T
			if err := json.Unmarshal(data, &req); err != nil {
				return &EventResponse{Result: ResultParamError}, nil
			}
		}

		return &EventResponse{
			Result: ResultSuccess,
			Output: json.RawMessage(fmt.Sprintf(`{"method": %q, "status": "received"}`, method)),
		}, nil
	}
}

// NoDataEventHandler creates a handler for events without data payload.
func NoDataEventHandler(method string) EventHandlerFunc {
	return func(_ context.Context, _ json.RawMessage) (*EventResponse, error) {
		return &EventResponse{
			Result: ResultSuccess,
			Output: json.RawMessage(fmt.Sprintf(`{"method": %q, "status": "received"}`, method)),
		}, nil
	}
}

// RegisterEventHandlers registers multiple event handlers and collects any errors.
// Returns nil if all registrations succeed, otherwise returns the first error.
func RegisterEventHandlers(r *EventRouter, handlers map[string]EventHandlerFunc) error {
	for method, handler := range handlers {
		if err := r.RegisterEventHandler(method, handler); err != nil {
			return &RegistrationError{Method: method, Err: err}
		}
	}
	return nil
}
