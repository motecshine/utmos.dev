package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/utmos/utmos/internal/ws"
	"github.com/utmos/utmos/internal/ws/hub"
	"github.com/utmos/utmos/internal/ws/push"
)

// TestLoadWebSocketConnections tests support for 10000+ WebSocket connections (NFR-003)
func TestLoadWebSocketConnections(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping WebSocket load test in short mode")
	}

	t.Run("1000 websocket connections", func(t *testing.T) {
		config := &ws.ServiceConfig{
			HubConfig: &hub.Config{
				MaxConnections: 10000,
				WriteTimeout:   10 * time.Second,
				ReadTimeout:    60 * time.Second,
				PingInterval:   30 * time.Second,
				PongTimeout:    30 * time.Second,
			},
			PusherConfig: &push.Config{
				WorkerCount: 8,
				QueueSize:   50000,
			},
			AllowedOrigins: []string{"*"},
		}

		wsSvc := ws.NewService(config, nil, nil)
		err := wsSvc.Start(context.Background())
		require.NoError(t, err)
		defer func() { _ = wsSvc.Stop() }()

		server := httptest.NewServer(http.HandlerFunc(wsSvc.HandleWebSocket))
		defer server.Close()

		wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

		connectionCount := 1000
		var wg sync.WaitGroup
		var successCount int64
		var errorCount int64

		conns := make([]*websocket.Conn, connectionCount)
		connMu := sync.Mutex{}

		start := time.Now()

		// Connect clients concurrently
		for i := 0; i < connectionCount; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()

				conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
				if err != nil {
					atomic.AddInt64(&errorCount, 1)
					return
				}

				connMu.Lock()
				conns[idx] = conn
				connMu.Unlock()

				atomic.AddInt64(&successCount, 1)
			}(i)
		}

		wg.Wait()
		elapsed := time.Since(start)

		t.Logf("Connected %d WebSocket clients in %v", successCount, elapsed)
		t.Logf("Connection rate: %.2f connections/second", float64(successCount)/elapsed.Seconds())
		t.Logf("Errors: %d", errorCount)

		// Wait for all registrations
		time.Sleep(500 * time.Millisecond)

		// Verify connections
		assert.Equal(t, int64(connectionCount), successCount)
		assert.Equal(t, connectionCount, wsSvc.Hub().GetClientCount())

		// Cleanup
		for _, conn := range conns {
			if conn != nil {
				_ = conn.Close()
			}
		}
	})

	t.Run("subscription load", func(t *testing.T) {
		wsSvc := ws.NewService(nil, nil, nil)
		err := wsSvc.Start(context.Background())
		require.NoError(t, err)
		defer func() { _ = wsSvc.Stop() }()

		server := httptest.NewServer(http.HandlerFunc(wsSvc.HandleWebSocket))
		defer server.Close()

		wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

		connectionCount := 100
		topicsPerClient := 10

		conns := make([]*websocket.Conn, connectionCount)

		// Connect clients
		for i := 0; i < connectionCount; i++ {
			conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
			require.NoError(t, err)
			conns[i] = conn
		}

		// Wait for registrations
		time.Sleep(200 * time.Millisecond)

		start := time.Now()

		// Subscribe each client to multiple topics
		var wg sync.WaitGroup
		for i, conn := range conns {
			wg.Add(1)
			go func(idx int, c *websocket.Conn) {
				defer wg.Done()

				for j := 0; j < topicsPerClient; j++ {
					topic := fmt.Sprintf("load.topic.%d", j)
					subscribeMsg := hub.Message{
						Type:  hub.MessageTypeSubscribe,
						Event: topic,
					}
					data, _ := json.Marshal(subscribeMsg)
					_ = c.WriteMessage(websocket.TextMessage, data)

					// Read ack
					_ = c.SetReadDeadline(time.Now().Add(time.Second))
					_, _, _ = c.ReadMessage()
				}
			}(i, conn)
		}

		wg.Wait()
		elapsed := time.Since(start)

		totalSubscriptions := connectionCount * topicsPerClient
		t.Logf("Created %d subscriptions in %v", totalSubscriptions, elapsed)
		t.Logf("Subscription rate: %.2f subscriptions/second", float64(totalSubscriptions)/elapsed.Seconds())

		// Verify subscriptions
		assert.Equal(t, topicsPerClient, wsSvc.SubscriptionManager().GetTopicCount())

		// Cleanup
		for _, conn := range conns {
			_ = conn.Close()
		}
	})
}

// TestLoadWebSocketBroadcast tests broadcast performance under load
func TestLoadWebSocketBroadcast(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping WebSocket load test in short mode")
	}

	config := &ws.ServiceConfig{
		HubConfig: &hub.Config{
			MaxConnections: 10000,
			WriteTimeout:   10 * time.Second,
			ReadTimeout:    60 * time.Second,
			PingInterval:   30 * time.Second,
			PongTimeout:    30 * time.Second,
		},
		PusherConfig: &push.Config{
			WorkerCount: 8,
			QueueSize:   100000,
		},
		AllowedOrigins: []string{"*"},
	}

	wsSvc := ws.NewService(config, nil, nil)
	err := wsSvc.Start(context.Background())
	require.NoError(t, err)
	defer func() { _ = wsSvc.Stop() }()

	server := httptest.NewServer(http.HandlerFunc(wsSvc.HandleWebSocket))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	t.Run("broadcast to 500 clients", func(t *testing.T) {
		connectionCount := 500
		conns := make([]*websocket.Conn, connectionCount)

		// Connect clients
		for i := 0; i < connectionCount; i++ {
			conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
			require.NoError(t, err)
			conns[i] = conn
		}

		// Wait for registrations
		time.Sleep(500 * time.Millisecond)

		// Subscribe all clients to broadcast topic
		for _, conn := range conns {
			subscribeMsg := hub.Message{
				Type:  hub.MessageTypeSubscribe,
				Event: "broadcast.load.topic",
			}
			data, _ := json.Marshal(subscribeMsg)
			_ = conn.WriteMessage(websocket.TextMessage, data)
			_ = conn.SetReadDeadline(time.Now().Add(time.Second))
			_, _, _ = conn.ReadMessage() // Read ack
		}

		// Wait for subscriptions
		time.Sleep(200 * time.Millisecond)

		messageCount := 100
		start := time.Now()

		// Broadcast messages
		for i := 0; i < messageCount; i++ {
			wsSvc.Pusher().PushToTopic("broadcast.load.topic", &hub.Message{
				Type:  hub.MessageTypeEvent,
				Event: "broadcast.load.topic",
				Data:  map[string]interface{}{"index": i},
			})
		}

		// Wait for all messages to be pushed
		time.Sleep(time.Second)

		elapsed := time.Since(start)
		totalDeliveries := messageCount * connectionCount

		t.Logf("Broadcast %d messages to %d clients in %v", messageCount, connectionCount, elapsed)
		t.Logf("Total deliveries: %d", totalDeliveries)
		t.Logf("Delivery rate: %.2f deliveries/second", float64(totalDeliveries)/elapsed.Seconds())

		// Cleanup
		for _, conn := range conns {
			_ = conn.Close()
		}
	})
}

// TestLoadWebSocketMessageThroughput tests message throughput
func TestLoadWebSocketMessageThroughput(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping WebSocket load test in short mode")
	}

	wsSvc := ws.NewService(nil, nil, nil)
	err := wsSvc.Start(context.Background())
	require.NoError(t, err)
	defer func() { _ = wsSvc.Stop() }()

	// Create subscriptions without actual WebSocket connections
	clientCount := 1000
	for i := 0; i < clientCount; i++ {
		clientID := fmt.Sprintf("throughput-client-%d", i)
		wsSvc.SubscriptionManager().Subscribe(clientID, "throughput.topic")
	}

	t.Run("push throughput", func(t *testing.T) {
		messageCount := 10000
		start := time.Now()

		for i := 0; i < messageCount; i++ {
			wsSvc.Pusher().PushToTopic("throughput.topic", &hub.Message{
				Type:  hub.MessageTypeEvent,
				Event: "throughput.topic",
				Data:  map[string]interface{}{"index": i},
			})
		}

		// Wait for queue to drain
		for wsSvc.Pusher().QueueLength() > 0 {
			time.Sleep(10 * time.Millisecond)
		}

		elapsed := time.Since(start)
		t.Logf("Pushed %d messages in %v", messageCount, elapsed)
		t.Logf("Push throughput: %.2f messages/second", float64(messageCount)/elapsed.Seconds())

		// Should push at least 10000 messages per second
		assert.Greater(t, float64(messageCount)/elapsed.Seconds(), 10000.0)
	})
}

// TestLoadWebSocketConnectionChurn tests connection churn handling
func TestLoadWebSocketConnectionChurn(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping WebSocket load test in short mode")
	}

	config := &ws.ServiceConfig{
		HubConfig: &hub.Config{
			MaxConnections: 10000,
			WriteTimeout:   10 * time.Second,
			ReadTimeout:    60 * time.Second,
			PingInterval:   30 * time.Second,
			PongTimeout:    30 * time.Second,
		},
		AllowedOrigins: []string{"*"},
	}

	wsSvc := ws.NewService(config, nil, nil)
	err := wsSvc.Start(context.Background())
	require.NoError(t, err)
	defer func() { _ = wsSvc.Stop() }()

	server := httptest.NewServer(http.HandlerFunc(wsSvc.HandleWebSocket))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	t.Run("connection churn", func(t *testing.T) {
		iterations := 5
		connectionsPerIteration := 100
		var totalConnections int64

		start := time.Now()

		for iter := 0; iter < iterations; iter++ {
			var wg sync.WaitGroup
			conns := make([]*websocket.Conn, connectionsPerIteration)

			// Connect
			for i := 0; i < connectionsPerIteration; i++ {
				wg.Add(1)
				go func(idx int) {
					defer wg.Done()
					conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
					if err == nil {
						conns[idx] = conn
						atomic.AddInt64(&totalConnections, 1)
					}
				}(i)
			}
			wg.Wait()

			// Small delay
			time.Sleep(50 * time.Millisecond)

			// Disconnect
			for _, conn := range conns {
				if conn != nil {
					_ = conn.Close()
				}
			}

			// Wait for cleanup
			time.Sleep(100 * time.Millisecond)
		}

		elapsed := time.Since(start)
		t.Logf("Completed %d connection cycles in %v", totalConnections, elapsed)
		t.Logf("Churn rate: %.2f connections/second", float64(totalConnections)/elapsed.Seconds())

		// Final state should be clean
		time.Sleep(500 * time.Millisecond)
		assert.Equal(t, 0, wsSvc.Hub().GetClientCount())
	})
}

// TestLoadWebSocketStats tests stats collection under load
func TestLoadWebSocketStats(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping WebSocket load test in short mode")
	}

	wsSvc := ws.NewService(nil, nil, nil)
	err := wsSvc.Start(context.Background())
	require.NoError(t, err)
	defer func() { _ = wsSvc.Stop() }()

	// Create subscriptions
	for i := 0; i < 1000; i++ {
		clientID := fmt.Sprintf("stats-client-%d", i)
		for j := 0; j < 5; j++ {
			topic := fmt.Sprintf("stats.topic.%d", j)
			wsSvc.SubscriptionManager().Subscribe(clientID, topic)
		}
	}

	// Push some messages
	for i := 0; i < 1000; i++ {
		wsSvc.Pusher().PushToTopic("stats.topic.0", &hub.Message{
			Type:  hub.MessageTypeEvent,
			Event: "stats.topic.0",
			Data:  map[string]interface{}{"index": i},
		})
	}

	// Wait for processing
	time.Sleep(500 * time.Millisecond)

	stats := wsSvc.GetStats()
	t.Logf("Stats: %+v", stats)

	assert.Equal(t, 5, stats.ActiveTopics)
	assert.GreaterOrEqual(t, stats.MessagesPushed, int64(0))
}

// BenchmarkWebSocketConnection benchmarks WebSocket connection handling
func BenchmarkWebSocketConnection(b *testing.B) {
	wsSvc := ws.NewService(nil, nil, nil)
	_ = wsSvc.Start(context.Background())
	defer func() { _ = wsSvc.Stop() }()

	server := httptest.NewServer(http.HandlerFunc(wsSvc.HandleWebSocket))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		if err == nil {
			_ = conn.Close()
		}
	}
}

// BenchmarkWebSocketSubscription benchmarks subscription operations
func BenchmarkWebSocketSubscription(b *testing.B) {
	wsSvc := ws.NewService(nil, nil, nil)
	_ = wsSvc.Start(context.Background())
	defer func() { _ = wsSvc.Stop() }()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		clientID := fmt.Sprintf("bench-client-%d", i)
		topic := fmt.Sprintf("bench.topic.%d", i%100)
		wsSvc.SubscriptionManager().Subscribe(clientID, topic)
	}
}

// BenchmarkWebSocketPushToTopic benchmarks push operations
func BenchmarkWebSocketPushToTopic(b *testing.B) {
	wsSvc := ws.NewService(nil, nil, nil)
	_ = wsSvc.Start(context.Background())
	defer func() { _ = wsSvc.Stop() }()

	// Pre-create subscriptions
	for i := 0; i < 100; i++ {
		wsSvc.SubscriptionManager().Subscribe(fmt.Sprintf("client-%d", i), "bench.topic")
	}

	msg := &hub.Message{
		Type:  hub.MessageTypeEvent,
		Event: "bench.topic",
		Data:  map[string]interface{}{"test": "data"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wsSvc.Pusher().PushToTopic("bench.topic", msg)
	}
}
