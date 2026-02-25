// Package handler provides message handlers for the DJI adapter.
package handler

import (
	"context"

	dji "github.com/utmos/utmos/pkg/adapter/dji"
	"github.com/utmos/utmos/pkg/rabbitmq"
)

// Handler defines the interface for DJI message handlers.
// Each handler is responsible for processing a specific topic type.
type Handler interface {
	// Handle processes a DJI message and returns a StandardMessage.
	Handle(ctx context.Context, msg *dji.Message, topic *dji.TopicInfo) (*rabbitmq.StandardMessage, error)

	// GetTopicType returns the topic type this handler processes.
	GetTopicType() dji.TopicType
}

// Func is a function type that implements Handler.
type Func func(ctx context.Context, msg *dji.Message, topic *dji.TopicInfo) (*rabbitmq.StandardMessage, error)
