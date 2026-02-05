package handler

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	dji "github.com/utmos/utmos/pkg/adapter/dji"
	"github.com/utmos/utmos/pkg/rabbitmq"
)

func TestNewRegistry(t *testing.T) {
	r := NewRegistry()
	assert.NotNil(t, r)
	assert.NotNil(t, r.handlers)
	assert.Empty(t, r.handlers)
}

func TestRegistry_Register(t *testing.T) {
	r := NewRegistry()

	handler := &mockHandler{topicType: dji.TopicTypeOSD}

	// First registration should succeed
	err := r.Register(handler)
	require.NoError(t, err)

	// Duplicate registration should fail
	err = r.Register(handler)
	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrHandlerAlreadyRegistered))
}

func TestRegistry_Get(t *testing.T) {
	r := NewRegistry()

	osdHandler := &mockHandler{topicType: dji.TopicTypeOSD}
	err := r.Register(osdHandler)
	require.NoError(t, err)

	// Get registered handler
	handler, err := r.Get(dji.TopicTypeOSD)
	require.NoError(t, err)
	assert.Equal(t, osdHandler, handler)

	// Get unregistered handler
	_, err = r.Get(dji.TopicTypeState)
	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrHandlerNotFound))
}

func TestRegistry_Has(t *testing.T) {
	r := NewRegistry()

	handler := &mockHandler{topicType: dji.TopicTypeOSD}
	err := r.Register(handler)
	require.NoError(t, err)

	assert.True(t, r.Has(dji.TopicTypeOSD))
	assert.False(t, r.Has(dji.TopicTypeState))
}

func TestRegistry_List(t *testing.T) {
	r := NewRegistry()

	handlers := []*mockHandler{
		{topicType: dji.TopicTypeOSD},
		{topicType: dji.TopicTypeState},
		{topicType: dji.TopicTypeEvents},
	}

	for _, h := range handlers {
		err := r.Register(h)
		require.NoError(t, err)
	}

	list := r.List()
	assert.Len(t, list, 3)
	assert.Contains(t, list, dji.TopicTypeOSD)
	assert.Contains(t, list, dji.TopicTypeState)
	assert.Contains(t, list, dji.TopicTypeEvents)
}

func TestRegistry_MustRegister(t *testing.T) {
	r := NewRegistry()

	handler := &mockHandler{topicType: dji.TopicTypeOSD}

	// Should not panic
	assert.NotPanics(t, func() {
		r.MustRegister(handler)
	})

	// Should panic on duplicate
	assert.Panics(t, func() {
		r.MustRegister(handler)
	})
}

func TestRegistry_ConcurrentAccess(t *testing.T) {
	r := NewRegistry()

	// Register handlers
	topicTypes := []dji.TopicType{
		dji.TopicTypeOSD,
		dji.TopicTypeState,
		dji.TopicTypeEvents,
		dji.TopicTypeServices,
		dji.TopicTypeStatus,
	}

	for _, tt := range topicTypes {
		err := r.Register(&mockHandler{topicType: tt})
		require.NoError(t, err)
	}

	// Concurrent reads
	done := make(chan bool)
	for i := 0; i < 100; i++ {
		go func() {
			_ = r.List()
			_ = r.Has(dji.TopicTypeOSD)
			_, _ = r.Get(dji.TopicTypeOSD)
			done <- true
		}()
	}

	for i := 0; i < 100; i++ {
		<-done
	}
}

func TestRegistry_Integration(t *testing.T) {
	r := NewRegistry()

	// Create a handler that returns a specific message
	expectedMsg := &rabbitmq.StandardMessage{
		TID:      "test-tid",
		DeviceSN: "test-device",
		Action:   "property.report",
	}

	osdHandler := &mockHandler{
		topicType: dji.TopicTypeOSD,
		handleFn: func(ctx context.Context, msg *dji.Message, topic *dji.TopicInfo) (*rabbitmq.StandardMessage, error) {
			return expectedMsg, nil
		},
	}

	err := r.Register(osdHandler)
	require.NoError(t, err)

	// Get handler and process message
	handler, err := r.Get(dji.TopicTypeOSD)
	require.NoError(t, err)

	result, err := handler.Handle(context.Background(), &dji.Message{}, &dji.TopicInfo{})
	require.NoError(t, err)
	assert.Equal(t, expectedMsg, result)
}
