package aircraft

// ===============================
// DJI Aircraft OSD and State Data
// ===============================
// Important: All fields are pointers because OSD and State messages
// contain PARTIAL UPDATES - not all fields are sent in every message.
// This allows for efficient bandwidth usage and proper JSON unmarshaling
// of partial data.

// PayloadInfo represents payload status information
type PayloadInfo struct {
	ControlSource   *string `json:"control_source,omitempty"`   // Payload control authority (A/B or browser UUID)
	PayloadIndex    *string `json:"payload_index,omitempty"`    // Payload index (format: type-subtype-gimbalindex)
	FirmwareVersion *string `json:"firmware_version,omitempty"` // Firmware version
	SN              *string `json:"sn,omitempty"`               // Payload serial number
}

// DistanceLimitStatus represents aircraft distance limit status
type DistanceLimitStatus struct {
	State               *int `json:"state,omitempty"`                  // Whether distance limit is enabled (0=not set, 1=set)
	DistanceLimit       *int `json:"distance_limit,omitempty"`         // Distance limit (m, 15-8000)
	IsNearDistanceLimit *int `json:"is_near_distance_limit,omitempty"` // Whether near distance limit (0=no, 1=yes)
}

// LiveviewWorldRegion represents the field of view in liveview
type LiveviewWorldRegion struct {
	Left   *float64 `json:"left,omitempty"`   // Left x coordinate
	Top    *float64 `json:"top,omitempty"`    // Top y coordinate
	Right  *float64 `json:"right,omitempty"`  // Right x coordinate
	Bottom *float64 `json:"bottom,omitempty"` // Bottom y coordinate
}

// IRMeteringPoint represents infrared temperature point
type IRMeteringPoint struct {
	X           *float64 `json:"x,omitempty"`           // Point x coordinate (0-1)
	Y           *float64 `json:"y,omitempty"`           // Point y coordinate (0-1)
	Temperature *float64 `json:"temperature,omitempty"` // Point temperature (°C)
}

// IRMeteringAreaPoint represents infrared temperature area point
type IRMeteringAreaPoint struct {
	X           *float64 `json:"x,omitempty"`           // Point x coordinate (0-1)
	Y           *float64 `json:"y,omitempty"`           // Point y coordinate (0-1)
	Temperature *float64 `json:"temperature,omitempty"` // Point temperature (°C)
}

// IRMeteringArea represents infrared temperature area
type IRMeteringArea struct {
	X                   *float64             `json:"x,omitempty"`                     // Area top-left x coordinate (0-1)
	Y                   *float64             `json:"y,omitempty"`                     // Area top-left y coordinate (0-1)
	Width               *float64             `json:"width,omitempty"`                 // Area width (0-1)
	Height              *float64             `json:"height,omitempty"`                // Area height (0-1)
	AverTemperature     *float64             `json:"aver_temperature,omitempty"`      // Average temperature (°C)
	MinTemperaturePoint *IRMeteringAreaPoint `json:"min_temperature_point,omitempty"` // Minimum temperature point
	MaxTemperaturePoint *IRMeteringAreaPoint `json:"max_temperature_point,omitempty"` // Maximum temperature point
}

// CameraInfo represents camera information
type CameraInfo struct {
	RemainPhotoNum                  *int                 `json:"remain_photo_num,omitempty"`                    // Remaining photo count
	RemainRecordDuration            *int                 `json:"remain_record_duration,omitempty"`              // Remaining recording time (s)
	RecordTime                      *int                 `json:"record_time,omitempty"`                         // Recording duration (s)
	PayloadIndex                    *string              `json:"payload_index,omitempty"`                       // Payload index (format: type-subtype-gimbalindex)
	CameraMode                      *int                 `json:"camera_mode,omitempty"`                         // Camera mode (0=photo, 1=video, 2=smart low light, 3=panorama, -1=unsupported)
	PhotoState                      *int                 `json:"photo_state,omitempty"`                         // Photo state (0=idle, 1=shooting)
	ScreenSplitEnable               *bool                `json:"screen_split_enable,omitempty"`                 // Screen split enable
	RecordingState                  *int                 `json:"recording_state,omitempty"`                     // Recording state (0=idle, 1=recording)
	ZoomFactor                      *float64             `json:"zoom_factor,omitempty"`                         // Zoom factor (2-200)
	IRZoomFactor                    *float64             `json:"ir_zoom_factor,omitempty"`                      // IR zoom factor (2-20)
	LiveviewWorldRegion             *LiveviewWorldRegion `json:"liveview_world_region,omitempty"`               // FOV region in liveview
	PhotoStorageSettings            []string             `json:"photo_storage_settings,omitempty"`              // Photo storage settings (current, wide, zoom, ir)
	VideoStorageSettings            []string             `json:"video_storage_settings,omitempty"`              // Video storage settings (current, wide, zoom, ir)
	WideExposureMode                *int                 `json:"wide_exposure_mode,omitempty"`                  // Wide lens exposure mode (1=auto, 2=shutter priority, 3=aperture priority, 4=manual)
	WideISO                         *int                 `json:"wide_iso,omitempty"`                            // Wide lens ISO
	WideShutterSpeed                *int                 `json:"wide_shutter_speed,omitempty"`                  // Wide lens shutter speed
	WideExposureValue               *int                 `json:"wide_exposure_value,omitempty"`                 // Wide lens exposure value
	ZoomExposureMode                *int                 `json:"zoom_exposure_mode,omitempty"`                  // Zoom lens exposure mode
	ZoomISO                         *int                 `json:"zoom_iso,omitempty"`                            // Zoom lens ISO
	ZoomShutterSpeed                *int                 `json:"zoom_shutter_speed,omitempty"`                  // Zoom lens shutter speed
	ZoomExposureValue               *int                 `json:"zoom_exposure_value,omitempty"`                 // Zoom lens exposure value
	ZoomFocusMode                   *int                 `json:"zoom_focus_mode,omitempty"`                     // Zoom lens focus mode (0=MF, 1=AFS, 2=AFC)
	ZoomFocusValue                  *int                 `json:"zoom_focus_value,omitempty"`                    // Zoom lens focus value
	ZoomMaxFocusValue               *int                 `json:"zoom_max_focus_value,omitempty"`                // Zoom lens max focus value
	ZoomMinFocusValue               *int                 `json:"zoom_min_focus_value,omitempty"`                // Zoom lens min focus value
	ZoomCalibrateFarthestFocusValue *int                 `json:"zoom_calibrate_farthest_focus_value,omitempty"` // Zoom lens calibrated farthest focus value
	ZoomCalibrateNearestFocusValue  *int                 `json:"zoom_calibrate_nearest_focus_value,omitempty"`  // Zoom lens calibrated nearest focus value
	ZoomFocusState                  *int                 `json:"zoom_focus_state,omitempty"`                    // Zoom lens focus state (0=idle, 1=focusing, 2=success, 3=failed)
	IRMeteringMode                  *int                 `json:"ir_metering_mode,omitempty"`                    // IR metering mode (0=off, 1=point, 2=area)
	IRMeteringPoint                 *IRMeteringPoint     `json:"ir_metering_point,omitempty"`                   // IR metering point
	IRMeteringArea                  *IRMeteringArea      `json:"ir_metering_area,omitempty"`                    // IR metering area
}

// BatteryDetail represents single battery details
//
// Struct shape similar to WirelessLink but fields are semantically different
type BatteryDetail struct {
	CapacityPercent        *int     `json:"capacity_percent,omitempty"`          // Battery remaining capacity (0-100)
	Index                  *int     `json:"index,omitempty"`                     // Battery index (0+)
	SN                     *string  `json:"sn,omitempty"`                        // Battery serial number
	Type                   *int     `json:"type,omitempty"`                      // Battery type
	SubType                *int     `json:"sub_type,omitempty"`                  // Battery sub-type
	FirmwareVersion        *string  `json:"firmware_version,omitempty"`          // Firmware version
	LoopTimes              *int     `json:"loop_times,omitempty"`                // Battery cycle count
	Voltage                *int     `json:"voltage,omitempty"`                   // Voltage (mV)
	Temperature            *float64 `json:"temperature,omitempty"`               // Temperature (°C)
	HighVoltageStorageDays *int     `json:"high_voltage_storage_days,omitempty"` // High voltage storage days
}

// BatteryInfo represents aircraft battery information
type BatteryInfo struct {
	CapacityPercent  *int            `json:"capacity_percent,omitempty"`   // Total remaining battery capacity (0-100)
	RemainFlightTime *int            `json:"remain_flight_time,omitempty"` // Remaining flight time (s)
	ReturnHomePower  *int            `json:"return_home_power,omitempty"`  // Return home required battery percentage (0-100)
	LandingPower     *int            `json:"landing_power,omitempty"`      // Forced landing battery percentage (0-100)
	Batteries        []BatteryDetail `json:"batteries,omitempty"`          // Battery details
}

// ObstacleAvoidance represents obstacle avoidance status
type ObstacleAvoidance struct {
	Horizon  *int `json:"horizon,omitempty"`  // Horizontal obstacle avoidance (0=off, 1=on)
	Upside   *int `json:"upside,omitempty"`   // Upward obstacle avoidance (0=off, 1=on)
	Downside *int `json:"downside,omitempty"` // Downward obstacle avoidance (0=off, 1=on)
}

// MaintainStatusItemAircraft represents aircraft maintenance status item
type MaintainStatusItemAircraft struct {
	State                     *int `json:"state,omitempty"`                        // Maintenance state (0=no maintenance, 1=needs maintenance)
	LastMaintainType          *int `json:"last_maintain_type,omitempty"`           // Last maintenance type (1=basic, 2=regular, 3=deep)
	LastMaintainTime          *int `json:"last_maintain_time,omitempty"`           // Last maintenance time (unix timestamp)
	LastMaintainFlightTime    *int `json:"last_maintain_flight_time,omitempty"`    // Last maintenance flight time (hours)
	LastMaintainFlightSorties *int `json:"last_maintain_flight_sorties,omitempty"` // Last maintenance flight sorties count
}

// MaintainStatusAircraft represents aircraft maintenance information
type MaintainStatusAircraft struct {
	MaintainStatusArray []MaintainStatusItemAircraft `json:"maintain_status_array,omitempty"` // Maintenance status array
}

// GimbalInfo represents gimbal and payload information (dynamic key: type-subtype-gimbalindex)
type GimbalInfo struct {
	GimbalPitch                   *float64 `json:"gimbal_pitch,omitempty"`                     // Gimbal pitch angle (-180 to 180)
	GimbalRoll                    *float64 `json:"gimbal_roll,omitempty"`                      // Gimbal roll angle (-180 to 180)
	GimbalYaw                     *float64 `json:"gimbal_yaw,omitempty"`                       // Gimbal yaw angle (-180 to 180)
	MeasureTargetLongitude        *float64 `json:"measure_target_longitude,omitempty"`         // Laser ranging target longitude (-180 to 180)
	MeasureTargetLatitude         *float64 `json:"measure_target_latitude,omitempty"`          // Laser ranging target latitude (-90 to 90)
	MeasureTargetAltitude         *float64 `json:"measure_target_altitude,omitempty"`          // Laser ranging target altitude (m)
	MeasureTargetDistance         *float64 `json:"measure_target_distance,omitempty"`          // Laser ranging distance (m)
	MeasureTargetErrorState       *int     `json:"measure_target_error_state,omitempty"`       // Laser ranging status (0=normal, 1=too close, 2=too far, 3=no signal)
	PayloadIndex                  *string  `json:"payload_index,omitempty"`                    // Payload index (format: type-subtype-gimbalindex)
	ZoomFactor                    *float64 `json:"zoom_factor,omitempty"`                      // Zoom factor
	ThermalCurrentPaletteStyle    *int     `json:"thermal_current_palette_style,omitempty"`    // Thermal palette style (0-13)
	ThermalSupportedPaletteStyles []int    `json:"thermal_supported_palette_styles,omitempty"` // Supported thermal palette styles
	ThermalGainMode               *int     `json:"thermal_gain_mode,omitempty"`                // Thermal gain mode (0=auto, 1=low gain, 2=high gain)
	ThermalIsothermState          *int     `json:"thermal_isotherm_state,omitempty"`           // Thermal isotherm enable (0=off, 1=on)
	ThermalIsothermUpperLimit     *int     `json:"thermal_isotherm_upper_limit,omitempty"`     // Thermal isotherm upper limit (°C)
	ThermalIsothermLowerLimit     *int     `json:"thermal_isotherm_lower_limit,omitempty"`     // Thermal isotherm lower limit (°C)
	ThermalGlobalTemperatureMin   *float64 `json:"thermal_global_temperature_min,omitempty"`   // Global min temperature (°C)
	ThermalGlobalTemperatureMax   *float64 `json:"thermal_global_temperature_max,omitempty"`   // Global max temperature (°C)
}

// PSDKUIResource represents PSDK UI resource information
type PSDKUIResource struct {
	PsdkIndex *int    `json:"psdk_index,omitempty"` // PSDK payload device index (0+)
	PsdkReady *int    `json:"psdk_ready,omitempty"` // PSDK ready status (0=not ready, 1=ready)
	ObjectKey *string `json:"object_key,omitempty"` // OSS object key
}

// SpeakerStatus represents speaker status
type SpeakerStatus struct {
	WorkMode     *int    `json:"work_mode,omitempty"`      // Speaker work mode (0=TTS, 1=audio)
	PlayMode     *int    `json:"play_mode,omitempty"`      // Speaker play mode (0=once, 1=loop)
	PlayVolume   *int    `json:"play_volume,omitempty"`    // Speaker volume (0-100)
	SystemState  *int    `json:"system_state,omitempty"`   // Speaker status (0=idle, 1=transmitting, 2=playing, 3=error, 4=TTS converting, 99=downloading)
	PlayFileName *string `json:"play_file_name,omitempty"` // Last played file name
	PlayFileMD5  *string `json:"play_file_md5,omitempty"`  // Last played file MD5
}

// PSDKWidgetValue represents PSDK widget value
type PSDKWidgetValue struct {
	Index *int `json:"index,omitempty"` // Widget index (0+)
	Value *int `json:"value,omitempty"` // Widget value
}

// PSDKWidgetValues represents PSDK payload device widget values
type PSDKWidgetValues struct {
	PsdkIndex      *int              `json:"psdk_index,omitempty"`       // PSDK payload device index (0+)
	PsdkName       *string           `json:"psdk_name,omitempty"`        // Device name
	PsdkSN         *string           `json:"psdk_sn,omitempty"`          // Device serial number
	PsdkVersion    *string           `json:"psdk_version,omitempty"`     // Device firmware version
	PsdkLibVersion *string           `json:"psdk_lib_version,omitempty"` // PSDK lib version
	Speaker        *SpeakerStatus    `json:"speaker,omitempty"`          // Speaker status
	Values         []PSDKWidgetValue `json:"values,omitempty"`           // Widget value list
}

// AircraftOSD represents the DJI Aircraft OSD and State data structure
// All fields are pointers to support partial updates
type AircraftOSD struct {
	// Payload information
	Payloads []PayloadInfo `json:"payloads,omitempty"` // Payload status list

	// Aircraft status
	ModeCode              *int    `json:"mode_code,omitempty"`               // Aircraft status (0-20)
	ModeCodeReason        *int    `json:"mode_code_reason,omitempty"`        // Reason for entering current status (0-23)
	Gear                  *int    `json:"gear,omitempty"`                    // Gear (0=A, 1=P, 2=NAV, 3=FPV, 4=FARM, 5=S, 6=F, 7=M, 8=G, 9=T)
	FirmwareVersion       *string `json:"firmware_version,omitempty"`        // Firmware version
	CompatibleStatus      *int    `json:"compatible_status,omitempty"`       // Firmware consistency (0=no upgrade needed, 1=upgrade needed)
	FirmwareUpgradeStatus *int    `json:"firmware_upgrade_status,omitempty"` // Firmware upgrade status (0=not upgrading, 1=upgrading)

	// Distance and height limits
	DistanceLimitStatus *DistanceLimitStatus `json:"distance_limit_status,omitempty"` // Distance limit status
	IsNearHeightLimit   *int                 `json:"is_near_height_limit,omitempty"`  // Whether near height limit (0=no, 1=yes)
	IsNearAreaLimit     *int                 `json:"is_near_area_limit,omitempty"`    // Whether near no-fly zone (0=no, 1=yes)
	HeightLimit         *int                 `json:"height_limit,omitempty"`          // Height limit (m, 20-1500)

	// Flight parameters
	WpmzVersion                *string  `json:"wpmz_version,omitempty"`                  // Wayline parsing library version
	RTHAltitude                *float64 `json:"rth_altitude,omitempty"`                  // Return home altitude (m, 20-500)
	RCLostAction               *int     `json:"rc_lost_action,omitempty"`                // RC lost action (0=hover, 1=land, 2=RTH)
	ExitWaylineWhenRCLost      *int     `json:"exit_wayline_when_rc_lost,omitempty"`     // [Deprecated] Wayline lost action
	CommanderModeLostAction    *int     `json:"commander_mode_lost_action,omitempty"`    // Commander flight lost action (0=continue, 1=exit and RTH)
	CurrentCommanderFlightMode *int     `json:"current_commander_flight_mode,omitempty"` // Current commander flight mode (0=smart height, 1=set height)
	CommanderFlightHeight      *float64 `json:"commander_flight_height,omitempty"`       // Commander flight height (m, 2-3000)
	CurrentRTHMode             *float64 `json:"current_rth_mode,omitempty"`              // Return home altitude mode (0=smart height, 1=set height)

	// Camera information
	Cameras []CameraInfo `json:"cameras,omitempty"` // Aircraft camera information

	// Location and navigation
	Country         *string  `json:"country,omitempty"`          // Country region code
	RIDState        *bool    `json:"rid_state,omitempty"`        // RID working status (false=abnormal, true=normal)
	HorizontalSpeed *float64 `json:"horizontal_speed,omitempty"` // Horizontal speed (m/s)
	VerticalSpeed   *float64 `json:"vertical_speed,omitempty"`   // Vertical speed (m/s)
	Longitude       *float64 `json:"longitude,omitempty"`        // Current longitude
	Latitude        *float64 `json:"latitude,omitempty"`         // Current latitude
	Height          *float64 `json:"height,omitempty"`           // Absolute height (m, ellipsoid height)
	Elevation       *float64 `json:"elevation,omitempty"`        // Relative takeoff point height (m)
	AttitudePitch   *float64 `json:"attitude_pitch,omitempty"`   // Pitch angle (degrees)
	AttitudeRoll    *float64 `json:"attitude_roll,omitempty"`    // Roll angle (degrees)
	AttitudeHead    *float64 `json:"attitude_head,omitempty"`    // Yaw angle (degrees, relative to true north)
	HomeLongitude   *float64 `json:"home_longitude,omitempty"`   // Home point longitude
	HomeLatitude    *float64 `json:"home_latitude,omitempty"`    // Home point latitude
	HomeDistance    *float64 `json:"home_distance,omitempty"`    // Distance to home point (m)

	// Wind information
	WindSpeed     *float64 `json:"wind_speed,omitempty"`     // Wind speed (m/s, estimated)
	WindDirection *int     `json:"wind_direction,omitempty"` // Wind direction (1-8: N, NE, E, SE, S, SW, W, NW)

	// Control
	ControlSource *string `json:"control_source,omitempty"` // Current control source (A/B or browser UUID)

	// Battery warnings
	LowBatteryWarningThreshold        *int `json:"low_battery_warning_threshold,omitempty"`         // Low battery warning threshold (%)
	SeriousLowBatteryWarningThreshold *int `json:"serious_low_battery_warning_threshold,omitempty"` // Serious low battery warning threshold (%)

	// Statistics
	TotalFlightTime     *float64 `json:"total_flight_time,omitempty"`     // Total flight time (s)
	TotalFlightDistance *float64 `json:"total_flight_distance,omitempty"` // Total flight distance (m)
	TotalFlightSorties  *float64 `json:"total_flight_sorties,omitempty"`  // Total flight sorties
	ActivationTime      *int     `json:"activation_time,omitempty"`       // Activation time (unix timestamp)

	// Battery and storage
	Battery *BatteryInfo `json:"battery,omitempty"` // Aircraft battery information
	Storage *Storage     `json:"storage,omitempty"` // Storage capacity (KB)

	// Positioning
	PositionState *PositionState `json:"position_state,omitempty"` // Satellite positioning status
	TrackID       *string        `json:"track_id,omitempty"`       // Track ID

	// Gimbal information (dynamic key: type-subtype-gimbalindex)
	// Note: In Go, we cannot have dynamic keys in struct, so gimbal info should be handled separately
	// or use a map[string]*GimbalInfo for dynamic payload indexes

	// Maintenance
	MaintainStatus *MaintainStatusAircraft `json:"maintain_status,omitempty"` // Maintenance information

	// Night lights and obstacle avoidance
	NightLightsState  *int               `json:"night_lights_state,omitempty"` // Night lights status (0=off, 1=on)
	ObstacleAvoidance *ObstacleAvoidance `json:"obstacle_avoidance,omitempty"` // Obstacle avoidance status

	// PSDK
	PsdkUIResource   []PSDKUIResource   `json:"psdk_ui_resource,omitempty"`   // PSDK UI resource list
	PsdkWidgetValues []PSDKWidgetValues `json:"psdk_widget_values,omitempty"` // PSDK widget values
}
