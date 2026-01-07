package errors

import "fmt"

// ErrorCode represents a specific error code
type ErrorCode string

const (
	// ErrCodeNotFound indicates resource not found
	ErrCodeNotFound ErrorCode = "NOT_FOUND"
	// ErrCodeInvalidInput indicates invalid input
	ErrCodeInvalidInput ErrorCode = "INVALID_INPUT"
	// ErrCodeInternal indicates internal server error
	ErrCodeInternal ErrorCode = "INTERNAL_ERROR"
	// ErrCodeUnauthorized indicates unauthorized access
	ErrCodeUnauthorized ErrorCode = "UNAUTHORIZED"
	// ErrCodeForbidden indicates forbidden access
	ErrCodeForbidden ErrorCode = "FORBIDDEN"
	// ErrCodeConflict indicates resource conflict
	ErrCodeConflict ErrorCode = "CONFLICT"
)

// AppError represents an application error with code and message
type AppError struct {
	Code    ErrorCode
	Message string
	Err     error
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap returns the underlying error
func (e *AppError) Unwrap() error {
	return e.Err
}

// New creates a new AppError
func New(code ErrorCode, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// Wrap wraps an existing error with an AppError
func Wrap(code ErrorCode, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// IsNotFound checks if error is not found
func IsNotFound(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == ErrCodeNotFound
	}
	return false
}

// IsInvalidInput checks if error is invalid input
func IsInvalidInput(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == ErrCodeInvalidInput
	}
	return false
}

