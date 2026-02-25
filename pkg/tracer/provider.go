// Package tracer provides distributed tracing using OpenTelemetry.
package tracer

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"

	"github.com/utmos/utmos/pkg/config"
)

// Provider wraps the OpenTelemetry TracerProvider.
type Provider struct {
	provider trace.TracerProvider
	shutdown func(context.Context) error
}

// NewProvider creates a new tracer provider based on configuration.
func NewProvider(cfg *config.TracerConfig) (*Provider, error) {
	if !cfg.Enabled {
		// Return a noop provider when tracing is disabled
		noopProvider := noop.NewTracerProvider()
		return &Provider{
			provider: noopProvider,
			shutdown: func(_ context.Context) error { return nil },
		}, nil
	}

	if cfg.ServiceName == "" {
		return nil, errors.New("service name is required for tracing")
	}

	// Create OTLP HTTP exporter
	exporter, err := otlptracehttp.New(
		context.Background(),
		otlptracehttp.WithEndpointURL(cfg.Endpoint),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	// Create resource with service information
	res, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", cfg.ServiceName),
		),
	)
	if err != nil {
		return nil, err
	}

	// Configure sampler based on sampling rate
	var sampler sdktrace.Sampler
	if cfg.SamplingRate >= 1.0 {
		sampler = sdktrace.AlwaysSample()
	} else if cfg.SamplingRate <= 0 {
		sampler = sdktrace.NeverSample()
	} else {
		sampler = sdktrace.TraceIDRatioBased(cfg.SamplingRate)
	}

	// Create TracerProvider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sampler),
	)

	// Set global TracerProvider and propagator
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return &Provider{
		provider: tp,
		shutdown: tp.Shutdown,
	}, nil
}

// Tracer returns a tracer with the given name.
func (p *Provider) Tracer(name string) trace.Tracer {
	return p.provider.Tracer(name)
}

// Shutdown gracefully shuts down the tracer provider.
func (p *Provider) Shutdown(ctx context.Context) error {
	if p.shutdown != nil {
		return p.shutdown(ctx)
	}
	return nil
}

// getSpanContextField extracts a field from the span context using the provided extractor function.
// Returns an empty string if the span context is invalid.
func getSpanContextField(ctx context.Context, extractor func(trace.SpanContext) string) string {
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.IsValid() {
		return extractor(spanCtx)
	}
	return ""
}

// GetTraceID extracts the trace ID from the context.
//
// minimal wrapper; cannot be reduced further
func GetTraceID(ctx context.Context) string {
	return getSpanContextField(ctx, func(sc trace.SpanContext) string { return sc.TraceID().String() })
}

// GetSpanID extracts the span ID from the context.
//
// minimal wrapper; cannot be reduced further
func GetSpanID(ctx context.Context) string {
	return getSpanContextField(ctx, func(sc trace.SpanContext) string { return sc.SpanID().String() })
}
