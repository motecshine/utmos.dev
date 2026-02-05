package init

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	dji "github.com/utmos/utmos/pkg/adapter/dji"
)

func TestNewInitializedAdapter(t *testing.T) {
	adapter := NewInitializedAdapter()
	require.NotNil(t, adapter)
	assert.Equal(t, "dji", adapter.GetVendor())
}

func TestInitializeAdapter(t *testing.T) {
	adapter := dji.NewAdapter()
	err := InitializeAdapter(adapter)
	require.NoError(t, err)

	// Test that HandleMessage works with OSD topic
	ctx := context.Background()
	topic := "thing/product/DOCK123/osd"
	payload := json.RawMessage(`{
		"tid": "test-tid",
		"bid": "test-bid",
		"timestamp": 1234567890,
		"data": {
			"mode_code": 1,
			"cover_state": 0
		}
	}`)

	sm, err := adapter.HandleMessage(ctx, topic, payload)
	require.NoError(t, err)
	require.NotNil(t, sm)
	assert.Equal(t, "test-tid", sm.TID)
	assert.Equal(t, "test-bid", sm.BID)
}

func TestInitializeAdapter_StateHandler(t *testing.T) {
	adapter := NewInitializedAdapter()

	ctx := context.Background()
	topic := "thing/product/DEVICE123/state"
	payload := json.RawMessage(`{
		"tid": "state-tid",
		"bid": "state-bid",
		"timestamp": 1234567890,
		"data": {
			"firmware_version": "1.0.0"
		}
	}`)

	sm, err := adapter.HandleMessage(ctx, topic, payload)
	require.NoError(t, err)
	require.NotNil(t, sm)
	assert.Equal(t, "state-tid", sm.TID)
}

func TestInitializeAdapter_StatusHandler(t *testing.T) {
	adapter := NewInitializedAdapter()

	ctx := context.Background()
	topic := "sys/product/DEVICE123/status"
	payload := json.RawMessage(`{
		"tid": "status-tid",
		"bid": "status-bid",
		"timestamp": 1234567890,
		"data": {
			"online": true
		}
	}`)

	sm, err := adapter.HandleMessage(ctx, topic, payload)
	require.NoError(t, err)
	require.NotNil(t, sm)
	assert.Equal(t, "status-tid", sm.TID)
}

func TestInitializeAdapter_ServiceHandler(t *testing.T) {
	adapter := NewInitializedAdapter()

	ctx := context.Background()
	topic := "thing/product/DEVICE123/services_reply"
	payload := json.RawMessage(`{
		"tid": "service-tid",
		"bid": "service-bid",
		"timestamp": 1234567890,
		"method": "cover_open",
		"data": {
			"result": 0
		}
	}`)

	sm, err := adapter.HandleMessage(ctx, topic, payload)
	require.NoError(t, err)
	require.NotNil(t, sm)
	assert.Equal(t, "service-tid", sm.TID)
}

func TestInitializeAdapter_EventHandler(t *testing.T) {
	adapter := NewInitializedAdapter()

	ctx := context.Background()
	topic := "thing/product/DEVICE123/events"
	payload := json.RawMessage(`{
		"tid": "event-tid",
		"bid": "event-bid",
		"timestamp": 1234567890,
		"method": "hms",
		"data": {
			"list": []
		}
	}`)

	sm, err := adapter.HandleMessage(ctx, topic, payload)
	require.NoError(t, err)
	require.NotNil(t, sm)
	assert.Equal(t, "event-tid", sm.TID)
}

func TestMustInitializeAdapter(t *testing.T) {
	adapter := dji.NewAdapter()

	// Should not panic
	assert.NotPanics(t, func() {
		MustInitializeAdapter(adapter)
	})
}
