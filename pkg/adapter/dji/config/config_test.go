package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestServiceCallTimeout(t *testing.T) {
	assert.Equal(t, 30*time.Second, ServiceCallTimeout)
}

func TestDRCHeartbeatTimeout(t *testing.T) {
	assert.Equal(t, 3*time.Second, DRCHeartbeatTimeout)
}

func TestDefaultUnknownDevicePolicy(t *testing.T) {
	assert.Equal(t, PolicyDiscard, DefaultUnknownDevicePolicy)
}

func TestUnknownDevicePolicyValues(t *testing.T) {
	assert.Equal(t, UnknownDevicePolicy("discard"), PolicyDiscard)
	assert.Equal(t, UnknownDevicePolicy("forward"), PolicyForward)
	assert.Equal(t, UnknownDevicePolicy("dlq"), PolicyDLQ)
}
