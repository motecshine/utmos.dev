package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/utmos/utmos/pkg/models"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&models.Device{})
	require.NoError(t, err)

	return db
}

func setupTestRouter(handler *Device) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.POST("/api/v1/devices", handler.Create)
	router.GET("/api/v1/devices", handler.List)
	router.GET("/api/v1/devices/:id", handler.Get)
	router.GET("/api/v1/devices/sn/:sn", handler.GetBySN)
	router.PUT("/api/v1/devices/:id", handler.Update)
	router.DELETE("/api/v1/devices/:id", handler.Delete)

	return router
}

func TestNewDevice(t *testing.T) {
	db := setupTestDB(t)
	handler := NewDevice(db, nil)

	require.NotNil(t, handler)
	assert.Equal(t, db, handler.db)
}

func TestDevice_Create(t *testing.T) {
	db := setupTestDB(t)
	handler := NewDevice(db, nil)
	router := setupTestRouter(handler)

	t.Run("successful creation", func(t *testing.T) {
		req := CreateDeviceRequest{
			DeviceSN:   "DEVICE001",
			DeviceName: "Test Device",
			DeviceType: "drone",
			Vendor:     "dji",
		}
		body, _ := json.Marshal(req)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api/v1/devices", bytes.NewReader(body))
		r.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusCreated, w.Code)

		var resp DeviceResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Equal(t, "DEVICE001", resp.DeviceSN)
		assert.Equal(t, "Test Device", resp.DeviceName)
		assert.Equal(t, "dji", resp.Vendor)
	})

	t.Run("duplicate device", func(t *testing.T) {
		req := CreateDeviceRequest{
			DeviceSN:   "DEVICE001",
			DeviceName: "Duplicate Device",
			DeviceType: "drone",
		}
		body, _ := json.Marshal(req)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api/v1/devices", bytes.NewReader(body))
		r.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusConflict, w.Code)
	})

	t.Run("invalid request", func(t *testing.T) {
		body := []byte(`{"device_name": "Missing SN"}`)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api/v1/devices", bytes.NewReader(body))
		r.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("default vendor", func(t *testing.T) {
		req := CreateDeviceRequest{
			DeviceSN:   "DEVICE002",
			DeviceName: "Generic Device",
			DeviceType: "sensor",
		}
		body, _ := json.Marshal(req)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api/v1/devices", bytes.NewReader(body))
		r.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusCreated, w.Code)

		var resp DeviceResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Equal(t, "generic", resp.Vendor)
	})
}

func TestDevice_Get(t *testing.T) {
	db := setupTestDB(t)
	handler := NewDevice(db, nil)
	router := setupTestRouter(handler)

	// Create a device first
	device := &models.Device{
		DeviceSN:   "DEVICE001",
		DeviceName: "Test Device",
		DeviceType: "drone",
		Vendor:     "dji",
		Status:     models.DeviceStatusOnline,
	}
	db.Create(device)

	t.Run("get existing device", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/api/v1/devices/1", nil)
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp DeviceResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Equal(t, "DEVICE001", resp.DeviceSN)
	})

	t.Run("device not found", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/api/v1/devices/999", nil)
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("invalid id", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/api/v1/devices/invalid", nil)
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestDevice_GetBySN(t *testing.T) {
	db := setupTestDB(t)
	handler := NewDevice(db, nil)
	router := setupTestRouter(handler)

	// Create a device first
	device := &models.Device{
		DeviceSN:   "DEVICE001",
		DeviceName: "Test Device",
		DeviceType: "drone",
		Vendor:     "dji",
	}
	db.Create(device)

	t.Run("get by serial number", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/api/v1/devices/sn/DEVICE001", nil)
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp DeviceResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Equal(t, "DEVICE001", resp.DeviceSN)
	})

	t.Run("device not found", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/api/v1/devices/sn/NONEXISTENT", nil)
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestDevice_List(t *testing.T) {
	db := setupTestDB(t)
	handler := NewDevice(db, nil)
	router := setupTestRouter(handler)

	// Create multiple devices
	devices := []models.Device{
		{DeviceSN: "DEVICE001", DeviceName: "Device 1", DeviceType: "drone", Vendor: "dji", Status: models.DeviceStatusOnline},
		{DeviceSN: "DEVICE002", DeviceName: "Device 2", DeviceType: "drone", Vendor: "dji", Status: models.DeviceStatusOffline},
		{DeviceSN: "DEVICE003", DeviceName: "Device 3", DeviceType: "sensor", Vendor: "generic", Status: models.DeviceStatusOnline},
	}
	for _, d := range devices {
		db.Create(&d)
	}

	t.Run("list all devices", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/api/v1/devices", nil)
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp ListDevicesResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Equal(t, int64(3), resp.Total)
		assert.Len(t, resp.Devices, 3)
	})

	t.Run("filter by vendor", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/api/v1/devices?vendor=dji", nil)
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp ListDevicesResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Equal(t, int64(2), resp.Total)
	})

	t.Run("filter by status", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/api/v1/devices?status=online", nil)
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp ListDevicesResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Equal(t, int64(2), resp.Total)
	})

	t.Run("pagination", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/api/v1/devices?page=1&page_size=2", nil)
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp ListDevicesResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Len(t, resp.Devices, 2)
		assert.Equal(t, 1, resp.Page)
		assert.Equal(t, 2, resp.PageSize)
		assert.Equal(t, 2, resp.TotalPages)
	})
}

func TestDevice_Update(t *testing.T) {
	db := setupTestDB(t)
	handler := NewDevice(db, nil)
	router := setupTestRouter(handler)

	// Create a device first
	device := &models.Device{
		DeviceSN:   "DEVICE001",
		DeviceName: "Test Device",
		DeviceType: "drone",
		Vendor:     "dji",
	}
	db.Create(device)

	t.Run("update device", func(t *testing.T) {
		newName := "Updated Device"
		req := UpdateDeviceRequest{
			DeviceName: &newName,
		}
		body, _ := json.Marshal(req)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("PUT", "/api/v1/devices/1", bytes.NewReader(body))
		r.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp DeviceResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Equal(t, "Updated Device", resp.DeviceName)
	})

	t.Run("update status", func(t *testing.T) {
		status := models.DeviceStatusOnline
		req := UpdateDeviceRequest{
			Status: &status,
		}
		body, _ := json.Marshal(req)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("PUT", "/api/v1/devices/1", bytes.NewReader(body))
		r.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp DeviceResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Equal(t, models.DeviceStatusOnline, resp.Status)
	})

	t.Run("device not found", func(t *testing.T) {
		newName := "Updated"
		req := UpdateDeviceRequest{DeviceName: &newName}
		body, _ := json.Marshal(req)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("PUT", "/api/v1/devices/999", bytes.NewReader(body))
		r.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestDevice_Delete(t *testing.T) {
	db := setupTestDB(t)
	handler := NewDevice(db, nil)
	router := setupTestRouter(handler)

	// Create a device first
	device := &models.Device{
		DeviceSN:   "DEVICE001",
		DeviceName: "Test Device",
		DeviceType: "drone",
		Vendor:     "dji",
	}
	db.Create(device)

	t.Run("delete device", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("DELETE", "/api/v1/devices/1", nil)
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusNoContent, w.Code)

		// Verify deletion
		var count int64
		db.Model(&models.Device{}).Count(&count)
		assert.Equal(t, int64(0), count)
	})

	t.Run("device not found", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("DELETE", "/api/v1/devices/999", nil)
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestToDeviceResponse(t *testing.T) {
	device := &models.Device{
		ID:         1,
		DeviceSN:   "DEVICE001",
		DeviceName: "Test Device",
		DeviceType: "drone",
		Vendor:     "dji",
		Status:     models.DeviceStatusOnline,
	}

	resp := toDeviceResponse(device)

	assert.Equal(t, uint(1), resp.ID)
	assert.Equal(t, "DEVICE001", resp.DeviceSN)
	assert.Equal(t, "Test Device", resp.DeviceName)
	assert.Equal(t, "drone", resp.DeviceType)
	assert.Equal(t, "dji", resp.Vendor)
	assert.Equal(t, models.DeviceStatusOnline, resp.Status)
}
