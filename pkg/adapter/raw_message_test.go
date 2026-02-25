package adapter

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRawMessage(t *testing.T) {
	payload := []byte(`{"tid": "123", "data": {"temp": 25}}`)
	msg := NewRawMessage("dji", "thing/product/ABC123/osd", payload, 1)

	assert.Equal(t, "dji", msg.Vendor)
	assert.Equal(t, "thing/product/ABC123/osd", msg.Topic)
	assert.Equal(t, payload, msg.Payload)
	assert.Equal(t, 1, msg.QoS)
	assert.True(t, msg.Timestamp > 0)
	assert.NotNil(t, msg.Headers)
}

func TestRawMessage_Validate(t *testing.T) {
	tests := []struct {
		name    string
		msg     *RawMessage
		wantErr bool
	}{
		{
			name: "valid message",
			msg: &RawMessage{
				Vendor:    "dji",
				Topic:     "thing/product/ABC123/osd",
				Payload:   []byte(`{"data": {}}`),
				Timestamp: time.Now().UnixMilli(),
			},
			wantErr: false,
		},
		{
			name: "missing vendor",
			msg: &RawMessage{
				Topic:     "thing/product/ABC123/osd",
				Payload:   []byte(`{"data": {}}`),
				Timestamp: time.Now().UnixMilli(),
			},
			wantErr: true,
		},
		{
			name: "missing topic",
			msg: &RawMessage{
				Vendor:    "dji",
				Payload:   []byte(`{"data": {}}`),
				Timestamp: time.Now().UnixMilli(),
			},
			wantErr: true,
		},
		{
			name: "empty payload",
			msg: &RawMessage{
				Vendor:    "dji",
				Topic:     "thing/product/ABC123/osd",
				Payload:   []byte{},
				Timestamp: time.Now().UnixMilli(),
			},
			wantErr: true,
		},
		{
			name: "nil payload",
			msg: &RawMessage{
				Vendor:    "dji",
				Topic:     "thing/product/ABC123/osd",
				Payload:   nil,
				Timestamp: time.Now().UnixMilli(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRawMessage_JSON(t *testing.T) {
	original := &RawMessage{
		Vendor:    "dji",
		Topic:     "thing/product/ABC123/osd",
		Payload:   []byte(`{"latitude": 39.9042}`),
		QoS:       1,
		Timestamp: 1234567890123,
		Headers: map[string]string{
			"traceparent": "00-abc123-def456-01",
		},
	}

	// Test marshal
	data, err := json.Marshal(original)
	require.NoError(t, err)

	// Test unmarshal
	var parsed RawMessage
	err = json.Unmarshal(data, &parsed)
	require.NoError(t, err)

	assert.Equal(t, original.Vendor, parsed.Vendor)
	assert.Equal(t, original.Topic, parsed.Topic)
	assert.Equal(t, original.QoS, parsed.QoS)
	assert.Equal(t, original.Timestamp, parsed.Timestamp)
	assert.Equal(t, original.Headers["traceparent"], parsed.Headers["traceparent"])
}

func TestRawMessage_WithHeader(t *testing.T) {
	msg := NewRawMessage("dji", "thing/product/ABC123/osd", []byte(`{}`), 0)
	msg.WithHeader("traceparent", "00-abc-def-01")
	msg.WithHeader("vendor", "dji")

	assert.Equal(t, "00-abc-def-01", msg.Headers["traceparent"])
	assert.Equal(t, "dji", msg.Headers["vendor"])
}
