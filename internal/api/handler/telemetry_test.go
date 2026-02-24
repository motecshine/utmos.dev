package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTelemetryTestRouter(handler *Telemetry) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.GET("/api/v1/telemetry/:device_sn", handler.Query)
	router.GET("/api/v1/telemetry/:device_sn/latest", handler.Latest)
	router.GET("/api/v1/telemetry/:device_sn/aggregate", handler.Aggregate)

	return router
}

func TestNewTelemetry(t *testing.T) {
	config := &TelemetryConfig{
		Org:    "test-org",
		Bucket: "test-bucket",
	}
	handler := NewTelemetry(config, nil)

	require.NotNil(t, handler)
	assert.Equal(t, "test-org", handler.org)
	assert.Equal(t, "test-bucket", handler.bucket)
}

func TestTelemetry_Query_NoService(t *testing.T) {
	config := &TelemetryConfig{
		Org:    "test-org",
		Bucket: "test-bucket",
	}
	handler := NewTelemetry(config, nil)
	router := setupTelemetryTestRouter(handler)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/api/v1/telemetry/DEVICE001", nil)
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
}

func TestTelemetryHandler_Query_InvalidDeviceSN(t *testing.T) {
	config := &TelemetryConfig{
		Org:    "test-org",
		Bucket: "test-bucket",
	}
	handler := NewTelemetry(config, nil)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/api/v1/telemetry/:device_sn", func(c *gin.Context) {
		c.Params = gin.Params{{Key: "device_sn", Value: ""}}
		handler.Query(c)
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/api/v1/telemetry/", nil)
	router.ServeHTTP(w, r)

	// Will get 404 because route doesn't match
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestTelemetryHandler_Latest_NoService(t *testing.T) {
	config := &TelemetryConfig{
		Org:    "test-org",
		Bucket: "test-bucket",
	}
	handler := NewTelemetry(config, nil)
	router := setupTelemetryTestRouter(handler)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/api/v1/telemetry/DEVICE001/latest", nil)
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
}

func TestTelemetryHandler_Aggregate_NoService(t *testing.T) {
	config := &TelemetryConfig{
		Org:    "test-org",
		Bucket: "test-bucket",
	}
	handler := NewTelemetry(config, nil)
	router := setupTelemetryTestRouter(handler)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/api/v1/telemetry/DEVICE001/aggregate?field=temperature", nil)
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
}

func TestTelemetryHandler_Aggregate_MissingField(t *testing.T) {
	config := &TelemetryConfig{
		Org:    "test-org",
		Bucket: "test-bucket",
	}
	handler := NewTelemetry(config, nil)
	router := setupTelemetryTestRouter(handler)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/api/v1/telemetry/DEVICE001/aggregate", nil)
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestTelemetryHandler_Aggregate_InvalidFunction(t *testing.T) {
	config := &TelemetryConfig{
		URL:    "http://localhost:8086",
		Token:  "test-token",
		Org:    "test-org",
		Bucket: "test-bucket",
	}
	handler := NewTelemetry(config, nil)
	router := setupTelemetryTestRouter(handler)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/api/v1/telemetry/DEVICE001/aggregate?field=temp&fn=invalid", nil)
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestTelemetryHandler_BuildQuery(t *testing.T) {
	config := &TelemetryConfig{
		Org:    "test-org",
		Bucket: "test-bucket",
	}
	handler := NewTelemetry(config, nil)

	t.Run("basic query", func(t *testing.T) {
		query := handler.buildQuery("DEVICE001", "-1h", "now()", "", 100)
		assert.Contains(t, query, "test-bucket")
		assert.Contains(t, query, "DEVICE001")
		assert.Contains(t, query, "limit(n: 100)")
	})

	t.Run("with measurement", func(t *testing.T) {
		query := handler.buildQuery("DEVICE001", "-1h", "now()", "temperature", 50)
		assert.Contains(t, query, "temperature")
		assert.Contains(t, query, "limit(n: 50)")
	})
}

func TestTelemetryPoint(t *testing.T) {
	point := TelemetryPoint{
		Time: "2024-01-01T00:00:00Z",
		Fields: map[string]any{
			"temperature": 25.5,
			"humidity":    60.0,
		},
		Tags: map[string]string{
			"device_sn": "DEVICE001",
			"vendor":    "dji",
		},
	}

	assert.Equal(t, "2024-01-01T00:00:00Z", point.Time)
	assert.Equal(t, 25.5, point.Fields["temperature"])
	assert.Equal(t, "DEVICE001", point.Tags["device_sn"])
}

func TestTelemetryQueryResponse(t *testing.T) {
	resp := TelemetryQueryResponse{
		DeviceSN: "DEVICE001",
		Points: []TelemetryPoint{
			{Time: "2024-01-01T00:00:00Z", Fields: map[string]any{"temp": 25.0}},
			{Time: "2024-01-01T01:00:00Z", Fields: map[string]any{"temp": 26.0}},
		},
		Total: 2,
	}

	assert.Equal(t, "DEVICE001", resp.DeviceSN)
	assert.Len(t, resp.Points, 2)
	assert.Equal(t, 2, resp.Total)
}

func TestLatestTelemetryResponse(t *testing.T) {
	resp := LatestTelemetryResponse{
		DeviceSN:  "DEVICE001",
		Timestamp: "2024-01-01T00:00:00Z",
		Data: map[string]any{
			"temperature": 25.5,
			"humidity":    60.0,
		},
	}

	assert.Equal(t, "DEVICE001", resp.DeviceSN)
	assert.Equal(t, 25.5, resp.Data["temperature"])
}
