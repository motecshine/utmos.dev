package mqtt

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseTopic(t *testing.T) {
	tests := []struct {
		name     string
		topic    string
		expected *TopicInfo
	}{
		{
			name:  "DJI thing topic",
			topic: "thing/product/device-001/osd",
			expected: &TopicInfo{
				Vendor:    "dji",
				ProductID: "product",
				DeviceSN:  "device-001",
				Service:   "osd",
				Raw:       "thing/product/device-001/osd",
			},
		},
		{
			name:  "DJI sys topic",
			topic: "sys/product/device-001/status",
			expected: &TopicInfo{
				Vendor:    "dji",
				ProductID: "product",
				DeviceSN:  "device-001",
				Service:   "status",
				Raw:       "sys/product/device-001/status",
			},
		},
		{
			name:  "DJI topic with method",
			topic: "thing/product/device-001/services/flighttask_prepare",
			expected: &TopicInfo{
				Vendor:    "dji",
				ProductID: "product",
				DeviceSN:  "device-001",
				Service:   "services",
				Method:    "flighttask_prepare",
				Raw:       "thing/product/device-001/services/flighttask_prepare",
			},
		},
		{
			name:  "Generic vendor topic",
			topic: "vendor/thing/product/device-001/telemetry",
			expected: &TopicInfo{
				Vendor:    "vendor",
				ProductID: "product",
				DeviceSN:  "device-001",
				Service:   "telemetry",
				Raw:       "vendor/thing/product/device-001/telemetry",
			},
		},
		{
			name:  "Short topic",
			topic: "a/b",
			expected: &TopicInfo{
				Vendor: "",
				Raw:    "a/b",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseTopic(tt.topic)
			assert.Equal(t, tt.expected.Vendor, result.Vendor)
			assert.Equal(t, tt.expected.ProductID, result.ProductID)
			assert.Equal(t, tt.expected.DeviceSN, result.DeviceSN)
			assert.Equal(t, tt.expected.Service, result.Service)
			assert.Equal(t, tt.expected.Method, result.Method)
			assert.Equal(t, tt.expected.Raw, result.Raw)
		})
	}
}

func TestMatchTopic(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		topic   string
		match   bool
	}{
		{
			name:    "exact match",
			pattern: "thing/product/device-001/osd",
			topic:   "thing/product/device-001/osd",
			match:   true,
		},
		{
			name:    "no match",
			pattern: "thing/product/device-001/osd",
			topic:   "thing/product/device-002/osd",
			match:   false,
		},
		{
			name:    "single level wildcard",
			pattern: "thing/product/+/osd",
			topic:   "thing/product/device-001/osd",
			match:   true,
		},
		{
			name:    "multi level wildcard",
			pattern: "thing/product/#",
			topic:   "thing/product/device-001/osd",
			match:   true,
		},
		{
			name:    "multi level wildcard at end",
			pattern: "thing/#",
			topic:   "thing/product/device-001/osd/data",
			match:   true,
		},
		{
			name:    "multiple single level wildcards",
			pattern: "thing/+/+/osd",
			topic:   "thing/product/device-001/osd",
			match:   true,
		},
		{
			name:    "wildcard no match",
			pattern: "thing/product/+/status",
			topic:   "thing/product/device-001/osd",
			match:   false,
		},
		{
			name:    "all wildcard",
			pattern: "#",
			topic:   "any/topic/here",
			match:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := matchTopic(tt.pattern, tt.topic)
			assert.Equal(t, tt.match, result)
		})
	}
}

func TestNewHandler(t *testing.T) {
	handler := NewHandler(nil)
	require.NotNil(t, handler)
	assert.NotNil(t, handler.logger)
	assert.NotNil(t, handler.processors)
}

func TestHandler_RegisterProcessor(t *testing.T) {
	handler := NewHandler(nil)

	processed := false
	processor := NewSimpleProcessor("thing/product/#", func(msg *Message, topicInfo *TopicInfo) error {
		processed = true
		return nil
	})

	handler.RegisterProcessor(processor)

	handler.mu.RLock()
	_, exists := handler.processors["thing/product/#"]
	handler.mu.RUnlock()

	assert.True(t, exists)
	assert.False(t, processed) // Not processed yet
}

func TestHandler_UnregisterProcessor(t *testing.T) {
	handler := NewHandler(nil)

	processor := NewSimpleProcessor("thing/product/#", func(msg *Message, topicInfo *TopicInfo) error {
		return nil
	})

	handler.RegisterProcessor(processor)
	handler.UnregisterProcessor("thing/product/#")

	handler.mu.RLock()
	_, exists := handler.processors["thing/product/#"]
	handler.mu.RUnlock()

	assert.False(t, exists)
}

func TestSimpleProcessor(t *testing.T) {
	called := false
	var receivedMsg *Message
	var receivedInfo *TopicInfo

	processor := NewSimpleProcessor("test/+/topic", func(msg *Message, topicInfo *TopicInfo) error {
		called = true
		receivedMsg = msg
		receivedInfo = topicInfo
		return nil
	})

	assert.Equal(t, "test/+/topic", processor.Pattern())

	msg := &Message{
		Topic:   "test/device/topic",
		Payload: json.RawMessage(`{"key": "value"}`),
	}
	topicInfo := ParseTopic(msg.Topic)

	err := processor.Process(msg, topicInfo)
	require.NoError(t, err)
	assert.True(t, called)
	assert.Equal(t, msg, receivedMsg)
	assert.Equal(t, topicInfo, receivedInfo)
}

func TestMessage_JSON(t *testing.T) {
	msg := &Message{
		Topic:     "test/topic",
		Payload:   json.RawMessage(`{"data": "test"}`),
		QoS:       1,
		Retained:  false,
		MessageID: 123,
	}

	data, err := json.Marshal(msg)
	require.NoError(t, err)

	var decoded Message
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, msg.Topic, decoded.Topic)
	assert.Equal(t, msg.QoS, decoded.QoS)
	assert.Equal(t, msg.MessageID, decoded.MessageID)
}
