package dji

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/utmos/utmos/pkg/adapter"
	"github.com/utmos/utmos/pkg/rabbitmq"
)

func TestAdapter_GetVendor(t *testing.T) {
	a := NewAdapter()
	assert.Equal(t, VendorDJI, a.GetVendor())
}

func TestAdapter_ParseRawMessage(t *testing.T) {
	a := NewAdapter()

	tests := []struct {
		name     string
		topic    string
		payload  []byte
		wantType adapter.MessageType
		wantErr  bool
	}{
		{
			name:  "parse OSD message",
			topic: "thing/product/ABC123/osd",
			payload: []byte(`{
				"tid": "tid-123",
				"bid": "bid-456",
				"timestamp": 1234567890123,
				"data": {"latitude": 39.9042}
			}`),
			wantType: adapter.MessageTypeProperty,
			wantErr:  false,
		},
		{
			name:  "parse events message",
			topic: "thing/product/GW789/events",
			payload: []byte(`{
				"tid": "tid-789",
				"bid": "bid-012",
				"method": "fly_to_point_progress",
				"data": {"progress": 50}
			}`),
			wantType: adapter.MessageTypeEvent,
			wantErr:  false,
		},
		{
			name:    "invalid topic",
			topic:   "invalid/topic",
			payload: []byte(`{"tid": "123", "bid": "456"}`),
			wantErr: true,
		},
		{
			name:    "invalid payload",
			topic:   "thing/product/ABC123/osd",
			payload: []byte(`{invalid json}`),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pm, err := a.ParseRawMessage(tt.topic, tt.payload)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, VendorDJI, pm.Vendor)
			assert.Equal(t, tt.topic, pm.Topic)
			assert.Equal(t, tt.wantType, pm.MessageType)
		})
	}
}

func TestAdapter_ToStandardMessage(t *testing.T) {
	a := NewAdapter()

	pm := &adapter.ProtocolMessage{
		Vendor:      VendorDJI,
		Topic:       "thing/product/ABC123/osd",
		DeviceSN:    "ABC123",
		MessageType: adapter.MessageTypeProperty,
		TID:         "tid-123",
		BID:         "bid-456",
		Timestamp:   1234567890123,
		Data:        json.RawMessage(`{"latitude": 39.9042}`),
	}

	sm, err := a.ToStandardMessage(pm)
	require.NoError(t, err)

	assert.Equal(t, "tid-123", sm.TID)
	assert.Equal(t, "bid-456", sm.BID)
	assert.Equal(t, "ABC123", sm.DeviceSN)
	assert.Equal(t, "property.report", sm.Action)
	assert.Equal(t, "dji-adapter", sm.Service)
	assert.NotNil(t, sm.ProtocolMeta)
	assert.Equal(t, VendorDJI, sm.ProtocolMeta.Vendor)
}

func TestAdapter_FromStandardMessage(t *testing.T) {
	a := NewAdapter()

	qos := 1
	stdMsg := &rabbitmq.StandardMessage{
		TID:      "tid-123",
		BID:      "bid-456",
		Service:  "dji-adapter",
		Action:   "service.call",
		DeviceSN: "ABC123",
		Data:     json.RawMessage(`{"method": "takeoff"}`),
		ProtocolMeta: &rabbitmq.ProtocolMeta{
			Vendor: VendorDJI,
			Method: "takeoff",
			QoS:    &qos,
		},
	}

	pm, err := a.FromStandardMessage(stdMsg)
	require.NoError(t, err)

	assert.Equal(t, VendorDJI, pm.Vendor)
	assert.Equal(t, "ABC123", pm.DeviceSN)
	assert.Equal(t, adapter.MessageTypeService, pm.MessageType)
}

func TestAdapter_GetRawPayload(t *testing.T) {
	a := NewAdapter()

	pm := &adapter.ProtocolMessage{
		TID:       "tid-123",
		BID:       "bid-456",
		Timestamp: 1234567890123,
		Method:    "takeoff",
		Data:      json.RawMessage(`{"param": "value"}`),
	}

	payload, err := a.GetRawPayload(pm)
	require.NoError(t, err)

	// Verify the payload is valid JSON
	var parsed map[string]any
	err = json.Unmarshal(payload, &parsed)
	require.NoError(t, err)

	assert.Equal(t, "tid-123", parsed["tid"])
	assert.Equal(t, "bid-456", parsed["bid"])
	assert.Equal(t, "takeoff", parsed["method"])
}

func TestAdapter_ImplementsInterface(_ *testing.T) {
	// Compile-time check that Adapter implements ProtocolAdapter
	var _ adapter.ProtocolAdapter = (*Adapter)(nil)
}

func TestRegister(t *testing.T) {
	// Reset global registry for test
	adapter.Unregister(VendorDJI)

	Register()

	retrieved, err := adapter.Get(VendorDJI)
	require.NoError(t, err)
	assert.Equal(t, VendorDJI, retrieved.GetVendor())

	// Cleanup
	adapter.Unregister(VendorDJI)
}
