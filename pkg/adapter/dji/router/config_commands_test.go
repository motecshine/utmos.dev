package router

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterConfigAndLiveCommands_ConfigMethods(t *testing.T) {
	r := NewServiceRouter()
	err := RegisterConfigAndLiveCommands(r)
	require.NoError(t, err)

	// Verify all config commands are registered
	expectedMethods := []string{
		MethodConfig,
		MethodStorageConfigGet,
		MethodPhotoStorageSet,
		MethodVideoStorageSet,
	}

	for _, method := range expectedMethods {
		assert.True(t, r.Has(method), "method %s should be registered", method)
	}
}

func TestConfigCommands_Config(t *testing.T) {
	r := NewServiceRouter()
	err := RegisterConfigAndLiveCommands(r)
	require.NoError(t, err)

	tests := []struct {
		name string
		data string
	}{
		{
			name: "config request with json type",
			data: `{"config_type": "json", "config_scope": "product"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &ServiceRequest{
				Method: MethodConfig,
				Data:   json.RawMessage(tt.data),
			}
			resp, err := r.RouteService(context.Background(), req)
			require.NoError(t, err)
			assert.Equal(t, 0, resp.Result)
		})
	}
}

func TestConfigCommands_StorageConfigGet(t *testing.T) {
	r := NewServiceRouter()
	err := RegisterConfigAndLiveCommands(r)
	require.NoError(t, err)

	tests := []struct {
		name string
		data string
	}{
		{
			name: "get media storage config",
			data: `{"module": 0}`,
		},
		{
			name: "get psdk ui resource config",
			data: `{"module": 1}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &ServiceRequest{
				Method: MethodStorageConfigGet,
				Data:   json.RawMessage(tt.data),
			}
			resp, err := r.RouteService(context.Background(), req)
			require.NoError(t, err)
			assert.Equal(t, 0, resp.Result)
		})
	}
}

func TestConfigCommands_PhotoStorageSet(t *testing.T) {
	r := NewServiceRouter()
	err := RegisterConfigAndLiveCommands(r)
	require.NoError(t, err)

	tests := []struct {
		name string
		data string
	}{
		{
			name: "set photo storage settings",
			data: `{"payload_index": "39-0-7", "photo_storage_settings": ["current", "wide"]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &ServiceRequest{
				Method: MethodPhotoStorageSet,
				Data:   json.RawMessage(tt.data),
			}
			resp, err := r.RouteService(context.Background(), req)
			require.NoError(t, err)
			assert.Equal(t, 0, resp.Result)
		})
	}
}

func TestConfigCommands_VideoStorageSet(t *testing.T) {
	r := NewServiceRouter()
	err := RegisterConfigAndLiveCommands(r)
	require.NoError(t, err)

	tests := []struct {
		name string
		data string
	}{
		{
			name: "set video storage settings",
			data: `{"payload_index": "39-0-7", "video_storage_settings": ["current", "zoom"]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &ServiceRequest{
				Method: MethodVideoStorageSet,
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
	err := RegisterConfigAndLiveCommands(r)
	require.NoError(t, err)

	req := &ServiceRequest{
		Method: MethodConfig,
		Data:   json.RawMessage(`{invalid json}`),
	}

	resp, err := r.RouteService(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, 314000, resp.Result) // Parameter error
}
