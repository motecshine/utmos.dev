package bridge

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/utmos/utmos/internal/gateway/mqtt"
)

func TestDefaultUplinkBridgeConfig(t *testing.T) {
	config := DefaultUplinkBridgeConfig()
	assert.Equal(t, "iot.topic", config.Exchange)
}

func TestNewUplinkBridge(t *testing.T) {
	t.Run("with nil config", func(t *testing.T) {
		bridge := NewUplinkBridge(nil, nil, nil)
		require.NotNil(t, bridge)
		assert.Equal(t, "iot.topic", bridge.exchange)
	})

	t.Run("with custom config", func(t *testing.T) {
		config := &UplinkBridgeConfig{
			Exchange: "custom.exchange",
		}
		bridge := NewUplinkBridge(nil, config, nil)
		require.NotNil(t, bridge)
		assert.Equal(t, "custom.exchange", bridge.exchange)
	})
}

func TestUplinkBridge_Bridge_NoPublisher(t *testing.T) {
	bridge := NewUplinkBridge(nil, nil, nil)

	msg := &mqtt.Message{
		Topic:     "thing/product/device-001/osd",
		Payload:   json.RawMessage(`{"data": "test"}`),
		Timestamp: time.Now(),
		TraceID:   "trace-123",
		SpanID:    "span-456",
	}
	topicInfo := mqtt.ParseTopic(msg.Topic)

	err := bridge.Bridge(context.Background(), msg, topicInfo)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "publisher not initialized")
}

func TestGetUplinkRoutingKey(t *testing.T) {
	tests := []struct {
		vendor   string
		expected string
	}{
		{"dji", "iot.raw.dji.uplink"},
		{"vendor", "iot.raw.vendor.uplink"},
		{"test", "iot.raw.test.uplink"},
	}

	for _, tt := range tests {
		t.Run(tt.vendor, func(t *testing.T) {
			result := GetUplinkRoutingKey(tt.vendor)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseRawUplinkMessage(t *testing.T) {
	t.Run("valid message", func(t *testing.T) {
		data := []byte(`{
			"vendor": "dji",
			"topic": "thing/product/device-001/osd",
			"payload": {"key": "value"},
			"qos": 1,
			"timestamp": 1234567890,
			"trace_id": "trace-123",
			"span_id": "span-456"
		}`)

		msg, err := ParseRawUplinkMessage(data)
		require.NoError(t, err)
		assert.Equal(t, "dji", msg.Vendor)
		assert.Equal(t, "thing/product/device-001/osd", msg.Topic)
		assert.Equal(t, 1, msg.QoS)
		assert.Equal(t, int64(1234567890), msg.Timestamp)
		assert.Equal(t, "trace-123", msg.TraceID)
	})

	t.Run("invalid JSON", func(t *testing.T) {
		data := []byte(`invalid json`)
		msg, err := ParseRawUplinkMessage(data)
		assert.Error(t, err)
		assert.Nil(t, msg)
	})
}

func TestNewRawUplinkMessage(t *testing.T) {
	payload := json.RawMessage(`{"data": "test"}`)
	msg := NewRawUplinkMessage("dji", "thing/product/device-001/osd", payload, 1)

	assert.Equal(t, "dji", msg.Vendor)
	assert.Equal(t, "thing/product/device-001/osd", msg.Topic)
	assert.Equal(t, payload, msg.Payload)
	assert.Equal(t, 1, msg.QoS)
	assert.NotEmpty(t, msg.TraceID)
	assert.NotEmpty(t, msg.SpanID)
	assert.Greater(t, msg.Timestamp, int64(0))
}

func TestUplinkBridge_CreateProcessor(t *testing.T) {
	bridge := NewUplinkBridge(nil, nil, nil)
	processor := bridge.CreateProcessor("thing/product/#")

	assert.NotNil(t, processor)
	assert.Equal(t, "thing/product/#", processor.Pattern())
}

func TestRawUplinkMessage_JSON(t *testing.T) {
	original := &RawUplinkMessage{
		Vendor:    "dji",
		Topic:     "thing/product/device-001/osd",
		Payload:   json.RawMessage(`{"key": "value"}`),
		QoS:       1,
		Timestamp: 1234567890,
		TraceID:   "trace-123",
		SpanID:    "span-456",
	}

	data, err := json.Marshal(original)
	require.NoError(t, err)

	var decoded RawUplinkMessage
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, original.Vendor, decoded.Vendor)
	assert.Equal(t, original.Topic, decoded.Topic)
	assert.Equal(t, original.QoS, decoded.QoS)
	assert.Equal(t, original.Timestamp, decoded.Timestamp)
	assert.Equal(t, original.TraceID, decoded.TraceID)
	assert.Equal(t, original.SpanID, decoded.SpanID)
}
