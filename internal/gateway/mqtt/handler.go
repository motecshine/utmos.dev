// Package mqtt provides MQTT message handling functionality
package mqtt

import (
	"context"
	"encoding/json"
	"strings"
	"sync"
	"time"

	pahomqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	pkgtracer "github.com/utmos/utmos/pkg/tracer"
)

// Message represents a parsed MQTT message
type Message struct {
	Topic     string          `json:"topic"`
	Payload   json.RawMessage `json:"payload"`
	QoS       byte            `json:"qos"`
	Retained  bool            `json:"retained"`
	MessageID uint16          `json:"message_id"`
	Timestamp time.Time       `json:"timestamp"`
	TraceID   string          `json:"trace_id"`
	SpanID    string          `json:"span_id"`
}

// TopicInfo contains parsed topic information
type TopicInfo struct {
	Vendor    string
	ProductID string
	DeviceSN  string
	Service   string
	Method    string
	Raw       string
}

// Handler processes MQTT messages
type Handler struct {
	logger     *logrus.Entry
	processors map[string]MessageProcessor
	mu         sync.RWMutex
}

// MessageProcessor processes messages for a specific topic pattern
type MessageProcessor interface {
	Process(ctx context.Context, msg *Message, topicInfo *TopicInfo) error
	Pattern() string
}

// MessageCallback is called when a message is processed
type MessageCallback func(ctx context.Context, msg *Message, topicInfo *TopicInfo) error

// NewHandler creates a new message handler
func NewHandler(logger *logrus.Entry) *Handler {
	if logger == nil {
		logger = logrus.NewEntry(logrus.StandardLogger())
	}
	return &Handler{
		logger:     logger.WithField("component", "mqtt-handler"),
		processors: make(map[string]MessageProcessor),
	}
}

// RegisterProcessor registers a message processor for a topic pattern
func (h *Handler) RegisterProcessor(processor MessageProcessor) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.processors[processor.Pattern()] = processor
	h.logger.WithField("pattern", processor.Pattern()).Debug("Registered message processor")
}

// UnregisterProcessor unregisters a message processor
func (h *Handler) UnregisterProcessor(pattern string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.processors, pattern)
}

// Handle processes an incoming MQTT message
func (h *Handler) Handle(client pahomqtt.Client, mqttMsg pahomqtt.Message) {
	msg := &Message{
		Topic:     mqttMsg.Topic(),
		Payload:   mqttMsg.Payload(),
		QoS:       mqttMsg.Qos(),
		Retained:  mqttMsg.Retained(),
		MessageID: mqttMsg.MessageID(),
		Timestamp: time.Now(),
	}

	topicInfo := ParseTopic(msg.Topic)

	// Create a root span for the MQTT message
	tr := otel.Tracer("iot-gateway")
	ctx, span := tr.Start(context.Background(), "mqtt.message.received",
		trace.WithSpanKind(trace.SpanKindConsumer),
		trace.WithAttributes(
			attribute.String("mqtt.topic", msg.Topic),
			attribute.String("device_sn", topicInfo.DeviceSN),
		),
	)
	defer span.End()

	msg.TraceID = pkgtracer.GetTraceID(ctx)
	msg.SpanID = pkgtracer.GetSpanID(ctx)

	h.logger.WithFields(logrus.Fields{
		"topic":      msg.Topic,
		"qos":        msg.QoS,
		"message_id": msg.MessageID,
		"trace_id":   msg.TraceID,
		"device_sn":  topicInfo.DeviceSN,
	}).Debug("Received MQTT message")

	h.mu.RLock()
	defer h.mu.RUnlock()

	for pattern, processor := range h.processors {
		if matchTopic(pattern, msg.Topic) {
			if err := processor.Process(ctx, msg, topicInfo); err != nil {
				h.logger.WithError(err).WithFields(logrus.Fields{
					"topic":    msg.Topic,
					"pattern":  pattern,
					"trace_id": msg.TraceID,
				}).Error("Failed to process message")
			}
		}
	}
}

// extractTopicParts extracts topic components from parts starting at the given offset
func extractTopicParts(parts []string, offset int, info *TopicInfo) {
	if len(parts) >= offset+3 {
		info.ProductID = parts[offset]
		info.DeviceSN = parts[offset+1]
		info.Service = parts[offset+2]
		if len(parts) >= offset+4 {
			info.Method = parts[offset+3]
		}
	}
}

// ParseTopic parses an MQTT topic into its components
// Expected formats:
// - thing/product/{product_id}/{device_sn}/{service}
// - sys/product/{device_sn}/{service}
// - {vendor}/thing/product/{device_sn}/{service}
func ParseTopic(topic string) *TopicInfo {
	info := &TopicInfo{Raw: topic}
	parts := strings.Split(topic, "/")

	if len(parts) < 3 {
		return info
	}

	// Detect vendor prefix
	switch parts[0] {
	case "thing", "sys":
		// DJI format: thing/product/{device_sn}/{service}
		info.Vendor = "dji"
		extractTopicParts(parts, 1, info)
	default:
		// Generic format: {vendor}/thing/product/{device_sn}/{service}
		info.Vendor = parts[0]
		extractTopicParts(parts, 2, info)
	}

	return info
}

// matchTopic checks if a topic matches a pattern
// Supports wildcards: + (single level) and # (multi level)
func matchTopic(pattern, topic string) bool {
	patternParts := strings.Split(pattern, "/")
	topicParts := strings.Split(topic, "/")

	patternIdx := 0
	topicIdx := 0

	for patternIdx < len(patternParts) && topicIdx < len(topicParts) {
		switch patternParts[patternIdx] {
		case "#":
			// # matches everything from here
			return true
		case "+":
			// + matches exactly one level
			patternIdx++
			topicIdx++
		default:
			if patternParts[patternIdx] != topicParts[topicIdx] {
				return false
			}
			patternIdx++
			topicIdx++
		}
	}

	// Check if we've consumed all parts
	if patternIdx == len(patternParts) && topicIdx == len(topicParts) {
		return true
	}

	// Handle trailing #
	if patternIdx == len(patternParts)-1 && patternParts[patternIdx] == "#" {
		return true
	}

	return false
}

// SimpleProcessor is a simple message processor implementation
type SimpleProcessor struct {
	pattern  string
	callback MessageCallback
}

// NewSimpleProcessor creates a new simple processor
func NewSimpleProcessor(pattern string, callback MessageCallback) *SimpleProcessor {
	return &SimpleProcessor{
		pattern:  pattern,
		callback: callback,
	}
}

// Process processes a message
func (p *SimpleProcessor) Process(ctx context.Context, msg *Message, topicInfo *TopicInfo) error {
	return p.callback(ctx, msg, topicInfo)
}

// Pattern returns the topic pattern
func (p *SimpleProcessor) Pattern() string {
	return p.pattern
}
