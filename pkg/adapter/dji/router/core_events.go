package router

import (
	"github.com/utmos/utmos/pkg/adapter/dji/protocol/aircraft"
)

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

// RegisterCoreEvents registers all core event handlers to the router.
// Returns an error if any handler registration fails.
func RegisterCoreEvents(r *EventRouter) error {
	handlers := map[string]EventHandlerFunc{
		// HMS event - use type from protocol/aircraft
		MethodHMS: SimpleEventHandler[aircraft.HMSData](MethodHMS),

		// Device events
		MethodDeviceExitHomingNotify:  NoDataEventHandler(MethodDeviceExitHomingNotify),
		MethodDeviceTempNtfyNeedClear: NoDataEventHandler(MethodDeviceTempNtfyNeedClear),

		// Control events
		MethodControlSourceChange:    NoDataEventHandler(MethodControlSourceChange),
		MethodFlyToPointProgress:     NoDataEventHandler(MethodFlyToPointProgress),
		MethodTakeoffToPointProgress: NoDataEventHandler(MethodTakeoffToPointProgress),
	}

	return RegisterEventHandlers(r, handlers)
}
