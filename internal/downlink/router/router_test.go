package router

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/utmos/utmos/internal/downlink/dispatcher"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	assert.Equal(t, RoutingKeyGatewayDownlink, config.DefaultRoutingKey)
	assert.True(t, config.EnableMetrics)
}

func TestNewRouter(t *testing.T) {
	t.Run("with config", func(t *testing.T) {
		config := &Config{
			DefaultRoutingKey: "custom.key",
		}
		router := NewRouter(nil, config, nil)

		require.NotNil(t, router)
		assert.Equal(t, "custom.key", router.config.DefaultRoutingKey)
	})

	t.Run("without config", func(t *testing.T) {
		router := NewRouter(nil, nil, nil)

		require.NotNil(t, router)
		assert.Equal(t, RoutingKeyGatewayDownlink, router.config.DefaultRoutingKey)
	})
}

func TestRouter_Route_NilCall(t *testing.T) {
	router := NewRouter(nil, nil, nil)

	_, err := router.Route(context.Background(), nil, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "service call is nil")
}

func TestRouter_Route_NoPublisher(t *testing.T) {
	router := NewRouter(nil, nil, nil)

	call := &dispatcher.ServiceCall{
		ID:       "call-001",
		DeviceSN: "DEVICE001",
		Vendor:   "dji",
		Method:   "takeoff",
	}

	_, err := router.Route(context.Background(), call, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "publisher not initialized")
}

func TestRouter_GetRoutingKey(t *testing.T) {
	router := NewRouter(nil, nil, nil)

	testCases := []struct {
		callType dispatcher.ServiceCallType
		expected string
	}{
		{dispatcher.ServiceCallTypeCommand, RoutingKeyGatewayCommand},
		{dispatcher.ServiceCallTypeProperty, RoutingKeyGatewayProperty},
		{dispatcher.ServiceCallTypeConfig, RoutingKeyGatewayDownlink}, // Default
	}

	for _, tc := range testCases {
		call := &dispatcher.ServiceCall{CallType: tc.callType}
		result := router.getRoutingKey(call)
		assert.Equal(t, tc.expected, result, "call type %s", tc.callType)
	}
}

func TestRouter_GetAction(t *testing.T) {
	router := NewRouter(nil, nil, nil)

	testCases := []struct {
		callType dispatcher.ServiceCallType
		expected string
	}{
		{dispatcher.ServiceCallTypeCommand, "command.send"},
		{dispatcher.ServiceCallTypeProperty, "property.set"},
		{dispatcher.ServiceCallTypeConfig, "config.update"},
	}

	for _, tc := range testCases {
		result := router.getAction(tc.callType)
		assert.Equal(t, tc.expected, result, "call type %s", tc.callType)
	}
}

func TestRouter_CreateGatewayMessage(t *testing.T) {
	router := NewRouter(nil, nil, nil)

	paramsJSON, _ := json.Marshal(map[string]any{"height": 50})
	call := &dispatcher.ServiceCall{
		ID:       "call-001",
		DeviceSN: "DEVICE001",
		Vendor:   "dji",
		Method:   "takeoff",
		Params:   paramsJSON,
		CallType: dispatcher.ServiceCallTypeCommand,
		TID:      "tid-001",
		BID:      "bid-001",
	}

	result := &dispatcher.DispatchResult{
		Success:   true,
		MessageID: "msg-001",
		SentAt:    time.Now(),
	}

	msg, err := router.createGatewayMessage(call, result)
	require.NoError(t, err)
	require.NotNil(t, msg)

	assert.Equal(t, "DEVICE001", msg.DeviceSN)
	assert.Equal(t, "tid-001", msg.TID)
	assert.Equal(t, "bid-001", msg.BID)
	assert.Equal(t, "iot-downlink", msg.Service)
	assert.Equal(t, "command.send", msg.Action)
	assert.NotNil(t, msg.ProtocolMeta)
	assert.Equal(t, "dji", msg.ProtocolMeta.Vendor)
	assert.Equal(t, "takeoff", msg.ProtocolMeta.Method)
}

func TestRouter_CreateGatewayMessage_WithError(t *testing.T) {
	router := NewRouter(nil, nil, nil)

	call := &dispatcher.ServiceCall{
		ID:       "call-001",
		DeviceSN: "DEVICE001",
		Vendor:   "dji",
		Method:   "takeoff",
		CallType: dispatcher.ServiceCallTypeCommand,
		TID:      "tid-001",
		BID:      "bid-001",
	}

	result := &dispatcher.DispatchResult{
		Success: false,
		Error:   assert.AnError,
		SentAt:  time.Now(),
	}

	msg, err := router.createGatewayMessage(call, result)
	require.NoError(t, err)
	require.NotNil(t, msg)

	var data gatewayPayload
	err = json.Unmarshal(msg.Data, &data)
	require.NoError(t, err)

	require.NotNil(t, data.DispatchResult)
	assert.False(t, data.DispatchResult.Success)
	assert.NotEmpty(t, data.DispatchResult.Error)
}

func TestRouter_Metrics(t *testing.T) {
	router := NewRouter(nil, nil, nil)

	// Initial metrics should be zero
	routed, failed := router.GetMetrics()
	assert.Equal(t, int64(0), routed)
	assert.Equal(t, int64(0), failed)

	// Increment counters
	router.incrementRouted()
	router.incrementRouted()
	router.incrementFailed()

	routed, failed = router.GetMetrics()
	assert.Equal(t, int64(2), routed)
	assert.Equal(t, int64(1), failed)

	// Reset metrics
	router.ResetMetrics()
	routed, failed = router.GetMetrics()
	assert.Equal(t, int64(0), routed)
	assert.Equal(t, int64(0), failed)
}

func TestRoutingKeyConstants(t *testing.T) {
	assert.Equal(t, "iot.gateway.downlink", RoutingKeyGatewayDownlink)
	assert.Equal(t, "iot.gateway.command", RoutingKeyGatewayCommand)
	assert.Equal(t, "iot.gateway.property", RoutingKeyGatewayProperty)
}

func TestNewBatchRouter(t *testing.T) {
	router := NewRouter(nil, nil, nil)
	batchRouter := NewBatchRouter(router, nil)

	require.NotNil(t, batchRouter)
	assert.Equal(t, router, batchRouter.router)
}

func TestBatchRouter_RouteBatch_NoPublisher(t *testing.T) {
	router := NewRouter(nil, nil, nil)
	batchRouter := NewBatchRouter(router, nil)

	calls := []*dispatcher.ServiceCall{
		{
			ID:       "call-001",
			DeviceSN: "DEVICE001",
			Vendor:   "dji",
			Method:   "takeoff",
		},
		{
			ID:       "call-002",
			DeviceSN: "DEVICE002",
			Vendor:   "dji",
			Method:   "land",
		},
	}

	result := batchRouter.RouteBatch(context.Background(), calls, nil)

	assert.Equal(t, 2, result.Total)
	assert.Equal(t, 0, result.Succeeded)
	assert.Equal(t, 2, result.Failed)
	assert.Len(t, result.Results, 2)
}

func TestRouteResult(t *testing.T) {
	result := &RouteResult{
		Success:    true,
		RoutingKey: "iot.gateway.command",
		Error:      nil,
	}

	assert.True(t, result.Success)
	assert.Equal(t, "iot.gateway.command", result.RoutingKey)
	assert.Nil(t, result.Error)
}

func TestBatchRouteResult(t *testing.T) {
	result := &BatchRouteResult{
		Total:     5,
		Succeeded: 3,
		Failed:    2,
		Results: []*RouteResult{
			{Success: true},
			{Success: true},
			{Success: true},
			{Success: false},
			{Success: false},
		},
	}

	assert.Equal(t, 5, result.Total)
	assert.Equal(t, 3, result.Succeeded)
	assert.Equal(t, 2, result.Failed)
	assert.Len(t, result.Results, 5)
}
