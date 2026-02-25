package observability

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/trace"
)

// HandlerObserver provides observability for message handlers.
type HandlerObserver struct {
	metrics *Metrics
	tracer  *Tracer
	logger  *Logger
}

// NewHandlerObserver creates a new handler observer.
func NewHandlerObserver(metrics *Metrics, tracer *Tracer, logger *Logger) *HandlerObserver {
	return &HandlerObserver{
		metrics: metrics,
		tracer:  tracer,
		logger:  logger,
	}
}

// ObserveResult holds the result of an observed operation.
type ObserveResult struct {
	ctx       context.Context
	span      trace.Span
	startTime time.Time
	observer  *HandlerObserver
	msgType   string
	deviceSN  string
	method    string
}

// StartObserve starts observing a message processing operation.
func (h *HandlerObserver) StartObserve(ctx context.Context, messageType, method, deviceSN string) *ObserveResult {
	ctx, span := h.tracer.StartMessageSpan(ctx, messageType, method, deviceSN)

	h.logger.WithMessage(ctx, deviceSN, messageType, method).Debug("processing message")

	return &ObserveResult{
		ctx:       ctx,
		span:      span,
		startTime: time.Now(),
		observer:  h,
		msgType:   messageType,
		deviceSN:  deviceSN,
		method:    method,
	}
}

// Context returns the context with span.
func (r *ObserveResult) Context() context.Context {
	return r.ctx
}

// End ends the observation with success.
func (r *ObserveResult) End() {
	duration := time.Since(r.startTime).Seconds()

	r.observer.metrics.RecordMessageReceived(r.msgType, "success")
	r.observer.metrics.RecordProcessingDuration(r.msgType, duration)
	r.observer.tracer.SetSuccess(r.span)

	r.observer.logger.WithMessage(r.ctx, r.deviceSN, r.msgType, r.method).
		WithField("duration_ms", duration*1000).
		Debug("message processed successfully")

	r.span.End()
}

// EndWithError ends the observation with an error.
func (r *ObserveResult) EndWithError(err error, errorType string) {
	duration := time.Since(r.startTime).Seconds()

	r.observer.metrics.RecordMessageReceived(r.msgType, "error")
	r.observer.metrics.RecordError(r.msgType, errorType)
	r.observer.metrics.RecordProcessingDuration(r.msgType, duration)
	r.observer.tracer.RecordError(r.span, err)

	r.observer.logger.WithMessage(r.ctx, r.deviceSN, r.msgType, r.method).
		WithError(err).
		WithField("duration_ms", duration*1000).
		Error("message processing failed")

	r.span.End()
}
