package handler

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	dji "github.com/utmos/utmos/pkg/adapter/dji"
)

func TestNewOSDHandler(t *testing.T) {
	handler := NewOSDHandler()
	assert.NotNil(t, handler)
	assert.NotNil(t, handler.parser)
}

func TestOSDHandler_GetTopicType(t *testing.T) {
	handler := NewOSDHandler()
	assert.Equal(t, dji.TopicTypeOSD, handler.GetTopicType())
}

func TestOSDHandler_Handle_NilMessage(t *testing.T) {
	handler := NewOSDHandler()
	topic := &dji.TopicInfo{
		Type:     dji.TopicTypeOSD,
		DeviceSN: "TEST-SN",
	}

	_, err := handler.Handle(context.Background(), nil, topic)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "nil message")
}

func TestOSDHandler_Handle_NilTopic(t *testing.T) {
	handler := NewOSDHandler()
	msg := &dji.Message{
		TID:  "test-tid",
		BID:  "test-bid",
		Data: json.RawMessage(`{"mode_code": 0}`),
	}

	_, err := handler.Handle(context.Background(), msg, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "nil topic")
}

func TestOSDHandler_Handle_AircraftOSD(t *testing.T) {
	handler := NewOSDHandler()

	msg := &dji.Message{
		TID:       "test-tid-001",
		BID:       "test-bid-001",
		Timestamp: 1706000000000,
		Data: json.RawMessage(`{
			"mode_code": 0,
			"longitude": 116.397128,
			"latitude": 39.916527,
			"height": 100.5,
			"elevation": 50.0,
			"horizontal_speed": 5.5,
			"vertical_speed": 1.2,
			"battery": {
				"capacity_percent": 85
			},
			"payloads": [
				{"payload_index": "39-0-7", "sn": "PAYLOAD-001"}
			]
		}`),
	}

	topic := &dji.TopicInfo{
		Type:      dji.TopicTypeOSD,
		DeviceSN:  "AIRCRAFT-SN-001",
		GatewaySN: "DOCK-SN-001",
		Raw:       "thing/product/DOCK-SN-001/osd",
	}

	sm, err := handler.Handle(context.Background(), msg, topic)
	require.NoError(t, err)
	require.NotNil(t, sm)

	// Verify StandardMessage fields
	assert.Equal(t, "test-tid-001", sm.TID)
	assert.Equal(t, "test-bid-001", sm.BID)
	assert.Equal(t, int64(1706000000000), sm.Timestamp)
	assert.Equal(t, "AIRCRAFT-SN-001", sm.DeviceSN)
	assert.Equal(t, dji.VendorDJI, sm.Service)
	assert.Equal(t, dji.ActionPropertyReport, sm.Action)

	// Verify ProtocolMeta
	require.NotNil(t, sm.ProtocolMeta)
	assert.Equal(t, dji.VendorDJI, sm.ProtocolMeta.Vendor)
	assert.Equal(t, "thing/product/DOCK-SN-001/osd", sm.ProtocolMeta.OriginalTopic)

	// Verify data
	var data map[string]interface{}
	err = json.Unmarshal(sm.Data, &data)
	require.NoError(t, err)

	assert.Equal(t, "aircraft", data["osd_type"])
	assert.Equal(t, "AIRCRAFT-SN-001", data["device_sn"])
	assert.InDelta(t, 116.397128, data["longitude"].(float64), 0.0001)
	assert.InDelta(t, 39.916527, data["latitude"].(float64), 0.0001)
	assert.InDelta(t, 100.5, data["height"].(float64), 0.01)
	assert.Equal(t, float64(0), data["mode_code"])
	assert.Equal(t, float64(85), data["battery_percent"])
}

func TestOSDHandler_Handle_DockOSD(t *testing.T) {
	handler := NewOSDHandler()

	msg := &dji.Message{
		TID:       "test-tid-002",
		BID:       "test-bid-002",
		Timestamp: 1706000000000,
		Data: json.RawMessage(`{
			"mode_code": 0,
			"cover_state": 0,
			"putter_state": 0,
			"drone_in_dock": 1,
			"longitude": 116.397128,
			"latitude": 39.916527,
			"height": 50.0,
			"environment_temperature": 25.5
		}`),
	}

	topic := &dji.TopicInfo{
		Type:      dji.TopicTypeOSD,
		DeviceSN:  "DOCK-SN-001",
		GatewaySN: "DOCK-SN-001",
		Raw:       "thing/product/DOCK-SN-001/osd",
	}

	sm, err := handler.Handle(context.Background(), msg, topic)
	require.NoError(t, err)
	require.NotNil(t, sm)

	// Verify data
	var data map[string]interface{}
	err = json.Unmarshal(sm.Data, &data)
	require.NoError(t, err)

	assert.Equal(t, "dock", data["osd_type"])
	assert.Equal(t, "DOCK-SN-001", data["device_sn"])
	assert.Equal(t, float64(0), data["cover_state"])
	assert.Equal(t, float64(1), data["drone_in_dock"])
}

func TestOSDHandler_Handle_RCOSD(t *testing.T) {
	handler := NewOSDHandler()

	msg := &dji.Message{
		TID:       "test-tid-003",
		BID:       "test-bid-003",
		Timestamp: 1706000000000,
		Data: json.RawMessage(`{
			"capacity_percent": 80,
			"longitude": 116.397128,
			"latitude": 39.916527,
			"height": 50.0,
			"wireless_link": {
				"sdr_quality": 4
			}
		}`),
	}

	topic := &dji.TopicInfo{
		Type:      dji.TopicTypeOSD,
		DeviceSN:  "RC-SN-001",
		GatewaySN: "RC-SN-001",
		Raw:       "thing/product/RC-SN-001/osd",
	}

	sm, err := handler.Handle(context.Background(), msg, topic)
	require.NoError(t, err)
	require.NotNil(t, sm)

	// Verify data
	var data map[string]interface{}
	err = json.Unmarshal(sm.Data, &data)
	require.NoError(t, err)

	assert.Equal(t, "rc", data["osd_type"])
	assert.Equal(t, "RC-SN-001", data["device_sn"])
	assert.Equal(t, float64(80), data["capacity_percent"])
}

func TestOSDHandler_Handle_InvalidData(t *testing.T) {
	handler := NewOSDHandler()

	msg := &dji.Message{
		TID:  "test-tid",
		BID:  "test-bid",
		Data: json.RawMessage(`{invalid json}`),
	}

	topic := &dji.TopicInfo{
		Type:     dji.TopicTypeOSD,
		DeviceSN: "TEST-SN",
	}

	_, err := handler.Handle(context.Background(), msg, topic)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse OSD")
}

func TestOSDHandler_Handle_EmptyData(t *testing.T) {
	handler := NewOSDHandler()

	msg := &dji.Message{
		TID:  "test-tid",
		BID:  "test-bid",
		Data: json.RawMessage(``),
	}

	topic := &dji.TopicInfo{
		Type:     dji.TopicTypeOSD,
		DeviceSN: "TEST-SN",
	}

	_, err := handler.Handle(context.Background(), msg, topic)
	require.Error(t, err)
}

func TestOSDHandler_Handle_TimestampDefault(t *testing.T) {
	handler := NewOSDHandler()

	msg := &dji.Message{
		TID:       "test-tid",
		BID:       "test-bid",
		Timestamp: 0, // No timestamp
		Data:      json.RawMessage(`{"mode_code": 0, "payloads": []}`),
	}

	topic := &dji.TopicInfo{
		Type:     dji.TopicTypeOSD,
		DeviceSN: "TEST-SN",
	}

	sm, err := handler.Handle(context.Background(), msg, topic)
	require.NoError(t, err)
	require.NotNil(t, sm)

	// Timestamp should be set to current time
	assert.Greater(t, sm.Timestamp, int64(0))
}

func TestOSDHandler_ImplementsHandler(_ *testing.T) {
	var _ Handler = (*OSDHandler)(nil)
}
