package router

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterWaylineEvents(t *testing.T) {
	r := NewEventRouter()
	err := RegisterWaylineEvents(r)
	require.NoError(t, err)

	// Verify all wayline events are registered
	expectedMethods := []string{
		MethodFlighttaskProgress,
		MethodFlighttaskReady,
		MethodReturnHomeInfo,
	}

	for _, method := range expectedMethods {
		assert.True(t, r.Has(method), "method %s should be registered", method)
	}
}

func TestWaylineEvents_FlighttaskProgress(t *testing.T) {
	r := NewEventRouter()
	err := RegisterWaylineEvents(r)
	require.NoError(t, err)

	tests := []struct {
		name string
		data string
	}{
		{
			name: "executing",
			data: `{"flight_id": "flight-001", "status": "executing", "progress": 50}`,
		},
		{
			name: "completed",
			data: `{"flight_id": "flight-001", "status": "completed", "progress": 100, "result": 0}`,
		},
		{
			name: "failed",
			data: `{"flight_id": "flight-001", "status": "failed", "progress": 30, "result": 314001}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &EventRequest{
				Method: MethodFlighttaskProgress,
				Data:   json.RawMessage(tt.data),
			}
			resp, err := r.RouteEvent(context.Background(), req)
			require.NoError(t, err)
			assert.Equal(t, 0, resp.Result)
		})
	}
}

func TestWaylineEvents_FlighttaskReady(t *testing.T) {
	r := NewEventRouter()
	err := RegisterWaylineEvents(r)
	require.NoError(t, err)

	needReply := 1
	req := &EventRequest{
		Method:    MethodFlighttaskReady,
		NeedReply: &needReply,
		Data:      json.RawMessage(`{"flight_id": "flight-001"}`),
	}

	resp, err := r.RouteEvent(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Result)
}

func TestWaylineEvents_ReturnHomeInfo(t *testing.T) {
	r := NewEventRouter()
	err := RegisterWaylineEvents(r)
	require.NoError(t, err)

	req := &EventRequest{
		Method: MethodReturnHomeInfo,
		Data: json.RawMessage(`{
			"flight_id": "flight-001",
			"last_point_type": 1,
			"planned_path_points": [
				{"latitude": 22.5, "longitude": 113.9, "height": 100}
			]
		}`),
	}

	resp, err := r.RouteEvent(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Result)
}
