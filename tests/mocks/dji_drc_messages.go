package mocks

import "encoding/json"

// DRCModeEnter is a mock drc_mode_enter service request.
var DRCModeEnter = json.RawMessage(`{
	"tid": "drc-tid-001",
	"bid": "drc-bid-001",
	"timestamp": 1706000000000,
	"method": "drc_mode_enter",
	"data": {
		"mqtt_broker": "mqtt://broker.example.com:1883",
		"client_id": "drc-client-001"
	}
}`)

// DRCModeExit is a mock drc_mode_exit service request.
var DRCModeExit = json.RawMessage(`{
	"tid": "drc-tid-002",
	"bid": "drc-bid-002",
	"timestamp": 1706000000000,
	"method": "drc_mode_exit",
	"data": {}
}`)

// DRCDroneControl is a mock drone_control command.
var DRCDroneControl = json.RawMessage(`{
	"tid": "drc-tid-003",
	"bid": "drc-bid-003",
	"timestamp": 1706000000000,
	"method": "drone_control",
	"data": {
		"x": 0.5,
		"y": 0.3,
		"h": 0.0,
		"w": 0.1,
		"seq": 100
	}
}`)

// DRCDroneEmergencyStop is a mock drone_emergency_stop command.
var DRCDroneEmergencyStop = json.RawMessage(`{
	"tid": "drc-tid-004",
	"bid": "drc-bid-004",
	"timestamp": 1706000000000,
	"method": "drone_emergency_stop",
	"data": {}
}`)

// DRCHeart is a mock heart (heartbeat) command.
var DRCHeart = json.RawMessage(`{
	"tid": "drc-tid-005",
	"bid": "drc-bid-005",
	"timestamp": 1706000000000,
	"method": "heart",
	"data": {
		"seq": 1,
		"timestamp": 1706000000000
	}
}`)

// DRCJoystickInvalidNotify is a mock joystick_invalid_notify event.
var DRCJoystickInvalidNotify = json.RawMessage(`{
	"tid": "drc-evt-001",
	"bid": "drc-evt-001",
	"timestamp": 1706000000000,
	"method": "joystick_invalid_notify",
	"data": {
		"reason": "timeout"
	}
}`)

// DRCStatusNotify is a mock drc_status_notify event.
var DRCStatusNotify = json.RawMessage(`{
	"tid": "drc-evt-002",
	"bid": "drc-evt-002",
	"timestamp": 1706000000000,
	"method": "drc_status_notify",
	"data": {
		"status": 1
	}
}`)

// DRCServiceReplySuccess is a mock successful DRC service reply.
var DRCServiceReplySuccess = json.RawMessage(`{
	"tid": "drc-tid-001",
	"bid": "drc-bid-001",
	"timestamp": 1706000001000,
	"method": "drc_mode_enter",
	"data": {
		"result": 0,
		"output": {
			"status": "connected"
		}
	}
}`)

// DRCServiceReplyError is a mock error DRC service reply.
var DRCServiceReplyError = json.RawMessage(`{
	"tid": "drc-tid-001",
	"bid": "drc-bid-001",
	"timestamp": 1706000001000,
	"method": "drc_mode_enter",
	"data": {
		"result": 314001,
		"output": {
			"message": "DRC mode not available"
		}
	}
}`)
