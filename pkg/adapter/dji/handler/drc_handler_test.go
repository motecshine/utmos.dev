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

func TestNewDRCHandler(t *testing.T) {
	sr := router.NewServiceRouter()
	er := router.NewEventRouter()
	handler := NewDRCHandler(sr, er)
	assert.NotNil(t, handler)
	assert.NotNil(t, handler.serviceRouter)
	assert.NotNil(t, handler.eventRouter)
	assert.Equal(t, config.DRCHeartbeatTimeout, handler.heartbeatTimeout)
}

func TestDRCHandler_GetTopicType(t *testing.T) {
	sr := router.NewServiceRouter()
	er := router.NewEventRouter()
	handler := NewDRCHandler(sr, er)
	assert.Equal(t, dji.TopicTypeDRCUp, handler.GetTopicType())
}

func TestDRCHandler_Handle_NilMessage(t *testing.T) {
	sr := router.NewServiceRouter()
	er := router.NewEventRouter()
	handler := NewDRCHandler(sr, er)
	topic := &dji.TopicInfo{
		Type:     dji.TopicTypeDRCUp,
		DeviceSN: "TEST-SN",
	}

	_, err := handler.Handle(context.Background(), nil, topic)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "nil message")
}

func TestDRCHandler_Handle_NilTopic(t *testing.T) {
	sr := router.NewServiceRouter()
	er := router.NewEventRouter()
	handler := NewDRCHandler(sr, er)
	msg := &dji.Message{
		TID:    "test-tid",
		BID:    "test-bid",
		Method: "drone_control",
	}

	_, err := handler.Handle(context.Background(), msg, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "nil topic")
}

func TestDRCHandler_Handle_DRCUp(t *testing.T) {
	sr := router.NewServiceRouter()
	er := router.NewEventRouter()
	handler := NewDRCHandler(sr, er)

	msg := &dji.Message{
		TID:       "drc-tid-001",
		BID:       "drc-bid-001",
		Timestamp: 1706000000000,
		Method:    "drone_control",
		Data:      json.RawMessage(`{"x": 0.5, "y": 0.5, "h": 0.0, "w": 0.0}`),
	}

	topic := &dji.TopicInfo{
		Type:      dji.TopicTypeDRCUp,
		DeviceSN:  "DOCK-SN-001",
		GatewaySN: "DOCK-SN-001",
		Raw:       "thing/product/DOCK-SN-001/drc/up",
	}

	sm, err := handler.Handle(context.Background(), msg, topic)
	require.NoError(t, err)
	require.NotNil(t, sm)

	// Verify StandardMessage fields
	assert.Equal(t, "drc-tid-001", sm.TID)
	assert.Equal(t, "drc-bid-001", sm.BID)
	assert.Equal(t, int64(1706000000000), sm.Timestamp)
	assert.Equal(t, "DOCK-SN-001", sm.DeviceSN)
	assert.Equal(t, dji.VendorDJI, sm.Service)
	assert.Equal(t, dji.ActionDRCCommand, sm.Action)

	// Verify ProtocolMeta
	require.NotNil(t, sm.ProtocolMeta)
	assert.Equal(t, dji.VendorDJI, sm.ProtocolMeta.Vendor)
	assert.Equal(t, "drone_control", sm.ProtocolMeta.Method)

	// Verify data
	var data map[string]interface{}
	err = json.Unmarshal(sm.Data, &data)
	require.NoError(t, err)

	assert.Equal(t, "drc_command", data["message_type"])
	assert.Equal(t, "drone_control", data["method"])
}

func TestDRCHandler_Handle_DRCDown(t *testing.T) {
	sr := router.NewServiceRouter()
	er := router.NewEventRouter()
	handler := NewDRCHandler(sr, er)

	msg := &dji.Message{
		TID:       "drc-tid-001",
		BID:       "drc-bid-001",
		Timestamp: 1706000000000,
		Method:    "joystick_invalid_notify",
		Data:      json.RawMessage(`{"reason": "timeout"}`),
	}

	topic := &dji.TopicInfo{
		Type:      dji.TopicTypeDRCDown,
		DeviceSN:  "DOCK-SN-001",
		GatewaySN: "DOCK-SN-001",
		Raw:       "thing/product/DOCK-SN-001/drc/down",
	}

	sm, err := handler.Handle(context.Background(), msg, topic)
	require.NoError(t, err)
	require.NotNil(t, sm)

	// Verify action is event
	assert.Equal(t, dji.ActionDRCEvent, sm.Action)

	// Verify data
	var data map[string]interface{}
	err = json.Unmarshal(sm.Data, &data)
	require.NoError(t, err)

	assert.Equal(t, "drc_event", data["message_type"])
}

func TestDRCHandler_Handle_Heartbeat(t *testing.T) {
	sr := router.NewServiceRouter()
	er := router.NewEventRouter()
	handler := NewDRCHandler(sr, er)

	msg := &dji.Message{
		TID:       "drc-tid-001",
		BID:       "drc-bid-001",
		Timestamp: 1706000000000,
		Method:    "heart",
		Data:      json.RawMessage(`{"seq": 1}`),
	}

	topic := &dji.TopicInfo{
		Type:      dji.TopicTypeDRCUp,
		DeviceSN:  "DOCK-SN-001",
		GatewaySN: "DOCK-SN-001",
		Raw:       "thing/product/DOCK-SN-001/drc/up",
	}

	sm, err := handler.Handle(context.Background(), msg, topic)
	require.NoError(t, err)
	require.NotNil(t, sm)

	// Verify data contains heartbeat info
	var data map[string]interface{}
	err = json.Unmarshal(sm.Data, &data)
	require.NoError(t, err)

	assert.Equal(t, "heart", data["method"])
}

func TestDRCHandler_SetHeartbeatTimeout(t *testing.T) {
	sr := router.NewServiceRouter()
	er := router.NewEventRouter()
	handler := NewDRCHandler(sr, er)

	newTimeout := 5 * time.Second
	handler.SetHeartbeatTimeout(newTimeout)
	assert.Equal(t, newTimeout, handler.heartbeatTimeout)
}

func TestDRCHandler_Handle_TimestampDefault(t *testing.T) {
	sr := router.NewServiceRouter()
	er := router.NewEventRouter()
	handler := NewDRCHandler(sr, er)

	msg := &dji.Message{
		TID:       "drc-tid",
		BID:       "drc-bid",
		Timestamp: 0,
		Method:    "drone_control",
	}

	topic := &dji.TopicInfo{
		Type:     dji.TopicTypeDRCUp,
		DeviceSN: "TEST-SN",
	}

	sm, err := handler.Handle(context.Background(), msg, topic)
	require.NoError(t, err)
	require.NotNil(t, sm)

	assert.Greater(t, sm.Timestamp, int64(0))
}

func TestDRCHandler_ImplementsHandler(t *testing.T) {
	var _ Handler = (*DRCHandler)(nil)
}
