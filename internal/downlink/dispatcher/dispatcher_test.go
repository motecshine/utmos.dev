package dispatcher

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockDispatcher is a test dispatcher implementation
type mockDispatcher struct {
	*BaseDispatcher
	canDispatchFunc func(call *ServiceCall) bool
	dispatchFunc    func(ctx context.Context, call *ServiceCall) (*DispatchResult, error)
}

func newMockDispatcher(vendor string) *mockDispatcher {
	return &mockDispatcher{
		BaseDispatcher: NewBaseDispatcher(vendor, nil, nil),
	}
}

func (d *mockDispatcher) CanDispatch(call *ServiceCall) bool {
	if d.canDispatchFunc != nil {
		return d.canDispatchFunc(call)
	}
	return call.Vendor == d.vendor
}

func (d *mockDispatcher) Dispatch(ctx context.Context, call *ServiceCall) (*DispatchResult, error) {
	if d.dispatchFunc != nil {
		return d.dispatchFunc(ctx, call)
	}
	return &DispatchResult{
		Success:   true,
		MessageID: "msg-" + call.ID,
		SentAt:    time.Now(),
	}, nil
}

func TestNewRegistry(t *testing.T) {
	registry := NewRegistry(nil)
	require.NotNil(t, registry)
	assert.Equal(t, 0, registry.Count())
}

func TestRegistry_Register(t *testing.T) {
	registry := NewRegistry(nil)
	dispatcher := newMockDispatcher("test-vendor")

	registry.Register(dispatcher)

	d, exists := registry.Get("test-vendor")
	assert.True(t, exists)
	assert.Equal(t, "test-vendor", d.GetVendor())
}

func TestRegistry_Unregister(t *testing.T) {
	registry := NewRegistry(nil)
	dispatcher := newMockDispatcher("test-vendor")

	registry.Register(dispatcher)
	registry.Unregister("test-vendor")

	_, exists := registry.Get("test-vendor")
	assert.False(t, exists)
}

func TestRegistry_GetForCall(t *testing.T) {
	registry := NewRegistry(nil)
	djiDispatcher := newMockDispatcher("dji")
	customDispatcher := newMockDispatcher("custom")

	registry.Register(djiDispatcher)
	registry.Register(customDispatcher)

	t.Run("match by vendor", func(t *testing.T) {
		call := &ServiceCall{
			DeviceSN: "DEVICE001",
			Vendor:   "dji",
			Method:   "takeoff",
		}

		dispatcher, found := registry.GetForCall(call)
		assert.True(t, found)
		assert.Equal(t, "dji", dispatcher.GetVendor())
	})

	t.Run("no match", func(t *testing.T) {
		call := &ServiceCall{
			DeviceSN: "DEVICE001",
			Vendor:   "unknown",
			Method:   "takeoff",
		}

		_, found := registry.GetForCall(call)
		assert.False(t, found)
	})

	t.Run("match by CanDispatch", func(t *testing.T) {
		customDispatcher.canDispatchFunc = func(call *ServiceCall) bool {
			return call.Method == "custom_method"
		}

		call := &ServiceCall{
			DeviceSN: "DEVICE001",
			Vendor:   "",
			Method:   "custom_method",
		}

		dispatcher, found := registry.GetForCall(call)
		assert.True(t, found)
		assert.Equal(t, "custom", dispatcher.GetVendor())
	})
}

func TestRegistry_ListVendors(t *testing.T) {
	registry := NewRegistry(nil)
	registry.Register(newMockDispatcher("vendor1"))
	registry.Register(newMockDispatcher("vendor2"))
	registry.Register(newMockDispatcher("vendor3"))

	vendors := registry.ListVendors()
	assert.Len(t, vendors, 3)
	assert.Contains(t, vendors, "vendor1")
	assert.Contains(t, vendors, "vendor2")
	assert.Contains(t, vendors, "vendor3")
}

func TestNewDispatchHandler(t *testing.T) {
	registry := NewRegistry(nil)
	handler := NewDispatchHandler(registry, nil)

	require.NotNil(t, handler)
	assert.NotNil(t, handler.registry)
}

func TestDispatchHandler_Handle(t *testing.T) {
	registry := NewRegistry(nil)
	dispatcher := newMockDispatcher("dji")
	registry.Register(dispatcher)

	handler := NewDispatchHandler(registry, nil)

	t.Run("nil call", func(t *testing.T) {
		_, err := handler.Handle(context.Background(), nil)
		assert.Error(t, err)
	})

	t.Run("no dispatcher found", func(t *testing.T) {
		call := &ServiceCall{
			ID:       "call-001",
			DeviceSN: "DEVICE001",
			Vendor:   "unknown",
			Method:   "takeoff",
		}

		_, err := handler.Handle(context.Background(), call)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no dispatcher found")
	})

	t.Run("successful dispatch", func(t *testing.T) {
		paramsJSON, _ := json.Marshal(map[string]any{"height": 50})
		call := &ServiceCall{
			ID:       "call-002",
			DeviceSN: "DEVICE001",
			Vendor:   "dji",
			Method:   "takeoff",
			Params:   paramsJSON,
		}

		var dispatchedCall *ServiceCall
		handler.SetOnDispatched(func(ctx context.Context, c *ServiceCall, result *DispatchResult) error {
			dispatchedCall = c
			return nil
		})

		result, err := handler.Handle(context.Background(), call)
		assert.NoError(t, err)
		require.NotNil(t, result)
		assert.True(t, result.Success)
		assert.NotEmpty(t, result.MessageID)
		assert.Equal(t, call, dispatchedCall)
	})
}

func TestDispatchHandler_RegisterDispatcher(t *testing.T) {
	registry := NewRegistry(nil)
	handler := NewDispatchHandler(registry, nil)

	dispatcher := newMockDispatcher("new-vendor")
	handler.RegisterDispatcher(dispatcher)

	d, exists := registry.Get("new-vendor")
	assert.True(t, exists)
	assert.Equal(t, "new-vendor", d.GetVendor())
}

func TestDispatchHandler_UnregisterDispatcher(t *testing.T) {
	registry := NewRegistry(nil)
	handler := NewDispatchHandler(registry, nil)

	dispatcher := newMockDispatcher("test-vendor")
	handler.RegisterDispatcher(dispatcher)
	handler.UnregisterDispatcher("test-vendor")

	_, exists := registry.Get("test-vendor")
	assert.False(t, exists)
}

func TestBaseDispatcher(t *testing.T) {
	dispatcher := NewBaseDispatcher("test-vendor", nil, nil)

	assert.Equal(t, "test-vendor", dispatcher.GetVendor())
	assert.Nil(t, dispatcher.Publisher())
	assert.NotNil(t, dispatcher.Logger())
}

func TestServiceCallType(t *testing.T) {
	assert.Equal(t, ServiceCallType("command"), ServiceCallTypeCommand)
	assert.Equal(t, ServiceCallType("property"), ServiceCallTypeProperty)
	assert.Equal(t, ServiceCallType("config"), ServiceCallTypeConfig)
}

func TestServiceCallStatus(t *testing.T) {
	assert.Equal(t, ServiceCallStatus("pending"), ServiceCallStatusPending)
	assert.Equal(t, ServiceCallStatus("sent"), ServiceCallStatusSent)
	assert.Equal(t, ServiceCallStatus("success"), ServiceCallStatusSuccess)
	assert.Equal(t, ServiceCallStatus("failed"), ServiceCallStatusFailed)
	assert.Equal(t, ServiceCallStatus("timeout"), ServiceCallStatusTimeout)
	assert.Equal(t, ServiceCallStatus("retrying"), ServiceCallStatusRetrying)
}

func TestNewServiceCall(t *testing.T) {
	params := map[string]any{
		"height": 50,
		"speed":  10,
	}
	paramsJSON, _ := json.Marshal(params)

	call := NewServiceCall("DEVICE001", "dji", "takeoff", paramsJSON)

	assert.Equal(t, "DEVICE001", call.DeviceSN)
	assert.Equal(t, "dji", call.Vendor)
	assert.Equal(t, "takeoff", call.Method)
	assert.Equal(t, json.RawMessage(paramsJSON), call.Params)
	assert.Equal(t, ServiceCallTypeCommand, call.CallType)
	assert.Equal(t, ServiceCallStatusPending, call.Status)
	assert.Equal(t, 3, call.MaxRetries)
	assert.False(t, call.CreatedAt.IsZero())
}

func TestServiceCall(t *testing.T) {
	now := time.Now()
	sentAt := now.Add(time.Second)
	completedAt := now.Add(2 * time.Second)
	paramsJSON, _ := json.Marshal(map[string]any{"height": 50})

	call := &ServiceCall{
		ID:          "call-001",
		DeviceSN:    "DEVICE001",
		Vendor:      "dji",
		Method:      "takeoff",
		Params:      paramsJSON,
		CallType:    ServiceCallTypeCommand,
		Status:      ServiceCallStatusSuccess,
		TID:         "tid-001",
		BID:         "bid-001",
		CreatedAt:   now,
		SentAt:      &sentAt,
		CompletedAt: &completedAt,
		RetryCount:  1,
		MaxRetries:  3,
	}

	assert.Equal(t, "call-001", call.ID)
	assert.Equal(t, "DEVICE001", call.DeviceSN)
	assert.Equal(t, "dji", call.Vendor)
	assert.Equal(t, "takeoff", call.Method)
	assert.Equal(t, ServiceCallStatusSuccess, call.Status)
	assert.Equal(t, 1, call.RetryCount)
}

func TestDispatchResult(t *testing.T) {
	result := &DispatchResult{
		Success:   true,
		MessageID: "msg-001",
		SentAt:    time.Now(),
	}

	assert.True(t, result.Success)
	assert.Equal(t, "msg-001", result.MessageID)
	assert.Nil(t, result.Error)
}
