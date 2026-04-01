package rabbitmq

import (
	"context"
	"encoding/json"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/utmos/utmos/pkg/metrics"
	"github.com/utmos/utmos/pkg/tracer"
)

// Publisher provides message publishing functionality.
type Publisher struct {
	client       *Client
	exchangeName string
	rmqMetrics  *metrics.RabbitMQMetrics
	serviceName string
}

// NewPublisher creates a new Publisher.
func NewPublisher(client *Client) *Publisher {
	return &Publisher{
		client:       client,
		exchangeName: client.cfg.ExchangeName,
	}
}

// NewPublisherWithMetrics creates a new Publisher with RabbitMQ metrics.
func NewPublisherWithMetrics(client *Client, rmqMetrics *metrics.RabbitMQMetrics, serviceName string) *Publisher {
	return &Publisher{
		client:       client,
		exchangeName: client.cfg.ExchangeName,
		rmqMetrics:  rmqMetrics,
		serviceName: serviceName,
	}
}

// Publish publishes a message with the given routing key.
// It automatically injects W3C Trace Context into the message headers.
func (p *Publisher) Publish(ctx context.Context, routingKey string, msg *StandardMessage) error {
	start := time.Now()
	if !p.client.IsConnected() {
		p.recordPublishError("not_connected")
		return ErrNotConnected
	}

	// Validate message
	if err := msg.Validate(); err != nil {
		p.recordPublishError("validation_error")
		return err
	}

	// Marshal message body
	body, err := json.Marshal(msg)
	if err != nil {
		p.recordPublishError("marshal_error")
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
	err = p.client.Channel().PublishWithContext(
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
	if err != nil {
		p.recordPublishError("publish_error")
	} else {
		p.recordPublishSuccess(routingKey, time.Since(start))
	}
	return err
}

// PublishRaw publishes raw bytes with custom headers for vendor adapter traffic.
// This is used for raw uplink/downlink messages where the body contains
// vendor-specific protocol bytes, not a StandardMessage JSON.
// Headers must include at least "original_topic" and "device_sn".
// W3C Trace Context is automatically injected.
func (p *Publisher) PublishRaw(ctx context.Context, routingKey string, body []byte, headers amqp.Table) error {
	start := time.Now()
	if !p.client.IsConnected() {
		p.recordPublishError("not_connected")
		return ErrNotConnected
	}

	if headers == nil {
		headers = make(amqp.Table)
	}

	// Inject trace context into headers
	headerMap := make(map[string]any)
	tracer.InjectContext(ctx, headerMap)
	for k, v := range headerMap {
		headers[k] = v
	}

	err := p.client.Channel().PublishWithContext(
		ctx,
		p.exchangeName,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/octet-stream",
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
			Headers:      headers,
			Body:         body,
		},
	)
	if err != nil {
		p.recordPublishError("publish_error")
	} else {
		p.recordPublishSuccess(routingKey, time.Since(start))
	}
	return err
}

// PublishWithVendor publishes a message using vendor-based routing key.
func (p *Publisher) PublishWithVendor(ctx context.Context, vendor, service, action string, msg *StandardMessage) error {
	rk := NewRoutingKey(vendor, service, action)
	return p.Publish(ctx, rk.String(), msg)
}

// recordPublishSuccess records successful publish metrics
func (p *Publisher) recordPublishSuccess(routingKey string, duration time.Duration) {
	if p.rmqMetrics == nil {
		return
	}
	p.rmqMetrics.PublishTotal.WithLabelValues(p.serviceName, "success").Inc()
	p.rmqMetrics.MessageDuration.WithLabelValues(p.serviceName).Observe(duration.Seconds())
}

// recordPublishError records failed publish metrics
func (p *Publisher) recordPublishError(status string) {
	if p.rmqMetrics == nil {
		return
	}
	p.rmqMetrics.PublishTotal.WithLabelValues(p.serviceName, status).Inc()
	p.rmqMetrics.ErrorTotal.WithLabelValues(p.serviceName).Inc()
}
