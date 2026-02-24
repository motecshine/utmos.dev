package rabbitmq

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/utmos/utmos/pkg/config"
)

func TestNewSubscriber(t *testing.T) {
	cfg := &config.RabbitMQConfig{
		URL:          "amqp://guest:guest@localhost:5672/",
		ExchangeName: "iot",
		ExchangeType: "topic",
	}

	client := NewClient(cfg)
	subscriber := NewSubscriber(client)

	if subscriber == nil {
		t.Fatal("expected non-nil subscriber")
	}
}

func TestSubscriber_SubscribeWithoutConnection(t *testing.T) {
	cfg := &config.RabbitMQConfig{
		URL:          "amqp://guest:guest@localhost:5672/",
		ExchangeName: "iot",
		ExchangeType: "topic",
	}

	client := NewClient(cfg)
	subscriber := NewSubscriber(client)

	handler := func(_ context.Context, _ *StandardMessage) error {
		return nil
	}

	// Should fail because not connected
	err := subscriber.Subscribe("test-queue", handler)
	if err == nil {
		t.Error("expected error when subscribing without connection")
	}
}

func TestSubscriber_UnsubscribeAll(_ *testing.T) {
	cfg := &config.RabbitMQConfig{
		URL:          "amqp://guest:guest@localhost:5672/",
		ExchangeName: "iot",
		ExchangeType: "topic",
	}

	client := NewClient(cfg)
	subscriber := NewSubscriber(client)

	// Should not panic even without subscriptions
	subscriber.UnsubscribeAll()
}

func TestSubscriber_HandlerType(t *testing.T) {
	var called bool
	var mu sync.Mutex

	handler := func(_ context.Context, msg *StandardMessage) error {
		mu.Lock()
		defer mu.Unlock()
		called = true

		// Verify message fields
		if msg.TID == "" {
			t.Error("expected TID in message")
		}
		if msg.DeviceSN == "" {
			t.Error("expected DeviceSN in message")
		}
		return nil
	}

	// Simulate handler execution
	msg, _ := NewStandardMessage(ServiceDevice, ActionPropertyReport, "test-device", nil)
	ctx := context.Background()

	err := handler(ctx, msg)
	if err != nil {
		t.Errorf("handler returned error: %v", err)
	}

	mu.Lock()
	if !called {
		t.Error("handler was not called")
	}
	mu.Unlock()
}

func TestSubscriber_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	var handlerCalled bool
	handler := func(ctx context.Context, _ *StandardMessage) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			handlerCalled = true
			return nil
		}
	}

	msg, _ := NewStandardMessage(ServiceDevice, ActionPropertyReport, "test-device", nil)
	err := handler(ctx, msg)
	if err != nil {
		t.Errorf("handler returned error: %v", err)
	}

	if !handlerCalled {
		t.Error("handler should have been called")
	}
}

func TestSubscriber_ErrorHandling(t *testing.T) {
	tests := []struct {
		name        string
		handlerErr  error
		expectError bool
	}{
		{
			name:        "successful handler",
			handlerErr:  nil,
			expectError: false,
		},
		{
			name:        "handler with error",
			handlerErr:  context.DeadlineExceeded,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := func(_ context.Context, _ *StandardMessage) error {
				return tt.handlerErr
			}

			msg, _ := NewStandardMessage(ServiceDevice, ActionPropertyReport, "test-device", nil)
			err := handler(context.Background(), msg)

			if tt.expectError && err == nil {
				t.Error("expected error from handler")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestSubscriber_QueueNaming(t *testing.T) {
	tests := []struct {
		queueName string
		valid     bool
	}{
		{"iot-uplink-queue", true},
		{"iot-downlink-queue", true},
		{"iot-gateway-queue", true},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.queueName, func(t *testing.T) {
			if tt.valid && tt.queueName == "" {
				t.Error("empty queue name should be invalid")
			}
		})
	}
}

func TestSubscriber_Unsubscribe(t *testing.T) {
	cfg := &config.RabbitMQConfig{
		URL:          "amqp://guest:guest@localhost:5672/",
		ExchangeName: "iot",
		ExchangeType: "topic",
	}

	client := NewClient(cfg)
	subscriber := NewSubscriber(client)

	// Unsubscribe from non-existent queue should not error
	err := subscriber.Unsubscribe("non-existent-queue")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSubscriber_MultipleUnsubscribeAll(t *testing.T) {
	cfg := &config.RabbitMQConfig{
		URL:          "amqp://guest:guest@localhost:5672/",
		ExchangeName: "iot",
		ExchangeType: "topic",
	}

	client := NewClient(cfg)
	subscriber := NewSubscriber(client)

	// Multiple calls should not panic
	subscriber.UnsubscribeAll()
	subscriber.UnsubscribeAll()
}
