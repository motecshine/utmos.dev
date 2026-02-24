// Package hub provides WebSocket connection management
package hub

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// Default configuration values
const (
	// DefaultMaxConnections is the default maximum number of WebSocket connections
	DefaultMaxConnections = 10000
	// DefaultWriteTimeout is the default timeout for writing to a connection
	DefaultWriteTimeout = 10 * time.Second
	// DefaultReadTimeout is the default timeout for reading from a connection
	DefaultReadTimeout = 60 * time.Second
	// DefaultPingInterval is the default interval for sending ping messages
	DefaultPingInterval = 30 * time.Second
	// DefaultPongTimeout is the default timeout for receiving pong responses
	DefaultPongTimeout = 30 * time.Second
	// DefaultBroadcastBufferSize is the default buffer size for broadcast channel
	DefaultBroadcastBufferSize = 1024
	// DefaultRegisterBufferSize is the default buffer size for register/unregister channels
	DefaultRegisterBufferSize = 256
)

// JSONData represents arbitrary JSON data that can be sent/received via WebSocket.
// Using json.RawMessage allows for type-safe handling while deferring parsing.
type JSONData = json.RawMessage

// Metadata represents client metadata as key-value pairs.
// Values are stored as strings for type safety.
type Metadata map[string]string

// Config holds hub configuration
type Config struct {
	// MaxConnections is the maximum number of connections allowed
	MaxConnections int
	// WriteTimeout is the timeout for writing to a connection
	WriteTimeout time.Duration
	// ReadTimeout is the timeout for reading from a connection
	ReadTimeout time.Duration
	// PingInterval is the interval for sending ping messages
	PingInterval time.Duration
	// PongTimeout is the timeout for receiving pong responses
	PongTimeout time.Duration
}

// DefaultConfig returns default hub configuration
func DefaultConfig() *Config {
	return &Config{
		MaxConnections: DefaultMaxConnections,
		WriteTimeout:   DefaultWriteTimeout,
		ReadTimeout:    DefaultReadTimeout,
		PingInterval:   DefaultPingInterval,
		PongTimeout:    DefaultPongTimeout,
	}
}

// Hub manages WebSocket connections
type Hub struct {
	config  *Config
	logger  *logrus.Entry
	clients map[string]*Client
	mu      sync.RWMutex

	// Channels for client management
	register   chan *Client
	unregister chan *Client
	broadcast  chan *Message

	// Callbacks
	onConnect    func(client *Client)
	onDisconnect func(client *Client)
	onMessage    func(client *Client, msg *Message)

	// State
	running bool
	done    chan struct{}
}

// Message represents a WebSocket message
type Message struct {
	Type    MessageType `json:"type"`
	Event   string      `json:"event,omitempty"`
	// Data holds arbitrary JSON data. Using any is intentional here
	// as WebSocket messages can contain any valid JSON structure.
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
	TraceID string `json:"trace_id,omitempty"`
}

// MessageType represents the type of WebSocket message
type MessageType string

const (
	MessageTypeEvent     MessageType = "event"
	MessageTypeSubscribe MessageType = "subscribe"
	MessageTypeUnsubscribe MessageType = "unsubscribe"
	MessageTypePing      MessageType = "ping"
	MessageTypePong      MessageType = "pong"
	MessageTypeError     MessageType = "error"
	MessageTypeAck       MessageType = "ack"
)

// NewHub creates a new WebSocket hub
func NewHub(config *Config, logger *logrus.Entry) *Hub {
	if config == nil {
		config = DefaultConfig()
	}
	if logger == nil {
		logger = logrus.NewEntry(logrus.StandardLogger())
	}

	return &Hub{
		config:     config,
		logger:     logger.WithField("component", "ws-hub"),
		clients:    make(map[string]*Client),
		register:   make(chan *Client, DefaultRegisterBufferSize),
		unregister: make(chan *Client, DefaultRegisterBufferSize),
		broadcast:  make(chan *Message, DefaultBroadcastBufferSize),
		done:       make(chan struct{}),
	}
}

// SetOnConnect sets the callback for new connections
func (h *Hub) SetOnConnect(callback func(client *Client)) {
	h.onConnect = callback
}

// SetOnDisconnect sets the callback for disconnections
func (h *Hub) SetOnDisconnect(callback func(client *Client)) {
	h.onDisconnect = callback
}

// SetOnMessage sets the callback for incoming messages
func (h *Hub) SetOnMessage(callback func(client *Client, msg *Message)) {
	h.onMessage = callback
}

// Start starts the hub
func (h *Hub) Start() {
	h.mu.Lock()
	if h.running {
		h.mu.Unlock()
		return
	}
	h.running = true
	h.done = make(chan struct{})
	h.mu.Unlock()

	h.logger.Info("Starting WebSocket hub")

	go h.run()
}

// Stop stops the hub
func (h *Hub) Stop() {
	h.mu.Lock()
	if !h.running {
		h.mu.Unlock()
		return
	}
	h.running = false
	h.mu.Unlock()

	close(h.done)

	// Close all client connections
	h.mu.Lock()
	for _, client := range h.clients {
		client.Close()
	}
	h.clients = make(map[string]*Client)
	h.mu.Unlock()

	h.logger.Info("WebSocket hub stopped")
}

// run is the main hub loop
func (h *Hub) run() {
	for {
		select {
		case <-h.done:
			return

		case client := <-h.register:
			h.handleRegister(client)

		case client := <-h.unregister:
			h.handleUnregister(client)

		case msg := <-h.broadcast:
			h.handleBroadcast(msg)
		}
	}
}

// handleRegister handles client registration
func (h *Hub) handleRegister(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Check max connections
	if len(h.clients) >= h.config.MaxConnections {
		h.logger.WithField("client_id", client.ID).Warn("Max connections reached, rejecting client")
		client.Close()
		return
	}

	h.clients[client.ID] = client

	h.logger.WithFields(logrus.Fields{
		"client_id":   client.ID,
		"total_count": len(h.clients),
	}).Debug("Client registered")

	if h.onConnect != nil {
		go h.onConnect(client)
	}
}

// handleUnregister handles client unregistration
func (h *Hub) handleUnregister(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, exists := h.clients[client.ID]; exists {
		delete(h.clients, client.ID)
		client.Close()

		h.logger.WithFields(logrus.Fields{
			"client_id":   client.ID,
			"total_count": len(h.clients),
		}).Debug("Client unregistered")

		if h.onDisconnect != nil {
			go h.onDisconnect(client)
		}
	}
}

// handleBroadcast handles broadcasting messages to all clients
func (h *Hub) handleBroadcast(msg *Message) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, client := range h.clients {
		select {
		case client.send <- msg:
		default:
			// Client send buffer full, skip
			h.logger.WithField("client_id", client.ID).Warn("Client send buffer full")
		}
	}
}

// Register registers a client with the hub
func (h *Hub) Register(client *Client) {
	select {
	case h.register <- client:
	default:
		h.logger.Warn("Register channel full")
	}
}

// Unregister unregisters a client from the hub
func (h *Hub) Unregister(client *Client) {
	select {
	case h.unregister <- client:
	default:
		h.logger.Warn("Unregister channel full")
	}
}

// Broadcast broadcasts a message to all clients
func (h *Hub) Broadcast(msg *Message) {
	select {
	case h.broadcast <- msg:
	default:
		h.logger.Warn("Broadcast channel full")
	}
}

// SendToClient sends a message to a specific client
func (h *Hub) SendToClient(clientID string, msg *Message) bool {
	h.mu.RLock()
	client, exists := h.clients[clientID]
	h.mu.RUnlock()

	if !exists {
		return false
	}

	select {
	case client.send <- msg:
		return true
	default:
		return false
	}
}

// GetClient returns a client by ID
func (h *Hub) GetClient(clientID string) (*Client, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	client, exists := h.clients[clientID]
	return client, exists
}

// GetClientCount returns the number of connected clients
func (h *Hub) GetClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

// GetClients returns all connected clients
func (h *Hub) GetClients() []*Client {
	h.mu.RLock()
	defer h.mu.RUnlock()

	clients := make([]*Client, 0, len(h.clients))
	for _, client := range h.clients {
		clients = append(clients, client)
	}
	return clients
}

// IsRunning returns whether the hub is running
func (h *Hub) IsRunning() bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.running
}

// Config returns the hub configuration
func (h *Hub) Config() *Config {
	return h.config
}
