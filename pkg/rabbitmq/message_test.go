package rabbitmq

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewStandardMessage(t *testing.T) {
	data := map[string]any{
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
	var parsedData map[string]any
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

func TestProtocolMeta(t *testing.T) {
	qos := 1
	needReply := true

	meta := &ProtocolMeta{
		Vendor:        "dji",
		OriginalTopic: "thing/product/ABC123/osd",
		QoS:           &qos,
		Method:        "thing.property.post",
		NeedReply:     &needReply,
	}

	assert.Equal(t, "dji", meta.Vendor)
	assert.Equal(t, "thing/product/ABC123/osd", meta.OriginalTopic)
	assert.Equal(t, 1, *meta.QoS)
	assert.Equal(t, "thing.property.post", meta.Method)
	assert.True(t, *meta.NeedReply)
}

func TestStandardMessage_WithProtocolMeta(t *testing.T) {
	qos := 1

	msg := &StandardMessage{
		TID:       "tid-123",
		BID:       "bid-456",
		Timestamp: 1234567890123,
		Service:   "dji-adapter",
		Action:    "property.report",
		DeviceSN:  "ABC123",
		Data:      json.RawMessage(`{"latitude": 39.9042, "longitude": 116.4074}`),
		ProtocolMeta: &ProtocolMeta{
			Vendor:        "dji",
			OriginalTopic: "thing/product/ABC123/osd",
			QoS:           &qos,
		},
	}

	// Test serialization
	data, err := json.Marshal(msg)
	require.NoError(t, err)

	// Verify protocol_meta is included
	assert.Contains(t, string(data), "protocol_meta")
	assert.Contains(t, string(data), "original_topic")

	// Test deserialization
	var parsed StandardMessage
	err = json.Unmarshal(data, &parsed)
	require.NoError(t, err)

	require.NotNil(t, parsed.ProtocolMeta)
	assert.Equal(t, "dji", parsed.ProtocolMeta.Vendor)
	assert.Equal(t, "thing/product/ABC123/osd", parsed.ProtocolMeta.OriginalTopic)
	assert.Equal(t, 1, *parsed.ProtocolMeta.QoS)
}

func TestStandardMessage_WithoutProtocolMeta_BackwardCompatible(t *testing.T) {
	// Test that messages without protocol_meta still work (backward compatibility)
	jsonData := `{
		"tid": "tid-123",
		"bid": "bid-456",
		"timestamp": 1234567890123,
		"service": "iot-gateway",
		"action": "property.report",
		"device_sn": "device-001",
		"data": {"temp": 25}
	}`

	var msg StandardMessage
	err := json.Unmarshal([]byte(jsonData), &msg)
	require.NoError(t, err)

	assert.Equal(t, "tid-123", msg.TID)
	assert.Equal(t, "iot-gateway", msg.Service)
	assert.Nil(t, msg.ProtocolMeta) // Should be nil for backward compatibility
}

func TestFromBytes(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		wantErr bool
	}{
		{
			name: "valid message",
			data: []byte(`{
				"tid": "tid-123",
				"bid": "bid-456",
				"timestamp": 1234567890123,
				"service": "iot-gateway",
				"action": "property.report",
				"device_sn": "device-001",
				"data": {"temp": 25}
			}`),
			wantErr: false,
		},
		{
			name:    "invalid JSON",
			data:    []byte(`not json`),
			wantErr: true,
		},
		{
			name:    "empty data",
			data:    []byte(``),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg, err := FromBytes(tt.data)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, msg)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, msg)
			}
		})
	}
}

func TestStandardMessage_ToBytes(t *testing.T) {
	msg := &StandardMessage{
		TID:       "tid-123",
		BID:       "bid-456",
		Timestamp: 1234567890123,
		Service:   "iot-gateway",
		Action:    "property.report",
		DeviceSN:  "device-001",
		Data:      json.RawMessage(`{"temperature": 25.5}`),
	}

	data, err := msg.ToBytes()
	require.NoError(t, err)
	assert.NotEmpty(t, data)

	// Verify round-trip
	parsed, err := FromBytes(data)
	require.NoError(t, err)
	assert.Equal(t, msg.TID, parsed.TID)
	assert.Equal(t, msg.BID, parsed.BID)
}

func TestStandardMessage_GetSetData(t *testing.T) {
	msg := &StandardMessage{
		TID:       "tid-123",
		BID:       "bid-456",
		Timestamp: 1234567890123,
		Service:   "iot-gateway",
		Action:    "property.report",
		DeviceSN:  "device-001",
	}

	// Test SetData
	testData := map[string]any{
		"temperature": 25.5,
		"humidity":    60,
	}
	err := msg.SetData(testData)
	require.NoError(t, err)

	// Test GetData
	var retrieved map[string]any
	err = msg.GetData(&retrieved)
	require.NoError(t, err)
	assert.Equal(t, 25.5, retrieved["temperature"])
	assert.Equal(t, float64(60), retrieved["humidity"])
}

func TestNewStandardMessageWithIDs(t *testing.T) {
	data := map[string]any{
		"temperature": 25.5,
	}

	msg, err := NewStandardMessageWithIDs("custom-tid", "custom-bid", "iot-gateway", "property.report", "device-001", data)
	require.NoError(t, err)

	assert.Equal(t, "custom-tid", msg.TID)
	assert.Equal(t, "custom-bid", msg.BID)
	assert.Equal(t, "iot-gateway", msg.Service)
	assert.Equal(t, "property.report", msg.Action)
	assert.Equal(t, "device-001", msg.DeviceSN)
	assert.True(t, msg.Timestamp > 0)
}
