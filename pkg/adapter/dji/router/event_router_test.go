package router

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewEventRouter(t *testing.T) {
	r := NewEventRouter()
	assert.NotNil(t, r)
}

func TestEventRouter_RegisterEventHandler(t *testing.T) {
	r := NewEventRouter()

	handler := func(_ context.Context, _ json.RawMessage) (*EventResponse, error) {
		return &EventResponse{Result: 0}, nil
	}

	err := r.RegisterEventHandler("test_event", handler)
	assert.NoError(t, err)
	assert.True(t, r.Has("test_event"))
}

func TestEventRouter_RegisterEventHandler_Duplicate(t *testing.T) {
	r := NewEventRouter()

	handler := func(_ context.Context, _ json.RawMessage) (*EventResponse, error) {
		return &EventResponse{Result: 0}, nil
	}

	err := r.RegisterEventHandler("test_event", handler)
	require.NoError(t, err)

	err = r.RegisterEventHandler("test_event", handler)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already registered")
}

func TestEventRouter_RouteEvent(t *testing.T) {
	r := NewEventRouter()

	handler := func(_ context.Context, _ json.RawMessage) (*EventResponse, error) {
		return &EventResponse{
			Result: 0,
			Output: json.RawMessage(`{"status": "processed"}`),
		}, nil
	}

	err := r.RegisterEventHandler("test_event", handler)
	require.NoError(t, err)

	req := &EventRequest{
		Method: "test_event",
		Data:   json.RawMessage(`{"key": "value"}`),
	}

	resp, err := r.RouteEvent(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Result)
}

func TestEventRouter_RouteEvent_UnknownMethod(t *testing.T) {
	r := NewEventRouter()

	req := &EventRequest{
		Method: "unknown_event",
	}

	_, err := r.RouteEvent(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown event method")
}

func TestEventRouter_RouteEvent_HandlerError(t *testing.T) {
	r := NewEventRouter()

	handler := func(_ context.Context, _ json.RawMessage) (*EventResponse, error) {
		return nil, errors.New("handler error")
	}

	err := r.RegisterEventHandler("error_event", handler)
	require.NoError(t, err)

	req := &EventRequest{
		Method: "error_event",
	}

	_, err = r.RouteEvent(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "handler error")
}

func TestEventRouter_RouteEvent_NeedReply(t *testing.T) {
	r := NewEventRouter()

	handler := func(_ context.Context, _ json.RawMessage) (*EventResponse, error) {
		return &EventResponse{
			Result: 0,
			Output: json.RawMessage(`{"ack": true}`),
		}, nil
	}

	err := r.RegisterEventHandler("need_reply_event", handler)
	require.NoError(t, err)

	needReply := 1
	req := &EventRequest{
		Method:    "need_reply_event",
		NeedReply: &needReply,
		Data:      json.RawMessage(`{}`),
	}

	resp, err := r.RouteEvent(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Result)
}

func TestEventRouter_List(t *testing.T) {
	r := NewEventRouter()

	handler := func(_ context.Context, _ json.RawMessage) (*EventResponse, error) {
		return &EventResponse{Result: 0}, nil
	}

	err := r.RegisterEventHandler("event_a", handler)
	require.NoError(t, err)
	err = r.RegisterEventHandler("event_b", handler)
	require.NoError(t, err)
	err = r.RegisterEventHandler("event_c", handler)
	require.NoError(t, err)

	list := r.List()
	assert.Len(t, list, 3)
	assert.Contains(t, list, "event_a")
	assert.Contains(t, list, "event_b")
	assert.Contains(t, list, "event_c")
}

func TestEventRouter_Has(t *testing.T) {
	r := NewEventRouter()

	handler := func(_ context.Context, _ json.RawMessage) (*EventResponse, error) {
		return &EventResponse{Result: 0}, nil
	}

	err := r.RegisterEventHandler("existing_event", handler)
	require.NoError(t, err)

	assert.True(t, r.Has("existing_event"))
	assert.False(t, r.Has("non_existing_event"))
}

func TestEventRequest_NeedReplyBool(t *testing.T) {
	tests := []struct {
		name      string
		needReply *int
		expected  bool
	}{
		{
			name:      "nil",
			needReply: nil,
			expected:  false,
		},
		{
			name:      "zero",
			needReply: intPtr(0),
			expected:  false,
		},
		{
			name:      "one",
			needReply: intPtr(1),
			expected:  true,
		},
		{
			name:      "other positive",
			needReply: intPtr(2),
			expected:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &EventRequest{NeedReply: tt.needReply}
			assert.Equal(t, tt.expected, req.NeedReplyBool())
		})
	}
}

func intPtr(i int) *int {
	return &i
}
