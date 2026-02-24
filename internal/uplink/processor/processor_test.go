package processor

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/utmos/utmos/pkg/adapter"
	"github.com/utmos/utmos/pkg/rabbitmq"
)

// mockProcessor is a test processor implementation
type mockProcessor struct {
	*BaseProcessor
	canProcessFunc func(msg *rabbitmq.StandardMessage) bool
	processFunc    func(ctx context.Context, msg *rabbitmq.StandardMessage) (*adapter.ProcessedMessage, error)
}

func newMockProcessor(vendor string) *mockProcessor {
	return &mockProcessor{
		BaseProcessor: NewBaseProcessor(vendor, nil),
	}
}

func (p *mockProcessor) CanProcess(msg *rabbitmq.StandardMessage) bool {
	if p.canProcessFunc != nil {
		return p.canProcessFunc(msg)
	}
	return msg.ProtocolMeta != nil && msg.ProtocolMeta.Vendor == p.vendor
}

func (p *mockProcessor) Process(ctx context.Context, msg *rabbitmq.StandardMessage) (*adapter.ProcessedMessage, error) {
	if p.processFunc != nil {
		return p.processFunc(ctx, msg)
	}
	return &adapter.ProcessedMessage{
		Original:    msg,
		MessageType: adapter.MessageTypeProperty,
		DeviceSN:    msg.DeviceSN,
		Vendor:      p.vendor,
		Properties:  make(map[string]any),
		Timestamp:   msg.Timestamp,
	}, nil
}

func TestNewProcessorRegistry(t *testing.T) {
	registry := NewProcessorRegistry(nil)
	require.NotNil(t, registry)
	assert.Equal(t, 0, registry.Count())
}

func TestProcessorRegistry_Register(t *testing.T) {
	registry := NewProcessorRegistry(nil)
	processor := newMockProcessor("test-vendor")

	registry.Register(processor)

	p, exists := registry.Get("test-vendor")
	assert.True(t, exists)
	assert.Equal(t, "test-vendor", p.GetVendor())
}

func TestProcessorRegistry_Unregister(t *testing.T) {
	registry := NewProcessorRegistry(nil)
	processor := newMockProcessor("test-vendor")

	registry.Register(processor)
	registry.Unregister("test-vendor")

	_, exists := registry.Get("test-vendor")
	assert.False(t, exists)
}

func TestProcessorRegistry_GetForMessage(t *testing.T) {
	registry := NewProcessorRegistry(nil)
	djiProcessor := newMockProcessor("dji")
	customProcessor := newMockProcessor("custom")

	registry.Register(djiProcessor)
	registry.Register(customProcessor)

	t.Run("match by protocol meta vendor", func(t *testing.T) {
		msg := &rabbitmq.StandardMessage{
			DeviceSN: "DEVICE001",
			ProtocolMeta: &rabbitmq.ProtocolMeta{
				Vendor: "dji",
			},
		}

		processor, found := registry.GetForMessage(msg)
		assert.True(t, found)
		assert.Equal(t, "dji", processor.GetVendor())
	})

	t.Run("no match", func(t *testing.T) {
		msg := &rabbitmq.StandardMessage{
			DeviceSN: "DEVICE001",
			ProtocolMeta: &rabbitmq.ProtocolMeta{
				Vendor: "unknown",
			},
		}

		_, found := registry.GetForMessage(msg)
		assert.False(t, found)
	})

	t.Run("match by CanProcess", func(t *testing.T) {
		customProcessor.canProcessFunc = func(msg *rabbitmq.StandardMessage) bool {
			return msg.Action == "custom.action"
		}

		msg := &rabbitmq.StandardMessage{
			DeviceSN: "DEVICE001",
			Action:   "custom.action",
		}

		processor, found := registry.GetForMessage(msg)
		assert.True(t, found)
		assert.Equal(t, "custom", processor.GetVendor())
	})
}

func TestProcessorRegistry_ListVendors(t *testing.T) {
	registry := NewProcessorRegistry(nil)
	registry.Register(newMockProcessor("vendor1"))
	registry.Register(newMockProcessor("vendor2"))
	registry.Register(newMockProcessor("vendor3"))

	vendors := registry.ListVendors()
	assert.Len(t, vendors, 3)
	assert.Contains(t, vendors, "vendor1")
	assert.Contains(t, vendors, "vendor2")
	assert.Contains(t, vendors, "vendor3")
}

func TestNewMessageHandler(t *testing.T) {
	registry := NewProcessorRegistry(nil)
	handler := NewMessageHandler(registry, nil)

	require.NotNil(t, handler)
	assert.NotNil(t, handler.registry)
}

func TestMessageHandler_Handle(t *testing.T) {
	registry := NewProcessorRegistry(nil)
	processor := newMockProcessor("dji")
	registry.Register(processor)

	handler := NewMessageHandler(registry, nil)

	t.Run("nil message", func(t *testing.T) {
		err := handler.Handle(context.Background(), nil)
		assert.Error(t, err)
	})

	t.Run("no processor found", func(t *testing.T) {
		msg := &rabbitmq.StandardMessage{
			DeviceSN: "DEVICE001",
			Action:   "test.action",
			ProtocolMeta: &rabbitmq.ProtocolMeta{
				Vendor: "unknown",
			},
		}

		err := handler.Handle(context.Background(), msg)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no processor found")
	})

	t.Run("successful processing", func(t *testing.T) {
		msg := &rabbitmq.StandardMessage{
			TID:       "tid-001",
			BID:       "bid-001",
			DeviceSN:  "DEVICE001",
			Action:    "property.report",
			Timestamp: 1704067200000,
			Data:      json.RawMessage(`{"temperature": 25.5}`),
			ProtocolMeta: &rabbitmq.ProtocolMeta{
				Vendor: "dji",
			},
		}

		var processedMsg *adapter.ProcessedMessage
		handler.SetOnProcessed(func(ctx context.Context, processed *adapter.ProcessedMessage) error {
			processedMsg = processed
			return nil
		})

		err := handler.Handle(context.Background(), msg)
		assert.NoError(t, err)
		require.NotNil(t, processedMsg)
		assert.Equal(t, "DEVICE001", processedMsg.DeviceSN)
		assert.Equal(t, "dji", processedMsg.Vendor)
	})
}

func TestMessageHandler_RegisterProcessor(t *testing.T) {
	registry := NewProcessorRegistry(nil)
	handler := NewMessageHandler(registry, nil)

	processor := newMockProcessor("new-vendor")
	handler.RegisterProcessor(processor)

	p, exists := registry.Get("new-vendor")
	assert.True(t, exists)
	assert.Equal(t, "new-vendor", p.GetVendor())
}

func TestMessageHandler_UnregisterProcessor(t *testing.T) {
	registry := NewProcessorRegistry(nil)
	handler := NewMessageHandler(registry, nil)

	processor := newMockProcessor("test-vendor")
	handler.RegisterProcessor(processor)
	handler.UnregisterProcessor("test-vendor")

	_, exists := registry.Get("test-vendor")
	assert.False(t, exists)
}

func TestBaseProcessor(t *testing.T) {
	processor := NewBaseProcessor("test-vendor", nil)

	assert.Equal(t, "test-vendor", processor.GetVendor())
	assert.NotNil(t, processor.Logger())
}

func TestMessageType(t *testing.T) {
	assert.Equal(t, adapter.MessageType("property"), adapter.MessageTypeProperty)
	assert.Equal(t, adapter.MessageType("event"), adapter.MessageTypeEvent)
	assert.Equal(t, adapter.MessageType("service"), adapter.MessageTypeService)
	assert.Equal(t, adapter.MessageType("status"), adapter.MessageTypeStatus)
}

func TestProcessedMessage(t *testing.T) {
	msg := &adapter.ProcessedMessage{
		MessageType: adapter.MessageTypeProperty,
		DeviceSN:    "DEVICE001",
		Vendor:      "dji",
		Properties: map[string]any{
			"temperature": 25.5,
			"humidity":    60,
		},
		Events: []adapter.Event{
			{
				Name: "alarm",
				Params: map[string]any{
					"level": "warning",
				},
			},
		},
		Timestamp: 1704067200000,
	}

	assert.Equal(t, adapter.MessageTypeProperty, msg.MessageType)
	assert.Equal(t, "DEVICE001", msg.DeviceSN)
	assert.Len(t, msg.Properties, 2)
	assert.Len(t, msg.Events, 1)
}
