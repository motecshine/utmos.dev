package tracer

import (
	"context"

	"go.opentelemetry.io/otel"
)

// Ensure MessageCarrier implements propagation.TextMapCarrier interface
var _ interface {
	Get(key string) string
	Set(key, value string)
	Keys() []string
} = (*MessageCarrier)(nil)

// MessageCarrier implements propagation.TextMapCarrier for RabbitMQ message headers.
type MessageCarrier struct {
	Headers map[string]interface{}
}

// Get returns the value for a given key.
func (c *MessageCarrier) Get(key string) string {
	if c.Headers == nil {
		return ""
	}
	if v, ok := c.Headers[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// Set sets a key-value pair.
func (c *MessageCarrier) Set(key, value string) {
	if c.Headers == nil {
		c.Headers = make(map[string]interface{})
	}
	c.Headers[key] = value
}

// Keys returns all keys in the carrier.
func (c *MessageCarrier) Keys() []string {
	if c.Headers == nil {
		return nil
	}
	keys := make([]string, 0, len(c.Headers))
	for k := range c.Headers {
		keys = append(keys, k)
	}
	return keys
}

// InjectContext injects the trace context from ctx into the message headers.
func InjectContext(ctx context.Context, headers map[string]interface{}) {
	if headers == nil {
		return
	}
	propagator := otel.GetTextMapPropagator()
	carrier := &MessageCarrier{Headers: headers}
	propagator.Inject(ctx, carrier)
}

// ExtractContext extracts the trace context from message headers into a new context.
func ExtractContext(ctx context.Context, headers map[string]interface{}) context.Context {
	if headers == nil {
		return ctx
	}
	propagator := otel.GetTextMapPropagator()
	carrier := &MessageCarrier{Headers: headers}
	return propagator.Extract(ctx, carrier)
}
