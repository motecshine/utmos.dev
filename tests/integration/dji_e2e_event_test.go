package integration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	djiinit "github.com/utmos/utmos/pkg/adapter/dji/init"
)

// TestEventFlow_E2E tests the complete event processing flow end-to-end.
func TestEventFlow_E2E(t *testing.T) {
	adapter := djiinit.NewInitializedAdapter()
	ctx := context.Background()

	tests := []struct {
		name        string
		topic       string
		payload     string
		wantService string
		wantAction  string
		wantMethod  string
	}{
		{
			name:  "Event - flighttask_progress",
			topic: "thing/product/gateway-001/events",
			payload: `{
				"tid": "tid-evt-001",
				"bid": "bid-evt-001",
				"timestamp": 1234567890123,
				"gateway": "gateway-001",
				"method": "flighttask_progress",
				"need_reply": 0,
				"data": {
					"status": "executing",
					"progress": {
						"current_step": 5,
						"total_step": 10,
						"percent": 50
					}
				}
			}`,
			wantService: "dji",
			wantAction:  "event.report",
			wantMethod:  "flighttask_progress",
		},
		{
			name:  "Event - flighttask_ready (need_reply)",
			topic: "thing/product/gateway-001/events",
			payload: `{
				"tid": "tid-evt-002",
				"bid": "bid-evt-002",
				"timestamp": 1234567890124,
				"gateway": "gateway-001",
				"method": "flighttask_ready",
				"need_reply": 1,
				"data": {
					"flight_ids": ["flight-001", "flight-002"]
				}
			}`,
			wantService: "dji",
			wantAction:  "event.report",
			wantMethod:  "flighttask_ready",
		},
		{
			name:  "Event - hms (Health Management System)",
			topic: "thing/product/gateway-001/events",
			payload: `{
				"tid": "tid-evt-003",
				"bid": "bid-evt-003",
				"timestamp": 1234567890125,
				"gateway": "gateway-001",
				"method": "hms",
				"need_reply": 0,
				"data": {
					"list": [
						{
							"code": "0x16100001",
							"level": 1,
							"module": 1,
							"in_the_sky": 0
						}
					]
				}
			}`,
			wantService: "dji",
			wantAction:  "event.report",
			wantMethod:  "hms",
		},
		{
			name:  "Event - file_upload_callback",
			topic: "thing/product/gateway-001/events",
			payload: `{
				"tid": "tid-evt-004",
				"bid": "bid-evt-004",
				"timestamp": 1234567890126,
				"gateway": "gateway-001",
				"method": "file_upload_callback",
				"need_reply": 0,
				"data": {
					"file": {
						"path": "/media/photo_001.jpg",
						"name": "photo_001.jpg",
						"size": 1024000
					}
				}
			}`,
			wantService: "dji",
			wantAction:  "event.report",
			wantMethod:  "file_upload_callback",
		},
		{
			name:  "Event - device_exit_homing_notify",
			topic: "thing/product/dock-001/events",
			payload: `{
				"tid": "tid-evt-005",
				"bid": "bid-evt-005",
				"timestamp": 1234567890127,
				"gateway": "dock-001",
				"method": "device_exit_homing_notify",
				"need_reply": 0,
				"data": {
					"action": 0,
					"reason": 1
				}
			}`,
			wantService: "dji",
			wantAction:  "event.report",
			wantMethod:  "device_exit_homing_notify",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := adapter.HandleMessage(ctx, tt.topic, []byte(tt.payload))
			require.NoError(t, err)
			require.NotNil(t, result)

			assert.Equal(t, tt.wantService, result.Service)
			assert.Equal(t, tt.wantAction, result.Action)
			assert.Equal(t, tt.wantMethod, result.ProtocolMeta.Method)
			assert.Equal(t, "dji", result.ProtocolMeta.Vendor)
		})
	}
}

// TestEventFlow_Reply tests event reply generation.
func TestEventFlow_Reply(t *testing.T) {
	adapter := djiinit.NewInitializedAdapter()
	ctx := context.Background()

	// Event that requires reply
	payload := `{
		"tid": "tid-reply",
		"bid": "bid-reply",
		"timestamp": 1234567890123,
		"gateway": "gateway-001",
		"method": "flighttask_ready",
		"need_reply": 1,
		"data": {
			"flight_ids": ["flight-001"]
		}
	}`

	result, err := adapter.HandleMessage(ctx, "thing/product/gateway-001/events", []byte(payload))
	require.NoError(t, err)
	require.NotNil(t, result)

	// Verify event was processed
	assert.Equal(t, "dji", result.Service)
	assert.Equal(t, "event.report", result.Action)
}

// TestEventFlow_StatusChange tests device status change events.
func TestEventFlow_StatusChange(t *testing.T) {
	adapter := djiinit.NewInitializedAdapter()
	ctx := context.Background()

	tests := []struct {
		name       string
		topic      string
		payload    string
		wantAction string
	}{
		{
			name:  "Device online",
			topic: "sys/product/gateway-001/status",
			payload: `{
				"tid": "tid-status-001",
				"bid": "bid-status-001",
				"timestamp": 1234567890123,
				"gateway": "gateway-001",
				"data": {
					"online": true,
					"sub_devices": [
						{
							"device_sn": "aircraft-001",
							"product_type": "aircraft",
							"online": true
						}
					]
				}
			}`,
			wantAction: "device.online",
		},
		{
			name:  "Device offline",
			topic: "sys/product/gateway-001/status",
			payload: `{
				"tid": "tid-status-002",
				"bid": "bid-status-002",
				"timestamp": 1234567890124,
				"gateway": "gateway-001",
				"data": {
					"online": false
				}
			}`,
			wantAction: "device.offline",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := adapter.HandleMessage(ctx, tt.topic, []byte(tt.payload))
			require.NoError(t, err)
			require.NotNil(t, result)

			assert.Equal(t, "dji", result.Service)
			assert.Equal(t, tt.wantAction, result.Action)
		})
	}
}

// TestEventFlow_StateChange tests device state change events.
func TestEventFlow_StateChange(t *testing.T) {
	adapter := djiinit.NewInitializedAdapter()
	ctx := context.Background()

	payload := `{
		"tid": "tid-state",
		"bid": "bid-state",
		"timestamp": 1234567890123,
		"gateway": "gateway-001",
		"data": {
			"flight_mode": 6,
			"gear": 1,
			"battery_percent": 80
		}
	}`

	result, err := adapter.HandleMessage(ctx, "thing/product/gateway-001/state", []byte(payload))
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, "dji", result.Service)
	assert.Equal(t, "property.report", result.Action)
}
