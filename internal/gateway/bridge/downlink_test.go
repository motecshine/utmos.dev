package bridge

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultDownlinkBridgeConfig(t *testing.T) {
	config := DefaultDownlinkBridgeConfig()
	assert.Equal(t, "iot.topic", config.Exchange)
	assert.Equal(t, "iot.gateway.downlink", config.Queue)
	assert.Equal(t, "iot.raw.*.downlink", config.RoutingKey)
}

func TestNewDownlinkBridge(t *testing.T) {
	t.Run("with nil config", func(t *testing.T) {
		bridge := NewDownlinkBridge(nil, nil, nil, nil)
		require.NotNil(t, bridge)
		assert.Equal(t, "iot.topic", bridge.exchange)
		assert.Equal(t, "iot.gateway.downlink", bridge.queue)
	})

	t.Run("with custom config", func(t *testing.T) {
		config := &DownlinkBridgeConfig{
			Exchange:   "custom.exchange",
			Queue:      "custom.queue",
			RoutingKey: "custom.routing.key",
		}
		bridge := NewDownlinkBridge(nil, nil, config, nil)
		require.NotNil(t, bridge)
		assert.Equal(t, "custom.exchange", bridge.exchange)
		assert.Equal(t, "custom.queue", bridge.queue)
		assert.Equal(t, "custom.routing.key", bridge.routingKey)
	})
}

func TestDownlinkBridge_Bridge_NoClient(t *testing.T) {
	bridge := NewDownlinkBridge(nil, nil, nil, nil)

	msg := &RawDownlinkMessage{
		DeviceSN: "device-001",
		Topic:    "thing/product/device-001/services_reply/test",
		Payload:  json.RawMessage(`{"result": 0}`),
		QoS:      1,
		TraceID:  "trace-123",
		SpanID:   "span-456",
	}

	err := bridge.Bridge(context.Background(), msg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "MQTT client not initialized")
}

func TestDownlinkBridge_Start_NoConsumer(t *testing.T) {
	bridge := NewDownlinkBridge(nil, nil, nil, nil)

	err := bridge.Start(context.Background())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "subscriber not initialized")
}

func TestDownlinkBridge_IsRunning(t *testing.T) {
	bridge := NewDownlinkBridge(nil, nil, nil, nil)
	assert.False(t, bridge.IsRunning())
}

func TestDownlinkBridge_Stop(t *testing.T) {
	bridge := NewDownlinkBridge(nil, nil, nil, nil)
	bridge.mu.Lock()
	bridge.running = true
	bridge.mu.Unlock()

	bridge.Stop()
	assert.False(t, bridge.IsRunning())
}

func TestGetDownlinkRoutingKey(t *testing.T) {
	tests := []struct {
		vendor   string
		expected string
	}{
		{"dji", "iot.raw.dji.downlink"},
		{"vendor", "iot.raw.vendor.downlink"},
		{"test", "iot.raw.test.downlink"},
	}

	for _, tt := range tests {
		t.Run(tt.vendor, func(t *testing.T) {
			result := GetDownlinkRoutingKey(tt.vendor)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseRawDownlinkMessage(t *testing.T) {
	t.Run("valid message", func(t *testing.T) {
		data := []byte(`{
			"device_sn": "device-001",
			"topic": "thing/product/device-001/services_reply/test",
			"payload": {"result": 0},
			"qos": 1,
			"retained": false,
			"trace_id": "trace-123",
			"span_id": "span-456"
		}`)

		msg, err := ParseRawDownlinkMessage(data)
		require.NoError(t, err)
		assert.Equal(t, "device-001", msg.DeviceSN)
		assert.Equal(t, "thing/product/device-001/services_reply/test", msg.Topic)
		assert.Equal(t, 1, msg.QoS)
		assert.False(t, msg.Retained)
		assert.Equal(t, "trace-123", msg.TraceID)
	})

	t.Run("invalid JSON", func(t *testing.T) {
		data := []byte(`invalid json`)
		msg, err := ParseRawDownlinkMessage(data)
		assert.Error(t, err)
		assert.Nil(t, msg)
	})
}

func TestNewRawDownlinkMessage(t *testing.T) {
	payload := json.RawMessage(`{"result": 0}`)
	msg := NewRawDownlinkMessage("device-001", "thing/product/device-001/services_reply/test", payload, 1, false, "trace-123", "span-456")

	assert.Equal(t, "device-001", msg.DeviceSN)
	assert.Equal(t, "thing/product/device-001/services_reply/test", msg.Topic)
	assert.Equal(t, payload, msg.Payload)
	assert.Equal(t, 1, msg.QoS)
	assert.False(t, msg.Retained)
	assert.Equal(t, "trace-123", msg.TraceID)
	assert.Equal(t, "span-456", msg.SpanID)
}

func TestRawDownlinkMessage_JSON(t *testing.T) {
	original := &RawDownlinkMessage{
		DeviceSN: "device-001",
		Topic:    "thing/product/device-001/services_reply/test",
		Payload:  json.RawMessage(`{"result": 0}`),
		QoS:      1,
		Retained: false,
		TraceID:  "trace-123",
		SpanID:   "span-456",
	}

	data, err := json.Marshal(original)
	require.NoError(t, err)

	var decoded RawDownlinkMessage
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, original.DeviceSN, decoded.DeviceSN)
	assert.Equal(t, original.Topic, decoded.Topic)
	assert.Equal(t, original.QoS, decoded.QoS)
	assert.Equal(t, original.Retained, decoded.Retained)
	assert.Equal(t, original.TraceID, decoded.TraceID)
	assert.Equal(t, original.SpanID, decoded.SpanID)
}

func TestDownlinkBridge_HandleMessage(t *testing.T) {
	bridge := NewDownlinkBridge(nil, nil, nil, nil)

	t.Run("invalid message", func(t *testing.T) {
		err := bridge.HandleMessage(context.Background(), []byte(`invalid`))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse downlink message")
	})

	t.Run("valid message but no client", func(t *testing.T) {
		data := []byte(`{
			"device_sn": "device-001",
			"topic": "thing/product/device-001/services_reply/test",
			"payload": {"result": 0},
			"qos": 1
		}`)
		err := bridge.HandleMessage(context.Background(), data)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "MQTT client not initialized")
	})
}
