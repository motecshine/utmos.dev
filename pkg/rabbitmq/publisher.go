package rabbitmq

import (
	"context"
	"encoding/json"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/utmos/utmos/pkg/tracer"
)

// Publisher provides message publishing functionality.
type Publisher struct {
	client       *Client
	exchangeName string
}

// NewPublisher creates a new Publisher.
func NewPublisher(client *Client) *Publisher {
	return &Publisher{
		client:       client,
		exchangeName: client.cfg.ExchangeName,
	}
}

// Publish publishes a message with the given routing key.
// It automatically injects W3C Trace Context into the message headers.
func (p *Publisher) Publish(ctx context.Context, routingKey string, msg *StandardMessage) error {
	if !p.client.IsConnected() {
		return ErrNotConnected
	}

	// Validate message
	if err := msg.Validate(); err != nil {
		return err
	}

	// Marshal message body
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	// Create headers and inject trace context
	headers := make(amqp.Table)
	headerMap := make(map[string]any)
	tracer.InjectContext(ctx, headerMap)

	// Copy trace headers to AMQP table
	for k, v := range headerMap {
		headers[k] = v
	}

	// Add message metadata to headers
	headers["message_type"] = msg.Action
	headers["service"] = msg.Service

	// Publish message
	return p.client.Channel().PublishWithContext(
		ctx,
		p.exchangeName, // exchange
		routingKey,     // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
			Headers:      headers,
			Body:         body,
		},
	)
}

// PublishWithVendor publishes a message using vendor-based routing key.
func (p *Publisher) PublishWithVendor(ctx context.Context, vendor, service, action string, msg *StandardMessage) error {
	rk := NewRoutingKey(vendor, service, action)
	return p.Publish(ctx, rk.String(), msg)
}
