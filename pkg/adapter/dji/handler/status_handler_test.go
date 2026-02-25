package handler

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	dji "github.com/utmos/utmos/pkg/adapter/dji"
)

func TestNewStatusHandler(t *testing.T) {
	handler := NewStatusHandler()
	assert.NotNil(t, handler)
}

func TestStatusHandler_GetTopicType(t *testing.T) {
	handler := NewStatusHandler()
	assert.Equal(t, dji.TopicTypeStatus, handler.GetTopicType())
}

func TestStatusHandler_Handle_NilMessage(t *testing.T) {
	handler := NewStatusHandler()
	topic := &dji.TopicInfo{
		Type:     dji.TopicTypeStatus,
		DeviceSN: "TEST-SN",
	}

	_, err := handler.Handle(context.Background(), nil, topic)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "nil message")
}

func TestStatusHandler_Handle_NilTopic(t *testing.T) {
	handler := NewStatusHandler()
	msg := &dji.Message{
		TID:  "test-tid",
		BID:  "test-bid",
		Data: json.RawMessage(`{"online": true}`),
	}

	_, err := handler.Handle(context.Background(), msg, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "nil topic")
}

func TestStatusHandler_Handle_OnlineStatus(t *testing.T) {
	handler := NewStatusHandler()

	msg := &dji.Message{
		TID:       "test-tid-001",
		BID:       "test-bid-001",
		Timestamp: 1706000000000,
		Data: json.RawMessage(`{
			"online": true,
			"gateway_sn": "DOCK-SN-001",
			"gateway_type": "dock",
			"sub_devices": [
				{
					"device_sn": "AIRCRAFT-SN-001",
					"product_type": "0-67-0",
					"online": true
				}
			]
		}`),
	}

	topic := &dji.TopicInfo{
		Type:      dji.TopicTypeStatus,
		DeviceSN:  "DOCK-SN-001",
		GatewaySN: "DOCK-SN-001",
		Raw:       "sys/product/DOCK-SN-001/status",
	}

	sm, err := handler.Handle(context.Background(), msg, topic)
	require.NoError(t, err)
	require.NotNil(t, sm)

	// Verify StandardMessage fields
	assert.Equal(t, "test-tid-001", sm.TID)
	assert.Equal(t, "test-bid-001", sm.BID)
	assert.Equal(t, int64(1706000000000), sm.Timestamp)
	assert.Equal(t, "DOCK-SN-001", sm.DeviceSN)
	assert.Equal(t, dji.VendorDJI, sm.Service)
	assert.Equal(t, dji.ActionDeviceOnline, sm.Action)

	// Verify ProtocolMeta
	require.NotNil(t, sm.ProtocolMeta)
	assert.Equal(t, dji.VendorDJI, sm.ProtocolMeta.Vendor)
	assert.Equal(t, "sys/product/DOCK-SN-001/status", sm.ProtocolMeta.OriginalTopic)

	// Verify data
	var data map[string]any
	err = json.Unmarshal(sm.Data, &data)
	require.NoError(t, err)

	assert.Equal(t, "status", data["message_type"])
	assert.Equal(t, "DOCK-SN-001", data["device_sn"])
	assert.Equal(t, true, data["online"])

	// Verify topology
	topology := data["topology"].(map[string]any)
	assert.Equal(t, "DOCK-SN-001", topology["gateway_sn"])
	assert.Equal(t, "dock", topology["gateway_type"])

	subDevices := topology["sub_devices"].([]any)
	assert.Len(t, subDevices, 1)
	subDevice := subDevices[0].(map[string]any)
	assert.Equal(t, "AIRCRAFT-SN-001", subDevice["device_sn"])
	assert.Equal(t, true, subDevice["online"])
}

func TestStatusHandler_Handle_OfflineStatus(t *testing.T) {
	handler := NewStatusHandler()

	msg := &dji.Message{
		TID:       "test-tid-002",
		BID:       "test-bid-002",
		Timestamp: 1706000000000,
		Data: json.RawMessage(`{
			"online": false,
			"gateway_sn": "DOCK-SN-001"
		}`),
	}

	topic := &dji.TopicInfo{
		Type:      dji.TopicTypeStatus,
		DeviceSN:  "DOCK-SN-001",
		GatewaySN: "DOCK-SN-001",
		Raw:       "sys/product/DOCK-SN-001/status",
	}

	sm, err := handler.Handle(context.Background(), msg, topic)
	require.NoError(t, err)
	require.NotNil(t, sm)

	assert.Equal(t, dji.ActionDeviceOffline, sm.Action)

	var data map[string]any
	err = json.Unmarshal(sm.Data, &data)
	require.NoError(t, err)

	assert.Equal(t, false, data["online"])
}

func TestStatusHandler_Handle_OnlineAsNumber(t *testing.T) {
	handler := NewStatusHandler()

	msg := &dji.Message{
		TID:       "test-tid",
		BID:       "test-bid",
		Timestamp: 1706000000000,
		Data:      json.RawMessage(`{"online": 1}`),
	}

	topic := &dji.TopicInfo{
		Type:     dji.TopicTypeStatus,
		DeviceSN: "TEST-SN",
	}

	sm, err := handler.Handle(context.Background(), msg, topic)
	require.NoError(t, err)
	require.NotNil(t, sm)

	assert.Equal(t, dji.ActionDeviceOnline, sm.Action)
}

func TestStatusHandler_Handle_OfflineAsNumber(t *testing.T) {
	handler := NewStatusHandler()

	msg := &dji.Message{
		TID:       "test-tid",
		BID:       "test-bid",
		Timestamp: 1706000000000,
		Data:      json.RawMessage(`{"online": 0}`),
	}

	topic := &dji.TopicInfo{
		Type:     dji.TopicTypeStatus,
		DeviceSN: "TEST-SN",
	}

	sm, err := handler.Handle(context.Background(), msg, topic)
	require.NoError(t, err)
	require.NotNil(t, sm)

	assert.Equal(t, dji.ActionDeviceOffline, sm.Action)
}

func TestStatusHandler_Handle_DeviceTopology(t *testing.T) {
	handler := NewStatusHandler()

	msg := &dji.Message{
		TID:       "test-tid",
		BID:       "test-bid",
		Timestamp: 1706000000000,
		Data: json.RawMessage(`{
			"online": true,
			"gateway_sn": "DOCK-SN-001",
			"gateway_type": "dock2",
			"sub_devices": [
				{
					"device_sn": "AIRCRAFT-SN-001",
					"product_type": "0-67-0",
					"online": true
				},
				{
					"device_sn": "PAYLOAD-SN-001",
					"product_type": "1-39-0",
					"online": true
				}
			]
		}`),
	}

	topic := &dji.TopicInfo{
		Type:      dji.TopicTypeStatus,
		DeviceSN:  "DOCK-SN-001",
		GatewaySN: "DOCK-SN-001",
		Raw:       "sys/product/DOCK-SN-001/status",
	}

	sm, err := handler.Handle(context.Background(), msg, topic)
	require.NoError(t, err)
	require.NotNil(t, sm)

	var data map[string]any
	err = json.Unmarshal(sm.Data, &data)
	require.NoError(t, err)

	topology := data["topology"].(map[string]any)
	subDevices := topology["sub_devices"].([]any)
	assert.Len(t, subDevices, 2)
}

func TestStatusHandler_Handle_EmptyData(t *testing.T) {
	handler := NewStatusHandler()

	msg := &dji.Message{
		TID:       "test-tid",
		BID:       "test-bid",
		Timestamp: 1706000000000,
		Data:      json.RawMessage(``),
	}

	topic := &dji.TopicInfo{
		Type:     dji.TopicTypeStatus,
		DeviceSN: "TEST-SN",
	}

	_, err := handler.Handle(context.Background(), msg, topic)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "empty status data")
}

func TestStatusHandler_Handle_TimestampDefault(t *testing.T) {
	handler := NewStatusHandler()

	msg := &dji.Message{
		TID:       "test-tid",
		BID:       "test-bid",
		Timestamp: 0,
		Data:      json.RawMessage(`{"online": true}`),
	}

	topic := &dji.TopicInfo{
		Type:     dji.TopicTypeStatus,
		DeviceSN: "TEST-SN",
	}

	sm, err := handler.Handle(context.Background(), msg, topic)
	require.NoError(t, err)
	require.NotNil(t, sm)

	assert.Greater(t, sm.Timestamp, int64(0))
}

func TestStatusHandler_ImplementsHandler(t *testing.T) {
	var _ Handler = (*StatusHandler)(nil)
}
