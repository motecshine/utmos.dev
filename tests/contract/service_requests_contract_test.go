package contract

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

	"github.com/utmos/utmos/internal/api/handler"
	"github.com/utmos/utmos/internal/downlink/model"
)

func setupContractDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&model.ServiceCall{}))
	return db
}

func setupContractRouter(t *testing.T, db *gorm.DB) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	svc := handler.NewService(db, nil, nil)

	r.POST("/api/v1/services/call", svc.Call)
	r.GET("/api/v1/services/calls/:id", svc.Get)
	r.GET("/api/v1/services/calls/device/:device_sn", svc.ListByDevice)
	return r
}

// TestServiceRequestContract validates the POST /api/v1/services/call
// endpoint against the contract: required fields, status code, and
// response shape.
func TestServiceRequestContract(t *testing.T) {
	db := setupContractDB(t)
	router := setupContractRouter(t, db)

	t.Run("valid request returns 202 with required response fields", func(t *testing.T) {
		body, _ := json.Marshal(map[string]any{
			"device_sn": "CONTRACT-DEV-001",
			"vendor":    "dji",
			"method":    "takeoff",
			"params":    map[string]any{"height": 50.0},
			"call_type": "command",
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/v1/services/call", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusAccepted, w.Code)

		var resp map[string]any
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))

		// Contract: response MUST include these fields
		assert.Contains(t, resp, "id", "response must include id field")
		assert.Equal(t, "CONTRACT-DEV-001", resp["device_sn"])
		assert.Equal(t, "dji", resp["vendor"])
		assert.Equal(t, "takeoff", resp["method"])
		assert.Equal(t, "command", resp["call_type"])
		assert.Contains(t, []string{"pending", "sent"}, resp["status"])
		assert.NotEmpty(t, resp["created_at"], "response must include created_at")
	})

	t.Run("default call_type is command", func(t *testing.T) {
		body, _ := json.Marshal(map[string]any{
			"device_sn": "CONTRACT-DEV-002",
			"vendor":    "dji",
			"method":    "land",
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/v1/services/call", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusAccepted, w.Code)

		var resp map[string]any
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
		assert.Equal(t, "command", resp["call_type"])
	})

	t.Run("missing required field returns 400", func(t *testing.T) {
		// Missing vendor and method
		body, _ := json.Marshal(map[string]any{
			"device_sn": "CONTRACT-DEV-001",
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/v1/services/call", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var errResp map[string]any
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &errResp))
		assert.NotEmpty(t, errResp["code"])
		assert.NotEmpty(t, errResp["message"])
	})

	t.Run("empty body returns 400", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/v1/services/call", bytes.NewReader([]byte(`{}`)))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("lifecycle status values conform to contract enum", func(t *testing.T) {
		validStatuses := map[string]bool{
			"pending":  true,
			"sent":     true,
			"success":  true,
			"failed":   true,
			"timeout":  true,
			"retrying": true,
		}

		// Create a call and verify its status is in the valid set
		body, _ := json.Marshal(map[string]any{
			"device_sn": "CONTRACT-DEV-003",
			"vendor":    "dji",
			"method":    "hover",
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/v1/services/call", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		var resp map[string]any
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
		status, ok := resp["status"].(string)
		require.True(t, ok, "status must be a string")
		assert.True(t, validStatuses[status], "status %q not in contract enum", status)
	})
}

// TestServiceRequestStatusContract validates the GET /api/v1/services/calls/:id
// endpoint against the contract: 200 with full response for existing
// requests, 404 with error body for missing requests.
func TestServiceRequestStatusContract(t *testing.T) {
	db := setupContractDB(t)
	router := setupContractRouter(t, db)

	// Seed a service call directly in the database
	call := &model.ServiceCall{
		ID:       "contract-call-001",
		DeviceSN: "CONTRACT-DEV-001",
		Vendor:   "dji",
		Method:   "takeoff",
		CallType: model.ServiceCallTypeCommand,
		Status:   model.ServiceCallStatusSuccess,
		TID:      "tid-contract-001",
		BID:      "bid-contract-001",
	}
	require.NoError(t, db.Create(call).Error)

	t.Run("existing request returns 200 with required fields", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/v1/services/calls/contract-call-001", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]any
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))

		// Contract: response MUST include these fields
		assert.Equal(t, "contract-call-001", resp["id"])
		assert.Equal(t, "CONTRACT-DEV-001", resp["device_sn"])
		assert.Equal(t, "dji", resp["vendor"])
		assert.Equal(t, "takeoff", resp["method"])
		assert.Equal(t, "command", resp["call_type"])
		assert.Equal(t, "success", resp["status"])
		assert.Equal(t, "tid-contract-001", resp["tid"])
		assert.Equal(t, "bid-contract-001", resp["bid"])
		assert.NotNil(t, resp["created_at"])
	})

	t.Run("non-existent request returns 404 with error body", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/v1/services/calls/does-not-exist", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var errResp map[string]any
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &errResp))
		assert.NotEmpty(t, errResp["code"])
		assert.NotEmpty(t, errResp["message"])
	})

	t.Run("status field values match lifecycle enum", func(t *testing.T) {
		// Seed calls in various lifecycle states
		states := []model.ServiceCallStatus{
			model.ServiceCallStatusPending,
			model.ServiceCallStatusSent,
			model.ServiceCallStatusFailed,
			model.ServiceCallStatusTimeout,
			model.ServiceCallStatusRetrying,
		}

		for i, status := range states {
			id := "lifecycle-" + string(status)
			seedCall := &model.ServiceCall{
				ID:       id,
				DeviceSN: "CONTRACT-DEV-001",
				Vendor:   "dji",
				Method:   "test",
				Status:   status,
			}
			require.NoError(t, db.Create(seedCall).Error, "seed call %d", i)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/api/v1/services/calls/"+id, nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			var resp map[string]any
			require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
			assert.Equal(t, string(status), resp["status"],
				"status mismatch for seeded call with status %s", status)
		}
	})

	t.Run("list by device returns matching calls", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/v1/services/calls/device/CONTRACT-DEV-001", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]any
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))

		total, ok := resp["total"].(float64)
		require.True(t, ok)
		assert.Greater(t, total, float64(0), "should return seeded calls")

		calls, ok := resp["service_calls"].([]any)
		require.True(t, ok)
		assert.NotEmpty(t, calls)
	})
}
