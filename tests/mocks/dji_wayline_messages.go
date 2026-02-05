package mocks

import "encoding/json"

// WaylineFlighttaskCreate is a mock flighttask_create service request.
var WaylineFlighttaskCreate = json.RawMessage(`{
	"tid": "wl-tid-001",
	"bid": "wl-bid-001",
	"timestamp": 1706000000000,
	"method": "flighttask_create",
	"data": {
		"flight_id": "flight-001",
		"type": 0,
		"file": {
			"url": "https://example.com/wayline.kmz",
			"fingerprint": "abc123def456"
		},
		"rth_altitude": 100.0
	}
}`)

// WaylineFlighttaskPrepare is a mock flighttask_prepare service request.
var WaylineFlighttaskPrepare = json.RawMessage(`{
	"tid": "wl-tid-002",
	"bid": "wl-bid-002",
	"timestamp": 1706000000000,
	"method": "flighttask_prepare",
	"need_reply": 1,
	"data": {
		"flight_id": "flight-001"
	}
}`)

// WaylineFlighttaskExecute is a mock flighttask_execute service request.
var WaylineFlighttaskExecute = json.RawMessage(`{
	"tid": "wl-tid-003",
	"bid": "wl-bid-003",
	"timestamp": 1706000000000,
	"method": "flighttask_execute",
	"data": {
		"flight_id": "flight-001"
	}
}`)

// WaylineFlighttaskPause is a mock flighttask_pause service request.
var WaylineFlighttaskPause = json.RawMessage(`{
	"tid": "wl-tid-004",
	"bid": "wl-bid-004",
	"timestamp": 1706000000000,
	"method": "flighttask_pause",
	"data": {}
}`)

// WaylineFlighttaskRecovery is a mock flighttask_recovery service request.
var WaylineFlighttaskRecovery = json.RawMessage(`{
	"tid": "wl-tid-005",
	"bid": "wl-bid-005",
	"timestamp": 1706000000000,
	"method": "flighttask_recovery",
	"data": {}
}`)

// WaylineFlighttaskUndo is a mock flighttask_undo service request.
var WaylineFlighttaskUndo = json.RawMessage(`{
	"tid": "wl-tid-006",
	"bid": "wl-bid-006",
	"timestamp": 1706000000000,
	"method": "flighttask_undo",
	"data": {
		"flight_ids": ["flight-001", "flight-002"]
	}
}`)

// WaylineReturnHome is a mock return_home service request.
var WaylineReturnHome = json.RawMessage(`{
	"tid": "wl-tid-007",
	"bid": "wl-bid-007",
	"timestamp": 1706000000000,
	"method": "return_home",
	"data": {}
}`)

// WaylineReturnHomeCancel is a mock return_home_cancel service request.
var WaylineReturnHomeCancel = json.RawMessage(`{
	"tid": "wl-tid-008",
	"bid": "wl-bid-008",
	"timestamp": 1706000000000,
	"method": "return_home_cancel",
	"data": {}
}`)

// WaylineFlighttaskProgressEvent is a mock flighttask_progress event.
var WaylineFlighttaskProgressEvent = json.RawMessage(`{
	"tid": "wl-evt-001",
	"bid": "wl-evt-001",
	"timestamp": 1706000000000,
	"method": "flighttask_progress",
	"data": {
		"flight_id": "flight-001",
		"status": "executing",
		"progress": 50
	}
}`)

// WaylineFlighttaskReadyEvent is a mock flighttask_ready event (requires reply).
var WaylineFlighttaskReadyEvent = json.RawMessage(`{
	"tid": "wl-evt-002",
	"bid": "wl-evt-002",
	"timestamp": 1706000000000,
	"method": "flighttask_ready",
	"need_reply": 1,
	"data": {
		"flight_id": "flight-001"
	}
}`)

// WaylineReturnHomeInfoEvent is a mock return_home_info event.
var WaylineReturnHomeInfoEvent = json.RawMessage(`{
	"tid": "wl-evt-003",
	"bid": "wl-evt-003",
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

// WaylineServiceReplySuccess is a mock successful wayline service reply.
var WaylineServiceReplySuccess = json.RawMessage(`{
	"tid": "wl-tid-001",
	"bid": "wl-bid-001",
	"timestamp": 1706000001000,
	"method": "flighttask_create",
	"data": {
		"result": 0,
		"output": {
			"status": "success"
		}
	}
}`)

// WaylineServiceReplyError is a mock error wayline service reply.
var WaylineServiceReplyError = json.RawMessage(`{
	"tid": "wl-tid-001",
	"bid": "wl-bid-001",
	"timestamp": 1706000001000,
	"method": "flighttask_create",
	"data": {
		"result": 314001,
		"output": {
			"message": "wayline file not found"
		}
	}
}`)
