package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	dji "github.com/utmos/utmos/pkg/adapter/dji"
	"github.com/utmos/utmos/pkg/rabbitmq"
)

// StatusHandler handles Status (device online/offline) messages.
type StatusHandler struct{}

// NewStatusHandler creates a new Status handler.
func NewStatusHandler() *StatusHandler {
	return &StatusHandler{}
}

// Handle processes a Status message and returns a StandardMessage.
func (h *StatusHandler) Handle(ctx context.Context, msg *dji.Message, topic *dji.TopicInfo) (*rabbitmq.StandardMessage, error) {
	if msg == nil {
		return nil, fmt.Errorf("nil message")
	}
	if topic == nil {
		return nil, fmt.Errorf("nil topic info")
	}

	// Parse status data to determine online/offline
	isOnline, topology, err := h.parseStatusData(msg.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse status data: %w", err)
	}

	// Determine action based on online status
	action := dji.ActionDeviceOnline
	if !isOnline {
		action = dji.ActionDeviceOffline
	}

	// Build StandardMessage
	sm := &rabbitmq.StandardMessage{
		TID:       msg.TID,
		BID:       msg.BID,
		Timestamp: msg.Timestamp,
		DeviceSN:  topic.DeviceSN,
		Service:   dji.VendorDJI,
		Action:    action,
		ProtocolMeta: &rabbitmq.ProtocolMeta{
			Vendor:        dji.VendorDJI,
			OriginalTopic: topic.Raw,
			Method:        msg.Method,
		},
	}

	// Set timestamp if not provided
	if sm.Timestamp == 0 {
		sm.Timestamp = time.Now().UnixMilli()
	}

	// Build status data
	data, err := h.buildStatusData(isOnline, topology, topic)
	if err != nil {
		return nil, fmt.Errorf("failed to build status data: %w", err)
	}
	sm.Data = data

	return sm, nil
}

// GetTopicType returns the topic type this handler processes.
func (h *StatusHandler) GetTopicType() dji.TopicType {
	return dji.TopicTypeStatus
}

// DeviceTopology represents the device topology from status message.
type DeviceTopology struct {
	GatewaySN   string          `json:"gateway_sn"`
	GatewayType string          `json:"gateway_type,omitempty"`
	SubDevices  []SubDeviceInfo `json:"sub_devices,omitempty"`
}

// SubDeviceInfo represents a sub-device in the topology.
type SubDeviceInfo struct {
	DeviceSN    string `json:"device_sn"`
	ProductType string `json:"product_type,omitempty"`
	Online      bool   `json:"online"`
}

// parseStatusData parses the status data and extracts online status and topology.
func (h *StatusHandler) parseStatusData(data json.RawMessage) (bool, *DeviceTopology, error) {
	if len(data) == 0 {
		return false, nil, fmt.Errorf("empty status data")
	}

	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return false, nil, fmt.Errorf("failed to parse status JSON: %w", err)
	}

	// Check online status
	isOnline := true
	if online, ok := raw["online"].(bool); ok {
		isOnline = online
	} else if onlineNum, ok := raw["online"].(float64); ok {
		isOnline = onlineNum == 1
	}

	// Extract topology
	topology := &DeviceTopology{}

	if gatewaySN, ok := raw["gateway_sn"].(string); ok {
		topology.GatewaySN = gatewaySN
	}

	if gatewayType, ok := raw["gateway_type"].(string); ok {
		topology.GatewayType = gatewayType
	}

	// Extract sub-devices
	if subDevices, ok := raw["sub_devices"].([]any); ok {
		for _, sd := range subDevices {
			if sdMap, ok := sd.(map[string]any); ok {
				subDevice := SubDeviceInfo{}
				if sn, ok := sdMap["device_sn"].(string); ok {
					subDevice.DeviceSN = sn
				}
				if pt, ok := sdMap["product_type"].(string); ok {
					subDevice.ProductType = pt
				}
				if online, ok := sdMap["online"].(bool); ok {
					subDevice.Online = online
				} else if onlineNum, ok := sdMap["online"].(float64); ok {
					subDevice.Online = onlineNum == 1
				}
				topology.SubDevices = append(topology.SubDevices, subDevice)
			}
		}
	}

	return isOnline, topology, nil
}

// buildStatusData converts status data to a data map for StandardMessage.
func (h *StatusHandler) buildStatusData(isOnline bool, topology *DeviceTopology, topic *dji.TopicInfo) (json.RawMessage, error) {
	result := make(map[string]any)

	result["device_sn"] = topic.DeviceSN
	result["gateway_sn"] = topic.GatewaySN
	result["message_type"] = "status"
	result["online"] = isOnline

	if topology != nil {
		result["topology"] = topology
	}

	return json.Marshal(result)
}

// Ensure StatusHandler implements Handler interface.
var _ Handler = (*StatusHandler)(nil)
