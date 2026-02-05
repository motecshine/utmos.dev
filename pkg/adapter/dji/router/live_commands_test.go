package router

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterLiveCommands(t *testing.T) {
	r := NewServiceRouter()
	err := RegisterLiveCommands(r)
	require.NoError(t, err)

	// Verify all live commands are registered
	expectedMethods := []string{
		MethodLiveStartPush,
		MethodLiveStopPush,
		MethodLiveSetQuality,
		MethodLiveLensChange,
	}

	for _, method := range expectedMethods {
		assert.True(t, r.Has(method), "method %s should be registered", method)
	}
}

func TestLiveCommands_LiveStartPush(t *testing.T) {
	r := NewServiceRouter()
	err := RegisterLiveCommands(r)
	require.NoError(t, err)

	req := &ServiceRequest{
		Method: MethodLiveStartPush,
		Data: json.RawMessage(`{
			"url_type": 1,
			"url": "rtmp://live.example.com/stream",
			"video_id": "SN001/39-0-7/normal-0",
			"video_quality": 0
		}`),
	}

	resp, err := r.RouteService(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Result)
}

func TestLiveCommands_LiveStopPush(t *testing.T) {
	r := NewServiceRouter()
	err := RegisterLiveCommands(r)
	require.NoError(t, err)

	req := &ServiceRequest{
		Method: MethodLiveStopPush,
		Data: json.RawMessage(`{
			"video_id": "SN001/39-0-7/normal-0"
		}`),
	}

	resp, err := r.RouteService(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Result)
}

func TestLiveCommands_LiveSetQuality(t *testing.T) {
	r := NewServiceRouter()
	err := RegisterLiveCommands(r)
	require.NoError(t, err)

	req := &ServiceRequest{
		Method: MethodLiveSetQuality,
		Data: json.RawMessage(`{
			"video_id": "SN001/39-0-7/normal-0",
			"video_quality": 3
		}`),
	}

	resp, err := r.RouteService(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Result)
}

func TestLiveCommands_LiveLensChange(t *testing.T) {
	r := NewServiceRouter()
	err := RegisterLiveCommands(r)
	require.NoError(t, err)

	req := &ServiceRequest{
		Method: MethodLiveLensChange,
		Data: json.RawMessage(`{
			"video_type": "wide"
		}`),
	}

	resp, err := r.RouteService(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Result)
}

func TestLiveCommands_InvalidData(t *testing.T) {
	r := NewServiceRouter()
	err := RegisterLiveCommands(r)
	require.NoError(t, err)

	req := &ServiceRequest{
		Method: MethodLiveStartPush,
		Data:   json.RawMessage(`{invalid json}`),
	}

	resp, err := r.RouteService(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, 314000, resp.Result) // Parameter error
}
