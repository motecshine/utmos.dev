// Package errors provides standardized error handling for UMOS IoT services.
package errors

import (
	"errors"
	"fmt"
)

// ErrorCode represents a standardized error code.
type ErrorCode int

const (
	// ErrInternal indicates an internal server error (1000).
	ErrInternal ErrorCode = 1000
	// ErrInvalidParameter indicates an invalid parameter (1001).
	ErrInvalidParameter ErrorCode = 1001
	// ErrNotFound indicates a resource was not found (1002).
	ErrNotFound ErrorCode = 1002
	// ErrAlreadyExists indicates a resource already exists (1003).
	ErrAlreadyExists ErrorCode = 1003
	// ErrUnauthorized indicates an unauthorized request (1004).
	ErrUnauthorized ErrorCode = 1004
	// ErrForbidden indicates a forbidden request (1005).
	ErrForbidden ErrorCode = 1005

	// ErrDeviceNotFound indicates a device was not found (2000).
	ErrDeviceNotFound ErrorCode = 2000
	// ErrDeviceOffline indicates a device is offline (2001).
	ErrDeviceOffline ErrorCode = 2001
	// ErrDeviceNotReady indicates a device is not ready (2002).
	ErrDeviceNotReady ErrorCode = 2002

	// ErrInvalidMessage indicates an invalid message format (3000).
	ErrInvalidMessage ErrorCode = 3000
	// ErrInvalidRoutingKey indicates an invalid routing key (3001).
	ErrInvalidRoutingKey ErrorCode = 3001
	// ErrMessageTimeout indicates a message timeout (3002).
	ErrMessageTimeout ErrorCode = 3002
	// ErrTraceContextMissing indicates a missing trace context (3003).
	ErrTraceContextMissing ErrorCode = 3003

	// ErrRabbitMQConnection indicates a RabbitMQ connection error (4000).
	ErrRabbitMQConnection ErrorCode = 4000
	// ErrDatabaseConnection indicates a database connection error (4001).
	ErrDatabaseConnection ErrorCode = 4001
	// ErrInfluxDBConnection indicates an InfluxDB connection error (4002).
	ErrInfluxDBConnection ErrorCode = 4002
)

// codeMessages maps error codes to their default messages.
var codeMessages = map[ErrorCode]string{
	ErrInternal:            "internal server error",
	ErrInvalidParameter:    "invalid parameter",
	ErrNotFound:            "resource not found",
	ErrAlreadyExists:       "resource already exists",
	ErrUnauthorized:        "unauthorized",
	ErrForbidden:           "forbidden",
	ErrDeviceNotFound:      "device not found",
	ErrDeviceOffline:       "device is offline",
	ErrDeviceNotReady:      "device is not ready",
	ErrInvalidMessage:      "invalid message format",
	ErrInvalidRoutingKey:   "invalid routing key",
	ErrMessageTimeout:      "message timeout",
	ErrTraceContextMissing: "trace context missing",
	ErrRabbitMQConnection:  "RabbitMQ connection error",
	ErrDatabaseConnection:  "database connection error",
	ErrInfluxDBConnection:  "InfluxDB connection error",
}

// Error represents a business error with code and message.
type Error struct {
	cause   error
	Message string    `json:"message"`
	Details string    `json:"details,omitempty"`
	Code    ErrorCode `json:"code"`
}

// Error implements the error interface.
func (e *Error) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("[%d] %s: %s", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// Unwrap returns the underlying cause of the error.
func (e *Error) Unwrap() error {
	return e.cause
}

// New creates a new Error with the given code and message.
func New(code ErrorCode, message string) *Error {
	if message == "" {
		message = codeMessages[code]
	}
	return &Error{
		Code:    code,
		Message: message,
	}
}

// NewWithDetails creates a new Error with the given code, message, and details.
func NewWithDetails(code ErrorCode, message, details string) *Error {
	if message == "" {
		message = codeMessages[code]
	}
	return &Error{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// Wrap wraps an existing error with a code and message.
func Wrap(err error, code ErrorCode, message string) *Error {
	if message == "" {
		message = codeMessages[code]
	}
	return &Error{
		Code:    code,
		Message: message,
		Details: err.Error(),
		cause:   err,
	}
}

// Is checks if the error matches the given error code.
func Is(err error, code ErrorCode) bool {
	var e *Error
	if errors.As(err, &e) {
		return e.Code == code
	}
	return false
}

// GetCode returns the error code from an error, or ErrInternal if not an Error.
func GetCode(err error) ErrorCode {
	var e *Error
	if errors.As(err, &e) {
		return e.Code
	}
	return ErrInternal
}
