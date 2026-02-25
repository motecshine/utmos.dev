package ws

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/utmos/utmos/internal/ws/hub"
	"github.com/utmos/utmos/pkg/rabbitmq"
)

func TestDefaultServiceConfig(t *testing.T) {
	config := DefaultServiceConfig()

	assert.NotNil(t, config.HubConfig)
	assert.NotNil(t, config.ClientConfig)
	assert.NotNil(t, config.PusherConfig)
	assert.Contains(t, config.AllowedOrigins, "*")
}

func TestNewService(t *testing.T) {
	t.Run("with nil config", func(t *testing.T) {
		svc := NewService(nil, nil, nil)
		require.NotNil(t, svc)
		assert.NotNil(t, svc.hub)
		assert.NotNil(t, svc.subManager)
		assert.NotNil(t, svc.pusher)
	})

	t.Run("with custom config", func(t *testing.T) {
		config := &ServiceConfig{
			AllowedOrigins: []string{"http://localhost:3000"},
		}
		svc := NewService(config, nil, nil)
		assert.Equal(t, config.AllowedOrigins, svc.config.AllowedOrigins)
	})
}

func TestService_StartStop(t *testing.T) {
	svc := NewService(nil, nil, nil)

	assert.False(t, svc.IsRunning())

	err := svc.Start(context.Background())
	require.NoError(t, err)
	assert.True(t, svc.IsRunning())

	// Starting again should be no-op
	err = svc.Start(context.Background())
	require.NoError(t, err)
	assert.True(t, svc.IsRunning())

	err = svc.Stop()
	require.NoError(t, err)
	assert.False(t, svc.IsRunning())

	// Stopping again should be no-op
	err = svc.Stop()
	require.NoError(t, err)
	assert.False(t, svc.IsRunning())
}

func TestService_HandleWebSocket_NotRunning(t *testing.T) {
	svc := NewService(nil, nil, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/ws", nil)

	svc.HandleWebSocket(w, r)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
}

func TestService_HandleWebSocket(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	svc := NewService(nil, nil, nil)
	err := svc.Start(context.Background())
	require.NoError(t, err)
	defer func() { _ = svc.Stop() }()

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(svc.HandleWebSocket))
	defer server.Close()

	// Connect WebSocket client
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "?device_sn=DRONE001&user_id=user123"
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)
	defer func() { _ = conn.Close() }()

	// Wait for registration
	time.Sleep(100 * time.Millisecond)

	assert.Equal(t, 1, svc.Hub().GetClientCount())
}

func TestService_WebSocketSubscription(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	svc := NewService(nil, nil, nil)
	err := svc.Start(context.Background())
	require.NoError(t, err)
	defer func() { _ = svc.Stop() }()

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(svc.HandleWebSocket))
	defer server.Close()

	// Connect WebSocket client
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)
	defer func() { _ = conn.Close() }()

	// Wait for registration
	time.Sleep(100 * time.Millisecond)

	// Subscribe to topic
	subscribeMsg := hub.Message{
		Type:  hub.MessageTypeSubscribe,
		Event: "device.telemetry",
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

	// Wait for subscription to be processed
	// Use polling with timeout instead of fixed sleep to handle race conditions
	var topicCount int
	for i := 0; i < 20; i++ {
		time.Sleep(50 * time.Millisecond)
		topicCount = svc.SubscriptionManager().GetTopicCount()
		if topicCount > 0 {
			break
		}
	}

	// Verify subscription in manager
	assert.Equal(t, 1, topicCount, "Expected 1 topic subscription")
}

func TestService_GetStats(t *testing.T) {
	svc := NewService(nil, nil, nil)
	err := svc.Start(context.Background())
	require.NoError(t, err)
	defer func() { _ = svc.Stop() }()

	stats := svc.GetStats()

	assert.Equal(t, 0, stats.ConnectedClients)
	assert.Equal(t, 0, stats.ActiveTopics)
	assert.Equal(t, int64(0), stats.MessagesPushed)
	assert.Equal(t, int64(0), stats.MessagesDropped)
	assert.Equal(t, 0, stats.QueueLength)
}

func TestService_ActionToTopic(t *testing.T) {
	svc := NewService(nil, nil, nil)

	tests := []struct {
		action   string
		expected string
	}{
		{"telemetry.update", "telemetry.update"},
		{"device.status", "device.status"},
		{"command.response", "command.response"},
	}

	for _, tt := range tests {
		t.Run(tt.action, func(t *testing.T) {
			result := svc.actionToTopic(tt.action)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestService_HandleRabbitMQMessage(t *testing.T) {
	svc := NewService(nil, nil, nil)
	err := svc.Start(context.Background())
	require.NoError(t, err)
	defer func() { _ = svc.Stop() }()

	// Subscribe a mock client
	svc.SubscriptionManager().Subscribe("client1", "telemetry.update")

	t.Run("nil message", func(t *testing.T) {
		// Should not panic
		svc.handleRabbitMQMessage(context.Background(), nil)
	})

	t.Run("valid message", func(t *testing.T) {
		msg := &rabbitmq.StandardMessage{
			DeviceSN: "DRONE001",
			Action:   "telemetry.update",
			TID:      "trace-123",
			Data:     json.RawMessage(`{"temperature": 25.5}`),
		}
		svc.handleRabbitMQMessage(context.Background(), msg)

		// Wait for processing
		time.Sleep(50 * time.Millisecond)
	})

	t.Run("message with invalid JSON data", func(t *testing.T) {
		msg := &rabbitmq.StandardMessage{
			DeviceSN: "DRONE001",
			Action:   "telemetry.update",
			Data:     json.RawMessage(`invalid json`),
		}
		svc.handleRabbitMQMessage(context.Background(), msg)
		// Should handle gracefully
	})
}

func TestService_Accessors(t *testing.T) {
	svc := NewService(nil, nil, nil)

	assert.NotNil(t, svc.Hub())
	assert.NotNil(t, svc.SubscriptionManager())
	assert.NotNil(t, svc.Pusher())
}

func TestService_OnClientCallbacks(t *testing.T) {
	svc := NewService(nil, nil, nil)
	err := svc.Start(context.Background())
	require.NoError(t, err)
	defer func() { _ = svc.Stop() }()

	// Create a mock client
	client := &hub.Client{
		ID:       "test-client",
		DeviceSN: "DRONE001",
	}

	// Test connect callback
	svc.onClientConnect(client)

	// Subscribe client to topic
	svc.SubscriptionManager().Subscribe(client.ID, "test.topic")
	assert.True(t, svc.SubscriptionManager().IsSubscribed(client.ID, "test.topic"))

	// Test disconnect callback
	svc.onClientDisconnect(client)

	// Verify unsubscribed
	assert.False(t, svc.SubscriptionManager().IsSubscribed(client.ID, "test.topic"))
}

func TestService_OnClientMessage(t *testing.T) {
	svc := NewService(nil, nil, nil)
	err := svc.Start(context.Background())
	require.NoError(t, err)
	defer func() { _ = svc.Stop() }()

	client := &hub.Client{
		ID: "test-client",
	}

	t.Run("subscribe message", func(t *testing.T) {
		msg := &hub.Message{
			Type:  hub.MessageTypeSubscribe,
			Event: "device.telemetry",
		}
		svc.onClientMessage(client, msg)

		assert.True(t, svc.SubscriptionManager().IsSubscribed(client.ID, "device.telemetry"))
	})

	t.Run("unsubscribe message", func(t *testing.T) {
		msg := &hub.Message{
			Type:  hub.MessageTypeUnsubscribe,
			Event: "device.telemetry",
		}
		svc.onClientMessage(client, msg)

		assert.False(t, svc.SubscriptionManager().IsSubscribed(client.ID, "device.telemetry"))
	})

	t.Run("other message type", func(t *testing.T) {
		msg := &hub.Message{
			Type:  hub.MessageTypeEvent,
			Event: "custom.event",
		}
		// Should not panic
		svc.onClientMessage(client, msg)
	})
}

func TestService_OriginCheck(t *testing.T) {
	t.Run("allow all origins", func(t *testing.T) {
		config := &ServiceConfig{
			AllowedOrigins: []string{"*"},
		}
		svc := NewService(config, nil, nil)

		r := httptest.NewRequest("GET", "/ws", nil)
		r.Header.Set("Origin", "http://example.com")

		assert.True(t, svc.upgrader.CheckOrigin(r))
	})

	t.Run("specific origin allowed", func(t *testing.T) {
		config := &ServiceConfig{
			AllowedOrigins: []string{"http://localhost:3000"},
		}
		svc := NewService(config, nil, nil)

		r := httptest.NewRequest("GET", "/ws", nil)
		r.Header.Set("Origin", "http://localhost:3000")

		assert.True(t, svc.upgrader.CheckOrigin(r))
	})

	t.Run("origin not allowed", func(t *testing.T) {
		config := &ServiceConfig{
			AllowedOrigins: []string{"http://localhost:3000"},
		}
		svc := NewService(config, nil, nil)

		r := httptest.NewRequest("GET", "/ws", nil)
		r.Header.Set("Origin", "http://evil.com")

		assert.False(t, svc.upgrader.CheckOrigin(r))
	})

	t.Run("empty allowed origins", func(t *testing.T) {
		config := &ServiceConfig{
			AllowedOrigins: []string{},
		}
		svc := NewService(config, nil, nil)

		r := httptest.NewRequest("GET", "/ws", nil)
		r.Header.Set("Origin", "http://any.com")

		assert.True(t, svc.upgrader.CheckOrigin(r))
	})
}
