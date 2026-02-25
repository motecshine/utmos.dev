package mocks

import "encoding/json"

// EventHMS is a mock HMS event message.
var EventHMS = json.RawMessage(`{
	"tid": "evt-tid-001",
	"bid": "evt-bid-001",
	"timestamp": 1706000000000,
	"method": "hms",
	"data": {
		"list": [
			{
				"code": "0x16100001",
				"level": 0,
				"module": 3,
				"in_the_sky": 0,
				"args": {
					"component_index": 0,
					"sensor_index": 0
				}
			}
		]
	}
}`)

// EventHMSMultiple is a mock HMS event with multiple items.
var EventHMSMultiple = json.RawMessage(`{
	"tid": "evt-tid-002",
	"bid": "evt-bid-002",
	"timestamp": 1706000000000,
	"method": "hms",
	"data": {
		"list": [
			{
				"code": "0x16100001",
				"level": 0,
				"module": 3,
				"in_the_sky": 0
			},
			{
				"code": "0x16100002",
				"level": 1,
				"module": 3,
				"in_the_sky": 1
			},
			{
				"code": "0x16100003",
				"level": 2,
				"module": 5,
				"in_the_sky": 1
			}
		]
	}
}`)

// EventFileUploadCallback is a mock file upload callback event.
var EventFileUploadCallback = json.RawMessage(`{
	"tid": "evt-tid-003",
	"bid": "evt-bid-003",
	"timestamp": 1706000000000,
	"method": "file_upload_callback",
	"data": {
		"file": {
			"path": "/media/DJI_001.jpg",
			"name": "DJI_001.jpg",
			"size": 1024000,
			"fingerprint": "abc123def456"
		}
	}
}`)

// EventFlighttaskProgress is a mock flight task progress event.
var EventFlighttaskProgress = json.RawMessage(`{
	"tid": "evt-tid-004",
	"bid": "evt-bid-004",
	"timestamp": 1706000000000,
	"method": "flighttask_progress",
	"data": {
		"flight_id": "flight-001",
		"status": "executing",
		"progress": 50
	}
}`)

// EventFlighttaskProgressComplete is a mock completed flight task progress event.
var EventFlighttaskProgressComplete = json.RawMessage(`{
	"tid": "evt-tid-005",
	"bid": "evt-bid-005",
	"timestamp": 1706000000000,
	"method": "flighttask_progress",
	"data": {
		"flight_id": "flight-001",
		"status": "completed",
		"progress": 100,
		"result": 0
	}
}`)

// EventFlighttaskReady is a mock flight task ready event (requires reply).
var EventFlighttaskReady = json.RawMessage(`{
	"tid": "evt-tid-006",
	"bid": "evt-bid-006",
	"timestamp": 1706000000000,
	"method": "flighttask_ready",
	"need_reply": 1,
	"data": {
		"flight_id": "flight-001"
	}
}`)

// EventReturnHomeInfo is a mock return home info event.
var EventReturnHomeInfo = json.RawMessage(`{
	"tid": "evt-tid-007",
	"bid": "evt-bid-007",
	"timestamp": 1706000000000,
	"method": "return_home_info",
	"data": {
		"flight_id": "flight-001",
		"last_point_type": 1,
		"planned_path_points": [
			{"latitude": 22.5431, "longitude": 113.9472, "height": 100},
			{"latitude": 22.5432, "longitude": 113.9473, "height": 80},
			{"latitude": 22.5433, "longitude": 113.9474, "height": 50}
		]
	}
}`)

// EventDeviceExitHomingNotify is a mock device exit homing notify event.
var EventDeviceExitHomingNotify = json.RawMessage(`{
	"tid": "evt-tid-008",
	"bid": "evt-bid-008",
	"timestamp": 1706000000000,
	"method": "device_exit_homing_notify",
	"data": {}
}`)

// EventDeviceTempNtfyNeedClear is a mock device temp notify need clear event.
var EventDeviceTempNtfyNeedClear = json.RawMessage(`{
	"tid": "evt-tid-009",
	"bid": "evt-bid-009",
	"timestamp": 1706000000000,
	"method": "device_temp_ntfy_need_clear",
	"data": {}
}`)

// EventDRCStatusNotify is a mock DRC status notify event.
var EventDRCStatusNotify = json.RawMessage(`{
	"tid": "evt-tid-010",
	"bid": "evt-bid-010",
	"timestamp": 1706000000000,
	"method": "drc_status_notify",
	"data": {
		"status": 1
	}
}`)

// EventJoystickInvalidNotify is a mock joystick invalid notify event.
var EventJoystickInvalidNotify = json.RawMessage(`{
	"tid": "evt-tid-011",
	"bid": "evt-bid-011",
	"timestamp": 1706000000000,
	"method": "joystick_invalid_notify",
	"data": {}
}`)

// EventOTAProgress is a mock OTA progress event.
var EventOTAProgress = json.RawMessage(`{
	"tid": "evt-tid-012",
	"bid": "evt-bid-012",
	"timestamp": 1706000000000,
	"method": "ota_progress",
	"data": {
		"progress": 50,
		"status": "downloading"
	}
}`)

// EventReplySuccess is a mock successful event reply.
var EventReplySuccess = json.RawMessage(`{
	"tid": "evt-tid-006",
	"bid": "evt-bid-006",
	"timestamp": 1706000001000,
	"method": "flighttask_ready",
	"data": {
		"result": 0
	}
}`)

// EventReplyError is a mock error event reply.
var EventReplyError = json.RawMessage(`{
	"tid": "evt-tid-006",
	"bid": "evt-bid-006",
	"timestamp": 1706000001000,
	"method": "flighttask_ready",
	"data": {
		"result": 314001,
		"output": {
			"message": "flight task not found"
		}
	}
}`)

// EventHighestPriorityUpload is a mock highest priority upload event.
var EventHighestPriorityUpload = json.RawMessage(`{
	"tid": "evt-tid-013",
	"bid": "evt-bid-013",
	"timestamp": 1706000000000,
	"method": "highest_priority_upload_flighttask_media",
	"need_reply": 1,
	"data": {
		"flight_id": "flight-001",
		"file_list": [
			{
				"path": "/media/DJI_001.jpg",
				"name": "DJI_001.jpg"
			}
		]
	}
}`)

// EventFileUploadProgress is a mock file upload progress event.
var EventFileUploadProgress = json.RawMessage(`{
	"tid": "evt-tid-014",
	"bid": "evt-bid-014",
	"timestamp": 1706000000000,
	"method": "file_upload_progress",
	"data": {
		"file_path": "/media/DJI_001.jpg",
		"progress": 75,
		"upload_rate": 1024
	}
}`)
