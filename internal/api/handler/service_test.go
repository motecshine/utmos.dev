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

	"github.com/utmos/utmos/internal/downlink/model"
)

func setupServiceTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&model.ServiceCall{})
	require.NoError(t, err)

	return db
}

func setupServiceTestRouter(handler *Service) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.POST("/api/v1/services/call", handler.Call)
	router.GET("/api/v1/services/calls/:id", handler.Get)
	router.GET("/api/v1/services/calls/device/:device_sn", handler.ListByDevice)
	router.POST("/api/v1/services/calls/:id/cancel", handler.Cancel)

	return router
}

func TestNewService(t *testing.T) {
	db := setupServiceTestDB(t)
	handler := NewService(db, nil, nil)

	require.NotNil(t, handler)
	assert.NotNil(t, handler.repository)
}

func TestService_Call(t *testing.T) {
	db := setupServiceTestDB(t)
	handler := NewService(db, nil, nil)
	router := setupServiceTestRouter(handler)

	t.Run("successful call", func(t *testing.T) {
		req := ServiceCallRequest{
			DeviceSN: "DEVICE001",
			Vendor:   "dji",
			Method:   "takeoff",
			Params:   map[string]any{"height": 50.0},
		}
		body, _ := json.Marshal(req)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api/v1/services/call", bytes.NewReader(body))
		r.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusAccepted, w.Code)

		var resp ServiceCallResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Equal(t, "DEVICE001", resp.DeviceSN)
		assert.Equal(t, "dji", resp.Vendor)
		assert.Equal(t, "takeoff", resp.Method)
	})

	t.Run("invalid request", func(t *testing.T) {
		body := []byte(`{"device_sn": "DEVICE001"}`)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api/v1/services/call", bytes.NewReader(body))
		r.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestService_Get(t *testing.T) {
	db := setupServiceTestDB(t)
	handler := NewService(db, nil, nil)
	router := setupServiceTestRouter(handler)

	// Create a service call first
	call := &model.ServiceCall{
		ID:       "call-001",
		DeviceSN: "DEVICE001",
		Vendor:   "dji",
		Method:   "takeoff",
		Status:   model.ServiceCallStatusPending,
	}
	db.Create(call)

	t.Run("get existing call", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/api/v1/services/calls/call-001", nil)
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp ServiceCallResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Equal(t, "call-001", resp.ID)
		assert.Equal(t, "DEVICE001", resp.DeviceSN)
	})

	t.Run("call not found", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/api/v1/services/calls/nonexistent", nil)
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestService_ListByDevice(t *testing.T) {
	db := setupServiceTestDB(t)
	handler := NewService(db, nil, nil)
	router := setupServiceTestRouter(handler)

	// Create service calls
	calls := []model.ServiceCall{
		{ID: "call-001", DeviceSN: "DEVICE001", Vendor: "dji", Method: "takeoff", Status: model.ServiceCallStatusSuccess},
		{ID: "call-002", DeviceSN: "DEVICE001", Vendor: "dji", Method: "land", Status: model.ServiceCallStatusPending},
		{ID: "call-003", DeviceSN: "DEVICE002", Vendor: "dji", Method: "takeoff", Status: model.ServiceCallStatusPending},
	}
	for _, c := range calls {
		db.Create(&c)
	}

	t.Run("list by device", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/api/v1/services/calls/device/DEVICE001", nil)
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp ListServiceCallsResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Equal(t, int64(2), resp.Total)
	})

	t.Run("with limit", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/api/v1/services/calls/device/DEVICE001?limit=1", nil)
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp ListServiceCallsResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Len(t, resp.ServiceCalls, 1)
	})
}

func TestService_Cancel(t *testing.T) {
	db := setupServiceTestDB(t)
	handler := NewService(db, nil, nil)
	router := setupServiceTestRouter(handler)

	t.Run("cancel pending call", func(t *testing.T) {
		call := &model.ServiceCall{
			ID:       "call-cancel-001",
			DeviceSN: "DEVICE001",
			Vendor:   "dji",
			Method:   "takeoff",
			Status:   model.ServiceCallStatusPending,
		}
		db.Create(call)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api/v1/services/calls/call-cancel-001/cancel", nil)
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp ServiceCallResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Equal(t, "cancelled", resp.Status)
	})

	t.Run("cannot cancel completed call", func(t *testing.T) {
		call := &model.ServiceCall{
			ID:       "call-cancel-002",
			DeviceSN: "DEVICE001",
			Vendor:   "dji",
			Method:   "takeoff",
			Status:   model.ServiceCallStatusSuccess,
		}
		db.Create(call)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api/v1/services/calls/call-cancel-002/cancel", nil)
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("call not found", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api/v1/services/calls/nonexistent/cancel", nil)
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestToDispatcherServiceCall(t *testing.T) {
	req := &ServiceCallRequest{
		DeviceSN:   "DEVICE001",
		Vendor:     "dji",
		Method:     "takeoff",
		Params:     map[string]any{"height": 50.0},
		CallType:   "command",
		MaxRetries: 5,
	}

	call := toDispatcherServiceCall(req)

	assert.Equal(t, "DEVICE001", call.DeviceSN)
	assert.Equal(t, "dji", call.Vendor)
	assert.Equal(t, "takeoff", call.Method)
	assert.Equal(t, 5, call.MaxRetries)
}

func TestToServiceCallResponse(t *testing.T) {
	call := &model.ServiceCall{
		ID:       "call-001",
		DeviceSN: "DEVICE001",
		Vendor:   "dji",
		Method:   "takeoff",
		CallType: model.ServiceCallTypeCommand,
		Status:   model.ServiceCallStatusSuccess,
		TID:      "tid-001",
	}

	resp := toServiceCallResponse(call)

	assert.Equal(t, "call-001", resp.ID)
	assert.Equal(t, "DEVICE001", resp.DeviceSN)
	assert.Equal(t, "dji", resp.Vendor)
	assert.Equal(t, "takeoff", resp.Method)
	assert.Equal(t, "command", resp.CallType)
	assert.Equal(t, "success", resp.Status)
}
