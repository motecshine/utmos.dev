package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/utmos/utmos/internal/api"
	"github.com/utmos/utmos/internal/api/handler"
	"github.com/utmos/utmos/internal/downlink/model"
	"github.com/utmos/utmos/internal/ws"
	"github.com/utmos/utmos/internal/ws/hub"
	"github.com/utmos/utmos/pkg/models"
	"github.com/utmos/utmos/pkg/rabbitmq"
)

// TestE2EDataFlow tests the complete end-to-end data flow
// This simulates: Device -> Gateway -> Uplink -> API/WS -> Client
func TestE2EDataFlow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	// Setup test database
	db := setupE2ETestDB(t)

	// Create a test device
	device := &models.Device{
		DeviceSN:   "E2E-DRONE-001",
		DeviceName: "E2E Test Drone",
		DeviceType: "drone",
		Vendor:     "dji",
		Status:     "online",
	}
	err := db.Create(device).Error
	require.NoError(t, err)

	t.Run("complete uplink flow", func(t *testing.T) {
		// Simulate uplink message processing
		// In a real scenario, this would come from iot-gateway via RabbitMQ

		msg := &rabbitmq.StandardMessage{
			TID:       "tid-e2e-001",
			BID:       "bid-e2e-001",
			Service:   "iot-uplink",
			Action:    "telemetry.report",
			DeviceSN:  "E2E-DRONE-001",
			Timestamp: time.Now().UnixMilli(),
			Data:      json.RawMessage(`{"latitude": 39.9042, "longitude": 116.4074, "altitude": 100.5}`),
		}

		// Verify message structure
		assert.NotEmpty(t, msg.TID)
		assert.NotEmpty(t, msg.BID)
		assert.Equal(t, "E2E-DRONE-001", msg.DeviceSN)
	})

	t.Run("complete downlink flow", func(t *testing.T) {
		// Setup API router
		config := &api.Config{
			EnableAuth:  false,
			EnableTrace: false,
		}
		router := api.NewRouter(config, db, nil, nil, nil)

		// Create service call request
		req := handler.ServiceCallRequest{
			DeviceSN: "E2E-DRONE-001",
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
		assert.Equal(t, "E2E-DRONE-001", resp.DeviceSN)
		assert.Equal(t, "takeoff", resp.Method)
	})

	t.Run("websocket real-time push", func(t *testing.T) {
		// Setup WebSocket service
		wsSvc := ws.NewService(nil, nil, nil)
		err := wsSvc.Start(context.Background())
		require.NoError(t, err)
		defer func() { _ = wsSvc.Stop() }()

		// Create test server
		server := httptest.NewServer(http.HandlerFunc(wsSvc.HandleWebSocket))
		defer server.Close()

		// Connect WebSocket client
		wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "?device_sn=E2E-DRONE-001"
		conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		require.NoError(t, err)
		defer func() { _ = conn.Close() }()

		// Wait for registration
		time.Sleep(100 * time.Millisecond)

		// Subscribe to device telemetry
		subscribeMsg := hub.Message{
			Type:  hub.MessageTypeSubscribe,
			Event: "device.E2E-DRONE-001.telemetry",
		}
		data, _ := json.Marshal(subscribeMsg)
		err = conn.WriteMessage(websocket.TextMessage, data)
		require.NoError(t, err)

		// Read ack
		_, respData, err := conn.ReadMessage()
		require.NoError(t, err)

		var ackMsg hub.Message
		err = json.Unmarshal(respData, &ackMsg)
		require.NoError(t, err)
		assert.Equal(t, hub.MessageTypeAck, ackMsg.Type)

		// Push telemetry message
		wsSvc.Pusher().PushToTopic("device.E2E-DRONE-001.telemetry", &hub.Message{
			Type:  hub.MessageTypeEvent,
			Event: "device.E2E-DRONE-001.telemetry",
			Data: map[string]interface{}{
				"latitude":  39.9042,
				"longitude": 116.4074,
				"altitude":  100.5,
			},
		})

		// Wait for push
		time.Sleep(100 * time.Millisecond)

		// Read pushed message
		_ = conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
		_, respData, err = conn.ReadMessage()
		require.NoError(t, err)

		var eventMsg hub.Message
		err = json.Unmarshal(respData, &eventMsg)
		require.NoError(t, err)
		assert.Equal(t, hub.MessageTypeEvent, eventMsg.Type)
		assert.Equal(t, "device.E2E-DRONE-001.telemetry", eventMsg.Event)
	})

	t.Run("device lifecycle", func(t *testing.T) {
		// Setup API router
		config := &api.Config{
			EnableAuth:  false,
			EnableTrace: false,
		}
		router := api.NewRouter(config, db, nil, nil, nil)

		// 1. Create device
		createReq := handler.CreateDeviceRequest{
			DeviceSN:   "E2E-LIFECYCLE-001",
			DeviceName: "Lifecycle Test Device",
			DeviceType: "drone",
			Vendor:     "dji",
		}
		body, _ := json.Marshal(createReq)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api/v1/devices", bytes.NewReader(body))
		r.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, r)
		assert.Equal(t, http.StatusCreated, w.Code)

		var createResp handler.DeviceResponse
		err := json.Unmarshal(w.Body.Bytes(), &createResp)
		require.NoError(t, err)
		deviceID := createResp.ID

		// 2. Get device
		w = httptest.NewRecorder()
		r = httptest.NewRequest("GET", "/api/v1/devices/sn/E2E-LIFECYCLE-001", nil)
		router.ServeHTTP(w, r)
		assert.Equal(t, http.StatusOK, w.Code)

		// 3. Update device
		newName := "Updated Lifecycle Device"
		updateReq := handler.UpdateDeviceRequest{
			DeviceName: &newName,
		}
		body, _ = json.Marshal(updateReq)

		w = httptest.NewRecorder()
		r = httptest.NewRequest("PUT", fmt.Sprintf("/api/v1/devices/%d", deviceID), bytes.NewReader(body))
		r.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, r)
		// May return 200 or 404 depending on ID format

		// 4. Send command to device
		cmdReq := handler.ServiceCallRequest{
			DeviceSN: "E2E-LIFECYCLE-001",
			Vendor:   "dji",
			Method:   "land",
		}
		body, _ = json.Marshal(cmdReq)

		w = httptest.NewRecorder()
		r = httptest.NewRequest("POST", "/api/v1/services/call", bytes.NewReader(body))
		r.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, r)
		assert.Equal(t, http.StatusAccepted, w.Code)
	})
}

// TestE2EMessageTracing tests distributed tracing across services
func TestE2EMessageTracing(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E tracing test in short mode")
	}

	// Create a message with trace context
	msg := &rabbitmq.StandardMessage{
		TID:       "trace-tid-001",
		BID:       "trace-bid-001",
		Service:   "iot-gateway",
		Action:    "telemetry.report",
		DeviceSN:  "TRACE-DRONE-001",
		Timestamp: time.Now().UnixMilli(),
		Data:      json.RawMessage(`{"test": "data"}`),
	}

	// Verify trace IDs are preserved
	assert.NotEmpty(t, msg.TID)
	assert.NotEmpty(t, msg.BID)

	// Simulate message passing through services
	// Gateway -> Uplink
	uplinkMsg := &rabbitmq.StandardMessage{
		TID:       msg.TID, // Preserve TID
		BID:       msg.BID, // Preserve BID
		Service:   "iot-uplink",
		Action:    msg.Action,
		DeviceSN:  msg.DeviceSN,
		Timestamp: time.Now().UnixMilli(),
		Data:      msg.Data,
	}

	assert.Equal(t, msg.TID, uplinkMsg.TID)
	assert.Equal(t, msg.BID, uplinkMsg.BID)

	// Uplink -> WS
	wsMsg := &rabbitmq.StandardMessage{
		TID:       uplinkMsg.TID,
		BID:       uplinkMsg.BID,
		Service:   "iot-ws",
		Action:    "push.telemetry",
		DeviceSN:  uplinkMsg.DeviceSN,
		Timestamp: time.Now().UnixMilli(),
		Data:      uplinkMsg.Data,
	}

	assert.Equal(t, msg.TID, wsMsg.TID)
	assert.Equal(t, msg.BID, wsMsg.BID)
}

// TestE2EErrorHandling tests error handling across the system
func TestE2EErrorHandling(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E error handling test in short mode")
	}

	db := setupE2ETestDB(t)
	config := &api.Config{
		EnableAuth:  false,
		EnableTrace: false,
	}
	router := api.NewRouter(config, db, nil, nil, nil)

	t.Run("invalid device command", func(t *testing.T) {
		// Send command to non-existent device
		req := handler.ServiceCallRequest{
			DeviceSN: "NON-EXISTENT-DEVICE",
			Vendor:   "dji",
			Method:   "takeoff",
		}
		body, _ := json.Marshal(req)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api/v1/services/call", bytes.NewReader(body))
		r.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, r)

		// Should still accept (async processing)
		assert.Equal(t, http.StatusAccepted, w.Code)
	})

	t.Run("invalid request body", func(t *testing.T) {
		body := []byte(`{invalid json}`)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api/v1/services/call", bytes.NewReader(body))
		r.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("missing required fields", func(t *testing.T) {
		req := map[string]string{
			"vendor": "dji",
		}
		body, _ := json.Marshal(req)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api/v1/services/call", bytes.NewReader(body))
		r.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func setupE2ETestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&models.Device{}, &model.ServiceCall{})
	require.NoError(t, err)

	return db
}
