package handler

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	dji "github.com/utmos/utmos/pkg/adapter/dji"
	"github.com/utmos/utmos/pkg/adapter/dji/router"
)

func TestNewRequestHandler(t *testing.T) {
	r := router.NewServiceRouter()
	handler := NewRequestHandler(r)
	assert.NotNil(t, handler)
	assert.NotNil(t, handler.router)
}

func TestRequestHandler_GetTopicType(t *testing.T) {
	r := router.NewServiceRouter()
	handler := NewRequestHandler(r)
	assert.Equal(t, dji.TopicTypeRequests, handler.GetTopicType())
}

func TestRequestHandler_Handle_NilMessage(t *testing.T) {
	r := router.NewServiceRouter()
	handler := NewRequestHandler(r)
	topic := &dji.TopicInfo{
		Type:     dji.TopicTypeRequests,
		DeviceSN: "TEST-SN",
	}

	_, err := handler.Handle(context.Background(), nil, topic)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "nil message")
}

func TestRequestHandler_Handle_NilTopic(t *testing.T) {
	r := router.NewServiceRouter()
	handler := NewRequestHandler(r)
	msg := &dji.Message{
		TID:    "test-tid",
		BID:    "test-bid",
		Method: "config_get",
	}

	_, err := handler.Handle(context.Background(), msg, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "nil topic")
}

func TestRequestHandler_Handle_Request(t *testing.T) {
	r := router.NewServiceRouter()
	handler := NewRequestHandler(r)

	msg := &dji.Message{
		TID:       "req-tid-001",
		BID:       "req-bid-001",
		Timestamp: 1706000000000,
		Method:    "config_get",
		Data:      json.RawMessage(`{"config_type": "basic_device_info"}`),
	}

	topic := &dji.TopicInfo{
		Type:      dji.TopicTypeRequests,
		DeviceSN:  "DOCK-SN-001",
		GatewaySN: "DOCK-SN-001",
		Raw:       "thing/product/DOCK-SN-001/requests",
	}

	sm, err := handler.Handle(context.Background(), msg, topic)
	require.NoError(t, err)
	require.NotNil(t, sm)

	// Verify StandardMessage fields
	assert.Equal(t, "req-tid-001", sm.TID)
	assert.Equal(t, "req-bid-001", sm.BID)
	assert.Equal(t, int64(1706000000000), sm.Timestamp)
	assert.Equal(t, "DOCK-SN-001", sm.DeviceSN)
	assert.Equal(t, dji.VendorDJI, sm.Service)
	assert.Equal(t, dji.ActionDeviceRequest, sm.Action)

	// Verify ProtocolMeta
	require.NotNil(t, sm.ProtocolMeta)
	assert.Equal(t, dji.VendorDJI, sm.ProtocolMeta.Vendor)
	assert.Equal(t, "config_get", sm.ProtocolMeta.Method)

	// Verify data
	var data map[string]interface{}
	err = json.Unmarshal(sm.Data, &data)
	require.NoError(t, err)

	assert.Equal(t, "device_request", data["message_type"])
	assert.Equal(t, "config_get", data["method"])
	assert.Equal(t, "DOCK-SN-001", data["device_sn"])
}

func TestRequestHandler_Handle_Reply(t *testing.T) {
	r := router.NewServiceRouter()
	handler := NewRequestHandler(r)

	msg := &dji.Message{
		TID:       "req-tid-001",
		BID:       "req-bid-001",
		Timestamp: 1706000000000,
		Method:    "config_get",
		Data: json.RawMessage(`{
			"result": 0,
			"output": {"name": "Dock-001"}
		}`),
	}

	topic := &dji.TopicInfo{
		Type:      dji.TopicTypeRequestsReply,
		DeviceSN:  "DOCK-SN-001",
		GatewaySN: "DOCK-SN-001",
		Raw:       "thing/product/DOCK-SN-001/requests_reply",
	}

	sm, err := handler.Handle(context.Background(), msg, topic)
	require.NoError(t, err)
	require.NotNil(t, sm)

	// Verify action is reply
	assert.Equal(t, dji.ActionDeviceRequestReply, sm.Action)

	// Verify data
	var data map[string]interface{}
	err = json.Unmarshal(sm.Data, &data)
	require.NoError(t, err)

	assert.Equal(t, "device_request_reply", data["message_type"])
	assert.Equal(t, "config_get", data["method"])
	assert.Equal(t, float64(0), data["result"])
}

func TestRequestHandler_Handle_TimestampDefault(t *testing.T) {
	r := router.NewServiceRouter()
	handler := NewRequestHandler(r)

	msg := &dji.Message{
		TID:       "req-tid",
		BID:       "req-bid",
		Timestamp: 0,
		Method:    "config_get",
	}

	topic := &dji.TopicInfo{
		Type:     dji.TopicTypeRequests,
		DeviceSN: "TEST-SN",
	}

	sm, err := handler.Handle(context.Background(), msg, topic)
	require.NoError(t, err)
	require.NotNil(t, sm)

	assert.Greater(t, sm.Timestamp, int64(0))
}

func TestRequestHandler_GetRouter(t *testing.T) {
	r := router.NewServiceRouter()
	handler := NewRequestHandler(r)

	assert.Equal(t, r, handler.GetRouter())
}

func TestRequestHandler_ImplementsHandler(_ *testing.T) {
	var _ Handler = (*RequestHandler)(nil)
}
