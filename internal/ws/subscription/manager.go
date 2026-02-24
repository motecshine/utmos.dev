// Package subscription provides subscription management for WebSocket clients
package subscription

import (
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

// deleteFromNestedMap removes innerKey from the inner map at key, and removes the outer entry if empty.
func deleteFromNestedMap[K comparable, V comparable](m map[K]map[V]bool, key K, innerKey V) {
	if m[key] != nil {
		delete(m[key], innerKey)
		if len(m[key]) == 0 {
			delete(m, key)
		}
	}
}

// Manager manages client subscriptions to topics
type Manager struct {
	// topic -> set of client IDs
	topics map[string]map[string]bool
	// client ID -> set of topics
	clients map[string]map[string]bool
	mu      sync.RWMutex
	logger  *logrus.Entry
}

// NewManager creates a new subscription manager
func NewManager(logger *logrus.Entry) *Manager {
	if logger == nil {
		logger = logrus.NewEntry(logrus.StandardLogger())
	}

	return &Manager{
		topics:  make(map[string]map[string]bool),
		clients: make(map[string]map[string]bool),
		logger:  logger.WithField("component", "subscription-manager"),
	}
}

// Subscribe subscribes a client to a topic
func (m *Manager) Subscribe(clientID, topic string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Add to topics map
	if m.topics[topic] == nil {
		m.topics[topic] = make(map[string]bool)
	}
	m.topics[topic][clientID] = true

	// Add to clients map
	if m.clients[clientID] == nil {
		m.clients[clientID] = make(map[string]bool)
	}
	m.clients[clientID][topic] = true

	m.logger.WithFields(logrus.Fields{
		"client_id": clientID,
		"topic":     topic,
	}).Debug("Client subscribed to topic")
}

// Unsubscribe unsubscribes a client from a topic
func (m *Manager) Unsubscribe(clientID, topic string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Remove from topics map
	deleteFromNestedMap(m.topics, topic, clientID)

	// Remove from clients map
	deleteFromNestedMap(m.clients, clientID, topic)

	m.logger.WithFields(logrus.Fields{
		"client_id": clientID,
		"topic":     topic,
	}).Debug("Client unsubscribed from topic")
}

// UnsubscribeAll unsubscribes a client from all topics
func (m *Manager) UnsubscribeAll(clientID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	topics, exists := m.clients[clientID]
	if !exists {
		return
	}

	// Remove client from all topics
	for topic := range topics {
		deleteFromNestedMap(m.topics, topic, clientID)
	}

	// Remove client entry
	delete(m.clients, clientID)

	m.logger.WithField("client_id", clientID).Debug("Client unsubscribed from all topics")
}

// getKeysFromMap looks up a key in a map[string]map[string]bool and returns the inner keys.
func (m *Manager) getKeysFromMap(lookup map[string]map[string]bool, key string) []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	innerMap, exists := lookup[key]
	if !exists {
		return nil
	}

	result := make([]string, 0, len(innerMap))
	for k := range innerMap {
		result = append(result, k)
	}
	return result
}

// GetSubscribers returns all client IDs subscribed to a topic
func (m *Manager) GetSubscribers(topic string) []string {
	return m.getKeysFromMap(m.topics, topic)
}

// GetSubscribersMatching returns all client IDs subscribed to topics matching a pattern
// Supports wildcard patterns:
// - "device.*" matches "device.telemetry", "device.status", etc.
// - "device.*.property" matches "device.drone1.property", "device.drone2.property", etc.
// - "*" matches any single segment
// - "**" matches any number of segments (greedy)
func (m *Manager) GetSubscribersMatching(pattern string) []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	clientSet := make(map[string]bool)

	for topic, clients := range m.topics {
		if matchTopic(pattern, topic) {
			for clientID := range clients {
				clientSet[clientID] = true
			}
		}
	}

	result := make([]string, 0, len(clientSet))
	for clientID := range clientSet {
		result = append(result, clientID)
	}
	return result
}

// GetTopics returns all topics a client is subscribed to
func (m *Manager) GetTopics(clientID string) []string {
	return m.getKeysFromMap(m.clients, clientID)
}

// IsSubscribed checks if a client is subscribed to a topic
func (m *Manager) IsSubscribed(clientID, topic string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	topics, exists := m.clients[clientID]
	if !exists {
		return false
	}
	return topics[topic]
}

// GetTopicCount returns the number of active topics
func (m *Manager) GetTopicCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.topics)
}

// GetClientCount returns the number of clients with subscriptions
func (m *Manager) GetClientCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.clients)
}

// GetSubscriberCount returns the number of subscribers for a topic
func (m *Manager) GetSubscriberCount(topic string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	clients, exists := m.topics[topic]
	if !exists {
		return 0
	}
	return len(clients)
}

// matchTopic checks if a topic matches a pattern
func matchTopic(pattern, topic string) bool {
	patternParts := strings.Split(pattern, ".")
	topicParts := strings.Split(topic, ".")

	return matchParts(patternParts, topicParts)
}

// matchParts recursively matches pattern parts against topic parts
func matchParts(pattern, topic []string) bool {
	if len(pattern) == 0 && len(topic) == 0 {
		return true
	}

	if len(pattern) == 0 {
		return false
	}

	if pattern[0] == "**" {
		// ** matches zero or more segments
		if len(pattern) == 1 {
			return true
		}
		// Try matching ** with 0, 1, 2, ... segments
		for i := 0; i <= len(topic); i++ {
			if matchParts(pattern[1:], topic[i:]) {
				return true
			}
		}
		return false
	}

	if len(topic) == 0 {
		return false
	}

	if pattern[0] == "*" || pattern[0] == topic[0] {
		return matchParts(pattern[1:], topic[1:])
	}

	return false
}

// TopicStats holds statistics for a topic
type TopicStats struct {
	Topic           string
	SubscriberCount int
}

// GetStats returns statistics for all topics
func (m *Manager) GetStats() []TopicStats {
	m.mu.RLock()
	defer m.mu.RUnlock()

	stats := make([]TopicStats, 0, len(m.topics))
	for topic, clients := range m.topics {
		stats = append(stats, TopicStats{
			Topic:           topic,
			SubscriberCount: len(clients),
		})
	}
	return stats
}
