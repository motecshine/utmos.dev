// Package integration provides protocol integration utilities for the DJI adapter.
package integration

import (
	"encoding/json"
	"fmt"

	"github.com/utmos/utmos/pkg/adapter/dji/protocol/aircraft"
)

// OSDParser parses OSD data from DJI messages.
type OSDParser struct{}

// NewOSDParser creates a new OSD parser.
func NewOSDParser() *OSDParser {
	return &OSDParser{}
}

// unmarshalOSD is a generic helper that unmarshals JSON data into a typed OSD struct.
func unmarshalOSD[T any](data json.RawMessage, typeName string) (*T, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty OSD data")
	}

	var osd T
	if err := json.Unmarshal(data, &osd); err != nil {
		return nil, fmt.Errorf("failed to parse %s OSD: %w", typeName, err)
	}

	return &osd, nil
}

// ParseAircraftOSD parses aircraft OSD data from raw JSON.
func (p *OSDParser) ParseAircraftOSD(data json.RawMessage) (*aircraft.AircraftOSD, error) {
	return unmarshalOSD[aircraft.AircraftOSD](data, "aircraft")
}

// ParseDockOSD parses dock OSD data from raw JSON.
func (p *OSDParser) ParseDockOSD(data json.RawMessage) (*aircraft.DockOSD, error) {
	return unmarshalOSD[aircraft.DockOSD](data, "dock")
}

// ParseRCOSD parses remote controller OSD data from raw JSON.
func (p *OSDParser) ParseRCOSD(data json.RawMessage) (*aircraft.RCOSD, error) {
	return unmarshalOSD[aircraft.RCOSD](data, "RC")
}

// OSDType represents the type of OSD data.
type OSDType string

const (
	OSDTypeAircraft OSDType = "aircraft"
	OSDTypeDock     OSDType = "dock"
	OSDTypeRC       OSDType = "rc"
)

// DetectOSDType attempts to detect the OSD type from raw JSON data.
// This is a heuristic based on field presence.
func (p *OSDParser) DetectOSDType(data json.RawMessage) OSDType {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return OSDTypeAircraft // default
	}

	// Dock-specific fields
	if _, ok := raw["cover_state"]; ok {
		return OSDTypeDock
	}
	if _, ok := raw["drone_in_dock"]; ok {
		return OSDTypeDock
	}
	if _, ok := raw["putter_state"]; ok {
		return OSDTypeDock
	}

	// RC-specific fields
	if _, ok := raw["wireless_link"]; ok {
		// Check if it's RC wireless link (has specific structure)
		if _, hasCapacity := raw["capacity_percent"]; hasCapacity {
			if _, hasPayloads := raw["payloads"]; !hasPayloads {
				return OSDTypeRC
			}
		}
	}

	// Default to aircraft
	return OSDTypeAircraft
}

// ParsedOSD holds the parsed OSD data with type information.
type ParsedOSD struct {
	Type     OSDType
	Aircraft *aircraft.AircraftOSD
	Dock     *aircraft.DockOSD
	RC       *aircraft.RCOSD
}

// ParseOSD parses OSD data and automatically detects the type.
func (p *OSDParser) ParseOSD(data json.RawMessage) (*ParsedOSD, error) {
	osdType := p.DetectOSDType(data)

	result := &ParsedOSD{Type: osdType}

	switch osdType {
	case OSDTypeDock:
		dock, err := p.ParseDockOSD(data)
		if err != nil {
			return nil, err
		}
		result.Dock = dock
	case OSDTypeRC:
		rc, err := p.ParseRCOSD(data)
		if err != nil {
			return nil, err
		}
		result.RC = rc
	default:
		aircraft, err := p.ParseAircraftOSD(data)
		if err != nil {
			return nil, err
		}
		result.Aircraft = aircraft
	}

	return result, nil
}
