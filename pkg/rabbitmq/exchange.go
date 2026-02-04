package rabbitmq

// Dead letter queue constants
const (
	// DeadLetterExchange is the name of the dead letter exchange
	DeadLetterExchange = "iot.dlx"

	// DeadLetterQueuePrefix is the prefix for dead letter queues
	DeadLetterQueuePrefix = "iot.dlq."
)

// DefaultExchangeName is the default topic exchange name
const DefaultExchangeName = "iot"

// DefaultExchangeType is the default exchange type
const DefaultExchangeType = "topic"

// SetupExchange sets up the main IoT topic exchange.
func (c *Client) SetupExchange() error {
	return c.DeclareExchange(c.cfg.ExchangeName, c.cfg.ExchangeType)
}

// SetupDeadLetterExchange sets up the dead letter exchange.
func (c *Client) SetupDeadLetterExchange() error {
	return c.DeclareExchange(DeadLetterExchange, "topic")
}

// SetupQueueWithBinding creates a queue and binds it to the exchange.
func (c *Client) SetupQueueWithBinding(queueName, routingKeyPattern string) error {
	// Declare the queue
	_, err := c.DeclareQueue(queueName, true)
	if err != nil {
		return err
	}

	// Bind the queue to the exchange
	return c.BindQueue(queueName, routingKeyPattern, c.cfg.ExchangeName)
}

// SetupQueueWithDLQ creates a queue with dead letter support and binds it.
func (c *Client) SetupQueueWithDLQ(queueName, routingKeyPattern string) error {
	// Ensure DLX exists
	if err := c.SetupDeadLetterExchange(); err != nil {
		return err
	}

	// Declare DLQ
	dlqName := DeadLetterQueuePrefix + queueName
	_, err := c.DeclareQueue(dlqName, true)
	if err != nil {
		return err
	}

	// Bind DLQ to DLX
	if err := c.BindQueue(dlqName, routingKeyPattern, DeadLetterExchange); err != nil {
		return err
	}

	// Declare main queue with DLX
	_, err = c.DeclareQueueWithDLQ(queueName, DeadLetterExchange)
	if err != nil {
		return err
	}

	// Bind main queue to main exchange
	return c.BindQueue(queueName, routingKeyPattern, c.cfg.ExchangeName)
}
