package aircraft

import "github.com/utmos/utmos/pkg/adapter/dji/protocol/common"

// ===============================
// Health Management System Events
// ===============================

// HMSArgs represents the HMS alert arguments
type HMSArgs struct {
	ComponentIndex int `json:"component_index"` // Component index variable for alert text
	SensorIndex    int `json:"sensor_index"`    // Sensor index variable for alert text
}

// HMSItem represents a single HMS alert item
type HMSItem struct {
	Level      int     `json:"level"`       // Alert level (0=notification, 1=reminder, 2=warning)
	Module     int     `json:"module"`      // Event module (0=flight task, 1=device management, 2=media, 3=hms)
	InTheSky   int     `json:"in_the_sky"`  // Whether flying (0=on ground, 1=in the air)
	Code       string  `json:"code"`        // Alert code
	DeviceType string  `json:"device_type"` // Device type (format: domain-type-subtype)
	Imminent   int     `json:"imminent"`    // Whether it is an immediate alert (0=no, 1=yes)
	Args       HMSArgs `json:"args"`        // Arguments for alert text variables
}

// HMSData represents the HMS health alert data
type HMSData struct {
	List []HMSItem `json:"list"` // Health alert list (max 20 items)
}

// HMSEvent represents the HMS health alert event
type HMSEvent struct {
	common.Header
	MethodName string  `json:"method"`
	DataValue  HMSData `json:"data"`
}

func (e *HMSEvent) Method() string            { return e.MethodName }
func (e *HMSEvent) Data() any                 { return e.DataValue }
func (e *HMSEvent) GetHeader() *common.Header { return &e.Header }
