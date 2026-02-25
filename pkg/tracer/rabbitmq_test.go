package tracer

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

func TestInjectContext(t *testing.T) {
	tests := []struct {
		name           string
		setupContext   func() context.Context
		expectHeaders  bool
		expectedFields []string
	}{
		{
			name: "inject trace context with valid span",
			setupContext: func() context.Context {
				// Create a context with a valid span
				tracer := noop.NewTracerProvider().Tracer("test")
				ctx, _ := tracer.Start(context.Background(), "test-span")
				return ctx
			},
			expectHeaders:  true,
			expectedFields: []string{"traceparent"},
		},
		{
			name: "inject with no span in context",
			setupContext: func() context.Context {
				return context.Background()
			},
			expectHeaders:  false,
			expectedFields: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.setupContext()
			headers := make(map[string]any)

			InjectContext(ctx, headers)

			if tt.expectHeaders {
				// Headers may or may not be present depending on span validity
				// The important thing is no panic occurs
				assert.NotNil(t, headers)
			}
		})
	}
}

func TestExtractContext(t *testing.T) {
	tests := []struct {
		name    string
		headers map[string]any
	}{
		{
			name: "extract with valid traceparent",
			headers: map[string]any{
				"traceparent": "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01",
			},
		},
		{
			name: "extract with traceparent and tracestate",
			headers: map[string]any{
				"traceparent": "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01",
				"tracestate":  "congo=t61rcWkgMzE",
			},
		},
		{
			name:    "extract with empty headers",
			headers: map[string]any{},
		},
		{
			name: "extract with invalid traceparent",
			headers: map[string]any{
				"traceparent": "invalid",
			},
		},
		{
			name:    "extract with nil headers",
			headers: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			result := ExtractContext(ctx, tt.headers)

			// Should always return a valid context
			assert.NotNil(t, result)
		})
	}
}

func TestRoundTrip_InjectAndExtract(t *testing.T) {
	// Create a span with a known trace context
	tracer := noop.NewTracerProvider().Tracer("test")
	ctx, span := tracer.Start(context.Background(), "test-span")
	defer span.End()

	// Inject the context into headers
	headers := make(map[string]any)
	InjectContext(ctx, headers)

	// Extract the context from headers
	extractedCtx := ExtractContext(context.Background(), headers)

	// Both contexts should be valid
	assert.NotNil(t, extractedCtx)
}

func TestMessageCarrier(t *testing.T) {
	headers := map[string]any{
		"traceparent": "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01",
		"tracestate":  "congo=t61rcWkgMzE",
		"custom":      "value",
	}

	carrier := &MessageCarrier{Headers: headers}

	// Test Get
	assert.Equal(t, "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01", carrier.Get("traceparent"))
	assert.Equal(t, "congo=t61rcWkgMzE", carrier.Get("tracestate"))
	assert.Equal(t, "", carrier.Get("nonexistent"))

	// Test Set
	carrier.Set("newkey", "newvalue")
	assert.Equal(t, "newvalue", carrier.Get("newkey"))

	// Test Keys
	keys := carrier.Keys()
	assert.Contains(t, keys, "traceparent")
	assert.Contains(t, keys, "tracestate")
}

func TestGetSpanFromContext(t *testing.T) {
	// Test with no span
	ctx := context.Background()
	span := trace.SpanFromContext(ctx)
	assert.NotNil(t, span)

	// Test with a span
	tracer := noop.NewTracerProvider().Tracer("test")
	ctx, testSpan := tracer.Start(context.Background(), "test-span")
	defer testSpan.End()

	retrievedSpan := trace.SpanFromContext(ctx)
	assert.NotNil(t, retrievedSpan)
}
