// Package router provides service and event routing for the DJI adapter.
package router

import "fmt"

// ErrMethodNotFound is returned when a method is not registered.
var ErrMethodNotFound = fmt.Errorf("method not found")

// ErrMethodAlreadyRegistered is returned when trying to register a duplicate method.
var ErrMethodAlreadyRegistered = fmt.Errorf("method already registered")
