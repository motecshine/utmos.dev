package firmware

import "github.com/utmos/utmos/pkg/adapter/dji/protocol/common"

// ===============================
// Firmware Upgrade Events
// ===============================

// OTAProgressProgress represents the firmware upgrade progress information
type OTAProgressProgress struct {
	Percent int    `json:"percent"`  // Progress percentage (0-100)
	StepKey string `json:"step_key"` // Current step (download_firmware, upgrade_firmware)
}

// OTAProgressOutput represents the firmware upgrade output
type OTAProgressOutput struct {
	Status   string              `json:"status"`   // Task status
	Progress OTAProgressProgress `json:"progress"` // Progress information
}

// OTAProgressData represents the firmware upgrade progress data
type OTAProgressData struct {
	Result int               `json:"result"` // Return code (0=success)
	Output OTAProgressOutput `json:"output"` // Output data
}

// OTAProgressEvent represents the firmware upgrade progress event
type OTAProgressEvent struct {
	common.Header
	MethodName string          `json:"method"`
	DataValue  OTAProgressData `json:"data"`
}

// Method returns the method name.
func (e *OTAProgressEvent) Method() string { return e.MethodName }

// Data returns the command/event data.
func (e *OTAProgressEvent) Data() any { return e.DataValue }

// GetHeader returns the event header.
func (e *OTAProgressEvent) GetHeader() *common.Header { return &e.Header }
