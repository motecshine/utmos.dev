package mocks

import "encoding/json"

// CameraModeSwitch is a mock camera_mode_switch service request.
var CameraModeSwitch = json.RawMessage(`{
	"tid": "cam-tid-001",
	"bid": "cam-bid-001",
	"timestamp": 1706000000000,
	"method": "camera_mode_switch",
	"data": {
		"payload_index": "39-0-7",
		"camera_mode": 0
	}
}`)

// CameraPhotoTake is a mock camera_photo_take service request.
var CameraPhotoTake = json.RawMessage(`{
	"tid": "cam-tid-002",
	"bid": "cam-bid-002",
	"timestamp": 1706000000000,
	"method": "camera_photo_take",
	"data": {
		"payload_index": "39-0-7"
	}
}`)

// CameraRecordingStart is a mock camera_recording_start service request.
var CameraRecordingStart = json.RawMessage(`{
	"tid": "cam-tid-003",
	"bid": "cam-bid-003",
	"timestamp": 1706000000000,
	"method": "camera_recording_start",
	"data": {
		"payload_index": "39-0-7"
	}
}`)

// CameraRecordingStop is a mock camera_recording_stop service request.
var CameraRecordingStop = json.RawMessage(`{
	"tid": "cam-tid-004",
	"bid": "cam-bid-004",
	"timestamp": 1706000000000,
	"method": "camera_recording_stop",
	"data": {
		"payload_index": "39-0-7"
	}
}`)

// CameraAim is a mock camera_aim service request.
var CameraAim = json.RawMessage(`{
	"tid": "cam-tid-005",
	"bid": "cam-bid-005",
	"timestamp": 1706000000000,
	"method": "camera_aim",
	"data": {
		"payload_index": "39-0-7",
		"camera_type": "wide",
		"locked": true,
		"x": 0.5,
		"y": 0.5
	}
}`)

// CameraFocalLengthSet is a mock camera_focal_length_set service request.
var CameraFocalLengthSet = json.RawMessage(`{
	"tid": "cam-tid-006",
	"bid": "cam-bid-006",
	"timestamp": 1706000000000,
	"method": "camera_focal_length_set",
	"data": {
		"payload_index": "39-0-7",
		"camera_type": "zoom",
		"zoom_factor": 5.0
	}
}`)

// GimbalReset is a mock gimbal_reset service request.
var GimbalReset = json.RawMessage(`{
	"tid": "cam-tid-007",
	"bid": "cam-bid-007",
	"timestamp": 1706000000000,
	"method": "gimbal_reset",
	"data": {
		"payload_index": "39-0-7",
		"reset_mode": 0
	}
}`)

// IRMeteringPoint is a mock ir_metering_point service request.
var IRMeteringPoint = json.RawMessage(`{
	"tid": "cam-tid-008",
	"bid": "cam-bid-008",
	"timestamp": 1706000000000,
	"method": "ir_metering_point",
	"data": {
		"payload_index": "39-0-7",
		"x": 0.5,
		"y": 0.5
	}
}`)

// IRMeteringArea is a mock ir_metering_area service request.
var IRMeteringArea = json.RawMessage(`{
	"tid": "cam-tid-009",
	"bid": "cam-bid-009",
	"timestamp": 1706000000000,
	"method": "ir_metering_area",
	"data": {
		"payload_index": "39-0-7",
		"x": 0.3,
		"y": 0.3,
		"width": 0.4,
		"height": 0.4
	}
}`)

// CameraServiceReplySuccess is a mock successful camera service reply.
var CameraServiceReplySuccess = json.RawMessage(`{
	"tid": "cam-tid-001",
	"bid": "cam-bid-001",
	"timestamp": 1706000001000,
	"method": "camera_mode_switch",
	"data": {
		"result": 0,
		"output": {
			"status": "success"
		}
	}
}`)
