// Package config provides configuration constants for the DJI adapter.
package config

import "time"

// Service call configuration.
const (
	// ServiceCallTimeout is the default timeout for service calls.
	ServiceCallTimeout = 30 * time.Second

	// DRCHeartbeatTimeout is the timeout for DRC heartbeat.
	DRCHeartbeatTimeout = 3 * time.Second
)

// UnknownDevicePolicy defines how to handle messages from unknown devices.
type UnknownDevicePolicy string

const (
	// PolicyDiscard discards messages from unknown devices.
	PolicyDiscard UnknownDevicePolicy = "discard"

	// PolicyForward forwards messages from unknown devices.
	PolicyForward UnknownDevicePolicy = "forward"

	// PolicyDLQ sends messages from unknown devices to dead letter queue.
	PolicyDLQ UnknownDevicePolicy = "dlq"
)

// DefaultUnknownDevicePolicy is the default policy for unknown devices.
const DefaultUnknownDevicePolicy = PolicyDiscard
