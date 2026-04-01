package router

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/utmos/utmos/pkg/adapter"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	assert.Equal(t, "iot.topic", config.Exchange)
	assert.True(t, config.EnableWSRouting)
	assert.True(t, config.EnableAPIRouting)
}

func TestNewRouter(t *testing.T) {
	t.Run("with nil config", func(t *testing.T) {
		router := NewRouter(nil, nil, nil)
		require.NotNil(t, router)
		assert.NotNil(t, router.config)
	})

	t.Run("with custom config", func(t *testing.T) {
		config := &Config{
			Exchange:         "custom.exchange",
			EnableWSRouting:  false,
			EnableAPIRouting: true,
		}

		router := NewRouter(nil, config, nil)
		require.NotNil(t, router)
		assert.Equal(t, "custom.exchange", router.config.Exchange)
		assert.False(t, router.config.EnableWSRouting)
	})
}

func TestRouter_StartStop(t *testing.T) {
	router := NewRouter(nil, nil, nil)

	t.Run("start", func(t *testing.T) {
		err := router.Start()
		assert.NoError(t, err)
		assert.True(t, router.IsRunning())
	})

	t.Run("double start", func(t *testing.T) {
		err := router.Start()
		assert.Error(t, err)
	})

	t.Run("stop", func(t *testing.T) {
		err := router.Stop()
		assert.NoError(t, err)
		assert.False(t, router.IsRunning())
	})

	t.Run("double stop", func(t *testing.T) {
		err := router.Stop()
		assert.NoError(t, err)
	})
}

func TestRouter_Route_NilMessage(t *testing.T) {
	router := NewRouter(nil, nil, nil)

	err := router.Route(context.Background(), nil)
	assert.Error(t, err)
}

func TestRouter_Route_NoPublisher(t *testing.T) {
	router := NewRouter(nil, nil, nil)

	msg := &adapter.ProcessedMessage{
		DeviceSN: "DEVICE001",
		Vendor:   "dji",
	}

	err := router.Route(context.Background(), msg)
	assert.Error(t, err)
}

func TestRouter_GetWSRoutingKey(t *testing.T) {
	router := NewRouter(nil, nil, nil)

	testCases := []struct {
		vendor   string
		msgType  adapter.MessageType
		expected string
	}{
		{"dji", adapter.MessageTypeProperty, "iot.dji.ws.property.report"},
		{"dji", adapter.MessageTypeEvent, "iot.dji.ws.event.notify"},
		{"dji", adapter.MessageTypeStatus, "iot.dji.ws.status.report"},
		{"generic", adapter.MessageTypeProperty, "iot.generic.ws.property.report"},
		{"tuya", adapter.MessageTypeEvent, "iot.tuya.ws.event.notify"},
		{"dji", adapter.MessageType("unknown"), "iot.dji.ws.property.report"}, // default
	}

	for _, tc := range testCases {
		t.Run(tc.vendor+"_"+string(tc.msgType), func(t *testing.T) {
			result := router.getWSRoutingKey(tc.msgType, tc.vendor)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestRouter_GetAPIRoutingKey(t *testing.T) {
	router := NewRouter(nil, nil, nil)

	testCases := []struct {
		vendor   string
		msgType  adapter.MessageType
		expected string
	}{
		{"dji", adapter.MessageTypeProperty, "iot.dji.api.property.report"},
		{"dji", adapter.MessageTypeEvent, "iot.dji.api.event.notify"},
		{"generic", adapter.MessageTypeProperty, "iot.generic.api.property.report"},
		{"tuya", adapter.MessageTypeEvent, "iot.tuya.api.event.notify"},
		{"dji", adapter.MessageType("unknown"), "iot.dji.api.property.report"}, // default
	}

	for _, tc := range testCases {
		t.Run(tc.vendor+"_"+string(tc.msgType), func(t *testing.T) {
			result := router.getAPIRoutingKey(tc.msgType, tc.vendor)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestRouter_GetAction(t *testing.T) {
	router := NewRouter(nil, nil, nil)

	testCases := []struct {
		msgType  adapter.MessageType
		expected string
	}{
		{adapter.MessageTypeProperty, "property.report"},
		{adapter.MessageTypeEvent, "event.notify"},
		{adapter.MessageTypeService, "service.reply"},
		{adapter.MessageTypeStatus, "status.report"},
		{adapter.MessageType("unknown"), "message.report"},
	}

	for _, tc := range testCases {
		t.Run(string(tc.msgType), func(t *testing.T) {
			result := router.getAction(tc.msgType)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestRouter_CreateStandardMessage(t *testing.T) {
	router := NewRouter(nil, nil, nil)

	msg := &adapter.ProcessedMessage{
		MessageType: adapter.MessageTypeProperty,
		DeviceSN:    "DEVICE001",
		Vendor:      "dji",
		Properties: map[string]any{
			"temperature": 25.5,
		},
		Events:    []adapter.Event{},
		Timestamp: 1704067200000,
	}

	stdMsg, err := router.createStandardMessage(msg)
	require.NoError(t, err)
	require.NotNil(t, stdMsg)

	assert.Equal(t, "DEVICE001", stdMsg.DeviceSN)
	assert.Equal(t, "iot-uplink", stdMsg.Service)
	assert.Equal(t, "property.report", stdMsg.Action)
	assert.NotNil(t, stdMsg.ProtocolMeta)
	assert.Equal(t, "dji", stdMsg.ProtocolMeta.Vendor)
}

func TestNewMultiRouter(t *testing.T) {
	router := NewMultiRouter(nil)
	require.NotNil(t, router)
	assert.Empty(t, router.routers)
}

func TestMultiRouter_AddRouter(t *testing.T) {
	router := NewMultiRouter(nil)

	router.AddRouter(func(ctx context.Context, msg *adapter.ProcessedMessage) error {
		return nil
	})

	assert.Len(t, router.routers, 1)
}

func TestMultiRouter_Route(t *testing.T) {
	router := NewMultiRouter(nil)

	var callOrder []int

	router.AddRouter(func(ctx context.Context, msg *adapter.ProcessedMessage) error {
		callOrder = append(callOrder, 1)
		return nil
	})

	router.AddRouter(func(ctx context.Context, msg *adapter.ProcessedMessage) error {
		callOrder = append(callOrder, 2)
		return nil
	})

	router.AddRouter(func(ctx context.Context, msg *adapter.ProcessedMessage) error {
		callOrder = append(callOrder, 3)
		return nil
	})

	msg := &adapter.ProcessedMessage{
		DeviceSN: "DEVICE001",
	}

	err := router.Route(context.Background(), msg)
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3}, callOrder)
}

func TestMultiRouter_Route_WithErrors(t *testing.T) {
	router := NewMultiRouter(nil)

	router.AddRouter(func(ctx context.Context, msg *adapter.ProcessedMessage) error {
		return nil
	})

	router.AddRouter(func(ctx context.Context, msg *adapter.ProcessedMessage) error {
		return errors.New("router 2 error")
	})

	router.AddRouter(func(ctx context.Context, msg *adapter.ProcessedMessage) error {
		return errors.New("router 3 error")
	})

	msg := &adapter.ProcessedMessage{
		DeviceSN: "DEVICE001",
	}

	err := router.Route(context.Background(), msg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "router 2 error")
	assert.Contains(t, err.Error(), "router 3 error")
}

func TestServiceConstants(t *testing.T) {
	assert.Equal(t, "ws", ServiceWS)
	assert.Equal(t, "api", ServiceAPI)
}
