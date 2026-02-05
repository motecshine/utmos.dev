package dji

import (
	"fmt"
	"strings"
)

// TopicInfo contains parsed information from a DJI MQTT topic.
type TopicInfo struct {
	Type      TopicType // Topic type (osd, state, services, events, status)
	DeviceSN  string    // Device serial number
	GatewaySN string    // Gateway serial number (same as DeviceSN for gateway topics)
	Raw       string    // Original raw topic string
}

// IsUplink returns true if this is an uplink (device to cloud) topic.
func (ti *TopicInfo) IsUplink() bool {
	switch ti.Type {
	case TopicTypeOSD, TopicTypeState, TopicTypeEvents, TopicTypeStatus, TopicTypeServicesReply:
		return true
	default:
		return false
	}
}

// IsDownlink returns true if this is a downlink (cloud to device) topic.
func (ti *TopicInfo) IsDownlink() bool {
	switch ti.Type {
	case TopicTypeServices, TopicTypeStatusReply:
		return true
	default:
		return false
	}
}

// ParseTopic parses a DJI MQTT topic string into TopicInfo.
// Supported topic patterns:
//   - thing/product/{device_sn}/osd
//   - thing/product/{device_sn}/state
//   - thing/product/{gateway_sn}/services
//   - thing/product/{gateway_sn}/services_reply
//   - thing/product/{gateway_sn}/events
//   - sys/product/{gateway_sn}/status
//   - sys/product/{gateway_sn}/status_reply
func ParseTopic(topic string) (*TopicInfo, error) {
	if topic == "" {
		return nil, ErrEmptyTopic
	}

	parts := strings.Split(topic, "/")
	if len(parts) < 4 {
		return nil, fmt.Errorf("%w: expected at least 4 segments, got %d", ErrInvalidTopic, len(parts))
	}

	prefix := parts[0]
	if prefix != "thing" && prefix != "sys" {
		return nil, fmt.Errorf("%w: unknown prefix %q", ErrInvalidTopic, prefix)
	}

	if parts[1] != "product" {
		return nil, fmt.Errorf("%w: expected 'product' as second segment", ErrInvalidTopic)
	}

	deviceSN := parts[2]
	topicTypeStr := parts[3]

	topicType, err := parseTopicType(topicTypeStr)
	if err != nil {
		return nil, err
	}

	return &TopicInfo{
		Type:      topicType,
		DeviceSN:  deviceSN,
		GatewaySN: deviceSN, // For DJI, gateway_sn is often the same as device_sn
		Raw:       topic,
	}, nil
}

// parseTopicType converts a topic type string to TopicType.
func parseTopicType(s string) (TopicType, error) {
	switch s {
	case "osd":
		return TopicTypeOSD, nil
	case "state":
		return TopicTypeState, nil
	case "services":
		return TopicTypeServices, nil
	case "services_reply":
		return TopicTypeServicesReply, nil
	case "events":
		return TopicTypeEvents, nil
	case "status":
		return TopicTypeStatus, nil
	case "status_reply":
		return TopicTypeStatusReply, nil
	default:
		return "", fmt.Errorf("%w: %q", ErrUnknownTopicType, s)
	}
}

// BuildTopic constructs a DJI MQTT topic from components.
func BuildTopic(topicType TopicType, deviceSN string) string {
	prefix := "thing"
	if topicType == TopicTypeStatus || topicType == TopicTypeStatusReply {
		prefix = "sys"
	}
	return fmt.Sprintf("%s/product/%s/%s", prefix, deviceSN, topicType)
}
