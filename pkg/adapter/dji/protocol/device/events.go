package device

import "github.com/utmos/utmos/pkg/adapter/dji/protocol/common"

// ===============================
// Device Control Progress Events
// ===============================

// DeviceProgress represents the progress information
type DeviceProgress struct {
	Percent int    `json:"percent"`            // Progress percentage (0-100)
	StepKey string `json:"step_key,omitempty"` // Current step (optional, varies by operation)
}

// DeviceOutput represents the device operation output
type DeviceOutput struct {
	Status   string         `json:"status"`   // Operation status (sent, in_progress, ok, failed, canceled, paused, rejected, timeout)
	Progress DeviceProgress `json:"progress"` // Progress information
}

// DeviceProgressData represents the device operation progress data
type DeviceProgressData struct {
	Result int          `json:"result"` // Return code (0=success)
	Output DeviceOutput `json:"output"` // Output data
}

// CoverOpenProgressEvent represents the cover open progress event
type CoverOpenProgressEvent struct {
	common.Header
	MethodName string             `json:"method"`
	DataValue  DeviceProgressData `json:"data"`
}

func (e *CoverOpenProgressEvent) Method() string            { return e.MethodName }
func (e *CoverOpenProgressEvent) Data() any                 { return e.DataValue }
func (e *CoverOpenProgressEvent) GetHeader() *common.Header { return &e.Header }

// CoverCloseProgressEvent represents the cover close progress event
type CoverCloseProgressEvent struct {
	common.Header
	MethodName string             `json:"method"`
	DataValue  DeviceProgressData `json:"data"`
}

func (e *CoverCloseProgressEvent) Method() string            { return e.MethodName }
func (e *CoverCloseProgressEvent) Data() any                 { return e.DataValue }
func (e *CoverCloseProgressEvent) GetHeader() *common.Header { return &e.Header }

// CoverForceCloseProgressEvent represents the cover force close progress event
type CoverForceCloseProgressEvent struct {
	common.Header
	MethodName string             `json:"method"`
	DataValue  DeviceProgressData `json:"data"`
}

func (e *CoverForceCloseProgressEvent) Method() string            { return e.MethodName }
func (e *CoverForceCloseProgressEvent) Data() any                 { return e.DataValue }
func (e *CoverForceCloseProgressEvent) GetHeader() *common.Header { return &e.Header }

// DroneOpenProgressEvent represents the drone power on progress event
type DroneOpenProgressEvent struct {
	common.Header
	MethodName string             `json:"method"`
	DataValue  DeviceProgressData `json:"data"`
}

func (e *DroneOpenProgressEvent) Method() string            { return e.MethodName }
func (e *DroneOpenProgressEvent) Data() any                 { return e.DataValue }
func (e *DroneOpenProgressEvent) GetHeader() *common.Header { return &e.Header }

// DroneCloseProgressEvent represents the drone power off progress event
type DroneCloseProgressEvent struct {
	common.Header
	MethodName string             `json:"method"`
	DataValue  DeviceProgressData `json:"data"`
}

func (e *DroneCloseProgressEvent) Method() string            { return e.MethodName }
func (e *DroneCloseProgressEvent) Data() any                 { return e.DataValue }
func (e *DroneCloseProgressEvent) GetHeader() *common.Header { return &e.Header }

// ChargeOpenProgressEvent represents the charge start progress event
type ChargeOpenProgressEvent struct {
	common.Header
	MethodName string             `json:"method"`
	DataValue  DeviceProgressData `json:"data"`
}

func (e *ChargeOpenProgressEvent) Method() string            { return e.MethodName }
func (e *ChargeOpenProgressEvent) Data() any                 { return e.DataValue }
func (e *ChargeOpenProgressEvent) GetHeader() *common.Header { return &e.Header }

// ChargeCloseProgressEvent represents the charge stop progress event
type ChargeCloseProgressEvent struct {
	common.Header
	MethodName string             `json:"method"`
	DataValue  DeviceProgressData `json:"data"`
}

func (e *ChargeCloseProgressEvent) Method() string            { return e.MethodName }
func (e *ChargeCloseProgressEvent) Data() any                 { return e.DataValue }
func (e *ChargeCloseProgressEvent) GetHeader() *common.Header { return &e.Header }

// DeviceRebootProgressEvent represents the device reboot progress event
type DeviceRebootProgressEvent struct {
	common.Header
	MethodName string             `json:"method"`
	DataValue  DeviceProgressData `json:"data"`
}

func (e *DeviceRebootProgressEvent) Method() string            { return e.MethodName }
func (e *DeviceRebootProgressEvent) Data() any                 { return e.DataValue }
func (e *DeviceRebootProgressEvent) GetHeader() *common.Header { return &e.Header }

// DeviceFormatProgressEvent represents the dock data format progress event
type DeviceFormatProgressEvent struct {
	common.Header
	MethodName string             `json:"method"`
	DataValue  DeviceProgressData `json:"data"`
}

func (e *DeviceFormatProgressEvent) Method() string            { return e.MethodName }
func (e *DeviceFormatProgressEvent) Data() any                 { return e.DataValue }
func (e *DeviceFormatProgressEvent) GetHeader() *common.Header { return &e.Header }

// DroneFormatProgressEvent represents the drone data format progress event
type DroneFormatProgressEvent struct {
	common.Header
	MethodName string             `json:"method"`
	DataValue  DeviceProgressData `json:"data"`
}

func (e *DroneFormatProgressEvent) Method() string            { return e.MethodName }
func (e *DroneFormatProgressEvent) Data() any                 { return e.DataValue }
func (e *DroneFormatProgressEvent) GetHeader() *common.Header { return &e.Header }

// PutterOpenProgressEvent represents the putter open progress event
type PutterOpenProgressEvent struct {
	common.Header
	MethodName string             `json:"method"`
	DataValue  DeviceProgressData `json:"data"`
}

func (e *PutterOpenProgressEvent) Method() string            { return e.MethodName }
func (e *PutterOpenProgressEvent) Data() any                 { return e.DataValue }
func (e *PutterOpenProgressEvent) GetHeader() *common.Header { return &e.Header }

// PutterCloseProgressEvent represents the putter close progress event
type PutterCloseProgressEvent struct {
	common.Header
	MethodName string             `json:"method"`
	DataValue  DeviceProgressData `json:"data"`
}

func (e *PutterCloseProgressEvent) Method() string            { return e.MethodName }
func (e *PutterCloseProgressEvent) Data() any                 { return e.DataValue }
func (e *PutterCloseProgressEvent) GetHeader() *common.Header { return &e.Header }
