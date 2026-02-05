package router

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterConfigCommands(t *testing.T) {
	r := NewServiceRouter()
	err := RegisterConfigCommands(r)
	require.NoError(t, err)

	// Verify all config commands are registered
	expectedMethods := []string{
		MethodConfigGet,
		MethodConfigSet,
	}

	for _, method := range expectedMethods {
		assert.True(t, r.Has(method), "method %s should be registered", method)
	}
}

func TestConfigCommands_ConfigGet(t *testing.T) {
	r := NewServiceRouter()
	err := RegisterConfigCommands(r)
	require.NoError(t, err)

	tests := []struct {
		name string
		data string
	}{
		{
			name: "get single config",
			data: `{"config_type": "basic_device_info"}`,
		},
		{
			name: "get multiple configs",
			data: `{"config_type": "basic_device_info,firmware_version"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &ServiceRequest{
				Method: MethodConfigGet,
				Data:   json.RawMessage(tt.data),
			}
			resp, err := r.RouteService(context.Background(), req)
			require.NoError(t, err)
			assert.Equal(t, 0, resp.Result)
		})
	}
}

func TestConfigCommands_ConfigSet(t *testing.T) {
	r := NewServiceRouter()
	err := RegisterConfigCommands(r)
	require.NoError(t, err)

	tests := []struct {
		name string
		data string
	}{
		{
			name: "set basic config",
			data: `{"config_type": "basic_device_info", "config_value": {"name": "Dock-001"}}`,
		},
		{
			name: "set network config",
			data: `{"config_type": "network", "config_value": {"ssid": "test-wifi"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &ServiceRequest{
				Method: MethodConfigSet,
				Data:   json.RawMessage(tt.data),
			}
			resp, err := r.RouteService(context.Background(), req)
			require.NoError(t, err)
			assert.Equal(t, 0, resp.Result)
		})
	}
}

func TestConfigCommands_InvalidData(t *testing.T) {
	r := NewServiceRouter()
	err := RegisterConfigCommands(r)
	require.NoError(t, err)

	req := &ServiceRequest{
		Method: MethodConfigGet,
		Data:   json.RawMessage(`{invalid json}`),
	}

	resp, err := r.RouteService(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, 314000, resp.Result) // Parameter error
}
