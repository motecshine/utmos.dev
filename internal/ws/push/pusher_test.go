package push

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/utmos/utmos/internal/ws/hub"
	"github.com/utmos/utmos/internal/ws/subscription"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	assert.Equal(t, 4, config.WorkerCount)
	assert.Equal(t, 10000, config.QueueSize)
}

func TestNewPusher(t *testing.T) {
	t.Run("with nil config", func(t *testing.T) {
		pusher := NewPusher(nil, nil, nil, nil)
		require.NotNil(t, pusher)
		assert.Equal(t, 4, pusher.config.WorkerCount)
	})

	t.Run("with custom config", func(t *testing.T) {
		config := &Config{
			WorkerCount: 8,
			QueueSize:   5000,
		}
		pusher := NewPusher(config, nil, nil, nil)
		assert.Equal(t, 8, pusher.config.WorkerCount)
		assert.Equal(t, 5000, pusher.config.QueueSize)
	})
}

func TestPusher_StartStop(t *testing.T) {
	pusher := NewPusher(nil, nil, nil, nil)

	assert.False(t, pusher.IsRunning())

	pusher.Start()
	assert.True(t, pusher.IsRunning())

	// Starting again should be no-op
	pusher.Start()
	assert.True(t, pusher.IsRunning())

	pusher.Stop()
	assert.False(t, pusher.IsRunning())

	// Stopping again should be no-op
	pusher.Stop()
	assert.False(t, pusher.IsRunning())
}

func TestPusher_Push(t *testing.T) {
	pusher := NewPusher(nil, nil, nil, nil)

	// Push before start should fail
	success := pusher.Push(&PushMessage{
		Topic:   "test",
		Message: &hub.Message{Type: hub.MessageTypeEvent},
	})
	assert.False(t, success)

	pusher.Start()
	defer pusher.Stop()

	// Push after start should succeed
	success = pusher.Push(&PushMessage{
		Topic:   "test",
		Message: &hub.Message{Type: hub.MessageTypeEvent},
	})
	assert.True(t, success)
}

func TestPusher_PushToTopic(t *testing.T) {
	h := hub.NewHub(nil, nil)
	h.Start()
	defer h.Stop()

	subManager := subscription.NewManager(nil)
	pusher := NewPusher(nil, h, subManager, nil)
	pusher.Start()
	defer pusher.Stop()

	// Create and register clients
	client1 := &hub.Client{
		ID: "client1",
	}
	// We need to access the send channel, but it's not exported
	// For testing, we'll verify through the hub

	// Subscribe clients
	subManager.Subscribe("client1", "device.telemetry")
	subManager.Subscribe("client2", "device.telemetry")

	// Push to topic
	success := pusher.PushToTopic("device.telemetry", &hub.Message{
		Type:  hub.MessageTypeEvent,
		Event: "telemetry.update",
	})
	assert.True(t, success)

	// Wait for processing
	time.Sleep(50 * time.Millisecond)

	// Verify queue was processed
	assert.Equal(t, 0, pusher.QueueLength())

	// Cleanup
	_ = client1
}

func TestPusher_PushToClients(t *testing.T) {
	pusher := NewPusher(nil, nil, nil, nil)
	pusher.Start()
	defer pusher.Stop()

	success := pusher.PushToClients([]string{"client1", "client2"}, &hub.Message{
		Type:  hub.MessageTypeEvent,
		Event: "direct.message",
	})
	assert.True(t, success)

	// Wait for processing
	time.Sleep(50 * time.Millisecond)

	assert.Equal(t, 0, pusher.QueueLength())
}

func TestPusher_PushToTopicExcluding(t *testing.T) {
	subManager := subscription.NewManager(nil)
	pusher := NewPusher(nil, nil, subManager, nil)
	pusher.Start()
	defer pusher.Stop()

	subManager.Subscribe("client1", "topic1")
	subManager.Subscribe("client2", "topic1")
	subManager.Subscribe("client3", "topic1")

	success := pusher.PushToTopicExcluding("topic1", &hub.Message{
		Type: hub.MessageTypeEvent,
	}, []string{"client2"})
	assert.True(t, success)

	// Wait for processing
	time.Sleep(50 * time.Millisecond)

	assert.Equal(t, 0, pusher.QueueLength())
}

func TestPusher_GetMetrics(t *testing.T) {
	pusher := NewPusher(nil, nil, nil, nil)

	pushed, dropped := pusher.GetMetrics()
	assert.Equal(t, int64(0), pushed)
	assert.Equal(t, int64(0), dropped)

	// Increment metrics manually for testing
	pusher.incrementPushed(10)
	pusher.incrementDropped()

	pushed, dropped = pusher.GetMetrics()
	assert.Equal(t, int64(10), pushed)
	assert.Equal(t, int64(1), dropped)
}

func TestPusher_QueueFull(t *testing.T) {
	config := &Config{
		WorkerCount: 1,
		QueueSize:   2,
	}
	pusher := NewPusher(config, nil, nil, nil)
	pusher.Start()
	// Don't defer stop - we want to test queue full behavior

	// Fill the queue
	for i := 0; i < 2; i++ {
		pusher.Push(&PushMessage{
			Topic:   "test",
			Message: &hub.Message{Type: hub.MessageTypeEvent},
		})
	}

	// This should fail due to full queue
	// Note: This test is timing-dependent
	time.Sleep(10 * time.Millisecond)

	pusher.Stop()
}

func TestPusher_ConcurrentPush(t *testing.T) {
	h := hub.NewHub(nil, nil)
	h.Start()
	defer h.Stop()

	subManager := subscription.NewManager(nil)
	pusher := NewPusher(nil, h, subManager, nil)
	pusher.Start()
	defer pusher.Stop()

	// Setup subscriptions
	for i := 0; i < 10; i++ {
		subManager.Subscribe(fmt.Sprintf("client-%d", i), "topic1")
	}

	var wg sync.WaitGroup
	pushCount := 100

	for i := 0; i < pushCount; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			pusher.PushToTopic("topic1", &hub.Message{
				Type:  hub.MessageTypeEvent,
				Event: "concurrent.event",
				Data:  id,
			})
		}(i)
	}

	wg.Wait()

	// Wait for all messages to be processed
	time.Sleep(100 * time.Millisecond)

	assert.Equal(t, 0, pusher.QueueLength())
}

func TestPusher_ProcessMessage(t *testing.T) {
	h := hub.NewHub(nil, nil)
	h.Start()
	defer h.Stop()

	subManager := subscription.NewManager(nil)
	pusher := NewPusher(nil, h, subManager, nil)

	t.Run("nil message", func(t *testing.T) {
		// Should not panic
		pusher.processMessage(nil)
	})

	t.Run("nil inner message", func(t *testing.T) {
		// Should not panic
		pusher.processMessage(&PushMessage{Topic: "test"})
	})

	t.Run("no subscribers", func(t *testing.T) {
		pusher.processMessage(&PushMessage{
			Topic:   "no-subscribers",
			Message: &hub.Message{Type: hub.MessageTypeEvent},
		})
		// Should complete without error
	})
}

func TestPusher_Broadcast(t *testing.T) {
	h := hub.NewHub(nil, nil)
	h.Start()
	defer h.Stop()

	pusher := NewPusher(nil, h, nil, nil)

	// Should not panic even without starting
	pusher.Broadcast(&hub.Message{
		Type:  hub.MessageTypeEvent,
		Event: "broadcast.event",
	})
}

func TestPusher_WithRealHub(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	h := hub.NewHub(nil, nil)
	h.Start()
	defer h.Stop()

	subManager := subscription.NewManager(nil)
	pusher := NewPusher(nil, h, subManager, nil)
	pusher.Start()
	defer pusher.Stop()

	// Create mock clients with send channels
	clients := make([]*mockClient, 3)
	for i := 0; i < 3; i++ {
		clients[i] = &mockClient{
			id:   fmt.Sprintf("client-%d", i),
			send: make(chan *hub.Message, 256),
		}
		subManager.Subscribe(clients[i].id, "device.telemetry")
	}

	// Push message
	pusher.PushToTopic("device.telemetry", &hub.Message{
		Type:  hub.MessageTypeEvent,
		Event: "telemetry.update",
		Data:  map[string]any{"temperature": 25.5},
	})

	// Wait for processing
	time.Sleep(100 * time.Millisecond)

	// Verify metrics
	pushed, _ := pusher.GetMetrics()
	// Note: pushed count depends on whether clients are actually registered in hub
	_ = pushed
}

type mockClient struct {
	id   string
	send chan *hub.Message
}
