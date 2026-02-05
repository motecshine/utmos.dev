package dji

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseTopic(t *testing.T) {
	tests := []struct {
		name     string
		topic    string
		wantType TopicType
		wantSN   string
		wantErr  bool
	}{
		{
			name:     "osd topic",
			topic:    "thing/product/ABC123/osd",
			wantType: TopicTypeOSD,
			wantSN:   "ABC123",
			wantErr:  false,
		},
		{
			name:     "state topic",
			topic:    "thing/product/DEF456/state",
			wantType: TopicTypeState,
			wantSN:   "DEF456",
			wantErr:  false,
		},
		{
			name:     "services topic",
			topic:    "thing/product/GW789/services",
			wantType: TopicTypeServices,
			wantSN:   "GW789",
			wantErr:  false,
		},
		{
			name:     "services reply topic",
			topic:    "thing/product/GW789/services_reply",
			wantType: TopicTypeServicesReply,
			wantSN:   "GW789",
			wantErr:  false,
		},
		{
			name:     "events topic",
			topic:    "thing/product/GW789/events",
			wantType: TopicTypeEvents,
			wantSN:   "GW789",
			wantErr:  false,
		},
		{
			name:     "status topic",
			topic:    "sys/product/GW789/status",
			wantType: TopicTypeStatus,
			wantSN:   "GW789",
			wantErr:  false,
		},
		{
			name:     "status reply topic",
			topic:    "sys/product/GW789/status_reply",
			wantType: TopicTypeStatusReply,
			wantSN:   "GW789",
			wantErr:  false,
		},
		{
			name:    "invalid - empty topic",
			topic:   "",
			wantErr: true,
		},
		{
			name:    "invalid - wrong prefix",
			topic:   "mqtt/product/ABC123/osd",
			wantErr: true,
		},
		{
			name:    "invalid - too few segments",
			topic:   "thing/product/osd",
			wantErr: true,
		},
		{
			name:    "invalid - unknown type",
			topic:   "thing/product/ABC123/unknown",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info, err := ParseTopic(tt.topic)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantType, info.Type)
			assert.Equal(t, tt.wantSN, info.DeviceSN)
		})
	}
}

func TestTopicInfo_IsUplink(t *testing.T) {
	tests := []struct {
		topicType TopicType
		isUplink  bool
	}{
		{TopicTypeOSD, true},
		{TopicTypeState, true},
		{TopicTypeEvents, true},
		{TopicTypeStatus, true},
		{TopicTypeServicesReply, true},
		{TopicTypeServices, false},
		{TopicTypeStatusReply, false},
	}

	for _, tt := range tests {
		t.Run(string(tt.topicType), func(t *testing.T) {
			info := &TopicInfo{Type: tt.topicType}
			assert.Equal(t, tt.isUplink, info.IsUplink())
		})
	}
}

func TestTopicInfo_IsDownlink(t *testing.T) {
	tests := []struct {
		topicType  TopicType
		isDownlink bool
	}{
		{TopicTypeServices, true},
		{TopicTypeStatusReply, true},
		{TopicTypeOSD, false},
		{TopicTypeState, false},
		{TopicTypeEvents, false},
	}

	for _, tt := range tests {
		t.Run(string(tt.topicType), func(t *testing.T) {
			info := &TopicInfo{Type: tt.topicType}
			assert.Equal(t, tt.isDownlink, info.IsDownlink())
		})
	}
}

func TestBuildTopic(t *testing.T) {
	tests := []struct {
		name      string
		topicType TopicType
		deviceSN  string
		expected  string
	}{
		{
			name:      "build osd topic",
			topicType: TopicTypeOSD,
			deviceSN:  "ABC123",
			expected:  "thing/product/ABC123/osd",
		},
		{
			name:      "build services topic",
			topicType: TopicTypeServices,
			deviceSN:  "GW789",
			expected:  "thing/product/GW789/services",
		},
		{
			name:      "build status topic",
			topicType: TopicTypeStatus,
			deviceSN:  "GW789",
			expected:  "sys/product/GW789/status",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BuildTopic(tt.topicType, tt.deviceSN)
			assert.Equal(t, tt.expected, result)
		})
	}
}
