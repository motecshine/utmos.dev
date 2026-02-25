package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/utmos/utmos/internal/api/handler"
	"github.com/utmos/utmos/internal/downlink/model"
	"github.com/utmos/utmos/pkg/models"
)

func setupRouterTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&models.Device{}, &model.ServiceCall{})
	require.NoError(t, err)

	return db
}

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	assert.True(t, config.EnableAuth)
	assert.True(t, config.EnableTrace)
	assert.Equal(t, "iot-api", config.ServiceName)
}

func TestNewRouter(t *testing.T) {
	db := setupRouterTestDB(t)
	config := &Config{
		EnableAuth:  false,
		EnableTrace: false,
	}

	router := NewRouter(config, db, nil, nil, nil)

	require.NotNil(t, router)
	assert.NotNil(t, router.engine)
	assert.NotNil(t, router.deviceHandler)
	assert.NotNil(t, router.serviceHandler)
}

func TestRouter_HealthEndpoints(t *testing.T) {
	db := setupRouterTestDB(t)
	config := &Config{
		EnableAuth:  false,
		EnableTrace: false,
	}
	router := NewRouter(config, db, nil, nil, nil)

	t.Run("health endpoint", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/health", nil)
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Equal(t, "healthy", resp["status"])
	})

	t.Run("ready endpoint", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/ready", nil)
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]any
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, "ready", resp["status"])
		assert.NotNil(t, resp["checks"])

		checks := resp["checks"].(map[string]any)
		assert.Equal(t, "ok", checks["database"])
	})
}

func TestRouter_DeviceEndpoints(t *testing.T) {
	db := setupRouterTestDB(t)
	config := &Config{
		EnableAuth:  false,
		EnableTrace: false,
	}
	router := NewRouter(config, db, nil, nil, nil)

	t.Run("list devices", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/api/v1/devices", nil)
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("get device not found", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/api/v1/devices/999", nil)
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestRouter_ServiceEndpoints(t *testing.T) {
	db := setupRouterTestDB(t)
	config := &Config{
		EnableAuth:  false,
		EnableTrace: false,
	}
	router := NewRouter(config, db, nil, nil, nil)

	t.Run("get service call not found", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/api/v1/services/calls/nonexistent", nil)
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("list service calls by device", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/api/v1/services/calls/device/DEVICE001", nil)
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestRouter_WithAuth(t *testing.T) {
	db := setupRouterTestDB(t)
	config := &Config{
		EnableAuth:  true,
		EnableTrace: false,
		APIKeys:     []string{"valid-api-key"},
	}
	router := NewRouter(config, db, nil, nil, nil)

	t.Run("without API key", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/api/v1/devices", nil)
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("with valid API key", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/api/v1/devices", nil)
		r.Header.Set("X-API-Key", "valid-api-key")
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("health endpoint bypasses auth", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/health", nil)
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestRouter_WithTelemetry(t *testing.T) {
	db := setupRouterTestDB(t)
	config := &Config{
		EnableAuth:  false,
		EnableTrace: false,
		TelemetryConfig: &handler.TelemetryConfig{
			Org:    "test-org",
			Bucket: "test-bucket",
		},
	}
	router := NewRouter(config, db, nil, nil, nil)

	t.Run("telemetry endpoint exists", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/api/v1/telemetry/DEVICE001", nil)
		router.ServeHTTP(w, r)

		// Will return service unavailable because no InfluxDB connection
		assert.Equal(t, http.StatusServiceUnavailable, w.Code)
	})
}

func TestRouter_Engine(t *testing.T) {
	db := setupRouterTestDB(t)
	config := &Config{
		EnableAuth:  false,
		EnableTrace: false,
	}
	router := NewRouter(config, db, nil, nil, nil)

	engine := router.Engine()
	assert.NotNil(t, engine)
}

func TestRouter_Close(t *testing.T) {
	db := setupRouterTestDB(t)
	config := &Config{
		EnableAuth:  false,
		EnableTrace: false,
	}
	router := NewRouter(config, db, nil, nil, nil)

	// Should not panic
	router.Close()
}
