package router

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterFirmwareCommands(t *testing.T) {
	r := NewServiceRouter()
	err := RegisterFirmwareCommands(r)
	require.NoError(t, err)

	// Verify all firmware commands are registered
	expectedMethods := []string{
		MethodOTACreate,
	}

	for _, method := range expectedMethods {
		assert.True(t, r.Has(method), "method %s should be registered", method)
	}
}

func TestFirmwareCommands_OTACreate(t *testing.T) {
	r := NewServiceRouter()
	err := RegisterFirmwareCommands(r)
	require.NoError(t, err)

	req := &ServiceRequest{
		Method: MethodOTACreate,
		Data: json.RawMessage(`{
			"devices": [
				{
					"sn": "DOCK-SN-001",
					"product_version": "01.02.0400",
					"file_url": "https://firmware.example.com/update.bin",
					"md5": "abc123",
					"file_size": 1024000,
					"file_name": "firmware.bin",
					"firmware_upgrade_type": 3
				}
			]
		}`),
	}

	resp, err := r.RouteService(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Result)
}

func TestFirmwareCommands_InvalidData(t *testing.T) {
	r := NewServiceRouter()
	err := RegisterFirmwareCommands(r)
	require.NoError(t, err)

	req := &ServiceRequest{
		Method: MethodOTACreate,
		Data:   json.RawMessage(`{invalid json}`),
	}

	resp, err := r.RouteService(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, 314000, resp.Result) // Parameter error
}
