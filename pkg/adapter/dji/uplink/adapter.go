// Package uplink provides DJI uplink message processing functionality
package uplink

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/utmos/utmos/pkg/adapter"
	"github.com/utmos/utmos/pkg/rabbitmq"
)

// ProcessorAdapter adapts the DJI Processor to the public adapter.UplinkProcessor interface.
// This adapter does NOT depend on internal/ packages.
type ProcessorAdapter struct {
	processor *Processor
}

// NewProcessorAdapter creates a new processor adapter
func NewProcessorAdapter(logger *logrus.Entry) *ProcessorAdapter {
	return &ProcessorAdapter{
		processor: NewProcessor(logger),
	}
}

// GetVendor returns the vendor name
func (a *ProcessorAdapter) GetVendor() string {
	return a.processor.GetVendor()
}

// CanProcess checks if this processor can handle the given message
func (a *ProcessorAdapter) CanProcess(msg *rabbitmq.StandardMessage) bool {
	return a.processor.CanProcess(msg)
}

// Process processes a message and returns the result in the public adapter format
func (a *ProcessorAdapter) Process(ctx context.Context, msg *rabbitmq.StandardMessage) (*adapter.ProcessedMessage, error) {
	processed, err := a.processor.Process(ctx, msg)
	if err != nil {
		return nil, err
	}

	// Convert to public adapter format
	events := make([]adapter.Event, len(processed.Events))
	for i, e := range processed.Events {
		events[i] = adapter.Event{
			Name:   e.Name,
			Params: e.Params,
			Output: e.Output,
		}
	}

	return &adapter.ProcessedMessage{
		Original:    processed.Original,
		MessageType: adapter.MessageType(processed.MessageType),
		DeviceSN:    processed.DeviceSN,
		Vendor:      processed.Vendor,
		Properties:  processed.Properties,
		Events:      events,
		Timestamp:   processed.Timestamp,
	}, nil
}

// Ensure ProcessorAdapter implements adapter.UplinkProcessor
var _ adapter.UplinkProcessor = (*ProcessorAdapter)(nil)
