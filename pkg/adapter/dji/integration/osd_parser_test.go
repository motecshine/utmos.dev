package integration

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewOSDParser(t *testing.T) {
	parser := NewOSDParser()
	assert.NotNil(t, parser)
}

func TestOSDParser_ParseAircraftOSD(t *testing.T) {
	parser := NewOSDParser()

	tests := []struct {
		name    string
		data    string
		wantErr bool
		check   func(t *testing.T, osd interface{})
	}{
		{
			name: "full aircraft OSD",
			data: `{
				"mode_code": 0,
				"longitude": 116.397128,
				"latitude": 39.916527,
				"height": 100.5,
				"elevation": 50.0,
				"horizontal_speed": 5.5,
				"vertical_speed": 1.2,
				"attitude_pitch": 10.5,
				"attitude_roll": 5.0,
				"attitude_head": 180.0,
				"battery": {
					"capacity_percent": 85,
					"remain_flight_time": 1200
				}
			}`,
			wantErr: false,
			check: func(_ *testing.T, osd interface{}) {
				a := osd.(*struct {
					ModeCode        *int
					Longitude       *float64
					Latitude        *float64
					Height          *float64
					Elevation       *float64
					HorizontalSpeed *float64
					VerticalSpeed   *float64
				})
				_ = a // Type assertion for documentation
			},
		},
		{
			name: "partial aircraft OSD",
			data: `{
				"mode_code": 1,
				"battery": {
					"capacity_percent": 75
				}
			}`,
			wantErr: false,
		},
		{
			name:    "empty data",
			data:    "",
			wantErr: true,
		},
		{
			name:    "invalid JSON",
			data:    "{invalid}",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var data json.RawMessage
			if tt.data != "" {
				data = json.RawMessage(tt.data)
			}

			osd, err := parser.ParseAircraftOSD(data)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, osd)
		})
	}
}

func TestOSDParser_ParseAircraftOSD_FullData(t *testing.T) {
	parser := NewOSDParser()

	data := json.RawMessage(`{
		"mode_code": 0,
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
		"battery": {
			"capacity_percent": 85,
			"remain_flight_time": 1200,
			"return_home_power": 30,
			"landing_power": 10
		},
		"cameras": [
			{
				"payload_index": "39-0-7",
				"camera_mode": 0,
				"photo_state": 0,
				"recording_state": 0,
				"zoom_factor": 2.0
			}
		],
		"payloads": [
			{
				"payload_index": "39-0-7",
				"sn": "PAYLOAD-001"
			}
		]
	}`)

	osd, err := parser.ParseAircraftOSD(data)
	require.NoError(t, err)
	require.NotNil(t, osd)

	// Verify parsed values
	assert.NotNil(t, osd.ModeCode)
	assert.Equal(t, 0, *osd.ModeCode)

	assert.NotNil(t, osd.Longitude)
	assert.InDelta(t, 116.397128, *osd.Longitude, 0.0001)

	assert.NotNil(t, osd.Latitude)
	assert.InDelta(t, 39.916527, *osd.Latitude, 0.0001)

	assert.NotNil(t, osd.Height)
	assert.InDelta(t, 100.5, *osd.Height, 0.01)

	assert.NotNil(t, osd.Battery)
	assert.NotNil(t, osd.Battery.CapacityPercent)
	assert.Equal(t, 85, *osd.Battery.CapacityPercent)

	assert.Len(t, osd.Cameras, 1)
	assert.Len(t, osd.Payloads, 1)
}

func TestOSDParser_ParseDockOSD(t *testing.T) {
	parser := NewOSDParser()

	tests := []struct {
		name    string
		data    string
		wantErr bool
	}{
		{
			name: "full dock OSD",
			data: `{
				"mode_code": 0,
				"cover_state": 0,
				"putter_state": 0,
				"drone_in_dock": 1,
				"longitude": 116.397128,
				"latitude": 39.916527,
				"height": 50.0,
				"network_state": {
					"type": 1,
					"quality": 4
				},
				"drone_charge_state": {
					"capacity_percent": 100,
					"state": 0
				}
			}`,
			wantErr: false,
		},
		{
			name: "partial dock OSD",
			data: `{
				"mode_code": 1,
				"cover_state": 1
			}`,
			wantErr: false,
		},
		{
			name:    "empty data",
			data:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var data json.RawMessage
			if tt.data != "" {
				data = json.RawMessage(tt.data)
			}

			osd, err := parser.ParseDockOSD(data)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, osd)
		})
	}
}

func TestOSDParser_ParseDockOSD_FullData(t *testing.T) {
	parser := NewOSDParser()

	data := json.RawMessage(`{
		"mode_code": 0,
		"cover_state": 0,
		"putter_state": 0,
		"drone_in_dock": 1,
		"longitude": 116.397128,
		"latitude": 39.916527,
		"height": 50.0,
		"environment_temperature": 25.5,
		"humidity": 60.0,
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
		}
	}`)

	osd, err := parser.ParseDockOSD(data)
	require.NoError(t, err)
	require.NotNil(t, osd)

	assert.NotNil(t, osd.ModeCode)
	assert.Equal(t, 0, *osd.ModeCode)

	assert.NotNil(t, osd.CoverState)
	assert.Equal(t, 0, *osd.CoverState)

	assert.NotNil(t, osd.DroneInDock)
	assert.Equal(t, 1, *osd.DroneInDock)

	assert.NotNil(t, osd.NetworkState)
	assert.Equal(t, 1, *osd.NetworkState.Type)

	assert.NotNil(t, osd.DroneChargeState)
	assert.Equal(t, 100, *osd.DroneChargeState.CapacityPercent)
}

func TestOSDParser_ParseRCOSD(t *testing.T) {
	parser := NewOSDParser()

	tests := []struct {
		name    string
		data    string
		wantErr bool
	}{
		{
			name: "full RC OSD",
			data: `{
				"capacity_percent": 80,
				"longitude": 116.397128,
				"latitude": 39.916527,
				"height": 50.0,
				"wireless_link": {
					"dongle_number": 2,
					"sdr_link_state": 1,
					"sdr_quality": 4
				}
			}`,
			wantErr: false,
		},
		{
			name: "partial RC OSD",
			data: `{
				"capacity_percent": 75
			}`,
			wantErr: false,
		},
		{
			name:    "empty data",
			data:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var data json.RawMessage
			if tt.data != "" {
				data = json.RawMessage(tt.data)
			}

			osd, err := parser.ParseRCOSD(data)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, osd)
		})
	}
}

func TestOSDParser_ParseRCOSD_FullData(t *testing.T) {
	parser := NewOSDParser()

	data := json.RawMessage(`{
		"capacity_percent": 80,
		"longitude": 116.397128,
		"latitude": 39.916527,
		"height": 50.0,
		"country": "CN",
		"wireless_link": {
			"dongle_number": 2,
			"sdr_link_state": 1,
			"sdr_quality": 4,
			"4g_link_state": 1,
			"4g_quality": 3
		}
	}`)

	osd, err := parser.ParseRCOSD(data)
	require.NoError(t, err)
	require.NotNil(t, osd)

	assert.NotNil(t, osd.CapacityPercent)
	assert.Equal(t, 80, *osd.CapacityPercent)

	assert.NotNil(t, osd.Longitude)
	assert.InDelta(t, 116.397128, *osd.Longitude, 0.0001)

	assert.NotNil(t, osd.WirelessLink)
	assert.Equal(t, 2, *osd.WirelessLink.DongleNumber)
}

func TestOSDParser_DetectOSDType(t *testing.T) {
	parser := NewOSDParser()

	tests := []struct {
		name     string
		data     string
		expected OSDType
	}{
		{
			name:     "dock OSD with cover_state",
			data:     `{"cover_state": 0, "mode_code": 0}`,
			expected: OSDTypeDock,
		},
		{
			name:     "dock OSD with drone_in_dock",
			data:     `{"drone_in_dock": 1, "mode_code": 0}`,
			expected: OSDTypeDock,
		},
		{
			name:     "dock OSD with putter_state",
			data:     `{"putter_state": 0, "mode_code": 0}`,
			expected: OSDTypeDock,
		},
		{
			name:     "aircraft OSD",
			data:     `{"mode_code": 0, "longitude": 116.0, "payloads": []}`,
			expected: OSDTypeAircraft,
		},
		{
			name:     "RC OSD",
			data:     `{"capacity_percent": 80, "wireless_link": {}}`,
			expected: OSDTypeRC,
		},
		{
			name:     "default to aircraft",
			data:     `{"mode_code": 0}`,
			expected: OSDTypeAircraft,
		},
		{
			name:     "invalid JSON defaults to aircraft",
			data:     `{invalid}`,
			expected: OSDTypeAircraft,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := json.RawMessage(tt.data)
			osdType := parser.DetectOSDType(data)
			assert.Equal(t, tt.expected, osdType)
		})
	}
}

func TestOSDParser_ParseOSD(t *testing.T) {
	parser := NewOSDParser()

	tests := []struct {
		name         string
		data         string
		expectedType OSDType
		wantErr      bool
	}{
		{
			name:         "auto-detect dock",
			data:         `{"cover_state": 0, "mode_code": 0, "drone_in_dock": 1}`,
			expectedType: OSDTypeDock,
			wantErr:      false,
		},
		{
			name:         "auto-detect aircraft",
			data:         `{"mode_code": 0, "longitude": 116.0, "latitude": 39.0, "payloads": []}`,
			expectedType: OSDTypeAircraft,
			wantErr:      false,
		},
		{
			name:         "auto-detect RC",
			data:         `{"capacity_percent": 80, "wireless_link": {"sdr_quality": 4}}`,
			expectedType: OSDTypeRC,
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := json.RawMessage(tt.data)
			result, err := parser.ParseOSD(data)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expectedType, result.Type)

			switch tt.expectedType {
			case OSDTypeDock:
				assert.NotNil(t, result.Dock)
				assert.Nil(t, result.Aircraft)
				assert.Nil(t, result.RC)
			case OSDTypeAircraft:
				assert.NotNil(t, result.Aircraft)
				assert.Nil(t, result.Dock)
				assert.Nil(t, result.RC)
			case OSDTypeRC:
				assert.NotNil(t, result.RC)
				assert.Nil(t, result.Aircraft)
				assert.Nil(t, result.Dock)
			}
		})
	}
}
