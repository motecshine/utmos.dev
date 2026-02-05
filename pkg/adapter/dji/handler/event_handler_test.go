package handler

import (
	"context"
	"encoding/json"
	"testing"

	dji "github.com/utmos/utmos/pkg/adapter/dji"
	"github.com/utmos/utmos/pkg/adapter/dji/router"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewEventHandler(t *testing.T) {
	r := router.NewEventRouter()
	handler := NewEventHandler(r)
	assert.NotNil(t, handler)
	assert.NotNil(t, handler.router)
}

func TestEventHandler_GetTopicType(t *testing.T) {
	r := router.NewEventRouter()
	handler := NewEventHandler(r)
	assert.Equal(t, dji.TopicTypeEvents, handler.GetTopicType())
}

func TestEventHandler_Handle_NilMessage(t *testing.T) {
	r := router.NewEventRouter()
	handler := NewEventHandler(r)
	topic := &dji.TopicInfo{
		Type:     dji.TopicTypeEvents,
		DeviceSN: "TEST-SN",
	}

	_, err := handler.Handle(context.Background(), nil, topic)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "nil message")
}

func TestEventHandler_Handle_NilTopic(t *testing.T) {
	r := router.NewEventRouter()
	handler := NewEventHandler(r)
	msg := &dji.Message{
		TID:    "test-tid",
		BID:    "test-bid",
		Method: "hms",
	}

	_, err := handler.Handle(context.Background(), msg, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "nil topic")
}

func TestEventHandler_Handle_EventRequest(t *testing.T) {
	r := router.NewEventRouter()
	handler := NewEventHandler(r)

	msg := &dji.Message{
		TID:       "evt-tid-001",
		BID:       "evt-bid-001",
		Timestamp: 1706000000000,
		Method:    "hms",
		Data: json.RawMessage(`{
			"list": [
				{
					"code": "0x16100001",
					"level": 0,
					"module": 3,
					"in_the_sky": 0
				}
			]
		}`),
	}

	topic := &dji.TopicInfo{
		Type:      dji.TopicTypeEvents,
		DeviceSN:  "DOCK-SN-001",
		GatewaySN: "DOCK-SN-001",
		Raw:       "thing/product/DOCK-SN-001/events",
	}

	sm, err := handler.Handle(context.Background(), msg, topic)
	require.NoError(t, err)
	require.NotNil(t, sm)

	// Verify StandardMessage fields
	assert.Equal(t, "evt-tid-001", sm.TID)
	assert.Equal(t, "evt-bid-001", sm.BID)
	assert.Equal(t, int64(1706000000000), sm.Timestamp)
	assert.Equal(t, "DOCK-SN-001", sm.DeviceSN)
	assert.Equal(t, dji.VendorDJI, sm.Service)
	assert.Equal(t, dji.ActionEventReport, sm.Action)

	// Verify ProtocolMeta
	require.NotNil(t, sm.ProtocolMeta)
	assert.Equal(t, dji.VendorDJI, sm.ProtocolMeta.Vendor)
	assert.Equal(t, "hms", sm.ProtocolMeta.Method)

	// Verify data
	var data map[string]interface{}
	err = json.Unmarshal(sm.Data, &data)
	require.NoError(t, err)

	assert.Equal(t, "event_request", data["message_type"])
	assert.Equal(t, "hms", data["method"])
	assert.Equal(t, "DOCK-SN-001", data["device_sn"])
}

func TestEventHandler_Handle_EventReply(t *testing.T) {
	r := router.NewEventRouter()
	handler := NewEventHandler(r)

	msg := &dji.Message{
		TID:       "evt-tid-001",
		BID:       "evt-bid-001",
		Timestamp: 1706000000000,
		Method:    "hms",
		Data: json.RawMessage(`{
			"result": 0
		}`),
	}

	topic := &dji.TopicInfo{
		Type:      dji.TopicTypeEventsReply,
		DeviceSN:  "DOCK-SN-001",
		GatewaySN: "DOCK-SN-001",
		Raw:       "thing/product/DOCK-SN-001/events_reply",
	}

	sm, err := handler.Handle(context.Background(), msg, topic)
	require.NoError(t, err)
	require.NotNil(t, sm)

	// Verify action is reply
	assert.Equal(t, dji.ActionEventReply, sm.Action)

	// Verify data
	var data map[string]interface{}
	err = json.Unmarshal(sm.Data, &data)
	require.NoError(t, err)

	assert.Equal(t, "event_reply", data["message_type"])
	assert.Equal(t, "hms", data["method"])
	assert.Equal(t, float64(0), data["result"])
}

func TestEventHandler_Handle_NeedReply(t *testing.T) {
	r := router.NewEventRouter()
	handler := NewEventHandler(r)

	needReply := 1
	msg := &dji.Message{
		TID:       "evt-tid-002",
		BID:       "evt-bid-002",
		Timestamp: 1706000000000,
		Method:    "flighttask_ready",
		NeedReply: &needReply,
		Data:      json.RawMessage(`{"flight_id": "flight-001"}`),
	}

	topic := &dji.TopicInfo{
		Type:      dji.TopicTypeEvents,
		DeviceSN:  "DOCK-SN-001",
		GatewaySN: "DOCK-SN-001",
		Raw:       "thing/product/DOCK-SN-001/events",
	}

	sm, err := handler.Handle(context.Background(), msg, topic)
	require.NoError(t, err)
	require.NotNil(t, sm)

	var data map[string]interface{}
	err = json.Unmarshal(sm.Data, &data)
	require.NoError(t, err)

	assert.Equal(t, true, data["need_reply"])
}

func TestEventHandler_Handle_HMSEvent(t *testing.T) {
	r := router.NewEventRouter()
	handler := NewEventHandler(r)

	msg := &dji.Message{
		TID:       "evt-tid-003",
		BID:       "evt-bid-003",
		Timestamp: 1706000000000,
		Method:    "hms",
		Data: json.RawMessage(`{
			"list": [
				{
					"code": "0x16100001",
					"level": 0,
					"module": 3,
					"in_the_sky": 0,
					"args": {
						"component_index": 0,
						"sensor_index": 0
					}
				},
				{
					"code": "0x16100002",
					"level": 1,
					"module": 3,
					"in_the_sky": 1
				}
			]
		}`),
	}

	topic := &dji.TopicInfo{
		Type:      dji.TopicTypeEvents,
		DeviceSN:  "DOCK-SN-001",
		GatewaySN: "DOCK-SN-001",
		Raw:       "thing/product/DOCK-SN-001/events",
	}

	sm, err := handler.Handle(context.Background(), msg, topic)
	require.NoError(t, err)
	require.NotNil(t, sm)

	var data map[string]interface{}
	err = json.Unmarshal(sm.Data, &data)
	require.NoError(t, err)

	assert.Equal(t, "hms", data["method"])
	assert.NotNil(t, data["data"])
}

func TestEventHandler_Handle_FileUploadCallback(t *testing.T) {
	r := router.NewEventRouter()
	handler := NewEventHandler(r)

	msg := &dji.Message{
		TID:       "evt-tid-004",
		BID:       "evt-bid-004",
		Timestamp: 1706000000000,
		Method:    "file_upload_callback",
		Data: json.RawMessage(`{
			"file": {
				"path": "/media/DJI_001.jpg",
				"name": "DJI_001.jpg",
				"size": 1024000,
				"fingerprint": "abc123"
			}
		}`),
	}

	topic := &dji.TopicInfo{
		Type:      dji.TopicTypeEvents,
		DeviceSN:  "DOCK-SN-001",
		GatewaySN: "DOCK-SN-001",
		Raw:       "thing/product/DOCK-SN-001/events",
	}

	sm, err := handler.Handle(context.Background(), msg, topic)
	require.NoError(t, err)
	require.NotNil(t, sm)

	var data map[string]interface{}
	err = json.Unmarshal(sm.Data, &data)
	require.NoError(t, err)

	assert.Equal(t, "file_upload_callback", data["method"])
}

func TestEventHandler_Handle_TimestampDefault(t *testing.T) {
	r := router.NewEventRouter()
	handler := NewEventHandler(r)

	msg := &dji.Message{
		TID:       "evt-tid",
		BID:       "evt-bid",
		Timestamp: 0,
		Method:    "hms",
	}

	topic := &dji.TopicInfo{
		Type:     dji.TopicTypeEvents,
		DeviceSN: "TEST-SN",
	}

	sm, err := handler.Handle(context.Background(), msg, topic)
	require.NoError(t, err)
	require.NotNil(t, sm)

	assert.Greater(t, sm.Timestamp, int64(0))
}

func TestEventHandler_GetRouter(t *testing.T) {
	r := router.NewEventRouter()
	handler := NewEventHandler(r)

	assert.Equal(t, r, handler.GetRouter())
}

func TestEventHandler_ImplementsHandler(t *testing.T) {
	var _ Handler = (*EventHandler)(nil)
}
