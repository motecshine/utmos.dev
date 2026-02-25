package mocks

import "encoding/json"

// ServiceRequestCoverOpen is a mock service request for cover_open.
var ServiceRequestCoverOpen = json.RawMessage(`{
	"tid": "srv-tid-001",
	"bid": "srv-bid-001",
	"timestamp": 1706000000000,
	"method": "cover_open",
	"data": {}
}`)

// ServiceRequestCoverClose is a mock service request for cover_close.
var ServiceRequestCoverClose = json.RawMessage(`{
	"tid": "srv-tid-002",
	"bid": "srv-bid-002",
	"timestamp": 1706000000000,
	"method": "cover_close",
	"data": {}
}`)

// ServiceRequestDroneOpen is a mock service request for drone_open.
var ServiceRequestDroneOpen = json.RawMessage(`{
	"tid": "srv-tid-003",
	"bid": "srv-bid-003",
	"timestamp": 1706000000000,
	"method": "drone_open",
	"data": {}
}`)

// ServiceRequestDeviceReboot is a mock service request for device_reboot.
var ServiceRequestDeviceReboot = json.RawMessage(`{
	"tid": "srv-tid-004",
	"bid": "srv-bid-004",
	"timestamp": 1706000000000,
	"method": "device_reboot",
	"data": {}
}`)

// ServiceRequestBatteryMaintenanceSwitch is a mock service request for battery_maintenance_switch.
var ServiceRequestBatteryMaintenanceSwitch = json.RawMessage(`{
	"tid": "srv-tid-005",
	"bid": "srv-bid-005",
	"timestamp": 1706000000000,
	"method": "battery_maintenance_switch",
	"data": {
		"enable": 1
	}
}`)

// ServiceRequestAirConditionerModeSwitch is a mock service request for air_conditioner_mode_switch.
var ServiceRequestAirConditionerModeSwitch = json.RawMessage(`{
	"tid": "srv-tid-006",
	"bid": "srv-bid-006",
	"timestamp": 1706000000000,
	"method": "air_conditioner_mode_switch",
	"data": {
		"mode": 2
	}
}`)

// ServiceReplySuccess is a mock successful service reply.
var ServiceReplySuccess = json.RawMessage(`{
	"tid": "srv-tid-001",
	"bid": "srv-bid-001",
	"timestamp": 1706000001000,
	"method": "cover_open",
	"data": {
		"result": 0,
		"output": {
			"status": "success"
		}
	}
}`)

// ServiceReplyError is a mock error service reply.
var ServiceReplyError = json.RawMessage(`{
	"tid": "srv-tid-001",
	"bid": "srv-bid-001",
	"timestamp": 1706000001000,
	"method": "cover_open",
	"data": {
		"result": 314001,
		"output": {
			"message": "device offline"
		}
	}
}`)

// ServiceReplyTimeout is a mock timeout service reply.
var ServiceReplyTimeout = json.RawMessage(`{
	"tid": "srv-tid-001",
	"bid": "srv-bid-001",
	"timestamp": 1706000031000,
	"method": "cover_open",
	"data": {
		"result": 314002,
		"output": {
			"message": "operation timeout"
		}
	}
}`)

// ServiceRequestWithNeedReply is a mock service request with need_reply flag.
var ServiceRequestWithNeedReply = json.RawMessage(`{
	"tid": "srv-tid-007",
	"bid": "srv-bid-007",
	"timestamp": 1706000000000,
	"method": "flighttask_prepare",
	"need_reply": 1,
	"data": {
		"flight_id": "flight-001"
	}
}`)

// ServiceRequestFlighttaskCreate is a mock service request for flighttask_create.
var ServiceRequestFlighttaskCreate = json.RawMessage(`{
	"tid": "srv-tid-008",
	"bid": "srv-bid-008",
	"timestamp": 1706000000000,
	"method": "flighttask_create",
	"data": {
		"flight_id": "flight-001",
		"type": 0,
		"file": {
			"url": "https://example.com/wayline.kmz",
			"fingerprint": "abc123"
		}
	}
}`)

// ServiceRequestFlighttaskExecute is a mock service request for flighttask_execute.
var ServiceRequestFlighttaskExecute = json.RawMessage(`{
	"tid": "srv-tid-009",
	"bid": "srv-bid-009",
	"timestamp": 1706000000000,
	"method": "flighttask_execute",
	"data": {
		"flight_id": "flight-001"
	}
}`)
