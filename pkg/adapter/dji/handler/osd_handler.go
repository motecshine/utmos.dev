package handler

import (
	"context"
	"encoding/json"
	"fmt"

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
//
// Same handler flow as StateHandler but with OSD-specific parsing
func (h *OSDHandler) Handle(_ context.Context, msg *dji.Message, topic *dji.TopicInfo) (*rabbitmq.StandardMessage, error) {
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

	// Build StandardMessage using shared helper
	cfg := MessageConfig{
		RequestAction: dji.ActionPropertyReport,
	}
	sm := BuildStandardMessage(msg, topic, cfg)

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
	result := map[string]any{
		"osd_type":   string(osd.Type),
		"device_sn":  topic.DeviceSN,
		"gateway_sn": topic.GatewaySN,
	}

	switch osd.Type {
	case integration.OSDTypeDock:
		h.extractDockFields(osd, result)
	case integration.OSDTypeRC:
		h.extractRCFields(osd, result)
	case integration.OSDTypeAircraft:
		h.extractAircraftFields(osd, result)
	}

	return json.Marshal(result)
}

// extractDockFields extracts key fields from DockOSD.
func (h *OSDHandler) extractDockFields(osd *integration.ParsedOSD, result map[string]any) {
	if osd.Dock == nil {
		return
	}
	result["dock"] = osd.Dock
	setIfNotNil(result, "mode_code", osd.Dock.ModeCode)
	setIfNotNil(result, "cover_state", osd.Dock.CoverState)
	setIfNotNil(result, "drone_in_dock", osd.Dock.DroneInDock)
	setIfNotNil(result, "longitude", osd.Dock.Longitude)
	setIfNotNil(result, "latitude", osd.Dock.Latitude)
}

// extractRCFields extracts key fields from RCOSD.
func (h *OSDHandler) extractRCFields(osd *integration.ParsedOSD, result map[string]any) {
	if osd.RC == nil {
		return
	}
	result["rc"] = osd.RC
	setIfNotNil(result, "capacity_percent", osd.RC.CapacityPercent)
	setIfNotNil(result, "longitude", osd.RC.Longitude)
	setIfNotNil(result, "latitude", osd.RC.Latitude)
}

// extractAircraftFields extracts key fields from AircraftOSD.
func (h *OSDHandler) extractAircraftFields(osd *integration.ParsedOSD, result map[string]any) {
	if osd.Aircraft == nil {
		return
	}
	result["aircraft"] = osd.Aircraft
	setIfNotNil(result, "mode_code", osd.Aircraft.ModeCode)
	setIfNotNil(result, "longitude", osd.Aircraft.Longitude)
	setIfNotNil(result, "latitude", osd.Aircraft.Latitude)
	setIfNotNil(result, "height", osd.Aircraft.Height)
	setIfNotNil(result, "elevation", osd.Aircraft.Elevation)
	setIfNotNil(result, "horizontal_speed", osd.Aircraft.HorizontalSpeed)
	setIfNotNil(result, "vertical_speed", osd.Aircraft.VerticalSpeed)
	if osd.Aircraft.Battery != nil {
		setIfNotNil(result, "battery_percent", osd.Aircraft.Battery.CapacityPercent)
	}
}

// setIfNotNil sets a value in the map if the pointer is not nil.
func setIfNotNil[T any](m map[string]any, key string, ptr *T) {
	if ptr != nil {
		m[key] = *ptr
	}
}

// Ensure OSDHandler implements Handler interface.
var _ Handler = (*OSDHandler)(nil)
