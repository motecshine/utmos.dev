// Package push provides message pushing functionality for WebSocket clients
package push

import (
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/utmos/utmos/internal/ws/hub"
	"github.com/utmos/utmos/internal/ws/subscription"
)

// Pusher configuration defaults
const (
	// DefaultWorkerCount is the default number of worker goroutines
	DefaultWorkerCount = 4
	// DefaultQueueSize is the default size of the message queue
	DefaultQueueSize = 10000
)

// Config holds pusher configuration
type Config struct {
	// WorkerCount is the number of worker goroutines for pushing messages
	WorkerCount int
	// QueueSize is the size of the message queue
	QueueSize int
}

// DefaultConfig returns default pusher configuration
func DefaultConfig() *Config {
	return &Config{
		WorkerCount: DefaultWorkerCount,
		QueueSize:   DefaultQueueSize,
	}
}

// PushMessage represents a message to be pushed
type PushMessage struct {
	// Topic is the subscription topic
	Topic string
	// Message is the message to push
	Message *hub.Message
	// ClientIDs is an optional list of specific client IDs to push to
	// If empty, pushes to all subscribers of the topic
	ClientIDs []string
	// ExcludeClientIDs is a list of client IDs to exclude
	ExcludeClientIDs []string
}

// Pusher pushes messages to WebSocket clients
type Pusher struct {
	config      *Config
	hub         *hub.Hub
	subManager  *subscription.Manager
	logger      *logrus.Entry
	queue       chan *PushMessage
	done        chan struct{}
	wg          sync.WaitGroup
	running     bool
	runningMu   sync.RWMutex

	// Metrics
	messagesPushed   int64
	messagesDropped  int64
	metricsMu        sync.RWMutex
}

// NewPusher creates a new message pusher
func NewPusher(
	config *Config,
	h *hub.Hub,
	subManager *subscription.Manager,
	logger *logrus.Entry,
) *Pusher {
	if config == nil {
		config = DefaultConfig()
	}
	if logger == nil {
		logger = logrus.NewEntry(logrus.StandardLogger())
	}

	return &Pusher{
		config:     config,
		hub:        h,
		subManager: subManager,
		logger:     logger.WithField("component", "pusher"),
		queue:      make(chan *PushMessage, config.QueueSize),
		done:       make(chan struct{}),
	}
}

// Start starts the pusher workers
func (p *Pusher) Start() {
	p.runningMu.Lock()
	if p.running {
		p.runningMu.Unlock()
		return
	}
	p.running = true
	p.done = make(chan struct{})
	p.runningMu.Unlock()

	p.logger.WithField("workers", p.config.WorkerCount).Info("Starting pusher")

	for i := 0; i < p.config.WorkerCount; i++ {
		p.wg.Add(1)
		go p.worker(i)
	}
}

// Stop stops the pusher workers
func (p *Pusher) Stop() {
	p.runningMu.Lock()
	if !p.running {
		p.runningMu.Unlock()
		return
	}
	p.running = false
	p.runningMu.Unlock()

	close(p.done)
	p.wg.Wait()

	p.logger.Info("Pusher stopped")
}

// IsRunning returns whether the pusher is running
func (p *Pusher) IsRunning() bool {
	p.runningMu.RLock()
	defer p.runningMu.RUnlock()
	return p.running
}

// Push queues a message for pushing
func (p *Pusher) Push(msg *PushMessage) bool {
	if !p.IsRunning() {
		return false
	}

	select {
	case p.queue <- msg:
		return true
	default:
		p.incrementDropped()
		p.logger.WithField("topic", msg.Topic).Warn("Push queue full, dropping message")
		return false
	}
}

// PushToTopic pushes a message to all subscribers of a topic
func (p *Pusher) PushToTopic(topic string, msg *hub.Message) bool {
	return p.Push(&PushMessage{
		Topic:   topic,
		Message: msg,
	})
}

// PushToClients pushes a message to specific clients
func (p *Pusher) PushToClients(clientIDs []string, msg *hub.Message) bool {
	return p.Push(&PushMessage{
		ClientIDs: clientIDs,
		Message:   msg,
	})
}

// PushToTopicExcluding pushes a message to topic subscribers excluding specific clients
func (p *Pusher) PushToTopicExcluding(topic string, msg *hub.Message, excludeIDs []string) bool {
	return p.Push(&PushMessage{
		Topic:            topic,
		Message:          msg,
		ExcludeClientIDs: excludeIDs,
	})
}

// Broadcast pushes a message to all connected clients
func (p *Pusher) Broadcast(msg *hub.Message) {
	if p.hub != nil {
		p.hub.Broadcast(msg)
	}
}

// worker processes messages from the queue
func (p *Pusher) worker(id int) {
	defer p.wg.Done()

	logger := p.logger.WithField("worker_id", id)
	logger.Debug("Worker started")

	for {
		select {
		case <-p.done:
			logger.Debug("Worker stopping")
			return

		case msg := <-p.queue:
			p.processMessage(msg)
		}
	}
}

// processMessage processes a single push message
func (p *Pusher) processMessage(msg *PushMessage) {
	if msg == nil || msg.Message == nil {
		return
	}

	var clientIDs []string

	if len(msg.ClientIDs) > 0 {
		// Push to specific clients
		clientIDs = msg.ClientIDs
	} else if msg.Topic != "" && p.subManager != nil {
		// Get subscribers for topic
		clientIDs = p.subManager.GetSubscribers(msg.Topic)
	}

	if len(clientIDs) == 0 {
		return
	}

	// Build exclude set
	excludeSet := make(map[string]bool)
	for _, id := range msg.ExcludeClientIDs {
		excludeSet[id] = true
	}

	// Push to each client
	pushedCount := 0
	for _, clientID := range clientIDs {
		if excludeSet[clientID] {
			continue
		}

		if p.hub != nil {
			if p.hub.SendToClient(clientID, msg.Message) {
				pushedCount++
			}
		}
	}

	if pushedCount > 0 {
		p.incrementPushed(int64(pushedCount))
	}

	p.logger.WithFields(logrus.Fields{
		"topic":        msg.Topic,
		"pushed_count": pushedCount,
		"total_targets": len(clientIDs),
	}).Debug("Message pushed")
}

// incrementPushed increments the pushed counter
func (p *Pusher) incrementPushed(count int64) {
	p.metricsMu.Lock()
	p.messagesPushed += count
	p.metricsMu.Unlock()
}

// incrementDropped increments the dropped counter
func (p *Pusher) incrementDropped() {
	p.metricsMu.Lock()
	p.messagesDropped++
	p.metricsMu.Unlock()
}

// GetMetrics returns pusher metrics
func (p *Pusher) GetMetrics() (pushed, dropped int64) {
	p.metricsMu.RLock()
	defer p.metricsMu.RUnlock()
	return p.messagesPushed, p.messagesDropped
}

// QueueLength returns the current queue length
func (p *Pusher) QueueLength() int {
	return len(p.queue)
}
