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
		MethodFileUploadCallback,
		MethodFileUploadProgress,
		MethodHighestPriorityUpload,
		MethodDeviceExitHomingNotify,
		MethodDeviceTempNtfyNeedClear,
		MethodFlighttaskProgress,
		MethodFlighttaskReady,
		MethodReturnHomeInfo,
		MethodControlSourceChange,
		MethodFlyToPointProgress,
		MethodTakeoffToPointProgress,
		MethodDRCStatusNotify,
		MethodJoystickInvalidNotify,
		MethodOTAProgress,
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

func TestCoreEvents_FileUploadCallback(t *testing.T) {
	r := NewEventRouter()
	err := RegisterCoreEvents(r)
	require.NoError(t, err)

	req := &EventRequest{
		Method: MethodFileUploadCallback,
		Data: json.RawMessage(`{
			"file": {
				"path": "/media/DJI_001.jpg",
				"name": "DJI_001.jpg",
				"size": 1024000,
				"fingerprint": "abc123"
			}
		}`),
	}

	resp, err := r.RouteEvent(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Result)
}

func TestCoreEvents_FlighttaskProgress(t *testing.T) {
	r := NewEventRouter()
	err := RegisterCoreEvents(r)
	require.NoError(t, err)

	req := &EventRequest{
		Method: MethodFlighttaskProgress,
		Data: json.RawMessage(`{
			"flight_id": "flight-001",
			"status": "executing",
			"progress": 50
		}`),
	}

	resp, err := r.RouteEvent(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Result)
}

func TestCoreEvents_FlighttaskReady(t *testing.T) {
	r := NewEventRouter()
	err := RegisterCoreEvents(r)
	require.NoError(t, err)

	needReply := 1
	req := &EventRequest{
		Method:    MethodFlighttaskReady,
		NeedReply: &needReply,
		Data:      json.RawMessage(`{"flight_id": "flight-001"}`),
	}

	resp, err := r.RouteEvent(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Result)
}

func TestCoreEvents_ReturnHomeInfo(t *testing.T) {
	r := NewEventRouter()
	err := RegisterCoreEvents(r)
	require.NoError(t, err)

	req := &EventRequest{
		Method: MethodReturnHomeInfo,
		Data: json.RawMessage(`{
			"flight_id": "flight-001",
			"last_point_type": 1,
			"planned_path_points": [
				{"latitude": 22.5, "longitude": 113.9, "height": 100}
			]
		}`),
	}

	resp, err := r.RouteEvent(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Result)
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

func TestCoreEvents_DRCStatusNotify(t *testing.T) {
	r := NewEventRouter()
	err := RegisterCoreEvents(r)
	require.NoError(t, err)

	req := &EventRequest{
		Method: MethodDRCStatusNotify,
		Data:   json.RawMessage(`{"status": 1}`),
	}

	resp, err := r.RouteEvent(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Result)
}

func TestCoreEvents_JoystickInvalidNotify(t *testing.T) {
	r := NewEventRouter()
	err := RegisterCoreEvents(r)
	require.NoError(t, err)

	req := &EventRequest{
		Method: MethodJoystickInvalidNotify,
		Data:   json.RawMessage(`{}`),
	}

	resp, err := r.RouteEvent(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Result)
}

func TestCoreEvents_OTAProgress(t *testing.T) {
	r := NewEventRouter()
	err := RegisterCoreEvents(r)
	require.NoError(t, err)

	req := &EventRequest{
		Method: MethodOTAProgress,
		Data:   json.RawMessage(`{"progress": 50}`),
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
