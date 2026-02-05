package router

// Event method names.
const (
	MethodDeviceExitHomingNotify  = "device_exit_homing_notify"
	MethodDeviceTempNtfyNeedClear = "device_temp_ntfy_need_clear"
	MethodFileUploadCallback      = "file_upload_callback"
	MethodHMS                     = "hms"
	MethodFlighttaskProgress      = "flighttask_progress"
	MethodFlighttaskReady         = "flighttask_ready"
	MethodReturnHomeInfo          = "return_home_info"
	MethodControlSourceChange     = "control_source_change"
	MethodFlyToPointProgress      = "fly_to_point_progress"
	MethodTakeoffToPointProgress  = "takeoff_to_point_progress"
	MethodDRCStatusNotify         = "drc_status_notify"
	MethodJoystickInvalidNotify   = "joystick_invalid_notify"
	MethodOTAProgress             = "ota_progress"
	MethodFileUploadProgress      = "file_upload_progress"
	MethodHighestPriorityUpload   = "highest_priority_upload_flighttask_media"
)

// HMSData represents HMS (Health Management System) event data.
type HMSData struct {
	List []HMSItem `json:"list"`
}

// HMSItem represents a single HMS item.
type HMSItem struct {
	Code     string                 `json:"code"`
	Level    int                    `json:"level"`    // 0=notice, 1=caution, 2=warning
	Module   int                    `json:"module"`   // Module ID
	InTheSky int                    `json:"in_the_sky"` // 0=ground, 1=flying
	Args     map[string]interface{} `json:"args,omitempty"`
}

// FileUploadCallbackData represents file upload callback event data.
type FileUploadCallbackData struct {
	File FileInfo `json:"file"`
}

// FileInfo represents file information.
type FileInfo struct {
	Path        string `json:"path"`
	Name        string `json:"name"`
	Size        int64  `json:"size"`
	Fingerprint string `json:"fingerprint"`
}

// FlighttaskProgressData represents flight task progress event data.
type FlighttaskProgressData struct {
	FlightID string `json:"flight_id"`
	Status   string `json:"status"`
	Progress int    `json:"progress"` // 0-100
	Result   int    `json:"result,omitempty"`
}

// FlighttaskReadyData represents flight task ready event data.
type FlighttaskReadyData struct {
	FlightID string `json:"flight_id"`
}

// ReturnHomeInfoData represents return home info event data.
type ReturnHomeInfoData struct {
	PlannedPathPoints []PathPoint `json:"planned_path_points,omitempty"`
	LastPointType     int         `json:"last_point_type,omitempty"`
	FlightID          string      `json:"flight_id,omitempty"`
}

// PathPoint represents a path point.
type PathPoint struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Height    float64 `json:"height"`
}

// RegisterCoreEvents registers all core event handlers to the router.
// Returns an error if any handler registration fails.
func RegisterCoreEvents(r *EventRouter) error {
	handlers := map[string]EventHandlerFunc{
		// HMS event
		MethodHMS: SimpleEventHandler[HMSData](MethodHMS),

		// File events
		MethodFileUploadCallback:    SimpleEventHandler[FileUploadCallbackData](MethodFileUploadCallback),
		MethodFileUploadProgress:    NoDataEventHandler(MethodFileUploadProgress),
		MethodHighestPriorityUpload: NoDataEventHandler(MethodHighestPriorityUpload),

		// Device events
		MethodDeviceExitHomingNotify:  NoDataEventHandler(MethodDeviceExitHomingNotify),
		MethodDeviceTempNtfyNeedClear: NoDataEventHandler(MethodDeviceTempNtfyNeedClear),

		// Flight task events
		MethodFlighttaskProgress: SimpleEventHandler[FlighttaskProgressData](MethodFlighttaskProgress),
		MethodFlighttaskReady:    SimpleEventHandler[FlighttaskReadyData](MethodFlighttaskReady),
		MethodReturnHomeInfo:     SimpleEventHandler[ReturnHomeInfoData](MethodReturnHomeInfo),

		// Control events
		MethodControlSourceChange:    NoDataEventHandler(MethodControlSourceChange),
		MethodFlyToPointProgress:     NoDataEventHandler(MethodFlyToPointProgress),
		MethodTakeoffToPointProgress: NoDataEventHandler(MethodTakeoffToPointProgress),

		// DRC events
		MethodDRCStatusNotify:       NoDataEventHandler(MethodDRCStatusNotify),
		MethodJoystickInvalidNotify: NoDataEventHandler(MethodJoystickInvalidNotify),

		// OTA events
		MethodOTAProgress: NoDataEventHandler(MethodOTAProgress),
	}

	return RegisterEventHandlers(r, handlers)
}
