package router

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterDRCCommands(t *testing.T) {
	r := NewServiceRouter()
	err := RegisterDRCCommands(r)
	require.NoError(t, err)

	// Verify all DRC commands are registered
	expectedMethods := []string{
		MethodDRCModeEnter,
		MethodDRCModeExit,
		MethodDroneControl,
		MethodDroneEmergencyStop,
		MethodHeart,
	}

	for _, method := range expectedMethods {
		assert.True(t, r.Has(method), "method %s should be registered", method)
	}
}

func TestDRCCommands_DRCModeEnter(t *testing.T) {
	r := NewServiceRouter()
	err := RegisterDRCCommands(r)
	require.NoError(t, err)

	req := &ServiceRequest{
		Method: MethodDRCModeEnter,
		Data: json.RawMessage(`{
			"mqtt_broker": {
				"address": "mqtt://broker.example.com:1883",
				"client_id": "drc-client-001",
				"username": "user",
				"password": "pass",
				"expire_time": 3600,
				"enable_tls": false
			},
			"osd_frequency": 10,
			"hsi_frequency": 10
		}`),
	}

	resp, err := r.RouteService(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Result)
}

func TestDRCCommands_DRCModeExit(t *testing.T) {
	r := NewServiceRouter()
	err := RegisterDRCCommands(r)
	require.NoError(t, err)

	req := &ServiceRequest{
		Method: MethodDRCModeExit,
	}

	resp, err := r.RouteService(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Result)
}

func TestDRCCommands_DroneControl(t *testing.T) {
	r := NewServiceRouter()
	err := RegisterDRCCommands(r)
	require.NoError(t, err)

	tests := []struct {
		name string
		data string
	}{
		{
			name: "basic control",
			data: `{"seq": 1, "x": 0.5, "y": 0.5, "h": 0.0, "w": 0.0}`,
		},
		{
			name: "full control",
			data: `{"seq": 100, "x": 0.8, "y": -0.3, "h": 0.2, "w": 0.1}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &ServiceRequest{
				Method: MethodDroneControl,
				Data:   json.RawMessage(tt.data),
			}
			resp, err := r.RouteService(context.Background(), req)
			require.NoError(t, err)
			assert.Equal(t, 0, resp.Result)
		})
	}
}

func TestDRCCommands_DroneEmergencyStop(t *testing.T) {
	r := NewServiceRouter()
	err := RegisterDRCCommands(r)
	require.NoError(t, err)

	req := &ServiceRequest{
		Method: MethodDroneEmergencyStop,
	}

	resp, err := r.RouteService(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Result)
}

func TestDRCCommands_Heart(t *testing.T) {
	r := NewServiceRouter()
	err := RegisterDRCCommands(r)
	require.NoError(t, err)

	req := &ServiceRequest{
		Method: MethodHeart,
		Data:   json.RawMessage(`{"seq": 1, "timestamp": 1706000000000}`),
	}

	resp, err := r.RouteService(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Result)
}

func TestDRCCommands_InvalidData(t *testing.T) {
	r := NewServiceRouter()
	err := RegisterDRCCommands(r)
	require.NoError(t, err)

	req := &ServiceRequest{
		Method: MethodDroneControl,
		Data:   json.RawMessage(`{invalid json}`),
	}

	resp, err := r.RouteService(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, 314000, resp.Result) // Parameter error
}
