package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTraceTestRouter(middleware *TraceMiddleware) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.Handler())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	router.GET("/api/v1/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"trace_id":   GetTraceID(c),
			"request_id": GetRequestID(c),
		})
	})

	return router
}

func TestDefaultTraceConfig(t *testing.T) {
	config := DefaultTraceConfig()

	assert.Equal(t, "iot-api", config.ServiceName)
	assert.Equal(t, "iot-api", config.TracerName)
	assert.Contains(t, config.SkipPaths, "/health")
}

func TestNewTraceMiddleware(t *testing.T) {
	config := &TraceConfig{
		ServiceName: "test-service",
		TracerName:  "test-tracer",
	}
	middleware := NewTraceMiddleware(config, nil)

	require.NotNil(t, middleware)
	assert.Equal(t, "test-service", middleware.config.ServiceName)
}

func TestTraceMiddleware_SkipPaths(t *testing.T) {
	config := &TraceConfig{
		SkipPaths: []string{"/health"},
	}
	middleware := NewTraceMiddleware(config, nil)
	router := setupTraceTestRouter(middleware)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/health", nil)
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	// Skip paths don't get trace headers
}

func TestTraceMiddleware_AddsHeaders(t *testing.T) {
	middleware := NewTraceMiddleware(nil, nil)
	router := setupTraceTestRouter(middleware)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/api/v1/test", nil)
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, w.Header().Get(TraceIDHeader))
	assert.NotEmpty(t, w.Header().Get(RequestIDHeader))
}

func TestTraceMiddleware_PreservesRequestID(t *testing.T) {
	middleware := NewTraceMiddleware(nil, nil)
	router := setupTraceTestRouter(middleware)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/api/v1/test", nil)
	r.Header.Set(RequestIDHeader, "custom-request-id")
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "custom-request-id", w.Header().Get(RequestIDHeader))
}

func TestRequestLogger(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(RequestLogger(nil))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, w.Header().Get(RequestIDHeader))
}

func TestRequestLogger_PreservesRequestID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(RequestLogger(nil))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"request_id": GetRequestID(c)})
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/test", nil)
	r.Header.Set(RequestIDHeader, "my-request-id")
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "my-request-id")
}

func TestInjectTraceContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		c.Set("trace_id", "trace-123")
		c.Set("request_id", "request-456")

		headers := InjectTraceContext(c)

		assert.Equal(t, "trace-123", headers[TraceIDHeader])
		assert.Equal(t, "request-456", headers[RequestIDHeader])

		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetTraceID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("with trace ID", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Set("trace_id", "trace-123")

		traceID := GetTraceID(c)
		assert.Equal(t, "trace-123", traceID)
	})

	t.Run("without trace ID", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())

		traceID := GetTraceID(c)
		assert.Empty(t, traceID)
	})
}

func TestGetRequestID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("with request ID", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Set("request_id", "request-456")

		requestID := GetRequestID(c)
		assert.Equal(t, "request-456", requestID)
	})

	t.Run("without request ID", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())

		requestID := GetRequestID(c)
		assert.Empty(t, requestID)
	})
}
