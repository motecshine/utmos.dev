package adapter

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProtocolMessage_Validate(t *testing.T) {
	tests := []struct {
		name    string
		msg     *ProtocolMessage
		wantErr bool
	}{
		{
			name: "valid message",
			msg: &ProtocolMessage{
				Vendor:      "dji",
				Topic:       "thing/product/ABC123/osd",
				DeviceSN:    "ABC123",
				MessageType: MessageTypeProperty,
				Data:        json.RawMessage(`{"latitude": 39.9042}`),
			},
			wantErr: false,
		},
		{
			name: "missing vendor",
			msg: &ProtocolMessage{
				Topic:       "thing/product/ABC123/osd",
				DeviceSN:    "ABC123",
				MessageType: MessageTypeProperty,
				Data:        json.RawMessage(`{}`),
			},
			wantErr: true,
		},
		{
			name: "missing topic",
			msg: &ProtocolMessage{
				Vendor:      "dji",
				DeviceSN:    "ABC123",
				MessageType: MessageTypeProperty,
				Data:        json.RawMessage(`{}`),
			},
			wantErr: true,
		},
		{
			name: "missing device_sn",
			msg: &ProtocolMessage{
				Vendor:      "dji",
				Topic:       "thing/product/ABC123/osd",
				MessageType: MessageTypeProperty,
				Data:        json.RawMessage(`{}`),
			},
			wantErr: true,
		},
		{
			name: "missing message_type",
			msg: &ProtocolMessage{
				Vendor:   "dji",
				Topic:    "thing/product/ABC123/osd",
				DeviceSN: "ABC123",
				Data:     json.RawMessage(`{}`),
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

func TestProtocolMessage_JSON(t *testing.T) {
	msg := &ProtocolMessage{
		Vendor:      "dji",
		Topic:       "thing/product/ABC123/osd",
		DeviceSN:    "ABC123",
		GatewaySN:   "GW456",
		MessageType: MessageTypeProperty,
		Method:      "thing.property.post",
		TID:         "tid-123",
		BID:         "bid-456",
		Timestamp:   1234567890123,
		Data:        json.RawMessage(`{"latitude": 39.9042}`),
	}

	// Test marshal
	data, err := json.Marshal(msg)
	require.NoError(t, err)

	// Test unmarshal
	var parsed ProtocolMessage
	err = json.Unmarshal(data, &parsed)
	require.NoError(t, err)

	assert.Equal(t, msg.Vendor, parsed.Vendor)
	assert.Equal(t, msg.Topic, parsed.Topic)
	assert.Equal(t, msg.DeviceSN, parsed.DeviceSN)
	assert.Equal(t, msg.GatewaySN, parsed.GatewaySN)
	assert.Equal(t, msg.MessageType, parsed.MessageType)
	assert.Equal(t, msg.Method, parsed.Method)
	assert.Equal(t, msg.TID, parsed.TID)
	assert.Equal(t, msg.BID, parsed.BID)
	assert.Equal(t, msg.Timestamp, parsed.Timestamp)
}

func TestMessageType_String(t *testing.T) {
	tests := []struct {
		mt       MessageType
		expected string
	}{
		{MessageTypeProperty, "property"},
		{MessageTypeEvent, "event"},
		{MessageTypeService, "service"},
		{MessageTypeStatus, "status"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.mt.String())
		})
	}
}
