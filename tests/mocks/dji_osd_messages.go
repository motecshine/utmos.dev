// Package mocks provides mock data for testing DJI adapter.
package mocks

import "encoding/json"

// AircraftOSDFull is a full aircraft OSD message sample.
var AircraftOSDFull = json.RawMessage(`{
	"mode_code": 0,
	"mode_code_reason": 0,
	"gear": 0,
	"firmware_version": "01.00.0000",
	"longitude": 116.397128,
	"latitude": 39.916527,
	"height": 100.5,
	"elevation": 50.0,
	"horizontal_speed": 5.5,
	"vertical_speed": 1.2,
	"attitude_pitch": 10.5,
	"attitude_roll": 5.0,
	"attitude_head": 180.0,
	"home_longitude": 116.397000,
	"home_latitude": 39.916000,
	"home_distance": 100.0,
	"wind_speed": 3.5,
	"wind_direction": 1,
	"battery": {
		"capacity_percent": 85,
		"remain_flight_time": 1200,
		"return_home_power": 30,
		"landing_power": 10,
		"batteries": [
			{
				"capacity_percent": 85,
				"index": 0,
				"sn": "BAT-001",
				"voltage": 15200,
				"temperature": 25.5
			}
		]
	},
	"cameras": [
		{
			"payload_index": "39-0-7",
			"camera_mode": 0,
			"photo_state": 0,
			"recording_state": 0,
			"zoom_factor": 2.0,
			"remain_photo_num": 1000,
			"remain_record_duration": 3600
		}
	],
	"payloads": [
		{
			"payload_index": "39-0-7",
			"sn": "PAYLOAD-001",
			"firmware_version": "01.00.0000"
		}
	],
	"position_state": {
		"is_calibration": 1,
		"is_fixed": 2,
		"quality": 5,
		"gps_number": 20,
		"rtk_number": 15
	},
	"obstacle_avoidance": {
		"horizon": 1,
		"upside": 1,
		"downside": 1
	},
	"total_flight_time": 36000.0,
	"total_flight_distance": 50000.0,
	"total_flight_sorties": 100.0
}`)

// AircraftOSDPartial is a partial aircraft OSD message sample.
var AircraftOSDPartial = json.RawMessage(`{
	"mode_code": 1,
	"longitude": 116.398000,
	"latitude": 39.917000,
	"height": 105.0,
	"battery": {
		"capacity_percent": 80
	}
}`)

// DockOSDFull is a full dock OSD message sample.
var DockOSDFull = json.RawMessage(`{
	"mode_code": 0,
	"cover_state": 0,
	"putter_state": 0,
	"drone_in_dock": 1,
	"emergency_stop_state": 0,
	"longitude": 116.397128,
	"latitude": 39.916527,
	"height": 50.0,
	"firmware_version": "01.00.0000",
	"environment_temperature": 25.5,
	"temperature": 30.0,
	"humidity": 60.0,
	"rainfall": 0,
	"wind_speed": 3.5,
	"network_state": {
		"type": 1,
		"quality": 4,
		"rate": 1024.5
	},
	"drone_charge_state": {
		"capacity_percent": 100,
		"state": 0
	},
	"storage": {
		"total": 1048576,
		"used": 524288
	},
	"position_state": {
		"is_calibration": 1,
		"is_fixed": 2,
		"quality": 10,
		"gps_number": 20,
		"rtk_number": 15
	},
	"sub_device": {
		"device_sn": "AIRCRAFT-SN-001",
		"product_type": "0-67-0",
		"device_online_status": 1,
		"device_paired": 1
	},
	"backup_battery": {
		"switch": 1,
		"voltage": 12000,
		"temperature": 25.0
	},
	"job_number": 100,
	"acc_time": 360000,
	"maintain_status": {
		"maintain_status_array": [
			{
				"state": 0,
				"last_maintain_type": 1,
				"last_maintain_time": 1700000000
			}
		]
	}
}`)

// DockOSDPartial is a partial dock OSD message sample.
var DockOSDPartial = json.RawMessage(`{
	"mode_code": 1,
	"cover_state": 1,
	"drone_in_dock": 0
}`)

// RCOSDFull is a full RC OSD message sample.
var RCOSDFull = json.RawMessage(`{
	"capacity_percent": 80,
	"longitude": 116.397128,
	"latitude": 39.916527,
	"height": 50.0,
	"country": "CN",
	"wireless_link": {
		"dongle_number": 2,
		"4g_link_state": 1,
		"sdr_link_state": 1,
		"link_workmode": 1,
		"sdr_quality": 4,
		"4g_quality": 3,
		"4g_uav_quality": 4,
		"4g_gnd_quality": 3,
		"sdr_freq_band": 2.4,
		"4g_freq_band": 1.8
	}
}`)

// RCOSDPartial is a partial RC OSD message sample.
var RCOSDPartial = json.RawMessage(`{
	"capacity_percent": 75,
	"wireless_link": {
		"sdr_quality": 3
	}
}`)

// OSDMessageSamples contains all OSD message samples for testing.
var OSDMessageSamples = map[string]json.RawMessage{
	"aircraft_full":    AircraftOSDFull,
	"aircraft_partial": AircraftOSDPartial,
	"dock_full":        DockOSDFull,
	"dock_partial":     DockOSDPartial,
	"rc_full":          RCOSDFull,
	"rc_partial":       RCOSDPartial,
}
