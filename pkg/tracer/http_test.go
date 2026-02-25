package tracer

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace/noop"
)

func TestHTTPMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Use a noop tracer for testing
	tracer := noop.NewTracerProvider().Tracer("test")

	tests := []struct {
		name           string
		method         string
		path           string
		traceparent    string
		expectedStatus int
	}{
		{
			name:           "GET request without trace context",
			method:         http.MethodGet,
			path:           "/test",
			traceparent:    "",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "POST request without trace context",
			method:         http.MethodPost,
			path:           "/api/v1/devices",
			traceparent:    "",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "GET request with valid trace context",
			method:         http.MethodGet,
			path:           "/test",
			traceparent:    "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.Use(HTTPMiddleware(tracer))
			router.Any("/*path", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			req := httptest.NewRequest(tt.method, tt.path, nil)
			if tt.traceparent != "" {
				req.Header.Set("traceparent", tt.traceparent)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestHTTPMiddleware_SpanAttributes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tracer := noop.NewTracerProvider().Tracer("test")

	router := gin.New()
	router.Use(HTTPMiddleware(tracer))
	router.GET("/api/v1/devices/:id", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"id": c.Param("id")})
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/devices/123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHTTPMiddleware_ResponseHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Set up a global tracer provider
	otel.SetTracerProvider(noop.NewTracerProvider())
	tracer := otel.Tracer("test")

	router := gin.New()
	router.Use(HTTPMiddleware(tracer))
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestExtractTraceID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tracer := noop.NewTracerProvider().Tracer("test")

	router := gin.New()
	router.Use(HTTPMiddleware(tracer))
	router.GET("/test", func(c *gin.Context) {
		_ = GetTraceID(c.Request.Context())
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	// With noop tracer, trace ID will be empty/invalid
	// This test validates the function doesn't panic
}
