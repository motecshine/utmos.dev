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

// StartServiceCallSpan starts a new span for service call.
func (t *Tracer) StartServiceCallSpan(ctx context.Context, method, deviceSN string) (context.Context, trace.Span) {
	ctx, span := t.tracer.Start(ctx, "dji.service.call",
		trace.WithAttributes(
			attribute.String("messaging.system", "dji"),
			attribute.String("dji.service.method", method),
			attribute.String("dji.device_sn", deviceSN),
		),
	)
	return ctx, span
}

// StartEventSpan starts a new span for event processing.
func (t *Tracer) StartEventSpan(ctx context.Context, eventType, deviceSN string) (context.Context, trace.Span) {
	ctx, span := t.tracer.Start(ctx, "dji.event.process",
		trace.WithAttributes(
			attribute.String("messaging.system", "dji"),
			attribute.String("dji.event.type", eventType),
			attribute.String("dji.device_sn", deviceSN),
		),
	)
	return ctx, span
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

// GetTraceID extracts the trace ID from context.
func GetTraceID(ctx context.Context) string {
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.IsValid() {
		return spanCtx.TraceID().String()
	}
	return ""
}

// GetSpanID extracts the span ID from context.
func GetSpanID(ctx context.Context) string {
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.IsValid() {
		return spanCtx.SpanID().String()
	}
	return ""
}
