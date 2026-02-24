package router

import (
	"github.com/utmos/utmos/pkg/adapter/dji/protocol/camera"
)

// Camera command method names.
const (
	MethodCameraModeSwitch       = "camera_mode_switch"
	MethodCameraPhotoTake        = "camera_photo_take"
	MethodCameraRecordingStart   = "camera_recording_start"
	MethodCameraRecordingStop    = "camera_recording_stop"
	MethodCameraAim              = "camera_aim"
	MethodCameraFocalLengthSet   = "camera_focal_length_set"
	MethodGimbalReset            = "gimbal_reset"
	MethodCameraPointFocusAction = "camera_point_focus_action"
	MethodCameraScreenSplit      = "camera_screen_split"
	MethodIRMeteringPoint        = "ir_metering_point"
	MethodIRMeteringArea         = "ir_metering_area"
)

// RegisterCameraCommands registers all camera commands to the router.
// Returns an error if any handler registration fails.
//
// Registration files share structural pattern but register different command types
func RegisterCameraCommands(r *ServiceRouter) error {
	handlers := map[string]ServiceHandlerFunc{
		// Commands using types from protocol/camera package
		MethodCameraModeSwitch:       SimpleCommandHandler[camera.CameraModeSwitchData](MethodCameraModeSwitch),
		MethodCameraPhotoTake:        SimpleCommandHandler[camera.CameraPhotoTakeData](MethodCameraPhotoTake),
		MethodCameraRecordingStart:   SimpleCommandHandler[camera.CameraRecordingStartData](MethodCameraRecordingStart),
		MethodCameraRecordingStop:    SimpleCommandHandler[camera.CameraRecordingStopData](MethodCameraRecordingStop),
		MethodCameraAim:              SimpleCommandHandler[camera.CameraAimData](MethodCameraAim),
		MethodCameraFocalLengthSet:   SimpleCommandHandler[camera.CameraFocalLengthSetData](MethodCameraFocalLengthSet),
		MethodGimbalReset:            SimpleCommandHandler[camera.GimbalResetData](MethodGimbalReset),
		MethodCameraPointFocusAction: SimpleCommandHandler[camera.CameraPointFocusActionData](MethodCameraPointFocusAction),
		MethodCameraScreenSplit:      SimpleCommandHandler[camera.CameraScreenSplitData](MethodCameraScreenSplit),

		// IR metering commands - use types from protocol/camera
		MethodIRMeteringPoint: SimpleCommandHandler[camera.IRMeteringPointSetData](MethodIRMeteringPoint),
		MethodIRMeteringArea:  SimpleCommandHandler[camera.IRMeteringAreaSetData](MethodIRMeteringArea),
	}

	return RegisterHandlers(r, handlers)
}
