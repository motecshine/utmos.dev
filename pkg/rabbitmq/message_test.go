package rabbitmq

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewStandardMessage(t *testing.T) {
	data := map[string]interface{}{
		"temperature": 25.5,
		"humidity":    60,
	}

	msg, err := NewStandardMessage("iot-gateway", "property.report", "device-001", data)
	require.NoError(t, err)

	assert.NotEmpty(t, msg.TID)
	assert.NotEmpty(t, msg.BID)
	assert.Equal(t, "iot-gateway", msg.Service)
	assert.Equal(t, "property.report", msg.Action)
	assert.Equal(t, "device-001", msg.DeviceSN)
	assert.True(t, msg.Timestamp > 0)

	// Verify data is properly marshaled
	var parsedData map[string]interface{}
	err = json.Unmarshal(msg.Data, &parsedData)
	require.NoError(t, err)
	assert.Equal(t, 25.5, parsedData["temperature"])
}

func TestStandardMessage_Validate(t *testing.T) {
	tests := []struct {
		name    string
		msg     *StandardMessage
		wantErr bool
	}{
		{
			name: "valid message",
			msg: &StandardMessage{
				TID:       "tid-123",
				BID:       "bid-456",
				Timestamp: time.Now().UnixMilli(),
				Service:   "iot-gateway",
				Action:    "property.report",
				DeviceSN:  "device-001",
				Data:      json.RawMessage(`{"temp": 25}`),
			},
			wantErr: false,
		},
		{
			name: "missing TID",
			msg: &StandardMessage{
				BID:       "bid-456",
				Timestamp: time.Now().UnixMilli(),
				Service:   "iot-gateway",
				Action:    "property.report",
				DeviceSN:  "device-001",
				Data:      json.RawMessage(`{}`),
			},
			wantErr: true,
		},
		{
			name: "missing BID",
			msg: &StandardMessage{
				TID:       "tid-123",
				Timestamp: time.Now().UnixMilli(),
				Service:   "iot-gateway",
				Action:    "property.report",
				DeviceSN:  "device-001",
				Data:      json.RawMessage(`{}`),
			},
			wantErr: true,
		},
		{
			name: "missing Service",
			msg: &StandardMessage{
				TID:       "tid-123",
				BID:       "bid-456",
				Timestamp: time.Now().UnixMilli(),
				Action:    "property.report",
				DeviceSN:  "device-001",
				Data:      json.RawMessage(`{}`),
			},
			wantErr: true,
		},
		{
			name: "missing Action",
			msg: &StandardMessage{
				TID:       "tid-123",
				BID:       "bid-456",
				Timestamp: time.Now().UnixMilli(),
				Service:   "iot-gateway",
				DeviceSN:  "device-001",
				Data:      json.RawMessage(`{}`),
			},
			wantErr: true,
		},
		{
			name: "missing DeviceSN",
			msg: &StandardMessage{
				TID:       "tid-123",
				BID:       "bid-456",
				Timestamp: time.Now().UnixMilli(),
				Service:   "iot-gateway",
				Action:    "property.report",
				Data:      json.RawMessage(`{}`),
			},
			wantErr: true,
		},
		{
			name: "zero timestamp",
			msg: &StandardMessage{
				TID:      "tid-123",
				BID:      "bid-456",
				Service:  "iot-gateway",
				Action:   "property.report",
				DeviceSN: "device-001",
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

func TestMessageHeader(t *testing.T) {
	header := MessageHeader{
		Traceparent: "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01",
		Tracestate:  "congo=t61rcWkgMzE",
		MessageType: "property",
		Vendor:      "dji",
	}

	assert.Equal(t, "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01", header.Traceparent)
	assert.Equal(t, "congo=t61rcWkgMzE", header.Tracestate)
	assert.Equal(t, "property", header.MessageType)
	assert.Equal(t, "dji", header.Vendor)
}

func TestStandardMessage_Marshal(t *testing.T) {
	msg := &StandardMessage{
		TID:       "tid-123",
		BID:       "bid-456",
		Timestamp: 1234567890123,
		Service:   "iot-gateway",
		Action:    "property.report",
		DeviceSN:  "device-001",
		Data:      json.RawMessage(`{"temperature": 25.5}`),
	}

	data, err := json.Marshal(msg)
	require.NoError(t, err)

	var parsed StandardMessage
	err = json.Unmarshal(data, &parsed)
	require.NoError(t, err)

	assert.Equal(t, msg.TID, parsed.TID)
	assert.Equal(t, msg.BID, parsed.BID)
	assert.Equal(t, msg.Timestamp, parsed.Timestamp)
	assert.Equal(t, msg.Service, parsed.Service)
	assert.Equal(t, msg.Action, parsed.Action)
	assert.Equal(t, msg.DeviceSN, parsed.DeviceSN)
}
