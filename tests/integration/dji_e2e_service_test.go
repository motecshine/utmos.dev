package integration

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	dji "github.com/utmos/utmos/pkg/adapter/dji"
	"github.com/utmos/utmos/pkg/adapter/dji/config"
	djiinit "github.com/utmos/utmos/pkg/adapter/dji/init"
	"github.com/utmos/utmos/pkg/rabbitmq"
)

// TestServiceFlow_E2E tests the complete service call flow end-to-end.
func TestServiceFlow_E2E(t *testing.T) {
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
			name:  "Service call - flighttask_prepare",
			topic: "thing/product/gateway-001/services_reply",
			payload: `{
				"tid": "tid-svc-001",
				"bid": "bid-svc-001",
				"timestamp": 1234567890123,
				"gateway": "gateway-001",
				"method": "flighttask_prepare",
				"data": {
					"result": 0
				}
			}`,
			wantService: "dji-adapter",
			wantAction:  "service.reply",
			wantMethod:  "flighttask_prepare",
		},
		{
			name:  "Service call - cover_open",
			topic: "thing/product/dock-001/services_reply",
			payload: `{
				"tid": "tid-svc-002",
				"bid": "bid-svc-002",
				"timestamp": 1234567890124,
				"gateway": "dock-001",
				"method": "cover_open",
				"data": {
					"result": 0
				}
			}`,
			wantService: "dji-adapter",
			wantAction:  "service.reply",
			wantMethod:  "cover_open",
		},
		{
			name:  "Service call - camera_photo_take",
			topic: "thing/product/gateway-001/services_reply",
			payload: `{
				"tid": "tid-svc-003",
				"bid": "bid-svc-003",
				"timestamp": 1234567890125,
				"gateway": "gateway-001",
				"method": "camera_photo_take",
				"data": {
					"result": 0,
					"output": {
						"file_path": "/media/photo_001.jpg"
					}
				}
			}`,
			wantService: "dji-adapter",
			wantAction:  "service.reply",
			wantMethod:  "camera_photo_take",
		},
		{
			name:  "Service call - drc_mode_enter",
			topic: "thing/product/gateway-001/services_reply",
			payload: `{
				"tid": "tid-svc-004",
				"bid": "bid-svc-004",
				"timestamp": 1234567890126,
				"gateway": "gateway-001",
				"method": "drc_mode_enter",
				"data": {
					"result": 0
				}
			}`,
			wantService: "dji-adapter",
			wantAction:  "service.reply",
			wantMethod:  "drc_mode_enter",
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

// TestServiceFlow_Downlink tests converting StandardMessage to DJI service call format.
func TestServiceFlow_Downlink(t *testing.T) {
	adapter := djiinit.NewInitializedAdapter()

	tests := []struct {
		name       string
		deviceSN   string
		method     string
		params     map[string]interface{}
		wantTopic  string
		wantMethod string
	}{
		{
			name:     "Downlink - flighttask_execute",
			deviceSN: "gateway-001",
			method:   "flighttask_execute",
			params: map[string]interface{}{
				"flight_id": "flight-001",
			},
			wantTopic:  "thing/product/gateway-001/services",
			wantMethod: "flighttask_execute",
		},
		{
			name:       "Downlink - cover_close",
			deviceSN:   "dock-001",
			method:     "cover_close",
			params:     map[string]interface{}{},
			wantTopic:  "thing/product/dock-001/services",
			wantMethod: "cover_close",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create StandardMessage for downlink
			dataBytes, err := json.Marshal(tt.params)
			require.NoError(t, err)

			stdMsg := &rabbitmq.StandardMessage{
				TID:       "tid-downlink",
				BID:       "bid-downlink",
				Timestamp: 1234567890123,
				DeviceSN:  tt.deviceSN,
				Service:   "device",
				Action:    dji.ActionServiceCall,
				Data:      dataBytes,
				ProtocolMeta: &rabbitmq.ProtocolMeta{
					Vendor: "dji",
					Method: tt.method,
				},
			}

			// Convert to protocol message
			pm, err := adapter.FromStandardMessage(stdMsg)
			require.NoError(t, err)
			require.NotNil(t, pm)

			assert.Equal(t, tt.wantTopic, pm.Topic)
			assert.Equal(t, tt.wantMethod, pm.Method)
			assert.Equal(t, stdMsg.TID, pm.TID)
			assert.Equal(t, stdMsg.BID, pm.BID)
		})
	}
}

// TestServiceFlow_Timeout tests service call timeout configuration.
func TestServiceFlow_Timeout(t *testing.T) {
	// Verify timeout configuration from config package
	assert.Equal(t, 30, int(config.ServiceCallTimeout.Seconds()))
}

// TestServiceFlow_ErrorResponse tests handling of service call errors.
func TestServiceFlow_ErrorResponse(t *testing.T) {
	adapter := djiinit.NewInitializedAdapter()
	ctx := context.Background()

	// Service call with error response
	payload := `{
		"tid": "tid-err",
		"bid": "bid-err",
		"timestamp": 1234567890123,
		"gateway": "gateway-001",
		"method": "flighttask_execute",
		"data": {
			"result": 316001,
			"output": {
				"status": "rejected"
			}
		}
	}`

	result, err := adapter.HandleMessage(ctx, "thing/product/gateway-001/services_reply", []byte(payload))
	require.NoError(t, err)
	require.NotNil(t, result)

	// Error should be preserved in data
	assert.NotNil(t, result.Data)
}
