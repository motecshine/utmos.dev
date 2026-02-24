package hub

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

// Client configuration defaults
const (
	// DefaultClientWriteWait is the default time allowed to write a message
	DefaultClientWriteWait = 10 * time.Second
	// DefaultClientPongWait is the default time allowed to read the next pong message
	DefaultClientPongWait = 60 * time.Second
	// DefaultClientPingPeriod is the default period for sending ping messages (must be less than PongWait)
	DefaultClientPingPeriod = 54 * time.Second
	// DefaultClientMaxMessageSize is the default maximum message size (512KB)
	DefaultClientMaxMessageSize = 512 * 1024
	// DefaultClientSendBufferSize is the default size of the send channel buffer
	DefaultClientSendBufferSize = 256
)

// ClientConfig holds client configuration
type ClientConfig struct {
	// WriteWait is the time allowed to write a message to the peer
	WriteWait time.Duration
	// PongWait is the time allowed to read the next pong message from the peer
	PongWait time.Duration
	// PingPeriod is the period for sending ping messages (must be less than PongWait)
	PingPeriod time.Duration
	// MaxMessageSize is the maximum message size allowed from peer
	MaxMessageSize int64
	// SendBufferSize is the size of the send channel buffer
	SendBufferSize int
}

// DefaultClientConfig returns default client configuration
func DefaultClientConfig() *ClientConfig {
	return &ClientConfig{
		WriteWait:      DefaultClientWriteWait,
		PongWait:       DefaultClientPongWait,
		PingPeriod:     DefaultClientPingPeriod,
		MaxMessageSize: DefaultClientMaxMessageSize,
		SendBufferSize: DefaultClientSendBufferSize,
	}
}

// Client represents a WebSocket client connection
type Client struct {
	ID       string
	DeviceSN string
	UserID   string
	// Metadata stores client-specific key-value data
	Metadata Metadata

	hub    *Hub
	conn   *websocket.Conn
	config *ClientConfig
	logger *logrus.Entry
	send   chan *Message

	// Subscriptions
	subscriptions map[string]bool
	subMu         sync.RWMutex

	// State
	closed   bool
	closedMu sync.RWMutex
	done     chan struct{}

	// Callbacks
	onMessage func(client *Client, msg *Message)
}

// NewClient creates a new WebSocket client
func NewClient(id string, conn *websocket.Conn, hub *Hub, config *ClientConfig, logger *logrus.Entry) *Client {
	if config == nil {
		config = DefaultClientConfig()
	}
	if logger == nil {
		logger = logrus.NewEntry(logrus.StandardLogger())
	}

	return &Client{
		ID:            id,
		hub:           hub,
		conn:          conn,
		config:        config,
		logger:        logger.WithField("client_id", id),
		send:          make(chan *Message, config.SendBufferSize),
		subscriptions: make(map[string]bool),
		done:          make(chan struct{}),
		Metadata:      make(Metadata),
	}
}

// SetOnMessage sets the message callback
func (c *Client) SetOnMessage(callback func(client *Client, msg *Message)) {
	c.onMessage = callback
}

// Start starts the client read and write pumps
func (c *Client) Start() {
	go c.writePump()
	go c.readPump()
}

// Close closes the client connection
func (c *Client) Close() {
	c.closedMu.Lock()
	if c.closed {
		c.closedMu.Unlock()
		return
	}
	c.closed = true
	c.closedMu.Unlock()

	if c.done != nil {
		close(c.done)
	}

	if c.conn != nil {
		_ = c.conn.Close()
	}

	if c.logger != nil {
		c.logger.Debug("Client connection closed")
	}
}

// IsClosed returns whether the client is closed
func (c *Client) IsClosed() bool {
	c.closedMu.RLock()
	defer c.closedMu.RUnlock()
	return c.closed
}

// Send sends a message to the client
func (c *Client) Send(msg *Message) bool {
	if c.IsClosed() {
		return false
	}

	select {
	case c.send <- msg:
		return true
	default:
		c.logger.Warn("Send buffer full, dropping message")
		return false
	}
}

// SendJSON sends a JSON message to the client
func (c *Client) SendJSON(msgType MessageType, event string, data any) bool {
	return c.Send(&Message{
		Type:  msgType,
		Event: event,
		Data:  data,
	})
}

// SendError sends an error message to the client
func (c *Client) SendError(err string, traceID string) bool {
	return c.Send(&Message{
		Type:    MessageTypeError,
		Error:   err,
		TraceID: traceID,
	})
}

// Subscribe subscribes the client to a topic
func (c *Client) Subscribe(topic string) {
	c.subMu.Lock()
	defer c.subMu.Unlock()
	c.subscriptions[topic] = true
	c.logger.WithField("topic", topic).Debug("Subscribed to topic")
}

// Unsubscribe unsubscribes the client from a topic
func (c *Client) Unsubscribe(topic string) {
	c.subMu.Lock()
	defer c.subMu.Unlock()
	delete(c.subscriptions, topic)
	c.logger.WithField("topic", topic).Debug("Unsubscribed from topic")
}

// IsSubscribed checks if the client is subscribed to a topic
func (c *Client) IsSubscribed(topic string) bool {
	c.subMu.RLock()
	defer c.subMu.RUnlock()
	return c.subscriptions[topic]
}

// GetSubscriptions returns all subscriptions
func (c *Client) GetSubscriptions() []string {
	c.subMu.RLock()
	defer c.subMu.RUnlock()

	subs := make([]string, 0, len(c.subscriptions))
	for topic := range c.subscriptions {
		subs = append(subs, topic)
	}
	return subs
}

// readPump pumps messages from the WebSocket connection to the hub
func (c *Client) readPump() {
	defer func() {
		if c.hub != nil {
			c.hub.Unregister(c)
		}
		c.Close()
	}()

	if c.conn == nil {
		return
	}

	c.conn.SetReadLimit(c.config.MaxMessageSize)
	_ = c.conn.SetReadDeadline(time.Now().Add(c.config.PongWait))
	c.conn.SetPongHandler(func(string) error {
		_ = c.conn.SetReadDeadline(time.Now().Add(c.config.PongWait))
		return nil
	})

	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.logger.WithError(err).Warn("WebSocket read error")
			}
			return
		}

		var msg Message
		if err := json.Unmarshal(data, &msg); err != nil {
			c.logger.WithError(err).Warn("Failed to unmarshal message")
			c.SendError("invalid message format", "")
			continue
		}

		c.handleMessage(&msg)
	}
}

// writePump pumps messages from the hub to the WebSocket connection
func (c *Client) writePump() {
	ticker := time.NewTicker(c.config.PingPeriod)
	defer func() {
		ticker.Stop()
		c.Close()
	}()

	for {
		select {
		case <-c.done:
			return

		case msg, ok := <-c.send:
			if c.conn == nil {
				return
			}

			if err := c.setWriteDeadline("Failed to set write deadline"); err != nil {
				return
			}
			if !ok {
				// Channel closed
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			data, err := json.Marshal(msg)
			if err != nil {
				c.logger.WithError(err).Error("Failed to marshal message")
				continue
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
				c.logger.WithError(err).Warn("WebSocket write error")
				return
			}

		case <-ticker.C:
			if c.conn == nil {
				return
			}

			if err := c.setWriteDeadline("Failed to set write deadline for ping"); err != nil {
				return
			}
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				c.logger.WithError(err).Warn("Failed to send ping")
				return
			}
		}
	}
}

// handleSubUnsubMessage handles subscribe/unsubscribe messages by executing the action and sending an ack.
func (c *Client) handleSubUnsubMessage(msg *Message, action func(string)) {
	if msg.Event != "" {
		action(msg.Event)
		c.Send(&Message{
			Type:  MessageTypeAck,
			Event: msg.Event,
		})
		if c.hub != nil && c.hub.onMessage != nil {
			c.hub.onMessage(c, msg)
		}
	}
}

// handleMessage handles incoming messages
func (c *Client) handleMessage(msg *Message) {
	switch msg.Type {
	case MessageTypePing:
		c.Send(&Message{Type: MessageTypePong})

	case MessageTypeSubscribe:
		c.handleSubUnsubMessage(msg, c.Subscribe)

	case MessageTypeUnsubscribe:
		c.handleSubUnsubMessage(msg, c.Unsubscribe)

	default:
		// Forward to callback
		if c.onMessage != nil {
			c.onMessage(c, msg)
		}
		// Also forward to hub callback
		if c.hub != nil && c.hub.onMessage != nil {
			c.hub.onMessage(c, msg)
		}
	}
}

// setWriteDeadline sets the write deadline on the connection and logs a warning on failure.
func (c *Client) setWriteDeadline(msg string) error {
	if err := c.conn.SetWriteDeadline(time.Now().Add(c.config.WriteWait)); err != nil {
		c.logger.WithError(err).Warn(msg)
		return err
	}
	return nil
}

// Config returns the client configuration
func (c *Client) Config() *ClientConfig {
	return c.config
}
