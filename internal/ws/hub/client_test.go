package hub

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultClientConfig(t *testing.T) {
	config := DefaultClientConfig()

	assert.Equal(t, 10*time.Second, config.WriteWait)
	assert.Equal(t, 60*time.Second, config.PongWait)
	assert.Equal(t, 54*time.Second, config.PingPeriod)
	assert.Equal(t, int64(512*1024), config.MaxMessageSize)
	assert.Equal(t, 256, config.SendBufferSize)
}

func TestNewClient(t *testing.T) {
	t.Run("with nil config", func(t *testing.T) {
		client := NewClient("test-id", nil, nil, nil, nil)
		require.NotNil(t, client)
		assert.Equal(t, "test-id", client.ID)
		assert.NotNil(t, client.config)
		assert.NotNil(t, client.send)
		assert.NotNil(t, client.subscriptions)
	})

	t.Run("with custom config", func(t *testing.T) {
		config := &ClientConfig{
			WriteWait:      5 * time.Second,
			SendBufferSize: 100,
		}
		client := NewClient("test-id", nil, nil, config, nil)
		assert.Equal(t, 5*time.Second, client.config.WriteWait)
		assert.Equal(t, 100, client.config.SendBufferSize)
	})
}

func TestClient_CloseAndIsClosed(t *testing.T) {
	client := NewClient("test-id", nil, nil, nil, nil)

	assert.False(t, client.IsClosed())

	client.Close()
	assert.True(t, client.IsClosed())

	// Closing again should be no-op
	client.Close()
	assert.True(t, client.IsClosed())
}

func TestClient_Send(t *testing.T) {
	client := NewClient("test-id", nil, nil, nil, nil)

	msg := &Message{
		Type:  MessageTypeEvent,
		Event: "test-event",
	}

	// Send should succeed
	success := client.Send(msg)
	assert.True(t, success)

	// Verify message in channel
	select {
	case received := <-client.send:
		assert.Equal(t, "test-event", received.Event)
	default:
		t.Error("Message not in send channel")
	}

	// Send after close should fail
	client.Close()
	success = client.Send(msg)
	assert.False(t, success)
}

func TestClient_SendJSON(t *testing.T) {
	client := NewClient("test-id", nil, nil, nil, nil)

	success := client.SendJSON(MessageTypeEvent, "json-event", map[string]string{"key": "value"})
	assert.True(t, success)

	select {
	case received := <-client.send:
		assert.Equal(t, MessageTypeEvent, received.Type)
		assert.Equal(t, "json-event", received.Event)
	default:
		t.Error("Message not in send channel")
	}
}

func TestClient_SendError(t *testing.T) {
	client := NewClient("test-id", nil, nil, nil, nil)

	success := client.SendError("test error", "trace-123")
	assert.True(t, success)

	select {
	case received := <-client.send:
		assert.Equal(t, MessageTypeError, received.Type)
		assert.Equal(t, "test error", received.Error)
		assert.Equal(t, "trace-123", received.TraceID)
	default:
		t.Error("Message not in send channel")
	}
}

func TestClient_Subscriptions(t *testing.T) {
	client := NewClient("test-id", nil, nil, nil, nil)

	// Initially no subscriptions
	assert.Empty(t, client.GetSubscriptions())
	assert.False(t, client.IsSubscribed("topic1"))

	// Subscribe
	client.Subscribe("topic1")
	client.Subscribe("topic2")

	assert.True(t, client.IsSubscribed("topic1"))
	assert.True(t, client.IsSubscribed("topic2"))
	assert.False(t, client.IsSubscribed("topic3"))

	subs := client.GetSubscriptions()
	assert.Len(t, subs, 2)
	assert.Contains(t, subs, "topic1")
	assert.Contains(t, subs, "topic2")

	// Unsubscribe
	client.Unsubscribe("topic1")
	assert.False(t, client.IsSubscribed("topic1"))
	assert.True(t, client.IsSubscribed("topic2"))

	subs = client.GetSubscriptions()
	assert.Len(t, subs, 1)
}

func TestClient_Config(t *testing.T) {
	config := &ClientConfig{
		WriteWait: 5 * time.Second,
	}
	client := NewClient("test-id", nil, nil, config, nil)

	assert.Equal(t, config, client.Config())
}

func TestClient_Metadata(t *testing.T) {
	client := NewClient("test-id", nil, nil, nil, nil)

	client.Metadata["key1"] = "value1"
	client.Metadata["key2"] = "123"

	assert.Equal(t, "value1", client.Metadata["key1"])
	assert.Equal(t, "123", client.Metadata["key2"])
}

func TestClient_HandleMessage(t *testing.T) {
	client := NewClient("test-id", nil, nil, nil, nil)

	t.Run("ping message", func(t *testing.T) {
		client.handleMessage(&Message{Type: MessageTypePing})

		select {
		case msg := <-client.send:
			assert.Equal(t, MessageTypePong, msg.Type)
		default:
			t.Error("Pong not sent")
		}
	})

	t.Run("subscribe message", func(t *testing.T) {
		client.handleMessage(&Message{
			Type:  MessageTypeSubscribe,
			Event: "test-topic",
		})

		assert.True(t, client.IsSubscribed("test-topic"))

		select {
		case msg := <-client.send:
			assert.Equal(t, MessageTypeAck, msg.Type)
			assert.Equal(t, "test-topic", msg.Event)
		default:
			t.Error("Ack not sent")
		}
	})

	t.Run("unsubscribe message", func(t *testing.T) {
		client.Subscribe("unsub-topic")
		assert.True(t, client.IsSubscribed("unsub-topic"))

		client.handleMessage(&Message{
			Type:  MessageTypeUnsubscribe,
			Event: "unsub-topic",
		})

		assert.False(t, client.IsSubscribed("unsub-topic"))

		select {
		case msg := <-client.send:
			assert.Equal(t, MessageTypeAck, msg.Type)
		default:
			t.Error("Ack not sent")
		}
	})

	t.Run("event message with callback", func(t *testing.T) {
		var receivedMsg *Message
		client.SetOnMessage(func(c *Client, msg *Message) {
			receivedMsg = msg
		})

		client.handleMessage(&Message{
			Type:  MessageTypeEvent,
			Event: "custom-event",
			Data:  "test-data",
		})

		require.NotNil(t, receivedMsg)
		assert.Equal(t, "custom-event", receivedMsg.Event)
	})
}

// Integration test with actual WebSocket connection
func TestClient_WebSocketIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	hub := NewHub(nil, nil)
	hub.Start()
	defer hub.Stop()

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Logf("Upgrade error: %v", err)
			return
		}

		client := NewClient("ws-client", conn, hub, nil, nil)
		hub.Register(client)
		client.Start()
	}))
	defer server.Close()

	// Connect client
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)
	defer func() { _ = conn.Close() }()

	// Wait for registration
	time.Sleep(100 * time.Millisecond)

	assert.Equal(t, 1, hub.GetClientCount())

	// Send subscribe message
	subscribeMsg := Message{
		Type:  MessageTypeSubscribe,
		Event: "device.telemetry",
	}
	data, _ := json.Marshal(subscribeMsg)
	err = conn.WriteMessage(websocket.TextMessage, data)
	require.NoError(t, err)

	// Read ack response
	_, respData, err := conn.ReadMessage()
	require.NoError(t, err)

	var ackMsg Message
	err = json.Unmarshal(respData, &ackMsg)
	require.NoError(t, err)
	assert.Equal(t, MessageTypeAck, ackMsg.Type)
	assert.Equal(t, "device.telemetry", ackMsg.Event)

	// Send ping
	pingMsg := Message{Type: MessageTypePing}
	data, _ = json.Marshal(pingMsg)
	err = conn.WriteMessage(websocket.TextMessage, data)
	require.NoError(t, err)

	// Read pong response
	_, respData, err = conn.ReadMessage()
	require.NoError(t, err)

	var pongMsg Message
	err = json.Unmarshal(respData, &pongMsg)
	require.NoError(t, err)
	assert.Equal(t, MessageTypePong, pongMsg.Type)
}
