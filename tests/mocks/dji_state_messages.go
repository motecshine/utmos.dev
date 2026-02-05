// Package mocks provides mock data for testing DJI adapter.
package mocks

import "encoding/json"

// StatePropertyChange is a state message with property changes.
var StatePropertyChange = json.RawMessage(`{
	"mode_code": 1,
	"firmware_version": "01.00.0001",
	"gear": 2
}`)

// StateMultipleProperties is a state message with multiple property changes.
var StateMultipleProperties = json.RawMessage(`{
	"mode_code": 2,
	"gear": 1,
	"height_limit": 120,
	"distance_limit_status": {
		"state": 1,
		"distance_limit": 5000
	},
	"rth_altitude": 100.0,
	"rc_lost_action": 2
}`)

// StateDockProperties is a state message for dock property changes.
var StateDockProperties = json.RawMessage(`{
	"mode_code": 1,
	"cover_state": 1,
	"putter_state": 1,
	"drone_in_dock": 0,
	"alarm_state": 1
}`)

// StatusOnlineDock is a status message for dock coming online.
var StatusOnlineDock = json.RawMessage(`{
	"online": true,
	"gateway_sn": "DOCK-SN-001",
	"gateway_type": "dock",
	"sub_devices": [
		{
			"device_sn": "AIRCRAFT-SN-001",
			"product_type": "0-67-0",
			"online": true
		}
	]
}`)

// StatusOnlineDock2 is a status message for dock2 coming online.
var StatusOnlineDock2 = json.RawMessage(`{
	"online": true,
	"gateway_sn": "DOCK2-SN-001",
	"gateway_type": "dock2",
	"sub_devices": [
		{
			"device_sn": "AIRCRAFT-SN-002",
			"product_type": "0-89-0",
			"online": true
		},
		{
			"device_sn": "PAYLOAD-SN-001",
			"product_type": "1-39-0",
			"online": true
		}
	]
}`)

// StatusOnlineRC is a status message for RC coming online.
var StatusOnlineRC = json.RawMessage(`{
	"online": true,
	"gateway_sn": "RC-SN-001",
	"gateway_type": "rc",
	"sub_devices": [
		{
			"device_sn": "AIRCRAFT-SN-003",
			"product_type": "0-67-0",
			"online": true
		}
	]
}`)

// StatusOffline is a status message for device going offline.
var StatusOffline = json.RawMessage(`{
	"online": false,
	"gateway_sn": "DOCK-SN-001"
}`)

// StatusOnlineNumeric is a status message with numeric online field.
var StatusOnlineNumeric = json.RawMessage(`{
	"online": 1,
	"gateway_sn": "DOCK-SN-001"
}`)

// StatusOfflineNumeric is a status message with numeric offline field.
var StatusOfflineNumeric = json.RawMessage(`{
	"online": 0,
	"gateway_sn": "DOCK-SN-001"
}`)

// StateMessageSamples contains all state message samples for testing.
var StateMessageSamples = map[string]json.RawMessage{
	"property_change":     StatePropertyChange,
	"multiple_properties": StateMultipleProperties,
	"dock_properties":     StateDockProperties,
}

// StatusMessageSamples contains all status message samples for testing.
var StatusMessageSamples = map[string]json.RawMessage{
	"online_dock":      StatusOnlineDock,
	"online_dock2":     StatusOnlineDock2,
	"online_rc":        StatusOnlineRC,
	"offline":          StatusOffline,
	"online_numeric":   StatusOnlineNumeric,
	"offline_numeric":  StatusOfflineNumeric,
}
