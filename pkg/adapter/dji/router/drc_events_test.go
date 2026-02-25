package router

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterDRCEvents(t *testing.T) {
	r := NewEventRouter()
	err := RegisterDRCEvents(r)
	require.NoError(t, err)

	// Verify all DRC events are registered
	expectedMethods := []string{
		MethodJoystickInvalidNotify,
		MethodDRCStatusNotify,
	}

	for _, method := range expectedMethods {
		assert.True(t, r.Has(method), "method %s should be registered", method)
	}
}

func TestDRCEvents_JoystickInvalidNotify(t *testing.T) {
	r := NewEventRouter()
	err := RegisterDRCEvents(r)
	require.NoError(t, err)

	req := &EventRequest{
		Method: MethodJoystickInvalidNotify,
		Data:   json.RawMessage(`{"reason": "timeout"}`),
	}

	resp, err := r.RouteEvent(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Result)
}

func TestDRCEvents_DRCStatusNotify(t *testing.T) {
	r := NewEventRouter()
	err := RegisterDRCEvents(r)
	require.NoError(t, err)

	tests := []struct {
		name string
		data string
	}{
		{
			name: "connected",
			data: `{"status": 1}`,
		},
		{
			name: "disconnected",
			data: `{"status": 0}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &EventRequest{
				Method: MethodDRCStatusNotify,
				Data:   json.RawMessage(tt.data),
			}
			resp, err := r.RouteEvent(context.Background(), req)
			require.NoError(t, err)
			assert.Equal(t, 0, resp.Result)
		})
	}
}
