// Package bridge provides RabbitMQ to MQTT bridging functionality
package bridge

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/utmos/utmos/internal/gateway/mqtt"
	"github.com/utmos/utmos/pkg/rabbitmq"
)

// RawDownlinkMessage represents a raw downlink message to MQTT
type RawDownlinkMessage struct {
	DeviceSN string          `json:"device_sn"`
	Topic    string          `json:"topic"`
	Payload  json.RawMessage `json:"payload"`
	QoS      int             `json:"qos"`
	Retained bool            `json:"retained"`
	TraceID  string          `json:"trace_id"`
	SpanID   string          `json:"span_id"`
}

// DownlinkBridge bridges RabbitMQ messages to MQTT
type DownlinkBridge struct {
	mqttClient *mqtt.Client
	subscriber *rabbitmq.Subscriber
	logger     *logrus.Entry
	exchange   string
	queue      string
	routingKey string
	running    bool
	mu         sync.RWMutex
	stopCh     chan struct{}
}

// DownlinkBridgeConfig holds configuration for downlink bridge
type DownlinkBridgeConfig struct {
	Exchange   string
	Queue      string
	RoutingKey string
}

// DefaultDownlinkBridgeConfig returns default configuration
func DefaultDownlinkBridgeConfig() *DownlinkBridgeConfig {
	return &DownlinkBridgeConfig{
		Exchange:   "iot.topic",
		Queue:      "iot.gateway.downlink",
		RoutingKey: "iot.raw.*.downlink",
	}
}

// NewDownlinkBridge creates a new downlink bridge
func NewDownlinkBridge(mqttClient *mqtt.Client, subscriber *rabbitmq.Subscriber, config *DownlinkBridgeConfig, logger *logrus.Entry) *DownlinkBridge {
	if config == nil {
		config = DefaultDownlinkBridgeConfig()
	}
	if logger == nil {
		logger = logrus.NewEntry(logrus.StandardLogger())
	}

	return &DownlinkBridge{
		mqttClient: mqttClient,
		subscriber: subscriber,
		logger:     logger.WithField("component", "downlink-bridge"),
		exchange:   config.Exchange,
		queue:      config.Queue,
		routingKey: config.RoutingKey,
		stopCh:     make(chan struct{}),
	}
}

// Start starts consuming messages from RabbitMQ and forwarding to MQTT
func (b *DownlinkBridge) Start(ctx context.Context) error {
	b.mu.Lock()
	if b.running {
		b.mu.Unlock()
		return fmt.Errorf("bridge already running")
	}
	b.running = true
	b.mu.Unlock()

	if b.subscriber == nil {
		return fmt.Errorf("subscriber not initialized")
	}

	b.logger.Info("Starting downlink bridge")

	// Subscribe to downlink queue
	err := b.subscriber.Subscribe(b.queue, func(ctx context.Context, msg *rabbitmq.StandardMessage) error {
		return b.handleStandardMessage(ctx, msg)
	})
	if err != nil {
		b.mu.Lock()
		b.running = false
		b.mu.Unlock()
		return fmt.Errorf("failed to subscribe: %w", err)
	}

	return nil
}

// handleStandardMessage handles a StandardMessage from RabbitMQ
func (b *DownlinkBridge) handleStandardMessage(ctx context.Context, msg *rabbitmq.StandardMessage) error {
	// Extract downlink data from StandardMessage
	var dataMap map[string]any
	if err := msg.GetData(&dataMap); err != nil {
		return fmt.Errorf("failed to unmarshal data: %w", err)
	}

	topic, ok := dataMap["topic"].(string)
	if !ok {
		return fmt.Errorf("missing topic in message data")
	}

	payload, err := json.Marshal(dataMap["payload"])
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	qos := 1
	if q, ok := dataMap["qos"].(float64); ok {
		qos = int(q)
	}

	retained := false
	if r, ok := dataMap["retained"].(bool); ok {
		retained = r
	}

	downlinkMsg := &RawDownlinkMessage{
		DeviceSN: msg.DeviceSN,
		Topic:    topic,
		Payload:  payload,
		QoS:      qos,
		Retained: retained,
		TraceID:  msg.TID,
		SpanID:   msg.BID,
	}

	return b.Bridge(ctx, downlinkMsg)
}

// Stop stops the downlink bridge
func (b *DownlinkBridge) Stop() {
	b.mu.Lock()
	defer b.mu.Unlock()

	if !b.running {
		return
	}

	if b.subscriber != nil {
		_ = b.subscriber.Unsubscribe(b.queue)
	}

	close(b.stopCh)
	b.running = false
	b.logger.Info("Downlink bridge stopped")
}

// Bridge forwards a RabbitMQ message to MQTT
func (b *DownlinkBridge) Bridge(ctx context.Context, msg *RawDownlinkMessage) error {
	if b.mqttClient == nil {
		return fmt.Errorf("MQTT client not initialized")
	}

	if !b.mqttClient.IsConnected() {
		return fmt.Errorf("MQTT client not connected")
	}

	// Publish to MQTT
	err := b.mqttClient.Publish(msg.Topic, byte(msg.QoS), msg.Retained, msg.Payload)
	if err != nil {
		return fmt.Errorf("failed to publish to MQTT: %w", err)
	}

	b.logger.WithFields(logrus.Fields{
		"topic":     msg.Topic,
		"device_sn": msg.DeviceSN,
		"trace_id":  msg.TraceID,
	}).Debug("Bridged downlink message to MQTT")

	return nil
}

// HandleMessage handles a raw downlink message from RabbitMQ
func (b *DownlinkBridge) HandleMessage(ctx context.Context, data []byte) error {
	msg, err := ParseRawDownlinkMessage(data)
	if err != nil {
		return fmt.Errorf("failed to parse downlink message: %w", err)
	}

	return b.Bridge(ctx, msg)
}

// IsRunning returns whether the bridge is running
func (b *DownlinkBridge) IsRunning() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.running
}

// GetDownlinkRoutingKey returns the routing key for a given vendor
func GetDownlinkRoutingKey(vendor string) string {
	return fmt.Sprintf("iot.raw.%s.downlink", vendor)
}

// ParseRawDownlinkMessage parses a raw downlink message from bytes
func ParseRawDownlinkMessage(data []byte) (*RawDownlinkMessage, error) {
	return parseRawMessage[RawDownlinkMessage](data, "downlink")
}

// NewRawDownlinkMessage creates a new raw downlink message
func NewRawDownlinkMessage(deviceSN, topic string, payload json.RawMessage, qos int, retained bool, traceID, spanID string) *RawDownlinkMessage {
	return &RawDownlinkMessage{
		DeviceSN: deviceSN,
		Topic:    topic,
		Payload:  payload,
		QoS:      qos,
		Retained: retained,
		TraceID:  traceID,
		SpanID:   spanID,
	}
}
