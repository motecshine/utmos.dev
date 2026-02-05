package router

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterWaylineCommands(t *testing.T) {
	r := NewServiceRouter()
	err := RegisterWaylineCommands(r)
	require.NoError(t, err)

	// Verify all wayline commands are registered
	expectedMethods := []string{
		MethodFlighttaskCreate,
		MethodFlighttaskPrepare,
		MethodFlighttaskExecute,
		MethodFlighttaskPause,
		MethodFlighttaskRecovery,
		MethodFlighttaskUndo,
		MethodReturnHome,
		MethodReturnHomeCancel,
	}

	for _, method := range expectedMethods {
		assert.True(t, r.Has(method), "method %s should be registered", method)
	}
}

func TestWaylineCommands_FlighttaskCreate(t *testing.T) {
	r := NewServiceRouter()
	err := RegisterWaylineCommands(r)
	require.NoError(t, err)

	req := &ServiceRequest{
		Method: MethodFlighttaskCreate,
		Data: json.RawMessage(`{
			"flighttask_id": "flight-001",
			"task_type": "immediate",
			"wayline_type": "wayline",
			"out_of_control_action": "execute_go_home"
		}`),
	}

	resp, err := r.RouteService(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Result)
}

func TestWaylineCommands_FlighttaskPrepare(t *testing.T) {
	r := NewServiceRouter()
	err := RegisterWaylineCommands(r)
	require.NoError(t, err)

	req := &ServiceRequest{
		Method: MethodFlighttaskPrepare,
		Data: json.RawMessage(`{
			"flight_id": "flight-001",
			"execute_time": 1706000000,
			"task_type": 0,
			"file": {"url": "https://example.com/wayline.kmz", "fingerprint": "abc123"},
			"out_of_control_action": 0,
			"exit_wayline_when_rc_lost": 0,
			"rth_altitude": 100,
			"rth_mode": 0
		}`),
	}

	resp, err := r.RouteService(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Result)
}

func TestWaylineCommands_FlighttaskExecute(t *testing.T) {
	r := NewServiceRouter()
	err := RegisterWaylineCommands(r)
	require.NoError(t, err)

	req := &ServiceRequest{
		Method: MethodFlighttaskExecute,
		Data: json.RawMessage(`{
			"flight_id": "flight-001"
		}`),
	}

	resp, err := r.RouteService(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Result)
}

func TestWaylineCommands_FlighttaskPause(t *testing.T) {
	r := NewServiceRouter()
	err := RegisterWaylineCommands(r)
	require.NoError(t, err)

	req := &ServiceRequest{
		Method: MethodFlighttaskPause,
	}

	resp, err := r.RouteService(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Result)
}

func TestWaylineCommands_FlighttaskRecovery(t *testing.T) {
	r := NewServiceRouter()
	err := RegisterWaylineCommands(r)
	require.NoError(t, err)

	req := &ServiceRequest{
		Method: MethodFlighttaskRecovery,
	}

	resp, err := r.RouteService(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Result)
}

func TestWaylineCommands_FlighttaskUndo(t *testing.T) {
	r := NewServiceRouter()
	err := RegisterWaylineCommands(r)
	require.NoError(t, err)

	req := &ServiceRequest{
		Method: MethodFlighttaskUndo,
		Data: json.RawMessage(`{
			"flight_ids": ["flight-001", "flight-002"]
		}`),
	}

	resp, err := r.RouteService(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Result)
}

func TestWaylineCommands_ReturnHome(t *testing.T) {
	r := NewServiceRouter()
	err := RegisterWaylineCommands(r)
	require.NoError(t, err)

	req := &ServiceRequest{
		Method: MethodReturnHome,
	}

	resp, err := r.RouteService(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Result)
}

func TestWaylineCommands_ReturnHomeCancel(t *testing.T) {
	r := NewServiceRouter()
	err := RegisterWaylineCommands(r)
	require.NoError(t, err)

	req := &ServiceRequest{
		Method: MethodReturnHomeCancel,
	}

	resp, err := r.RouteService(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Result)
}

func TestWaylineCommands_InvalidData(t *testing.T) {
	r := NewServiceRouter()
	err := RegisterWaylineCommands(r)
	require.NoError(t, err)

	req := &ServiceRequest{
		Method: MethodFlighttaskCreate,
		Data:   json.RawMessage(`{invalid json}`),
	}

	resp, err := r.RouteService(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, 314000, resp.Result) // Parameter error
}
