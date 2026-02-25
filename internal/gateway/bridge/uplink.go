// Package bridge provides MQTT to RabbitMQ bridging functionality
package bridge

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/utmos/utmos/internal/gateway/mqtt"
	"github.com/utmos/utmos/pkg/rabbitmq"
)

// RawUplinkMessage represents a raw uplink message from MQTT
type RawUplinkMessage struct {
	Vendor    string          `json:"vendor"`
	Topic     string          `json:"topic"`
	Payload   json.RawMessage `json:"payload"`
	QoS       int             `json:"qos"`
	Timestamp int64           `json:"timestamp"`
	TraceID   string          `json:"trace_id"`
	SpanID    string          `json:"span_id"`
}

// UplinkBridge bridges MQTT messages to RabbitMQ
type UplinkBridge struct {
	publisher *rabbitmq.Publisher
	logger    *logrus.Entry
	exchange  string
}

// UplinkBridgeConfig holds configuration for uplink bridge
type UplinkBridgeConfig struct {
	Exchange string
}

// DefaultUplinkBridgeConfig returns default configuration
func DefaultUplinkBridgeConfig() *UplinkBridgeConfig {
	return &UplinkBridgeConfig{
		Exchange: "iot.topic",
	}
}

// NewUplinkBridge creates a new uplink bridge
func NewUplinkBridge(publisher *rabbitmq.Publisher, config *UplinkBridgeConfig, logger *logrus.Entry) *UplinkBridge {
	if config == nil {
		config = DefaultUplinkBridgeConfig()
	}
	if logger == nil {
		logger = logrus.NewEntry(logrus.StandardLogger())
	}

	return &UplinkBridge{
		publisher: publisher,
		logger:    logger.WithField("component", "uplink-bridge"),
		exchange:  config.Exchange,
	}
}

// Bridge forwards an MQTT message to RabbitMQ
func (b *UplinkBridge) Bridge(ctx context.Context, msg *mqtt.Message, topicInfo *mqtt.TopicInfo) error {
	if b.publisher == nil {
		return fmt.Errorf("publisher not initialized")
	}

	// Create a child span for the uplink bridge operation
	tr := otel.Tracer("iot-gateway")
	ctx, span := tr.Start(ctx, "gateway.bridge.uplink",
		trace.WithAttributes(
			attribute.String("device_sn", topicInfo.DeviceSN),
			attribute.String("mqtt.topic", msg.Topic),
			attribute.String("vendor", topicInfo.Vendor),
		),
	)
	defer span.End()

	// Determine routing key: iot.raw.{vendor}.uplink
	routingKey := fmt.Sprintf("iot.raw.%s.uplink", topicInfo.Vendor)

	// Create data payload
	dataPayload := map[string]any{
		"vendor":  topicInfo.Vendor,
		"topic":   msg.Topic,
		"payload": msg.Payload,
		"qos":     msg.QoS,
	}

	// Create StandardMessage for RabbitMQ
	qos := int(msg.QoS)
	stdMsg, err := rabbitmq.NewStandardMessageWithIDs(
		msg.TraceID,
		msg.SpanID,
		"raw",
		"uplink",
		topicInfo.DeviceSN,
		dataPayload,
	)
	if err != nil {
		return fmt.Errorf("failed to create standard message: %w", err)
	}

	// Set protocol metadata
	stdMsg.ProtocolMeta = &rabbitmq.ProtocolMeta{
		Vendor:        topicInfo.Vendor,
		OriginalTopic: msg.Topic,
		QoS:           &qos,
	}

	// Publish to RabbitMQ
	err = b.publisher.Publish(ctx, routingKey, stdMsg)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	b.logger.WithFields(logrus.Fields{
		"topic":       msg.Topic,
		"routing_key": routingKey,
		"trace_id":    msg.TraceID,
		"device_sn":   topicInfo.DeviceSN,
	}).Debug("Bridged uplink message to RabbitMQ")

	return nil
}

// CreateProcessor creates a message processor for the uplink bridge
func (b *UplinkBridge) CreateProcessor(pattern string) *mqtt.SimpleProcessor {
	return mqtt.NewSimpleProcessor(pattern, func(ctx context.Context, msg *mqtt.Message, topicInfo *mqtt.TopicInfo) error {
		return b.Bridge(ctx, msg, topicInfo)
	})
}

// GetUplinkRoutingKey returns the routing key for a given vendor
func GetUplinkRoutingKey(vendor string) string {
	return fmt.Sprintf("iot.raw.%s.uplink", vendor)
}

// ParseRawUplinkMessage parses a raw uplink message from bytes
func ParseRawUplinkMessage(data []byte) (*RawUplinkMessage, error) {
	return parseRawMessage[RawUplinkMessage](data, "uplink")
}

// NewRawUplinkMessage creates a new raw uplink message
func NewRawUplinkMessage(vendor, topic string, payload json.RawMessage, qos int) *RawUplinkMessage {
	return &RawUplinkMessage{
		Vendor:    vendor,
		Topic:     topic,
		Payload:   payload,
		QoS:       qos,
		Timestamp: time.Now().UnixMilli(),
		TraceID:   uuid.New().String(),
		SpanID:    uuid.New().String(),
	}
}
