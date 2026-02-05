package router

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSimpleCommandHandler(t *testing.T) {
	type TestData struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	handler := SimpleCommandHandler[TestData]("test_method")

	tests := []struct {
		name       string
		data       json.RawMessage
		wantResult int
	}{
		{
			name:       "valid data",
			data:       json.RawMessage(`{"name": "test", "value": 42}`),
			wantResult: ResultSuccess,
		},
		{
			name:       "empty data",
			data:       nil,
			wantResult: ResultSuccess,
		},
		{
			name:       "invalid json",
			data:       json.RawMessage(`{invalid}`),
			wantResult: ResultParamError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := handler(context.Background(), tt.data)
			require.NoError(t, err)
			assert.Equal(t, tt.wantResult, resp.Result)
		})
	}
}

func TestNoDataCommandHandler(t *testing.T) {
	handler := NoDataCommandHandler("test_method")

	resp, err := handler(context.Background(), nil)
	require.NoError(t, err)
	assert.Equal(t, ResultSuccess, resp.Result)
	assert.Contains(t, string(resp.Output), `"method": "test_method"`)
	assert.Contains(t, string(resp.Output), `"status": "accepted"`)
}

func TestRegisterHandlers(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		r := NewServiceRouter()
		handlers := map[string]ServiceHandlerFunc{
			"method1": NoDataCommandHandler("method1"),
			"method2": NoDataCommandHandler("method2"),
		}

		err := RegisterHandlers(r, handlers)
		require.NoError(t, err)
		assert.True(t, r.Has("method1"))
		assert.True(t, r.Has("method2"))
	})

	t.Run("duplicate method error", func(t *testing.T) {
		r := NewServiceRouter()
		// First registration
		err := r.RegisterServiceHandler("method1", NoDataCommandHandler("method1"))
		require.NoError(t, err)

		// Try to register again
		handlers := map[string]ServiceHandlerFunc{
			"method1": NoDataCommandHandler("method1"),
		}

		err = RegisterHandlers(r, handlers)
		require.Error(t, err)

		var regErr *RegistrationError
		assert.True(t, errors.As(err, &regErr))
		assert.Equal(t, "method1", regErr.Method)
	})
}

func TestRegistrationError(t *testing.T) {
	innerErr := errors.New("inner error")
	err := &RegistrationError{
		Method: "test_method",
		Err:    innerErr,
	}

	assert.Contains(t, err.Error(), "test_method")
	assert.Contains(t, err.Error(), "inner error")
	assert.Equal(t, innerErr, err.Unwrap())
}
