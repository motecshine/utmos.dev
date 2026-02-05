package handler

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	dji "github.com/utmos/utmos/pkg/adapter/dji"
)

func TestNewStateHandler(t *testing.T) {
	handler := NewStateHandler()
	assert.NotNil(t, handler)
}

func TestStateHandler_GetTopicType(t *testing.T) {
	handler := NewStateHandler()
	assert.Equal(t, dji.TopicTypeState, handler.GetTopicType())
}

func TestStateHandler_Handle_NilMessage(t *testing.T) {
	handler := NewStateHandler()
	topic := &dji.TopicInfo{
		Type:     dji.TopicTypeState,
		DeviceSN: "TEST-SN",
	}

	_, err := handler.Handle(context.Background(), nil, topic)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "nil message")
}

func TestStateHandler_Handle_NilTopic(t *testing.T) {
	handler := NewStateHandler()
	msg := &dji.Message{
		TID:  "test-tid",
		BID:  "test-bid",
		Data: json.RawMessage(`{"mode_code": 1}`),
	}

	_, err := handler.Handle(context.Background(), msg, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "nil topic")
}

func TestStateHandler_Handle_PropertyChange(t *testing.T) {
	handler := NewStateHandler()

	msg := &dji.Message{
		TID:       "test-tid-001",
		BID:       "test-bid-001",
		Timestamp: 1706000000000,
		Data: json.RawMessage(`{
			"mode_code": 1,
			"firmware_version": "01.00.0001"
		}`),
	}

	topic := &dji.TopicInfo{
		Type:      dji.TopicTypeState,
		DeviceSN:  "DEVICE-SN-001",
		GatewaySN: "GATEWAY-SN-001",
		Raw:       "thing/product/GATEWAY-SN-001/state",
	}

	sm, err := handler.Handle(context.Background(), msg, topic)
	require.NoError(t, err)
	require.NotNil(t, sm)

	// Verify StandardMessage fields
	assert.Equal(t, "test-tid-001", sm.TID)
	assert.Equal(t, "test-bid-001", sm.BID)
	assert.Equal(t, int64(1706000000000), sm.Timestamp)
	assert.Equal(t, "DEVICE-SN-001", sm.DeviceSN)
	assert.Equal(t, dji.VendorDJI, sm.Service)
	assert.Equal(t, dji.ActionPropertyReport, sm.Action)

	// Verify ProtocolMeta
	require.NotNil(t, sm.ProtocolMeta)
	assert.Equal(t, dji.VendorDJI, sm.ProtocolMeta.Vendor)
	assert.Equal(t, "thing/product/GATEWAY-SN-001/state", sm.ProtocolMeta.OriginalTopic)

	// Verify data
	var data map[string]interface{}
	err = json.Unmarshal(sm.Data, &data)
	require.NoError(t, err)

	assert.Equal(t, "state", data["message_type"])
	assert.Equal(t, "DEVICE-SN-001", data["device_sn"])

	// Verify properties
	props := data["properties"].(map[string]interface{})
	assert.Equal(t, float64(1), props["mode_code"])
	assert.Equal(t, "01.00.0001", props["firmware_version"])

	// Verify changed properties list
	changedProps := data["changed_properties"].([]interface{})
	assert.Len(t, changedProps, 2)
}

func TestStateHandler_Handle_MultipleProperties(t *testing.T) {
	handler := NewStateHandler()

	msg := &dji.Message{
		TID:       "test-tid",
		BID:       "test-bid",
		Timestamp: 1706000000000,
		Data: json.RawMessage(`{
			"mode_code": 2,
			"gear": 1,
			"height_limit": 120,
			"distance_limit_status": {
				"state": 1,
				"distance_limit": 5000
			}
		}`),
	}

	topic := &dji.TopicInfo{
		Type:      dji.TopicTypeState,
		DeviceSN:  "DEVICE-SN-001",
		GatewaySN: "GATEWAY-SN-001",
		Raw:       "thing/product/GATEWAY-SN-001/state",
	}

	sm, err := handler.Handle(context.Background(), msg, topic)
	require.NoError(t, err)
	require.NotNil(t, sm)

	var data map[string]interface{}
	err = json.Unmarshal(sm.Data, &data)
	require.NoError(t, err)

	changedProps := data["changed_properties"].([]interface{})
	assert.Len(t, changedProps, 4)
}

func TestStateHandler_Handle_EmptyData(t *testing.T) {
	handler := NewStateHandler()

	msg := &dji.Message{
		TID:       "test-tid",
		BID:       "test-bid",
		Timestamp: 1706000000000,
		Data:      json.RawMessage(`{}`),
	}

	topic := &dji.TopicInfo{
		Type:      dji.TopicTypeState,
		DeviceSN:  "DEVICE-SN-001",
		GatewaySN: "GATEWAY-SN-001",
		Raw:       "thing/product/GATEWAY-SN-001/state",
	}

	sm, err := handler.Handle(context.Background(), msg, topic)
	require.NoError(t, err)
	require.NotNil(t, sm)

	var data map[string]interface{}
	err = json.Unmarshal(sm.Data, &data)
	require.NoError(t, err)

	changedProps := data["changed_properties"].([]interface{})
	assert.Len(t, changedProps, 0)
}

func TestStateHandler_Handle_TimestampDefault(t *testing.T) {
	handler := NewStateHandler()

	msg := &dji.Message{
		TID:       "test-tid",
		BID:       "test-bid",
		Timestamp: 0,
		Data:      json.RawMessage(`{"mode_code": 0}`),
	}

	topic := &dji.TopicInfo{
		Type:     dji.TopicTypeState,
		DeviceSN: "TEST-SN",
	}

	sm, err := handler.Handle(context.Background(), msg, topic)
	require.NoError(t, err)
	require.NotNil(t, sm)

	assert.Greater(t, sm.Timestamp, int64(0))
}

func TestStateHandler_ImplementsHandler(t *testing.T) {
	var _ Handler = (*StateHandler)(nil)
}
