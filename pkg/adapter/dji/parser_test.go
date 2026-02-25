package dji

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseMessage(t *testing.T) {
	tests := []struct {
		name    string
		payload []byte
		wantTID string
		wantBID string
		wantErr bool
	}{
		{
			name: "valid OSD message",
			payload: []byte(`{
				"tid": "tid-123",
				"bid": "bid-456",
				"timestamp": 1234567890123,
				"data": {"latitude": 39.9042, "longitude": 116.4074}
			}`),
			wantTID: "tid-123",
			wantBID: "bid-456",
			wantErr: false,
		},
		{
			name: "valid service message with method",
			payload: []byte(`{
				"tid": "tid-789",
				"bid": "bid-012",
				"timestamp": 1234567890123,
				"method": "thing.property.post",
				"data": {"param": "value"}
			}`),
			wantTID: "tid-789",
			wantBID: "bid-012",
			wantErr: false,
		},
		{
			name: "message with need_reply",
			payload: []byte(`{
				"tid": "tid-abc",
				"bid": "bid-def",
				"timestamp": 1234567890123,
				"need_reply": 1,
				"data": {}
			}`),
			wantTID: "tid-abc",
			wantBID: "bid-def",
			wantErr: false,
		},
		{
			name:    "invalid JSON",
			payload: []byte(`{invalid json}`),
			wantErr: true,
		},
		{
			name:    "empty payload",
			payload: []byte{},
			wantErr: true,
		},
		{
			name:    "nil payload",
			payload: nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg, err := ParseMessage(tt.payload)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantTID, msg.TID)
			assert.Equal(t, tt.wantBID, msg.BID)
		})
	}
}

func TestMessage_Validate(t *testing.T) {
	tests := []struct {
		name    string
		msg     *Message
		wantErr bool
	}{
		{
			name: "valid message",
			msg: &Message{
				TID:       "tid-123",
				BID:       "bid-456",
				Timestamp: 1234567890123,
				Data:      json.RawMessage(`{}`),
			},
			wantErr: false,
		},
		{
			name: "missing TID",
			msg: &Message{
				BID:       "bid-456",
				Timestamp: 1234567890123,
				Data:      json.RawMessage(`{}`),
			},
			wantErr: true,
		},
		{
			name: "missing BID",
			msg: &Message{
				TID:       "tid-123",
				Timestamp: 1234567890123,
				Data:      json.RawMessage(`{}`),
			},
			wantErr: true,
		},
		{
			name: "zero timestamp is allowed",
			msg: &Message{
				TID:  "tid-123",
				BID:  "bid-456",
				Data: json.RawMessage(`{}`),
			},
			wantErr: false,
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

func TestMessage_NeedReplyBool(t *testing.T) {
	tests := []struct {
		name      string
		needReply *int
		expected  bool
	}{
		{
			name:      "nil need_reply",
			needReply: nil,
			expected:  false,
		},
		{
			name:      "need_reply = 0",
			needReply: intPtr(0),
			expected:  false,
		},
		{
			name:      "need_reply = 1",
			needReply: intPtr(1),
			expected:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := &Message{NeedReply: tt.needReply}
			assert.Equal(t, tt.expected, msg.NeedReplyBool())
		})
	}
}

func intPtr(i int) *int {
	return &i
}
