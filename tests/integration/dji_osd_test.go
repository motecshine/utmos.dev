package integration

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	djiinit "github.com/utmos/utmos/pkg/adapter/dji/init"
)

// TestOSDFlow_Integration tests the complete OSD message processing flow.
// This test verifies that OSD messages from DJI devices are correctly parsed
// and converted to StandardMessage format.
func TestOSDFlow_Integration(t *testing.T) {
	// Initialize adapter with all handlers
	adapter := djiinit.NewInitializedAdapter()
	require.NotNil(t, adapter)

	tests := []struct {
		name        string
		topic       string
		payload     string
		wantService string
		wantAction  string
		wantErr     bool
	}{
		{
			name:  "Aircraft OSD message",
			topic: "thing/product/gateway-001/osd",
			payload: `{
				"tid": "tid-001",
				"bid": "bid-001",
				"timestamp": 1234567890123,
				"gateway": "gateway-001",
				"data": {
					"host": {
						"latitude": 31.2304,
						"longitude": 121.4737,
						"altitude": 100.5,
						"height": 50.0,
						"attitude_pitch": 0.5,
						"attitude_roll": 0.2,
						"attitude_yaw": 180.0,
						"horizontal_speed": 5.0,
						"vertical_speed": 0.0,
						"battery": {
							"capacity_percent": 85,
							"voltage": 48000,
							"temperature": 25.0
						}
					}
				}
			}`,
			wantService: "dji",
			wantAction:  "property.report",
			wantErr:     false,
		},
		{
			name:  "Dock OSD message",
			topic: "thing/product/dock-001/osd",
			payload: `{
				"tid": "tid-002",
				"bid": "bid-002",
				"timestamp": 1234567890124,
				"gateway": "dock-001",
				"data": {
					"host": {
						"network_state": {
							"type": 2,
							"quality": 4,
							"rate": 100.0
						},
						"drone_in_dock": 1,
						"drone_charge_state": {
							"state": 1,
							"capacity_percent": 90
						}
					}
				}
			}`,
			wantService: "dji",
			wantAction:  "property.report",
			wantErr:     false,
		},
		{
			name:  "RC OSD message",
			topic: "thing/product/rc-001/osd",
			payload: `{
				"tid": "tid-003",
				"bid": "bid-003",
				"timestamp": 1234567890125,
				"gateway": "rc-001",
				"data": {
					"host": {
						"latitude": 31.2305,
						"longitude": 121.4738,
						"capacity_percent": 75
					}
				}
			}`,
			wantService: "dji",
			wantAction:  "property.report",
			wantErr:     false,
		},
		{
			name:    "Invalid topic",
			topic:   "invalid/topic",
			payload: `{"tid":"t","bid":"b","timestamp":123,"data":{}}`,
			wantErr: true,
		},
		{
			name:    "Malformed JSON",
			topic:   "thing/product/gateway-001/osd",
			payload: `{invalid json}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			// Handle message using the adapter's HandleMessage method
			result, err := adapter.HandleMessage(ctx, tt.topic, []byte(tt.payload))

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, result)

			// Verify result
			assert.Equal(t, tt.wantService, result.Service)
			assert.Equal(t, tt.wantAction, result.Action)
			assert.NotEmpty(t, result.TID)
			assert.NotEmpty(t, result.BID)
			assert.NotZero(t, result.Timestamp)
			assert.NotEmpty(t, result.DeviceSN)
			assert.Equal(t, "dji", result.ProtocolMeta.Vendor)
		})
	}
}

// TestOSDFlow_NestedDeviceData tests OSD messages with nested device data.
func TestOSDFlow_NestedDeviceData(t *testing.T) {
	adapter := djiinit.NewInitializedAdapter()
	ctx := context.Background()

	// OSD message with both gateway and sub-device data
	payload := `{
		"tid": "tid-nested",
		"bid": "bid-nested",
		"timestamp": 1234567890126,
		"gateway": "dock-001",
		"data": {
			"host": {
				"network_state": {"type": 2, "quality": 4},
				"drone_in_dock": 1
			},
			"0-0-0": {
				"latitude": 31.2304,
				"longitude": 121.4737,
				"altitude": 100.5,
				"battery": {"capacity_percent": 85}
			}
		}
	}`

	result, err := adapter.HandleMessage(ctx, "thing/product/dock-001/osd", []byte(payload))
	require.NoError(t, err)
	require.NotNil(t, result)

	// Verify nested data is preserved
	assert.NotNil(t, result.Data)
}

// TestOSDFlow_PartialData tests OSD messages with partial/optional fields.
func TestOSDFlow_PartialData(t *testing.T) {
	adapter := djiinit.NewInitializedAdapter()
	ctx := context.Background()

	// OSD message with minimal required fields
	payload := `{
		"tid": "tid-partial",
		"bid": "bid-partial",
		"timestamp": 1234567890127,
		"gateway": "gateway-001",
		"data": {
			"host": {
				"latitude": 31.2304,
				"longitude": 121.4737
			}
		}
	}`

	result, err := adapter.HandleMessage(ctx, "thing/product/gateway-001/osd", []byte(payload))
	require.NoError(t, err)
	require.NotNil(t, result)

	// Should handle partial data gracefully
	assert.NotNil(t, result.Data)
}

// TestOSDFlow_RabbitMQMessage simulates receiving OSD via RabbitMQ.
func TestOSDFlow_RabbitMQMessage(t *testing.T) {
	adapter := djiinit.NewInitializedAdapter()
	ctx := context.Background()

	// Simulate RabbitMQ message structure
	type RawUplinkMessage struct {
		Vendor    string          `json:"vendor"`
		Topic     string          `json:"topic"`
		Payload   json.RawMessage `json:"payload"`
		QoS       int             `json:"qos"`
		Timestamp int64           `json:"timestamp"`
		TraceID   string          `json:"trace_id"`
		SpanID    string          `json:"span_id"`
	}

	osdPayload := `{
		"tid": "tid-rmq",
		"bid": "bid-rmq",
		"timestamp": 1234567890128,
		"gateway": "gateway-001",
		"data": {
			"host": {
				"latitude": 31.2304,
				"longitude": 121.4737,
				"altitude": 100.5
			}
		}
	}`

	rmqMsg := RawUplinkMessage{
		Vendor:    "dji",
		Topic:     "thing/product/gateway-001/osd",
		Payload:   json.RawMessage(osdPayload),
		QoS:       1,
		Timestamp: time.Now().UnixMilli(),
		TraceID:   "trace-001",
		SpanID:    "span-001",
	}

	// Process as if received from RabbitMQ
	result, err := adapter.HandleMessage(ctx, rmqMsg.Topic, rmqMsg.Payload)
	require.NoError(t, err)
	require.NotNil(t, result)

	// Verify StandardMessage output
	assert.Equal(t, "dji", result.Service)
	assert.Equal(t, "property.report", result.Action)
	assert.Equal(t, "dji", result.ProtocolMeta.Vendor)
	assert.Equal(t, rmqMsg.Topic, result.ProtocolMeta.OriginalTopic)
}
