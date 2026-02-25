package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const (
	// TraceIDHeader is the header name for trace ID
	TraceIDHeader = "X-Trace-ID"
	// RequestIDHeader is the header name for request ID
	RequestIDHeader = "X-Request-ID"
)

// TraceConfig holds tracing configuration
type TraceConfig struct {
	// ServiceName is the name of the service
	ServiceName string
	// TracerName is the name of the tracer
	TracerName string
	// SkipPaths are paths that don't need tracing
	SkipPaths []string
}

// DefaultTraceConfig returns default trace configuration
func DefaultTraceConfig() *TraceConfig {
	return &TraceConfig{
		ServiceName: "iot-api",
		TracerName:  "iot-api",
		SkipPaths:   []string{"/health", "/ready", "/metrics"},
	}
}

// TraceMiddleware provides distributed tracing
type TraceMiddleware struct {
	config     *TraceConfig
	logger     *logrus.Entry
	tracer     trace.Tracer
	propagator propagation.TextMapPropagator
}

// NewTraceMiddleware creates a new trace middleware
func NewTraceMiddleware(config *TraceConfig, logger *logrus.Entry) *TraceMiddleware {
	if config == nil {
		config = DefaultTraceConfig()
	}
	if logger == nil {
		logger = logrus.NewEntry(logrus.StandardLogger())
	}

	return &TraceMiddleware{
		config:     config,
		logger:     logger.WithField("middleware", "trace"),
		tracer:     otel.Tracer(config.TracerName),
		propagator: otel.GetTextMapPropagator(),
	}
}

// Handler returns the Gin middleware handler
func (m *TraceMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if path should be skipped
		path := c.Request.URL.Path
		for _, skipPath := range m.config.SkipPaths {
			if path == skipPath {
				c.Next()
				return
			}
		}

		// Extract trace context from incoming request
		ctx := m.propagator.Extract(c.Request.Context(), propagation.HeaderCarrier(c.Request.Header))

		// Generate or get request ID
		requestID := c.GetHeader(RequestIDHeader)
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Start span
		spanName := c.Request.Method + " " + c.FullPath()
		if c.FullPath() == "" {
			spanName = c.Request.Method + " " + path
		}

		ctx, span := m.tracer.Start(ctx, spanName,
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(
				attribute.String("http.method", c.Request.Method),
				attribute.String("http.url", c.Request.URL.String()),
				attribute.String("http.path", path),
				attribute.String("http.host", c.Request.Host),
				attribute.String("http.user_agent", c.Request.UserAgent()),
				attribute.String("request.id", requestID),
			),
		)
		defer span.End()

		// Get trace ID
		traceID := span.SpanContext().TraceID().String()
		spanID := span.SpanContext().SpanID().String()

		// Set context and headers
		c.Request = c.Request.WithContext(ctx)
		c.Set("trace_id", traceID)
		c.Set("span_id", spanID)
		c.Set("request_id", requestID)

		// Set response headers
		c.Header(TraceIDHeader, traceID)
		c.Header(RequestIDHeader, requestID)

		// Record start time
		start := time.Now()

		// Process request
		c.Next()

		// Record response attributes
		duration := time.Since(start)
		statusCode := c.Writer.Status()

		span.SetAttributes(
			attribute.Int("http.status_code", statusCode),
			attribute.Int64("http.response_size", int64(c.Writer.Size())),
			attribute.Float64("http.duration_ms", float64(duration.Milliseconds())),
		)

		// Log request
		m.logger.WithFields(logrus.Fields{
			"trace_id":    traceID,
			"span_id":     spanID,
			"request_id":  requestID,
			"method":      c.Request.Method,
			"path":        path,
			"status":      statusCode,
			"duration_ms": duration.Milliseconds(),
			"client_ip":   c.ClientIP(),
		}).Info("Request completed")
	}
}

// RequestLogger returns a middleware that logs requests with trace context
func RequestLogger(logger *logrus.Entry) gin.HandlerFunc {
	if logger == nil {
		logger = logrus.NewEntry(logrus.StandardLogger())
	}

	return func(c *gin.Context) {
		start := time.Now()

		// Generate request ID if not present
		requestID := c.GetHeader(RequestIDHeader)
		if requestID == "" {
			requestID = uuid.New().String()
			c.Header(RequestIDHeader, requestID)
		}
		c.Set("request_id", requestID)

		// Process request
		c.Next()

		// Log after request
		duration := time.Since(start)
		logger.WithFields(logrus.Fields{
			"request_id":  requestID,
			"method":      c.Request.Method,
			"path":        c.Request.URL.Path,
			"status":      c.Writer.Status(),
			"duration_ms": duration.Milliseconds(),
			"client_ip":   c.ClientIP(),
		}).Info("Request completed")
	}
}

// InjectTraceContext injects trace context into outgoing requests
func InjectTraceContext(c *gin.Context) map[string]string {
	headers := make(map[string]string)

	if traceID := c.GetString("trace_id"); traceID != "" {
		headers[TraceIDHeader] = traceID
	}
	if requestID := c.GetString("request_id"); requestID != "" {
		headers[RequestIDHeader] = requestID
	}

	return headers
}

// GetTraceID returns the trace ID from context
func GetTraceID(c *gin.Context) string {
	return c.GetString("trace_id")
}

// GetRequestID returns the request ID from context
func GetRequestID(c *gin.Context) string {
	return c.GetString("request_id")
}
