package downlink

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/utmos/utmos/internal/downlink/dispatcher"
	"github.com/utmos/utmos/internal/downlink/retry"
	"github.com/utmos/utmos/internal/downlink/router"
	djidownlink "github.com/utmos/utmos/pkg/adapter/dji/downlink"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	assert.NotNil(t, config.RetryConfig)
	assert.NotNil(t, config.RouterConfig)
	assert.True(t, config.EnableRetry)
	assert.True(t, config.EnableRouting)
	assert.Equal(t, 5*time.Second, config.RetryWorkerInterval)
}

func TestNewService(t *testing.T) {
	t.Run("with config", func(t *testing.T) {
		config := &Config{
			EnableRetry:   true,
			EnableRouting: false,
			RetryConfig:   retry.DefaultConfig(),
		}
		svc := NewService(config, nil, nil, nil)

		require.NotNil(t, svc)
		assert.NotNil(t, svc.registry)
		assert.NotNil(t, svc.handler)
		assert.NotNil(t, svc.retryHandler)
		assert.Nil(t, svc.router) // Routing disabled
	})

	t.Run("without config", func(t *testing.T) {
		svc := NewService(nil, nil, nil, nil)

		require.NotNil(t, svc)
		assert.NotNil(t, svc.retryHandler)
		assert.NotNil(t, svc.router)
	})

	t.Run("retry disabled", func(t *testing.T) {
		config := &Config{
			EnableRetry:   false,
			EnableRouting: false,
		}
		svc := NewService(config, nil, nil, nil)

		require.NotNil(t, svc)
		assert.Nil(t, svc.retryHandler)
	})
}

func TestService_RegisterDispatcher(t *testing.T) {
	svc := NewService(nil, nil, nil, nil)

	// Create a mock dispatcher
	mockDispatcher := dispatcher.NewBaseDispatcher("test-vendor", nil, nil)

	// Register DJI dispatcher using the adapter
	djiDispatcher := djidownlink.NewDispatcherAdapter(nil, nil)
	svc.RegisterAdapterDispatcher(djiDispatcher)

	vendors := svc.GetRegisteredVendors()
	assert.Contains(t, vendors, "dji")

	// Verify mock dispatcher base
	assert.Equal(t, "test-vendor", mockDispatcher.GetVendor())
}

func TestService_StartStop(t *testing.T) {
	config := &Config{
		EnableRetry:         true,
		EnableRouting:       false,
		RetryConfig:         retry.DefaultConfig(),
		RetryWorkerInterval: 100 * time.Millisecond,
	}
	svc := NewService(config, nil, nil, nil)

	t.Run("start service", func(t *testing.T) {
		err := svc.Start(context.Background())
		assert.NoError(t, err)
		assert.True(t, svc.IsRunning())
	})

	t.Run("start already running", func(t *testing.T) {
		err := svc.Start(context.Background())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already running")
	})

	t.Run("stop service", func(t *testing.T) {
		err := svc.Stop()
		assert.NoError(t, err)
		assert.False(t, svc.IsRunning())
	})

	t.Run("stop already stopped", func(t *testing.T) {
		err := svc.Stop()
		assert.NoError(t, err)
	})
}

func TestService_Dispatch_NilCall(t *testing.T) {
	svc := NewService(nil, nil, nil, nil)

	_, err := svc.Dispatch(context.Background(), nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "service call is nil")
}

func TestService_Dispatch_NoDispatcher(t *testing.T) {
	config := &Config{
		EnableRetry:   true,
		EnableRouting: false,
		RetryConfig:   retry.DefaultConfig(),
	}
	svc := NewService(config, nil, nil, nil)

	call := &dispatcher.ServiceCall{
		ID:       "call-001",
		DeviceSN: "DEVICE001",
		Vendor:   "unknown",
		Method:   "takeoff",
	}

	_, err := svc.Dispatch(context.Background(), call)
	assert.Error(t, err)

	// Should have scheduled retry
	pending, _ := svc.GetRetryMetrics()
	assert.Equal(t, 1, pending)
}

func TestService_Dispatch_WithDJIDispatcher(t *testing.T) {
	config := &Config{
		EnableRetry:   true,
		EnableRouting: false,
		RetryConfig:   retry.DefaultConfig(),
	}
	svc := NewService(config, nil, nil, nil)
	djiDispatcher := djidownlink.NewDispatcherAdapter(nil, nil)
	svc.RegisterAdapterDispatcher(djiDispatcher)

	paramsJSON, _ := json.Marshal(map[string]any{"height": 50})
	call := &dispatcher.ServiceCall{
		ID:       "call-001",
		DeviceSN: "DEVICE001",
		Vendor:   "dji",
		Method:   "takeoff",
		Params:   paramsJSON,
	}

	// Will fail because no publisher, but dispatcher should be found
	_, err := svc.Dispatch(context.Background(), call)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "publisher not initialized")
}

func TestService_GetMetrics(t *testing.T) {
	svc := NewService(nil, nil, nil, nil)

	// Initial metrics
	processed, failed := svc.GetMetrics()
	assert.Equal(t, int64(0), processed)
	assert.Equal(t, int64(0), failed)

	// Increment counters
	svc.incrementProcessed()
	svc.incrementProcessed()
	svc.incrementFailed()

	processed, failed = svc.GetMetrics()
	assert.Equal(t, int64(2), processed)
	assert.Equal(t, int64(1), failed)
}

func TestService_GetRetryMetrics(t *testing.T) {
	t.Run("with retry enabled", func(t *testing.T) {
		config := &Config{
			EnableRetry:   true,
			EnableRouting: false,
			RetryConfig:   retry.DefaultConfig(),
		}
		svc := NewService(config, nil, nil, nil)

		pending, deadLetter := svc.GetRetryMetrics()
		assert.Equal(t, 0, pending)
		assert.Equal(t, 0, deadLetter)
	})

	t.Run("with retry disabled", func(t *testing.T) {
		config := &Config{
			EnableRetry:   false,
			EnableRouting: false,
		}
		svc := NewService(config, nil, nil, nil)

		pending, deadLetter := svc.GetRetryMetrics()
		assert.Equal(t, 0, pending)
		assert.Equal(t, 0, deadLetter)
	})
}

func TestService_GetRouterMetrics(t *testing.T) {
	t.Run("with routing enabled", func(t *testing.T) {
		config := &Config{
			EnableRetry:   false,
			EnableRouting: true,
			RouterConfig:  router.DefaultConfig(),
		}
		svc := NewService(config, nil, nil, nil)

		routed, failed := svc.GetRouterMetrics()
		assert.Equal(t, int64(0), routed)
		assert.Equal(t, int64(0), failed)
	})

	t.Run("with routing disabled", func(t *testing.T) {
		config := &Config{
			EnableRetry:   false,
			EnableRouting: false,
		}
		svc := NewService(config, nil, nil, nil)

		routed, failed := svc.GetRouterMetrics()
		assert.Equal(t, int64(0), routed)
		assert.Equal(t, int64(0), failed)
	})
}

func TestService_GetRegisteredVendors(t *testing.T) {
	svc := NewService(nil, nil, nil, nil)

	// Initially empty
	vendors := svc.GetRegisteredVendors()
	assert.Empty(t, vendors)

	// Register DJI
	djiDispatcher := djidownlink.NewDispatcherAdapter(nil, nil)
	svc.RegisterAdapterDispatcher(djiDispatcher)

	vendors = svc.GetRegisteredVendors()
	assert.Len(t, vendors, 1)
	assert.Contains(t, vendors, "dji")
}

func TestService_SetSubscriber(t *testing.T) {
	svc := NewService(nil, nil, nil, nil)

	assert.Nil(t, svc.subscriber)

	// SetSubscriber would be called with actual subscriber
	// For now just verify the method exists
	svc.SetSubscriber(nil)
	assert.Nil(t, svc.subscriber)
}

func TestService_OnDispatched(t *testing.T) {
	config := &Config{
		EnableRetry:   false,
		EnableRouting: false, // Disable routing to avoid nil publisher error
	}
	svc := NewService(config, nil, nil, nil)

	call := &dispatcher.ServiceCall{
		ID:       "call-001",
		DeviceSN: "DEVICE001",
		Vendor:   "dji",
		Method:   "takeoff",
	}

	result := &dispatcher.DispatchResult{
		Success:   true,
		MessageID: "msg-001",
		SentAt:    time.Now(),
	}

	// Should not error when routing is disabled
	err := svc.onDispatched(context.Background(), call, result)
	assert.NoError(t, err)
}

func TestService_OnRetry(t *testing.T) {
	config := &Config{
		EnableRetry:   false,
		EnableRouting: false,
	}
	svc := NewService(config, nil, nil, nil)
	djiDispatcher := djidownlink.NewDispatcherAdapter(nil, nil)
	svc.RegisterAdapterDispatcher(djiDispatcher)

	call := &dispatcher.ServiceCall{
		ID:         "call-001",
		DeviceSN:   "DEVICE001",
		Vendor:     "dji",
		Method:     "takeoff",
		RetryCount: 1,
	}

	// Will fail because no publisher
	err := svc.onRetry(context.Background(), call)
	assert.Error(t, err)
}

func TestService_OnDeadLetter(t *testing.T) {
	svc := NewService(nil, nil, nil, nil)

	entry := &retry.DeadLetterEntry{
		Call: &dispatcher.ServiceCall{
			ID:       "call-001",
			DeviceSN: "DEVICE001",
		},
		Error:    "max retries exceeded",
		FailedAt: time.Now(),
		Retries:  3,
	}

	// Should not panic
	svc.onDeadLetter(entry)
}
