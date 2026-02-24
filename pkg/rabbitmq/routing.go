// Package rabbitmq provides RabbitMQ client and messaging utilities.
package rabbitmq

import (
	"errors"
	"fmt"
	"strings"
)

// Routing key prefix
const RoutingKeyPrefix = "iot"

// VendorGeneric Predefined vendor constants
// Note: Vendor-specific constants (e.g., "dji", "tuya") should be defined
// in their respective adapter packages (pkg/adapter/{vendor}/), not here.
// This package only provides generic/common constants.
const (
	VendorGeneric = "generic"
)

// Predefined service constants
const (
	ServiceDevice  = "device"
	ServiceEvent   = "event"
	ServiceService = "service"
	ServiceRaw     = "raw" // Raw messages from/to protocol adapters
)

// Predefined action constants
const (
	ActionPropertyReport = "property.report"
	ActionPropertySet    = "property.set"
	ActionServiceCall    = "service.call"
	ActionServiceReply   = "service.reply"
	ActionEventReport    = "event.report"
	ActionEventNotify    = "event.notify"
	ActionDeviceOnline   = "device.online"
	ActionDeviceOffline  = "device.offline"
)

// Predefined direction constants for raw messages
const (
	DirectionUplink   = "uplink"
	DirectionDownlink = "downlink"
)

// RoutingKey represents a RabbitMQ routing key.
// Format: iot.{vendor}.{service}.{action}
type RoutingKey struct {
	Vendor  string
	Service string
	Action  string
}

// NewRoutingKey creates a new routing key.
func NewRoutingKey(vendor, service, action string) RoutingKey {
	return RoutingKey{
		Vendor:  vendor,
		Service: service,
		Action:  action,
	}
}

// String returns the routing key as a string.
// Format: iot.{vendor}.{service}.{action}
func (r RoutingKey) String() string {
	return fmt.Sprintf("%s.%s.%s.%s", RoutingKeyPrefix, r.Vendor, r.Service, r.Action)
}

// Validate validates the routing key.
func (r RoutingKey) Validate() error {
	if r.Vendor == "" {
		return errors.New("vendor is required")
	}
	if r.Service == "" {
		return errors.New("service is required")
	}
	if r.Action == "" {
		return errors.New("action is required")
	}
	return nil
}

// Parse parses a routing key string into a RoutingKey struct.
// Expected format: iot.{vendor}.{service}.{action}
// Action can contain dots (e.g., "property.report")
func Parse(key string) (*RoutingKey, error) {
	if key == "" {
		return nil, errors.New("routing key is empty")
	}

	parts := strings.Split(key, ".")
	if len(parts) < 4 {
		return nil, fmt.Errorf("invalid routing key format: expected at least 4 segments, got %d", len(parts))
	}

	if parts[0] != RoutingKeyPrefix {
		return nil, fmt.Errorf("invalid routing key prefix: expected '%s', got '%s'", RoutingKeyPrefix, parts[0])
	}

	// parts[0] = "iot"
	// parts[1] = vendor
	// parts[2] = service
	// parts[3:] = action (may contain dots)
	vendor := parts[1]
	service := parts[2]
	action := strings.Join(parts[3:], ".")

	rk := &RoutingKey{
		Vendor:  vendor,
		Service: service,
		Action:  action,
	}

	if err := rk.Validate(); err != nil {
		return nil, err
	}

	return rk, nil
}

// BuildBindingPattern creates a routing key pattern for queue binding.
// Use "*" for single word wildcard, "#" for multi-word wildcard.
func BuildBindingPattern(vendor, service, action string) string {
	v := vendor
	if v == "" {
		v = "*"
	}
	s := service
	if s == "" {
		s = "*"
	}
	a := action
	if a == "" {
		a = "#"
	}
	return fmt.Sprintf("%s.%s.%s.%s", RoutingKeyPrefix, v, s, a)
}

// RawRoutingKey represents a routing key for raw messages between gateway and adapters.
// Format: iot.raw.{vendor}.{direction}
type RawRoutingKey struct {
	Vendor    string
	Direction string // uplink or downlink
}

// NewRawRoutingKey creates a new raw routing key.
func NewRawRoutingKey(vendor, direction string) RawRoutingKey {
	return RawRoutingKey{
		Vendor:    vendor,
		Direction: direction,
	}
}

// String returns the raw routing key as a string.
// Format: iot.raw.{vendor}.{direction}
func (r RawRoutingKey) String() string {
	return fmt.Sprintf("%s.%s.%s.%s", RoutingKeyPrefix, ServiceRaw, r.Vendor, r.Direction)
}

// Validate validates the raw routing key.
func (r RawRoutingKey) Validate() error {
	if r.Vendor == "" {
		return errors.New("vendor is required")
	}
	if r.Direction != DirectionUplink && r.Direction != DirectionDownlink {
		return fmt.Errorf("direction must be '%s' or '%s'", DirectionUplink, DirectionDownlink)
	}
	return nil
}

// BuildRawBindingPattern creates a routing key pattern for raw message queue binding.
func BuildRawBindingPattern(vendor, direction string) string {
	v := vendor
	if v == "" {
		v = "*"
	}
	d := direction
	if d == "" {
		d = "*"
	}
	return fmt.Sprintf("%s.%s.%s.%s", RoutingKeyPrefix, ServiceRaw, v, d)
}
