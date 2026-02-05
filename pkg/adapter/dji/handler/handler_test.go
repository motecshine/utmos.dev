package handler

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	dji "github.com/utmos/utmos/pkg/adapter/dji"
	"github.com/utmos/utmos/pkg/rabbitmq"
)

// mockHandler is a test implementation of Handler.
type mockHandler struct {
	topicType dji.TopicType
	handleFn  func(ctx context.Context, msg *dji.Message, topic *dji.TopicInfo) (*rabbitmq.StandardMessage, error)
}

func (h *mockHandler) Handle(ctx context.Context, msg *dji.Message, topic *dji.TopicInfo) (*rabbitmq.StandardMessage, error) {
	if h.handleFn != nil {
		return h.handleFn(ctx, msg, topic)
	}
	return &rabbitmq.StandardMessage{}, nil
}

func (h *mockHandler) GetTopicType() dji.TopicType {
	return h.topicType
}

func TestHandlerInterface(_ *testing.T) {
	// Verify mockHandler implements Handler interface
	var _ Handler = (*mockHandler)(nil)
}

func TestMockHandler_Handle(t *testing.T) {
	expectedMsg := &rabbitmq.StandardMessage{
		TID:      "test-tid",
		DeviceSN: "test-device",
	}

	handler := &mockHandler{
		topicType: dji.TopicTypeOSD,
		handleFn: func(_ context.Context, _ *dji.Message, _ *dji.TopicInfo) (*rabbitmq.StandardMessage, error) {
			return expectedMsg, nil
		},
	}

	result, err := handler.Handle(context.Background(), &dji.Message{}, &dji.TopicInfo{})
	require.NoError(t, err)
	assert.Equal(t, expectedMsg, result)
}

func TestMockHandler_GetTopicType(t *testing.T) {
	tests := []struct {
		name      string
		topicType dji.TopicType
	}{
		{"OSD", dji.TopicTypeOSD},
		{"State", dji.TopicTypeState},
		{"Events", dji.TopicTypeEvents},
		{"Services", dji.TopicTypeServices},
		{"Status", dji.TopicTypeStatus},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &mockHandler{topicType: tt.topicType}
			assert.Equal(t, tt.topicType, handler.GetTopicType())
		})
	}
}
