package wayline

// ===============================
// Wayline Mission Data Types
// ===============================

// CreateData represents the create flight task data
type CreateData struct {
	FlighttaskID             string                `json:"flighttask_id"`
	ExecuteTime              *int64                `json:"execute_time,omitempty"`                // Optional execution time (timestamp)
	TaskType                 string                `json:"task_type"`                             // Task type: immediate, timed, conditional
	WaylineType              string                `json:"wayline_type"`                          // Wayline type: wayline, mapping_2d, mapping_3d, mapping_strip
	WaylineFile              *string               `json:"wayline_file,omitempty"`                // Wayline file URL
	ExecuteTimes             *int                  `json:"execute_times,omitempty"`               // Number of execution times
	OutOfControlAction       string                `json:"out_of_control_action"`                 // Action when out of control: continue, execute_go_home, hover
	RthAltitude              *float64              `json:"rth_altitude,omitempty"`                // Return to home altitude
	TaskName                 *string               `json:"task_name,omitempty"`                   // Task name
	ReadyConditions          *ReadyConditions      `json:"ready_conditions,omitempty"`            // Ready conditions (for conditional tasks)
	ExecutableConditions     *ExecutableConditions `json:"executable_conditions,omitempty"`       // Executable conditions
	SimulateMission          *SimulateMission      `json:"simulate_mission,omitempty"`            // Simulation settings
	BreakPoint               *BreakPoint           `json:"break_point,omitempty"`                 // Breakpoint information for resume
	RthMode                  *int                  `json:"rth_mode,omitempty"`                    // Return to home mode (0=smart, 1=preset)
	ExitWaylineWhenRCLost    *int                  `json:"exit_wayline_when_rc_lost,omitempty"`   // Exit wayline when RC lost (0=continue, 1=exit)
	WaylinePrecisionType     *int                  `json:"wayline_precision_type,omitempty"`      // Wayline precision type (0=GPS, 1=RTK)
	FlightSafetyAdvanceCheck *bool                 `json:"flight_safety_advance_check,omitempty"` // Flight safety advance check
}

// PrepareData represents the prepare flight task data
type PrepareData struct {
	BreakPoint            *BreakPoint           `json:"break_point,omitempty"`
	ExecutableConditions  *ExecutableConditions `json:"executable_conditions,omitempty"`
	ExecuteTime           int64                 `json:"execute_time"`
	ExitWaylineWhenRcLost int                   `json:"exit_wayline_when_rc_lost"`
	File                  File                  `json:"file"`
	FlightId              string                `json:"flight_id"`
	OutOfControlAction    int                   `json:"out_of_control_action"`
	ReadyConditions       *ReadyConditions      `json:"ready_conditions,omitempty"`
	RthAltitude           int                   `json:"rth_altitude"`
	SimulateMission       *SimulateMission      `json:"simulate_mission,omitempty"`
	TaskType              int                   `json:"task_type"`
	RthMode               int                   `json:"rth_mode"`
}

type File struct {
	Fingerprint string `json:"fingerprint"`
	Url         string `json:"url"`
}

// ExecuteData represents the execute flight task data
type ExecuteData struct {
	FlighttaskID  string         `json:"flight_id"`                 // Flight task ID
	MultiDockTask *MultiDockTask `json:"multi_dock_task,omitempty"` // Multi-dock task parameters (optional, for leapfrog missions)
}

// UndoData represents the undo flight task data
type UndoData struct {
	FlightIds []string `json:"flight_ids"` // Flight task ID
}

// ReadyConditions represents the task ready conditions
type ReadyConditions struct {
	BatteryCapacity int   `json:"battery_capacity"` // Battery capacity percentage threshold
	BeginTime       int64 `json:"begin_time"`       // Task executable period start time (millisecond timestamp)
	EndTime         int64 `json:"end_time"`         // Task executable period end time (millisecond timestamp)
}

// ExecutableConditions represents the task executable conditions
type ExecutableConditions struct {
	StorageCapacity int `json:"storage_capacity"` // Minimum storage capacity
}

// SimulateMission represents the simulate mission settings
type SimulateMission struct {
	IsEnable  int     `json:"is_enable"` // Enable simulator task (0=disable, 1=enable)
	Latitude  float64 `json:"latitude"`  // Latitude (-90.0 to 90.0)
	Longitude float64 `json:"longitude"` // Longitude (-180.0 to 180.0)
}

// BreakPoint represents the wayline breakpoint information
type BreakPoint struct {
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

// MultiDockTask represents the multi-dock (leapfrog) task parameters
type MultiDockTask struct {
	WirelessLinkTopo WirelessLinkTopo `json:"wireless_link_topo"` // Wireless link topology
	DockInfos        []DockInfo       `json:"dock_infos"`         // Dock information (max 2)
}

// WirelessLinkTopo represents the wireless link topology
type WirelessLinkTopo struct {
	SecretCode []int      `json:"secret_code"` // Encryption code (length 28)
	CenterNode CenterNode `json:"center_node"` // Aircraft pairing information
	LeafNodes  []LeafNode `json:"leaf_nodes"`  // Dock or remote controller pairing information
}

// CenterNode represents the center node (aircraft) in wireless link topology
type CenterNode struct {
	SDRID int    `json:"sdr_id"` // SDR ID (scrambling code)
	SN    string `json:"sn"`     // Device serial number
}

// LeafNode represents the leaf node (dock or remote controller) in wireless link topology
type LeafNode struct {
	SDRID              int    `json:"sdr_id"`               // SDR ID (scrambling code)
	SN                 string `json:"sn"`                   // Device serial number
	ControlSourceIndex int    `json:"control_source_index"` // Control source index (1-2)
}

// DockInfo represents the dock information
type DockInfo struct {
	DockType            string             `json:"dock_type"`              // Dock role (takeoff=takeoff dock, landing=landing dock)
	Latitude            float64            `json:"latitude"`               // Dock latitude
	Longitude           float64            `json:"longitude"`              // Dock longitude
	Height              float64            `json:"height"`                 // Dock ellipsoid height
	Heading             float64            `json:"heading"`                // Dock heading angle (-180 to 180 degrees)
	HomePositionIsValid int                `json:"home_position_is_valid"` // Home point validity (0=invalid, 1=valid)
	Index               int                `json:"index"`                  // Dock task unique identifier (1-31)
	SN                  string             `json:"sn"`                     // Dock serial number
	RTCMInfo            RTCMInfo           `json:"rtcm_info"`              // Dock RTK calibration source
	AlternateLandPoint  AlternateLandPoint `json:"alternate_land_point"`   // Alternate landing point
}

// RTCMInfo represents the RTCM information for dock
type RTCMInfo struct {
	MountPoint     string `json:"mount_point"`      // Network RTK mount point
	Port           string `json:"port"`             // Network port
	Host           string `json:"host"`             // Network host
	RTCMDeviceType int    `json:"rtcm_device_type"` // Device type (1=dock)
	SourceType     int    `json:"source_type"`      // Calibration type (0=uncalibrated, 1=self-converged, 2=manual, 3=network RTK)
}

// AlternateLandPoint represents the alternate landing point
type AlternateLandPoint struct {
	Longitude      float64 `json:"longitude"`        // Longitude
	Latitude       float64 `json:"latitude"`         // Latitude
	Height         float64 `json:"height"`           // Ellipsoid height
	SafeLandHeight float64 `json:"safe_land_height"` // Safe landing height (alternate transfer height)
	IsConfigured   int     `json:"is_configured"`    // Whether alternate landing point is configured (0=no, 1=yes)
}
