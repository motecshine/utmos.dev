package rabbitmq

import (
	"context"
	"testing"

	"github.com/utmos/utmos/internal/shared/config"
)

func TestNewPublisher(t *testing.T) {
	cfg := &config.RabbitMQConfig{
		URL:          "amqp://guest:guest@localhost:5672/",
		ExchangeName: "iot",
		ExchangeType: "topic",
	}

	client := NewClient(cfg)
	publisher := NewPublisher(client)

	if publisher == nil {
		t.Fatal("expected non-nil publisher")
	}
}

func TestPublisher_PublishWithoutConnection(t *testing.T) {
	cfg := &config.RabbitMQConfig{
		URL:          "amqp://guest:guest@localhost:5672/",
		ExchangeName: "iot",
		ExchangeType: "topic",
	}

	client := NewClient(cfg)
	publisher := NewPublisher(client)

	ctx := context.Background()
	routingKey := NewRoutingKey(VendorDJI, ServiceDevice, ActionPropertyReport)

	msg, err := NewStandardMessage(ServiceDevice, ActionPropertyReport, "test-device-001", map[string]interface{}{
		"temperature": 25.5,
	})
	if err != nil {
		t.Fatalf("failed to create message: %v", err)
	}

	// Should fail because not connected
	err = publisher.Publish(ctx, routingKey.String(), msg)
	if err == nil {
		t.Error("expected error when publishing without connection")
	}
}

func TestPublisher_PublishMessageValidation(t *testing.T) {
	cfg := &config.RabbitMQConfig{
		URL:          "amqp://guest:guest@localhost:5672/",
		ExchangeName: "iot",
		ExchangeType: "topic",
	}

	client := NewClient(cfg)
	publisher := NewPublisher(client)

	ctx := context.Background()

	// Test with nil message
	err := publisher.Publish(ctx, "iot.dji.device.property_report", nil)
	if err == nil {
		t.Error("expected error when publishing nil message")
	}
}

func TestPublisher_RoutingKeyFormats(t *testing.T) {
	tests := []struct {
		name       string
		routingKey RoutingKey
		expected   string
	}{
		{
			name:       "DJI property report",
			routingKey: NewRoutingKey(VendorDJI, ServiceDevice, ActionPropertyReport),
			expected:   "iot.dji.device.property_report",
		},
		{
			name:       "Generic device online",
			routingKey: NewRoutingKey(VendorGeneric, ServiceDevice, ActionDeviceOnline),
			expected:   "iot.generic.device.device_online",
		},
		{
			name:       "Tuya service call",
			routingKey: NewRoutingKey(VendorTuya, ServiceService, ActionServiceCall),
			expected:   "iot.tuya.service.service_call",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.routingKey.String() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, tt.routingKey.String())
			}
		})
	}
}

func TestPublisher_MessageWithTraceContext(t *testing.T) {
	cfg := &config.RabbitMQConfig{
		URL:          "amqp://guest:guest@localhost:5672/",
		ExchangeName: "iot",
		ExchangeType: "topic",
	}

	client := NewClient(cfg)
	publisher := NewPublisher(client)

	ctx := context.Background()
	msg, _ := NewStandardMessage(ServiceDevice, ActionPropertyReport, "test-device", map[string]interface{}{})

	// Verify message has required fields
	if msg.TID == "" {
		t.Error("expected TID to be set")
	}
	if msg.BID == "" {
		t.Error("expected BID to be set")
	}
	if msg.Timestamp == 0 {
		t.Error("expected Timestamp to be set")
	}

	// Publishing should fail (not connected) but message should be valid
	_ = publisher.Publish(ctx, "iot.dji.device.property_report", msg)
}
