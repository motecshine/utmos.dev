package router

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterFileEvents(t *testing.T) {
	r := NewEventRouter()
	err := RegisterFileEvents(r)
	require.NoError(t, err)

	// Verify all file events are registered
	expectedMethods := []string{
		MethodFileUploadCallback,
		MethodFileUploadProgress,
		MethodHighestPriorityUpload,
	}

	for _, method := range expectedMethods {
		assert.True(t, r.Has(method), "method %s should be registered", method)
	}
}

func TestFileEvents_FileUploadCallback(t *testing.T) {
	r := NewEventRouter()
	err := RegisterFileEvents(r)
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

func TestFileEvents_FileUploadProgress(t *testing.T) {
	r := NewEventRouter()
	err := RegisterFileEvents(r)
	require.NoError(t, err)

	req := &EventRequest{
		Method: MethodFileUploadProgress,
		Data: json.RawMessage(`{
			"file_path": "/media/DJI_001.jpg",
			"progress": 75,
			"upload_rate": 1024
		}`),
	}

	resp, err := r.RouteEvent(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Result)
}

func TestFileEvents_HighestPriorityUpload(t *testing.T) {
	r := NewEventRouter()
	err := RegisterFileEvents(r)
	require.NoError(t, err)

	needReply := 1
	req := &EventRequest{
		Method:    MethodHighestPriorityUpload,
		NeedReply: &needReply,
		Data: json.RawMessage(`{
			"flight_id": "flight-001",
			"file_list": [
				{"path": "/media/DJI_001.jpg", "name": "DJI_001.jpg"}
			]
		}`),
	}

	resp, err := r.RouteEvent(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Result)
}
