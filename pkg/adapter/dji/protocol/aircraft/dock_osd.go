package aircraft

// ===============================
// DJI Dock OSD and State Data
// ===============================
// Important: All fields are pointers because OSD and State messages
// contain PARTIAL UPDATES - not all fields are sent in every message.
// This allows for efficient bandwidth usage and proper JSON unmarshaling
// of partial data.

// SubDevice represents the sub-device (aircraft) status
type SubDevice struct {
	DeviceSN           *string `json:"device_sn,omitempty"`            // Sub-device serial number
	ProductType        *string `json:"product_type,omitempty"`         // Sub-device product type (format: domain-type-subtype)
	DeviceOnlineStatus *int    `json:"device_online_status,omitempty"` // Aircraft power status (0=off, 1=on)
	DevicePaired       *int    `json:"device_paired,omitempty"`        // Whether aircraft is paired with dock (0=unpaired, 1=paired)
}

// NetworkState represents the network status
type NetworkState struct {
	Type    *int     `json:"type,omitempty"`    // Network type (1=4G, 2=Ethernet)
	Quality *int     `json:"quality,omitempty"` // Network quality (0-5)
	Rate    *float64 `json:"rate,omitempty"`    // Network rate (KB/s)
}

// MediaFileDetail represents media file upload details
type MediaFileDetail struct {
	RemainUpload *int `json:"remain_upload,omitempty"` // Remaining upload count
}

// WirelessLink represents the wireless link topology
//
// Struct shape similar to BatteryDetail but fields are semantically different
type WirelessLink struct {
	DongleNumber   *int     `json:"dongle_number,omitempty"`  // Dongle count on aircraft
	Link4GState    *int     `json:"4g_link_state,omitempty"`  // 4G link state (0=disconnected, 1=connected)
	SDRLinkState   *int     `json:"sdr_link_state,omitempty"` // SDR link state (0=disconnected, 1=connected)
	LinkWorkmode   *int     `json:"link_workmode,omitempty"`  // Link work mode (0=SDR, 1=4G fusion)
	SDRQuality     *int     `json:"sdr_quality,omitempty"`    // SDR signal quality (0-5)
	Link4GQuality  *int     `json:"4g_quality,omitempty"`     // Overall 4G signal quality (0-5)
	UAV4GQuality   *int     `json:"4g_uav_quality,omitempty"` // UAV 4G signal quality (0-5)
	GND4GQuality   *int     `json:"4g_gnd_quality,omitempty"` // Ground 4G signal quality (0-5)
	SDRFreqBand    *float64 `json:"sdr_freq_band,omitempty"`  // SDR frequency band
	Link4GFreqBand *float64 `json:"4g_freq_band,omitempty"`   // 4G frequency band
}

// LiveStatusItem represents a single live stream status
type LiveStatusItem struct {
	VideoID      *string `json:"video_id,omitempty"`      // Video stream identifier (format: sn/camera_index/video_index)
	VideoType    *string `json:"video_type,omitempty"`    // Video type (normal/wide/zoom/infrared)
	VideoQuality *int    `json:"video_quality,omitempty"` // Video quality (0=adaptive, 1=smooth, 2=SD, 3=HD, 4=UHD)
	Status       *int    `json:"status,omitempty"`        // Live status (0=not live, 1=live)
	ErrorStatus  *int    `json:"error_status,omitempty"`  // Error code
}

// CameraListItem represents a camera in the device list
type CameraListItem struct {
	CameraIndex           *string     `json:"camera_index,omitempty"`             // Camera index (format: type-subtype-gimbalindex)
	AvailableVideoNumber  *int        `json:"available_video_number,omitempty"`   // Available video stream count
	CoexistVideoNumberMax *int        `json:"coexist_video_number_max,omitempty"` // Maximum concurrent video stream count
	VideoList             []VideoItem `json:"video_list,omitempty"`               // Video stream list
}

// VideoItem represents a video stream option
type VideoItem struct {
	VideoIndex           *string  `json:"video_index,omitempty"`            // Video stream index
	VideoType            *string  `json:"video_type,omitempty"`             // Video stream type
	SwitchableVideoTypes []string `json:"switchable_video_types,omitempty"` // Switchable video types
}

// DeviceListItem represents a video source device
type DeviceListItem struct {
	SN                    *string          `json:"sn,omitempty"`                       // Device serial number
	AvailableVideoNumber  *int             `json:"available_video_number,omitempty"`   // Available video stream count
	CoexistVideoNumberMax *int             `json:"coexist_video_number_max,omitempty"` // Maximum concurrent video stream count
	CameraList            []CameraListItem `json:"camera_list,omitempty"`              // Camera list
}

// LiveCapacity represents the live streaming capability
type LiveCapacity struct {
	AvailableVideoNumber  *int             `json:"available_video_number,omitempty"`   // Available video stream count
	CoexistVideoNumberMax *int             `json:"coexist_video_number_max,omitempty"` // Maximum concurrent video stream count
	DeviceList            []DeviceListItem `json:"device_list,omitempty"`              // Device list
}

// Storage represents storage capacity information
type Storage struct {
	Total float64 `json:"total,omitempty"` // Total capacity (KB)
	Used  float64 `json:"used,omitempty"`  // Used capacity (KB)
}

// AlternateLandPoint represents the alternate landing point
type AlternateLandPoint struct {
	Longitude      *float64 `json:"longitude,omitempty"`        // Longitude
	Latitude       *float64 `json:"latitude,omitempty"`         // Latitude
	SafeLandHeight *float64 `json:"safe_land_height,omitempty"` // Safe landing height
	IsConfigured   *int     `json:"is_configured,omitempty"`    // Whether configured (0=no, 1=yes)
	Height         *float64 `json:"height,omitempty"`           // Ellipsoid height
}

// BackupBattery represents backup battery information
type BackupBattery struct {
	Switch      *int     `json:"switch,omitempty"`      // Backup battery switch (0=off, 1=on)
	Voltage     *int     `json:"voltage,omitempty"`     // Backup battery voltage (mV, 0 when off)
	Temperature *float64 `json:"temperature,omitempty"` // Backup battery temperature (°C)
}

// DroneChargeState represents aircraft charging status
type DroneChargeState struct {
	CapacityPercent *int `json:"capacity_percent,omitempty"` // Battery percentage (0-100)
	State           *int `json:"state,omitempty"`            // Charging state (0=idle, 1=charging)
}

// PositionState represents positioning status
type PositionState struct {
	IsCalibration *int `json:"is_calibration,omitempty"` // Whether calibrated (0=no, 1=yes)
	IsFixed       *int `json:"is_fixed,omitempty"`       // Convergence status (0=not started, 1=converging, 2=success, 3=failed)
	Quality       *int `json:"quality,omitempty"`        // Positioning quality level (1-5, 10=RTK fixed)
	GPSNumber     *int `json:"gps_number,omitempty"`     // GPS satellite count
	RTKNumber     *int `json:"rtk_number,omitempty"`     // RTK satellite count
}

// MaintainStatusItem represents a maintenance status item
type MaintainStatusItem struct {
	State                   *int `json:"state,omitempty"`                      // Maintenance state (0=no maintenance, 1=needs maintenance)
	LastMaintainType        *int `json:"last_maintain_type,omitempty"`         // Last maintenance type
	LastMaintainTime        *int `json:"last_maintain_time,omitempty"`         // Last maintenance time (seconds)
	LastMaintainWorkSorties *int `json:"last_maintain_work_sorties,omitempty"` // Last maintenance work sorties count
}

// MaintainStatus represents maintenance information
type MaintainStatus struct {
	MaintainStatusArray []MaintainStatusItem `json:"maintain_status_array,omitempty"` // Maintenance status array
}

// AirConditioner represents air conditioner working status
type AirConditioner struct {
	AirConditionerState *int `json:"air_conditioner_state,omitempty"` // Air conditioner state (0=idle, 1-9=various modes)
	SwitchTime          *int `json:"switch_time,omitempty"`           // Remaining time to switch mode (seconds)
}

// DroneBatteryInfo represents drone battery detailed information
type DroneBatteryInfo struct {
	CapacityPercent *int     `json:"capacity_percent,omitempty"` // Battery remaining percentage (0-100)
	Index           *int     `json:"index,omitempty"`            // Battery index (0=left, 1=right)
	Voltage         *int     `json:"voltage,omitempty"`          // Voltage (mV)
	Temperature     *float64 `json:"temperature,omitempty"`      // Temperature (°C)
}

// DroneBatteryMaintenanceInfo represents drone battery maintenance information
type DroneBatteryMaintenanceInfo struct {
	MaintenanceState    *int               `json:"maintenance_state,omitempty"`     // Maintenance state (0=no need, 1=pending, 2=in progress)
	MaintenanceTimeLeft *int               `json:"maintenance_time_left,omitempty"` // Maintenance time left (hours)
	HeatState           *int               `json:"heat_state,omitempty"`            // Battery heating state (0=off, 1=heating, 2=warming)
	Batteries           []DroneBatteryInfo `json:"batteries,omitempty"`             // Battery details
}

// PayloadAuthorityInfo represents payload control authority information
type PayloadAuthorityInfo struct {
	ControlSource *string `json:"control_source,omitempty"` // Payload control authority ("A" or "B" or "")
	PayloadIndex  *string `json:"payload_index,omitempty"`  // Payload index (format: type-subtype-gimbalindex)
	SN            *string `json:"sn,omitempty"`             // Payload serial number
}

// DroneAuthorityInfo represents drone control authority status
type DroneAuthorityInfo struct {
	ControlSource *string                `json:"control_source,omitempty"` // Flight control authority ("A" or "B" or "")
	Locked        *bool                  `json:"locked,omitempty"`         // Whether control authority is locked
	Payloads      []PayloadAuthorityInfo `json:"payloads,omitempty"`       // Payload control authority list
}

// DockOSD represents the DJI Dock OSD and State data structure
// All fields are pointers to support partial updates
type DockOSD struct {
	// Basic location and status
	Longitude       *float64 `json:"longitude,omitempty"`        // Longitude (-180 to 180)
	Latitude        *float64 `json:"latitude,omitempty"`         // Latitude (-90 to 90)
	Height          *float64 `json:"height,omitempty"`           // Ellipsoid height (m)
	FirmwareVersion *string  `json:"firmware_version,omitempty"` // Firmware version

	// Status codes
	FirmwareUpgradeStatus *int `json:"firmware_upgrade_status,omitempty"` // Firmware upgrade status (0=not upgrading, 1=upgrading)
	ModeCode              *int `json:"mode_code,omitempty"`               // Dock status (0=idle, 1=field debug, 2=remote debug, 3=firmware upgrading, 4=working)
	FlighttaskStepCode    *int `json:"flighttask_step_code,omitempty"`    // Dock task status

	// Sub-device status
	SubDevice *SubDevice `json:"sub_device,omitempty"` // Sub-device (aircraft) status

	// Physical status
	CoverState           *int `json:"cover_state,omitempty"`            // Cover state (0=closed, 1=open, 2=half-open, 3=abnormal)
	SupplementLightState *int `json:"supplement_light_state,omitempty"` // Supplement light state (0=off, 1=on)
	PutterState          *int `json:"putter_state,omitempty"`           // Putter state (0=closed, 1=open, 2=half-open, 3=abnormal)
	DroneInDock          *int `json:"drone_in_dock,omitempty"`          // Whether aircraft is in dock (0=out, 1=in)
	EmergencyStopState   *int `json:"emergency_stop_state,omitempty"`   // Emergency stop button state (0=off, 1=on)

	// Network and connectivity
	NetworkState *NetworkState `json:"network_state,omitempty"` // Network status
	WirelessLink *WirelessLink `json:"wireless_link,omitempty"` // Wireless link status
	DRCState     *int          `json:"drc_state,omitempty"`     // DRC link state (0=disconnected, 1=connecting, 2=connected)

	// Media and live streaming
	MediaFileDetail *MediaFileDetail `json:"media_file_detail,omitempty"` // Media file upload details
	LiveStatus      []LiveStatusItem `json:"live_status,omitempty"`       // Live streaming status
	LiveCapacity    *LiveCapacity    `json:"live_capacity,omitempty"`     // Live streaming capability

	// Environmental sensors
	Rainfall               *int     `json:"rainfall,omitempty"`                // Rainfall level (0=none, 1=light, 2=moderate, 3=heavy)
	WindSpeed              *float64 `json:"wind_speed,omitempty"`              // Wind speed (m/s)
	EnvironmentTemperature *float64 `json:"environment_temperature,omitempty"` // Environment temperature (°C)
	Temperature            *float64 `json:"temperature,omitempty"`             // Internal temperature (°C)
	Humidity               *float64 `json:"humidity,omitempty"`                // Internal humidity (%RH)

	// Power system
	ElectricSupplyVoltage *int           `json:"electric_supply_voltage,omitempty"` // Mains voltage (V)
	WorkingVoltage        *int           `json:"working_voltage,omitempty"`         // Working voltage (mV)
	WorkingCurrent        *float64       `json:"working_current,omitempty"`         // Working current (mA)
	BackupBattery         *BackupBattery `json:"backup_battery,omitempty"`          // Backup battery information

	// Storage
	Storage *Storage `json:"storage,omitempty"` // Storage capacity

	// Operation statistics
	JobNumber      *int `json:"job_number,omitempty"`      // Cumulative job count
	AccTime        *int `json:"acc_time,omitempty"`        // Cumulative operation time (seconds)
	FirstPowerOn   *int `json:"first_power_on,omitempty"`  // First power-on time (ms timestamp)
	ActivationTime *int `json:"activation_time,omitempty"` // Activation time (unix timestamp)

	// Firmware and compatibility
	CompatibleStatus *int    `json:"compatible_status,omitempty"` // Firmware consistency (0=no upgrade needed, 1=upgrade needed)
	WPMZVersion      *string `json:"wpmz_version,omitempty"`      // Wayline parsing library version

	// Position and landing
	AlternateLandPoint *AlternateLandPoint `json:"alternate_land_point,omitempty"` // Alternate landing point
	PositionState      *PositionState      `json:"position_state,omitempty"`       // Positioning status

	// Battery management
	BatteryStoreMode            *int                         `json:"battery_store_mode,omitempty"`             // Battery operation mode (1=plan, 2=standby)
	DroneChargeState            *DroneChargeState            `json:"drone_charge_state,omitempty"`             // Aircraft charging status
	DroneBatteryMaintenanceInfo *DroneBatteryMaintenanceInfo `json:"drone_battery_maintenance_info,omitempty"` // Aircraft battery maintenance info

	// Maintenance
	MaintainStatus *MaintainStatus `json:"maintain_status,omitempty"` // Maintenance information

	// Control and settings
	AlarmState                *int                `json:"alarm_state,omitempty"`                 // Alarm state (0=off, 1=on)
	AirConditioner            *AirConditioner     `json:"air_conditioner,omitempty"`             // Air conditioner status
	UserExperienceImprovement *int                `json:"user_experience_improvement,omitempty"` // User experience improvement plan (0=initial, 1=declined, 2=agreed)
	SilentMode                *int                `json:"silent_mode,omitempty"`                 // Silent mode (0=off, 1=on)
	DroneAuthorityInfo        *DroneAuthorityInfo `json:"drone_authority_info,omitempty"`        // Drone control authority status
}
