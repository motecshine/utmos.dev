package rabbitmq

import (
	"context"
	"fmt"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/utmos/utmos/internal/shared/config"
)

// Client represents a RabbitMQ client.
type Client struct {
	cfg        *config.RabbitMQConfig
	conn       *amqp.Connection
	channel    *amqp.Channel
	mu         sync.RWMutex
	connected  bool
	closeChan  chan struct{}
}

// NewClient creates a new RabbitMQ client.
func NewClient(cfg *config.RabbitMQConfig) *Client {
	return &Client{
		cfg:       cfg,
		closeChan: make(chan struct{}),
	}
}

// Connect connects to RabbitMQ with exponential backoff retry.
func (c *Client) Connect(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.connected {
		return nil
	}

	var err error
	delay := c.cfg.Retry.InitialDelay
	maxRetries := c.cfg.Retry.MaxRetries

	for attempt := 0; attempt < maxRetries; attempt++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		c.conn, err = amqp.Dial(c.cfg.URL)
		if err == nil {
			c.channel, err = c.conn.Channel()
			if err == nil {
				// Set prefetch count
				if c.cfg.PrefetchCount > 0 {
					err = c.channel.Qos(c.cfg.PrefetchCount, 0, false)
					if err != nil {
						c.conn.Close()
						continue
					}
				}
				c.connected = true
				return nil
			}
			c.conn.Close()
		}

		// Exponential backoff
		if attempt < maxRetries-1 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
			}

			delay = time.Duration(float64(delay) * c.cfg.Retry.Multiplier)
			if delay > c.cfg.Retry.MaxDelay {
				delay = c.cfg.Retry.MaxDelay
			}
		}
	}

	return fmt.Errorf("failed to connect to RabbitMQ after %d attempts: %w", maxRetries, err)
}

// Close closes the RabbitMQ connection.
func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	close(c.closeChan)

	var errs []error
	if c.channel != nil {
		if err := c.channel.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	c.connected = false

	if len(errs) > 0 {
		return errs[0]
	}
	return nil
}

// IsConnected returns true if the client is connected.
func (c *Client) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.connected && c.conn != nil && !c.conn.IsClosed()
}

// DeclareExchange declares an exchange (idempotent).
func (c *Client) DeclareExchange(name, exchangeType string) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.connected {
		return ErrNotConnected
	}

	return c.channel.ExchangeDeclare(
		name,         // name
		exchangeType, // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
}

// DeclareQueue declares a queue (idempotent).
func (c *Client) DeclareQueue(name string, durable bool) (amqp.Queue, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.connected {
		return amqp.Queue{}, ErrNotConnected
	}

	return c.channel.QueueDeclare(
		name,    // name
		durable, // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
}

// DeclareQueueWithDLQ declares a queue with dead letter queue support.
func (c *Client) DeclareQueueWithDLQ(name string, dlxName string) (amqp.Queue, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.connected {
		return amqp.Queue{}, ErrNotConnected
	}

	args := amqp.Table{
		"x-dead-letter-exchange": dlxName,
	}

	return c.channel.QueueDeclare(
		name,  // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		args,  // arguments
	)
}

// BindQueue binds a queue to an exchange.
func (c *Client) BindQueue(queueName, routingKey, exchangeName string) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.connected {
		return ErrNotConnected
	}

	return c.channel.QueueBind(
		queueName,    // queue name
		routingKey,   // routing key
		exchangeName, // exchange
		false,        // no-wait
		nil,          // arguments
	)
}

// Channel returns the underlying AMQP channel.
func (c *Client) Channel() *amqp.Channel {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.channel
}

// ErrNotConnected is returned when an operation is attempted on a disconnected client.
var ErrNotConnected = fmt.Errorf("rabbitmq client is not connected")
