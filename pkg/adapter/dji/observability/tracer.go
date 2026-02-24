// Package observability provides metrics, logging, and tracing for the DJI adapter.
package observability

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const (
	// TracerName is the name of the DJI adapter tracer.
	TracerName = "dji-adapter"
)

// Tracer wraps OpenTelemetry tracer for DJI adapter.
type Tracer struct {
	tracer trace.Tracer
}

// NewTracer creates a new DJI adapter tracer.
func NewTracer() *Tracer {
	return &Tracer{
		tracer: otel.Tracer(TracerName),
	}
}

// StartSpan starts a new span for message processing.
func (t *Tracer) StartSpan(ctx context.Context, operationName string) (context.Context, trace.Span) {
	return t.tracer.Start(ctx, operationName)
}

// StartMessageSpan starts a new span for DJI message processing with common attributes.
func (t *Tracer) StartMessageSpan(ctx context.Context, messageType, method, deviceSN string) (context.Context, trace.Span) {
	ctx, span := t.tracer.Start(ctx, "dji.message.process",
		trace.WithAttributes(
			attribute.String("messaging.system", "dji"),
			attribute.String("messaging.message_type", messageType),
			attribute.String("dji.method", method),
			attribute.String("dji.device_sn", deviceSN),
		),
	)
	return ctx, span
}

// startSpan starts a new span with common DJI attributes.
func (t *Tracer) startSpan(ctx context.Context, spanName, attrKey, attrValue, deviceSN string) (context.Context, trace.Span) {
	return t.tracer.Start(ctx, spanName,
		trace.WithAttributes(
			attribute.String("messaging.system", "dji"),
			attribute.String(attrKey, attrValue),
			attribute.String("dji.device_sn", deviceSN),
		),
	)
}

// StartServiceCallSpan starts a new span for service call.
func (t *Tracer) StartServiceCallSpan(ctx context.Context, method, deviceSN string) (context.Context, trace.Span) {
	return t.startSpan(ctx, "dji.service.call", "dji.service.method", method, deviceSN)
}

// StartEventSpan starts a new span for event processing.
func (t *Tracer) StartEventSpan(ctx context.Context, eventType, deviceSN string) (context.Context, trace.Span) {
	return t.startSpan(ctx, "dji.event.process", "dji.event.type", eventType, deviceSN)
}

// RecordError records an error on the span.
func (t *Tracer) RecordError(span trace.Span, err error) {
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
}

// SetSuccess marks the span as successful.
func (t *Tracer) SetSuccess(span trace.Span) {
	span.SetStatus(codes.Ok, "success")
}

// AddAttribute adds an attribute to the span.
func (t *Tracer) AddAttribute(span trace.Span, key, value string) {
	span.SetAttributes(attribute.String(key, value))
}

// spanContextField extracts a field from the span context using the provided extractor function.
// Returns an empty string if the span context is invalid.
func spanContextField(ctx context.Context, extract func(trace.SpanContext) string) string {
	sc := trace.SpanContextFromContext(ctx)
	if sc.IsValid() {
		return extract(sc)
	}
	return ""
}

// GetTraceID extracts the trace ID from context.
func GetTraceID(ctx context.Context) string {
	return spanContextField(ctx, func(sc trace.SpanContext) string { return sc.TraceID().String() })
}

// GetSpanID extracts the span ID from context.
func GetSpanID(ctx context.Context) string {
	return spanContextField(ctx, func(sc trace.SpanContext) string { return sc.SpanID().String() })
}
