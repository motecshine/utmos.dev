package wpml

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConvertWaylinesToWPMLMission(t *testing.T) {
	tests := []struct {
		name        string
		waylines    *Waylines
		expectError bool
		errorMsg    string
	}{
		{
			name: "Valid basic mission",
			waylines: &Waylines{
				Name:                    "Test Mission",
				DroneModel:              DroneM3DSeries,
				PayloadModel:            PayloadMatrice3TD,
				TemplateType:            TemplateTypeWaypoint,
				GlobalHeight:            50.0,
				GlobalSpeed:             5.0,
				ClimbMode:               "vertical",
				SafeHeight:              30.0,
				GlobalRTHHeight:         100.0,
				AircraftYawMode:         "followWayline",
				GimbalPitchMode:         "usePointSetting",
				GlobalTransitionalSpeed: 5.0,
				Waypoints: []WaylinesWaypoint{
					{
						Latitude:    39.9042,
						Longitude:   116.4074,
						Height:      50.0,
						Speed:       5.0,
						TriggerType: "reachPoint",
					},
					{
						Latitude:    39.9052,
						Longitude:   116.4084,
						Height:      60.0,
						Speed:       5.0,
						TriggerType: "reachPoint",
					},
				},
			},
			expectError: false,
		},
		{
			name: "Invalid - missing name",
			waylines: &Waylines{
				DroneModel:   DroneM3DSeries,
				PayloadModel: PayloadMatrice3TD,
				TemplateType: TemplateTypeWaypoint,
				Waypoints: []WaylinesWaypoint{
					{
						Latitude:  39.9042,
						Longitude: 116.4074,
						Height:    50.0,
					},
				},
			},
			expectError: true,
		},
		{
			name: "Invalid - missing waypoints",
			waylines: &Waylines{
				Name:         "Test Mission",
				DroneModel:   DroneM3DSeries,
				PayloadModel: PayloadMatrice3TD,
				TemplateType: TemplateTypeWaypoint,
			},
			expectError: true,
		},
		{
			name: "Valid mission with actions",
			waylines: &Waylines{
				Name:                    "Action Mission",
				DroneModel:              DroneM3DSeries,
				PayloadModel:            PayloadMatrice3TD,
				TemplateType:            TemplateTypeWaypoint,
				GlobalHeight:            80.0,
				GlobalSpeed:             8.0,
				ClimbMode:               "vertical",
				SafeHeight:              30.0,
				GlobalRTHHeight:         100.0,
				AircraftYawMode:         "followWayline",
				GimbalPitchMode:         "usePointSetting",
				GlobalTransitionalSpeed: 8.0,
				Waypoints: []WaylinesWaypoint{
					{
						Latitude:    39.9042,
						Longitude:   116.4074,
						Height:      80.0,
						Speed:       8.0,
						TriggerType: "reachPoint",
						Actions: []ActionRequest{
							{
								Type:   ActionTypeTakePhoto,
								Action: &TakePhotoAction{PayloadPositionIndex: PayloadPosition0},
							},
						},
					},
				},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mission, err := ConvertWaylinesToWPMLMission(tt.waylines)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, mission)
			} else {
				require.NoError(t, err)
				require.NotNil(t, mission)

				// Verify basic mission properties
				assert.Equal(t, DefaultAuthor, mission.Template.Document.Author)
				assert.NotEmpty(t, mission.Template.Document.CreateTime)
				assert.NotNil(t, mission.Template.Document.MissionConfig)
				assert.NotEmpty(t, mission.Template.Document.Folders)
				assert.NotEmpty(t, mission.Waylines.Document.Folders)

				// Verify waypoints count
				waylineFolder := mission.Waylines.Document.Folders[0]
				expectedWaypoints := len(tt.waylines.Waypoints)
				assert.Equal(t, expectedWaypoints, len(waylineFolder.Placemarks))
			}
		})
	}
}

func TestConvertWaylinesToWPMLMissionWithJSON(t *testing.T) {
	jsonData := `{
		"name": "JSON Test Mission",
		"drone_model": 91,
		"payload_model": 81,
		"template_type": "waypoint",
		"global_height": 100,
		"global_speed": 10,
		"climb_mode": "vertical",
		"safe_height": 30,
		"global_rth_height": 120,
		"aircraft_yaw_mode": "followWayline",
		"gimbal_pitch_mode": "usePointSetting",
		"global_transitional_speed": 10,
		"waypoints": [
			{
				"latitude": 39.9042,
				"longitude": 116.4074,
				"height": 100,
				"speed": 10,
				"trigger_type": "reachPoint",
				"actions": [
					{
						"type": "takePhoto",
						"action": {
							"payload_position_index": 0
						}
					}
				]
			},
			{
				"latitude": 39.9052,
				"longitude": 116.4084,
				"height": 120,
				"speed": 10,
				"trigger_type": "reachPoint"
			}
		]
	}`

	var waylines Waylines
	err := json.Unmarshal([]byte(jsonData), &waylines)
	require.NoError(t, err)

	mission, err := ConvertWaylinesToWPMLMission(&waylines)
	require.NoError(t, err)
	require.NotNil(t, mission)

	assert.Equal(t, "JSON Test Mission", waylines.Name)
	assert.Equal(t, DroneM3DSeries, waylines.DroneModel)
	assert.Equal(t, PayloadMatrice3TD, waylines.PayloadModel)
	assert.Equal(t, 2, len(waylines.Waypoints))
	assert.Equal(t, 1, len(waylines.Waypoints[0].Actions))
}

func TestWaylinesValidation(t *testing.T) {
	tests := []struct {
		name        string
		waylines    *Waylines
		expectError bool
	}{
		{
			name: "Valid waylines",
			waylines: &Waylines{
				Name:                    "Valid Mission",
				DroneModel:              DroneM3DSeries,
				PayloadModel:            PayloadMatrice3TD,
				TemplateType:            TemplateTypeWaypoint,
				GlobalHeight:            50.0,
				GlobalSpeed:             5.0,
				ClimbMode:               "vertical",
				SafeHeight:              30.0,
				GlobalRTHHeight:         100.0,
				AircraftYawMode:         "followWayline",
				GimbalPitchMode:         "usePointSetting",
				GlobalTransitionalSpeed: 5.0,
				Waypoints: []WaylinesWaypoint{
					{
						Latitude:    39.9042,
						Longitude:   116.4074,
						Height:      50.0,
						Speed:       5.0,
						TriggerType: "reachPoint",
					},
				},
			},
			expectError: false,
		},
		{
			name: "Invalid - empty name",
			waylines: &Waylines{
				Name:         "",
				DroneModel:   DroneM3DSeries,
				PayloadModel: PayloadMatrice3TD,
				TemplateType: TemplateTypeWaypoint,
				Waypoints: []WaylinesWaypoint{
					{
						Latitude:  39.9042,
						Longitude: 116.4074,
						Height:    50.0,
					},
				},
			},
			expectError: true,
		},
		{
			name: "Invalid - invalid latitude",
			waylines: &Waylines{
				Name:         "Test Mission",
				DroneModel:   DroneM3DSeries,
				PayloadModel: PayloadMatrice3TD,
				TemplateType: TemplateTypeWaypoint,
				Waypoints: []WaylinesWaypoint{
					{
						Latitude:  -100.0, // Invalid latitude
						Longitude: 116.4074,
						Height:    50.0,
					},
				},
			},
			expectError: true,
		},
		{
			name: "Invalid - height too low",
			waylines: &Waylines{
				Name:         "Test Mission",
				DroneModel:   DroneM3DSeries,
				PayloadModel: PayloadMatrice3TD,
				TemplateType: TemplateTypeWaypoint,
				Waypoints: []WaylinesWaypoint{
					{
						Latitude:  39.9042,
						Longitude: 116.4074,
						Height:    1.0, // Too low
					},
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.waylines.Validate()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestWaylinesApplyDefaults(t *testing.T) {
	waylines := &Waylines{
		Name:         "Test Mission",
		DroneModel:   DroneM3DSeries,
		PayloadModel: PayloadMatrice3TD,
		TemplateType: TemplateTypeWaypoint,
		Waypoints: []WaylinesWaypoint{
			{
				Latitude:  39.9042,
				Longitude: 116.4074,
				Height:    50.0,
			},
		},
	}

	// Before applying defaults
	assert.Empty(t, waylines.HeightType)

	waylines.ApplyDefaults()

	// After applying defaults
	assert.Equal(t, HeightModeRelativeToStartPoint, waylines.HeightType)
}
