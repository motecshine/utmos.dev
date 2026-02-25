package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/utmos/utmos/internal/api"
	"github.com/utmos/utmos/internal/api/handler"
	"github.com/utmos/utmos/internal/downlink/model"
	"github.com/utmos/utmos/pkg/models"
)

func setupAPITestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&models.Device{}, &model.ServiceCall{})
	require.NoError(t, err)

	return db
}

// TestAPIIntegration tests the API service integration
func TestAPIIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db := setupAPITestDB(t)
	config := &api.Config{
		EnableAuth:  false,
		EnableTrace: false,
	}
	router := api.NewRouter(config, db, nil, nil, nil)

	t.Run("health check", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/health", nil)
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Equal(t, "healthy", resp["status"])
	})

	t.Run("ready check", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/ready", nil)
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

// TestDeviceAPIIntegration tests device API endpoints
func TestDeviceAPIIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db := setupAPITestDB(t)
	config := &api.Config{
		EnableAuth:  false,
		EnableTrace: false,
	}
	router := api.NewRouter(config, db, nil, nil, nil)

	var createdDeviceID uint

	t.Run("create device", func(t *testing.T) {
		req := handler.CreateDeviceRequest{
			DeviceSN:   "DEVICE001",
			DeviceName: "Test Drone",
			DeviceType: "drone",
			Vendor:     "dji",
		}
		body, _ := json.Marshal(req)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api/v1/devices", bytes.NewReader(body))
		r.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusCreated, w.Code)

		var resp handler.DeviceResponse
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, "DEVICE001", resp.DeviceSN)
		assert.Equal(t, "Test Drone", resp.DeviceName)
		assert.Equal(t, "dji", resp.Vendor)
		createdDeviceID = resp.ID
	})

	t.Run("get device by ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/api/v1/devices/1", nil)
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp handler.DeviceResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Equal(t, createdDeviceID, resp.ID)
		assert.Equal(t, "DEVICE001", resp.DeviceSN)
	})

	t.Run("get device by serial number", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/api/v1/devices/sn/DEVICE001", nil)
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp handler.DeviceResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Equal(t, "DEVICE001", resp.DeviceSN)
	})

	t.Run("list devices", func(t *testing.T) {
		// Create more devices
		for i := 2; i <= 5; i++ {
			req := handler.CreateDeviceRequest{
				DeviceSN:   fmt.Sprintf("DEVICE00%d", i),
				DeviceName: fmt.Sprintf("Test Device %d", i),
				DeviceType: "drone",
				Vendor:     "dji",
			}
			body, _ := json.Marshal(req)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/api/v1/devices", bytes.NewReader(body))
			r.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, r)
		}

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/api/v1/devices", nil)
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp handler.ListDevicesResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Equal(t, int64(5), resp.Total)
	})

	t.Run("list devices with pagination", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/api/v1/devices?page=1&page_size=2", nil)
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp handler.ListDevicesResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Len(t, resp.Devices, 2)
		assert.Equal(t, 1, resp.Page)
		assert.Equal(t, 2, resp.PageSize)
	})

	t.Run("list devices with filter", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/api/v1/devices?vendor=dji", nil)
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp handler.ListDevicesResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Equal(t, int64(5), resp.Total)
	})

	t.Run("update device", func(t *testing.T) {
		newName := "Updated Drone"
		req := handler.UpdateDeviceRequest{
			DeviceName: &newName,
		}
		body, _ := json.Marshal(req)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("PUT", "/api/v1/devices/1", bytes.NewReader(body))
		r.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp handler.DeviceResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Equal(t, "Updated Drone", resp.DeviceName)
	})

	t.Run("delete device", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("DELETE", "/api/v1/devices/1", nil)
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusNoContent, w.Code)

		// Verify deletion
		w = httptest.NewRecorder()
		r = httptest.NewRequest("GET", "/api/v1/devices/1", nil)
		router.ServeHTTP(w, r)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

// TestServiceCallAPIIntegration tests service call API endpoints
func TestServiceCallAPIIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db := setupAPITestDB(t)
	config := &api.Config{
		EnableAuth:  false,
		EnableTrace: false,
	}
	router := api.NewRouter(config, db, nil, nil, nil)

	t.Run("create service call", func(t *testing.T) {
		req := handler.ServiceCallRequest{
			DeviceSN: "DEVICE001",
			Vendor:   "dji",
			Method:   "takeoff",
			Params:   map[string]interface{}{"height": 50.0},
		}
		body, _ := json.Marshal(req)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api/v1/services/call", bytes.NewReader(body))
		r.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusAccepted, w.Code)

		var resp handler.ServiceCallResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Equal(t, "DEVICE001", resp.DeviceSN)
		assert.Equal(t, "dji", resp.Vendor)
		assert.Equal(t, "takeoff", resp.Method)
	})

	t.Run("list service calls by device", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/api/v1/services/calls/device/DEVICE001", nil)
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp handler.ListServiceCallsResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, resp.Total, int64(1))
	})

	t.Run("DJI takeoff command", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api/v1/services/dji/DEVICE001/takeoff?height=100", nil)
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusAccepted, w.Code)

		var resp handler.ServiceCallResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Equal(t, "takeoff", resp.Method)
	})

	t.Run("DJI land command", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api/v1/services/dji/DEVICE001/land", nil)
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusAccepted, w.Code)

		var resp handler.ServiceCallResponse
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, "land", resp.Method)
	})

	t.Run("DJI return home command", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api/v1/services/dji/DEVICE001/return-home", nil)
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusAccepted, w.Code)

		var resp handler.ServiceCallResponse
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, "return_home", resp.Method)
	})
}

// TestAPIAuthIntegration tests API authentication
func TestAPIAuthIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db := setupAPITestDB(t)
	config := &api.Config{
		EnableAuth:  true,
		EnableTrace: false,
		APIKeys:     []string{"valid-api-key"},
	}
	router := api.NewRouter(config, db, nil, nil, nil)

	t.Run("request without API key", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/api/v1/devices", nil)
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("request with invalid API key", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/api/v1/devices", nil)
		r.Header.Set("X-API-Key", "invalid-key")
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("request with valid API key", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/api/v1/devices", nil)
		r.Header.Set("X-API-Key", "valid-api-key")
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("request with Bearer token", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/api/v1/devices", nil)
		r.Header.Set("Authorization", "Bearer valid-api-key")
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

// TestAPIErrorHandling tests API error handling
func TestAPIErrorHandling(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db := setupAPITestDB(t)
	config := &api.Config{
		EnableAuth:  false,
		EnableTrace: false,
	}
	router := api.NewRouter(config, db, nil, nil, nil)

	t.Run("device not found", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/api/v1/devices/999", nil)
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var resp handler.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Equal(t, "DEVICE_NOT_FOUND", resp.Code)
	})

	t.Run("invalid device ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/api/v1/devices/invalid", nil)
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("invalid request body", func(t *testing.T) {
		body := []byte(`{"invalid": json}`)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api/v1/devices", bytes.NewReader(body))
		r.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("missing required fields", func(t *testing.T) {
		req := map[string]string{
			"device_name": "Test Device",
		}
		body, _ := json.Marshal(req)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api/v1/devices", bytes.NewReader(body))
		r.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
