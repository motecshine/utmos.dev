package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	dji "github.com/utmos/utmos/pkg/adapter/dji"
	"github.com/utmos/utmos/pkg/adapter/dji/integration"
	"github.com/utmos/utmos/pkg/rabbitmq"
)

// OSDHandler handles OSD (Operational Status Data) messages.
type OSDHandler struct {
	parser *integration.OSDParser
}

// NewOSDHandler creates a new OSD handler.
func NewOSDHandler() *OSDHandler {
	return &OSDHandler{
		parser: integration.NewOSDParser(),
	}
}

// Handle processes an OSD message and returns a StandardMessage.
func (h *OSDHandler) Handle(ctx context.Context, msg *dji.Message, topic *dji.TopicInfo) (*rabbitmq.StandardMessage, error) {
	if msg == nil {
		return nil, fmt.Errorf("nil message")
	}
	if topic == nil {
		return nil, fmt.Errorf("nil topic info")
	}

	// Parse OSD data with auto-detection
	parsedOSD, err := h.parser.ParseOSD(msg.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse OSD: %w", err)
	}

	// Build StandardMessage
	sm := &rabbitmq.StandardMessage{
		TID:       msg.TID,
		BID:       msg.BID,
		Timestamp: msg.Timestamp,
		DeviceSN:  topic.DeviceSN,
		Service:   dji.VendorDJI,
		Action:    dji.ActionPropertyReport,
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

	// Convert parsed OSD to data map
	data, err := h.buildOSDData(parsedOSD, topic)
	if err != nil {
		return nil, fmt.Errorf("failed to build OSD data: %w", err)
	}
	sm.Data = data

	return sm, nil
}

// GetTopicType returns the topic type this handler processes.
func (h *OSDHandler) GetTopicType() dji.TopicType {
	return dji.TopicTypeOSD
}

// buildOSDData converts parsed OSD to a data map for StandardMessage.
func (h *OSDHandler) buildOSDData(osd *integration.ParsedOSD, topic *dji.TopicInfo) (json.RawMessage, error) {
	result := make(map[string]interface{})

	result["osd_type"] = string(osd.Type)
	result["device_sn"] = topic.DeviceSN
	result["gateway_sn"] = topic.GatewaySN

	switch osd.Type {
	case integration.OSDTypeDock:
		if osd.Dock != nil {
			result["dock"] = osd.Dock
			// Extract key fields for quick access
			if osd.Dock.ModeCode != nil {
				result["mode_code"] = *osd.Dock.ModeCode
			}
			if osd.Dock.CoverState != nil {
				result["cover_state"] = *osd.Dock.CoverState
			}
			if osd.Dock.DroneInDock != nil {
				result["drone_in_dock"] = *osd.Dock.DroneInDock
			}
			if osd.Dock.Longitude != nil {
				result["longitude"] = *osd.Dock.Longitude
			}
			if osd.Dock.Latitude != nil {
				result["latitude"] = *osd.Dock.Latitude
			}
		}

	case integration.OSDTypeRC:
		if osd.RC != nil {
			result["rc"] = osd.RC
			if osd.RC.CapacityPercent != nil {
				result["capacity_percent"] = *osd.RC.CapacityPercent
			}
			if osd.RC.Longitude != nil {
				result["longitude"] = *osd.RC.Longitude
			}
			if osd.RC.Latitude != nil {
				result["latitude"] = *osd.RC.Latitude
			}
		}

	case integration.OSDTypeAircraft:
		if osd.Aircraft != nil {
			result["aircraft"] = osd.Aircraft
			// Extract key fields for quick access
			if osd.Aircraft.ModeCode != nil {
				result["mode_code"] = *osd.Aircraft.ModeCode
			}
			if osd.Aircraft.Longitude != nil {
				result["longitude"] = *osd.Aircraft.Longitude
			}
			if osd.Aircraft.Latitude != nil {
				result["latitude"] = *osd.Aircraft.Latitude
			}
			if osd.Aircraft.Height != nil {
				result["height"] = *osd.Aircraft.Height
			}
			if osd.Aircraft.Elevation != nil {
				result["elevation"] = *osd.Aircraft.Elevation
			}
			if osd.Aircraft.HorizontalSpeed != nil {
				result["horizontal_speed"] = *osd.Aircraft.HorizontalSpeed
			}
			if osd.Aircraft.VerticalSpeed != nil {
				result["vertical_speed"] = *osd.Aircraft.VerticalSpeed
			}
			if osd.Aircraft.Battery != nil && osd.Aircraft.Battery.CapacityPercent != nil {
				result["battery_percent"] = *osd.Aircraft.Battery.CapacityPercent
			}
		}
	}

	return json.Marshal(result)
}

// Ensure OSDHandler implements Handler interface.
var _ Handler = (*OSDHandler)(nil)
