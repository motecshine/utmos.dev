package rabbitmq

import (
	"context"
	"testing"
	"time"

	"github.com/utmos/utmos/internal/shared/config"
)

func TestNewClient(t *testing.T) {
	cfg := &config.RabbitMQConfig{
		URL:          "amqp://guest:guest@localhost:5672/",
		ExchangeName: "iot",
		ExchangeType: "topic",
	}

	client := NewClient(cfg)
	if client == nil {
		t.Fatal("expected non-nil client")
	}

	// Client should not be connected initially (expected behavior)
	_ = client.IsConnected()
}

func TestClientIsConnected(t *testing.T) {
	cfg := &config.RabbitMQConfig{
		URL:          "amqp://guest:guest@localhost:5672/",
		ExchangeName: "iot",
		ExchangeType: "topic",
	}

	client := NewClient(cfg)

	// Should not be connected initially
	if client.IsConnected() {
		t.Error("expected client to not be connected initially")
	}
}

func TestClientClose(t *testing.T) {
	cfg := &config.RabbitMQConfig{
		URL:          "amqp://guest:guest@localhost:5672/",
		ExchangeName: "iot",
		ExchangeType: "topic",
	}

	client := NewClient(cfg)

	// Close should not error even if not connected
	err := client.Close()
	if err != nil {
		t.Errorf("expected no error on close, got %v", err)
	}
}

func TestClientConnectTimeout(t *testing.T) {
	cfg := &config.RabbitMQConfig{
		URL:          "amqp://guest:guest@invalid-host-that-does-not-exist:5672/",
		ExchangeName: "iot",
		ExchangeType: "topic",
		Retry: config.RetryConfig{
			MaxRetries:   1,
			InitialDelay: 100 * time.Millisecond,
			MaxDelay:     100 * time.Millisecond,
			Multiplier:   1.0,
		},
	}

	client := NewClient(cfg)
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	err := client.Connect(ctx)
	if err == nil {
		t.Error("expected connection to fail")
		client.Close()
	}
}

func TestRetryBackoff(t *testing.T) {
	tests := []struct {
		attempt  int
		base     time.Duration
		max      time.Duration
		expected time.Duration
	}{
		{0, time.Second, 30 * time.Second, time.Second},
		{1, time.Second, 30 * time.Second, 2 * time.Second},
		{2, time.Second, 30 * time.Second, 4 * time.Second},
		{3, time.Second, 30 * time.Second, 8 * time.Second},
		{4, time.Second, 30 * time.Second, 16 * time.Second},
		{5, time.Second, 30 * time.Second, 30 * time.Second}, // capped at max
		{10, time.Second, 30 * time.Second, 30 * time.Second},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := calculateBackoff(tt.attempt, tt.base, tt.max)
			if result != tt.expected {
				t.Errorf("attempt %d: expected %v, got %v", tt.attempt, tt.expected, result)
			}
		})
	}
}

// calculateBackoff is a helper for testing exponential backoff
func calculateBackoff(attempt int, base, maxDelay time.Duration) time.Duration {
	backoff := base * (1 << attempt)
	if backoff > maxDelay {
		return maxDelay
	}
	return backoff
}

func TestClientDeclareExchangeNotConnected(t *testing.T) {
	cfg := &config.RabbitMQConfig{
		URL:          "amqp://guest:guest@localhost:5672/",
		ExchangeName: "iot",
		ExchangeType: "topic",
	}

	client := NewClient(cfg)

	// Should fail because not connected
	err := client.DeclareExchange("test-exchange", "topic")
	if err == nil {
		t.Error("expected error when declaring exchange without connection")
	}
	if err != ErrNotConnected {
		t.Errorf("expected ErrNotConnected, got %v", err)
	}
}

func TestClientDeclareQueueNotConnected(t *testing.T) {
	cfg := &config.RabbitMQConfig{
		URL:          "amqp://guest:guest@localhost:5672/",
		ExchangeName: "iot",
		ExchangeType: "topic",
	}

	client := NewClient(cfg)

	// Should fail because not connected
	_, err := client.DeclareQueue("test-queue", true)
	if err == nil {
		t.Error("expected error when declaring queue without connection")
	}
	if err != ErrNotConnected {
		t.Errorf("expected ErrNotConnected, got %v", err)
	}
}

func TestClientDeclareQueueWithDLQNotConnected(t *testing.T) {
	cfg := &config.RabbitMQConfig{
		URL:          "amqp://guest:guest@localhost:5672/",
		ExchangeName: "iot",
		ExchangeType: "topic",
	}

	client := NewClient(cfg)

	// Should fail because not connected
	_, err := client.DeclareQueueWithDLQ("test-queue", "dlx")
	if err == nil {
		t.Error("expected error when declaring queue with DLQ without connection")
	}
	if err != ErrNotConnected {
		t.Errorf("expected ErrNotConnected, got %v", err)
	}
}

func TestClientBindQueueNotConnected(t *testing.T) {
	cfg := &config.RabbitMQConfig{
		URL:          "amqp://guest:guest@localhost:5672/",
		ExchangeName: "iot",
		ExchangeType: "topic",
	}

	client := NewClient(cfg)

	// Should fail because not connected
	err := client.BindQueue("test-queue", "test.routing.key", "test-exchange")
	if err == nil {
		t.Error("expected error when binding queue without connection")
	}
	if err != ErrNotConnected {
		t.Errorf("expected ErrNotConnected, got %v", err)
	}
}

func TestClientChannelNotConnected(t *testing.T) {
	cfg := &config.RabbitMQConfig{
		URL:          "amqp://guest:guest@localhost:5672/",
		ExchangeName: "iot",
		ExchangeType: "topic",
	}

	client := NewClient(cfg)

	// Channel should be nil when not connected
	ch := client.Channel()
	if ch != nil {
		t.Error("expected nil channel when not connected")
	}
}

func TestClientConnectContextCanceled(t *testing.T) {
	cfg := &config.RabbitMQConfig{
		URL:          "amqp://guest:guest@invalid-host:5672/",
		ExchangeName: "iot",
		ExchangeType: "topic",
		Retry: config.RetryConfig{
			MaxRetries:   10,
			InitialDelay: 100 * time.Millisecond,
			MaxDelay:     1 * time.Second,
			Multiplier:   2.0,
		},
	}

	client := NewClient(cfg)

	// Cancel context immediately
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := client.Connect(ctx)
	if err == nil {
		t.Error("expected error when context is canceled")
		client.Close()
	}
}
