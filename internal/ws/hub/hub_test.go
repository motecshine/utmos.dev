package hub

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newTestClient creates a client for testing purposes
func newTestClient(id string) *Client {
	return &Client{
		ID:            id,
		send:          make(chan *Message, 256),
		subscriptions: make(map[string]bool),
		done:          make(chan struct{}),
		Metadata:      make(Metadata),
	}
}

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	assert.Equal(t, 10000, config.MaxConnections)
	assert.Equal(t, 10*time.Second, config.WriteTimeout)
	assert.Equal(t, 60*time.Second, config.ReadTimeout)
	assert.Equal(t, 30*time.Second, config.PingInterval)
	assert.Equal(t, 30*time.Second, config.PongTimeout)
}

func TestNewHub(t *testing.T) {
	t.Run("with nil config", func(t *testing.T) {
		hub := NewHub(nil, nil)
		require.NotNil(t, hub)
		assert.NotNil(t, hub.config)
		assert.Equal(t, 10000, hub.config.MaxConnections)
	})

	t.Run("with custom config", func(t *testing.T) {
		config := &Config{
			MaxConnections: 100,
			WriteTimeout:   5 * time.Second,
		}
		hub := NewHub(config, nil)
		require.NotNil(t, hub)
		assert.Equal(t, 100, hub.config.MaxConnections)
	})
}

func TestHub_StartStop(t *testing.T) {
	hub := NewHub(nil, nil)

	assert.False(t, hub.IsRunning())

	hub.Start()
	assert.True(t, hub.IsRunning())

	// Starting again should be no-op
	hub.Start()
	assert.True(t, hub.IsRunning())

	hub.Stop()
	assert.False(t, hub.IsRunning())

	// Stopping again should be no-op
	hub.Stop()
	assert.False(t, hub.IsRunning())
}

func TestHub_RegisterUnregister(t *testing.T) {
	hub := NewHub(nil, nil)
	hub.Start()
	defer hub.Stop()

	// Create a test client
	client := newTestClient("test-client-1")

	// Register client
	hub.Register(client)

	// Wait for registration to be processed
	time.Sleep(50 * time.Millisecond)

	assert.Equal(t, 1, hub.GetClientCount())

	// Get client
	retrieved, exists := hub.GetClient("test-client-1")
	assert.True(t, exists)
	assert.Equal(t, client.ID, retrieved.ID)

	// Unregister client
	hub.Unregister(client)

	// Wait for unregistration to be processed
	time.Sleep(50 * time.Millisecond)

	assert.Equal(t, 0, hub.GetClientCount())

	_, exists = hub.GetClient("test-client-1")
	assert.False(t, exists)
}

func TestHub_Broadcast(t *testing.T) {
	hub := NewHub(nil, nil)
	hub.Start()
	defer hub.Stop()

	// Create multiple clients
	clients := make([]*Client, 3)
	for i := 0; i < 3; i++ {
		clients[i] = newTestClient(fmt.Sprintf("client-%d", i))
		hub.Register(clients[i])
	}

	// Wait for registrations
	time.Sleep(50 * time.Millisecond)

	assert.Equal(t, 3, hub.GetClientCount())

	// Broadcast message
	msg := &Message{
		Type:  MessageTypeEvent,
		Event: "test-event",
		Data:  map[string]string{"key": "value"},
	}
	hub.Broadcast(msg)

	// Wait for broadcast
	time.Sleep(50 * time.Millisecond)

	// Check all clients received the message
	for _, client := range clients {
		select {
		case received := <-client.send:
			assert.Equal(t, MessageTypeEvent, received.Type)
			assert.Equal(t, "test-event", received.Event)
		default:
			t.Error("Client did not receive broadcast message")
		}
	}
}

func TestHub_SendToClient(t *testing.T) {
	hub := NewHub(nil, nil)
	hub.Start()
	defer hub.Stop()

	client := newTestClient("target-client")
	hub.Register(client)

	// Wait for registration
	time.Sleep(50 * time.Millisecond)

	msg := &Message{
		Type:  MessageTypeEvent,
		Event: "direct-message",
	}

	// Send to existing client
	success := hub.SendToClient("target-client", msg)
	assert.True(t, success)

	// Verify message received
	select {
	case received := <-client.send:
		assert.Equal(t, "direct-message", received.Event)
	case <-time.After(100 * time.Millisecond):
		t.Error("Client did not receive message")
	}

	// Send to non-existing client
	success = hub.SendToClient("non-existing", msg)
	assert.False(t, success)
}

func TestHub_GetClients(t *testing.T) {
	hub := NewHub(nil, nil)
	hub.Start()
	defer hub.Stop()

	// Register clients
	for i := 0; i < 5; i++ {
		client := newTestClient(fmt.Sprintf("client-%d", i))
		hub.Register(client)
	}

	// Wait for registrations
	time.Sleep(50 * time.Millisecond)

	clients := hub.GetClients()
	assert.Len(t, clients, 5)
}

func TestHub_MaxConnections(t *testing.T) {
	config := &Config{
		MaxConnections: 2,
		WriteTimeout:   10 * time.Second,
		ReadTimeout:    60 * time.Second,
		PingInterval:   30 * time.Second,
		PongTimeout:    30 * time.Second,
	}
	hub := NewHub(config, nil)
	hub.Start()
	defer hub.Stop()

	// Register max clients
	for i := 0; i < 2; i++ {
		client := newTestClient(fmt.Sprintf("client-%d", i))
		hub.Register(client)
	}

	// Wait for registrations
	time.Sleep(50 * time.Millisecond)

	assert.Equal(t, 2, hub.GetClientCount())

	// Try to register one more (should be rejected)
	extraClient := newTestClient("extra-client")
	hub.Register(extraClient)

	// Wait for registration attempt
	time.Sleep(50 * time.Millisecond)

	// Should still be 2
	assert.Equal(t, 2, hub.GetClientCount())
}

func TestHub_Callbacks(t *testing.T) {
	hub := NewHub(nil, nil)

	var connectCalled, disconnectCalled bool
	var mu sync.Mutex

	hub.SetOnConnect(func(client *Client) {
		mu.Lock()
		connectCalled = true
		mu.Unlock()
	})

	hub.SetOnDisconnect(func(client *Client) {
		mu.Lock()
		disconnectCalled = true
		mu.Unlock()
	})

	hub.Start()
	defer hub.Stop()

	client := newTestClient("callback-client")

	hub.Register(client)
	time.Sleep(50 * time.Millisecond)

	mu.Lock()
	assert.True(t, connectCalled)
	mu.Unlock()

	hub.Unregister(client)
	time.Sleep(50 * time.Millisecond)

	mu.Lock()
	assert.True(t, disconnectCalled)
	mu.Unlock()
}

func TestHub_Config(t *testing.T) {
	config := &Config{
		MaxConnections: 500,
	}
	hub := NewHub(config, nil)

	assert.Equal(t, config, hub.Config())
}

func TestHub_ConcurrentAccess(t *testing.T) {
	hub := NewHub(nil, nil)
	hub.Start()
	defer hub.Stop()

	var wg sync.WaitGroup
	clientCount := 100

	// Concurrent registrations
	for i := 0; i < clientCount; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			client := newTestClient(fmt.Sprintf("concurrent-client-%d", id))
			hub.Register(client)
		}(i)
	}

	wg.Wait()
	time.Sleep(100 * time.Millisecond)

	// All clients should be registered
	assert.Equal(t, clientCount, hub.GetClientCount())

	// Concurrent broadcasts
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			hub.Broadcast(&Message{
				Type:  MessageTypeEvent,
				Event: "concurrent-event",
			})
		}(i)
	}

	wg.Wait()
}
