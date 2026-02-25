package integration

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/utmos/utmos/internal/downlink"
	"github.com/utmos/utmos/internal/downlink/dispatcher"
	"github.com/utmos/utmos/internal/downlink/retry"
	"github.com/utmos/utmos/internal/downlink/router"
	djidownlink "github.com/utmos/utmos/pkg/adapter/dji/downlink"
)

// TestDownlinkServiceIntegration tests the downlink service integration
func TestDownlinkServiceIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("service lifecycle", func(t *testing.T) {
		config := &downlink.Config{
			RetryConfig: &retry.Config{
				MaxRetries:       3,
				InitialDelay:     100 * time.Millisecond,
				MaxDelay:         time.Second,
				Multiplier:       2.0,
				EnableDeadLetter: true,
			},
			RouterConfig: &router.Config{
				DefaultRoutingKey: router.RoutingKeyGatewayDownlink,
				EnableMetrics:     true,
			},
			EnableRetry:         true,
			EnableRouting:       false, // Disable routing for test (no publisher)
			RetryWorkerInterval: 100 * time.Millisecond,
		}

		svc := downlink.NewService(config, nil, nil, nil)
		require.NotNil(t, svc)

		// Register DJI dispatcher
		djiDispatcher := djidownlink.NewDispatcherAdapter(nil, nil)
		svc.RegisterAdapterDispatcher(djiDispatcher)

		// Verify vendor registration
		vendors := svc.GetRegisteredVendors()
		assert.Contains(t, vendors, "dji")

		// Start service
		err := svc.Start(context.Background())
		require.NoError(t, err)
		assert.True(t, svc.IsRunning())

		// Stop service
		err = svc.Stop()
		require.NoError(t, err)
		assert.False(t, svc.IsRunning())
	})
}

// TestRegistryIntegration tests dispatcher registry integration
func TestRegistryIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	registry := dispatcher.NewRegistry(nil)

	// Register DJI dispatcher
	djiDispatcher := djidownlink.NewDispatcherAdapter(nil, nil)
	registry.Register(dispatcher.NewAdapterDispatcher(djiDispatcher))

	t.Run("find dispatcher by vendor", func(t *testing.T) {
		d, found := registry.Get("dji")
		assert.True(t, found)
		assert.Equal(t, "dji", d.GetVendor())
	})

	t.Run("find dispatcher for call", func(t *testing.T) {
		call := &dispatcher.ServiceCall{
			DeviceSN: "DEVICE001",
			Vendor:   "dji",
			Method:   "takeoff",
		}

		d, found := registry.GetForCall(call)
		assert.True(t, found)
		assert.Equal(t, "dji", d.GetVendor())
	})

	t.Run("no dispatcher for unknown vendor", func(t *testing.T) {
		call := &dispatcher.ServiceCall{
			DeviceSN: "DEVICE001",
			Vendor:   "unknown",
			Method:   "takeoff",
		}

		_, found := registry.GetForCall(call)
		assert.False(t, found)
	})
}

// TestRetryHandlerIntegration tests retry handler integration
func TestRetryHandlerIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	config := &retry.Config{
		MaxRetries:       2,
		InitialDelay:     10 * time.Millisecond,
		MaxDelay:         100 * time.Millisecond,
		Multiplier:       2.0,
		EnableDeadLetter: true,
	}

	handler := retry.NewHandler(config, nil)

	t.Run("retry workflow", func(t *testing.T) {
		retryCount := 0
		handler.SetOnRetry(func(ctx context.Context, call *dispatcher.ServiceCall) error {
			retryCount++
			return assert.AnError // Keep failing
		})

		var deadLetterEntry *retry.DeadLetterEntry
		handler.SetOnDeadLetter(func(entry *retry.DeadLetterEntry) {
			deadLetterEntry = entry
		})

		call := &dispatcher.ServiceCall{
			ID:         "call-001",
			DeviceSN:   "DEVICE001",
			Vendor:     "dji",
			Method:     "takeoff",
			RetryCount: 0,
		}

		// Schedule first retry
		scheduled := handler.ScheduleRetry(call, "initial error")
		assert.True(t, scheduled)
		assert.Equal(t, 1, handler.GetPendingRetries())

		// Process retries until max retries exceeded
		ctx := context.Background()
		for i := 0; i < 5; i++ {
			time.Sleep(20 * time.Millisecond)
			handler.ProcessRetries(ctx)
		}

		// Should have moved to dead letter
		assert.Equal(t, 0, handler.GetPendingRetries())
		assert.Equal(t, 1, handler.GetDeadLetterCount())
		assert.NotNil(t, deadLetterEntry)
		assert.Equal(t, "call-001", deadLetterEntry.Call.ID)
	})

	t.Run("requeue from dead letter", func(t *testing.T) {
		// Requeue the entry
		requeued := handler.RequeueFromDeadLetter("call-001")
		assert.True(t, requeued)
		assert.Equal(t, 0, handler.GetDeadLetterCount())
		assert.Equal(t, 1, handler.GetPendingRetries())
	})
}

// TestDJIDispatcherIntegration tests DJI dispatcher integration
func TestDJIDispatcherIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	djiDispatcher := djidownlink.NewDispatcherAdapter(nil, nil)

	t.Run("can dispatch DJI calls", func(t *testing.T) {
		call := &dispatcher.ServiceCall{
			DeviceSN: "DEVICE001",
			Vendor:   "dji",
			Method:   "takeoff",
		}

		assert.True(t, djiDispatcher.CanDispatch(call.Vendor))
	})

	t.Run("cannot dispatch non-DJI calls", func(t *testing.T) {
		call := &dispatcher.ServiceCall{
			DeviceSN: "DEVICE001",
			Vendor:   "other",
			Method:   "takeoff",
		}

		assert.False(t, djiDispatcher.CanDispatch(call.Vendor))
	})

	t.Run("DJI service call helpers", func(t *testing.T) {
		// Test takeoff call
		takeoffCall := djidownlink.NewTakeoffCall("DEVICE001", 50.0)
		assert.Equal(t, "DEVICE001", takeoffCall.DeviceSN)
		assert.Equal(t, "dji", takeoffCall.Vendor)
		assert.Equal(t, "takeoff", takeoffCall.Method)
		assert.Equal(t, 50.0, takeoffCall.Params["height"])

		// Test land call
		landCall := djidownlink.NewLandCall("DEVICE001")
		assert.Equal(t, "land", landCall.Method)

		// Test return home call
		returnHomeCall := djidownlink.NewReturnHomeCall("DEVICE001")
		assert.Equal(t, "return_home", returnHomeCall.Method)

		// Test fly to point call
		flyToPointCall := djidownlink.NewFlyToPointCall("DEVICE001", 22.5431, 113.9234, 100.0, 10.0)
		assert.Equal(t, "fly_to_point", flyToPointCall.Method)
		assert.Equal(t, 22.5431, flyToPointCall.Params["latitude"])
		assert.Equal(t, 113.9234, flyToPointCall.Params["longitude"])
		assert.Equal(t, 100.0, flyToPointCall.Params["altitude"])
		assert.Equal(t, 10.0, flyToPointCall.Params["speed"])

		// Test gimbal rotate call
		gimbalCall := djidownlink.NewGimbalRotateCall("DEVICE001", -30.0, 45.0)
		assert.Equal(t, "gimbal_rotate", gimbalCall.Method)
		assert.Equal(t, -30.0, gimbalCall.Params["pitch"])
		assert.Equal(t, 45.0, gimbalCall.Params["yaw"])

		// Test camera photo call
		photoCall := djidownlink.NewCameraPhotoCall("DEVICE001")
		assert.Equal(t, "camera_photo_take", photoCall.Method)
	})
}

// TestDispatchHandlerIntegration tests dispatch handler integration
func TestDispatchHandlerIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	registry := dispatcher.NewRegistry(nil)
	handler := dispatcher.NewDispatchHandler(registry, nil)

	// Register DJI dispatcher
	djiDispatcher := djidownlink.NewDispatcherAdapter(nil, nil)
	handler.RegisterAdapterDispatcher(djiDispatcher)

	t.Run("dispatch to registered vendor", func(t *testing.T) {
		paramsJSON, _ := json.Marshal(map[string]interface{}{"height": 50})
		call := &dispatcher.ServiceCall{
			ID:       "call-001",
			DeviceSN: "DEVICE001",
			Vendor:   "dji",
			Method:   "takeoff",
			Params:   paramsJSON,
		}

		// Will fail because no publisher, but should find dispatcher
		_, err := handler.Handle(context.Background(), call)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "publisher not initialized")
	})

	t.Run("dispatch to unregistered vendor", func(t *testing.T) {
		call := &dispatcher.ServiceCall{
			ID:       "call-002",
			DeviceSN: "DEVICE001",
			Vendor:   "unknown",
			Method:   "takeoff",
		}

		_, err := handler.Handle(context.Background(), call)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no dispatcher found")
	})
}

// TestRouterIntegration tests router integration
func TestRouterIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	config := &router.Config{
		DefaultRoutingKey: router.RoutingKeyGatewayDownlink,
		EnableMetrics:     true,
	}

	r := router.NewRouter(nil, config, nil)

	t.Run("metrics tracking", func(t *testing.T) {
		// Initial metrics
		routed, failed := r.GetMetrics()
		assert.Equal(t, int64(0), routed)
		assert.Equal(t, int64(0), failed)

		// Route will fail (no publisher) - returns error before incrementing metrics
		call := &dispatcher.ServiceCall{
			ID:       "call-001",
			DeviceSN: "DEVICE001",
			Vendor:   "dji",
			Method:   "takeoff",
		}

		_, err := r.Route(context.Background(), call, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "publisher not initialized")

		// Metrics should still be zero since error occurred before routing attempt
		routed, failed = r.GetMetrics()
		assert.Equal(t, int64(0), routed)
		assert.Equal(t, int64(0), failed)
	})
}

// TestServiceCallCreation tests service call creation helpers
func TestServiceCallCreation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("new service call", func(t *testing.T) {
		params := map[string]interface{}{
			"height": 50,
			"speed":  10,
		}
		paramsJSON, _ := json.Marshal(params)

		call := dispatcher.NewServiceCall("DEVICE001", "dji", "takeoff", paramsJSON)

		assert.Equal(t, "DEVICE001", call.DeviceSN)
		assert.Equal(t, "dji", call.Vendor)
		assert.Equal(t, "takeoff", call.Method)
		assert.Equal(t, paramsJSON, []byte(call.Params))
		assert.Equal(t, dispatcher.ServiceCallTypeCommand, call.CallType)
		assert.Equal(t, dispatcher.ServiceCallStatusPending, call.Status)
		assert.Equal(t, 3, call.MaxRetries)
		assert.False(t, call.CreatedAt.IsZero())
	})
}

// TestEndToEndDispatchFlow tests the end-to-end dispatch flow
func TestEndToEndDispatchFlow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Create service with retry enabled but routing disabled (no publisher)
	config := &downlink.Config{
		RetryConfig: &retry.Config{
			MaxRetries:       2,
			InitialDelay:     10 * time.Millisecond,
			MaxDelay:         100 * time.Millisecond,
			Multiplier:       2.0,
			EnableDeadLetter: true,
		},
		EnableRetry:         true,
		EnableRouting:       false,
		RetryWorkerInterval: 50 * time.Millisecond,
	}

	svc := downlink.NewService(config, nil, nil, nil)
	djiDispatcher := djidownlink.NewDispatcherAdapter(nil, nil)
	svc.RegisterAdapterDispatcher(djiDispatcher)

	// Start service
	ctx := context.Background()
	err := svc.Start(ctx)
	require.NoError(t, err)
	defer func() { _ = svc.Stop() }()

	t.Run("dispatch fails and schedules retry", func(t *testing.T) {
		call := djidownlink.NewTakeoffCall("DEVICE001", 50.0)
		// Convert to dispatcher.ServiceCall
		paramsJSON, _ := json.Marshal(call.Params)
		dispatcherCall := &dispatcher.ServiceCall{
			DeviceSN: call.DeviceSN,
			Vendor:   call.Vendor,
			Method:   call.Method,
			Params:   paramsJSON,
		}

		// Dispatch will fail (no publisher)
		_, err := svc.Dispatch(ctx, dispatcherCall)
		assert.Error(t, err)

		// Should have scheduled retry
		pending, _ := svc.GetRetryMetrics()
		assert.Equal(t, 1, pending)

		// Check metrics
		processed, failed := svc.GetMetrics()
		assert.Equal(t, int64(0), processed)
		assert.Equal(t, int64(1), failed)
	})
}
