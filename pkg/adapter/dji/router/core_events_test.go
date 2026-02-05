package router

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterCoreEvents(t *testing.T) {
	r := NewEventRouter()
	err := RegisterCoreEvents(r)
	require.NoError(t, err)

	// Verify all core events are registered
	expectedMethods := []string{
		MethodHMS,
		MethodDeviceExitHomingNotify,
		MethodDeviceTempNtfyNeedClear,
		MethodControlSourceChange,
		MethodFlyToPointProgress,
		MethodTakeoffToPointProgress,
	}

	for _, method := range expectedMethods {
		assert.True(t, r.Has(method), "method %s should be registered", method)
	}
}

func TestCoreEvents_HMS(t *testing.T) {
	r := NewEventRouter()
	err := RegisterCoreEvents(r)
	require.NoError(t, err)

	tests := []struct {
		name string
		data string
	}{
		{
			name: "valid HMS data",
			data: `{
				"list": [
					{
						"code": "0x16100001",
						"level": 0,
						"module": 3,
						"in_the_sky": 0
					}
				]
			}`,
		},
		{
			name: "multiple HMS items",
			data: `{
				"list": [
					{"code": "0x16100001", "level": 0, "module": 3, "in_the_sky": 0},
					{"code": "0x16100002", "level": 1, "module": 3, "in_the_sky": 1}
				]
			}`,
		},
		{
			name: "HMS with args",
			data: `{
				"list": [
					{
						"code": "0x16100001",
						"level": 0,
						"module": 3,
						"in_the_sky": 0,
						"args": {"component_index": 0, "sensor_index": 0}
					}
				]
			}`,
		},
		{
			name: "empty list",
			data: `{"list": []}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &EventRequest{
				Method: MethodHMS,
				Data:   json.RawMessage(tt.data),
			}
			resp, err := r.RouteEvent(context.Background(), req)
			require.NoError(t, err)
			assert.Equal(t, 0, resp.Result)
		})
	}
}

func TestCoreEvents_DeviceExitHomingNotify(t *testing.T) {
	r := NewEventRouter()
	err := RegisterCoreEvents(r)
	require.NoError(t, err)

	req := &EventRequest{
		Method: MethodDeviceExitHomingNotify,
		Data:   json.RawMessage(`{}`),
	}

	resp, err := r.RouteEvent(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Result)
}

func TestCoreEvents_DeviceTempNtfyNeedClear(t *testing.T) {
	r := NewEventRouter()
	err := RegisterCoreEvents(r)
	require.NoError(t, err)

	req := &EventRequest{
		Method: MethodDeviceTempNtfyNeedClear,
		Data:   json.RawMessage(`{}`),
	}

	resp, err := r.RouteEvent(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Result)
}

func TestCoreEvents_ControlSourceChange(t *testing.T) {
	r := NewEventRouter()
	err := RegisterCoreEvents(r)
	require.NoError(t, err)

	req := &EventRequest{
		Method: MethodControlSourceChange,
		Data:   json.RawMessage(`{}`),
	}

	resp, err := r.RouteEvent(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Result)
}

func TestCoreEvents_FlyToPointProgress(t *testing.T) {
	r := NewEventRouter()
	err := RegisterCoreEvents(r)
	require.NoError(t, err)

	req := &EventRequest{
		Method: MethodFlyToPointProgress,
		Data:   json.RawMessage(`{}`),
	}

	resp, err := r.RouteEvent(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Result)
}

func TestCoreEvents_TakeoffToPointProgress(t *testing.T) {
	r := NewEventRouter()
	err := RegisterCoreEvents(r)
	require.NoError(t, err)

	req := &EventRequest{
		Method: MethodTakeoffToPointProgress,
		Data:   json.RawMessage(`{}`),
	}

	resp, err := r.RouteEvent(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Result)
}

func TestCoreEvents_InvalidData(t *testing.T) {
	r := NewEventRouter()
	err := RegisterCoreEvents(r)
	require.NoError(t, err)

	req := &EventRequest{
		Method: MethodHMS,
		Data:   json.RawMessage(`{invalid json}`),
	}

	resp, err := r.RouteEvent(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, 314000, resp.Result) // Parameter error
}
