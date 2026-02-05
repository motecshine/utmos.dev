package handler

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	dji "github.com/utmos/utmos/pkg/adapter/dji"
	"github.com/utmos/utmos/pkg/adapter/dji/config"
	"github.com/utmos/utmos/pkg/adapter/dji/router"
)

func TestNewServiceHandler(t *testing.T) {
	r := router.NewServiceRouter()
	handler := NewServiceHandler(r)
	assert.NotNil(t, handler)
	assert.NotNil(t, handler.router)
	assert.Equal(t, config.ServiceCallTimeout, handler.timeout)
}

func TestServiceHandler_GetTopicType(t *testing.T) {
	r := router.NewServiceRouter()
	handler := NewServiceHandler(r)
	assert.Equal(t, dji.TopicTypeServices, handler.GetTopicType())
}

func TestServiceHandler_Handle_NilMessage(t *testing.T) {
	r := router.NewServiceRouter()
	handler := NewServiceHandler(r)
	topic := &dji.TopicInfo{
		Type:     dji.TopicTypeServices,
		DeviceSN: "TEST-SN",
	}

	_, err := handler.Handle(context.Background(), nil, topic)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "nil message")
}

func TestServiceHandler_Handle_NilTopic(t *testing.T) {
	r := router.NewServiceRouter()
	handler := NewServiceHandler(r)
	msg := &dji.Message{
		TID:    "test-tid",
		BID:    "test-bid",
		Method: "cover_open",
	}

	_, err := handler.Handle(context.Background(), msg, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "nil topic")
}

func TestServiceHandler_Handle_ServiceRequest(t *testing.T) {
	r := router.NewServiceRouter()
	handler := NewServiceHandler(r)

	msg := &dji.Message{
		TID:       "test-tid-001",
		BID:       "test-bid-001",
		Timestamp: 1706000000000,
		Method:    "cover_open",
		Data:      json.RawMessage(`{"action": "open"}`),
	}

	topic := &dji.TopicInfo{
		Type:      dji.TopicTypeServices,
		DeviceSN:  "DOCK-SN-001",
		GatewaySN: "DOCK-SN-001",
		Raw:       "thing/product/DOCK-SN-001/services",
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
	assert.Equal(t, dji.ActionServiceCall, sm.Action)

	// Verify ProtocolMeta
	require.NotNil(t, sm.ProtocolMeta)
	assert.Equal(t, dji.VendorDJI, sm.ProtocolMeta.Vendor)
	assert.Equal(t, "cover_open", sm.ProtocolMeta.Method)

	// Verify data
	var data map[string]interface{}
	err = json.Unmarshal(sm.Data, &data)
	require.NoError(t, err)

	assert.Equal(t, "service_request", data["message_type"])
	assert.Equal(t, "cover_open", data["method"])
	assert.Equal(t, "DOCK-SN-001", data["device_sn"])
	assert.Equal(t, float64(30000), data["timeout_ms"]) // 30s default
}

func TestServiceHandler_Handle_ServiceReply(t *testing.T) {
	r := router.NewServiceRouter()
	handler := NewServiceHandler(r)

	msg := &dji.Message{
		TID:       "test-tid-001",
		BID:       "test-bid-001",
		Timestamp: 1706000000000,
		Method:    "cover_open",
		Data: json.RawMessage(`{
			"result": 0,
			"output": {"status": "success"}
		}`),
	}

	topic := &dji.TopicInfo{
		Type:      dji.TopicTypeServicesReply,
		DeviceSN:  "DOCK-SN-001",
		GatewaySN: "DOCK-SN-001",
		Raw:       "thing/product/DOCK-SN-001/services_reply",
	}

	sm, err := handler.Handle(context.Background(), msg, topic)
	require.NoError(t, err)
	require.NotNil(t, sm)

	// Verify action is reply
	assert.Equal(t, dji.ActionServiceReply, sm.Action)

	// Verify data
	var data map[string]interface{}
	err = json.Unmarshal(sm.Data, &data)
	require.NoError(t, err)

	assert.Equal(t, "service_reply", data["message_type"])
	assert.Equal(t, "cover_open", data["method"])
	assert.Equal(t, float64(0), data["result"])
}

func TestServiceHandler_Handle_ServiceReplyWithError(t *testing.T) {
	r := router.NewServiceRouter()
	handler := NewServiceHandler(r)

	msg := &dji.Message{
		TID:       "test-tid",
		BID:       "test-bid",
		Timestamp: 1706000000000,
		Method:    "cover_open",
		Data: json.RawMessage(`{
			"result": 314001,
			"output": {"message": "device offline"}
		}`),
	}

	topic := &dji.TopicInfo{
		Type:      dji.TopicTypeServicesReply,
		DeviceSN:  "DOCK-SN-001",
		GatewaySN: "DOCK-SN-001",
		Raw:       "thing/product/DOCK-SN-001/services_reply",
	}

	sm, err := handler.Handle(context.Background(), msg, topic)
	require.NoError(t, err)
	require.NotNil(t, sm)

	var data map[string]interface{}
	err = json.Unmarshal(sm.Data, &data)
	require.NoError(t, err)

	assert.Equal(t, float64(314001), data["result"])
}

func TestServiceHandler_Handle_NeedReply(t *testing.T) {
	r := router.NewServiceRouter()
	handler := NewServiceHandler(r)

	needReply := 1
	msg := &dji.Message{
		TID:       "test-tid",
		BID:       "test-bid",
		Timestamp: 1706000000000,
		Method:    "flighttask_prepare",
		NeedReply: &needReply,
		Data:      json.RawMessage(`{}`),
	}

	topic := &dji.TopicInfo{
		Type:      dji.TopicTypeServices,
		DeviceSN:  "DOCK-SN-001",
		GatewaySN: "DOCK-SN-001",
		Raw:       "thing/product/DOCK-SN-001/services",
	}

	sm, err := handler.Handle(context.Background(), msg, topic)
	require.NoError(t, err)
	require.NotNil(t, sm)

	var data map[string]interface{}
	err = json.Unmarshal(sm.Data, &data)
	require.NoError(t, err)

	assert.Equal(t, true, data["need_reply"])
}

func TestServiceHandler_SetTimeout(t *testing.T) {
	r := router.NewServiceRouter()
	handler := NewServiceHandler(r)

	newTimeout := 60 * time.Second
	handler.SetTimeout(newTimeout)
	assert.Equal(t, newTimeout, handler.timeout)
}

func TestServiceHandler_GetRouter(t *testing.T) {
	r := router.NewServiceRouter()
	handler := NewServiceHandler(r)

	assert.Equal(t, r, handler.GetRouter())
}

func TestServiceHandler_Handle_TimestampDefault(t *testing.T) {
	r := router.NewServiceRouter()
	handler := NewServiceHandler(r)

	msg := &dji.Message{
		TID:       "test-tid",
		BID:       "test-bid",
		Timestamp: 0,
		Method:    "cover_open",
	}

	topic := &dji.TopicInfo{
		Type:     dji.TopicTypeServices,
		DeviceSN: "TEST-SN",
	}

	sm, err := handler.Handle(context.Background(), msg, topic)
	require.NoError(t, err)
	require.NotNil(t, sm)

	assert.Greater(t, sm.Timestamp, int64(0))
}

func TestServiceHandler_ImplementsHandler(_ *testing.T) {
	var _ Handler = (*ServiceHandler)(nil)
}
