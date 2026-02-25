package router

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterDeviceCommands(t *testing.T) {
	r := NewServiceRouter()
	err := RegisterDeviceCommands(r)
	require.NoError(t, err)

	// Verify all device commands are registered
	expectedMethods := []string{
		MethodCoverOpen,
		MethodCoverClose,
		MethodCoverForceClose,
		MethodDroneOpen,
		MethodDroneClose,
		MethodChargeOpen,
		MethodChargeClose,
		MethodDeviceReboot,
		MethodDeviceFormat,
		MethodDroneFormat,
		MethodDebugModeOpen,
		MethodDebugModeClose,
		MethodBatteryMaintenanceSwitch,
		MethodAirConditionerModeSwitch,
		MethodAlarmStateSwitch,
		MethodSDRWorkmodeSwitch,
	}

	for _, method := range expectedMethods {
		assert.True(t, r.Has(method), "method %s should be registered", method)
	}
}

func TestDeviceCommands_CoverOpen(t *testing.T) {
	r := NewServiceRouter()
	err := RegisterDeviceCommands(r)
	require.NoError(t, err)

	req := &ServiceRequest{Method: MethodCoverOpen}
	resp, err := r.RouteService(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Result)
}

func TestDeviceCommands_CoverClose(t *testing.T) {
	r := NewServiceRouter()
	err := RegisterDeviceCommands(r)
	require.NoError(t, err)

	req := &ServiceRequest{Method: MethodCoverClose}
	resp, err := r.RouteService(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Result)
}

func TestDeviceCommands_DroneOpen(t *testing.T) {
	r := NewServiceRouter()
	err := RegisterDeviceCommands(r)
	require.NoError(t, err)

	req := &ServiceRequest{Method: MethodDroneOpen}
	resp, err := r.RouteService(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Result)
}

func TestDeviceCommands_DroneClose(t *testing.T) {
	r := NewServiceRouter()
	err := RegisterDeviceCommands(r)
	require.NoError(t, err)

	req := &ServiceRequest{Method: MethodDroneClose}
	resp, err := r.RouteService(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Result)
}

func TestDeviceCommands_DeviceReboot(t *testing.T) {
	r := NewServiceRouter()
	err := RegisterDeviceCommands(r)
	require.NoError(t, err)

	req := &ServiceRequest{Method: MethodDeviceReboot}
	resp, err := r.RouteService(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Result)
}

func TestDeviceCommands_BatteryMaintenanceSwitch(t *testing.T) {
	r := NewServiceRouter()
	err := RegisterDeviceCommands(r)
	require.NoError(t, err)

	tests := []struct {
		name string
		data string
	}{
		{
			name: "enable",
			data: `{"action": 1}`,
		},
		{
			name: "disable",
			data: `{"action": 0}`,
		},
		{
			name: "no data",
			data: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var data json.RawMessage
			if tt.data != "" {
				data = json.RawMessage(tt.data)
			}

			req := &ServiceRequest{
				Method: MethodBatteryMaintenanceSwitch,
				Data:   data,
			}
			resp, err := r.RouteService(context.Background(), req)
			require.NoError(t, err)
			assert.Equal(t, 0, resp.Result)
		})
	}
}

func TestDeviceCommands_AirConditionerModeSwitch(t *testing.T) {
	r := NewServiceRouter()
	err := RegisterDeviceCommands(r)
	require.NoError(t, err)

	req := &ServiceRequest{
		Method: MethodAirConditionerModeSwitch,
		Data:   json.RawMessage(`{"action": 1}`),
	}
	resp, err := r.RouteService(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Result)
}

func TestDeviceCommands_AlarmStateSwitch(t *testing.T) {
	r := NewServiceRouter()
	err := RegisterDeviceCommands(r)
	require.NoError(t, err)

	req := &ServiceRequest{
		Method: MethodAlarmStateSwitch,
		Data:   json.RawMessage(`{"action": 1}`),
	}
	resp, err := r.RouteService(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Result)
}

func TestDeviceCommands_SDRWorkmodeSwitch(t *testing.T) {
	r := NewServiceRouter()
	err := RegisterDeviceCommands(r)
	require.NoError(t, err)

	req := &ServiceRequest{
		Method: MethodSDRWorkmodeSwitch,
		Data:   json.RawMessage(`{"link_workmode": 0}`),
	}
	resp, err := r.RouteService(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Result)
}

func TestDeviceCommands_InvalidData(t *testing.T) {
	r := NewServiceRouter()
	err := RegisterDeviceCommands(r)
	require.NoError(t, err)

	req := &ServiceRequest{
		Method: MethodBatteryMaintenanceSwitch,
		Data:   json.RawMessage(`{invalid json}`),
	}
	resp, err := r.RouteService(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, 314000, resp.Result) // Parameter error
}
