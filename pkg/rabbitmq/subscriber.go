package rabbitmq

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/utmos/utmos/pkg/tracer"
)

// MessageHandler is a function that handles incoming messages.
// Return an error to Nack the message.
type MessageHandler func(ctx context.Context, msg *StandardMessage) error

// Subscriber provides message subscription functionality.
type Subscriber struct {
	client   *Client
	handlers map[string]chan struct{}
}

// NewSubscriber creates a new Subscriber.
func NewSubscriber(client *Client) *Subscriber {
	return &Subscriber{
		client:   client,
		handlers: make(map[string]chan struct{}),
	}
}

// Subscribe subscribes to messages from a queue.
// Messages are processed with manual acknowledgment.
func (s *Subscriber) Subscribe(queueName string, handler MessageHandler) error {
	if !s.client.IsConnected() {
		return ErrNotConnected
	}

	msgs, err := s.client.Channel().Consume(
		queueName, // queue
		"",        // consumer tag (auto-generated)
		false,     // auto-ack (we use manual ack)
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return err
	}

	stopChan := make(chan struct{})
	s.handlers[queueName] = stopChan

	go s.processMessages(msgs, handler, stopChan)

	return nil
}

// processMessages processes incoming messages.
func (s *Subscriber) processMessages(msgs <-chan amqp.Delivery, handler MessageHandler, stopChan chan struct{}) {
	for {
		select {
		case <-stopChan:
			return
		case delivery, ok := <-msgs:
			if !ok {
				return
			}
			s.handleDelivery(delivery, handler)
		}
	}
}

// handleDelivery handles a single message delivery.
func (s *Subscriber) handleDelivery(delivery amqp.Delivery, handler MessageHandler) {
	// Extract trace context from headers
	headerMap := make(map[string]interface{})
	for k, v := range delivery.Headers {
		headerMap[k] = v
	}
	ctx := tracer.ExtractContext(context.Background(), headerMap)

	// Parse message
	msg, err := FromBytes(delivery.Body)
	if err != nil {
		// Nack invalid messages without requeue
		_ = delivery.Nack(false, false)
		return
	}

	// Call handler
	if err := handler(ctx, msg); err != nil {
		// Nack with requeue on handler error
		_ = delivery.Nack(false, true)
		return
	}

	// Ack on success
	_ = delivery.Ack(false)
}

// Unsubscribe stops consuming from a queue.
func (s *Subscriber) Unsubscribe(queueName string) error {
	if stopChan, ok := s.handlers[queueName]; ok {
		close(stopChan)
		delete(s.handlers, queueName)
	}
	return nil
}

// UnsubscribeAll stops all subscriptions.
func (s *Subscriber) UnsubscribeAll() {
	for queueName := range s.handlers {
		_ = s.Unsubscribe(queueName)
	}
}
