// Package processor provides message processing functionality for iot-uplink
package processor

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/utmos/utmos/pkg/adapter"
	"github.com/utmos/utmos/pkg/rabbitmq"
	"github.com/utmos/utmos/pkg/registry"
)

// Handler handles incoming messages and routes them to appropriate processors
type Handler interface {
	// Handle processes an incoming message
	Handle(ctx context.Context, msg *rabbitmq.StandardMessage) error

	// RegisterProcessor registers a processor for a vendor
	RegisterProcessor(processor adapter.UplinkProcessor)

	// UnregisterProcessor unregisters a processor
	UnregisterProcessor(vendor string)
}

// Registry manages registered processors using the generic registry
type Registry struct {
	*registry.Registry[adapter.UplinkProcessor]
}

// NewRegistry creates a new processor registry
func NewRegistry(logger *logrus.Entry) *Registry {
	return &Registry{
		Registry: registry.New[adapter.UplinkProcessor]("processor-registry", logger),
	}
}

// GetForMessage returns a processor that can handle the given message
func (r *Registry) GetForMessage(msg *rabbitmq.StandardMessage) (adapter.UplinkProcessor, bool) {
	// First try to match by protocol meta vendor
	vendor := ""
	if msg.ProtocolMeta != nil {
		vendor = msg.ProtocolMeta.Vendor
	}

	// Use GetOrFind: first try vendor lookup, then predicate
	return r.GetOrFind(vendor, func(p adapter.UplinkProcessor) bool {
		return p.CanProcess(msg)
	})
}

// MessageHandler handles incoming messages using registered processors
type MessageHandler struct {
	registry    *Registry
	logger      *logrus.Entry
	onProcessed func(ctx context.Context, processed *adapter.ProcessedMessage) error
}

// NewMessageHandler creates a new message handler
func NewMessageHandler(registry *Registry, logger *logrus.Entry) *MessageHandler {
	if logger == nil {
		logger = logrus.NewEntry(logrus.StandardLogger())
	}
	return &MessageHandler{
		registry: registry,
		logger:   logger.WithField("component", "message-handler"),
	}
}

// SetOnProcessed sets the callback for processed messages
func (h *MessageHandler) SetOnProcessed(callback func(ctx context.Context, processed *adapter.ProcessedMessage) error) {
	h.onProcessed = callback
}

// Handle processes an incoming message
func (h *MessageHandler) Handle(ctx context.Context, msg *rabbitmq.StandardMessage) error {
	if msg == nil {
		return fmt.Errorf("message is nil")
	}

	tr := otel.Tracer("iot-uplink")
	ctx, span := tr.Start(ctx, "uplink.message.process",
		trace.WithAttributes(
			attribute.String("device_sn", msg.DeviceSN),
			attribute.String("action", msg.Action),
		),
	)
	defer span.End()

	// Find appropriate processor
	processor, found := h.registry.GetForMessage(msg)
	if !found {
		h.logger.WithFields(logrus.Fields{
			"device_sn": msg.DeviceSN,
			"action":    msg.Action,
			"tid":       msg.TID,
		}).Warn("No processor found for message")
		err := fmt.Errorf("no processor found for message")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	// Process the message
	processed, err := processor.Process(ctx, msg)
	if err != nil {
		h.logger.WithError(err).WithFields(logrus.Fields{
			"device_sn": msg.DeviceSN,
			"vendor":    processor.GetVendor(),
			"tid":       msg.TID,
		}).Error("Failed to process message")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return fmt.Errorf("failed to process message: %w", err)
	}

	h.logger.WithFields(logrus.Fields{
		"device_sn":    processed.DeviceSN,
		"vendor":       processed.Vendor,
		"message_type": processed.MessageType,
		"tid":          msg.TID,
	}).Debug("Message processed")

	// Call the processed callback if set
	if h.onProcessed != nil {
		if err := h.onProcessed(ctx, processed); err != nil {
			h.logger.WithError(err).WithField("tid", msg.TID).Error("Failed to handle processed message")
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return fmt.Errorf("failed to handle processed message: %w", err)
		}
	}

	return nil
}

// RegisterProcessor registers a processor
func (h *MessageHandler) RegisterProcessor(processor adapter.UplinkProcessor) {
	h.registry.Register(processor)
}

// UnregisterProcessor unregisters a processor
func (h *MessageHandler) UnregisterProcessor(vendor string) {
	h.registry.Unregister(vendor)
}

// BaseProcessor provides common functionality for processors
type BaseProcessor struct {
	vendor string
	logger *logrus.Entry
}

// NewBaseProcessor creates a new base processor
func NewBaseProcessor(vendor string, logger *logrus.Entry) *BaseProcessor {
	if logger == nil {
		logger = logrus.NewEntry(logrus.StandardLogger())
	}
	return &BaseProcessor{
		vendor: vendor,
		logger: logger.WithField("processor", vendor),
	}
}

// GetVendor returns the vendor name
func (p *BaseProcessor) GetVendor() string {
	return p.vendor
}

// Logger returns the processor's logger
func (p *BaseProcessor) Logger() *logrus.Entry {
	return p.logger
}
