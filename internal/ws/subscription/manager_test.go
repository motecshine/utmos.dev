package subscription

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewManager(t *testing.T) {
	manager := NewManager(nil)
	require.NotNil(t, manager)
	assert.NotNil(t, manager.topics)
	assert.NotNil(t, manager.clients)
}

func TestManager_SubscribeUnsubscribe(t *testing.T) {
	manager := NewManager(nil)

	// Subscribe
	manager.Subscribe("client1", "topic1")
	manager.Subscribe("client1", "topic2")
	manager.Subscribe("client2", "topic1")

	// Verify subscriptions
	assert.True(t, manager.IsSubscribed("client1", "topic1"))
	assert.True(t, manager.IsSubscribed("client1", "topic2"))
	assert.True(t, manager.IsSubscribed("client2", "topic1"))
	assert.False(t, manager.IsSubscribed("client2", "topic2"))

	// Unsubscribe
	manager.Unsubscribe("client1", "topic1")
	assert.False(t, manager.IsSubscribed("client1", "topic1"))
	assert.True(t, manager.IsSubscribed("client1", "topic2"))
}

func TestManager_UnsubscribeAll(t *testing.T) {
	manager := NewManager(nil)

	manager.Subscribe("client1", "topic1")
	manager.Subscribe("client1", "topic2")
	manager.Subscribe("client1", "topic3")
	manager.Subscribe("client2", "topic1")

	manager.UnsubscribeAll("client1")

	assert.False(t, manager.IsSubscribed("client1", "topic1"))
	assert.False(t, manager.IsSubscribed("client1", "topic2"))
	assert.False(t, manager.IsSubscribed("client1", "topic3"))
	assert.True(t, manager.IsSubscribed("client2", "topic1"))

	// Unsubscribe non-existing client should not panic
	manager.UnsubscribeAll("non-existing")
}

func TestManager_GetSubscribers(t *testing.T) {
	manager := NewManager(nil)

	manager.Subscribe("client1", "topic1")
	manager.Subscribe("client2", "topic1")
	manager.Subscribe("client3", "topic1")
	manager.Subscribe("client1", "topic2")

	subscribers := manager.GetSubscribers("topic1")
	assert.Len(t, subscribers, 3)
	assert.Contains(t, subscribers, "client1")
	assert.Contains(t, subscribers, "client2")
	assert.Contains(t, subscribers, "client3")

	subscribers = manager.GetSubscribers("topic2")
	assert.Len(t, subscribers, 1)
	assert.Contains(t, subscribers, "client1")

	subscribers = manager.GetSubscribers("non-existing")
	assert.Nil(t, subscribers)
}

func TestManager_GetTopics(t *testing.T) {
	manager := NewManager(nil)

	manager.Subscribe("client1", "topic1")
	manager.Subscribe("client1", "topic2")
	manager.Subscribe("client1", "topic3")

	topics := manager.GetTopics("client1")
	assert.Len(t, topics, 3)
	assert.Contains(t, topics, "topic1")
	assert.Contains(t, topics, "topic2")
	assert.Contains(t, topics, "topic3")

	topics = manager.GetTopics("non-existing")
	assert.Nil(t, topics)
}

func TestManager_Counts(t *testing.T) {
	manager := NewManager(nil)

	assert.Equal(t, 0, manager.GetTopicCount())
	assert.Equal(t, 0, manager.GetClientCount())
	assert.Equal(t, 0, manager.GetSubscriberCount("topic1"))

	manager.Subscribe("client1", "topic1")
	manager.Subscribe("client2", "topic1")
	manager.Subscribe("client1", "topic2")

	assert.Equal(t, 2, manager.GetTopicCount())
	assert.Equal(t, 2, manager.GetClientCount())
	assert.Equal(t, 2, manager.GetSubscriberCount("topic1"))
	assert.Equal(t, 1, manager.GetSubscriberCount("topic2"))
}

func TestManager_GetStats(t *testing.T) {
	manager := NewManager(nil)

	manager.Subscribe("client1", "topic1")
	manager.Subscribe("client2", "topic1")
	manager.Subscribe("client1", "topic2")

	stats := manager.GetStats()
	assert.Len(t, stats, 2)

	statsMap := make(map[string]int)
	for _, s := range stats {
		statsMap[s.Topic] = s.SubscriberCount
	}

	assert.Equal(t, 2, statsMap["topic1"])
	assert.Equal(t, 1, statsMap["topic2"])
}

func TestManager_GetSubscribersMatching(t *testing.T) {
	manager := NewManager(nil)

	// Setup subscriptions
	manager.Subscribe("client1", "device.telemetry")
	manager.Subscribe("client2", "device.status")
	manager.Subscribe("client3", "device.drone1.property")
	manager.Subscribe("client4", "device.drone2.property")
	manager.Subscribe("client5", "system.health")

	t.Run("exact match", func(t *testing.T) {
		subscribers := manager.GetSubscribersMatching("device.telemetry")
		assert.Len(t, subscribers, 1)
		assert.Contains(t, subscribers, "client1")
	})

	t.Run("single wildcard", func(t *testing.T) {
		subscribers := manager.GetSubscribersMatching("device.*")
		assert.Len(t, subscribers, 2)
		assert.Contains(t, subscribers, "client1")
		assert.Contains(t, subscribers, "client2")
	})

	t.Run("middle wildcard", func(t *testing.T) {
		subscribers := manager.GetSubscribersMatching("device.*.property")
		assert.Len(t, subscribers, 2)
		assert.Contains(t, subscribers, "client3")
		assert.Contains(t, subscribers, "client4")
	})

	t.Run("double wildcard", func(t *testing.T) {
		subscribers := manager.GetSubscribersMatching("device.**")
		assert.Len(t, subscribers, 4)
		assert.Contains(t, subscribers, "client1")
		assert.Contains(t, subscribers, "client2")
		assert.Contains(t, subscribers, "client3")
		assert.Contains(t, subscribers, "client4")
	})

	t.Run("no match", func(t *testing.T) {
		subscribers := manager.GetSubscribersMatching("unknown.*")
		assert.Empty(t, subscribers)
	})
}

func TestMatchTopic(t *testing.T) {
	tests := []struct {
		pattern string
		topic   string
		match   bool
	}{
		// Exact matches
		{"device.telemetry", "device.telemetry", true},
		{"device.telemetry", "device.status", false},

		// Single wildcard
		{"device.*", "device.telemetry", true},
		{"device.*", "device.status", true},
		{"device.*", "device.drone1.property", false},
		{"*.telemetry", "device.telemetry", true},
		{"*.telemetry", "system.telemetry", true},

		// Middle wildcard
		{"device.*.property", "device.drone1.property", true},
		{"device.*.property", "device.drone2.property", true},
		{"device.*.property", "device.telemetry", false},

		// Double wildcard
		{"device.**", "device.telemetry", true},
		{"device.**", "device.drone1.property", true},
		{"device.**", "device.drone1.property.value", true},
		{"**", "anything.goes.here", true},
		{"**.property", "device.drone1.property", true},
		{"device.**.value", "device.drone1.property.value", true},

		// Edge cases
		{"", "", true},
		{"*", "single", true},
		{"*", "multi.segment", false},
	}

	for _, tt := range tests {
		t.Run(tt.pattern+"_"+tt.topic, func(t *testing.T) {
			result := matchTopic(tt.pattern, tt.topic)
			assert.Equal(t, tt.match, result, "pattern=%s topic=%s", tt.pattern, tt.topic)
		})
	}
}

func TestManager_ConcurrentAccess(t *testing.T) {
	manager := NewManager(nil)

	var wg sync.WaitGroup
	clientCount := 100
	topicCount := 10

	// Concurrent subscriptions
	for i := 0; i < clientCount; i++ {
		wg.Add(1)
		go func(clientID int) {
			defer wg.Done()
			for j := 0; j < topicCount; j++ {
				manager.Subscribe(
					fmt.Sprintf("client-%d", clientID),
					fmt.Sprintf("topic-%d", j),
				)
			}
		}(i)
	}

	wg.Wait()

	// Concurrent reads
	for i := 0; i < clientCount; i++ {
		wg.Add(1)
		go func(clientID int) {
			defer wg.Done()
			manager.GetTopics(fmt.Sprintf("client-%d", clientID))
			manager.GetSubscribers("topic-0")
			manager.GetStats()
		}(i)
	}

	wg.Wait()

	// Concurrent unsubscriptions
	for i := 0; i < clientCount; i++ {
		wg.Add(1)
		go func(clientID int) {
			defer wg.Done()
			manager.UnsubscribeAll(fmt.Sprintf("client-%d", clientID))
		}(i)
	}

	wg.Wait()

	assert.Equal(t, 0, manager.GetClientCount())
}

func TestManager_CleanupEmptyTopics(t *testing.T) {
	manager := NewManager(nil)

	manager.Subscribe("client1", "topic1")
	assert.Equal(t, 1, manager.GetTopicCount())

	manager.Unsubscribe("client1", "topic1")
	assert.Equal(t, 0, manager.GetTopicCount())
}
