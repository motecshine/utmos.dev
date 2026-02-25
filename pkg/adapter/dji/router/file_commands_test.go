package router

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterFileAndFirmwareCommands_FileMethods(t *testing.T) {
	r := NewServiceRouter()
	err := RegisterFileAndFirmwareCommands(r)
	require.NoError(t, err)

	// Verify all file commands are registered
	expectedMethods := []string{
		MethodFileUploadStart,
		MethodFileUploadFinish,
		MethodFileUploadList,
	}

	for _, method := range expectedMethods {
		assert.True(t, r.Has(method), "method %s should be registered", method)
	}
}

func TestFileCommands_FileUploadStart(t *testing.T) {
	r := NewServiceRouter()
	err := RegisterFileAndFirmwareCommands(r)
	require.NoError(t, err)

	req := &ServiceRequest{
		Method: MethodFileUploadStart,
		Data: json.RawMessage(`{
			"bucket": "test-bucket",
			"region": "us-east-1",
			"endpoint": "https://s3.amazonaws.com",
			"provider": "aws",
			"credentials": {
				"access_key_id": "test-key",
				"access_key_secret": "test-secret",
				"expire": 3600,
				"security_token": "test-token"
			},
			"params": {
				"files": []
			}
		}`),
	}

	resp, err := r.RouteService(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Result)
}

func TestFileCommands_FileUploadFinish(t *testing.T) {
	r := NewServiceRouter()
	err := RegisterFileAndFirmwareCommands(r)
	require.NoError(t, err)

	req := &ServiceRequest{
		Method: MethodFileUploadFinish,
		Data: json.RawMessage(`{
			"status": "cancel",
			"module_list": ["0", "3"]
		}`),
	}

	resp, err := r.RouteService(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Result)
}

func TestFileCommands_FileUploadList(t *testing.T) {
	r := NewServiceRouter()
	err := RegisterFileAndFirmwareCommands(r)
	require.NoError(t, err)

	req := &ServiceRequest{
		Method: MethodFileUploadList,
		Data: json.RawMessage(`{
			"module_list": ["0", "3"]
		}`),
	}

	resp, err := r.RouteService(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Result)
}

func TestFileCommands_InvalidData(t *testing.T) {
	r := NewServiceRouter()
	err := RegisterFileAndFirmwareCommands(r)
	require.NoError(t, err)

	req := &ServiceRequest{
		Method: MethodFileUploadStart,
		Data:   json.RawMessage(`{invalid json}`),
	}

	resp, err := r.RouteService(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, 314000, resp.Result) // Parameter error
}
