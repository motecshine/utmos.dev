package dji

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/utmos/utmos/pkg/rabbitmq"
)

func TestToStandardMessage(t *testing.T) {
	tests := []struct {
		name       string
		djiMsg     *Message
		topicInfo  *TopicInfo
		wantAction string
		wantErr    bool
	}{
		{
			name: "OSD message to property.report",
			djiMsg: &Message{
				TID:       "tid-123",
				BID:       "bid-456",
				Timestamp: 1234567890123,
				Data:      json.RawMessage(`{"latitude": 39.9042}`),
			},
			topicInfo: &TopicInfo{
				Type:     TopicTypeOSD,
				DeviceSN: "ABC123",
			},
			wantAction: "property.report",
			wantErr:    false,
		},
		{
			name: "State message to property.report",
			djiMsg: &Message{
				TID:       "tid-123",
				BID:       "bid-456",
				Timestamp: 1234567890123,
				Data:      json.RawMessage(`{"battery": 80}`),
			},
			topicInfo: &TopicInfo{
				Type:     TopicTypeState,
				DeviceSN: "ABC123",
			},
			wantAction: "property.report",
			wantErr:    false,
		},
		{
			name: "Events message to event.report",
			djiMsg: &Message{
				TID:       "tid-123",
				BID:       "bid-456",
				Timestamp: 1234567890123,
				Method:    "fly_to_point_progress",
				Data:      json.RawMessage(`{"progress": 50}`),
			},
			topicInfo: &TopicInfo{
				Type:     TopicTypeEvents,
				DeviceSN: "GW789",
			},
			wantAction: "event.report",
			wantErr:    false,
		},
		{
			name: "Status message to device.online",
			djiMsg: &Message{
				TID:       "tid-123",
				BID:       "bid-456",
				Timestamp: 1234567890123,
				Data:      json.RawMessage(`{"status": "online"}`),
			},
			topicInfo: &TopicInfo{
				Type:     TopicTypeStatus,
				DeviceSN: "GW789",
			},
			wantAction: "device.online",
			wantErr:    false,
		},
		{
			name: "Services reply to service.reply",
			djiMsg: &Message{
				TID:       "tid-123",
				BID:       "bid-456",
				Timestamp: 1234567890123,
				Method:    "takeoff",
				Data:      json.RawMessage(`{"result": 0}`),
			},
			topicInfo: &TopicInfo{
				Type:     TopicTypeServicesReply,
				DeviceSN: "GW789",
			},
			wantAction: "service.reply",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			converter := NewConverter()
			stdMsg, err := converter.ToStandardMessage(tt.djiMsg, tt.topicInfo)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantAction, stdMsg.Action)
			assert.Equal(t, tt.djiMsg.TID, stdMsg.TID)
			assert.Equal(t, tt.djiMsg.BID, stdMsg.BID)
			assert.Equal(t, tt.topicInfo.DeviceSN, stdMsg.DeviceSN)
			assert.NotNil(t, stdMsg.ProtocolMeta)
			assert.Equal(t, VendorDJI, stdMsg.ProtocolMeta.Vendor)
		})
	}
}

func TestFromStandardMessage(t *testing.T) {
	tests := []struct {
		name    string
		stdMsg  *rabbitmq.StandardMessage
		wantErr bool
	}{
		{
			name: "service call message",
			stdMsg: &rabbitmq.StandardMessage{
				TID:       "tid-123",
				BID:       "bid-456",
				Timestamp: 1234567890123,
				Service:   "dji-adapter",
				Action:    "service.call",
				DeviceSN:  "ABC123",
				Data:      json.RawMessage(`{"method": "takeoff", "params": {}}`),
				ProtocolMeta: &rabbitmq.ProtocolMeta{
					Vendor:        "dji",
					OriginalTopic: "thing/product/ABC123/services",
				},
			},
			wantErr: false,
		},
		{
			name: "property set message",
			stdMsg: &rabbitmq.StandardMessage{
				TID:       "tid-789",
				BID:       "bid-012",
				Timestamp: 1234567890123,
				Service:   "dji-adapter",
				Action:    "property.set",
				DeviceSN:  "ABC123",
				Data:      json.RawMessage(`{"gimbal_pitch": -30}`),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			converter := NewConverter()
			djiMsg, err := converter.FromStandardMessage(tt.stdMsg)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.stdMsg.TID, djiMsg.TID)
			assert.Equal(t, tt.stdMsg.BID, djiMsg.BID)
		})
	}
}

func TestMapTopicTypeToAction(t *testing.T) {
	tests := []struct {
		topicType TopicType
		expected  string
	}{
		{TopicTypeOSD, "property.report"},
		{TopicTypeState, "property.report"},
		{TopicTypeEvents, "event.report"},
		{TopicTypeStatus, "device.online"},
		{TopicTypeServicesReply, "service.reply"},
	}

	for _, tt := range tests {
		t.Run(string(tt.topicType), func(t *testing.T) {
			result := MapTopicTypeToAction(tt.topicType)
			assert.Equal(t, tt.expected, result)
		})
	}
}
