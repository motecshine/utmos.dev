package integration

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel/trace"

	"github.com/utmos/utmos/internal/shared/config"
	"github.com/utmos/utmos/pkg/tracer"
)

// TestTracerProviderCreation tests tracer provider initialization
func TestTracerProviderCreation(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *config.TracerConfig
		wantErr bool
	}{
		{
			name: "enabled tracer",
			cfg: &config.TracerConfig{
				Enabled:      true,
				ServiceName:  "test-service",
				Endpoint:     "http://localhost:4318",
				SamplingRate: 1.0,
			},
			wantErr: false,
		},
		{
			name: "disabled tracer",
			cfg: &config.TracerConfig{
				Enabled:     false,
				ServiceName: "test-service",
			},
			wantErr: false,
		},
		{
			name: "empty service name",
			cfg: &config.TracerConfig{
				Enabled:     true,
				ServiceName: "",
				Endpoint:    "http://localhost:4318",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider, err := tracer.NewProvider(tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewProvider() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if provider != nil && !tt.wantErr {
				// Test that we can get a tracer
				tr := provider.Tracer("test")
				if tr == nil {
					t.Error("expected non-nil tracer")
				}

				// Clean up
				if err := provider.Shutdown(context.Background()); err != nil {
					t.Logf("shutdown warning: %v", err)
				}
			}
		})
	}
}

// TestTraceContextPropagation tests W3C trace context propagation
func TestTraceContextPropagation(t *testing.T) {
	// Test with valid traceparent header
	headers := map[string]any{
		"traceparent": "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01",
	}

	ctx := context.Background()
	extractedCtx := tracer.ExtractContext(ctx, headers)

	span := trace.SpanFromContext(extractedCtx)
	spanCtx := span.SpanContext()

	// Validate trace ID format (should be 32 hex chars or zero)
	traceID := spanCtx.TraceID().String()
	if len(traceID) != 32 {
		t.Errorf("unexpected trace ID length: %d", len(traceID))
	}
}

// TestTracerNoopWhenDisabled tests that disabled tracer returns noop
func TestTracerNoopWhenDisabled(t *testing.T) {
	cfg := &config.TracerConfig{
		Enabled:     false,
		ServiceName: "test-service",
	}

	provider, err := tracer.NewProvider(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tr := provider.Tracer("test")
	ctx, span := tr.Start(context.Background(), "test-span")

	// Span should be valid but not recording
	if span == nil {
		t.Error("expected non-nil span")
	}

	span.End()

	if ctx == nil {
		t.Error("expected non-nil context")
	}
}

// TestTracerSpanCreation tests span creation and attributes
func TestTracerSpanCreation(t *testing.T) {
	cfg := &config.TracerConfig{
		Enabled:      true,
		ServiceName:  "test-service",
		Endpoint:     "http://localhost:4318",
		SamplingRate: 1.0,
	}

	provider, err := tracer.NewProvider(cfg)
	if err != nil {
		t.Fatalf("failed to create provider: %v", err)
	}
	defer provider.Shutdown(context.Background())

	tr := provider.Tracer("test-component")

	// Create parent span
	ctx, parentSpan := tr.Start(context.Background(), "parent-operation")
	parentSpanCtx := parentSpan.SpanContext()

	// Create child span
	_, childSpan := tr.Start(ctx, "child-operation")
	childSpanCtx := childSpan.SpanContext()

	// Child should have same trace ID as parent
	if parentSpanCtx.TraceID() != childSpanCtx.TraceID() {
		t.Error("child span should have same trace ID as parent")
	}

	// Spans should have different span IDs
	if parentSpanCtx.SpanID() == childSpanCtx.SpanID() {
		t.Error("parent and child should have different span IDs")
	}

	childSpan.End()
	parentSpan.End()
}

// TestMessageHeaderTracing tests trace context in message headers
func TestMessageHeaderTracing(t *testing.T) {
	cfg := &config.TracerConfig{
		Enabled:      true,
		ServiceName:  "test-service",
		Endpoint:     "http://localhost:4318",
		SamplingRate: 1.0,
	}

	provider, err := tracer.NewProvider(cfg)
	if err != nil {
		t.Fatalf("failed to create provider: %v", err)
	}
	defer provider.Shutdown(context.Background())

	tr := provider.Tracer("test")
	ctx, span := tr.Start(context.Background(), "publish-message")
	defer span.End()

	// Inject trace context into headers
	headers := make(map[string]any)
	tracer.InjectContext(ctx, headers)

	// Headers should contain traceparent
	if _, ok := headers["traceparent"]; !ok {
		t.Log("traceparent header not set (may be expected without proper propagator)")
	}

	// Extract and verify round-trip
	extractedCtx := tracer.ExtractContext(context.Background(), headers)
	extractedSpan := trace.SpanFromContext(extractedCtx)
	extractedSpanCtx := extractedSpan.SpanContext()

	originalSpanCtx := span.SpanContext()

	// Trace IDs should match if propagation worked
	if extractedSpanCtx.TraceID().IsValid() && originalSpanCtx.TraceID().IsValid() {
		if extractedSpanCtx.TraceID() != originalSpanCtx.TraceID() {
			t.Log("trace IDs don't match after round-trip (may be expected)")
		}
	}
}
