// Package mqtt provides MQTT client functionality for iot-gateway
package mqtt

import (
	"context"
	"fmt"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
)

// MQTT client configuration defaults
const (
	// DefaultBroker is the default MQTT broker address
	DefaultBroker = "localhost"
	// DefaultPort is the default MQTT broker port
	DefaultPort = 1883
	// DefaultClientID is the default client identifier
	DefaultClientID = "iot-gateway"
	// DefaultConnectTimeout is the default connection timeout
	DefaultConnectTimeout = 30 * time.Second
	// DefaultKeepAlive is the default keep-alive interval
	DefaultKeepAlive = 60 * time.Second
	// DefaultPingTimeout is the default ping timeout
	DefaultPingTimeout = 10 * time.Second
	// DefaultMaxReconnectWait is the default maximum reconnect wait time
	DefaultMaxReconnectWait = 5 * time.Minute
	// DefaultQoS is the default quality of service level
	DefaultQoS = 1
)

// Config holds MQTT client configuration
type Config struct {
	Broker           string
	Port             int
	ClientID         string
	Username         string
	Password         string
	CleanSession     bool
	AutoReconnect    bool
	ConnectTimeout   time.Duration
	KeepAlive        time.Duration
	PingTimeout      time.Duration
	MaxReconnectWait time.Duration
	QoS              byte
}

// DefaultConfig returns default MQTT client configuration
func DefaultConfig() *Config {
	return &Config{
		Broker:           DefaultBroker,
		Port:             DefaultPort,
		ClientID:         DefaultClientID,
		CleanSession:     false,
		AutoReconnect:    true,
		ConnectTimeout:   DefaultConnectTimeout,
		KeepAlive:        DefaultKeepAlive,
		PingTimeout:      DefaultPingTimeout,
		MaxReconnectWait: DefaultMaxReconnectWait,
		QoS:              DefaultQoS,
	}
}

// Client wraps the MQTT client with additional functionality
type Client struct {
	config         *Config
	client         mqtt.Client
	logger         *logrus.Entry
	messageHandler MessageHandler
	connectHandler ConnectHandler
	lostHandler    ConnectionLostHandler
	mu             sync.RWMutex
	connected      bool
	subscriptions  map[string]byte
}

// MessageHandler handles incoming MQTT messages
type MessageHandler func(client *Client, msg mqtt.Message)

// ConnectHandler handles connection events
type ConnectHandler func(client *Client)

// ConnectionLostHandler handles connection lost events
type ConnectionLostHandler func(client *Client, err error)

// NewClient creates a new MQTT client
func NewClient(config *Config, logger *logrus.Entry) *Client {
	if config == nil {
		config = DefaultConfig()
	}
	if logger == nil {
		logger = logrus.NewEntry(logrus.StandardLogger())
	}

	return &Client{
		config:        config,
		logger:        logger.WithField("component", "mqtt-client"),
		subscriptions: make(map[string]byte),
	}
}

// SetMessageHandler sets the message handler
func (c *Client) SetMessageHandler(handler MessageHandler) {
	c.messageHandler = handler
}

// SetConnectHandler sets the connect handler
func (c *Client) SetConnectHandler(handler ConnectHandler) {
	c.connectHandler = handler
}

// SetConnectionLostHandler sets the connection lost handler
func (c *Client) SetConnectionLostHandler(handler ConnectionLostHandler) {
	c.lostHandler = handler
}

// Connect establishes connection to the MQTT broker
func (c *Client) Connect(ctx context.Context) error {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", c.config.Broker, c.config.Port))
	opts.SetClientID(c.config.ClientID)
	opts.SetCleanSession(c.config.CleanSession)
	opts.SetAutoReconnect(c.config.AutoReconnect)
	opts.SetConnectTimeout(c.config.ConnectTimeout)
	opts.SetKeepAlive(c.config.KeepAlive)
	opts.SetPingTimeout(c.config.PingTimeout)
	opts.SetMaxReconnectInterval(c.config.MaxReconnectWait)

	if c.config.Username != "" {
		opts.SetUsername(c.config.Username)
		opts.SetPassword(c.config.Password)
	}

	// Set default message handler
	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		if c.messageHandler != nil {
			c.messageHandler(c, msg)
		}
	})

	// Set connection handler
	opts.SetOnConnectHandler(func(client mqtt.Client) {
		c.mu.Lock()
		c.connected = true
		c.mu.Unlock()

		c.logger.Info("Connected to MQTT broker")

		// Resubscribe to topics after reconnection
		c.resubscribe()

		if c.connectHandler != nil {
			c.connectHandler(c)
		}
	})

	// Set connection lost handler
	opts.SetConnectionLostHandler(func(client mqtt.Client, err error) {
		c.mu.Lock()
		c.connected = false
		c.mu.Unlock()

		c.logger.WithError(err).Warn("Connection to MQTT broker lost")

		if c.lostHandler != nil {
			c.lostHandler(c, err)
		}
	})

	c.client = mqtt.NewClient(opts)

	token := c.client.Connect()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-token.Done():
		if token.Error() != nil {
			return fmt.Errorf("failed to connect to MQTT broker: %w", token.Error())
		}
	}

	c.logger.WithFields(logrus.Fields{
		"broker":   c.config.Broker,
		"port":     c.config.Port,
		"clientID": c.config.ClientID,
	}).Info("MQTT client connected")

	return nil
}

// Disconnect closes the MQTT connection
func (c *Client) Disconnect(quiesce uint) {
	if c.client != nil && c.client.IsConnected() {
		c.client.Disconnect(quiesce)
		c.mu.Lock()
		c.connected = false
		c.mu.Unlock()
		c.logger.Info("MQTT client disconnected")
	}
}

// Subscribe subscribes to a topic
func (c *Client) Subscribe(topic string, qos byte, handler mqtt.MessageHandler) error {
	if c.client == nil || !c.client.IsConnected() {
		return fmt.Errorf("client not connected")
	}

	token := c.client.Subscribe(topic, qos, handler)
	token.Wait()
	if token.Error() != nil {
		return fmt.Errorf("failed to subscribe to topic %s: %w", topic, token.Error())
	}

	c.mu.Lock()
	c.subscriptions[topic] = qos
	c.mu.Unlock()

	c.logger.WithFields(logrus.Fields{
		"topic": topic,
		"qos":   qos,
	}).Debug("Subscribed to topic")

	return nil
}

// Unsubscribe unsubscribes from a topic
func (c *Client) Unsubscribe(topics ...string) error {
	if c.client == nil || !c.client.IsConnected() {
		return fmt.Errorf("client not connected")
	}

	token := c.client.Unsubscribe(topics...)
	token.Wait()
	if token.Error() != nil {
		return fmt.Errorf("failed to unsubscribe from topics: %w", token.Error())
	}

	c.mu.Lock()
	for _, topic := range topics {
		delete(c.subscriptions, topic)
	}
	c.mu.Unlock()

	return nil
}

// Publish publishes a message to a topic
func (c *Client) Publish(topic string, qos byte, retained bool, payload any) error {
	if c.client == nil || !c.client.IsConnected() {
		return fmt.Errorf("client not connected")
	}

	token := c.client.Publish(topic, qos, retained, payload)
	token.Wait()
	if token.Error() != nil {
		return fmt.Errorf("failed to publish to topic %s: %w", topic, token.Error())
	}

	c.logger.WithFields(logrus.Fields{
		"topic":    topic,
		"qos":      qos,
		"retained": retained,
	}).Debug("Published message")

	return nil
}

// IsConnected returns the connection status
func (c *Client) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.connected && c.client != nil && c.client.IsConnected()
}

// resubscribe resubscribes to all topics after reconnection
func (c *Client) resubscribe() {
	c.mu.RLock()
	subs := make(map[string]byte, len(c.subscriptions))
	for topic, qos := range c.subscriptions {
		subs[topic] = qos
	}
	c.mu.RUnlock()

	for topic, qos := range subs {
		token := c.client.Subscribe(topic, qos, nil)
		token.Wait()
		if token.Error() != nil {
			c.logger.WithError(token.Error()).WithField("topic", topic).Error("Failed to resubscribe")
		} else {
			c.logger.WithField("topic", topic).Debug("Resubscribed to topic")
		}
	}
}

// GetConfig returns the client configuration
func (c *Client) GetConfig() *Config {
	return c.config
}
