package wayline

import "github.com/utmos/utmos/pkg/adapter/dji/protocol/common"

// ===============================
// Wayline Mission Events (Device → Cloud)
// ===============================

// FlightSetupExceptionEvent represents the flight setup exception notification event
type FlightSetupExceptionEvent struct {
	common.Header
	MethodName string                         `json:"method"`
	DataValue  FlightSetupExceptionNotifyData `json:"data"`
}

// Method returns the method name.
func (e *FlightSetupExceptionEvent) Method() string { return e.MethodName }

// Data returns the command/event data.
func (e *FlightSetupExceptionEvent) Data() any { return e.DataValue }

// GetHeader returns the event header.
func (e *FlightSetupExceptionEvent) GetHeader() *common.Header { return &e.Header }

// ExitHomingEvent represents the device exit homing notification event
type ExitHomingEvent struct {
	common.Header
	MethodName string               `json:"method"`
	DataValue  ExitHomingNotifyData `json:"data"`
}

// Method returns the method name.
func (e *ExitHomingEvent) Method() string { return e.MethodName }

// Data returns the command/event data.
func (e *ExitHomingEvent) Data() any { return e.DataValue }

// GetHeader returns the event header.
func (e *ExitHomingEvent) GetHeader() *common.Header { return &e.Header }

// ProgressEvent represents the flight task progress event
type ProgressEvent struct {
	common.Header
	MethodName string       `json:"method"`
	DataValue  ProgressData `json:"data"`
}

// Method returns the method name.
func (e *ProgressEvent) Method() string { return e.MethodName }

// Data returns the command/event data.
func (e *ProgressEvent) Data() any { return e.DataValue }

// GetHeader returns the event header.
func (e *ProgressEvent) GetHeader() *common.Header { return &e.Header }

// ReadyEvent represents the task ready notification event
type ReadyEvent struct {
	common.Header
	MethodName string    `json:"method"`
	DataValue  ReadyData `json:"data"`
}

// Method returns the method name.
func (e *ReadyEvent) Method() string { return e.MethodName }

// Data returns the command/event data.
func (e *ReadyEvent) Data() any { return e.DataValue }

// GetHeader returns the event header.
func (e *ReadyEvent) GetHeader() *common.Header { return &e.Header }

// ReturnHomeInfoEvent represents the return home information event
type ReturnHomeInfoEvent struct {
	common.Header
	MethodName string             `json:"method"`
	DataValue  ReturnHomeInfoData `json:"data"`
}

// Method returns the method name.
func (e *ReturnHomeInfoEvent) Method() string { return e.MethodName }

// Data returns the command/event data.
func (e *ReturnHomeInfoEvent) Data() any { return e.DataValue }

// GetHeader returns the event header.
func (e *ReturnHomeInfoEvent) GetHeader() *common.Header { return &e.Header }

// ===============================
// Event Data Types
// ===============================

// FlightSetupExceptionNotifyData represents the flight setup exception notification data
type FlightSetupExceptionNotifyData struct {
	SN          string  `json:"sn"`           // Dock serial number
	TimeoutTime int     `json:"timeout_time"` // Exception timeout time (minutes, 2-10)
	Timestamp   float64 `json:"timestamp"`    // Current UTC time
	FlightType  int     `json:"flight_type"`  // Task type (1=wayline task, 2=command flight task)
}

// ExitHomingNotifyData represents the device exit homing notification data
type ExitHomingNotifyData struct {
	SN     string `json:"sn"`     // Dock serial number
	Action int    `json:"action"` // Exit return-to-home notification type (0=exit "return-to-home exit state", 1=enter "return-to-home exit state")
	Reason int    `json:"reason"` // Exit return-to-home reason code
}

// ProgressData represents the flight task progress data
type ProgressData struct {
	Output ProgressOutput `json:"output"` // Output data
	Result int            `json:"result"` // Return code (0=success)
}

// ProgressOutput represents the output of flight task progress
type ProgressOutput struct {
	Ext      ProgressExt      `json:"ext"`      // Extended content
	Status   string           `json:"status"`   // Task status
	Progress ProgressProgress `json:"progress"` // Progress information
}

// ProgressExt represents the extended information in flight task progress
type ProgressExt struct {
	CurrentWaypointIndex int             `json:"current_waypoint_index"` // Current waypoint number
	WaylineMissionState  int             `json:"wayline_mission_state"`  // Wayline mission state
	MediaCount           int             `json:"media_count"`            // Media file count for this mission
	TrackID              string          `json:"track_id"`               // Track ID
	FlightID             string          `json:"flight_id"`              // Flight task ID
	BreakPoint           *BreakPointInfo `json:"break_point,omitempty"`  // Breakpoint information (optional)
	WaylineID            int             `json:"wayline_id"`             // Current wayline ID being executed
}

// ProgressProgress represents the progress information
type ProgressProgress struct {
	CurrentStep int `json:"current_step"` // Execution step
	Percent     int `json:"percent"`      // Progress percentage (0-100)
}

// BreakPointInfo represents the wayline breakpoint information
type BreakPointInfo struct {
	Index        int     `json:"index"`         // Breakpoint index
	State        int     `json:"state"`         // Breakpoint state (0=on segment, 1=on waypoint)
	Progress     float64 `json:"progress"`      // Current segment progress (0-1.0)
	WaylineID    int     `json:"wayline_id"`    // Wayline ID
	BreakReason  int     `json:"break_reason"`  // Break reason code
	Latitude     float64 `json:"latitude"`      // Breakpoint latitude
	Longitude    float64 `json:"longitude"`     // Breakpoint longitude
	Height       float64 `json:"height"`        // Breakpoint height relative to Earth ellipsoid
	AttitudeHead float64 `json:"attitude_head"` // Breakpoint yaw angle
}

// ReadyData represents the task ready notification data
type ReadyData struct {
	FlightIDs []string `json:"flight_ids"` // Task IDs that meet the ready conditions
}

// ReturnHomeInfoData represents the return home information data
type ReturnHomeInfoData struct {
	PlannedPathPoints []common.PlannedPathPoint `json:"planned_path_points"` // Planned return home trajectory points
	LastPointType     int                       `json:"last_point_type"`     // Last point type (0=above return point, 1=not above return point)
	FlightID          string                    `json:"flight_id"`           // Current flight task ID
}
