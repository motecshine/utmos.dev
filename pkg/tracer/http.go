package tracer

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// HTTPMiddleware returns a Gin middleware that adds distributed tracing to HTTP requests.
func HTTPMiddleware(tracer trace.Tracer) gin.HandlerFunc {
	propagator := otel.GetTextMapPropagator()

	return func(c *gin.Context) {
		// Extract trace context from incoming request headers
		ctx := propagator.Extract(c.Request.Context(), propagation.HeaderCarrier(c.Request.Header))

		// Create a new span for this request
		spanName := c.Request.Method + " " + c.FullPath()
		if c.FullPath() == "" {
			spanName = c.Request.Method + " " + c.Request.URL.Path
		}

		ctx, span := tracer.Start(ctx, spanName,
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(
				attribute.String("http.method", c.Request.Method),
				attribute.String("http.url", c.Request.URL.String()),
				attribute.String("http.target", c.Request.URL.Path),
				attribute.String("http.user_agent", c.Request.UserAgent()),
				attribute.String("net.host.name", c.Request.Host),
			),
		)
		defer span.End()

		// Update request with new context
		c.Request = c.Request.WithContext(ctx)

		// Add trace ID to response headers
		spanCtx := span.SpanContext()
		if spanCtx.IsValid() {
			c.Header("X-Trace-ID", spanCtx.TraceID().String())
		}

		// Process request
		c.Next()

		// Record response status
		status := c.Writer.Status()
		span.SetAttributes(attribute.Int("http.status_code", status))

		// Record errors if any
		if len(c.Errors) > 0 {
			span.SetAttributes(attribute.String("error", c.Errors.String()))
		}
	}
}
