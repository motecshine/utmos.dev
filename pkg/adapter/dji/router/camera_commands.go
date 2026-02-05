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

// IRMeteringPointData represents IR metering point data.
// Note: This type is specific to router and not in protocol/camera.
type IRMeteringPointData struct {
	PayloadIndex string  `json:"payload_index"`
	X            float64 `json:"x"`
	Y            float64 `json:"y"`
}

// IRMeteringAreaData represents IR metering area data.
// Note: This type is specific to router and not in protocol/camera.
type IRMeteringAreaData struct {
	PayloadIndex string  `json:"payload_index"`
	X            float64 `json:"x"`
	Y            float64 `json:"y"`
	Width        float64 `json:"width"`
	Height       float64 `json:"height"`
}

// RegisterCameraCommands registers all camera commands to the router.
// Returns an error if any handler registration fails.
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

		// IR metering commands - using local types as they're not in protocol/camera
		MethodIRMeteringPoint: SimpleCommandHandler[IRMeteringPointData](MethodIRMeteringPoint),
		MethodIRMeteringArea:  SimpleCommandHandler[IRMeteringAreaData](MethodIRMeteringArea),
	}

	return RegisterHandlers(r, handlers)
}
