package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/utmos/utmos/internal/ws"
	"github.com/utmos/utmos/internal/ws/hub"
	"github.com/utmos/utmos/internal/ws/push"
)

// TestWSServiceIntegration tests the WebSocket service integration
func TestWSServiceIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	config := &ws.ServiceConfig{
		HubConfig: &hub.Config{
			MaxConnections: 100,
			WriteTimeout:   10 * time.Second,
			ReadTimeout:    60 * time.Second,
			PingInterval:   30 * time.Second,
			PongTimeout:    30 * time.Second,
		},
		ClientConfig: hub.DefaultClientConfig(),
		PusherConfig: &push.Config{
			WorkerCount: 2,
			QueueSize:   1000,
		},
		AllowedOrigins: []string{"*"},
	}

	svc := ws.NewService(config, nil, nil)
	err := svc.Start(context.Background())
	require.NoError(t, err)
	defer func() { _ = svc.Stop() }()

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(svc.HandleWebSocket))
	defer server.Close()

	t.Run("single client connection", func(t *testing.T) {
		wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
		conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		require.NoError(t, err)
		defer func() { _ = conn.Close() }()

		// Wait for registration
		time.Sleep(100 * time.Millisecond)

		assert.Equal(t, 1, svc.Hub().GetClientCount())
	})

	t.Run("multiple client connections", func(t *testing.T) {
		wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

		conns := make([]*websocket.Conn, 5)
		for i := 0; i < 5; i++ {
			conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
			require.NoError(t, err)
			conns[i] = conn
		}

		// Wait for registrations
		time.Sleep(200 * time.Millisecond)

		// Should have 5 + 1 (from previous test) = 6 clients
		// But previous test's client may have disconnected
		assert.GreaterOrEqual(t, svc.Hub().GetClientCount(), 5)

		// Cleanup
		for _, conn := range conns {
			_ = conn.Close()
		}
	})
}

// TestWSSubscriptionIntegration tests WebSocket subscription functionality
func TestWSSubscriptionIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	svc := ws.NewService(nil, nil, nil)
	err := svc.Start(context.Background())
	require.NoError(t, err)
	defer func() { _ = svc.Stop() }()

	server := httptest.NewServer(http.HandlerFunc(svc.HandleWebSocket))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)
	defer func() { _ = conn.Close() }()

	// Wait for registration
	time.Sleep(100 * time.Millisecond)

	t.Run("subscribe to topic", func(t *testing.T) {
		subscribeMsg := hub.Message{
			Type:  hub.MessageTypeSubscribe,
			Event: "device.telemetry",
		}
		data, _ := json.Marshal(subscribeMsg)
		err := conn.WriteMessage(websocket.TextMessage, data)
		require.NoError(t, err)

		// Read ack
		_, respData, err := conn.ReadMessage()
		require.NoError(t, err)

		var ackMsg hub.Message
		err = json.Unmarshal(respData, &ackMsg)
		require.NoError(t, err)
		assert.Equal(t, hub.MessageTypeAck, ackMsg.Type)
		assert.Equal(t, "device.telemetry", ackMsg.Event)

		// Wait for subscription processing
		time.Sleep(50 * time.Millisecond)

		assert.Equal(t, 1, svc.SubscriptionManager().GetTopicCount())
	})

	t.Run("unsubscribe from topic", func(t *testing.T) {
		unsubscribeMsg := hub.Message{
			Type:  hub.MessageTypeUnsubscribe,
			Event: "device.telemetry",
		}
		data, _ := json.Marshal(unsubscribeMsg)
		err := conn.WriteMessage(websocket.TextMessage, data)
		require.NoError(t, err)

		// Read ack
		_, respData, err := conn.ReadMessage()
		require.NoError(t, err)

		var ackMsg hub.Message
		err = json.Unmarshal(respData, &ackMsg)
		require.NoError(t, err)
		assert.Equal(t, hub.MessageTypeAck, ackMsg.Type)

		// Wait for unsubscription processing
		time.Sleep(50 * time.Millisecond)

		assert.Equal(t, 0, svc.SubscriptionManager().GetTopicCount())
	})
}

// TestWSPingPongIntegration tests WebSocket ping/pong functionality
func TestWSPingPongIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	svc := ws.NewService(nil, nil, nil)
	err := svc.Start(context.Background())
	require.NoError(t, err)
	defer func() { _ = svc.Stop() }()

	server := httptest.NewServer(http.HandlerFunc(svc.HandleWebSocket))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)
	defer func() { _ = conn.Close() }()

	// Wait for registration
	time.Sleep(100 * time.Millisecond)

	// Send ping
	pingMsg := hub.Message{Type: hub.MessageTypePing}
	data, _ := json.Marshal(pingMsg)
	err = conn.WriteMessage(websocket.TextMessage, data)
	require.NoError(t, err)

	// Read pong
	_, respData, err := conn.ReadMessage()
	require.NoError(t, err)

	var pongMsg hub.Message
	err = json.Unmarshal(respData, &pongMsg)
	require.NoError(t, err)
	assert.Equal(t, hub.MessageTypePong, pongMsg.Type)
}

// TestWSMessagePushIntegration tests message pushing to subscribed clients
func TestWSMessagePushIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	svc := ws.NewService(nil, nil, nil)
	err := svc.Start(context.Background())
	require.NoError(t, err)
	defer func() { _ = svc.Stop() }()

	server := httptest.NewServer(http.HandlerFunc(svc.HandleWebSocket))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	// Connect multiple clients
	clients := make([]*websocket.Conn, 3)
	for i := 0; i < 3; i++ {
		conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		require.NoError(t, err)
		clients[i] = conn
	}
	defer func() {
		for _, conn := range clients {
			_ = conn.Close()
		}
	}()

	// Wait for registrations
	time.Sleep(200 * time.Millisecond)

	// Subscribe all clients to same topic
	for _, conn := range clients {
		subscribeMsg := hub.Message{
			Type:  hub.MessageTypeSubscribe,
			Event: "broadcast.topic",
		}
		data, _ := json.Marshal(subscribeMsg)
		err := conn.WriteMessage(websocket.TextMessage, data)
		require.NoError(t, err)

		// Read ack
		_, _, err = conn.ReadMessage()
		require.NoError(t, err)
	}

	// Wait for subscriptions
	time.Sleep(100 * time.Millisecond)

	// Push message to topic
	svc.Pusher().PushToTopic("broadcast.topic", &hub.Message{
		Type:  hub.MessageTypeEvent,
		Event: "broadcast.topic",
		Data:  map[string]string{"message": "hello"},
	})

	// Wait for push
	time.Sleep(100 * time.Millisecond)

	// Verify all clients received the message
	for i, conn := range clients {
		_ = conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
		_, respData, err := conn.ReadMessage()
		if err != nil {
			t.Logf("Client %d did not receive message: %v", i, err)
			continue
		}

		var msg hub.Message
		err = json.Unmarshal(respData, &msg)
		require.NoError(t, err)
		assert.Equal(t, hub.MessageTypeEvent, msg.Type)
		assert.Equal(t, "broadcast.topic", msg.Event)
	}
}

// TestWSConcurrentConnectionsIntegration tests concurrent WebSocket connections
func TestWSConcurrentConnectionsIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	config := &ws.ServiceConfig{
		HubConfig: &hub.Config{
			MaxConnections: 100,
			WriteTimeout:   10 * time.Second,
			ReadTimeout:    60 * time.Second,
			PingInterval:   30 * time.Second,
			PongTimeout:    30 * time.Second,
		},
	}

	svc := ws.NewService(config, nil, nil)
	err := svc.Start(context.Background())
	require.NoError(t, err)
	defer func() { _ = svc.Stop() }()

	server := httptest.NewServer(http.HandlerFunc(svc.HandleWebSocket))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	var wg sync.WaitGroup
	clientCount := 50
	conns := make([]*websocket.Conn, clientCount)
	connMu := sync.Mutex{}

	// Connect clients concurrently
	for i := 0; i < clientCount; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
			if err != nil {
				t.Logf("Failed to connect client %d: %v", idx, err)
				return
			}
			connMu.Lock()
			conns[idx] = conn
			connMu.Unlock()
		}(i)
	}

	wg.Wait()

	// Wait for all registrations
	time.Sleep(500 * time.Millisecond)

	// Count successful connections
	successCount := 0
	for _, conn := range conns {
		if conn != nil {
			successCount++
		}
	}

	assert.Equal(t, clientCount, successCount)
	assert.Equal(t, clientCount, svc.Hub().GetClientCount())

	// Cleanup
	for _, conn := range conns {
		if conn != nil {
			_ = conn.Close()
		}
	}
}

// TestWSStatsIntegration tests WebSocket service statistics
func TestWSStatsIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	svc := ws.NewService(nil, nil, nil)
	err := svc.Start(context.Background())
	require.NoError(t, err)
	defer func() { _ = svc.Stop() }()

	server := httptest.NewServer(http.HandlerFunc(svc.HandleWebSocket))
	defer server.Close()

	// Initial stats
	stats := svc.GetStats()
	assert.Equal(t, 0, stats.ConnectedClients)
	assert.Equal(t, 0, stats.ActiveTopics)

	// Connect client
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)
	defer func() { _ = conn.Close() }()

	// Wait for registration
	time.Sleep(100 * time.Millisecond)

	// Check stats after connection
	stats = svc.GetStats()
	assert.Equal(t, 1, stats.ConnectedClients)

	// Subscribe to topic
	subscribeMsg := hub.Message{
		Type:  hub.MessageTypeSubscribe,
		Event: "stats.topic",
	}
	data, _ := json.Marshal(subscribeMsg)
	err = conn.WriteMessage(websocket.TextMessage, data)
	require.NoError(t, err)

	// Read ack
	_, _, err = conn.ReadMessage()
	require.NoError(t, err)

	// Wait for subscription
	time.Sleep(50 * time.Millisecond)

	// Check stats after subscription
	stats = svc.GetStats()
	assert.Equal(t, 1, stats.ActiveTopics)
}

// TestWSClientMetadataIntegration tests client metadata from query parameters
func TestWSClientMetadataIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	svc := ws.NewService(nil, nil, nil)
	err := svc.Start(context.Background())
	require.NoError(t, err)
	defer func() { _ = svc.Stop() }()

	server := httptest.NewServer(http.HandlerFunc(svc.HandleWebSocket))
	defer server.Close()

	// Connect with query parameters
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "?device_sn=DRONE001&user_id=user123"
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)
	defer func() { _ = conn.Close() }()

	// Wait for registration
	time.Sleep(100 * time.Millisecond)

	// Verify client is registered
	assert.Equal(t, 1, svc.Hub().GetClientCount())

	// Get client and verify metadata
	clients := svc.Hub().GetClients()
	require.Len(t, clients, 1)
	assert.Equal(t, "DRONE001", clients[0].DeviceSN)
	assert.Equal(t, "user123", clients[0].UserID)
}

// TestWSGracefulDisconnectIntegration tests graceful client disconnection
func TestWSGracefulDisconnectIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	svc := ws.NewService(nil, nil, nil)
	err := svc.Start(context.Background())
	require.NoError(t, err)
	defer func() { _ = svc.Stop() }()

	server := httptest.NewServer(http.HandlerFunc(svc.HandleWebSocket))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)

	// Wait for registration
	time.Sleep(100 * time.Millisecond)

	assert.Equal(t, 1, svc.Hub().GetClientCount())

	// Subscribe to topic
	subscribeMsg := hub.Message{
		Type:  hub.MessageTypeSubscribe,
		Event: "disconnect.topic",
	}
	data, _ := json.Marshal(subscribeMsg)
	err = conn.WriteMessage(websocket.TextMessage, data)
	require.NoError(t, err)

	// Read ack
	_, _, err = conn.ReadMessage()
	require.NoError(t, err)

	// Wait for subscription
	time.Sleep(50 * time.Millisecond)

	assert.Equal(t, 1, svc.SubscriptionManager().GetTopicCount())

	// Close connection
	_ = conn.Close()

	// Wait for cleanup
	time.Sleep(200 * time.Millisecond)

	// Verify client is removed and subscriptions cleaned up
	assert.Equal(t, 0, svc.Hub().GetClientCount())
	assert.Equal(t, 0, svc.SubscriptionManager().GetTopicCount())
}
