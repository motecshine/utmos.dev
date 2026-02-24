package router

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewServiceRouter(t *testing.T) {
	r := NewServiceRouter()
	assert.NotNil(t, r)
	assert.NotNil(t, r.registry.handlers)
}

func TestServiceRouter_RegisterServiceHandler(t *testing.T) {
	r := NewServiceRouter()

	handler := func(_ context.Context, _ json.RawMessage) (*ServiceResponse, error) {
		return &ServiceResponse{Result: 0}, nil
	}

	// First registration should succeed
	err := r.RegisterServiceHandler("test_method", handler)
	require.NoError(t, err)

	// Duplicate registration should fail
	err = r.RegisterServiceHandler("test_method", handler)
	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrMethodAlreadyRegistered))
}

func TestServiceRouter_RouteService(t *testing.T) {
	r := NewServiceRouter()

	handler := func(_ context.Context, _ json.RawMessage) (*ServiceResponse, error) {
		return &ServiceResponse{
			Result: 0,
			Output: json.RawMessage(`{"status": "success"}`),
		}, nil
	}

	err := r.RegisterServiceHandler("test_method", handler)
	require.NoError(t, err)

	req := &ServiceRequest{
		Method: "test_method",
		Data:   json.RawMessage(`{"key": "value"}`),
	}

	resp, err := r.RouteService(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Result)
	assert.NotNil(t, resp.Output)
}

func TestServiceRouter_RouteService_NilRequest(t *testing.T) {
	r := NewServiceRouter()

	_, err := r.RouteService(context.Background(), nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "nil service request")
}

func TestServiceRouter_RouteService_UnknownMethod(t *testing.T) {
	r := NewServiceRouter()

	req := &ServiceRequest{
		Method: "unknown_method",
	}

	_, err := r.RouteService(context.Background(), req)
	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrMethodNotFound))
}

func TestServiceRouter_RouteService_HandlerError(t *testing.T) {
	r := NewServiceRouter()

	expectedErr := errors.New("handler error")
	handler := func(_ context.Context, _ json.RawMessage) (*ServiceResponse, error) {
		return nil, expectedErr
	}

	err := r.RegisterServiceHandler("error_method", handler)
	require.NoError(t, err)

	req := &ServiceRequest{
		Method: "error_method",
	}

	_, err = r.RouteService(context.Background(), req)
	require.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestServiceRouter_RouteService_WithData(t *testing.T) {
	r := NewServiceRouter()

	handler := func(_ context.Context, data json.RawMessage) (*ServiceResponse, error) {
		// Parse and echo back the data
		var input map[string]any
		if err := json.Unmarshal(data, &input); err != nil {
			return nil, err
		}

		output, _ := json.Marshal(map[string]any{
			"received": input,
		})

		return &ServiceResponse{
			Result: 0,
			Output: output,
		}, nil
	}

	err := r.RegisterServiceHandler("echo", handler)
	require.NoError(t, err)

	req := &ServiceRequest{
		Method: "echo",
		Data:   json.RawMessage(`{"message": "hello"}`),
	}

	resp, err := r.RouteService(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Result)

	var output map[string]any
	err = json.Unmarshal(resp.Output, &output)
	require.NoError(t, err)

	received := output["received"].(map[string]any)
	assert.Equal(t, "hello", received["message"])
}

func TestServiceRouter_MultipleHandlers(t *testing.T) {
	r := NewServiceRouter()

	methods := []string{
		"cover_open",
		"cover_close",
		"drone_open",
		"drone_close",
		"device_reboot",
	}

	for _, method := range methods {
		m := method // capture
		handler := func(_ context.Context, _ json.RawMessage) (*ServiceResponse, error) {
			return &ServiceResponse{
				Result: 0,
				Output: json.RawMessage(`{"method": "` + m + `"}`),
			}, nil
		}
		err := r.RegisterServiceHandler(method, handler)
		require.NoError(t, err)
	}

	// Verify all methods are registered
	list := r.List()
	assert.Len(t, list, len(methods))

	// Route to each method
	for _, method := range methods {
		req := &ServiceRequest{Method: method}
		resp, err := r.RouteService(context.Background(), req)
		require.NoError(t, err)
		assert.Equal(t, 0, resp.Result)
	}
}
