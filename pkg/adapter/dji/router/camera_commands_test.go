package router

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterCameraCommands(t *testing.T) {
	r := NewServiceRouter()
	err := RegisterCameraCommands(r)
	require.NoError(t, err)

	// Verify all camera commands are registered
	expectedMethods := []string{
		MethodCameraModeSwitch,
		MethodCameraPhotoTake,
		MethodCameraRecordingStart,
		MethodCameraRecordingStop,
		MethodCameraAim,
		MethodCameraFocalLengthSet,
		MethodGimbalReset,
		MethodCameraPointFocusAction,
		MethodCameraScreenSplit,
		MethodIRMeteringPoint,
		MethodIRMeteringArea,
	}

	for _, method := range expectedMethods {
		assert.True(t, r.Has(method), "method %s should be registered", method)
	}
}

func TestCameraCommands_CameraModeSwitch(t *testing.T) {
	r := NewServiceRouter()
	err := RegisterCameraCommands(r)
	require.NoError(t, err)

	tests := []struct {
		name string
		data string
	}{
		{
			name: "photo mode",
			data: `{"payload_index": "39-0-7", "camera_mode": 0}`,
		},
		{
			name: "video mode",
			data: `{"payload_index": "39-0-7", "camera_mode": 1}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &ServiceRequest{
				Method: MethodCameraModeSwitch,
				Data:   json.RawMessage(tt.data),
			}
			resp, err := r.RouteService(context.Background(), req)
			require.NoError(t, err)
			assert.Equal(t, 0, resp.Result)
		})
	}
}

func TestCameraCommands_CameraPhotoTake(t *testing.T) {
	r := NewServiceRouter()
	err := RegisterCameraCommands(r)
	require.NoError(t, err)

	req := &ServiceRequest{
		Method: MethodCameraPhotoTake,
		Data:   json.RawMessage(`{"payload_index": "39-0-7"}`),
	}

	resp, err := r.RouteService(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Result)
}

func TestCameraCommands_CameraRecordingStart(t *testing.T) {
	r := NewServiceRouter()
	err := RegisterCameraCommands(r)
	require.NoError(t, err)

	req := &ServiceRequest{
		Method: MethodCameraRecordingStart,
		Data:   json.RawMessage(`{"payload_index": "39-0-7"}`),
	}

	resp, err := r.RouteService(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Result)
}

func TestCameraCommands_CameraRecordingStop(t *testing.T) {
	r := NewServiceRouter()
	err := RegisterCameraCommands(r)
	require.NoError(t, err)

	req := &ServiceRequest{
		Method: MethodCameraRecordingStop,
		Data:   json.RawMessage(`{"payload_index": "39-0-7"}`),
	}

	resp, err := r.RouteService(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Result)
}

func TestCameraCommands_CameraAim(t *testing.T) {
	r := NewServiceRouter()
	err := RegisterCameraCommands(r)
	require.NoError(t, err)

	req := &ServiceRequest{
		Method: MethodCameraAim,
		Data: json.RawMessage(`{
			"payload_index": "39-0-7",
			"camera_type": "wide",
			"locked": true,
			"x": 0.5,
			"y": 0.5
		}`),
	}

	resp, err := r.RouteService(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Result)
}

func TestCameraCommands_CameraFocalLengthSet(t *testing.T) {
	r := NewServiceRouter()
	err := RegisterCameraCommands(r)
	require.NoError(t, err)

	req := &ServiceRequest{
		Method: MethodCameraFocalLengthSet,
		Data: json.RawMessage(`{
			"payload_index": "39-0-7",
			"camera_type": "zoom",
			"zoom_factor": 5.0
		}`),
	}

	resp, err := r.RouteService(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Result)
}

func TestCameraCommands_GimbalReset(t *testing.T) {
	r := NewServiceRouter()
	err := RegisterCameraCommands(r)
	require.NoError(t, err)

	req := &ServiceRequest{
		Method: MethodGimbalReset,
		Data: json.RawMessage(`{
			"payload_index": "39-0-7",
			"reset_mode": 0
		}`),
	}

	resp, err := r.RouteService(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Result)
}

func TestCameraCommands_IRMeteringPoint(t *testing.T) {
	r := NewServiceRouter()
	err := RegisterCameraCommands(r)
	require.NoError(t, err)

	req := &ServiceRequest{
		Method: MethodIRMeteringPoint,
		Data: json.RawMessage(`{
			"payload_index": "39-0-7",
			"x": 0.5,
			"y": 0.5
		}`),
	}

	resp, err := r.RouteService(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Result)
}

func TestCameraCommands_IRMeteringArea(t *testing.T) {
	r := NewServiceRouter()
	err := RegisterCameraCommands(r)
	require.NoError(t, err)

	req := &ServiceRequest{
		Method: MethodIRMeteringArea,
		Data: json.RawMessage(`{
			"payload_index": "39-0-7",
			"x": 0.3,
			"y": 0.3,
			"width": 0.4,
			"height": 0.4
		}`),
	}

	resp, err := r.RouteService(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Result)
}

func TestCameraCommands_InvalidData(t *testing.T) {
	r := NewServiceRouter()
	err := RegisterCameraCommands(r)
	require.NoError(t, err)

	req := &ServiceRequest{
		Method: MethodCameraModeSwitch,
		Data:   json.RawMessage(`{invalid json}`),
	}

	resp, err := r.RouteService(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, 314000, resp.Result) // Parameter error
}
