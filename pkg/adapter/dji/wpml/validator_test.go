package wpml

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewWPMLValidator(t *testing.T) {
	validator, _ := NewWPMLValidator()
	require.NotNil(t, validator)
}

func TestWPMLValidator_ValidateStruct(t *testing.T) {
	validator, _ := NewWPMLValidator()

	tests := []struct {
		name        string
		input       interface{}
		expectError bool
	}{
		{
			name: "Valid waylines",
			input: &Waylines{
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
			name: "Invalid - missing required fields",
			input: &Waylines{
				GlobalHeight: 50.0,
			},
			expectError: true,
		},
		{
			name: "Invalid - invalid height",
			input: &Waylines{
				Name:         "Test",
				DroneModel:   DroneM3DSeries,
				PayloadModel: PayloadMatrice3TD,
				TemplateType: TemplateTypeWaypoint,
				GlobalHeight: 2000.0, // Too high
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
			name: "Invalid - invalid speed",
			input: &Waylines{
				Name:         "Test",
				DroneModel:   DroneM3DSeries,
				PayloadModel: PayloadMatrice3TD,
				TemplateType: TemplateTypeWaypoint,
				GlobalSpeed:  50.0, // Too fast
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateStruct(tt.input)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDroneModelValidation(t *testing.T) {
	validator, _ := NewWPMLValidator()

	tests := []struct {
		name        string
		droneModel  DroneModel
		expectError bool
	}{
		{
			name:        "Valid M3D Series",
			droneModel:  DroneM3DSeries,
			expectError: false,
		},
		{
			name:        "Valid M300RTK",
			droneModel:  DroneM300RTK,
			expectError: false,
		},
		{
			name:        "Invalid drone model",
			droneModel:  DroneModel(9999),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			waylines := &Waylines{
				Name:                    "Test Mission",
				DroneModel:              tt.droneModel,
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
			}

			err := validator.ValidateStruct(waylines)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPayloadModelValidation(t *testing.T) {
	validator, _ := NewWPMLValidator()

	tests := []struct {
		name         string
		payloadModel PayloadModel
		expectError  bool
	}{
		{
			name:         "Valid Matrice3TD",
			payloadModel: PayloadMatrice3TD,
			expectError:  false,
		},
		{
			name:         "Valid H20T",
			payloadModel: PayloadH20T,
			expectError:  false,
		},
		{
			name:         "Invalid payload model",
			payloadModel: PayloadModel(9999),
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			waylines := &Waylines{
				Name:                    "Test Mission",
				DroneModel:              DroneM3DSeries,
				PayloadModel:            tt.payloadModel,
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
			}

			err := validator.ValidateStruct(waylines)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestWaypointValidation(t *testing.T) {
	validator, _ := NewWPMLValidator()

	tests := []struct {
		name        string
		waypoint    WaylinesWaypoint
		expectError bool
	}{
		{
			name: "Valid waypoint",
			waypoint: WaylinesWaypoint{
				Latitude:    39.9042,
				Longitude:   116.4074,
				Height:      50.0,
				Speed:       5.0,
				TriggerType: "reachPoint",
			},
			expectError: false,
		},
		{
			name: "Invalid latitude - too high",
			waypoint: WaylinesWaypoint{
				Latitude:  95.0, // > 90
				Longitude: 116.4074,
				Height:    50.0,
			},
			expectError: true,
		},
		{
			name: "Invalid latitude - too low",
			waypoint: WaylinesWaypoint{
				Latitude:  -95.0, // < -90
				Longitude: 116.4074,
				Height:    50.0,
			},
			expectError: true,
		},
		{
			name: "Invalid longitude - too high",
			waypoint: WaylinesWaypoint{
				Latitude:  39.9042,
				Longitude: 185.0, // > 180
				Height:    50.0,
			},
			expectError: true,
		},
		{
			name: "Invalid longitude - too low",
			waypoint: WaylinesWaypoint{
				Latitude:  39.9042,
				Longitude: -185.0, // < -180
				Height:    50.0,
			},
			expectError: true,
		},
		{
			name: "Invalid height - too low",
			waypoint: WaylinesWaypoint{
				Latitude:  39.9042,
				Longitude: 116.4074,
				Height:    2.0, // < 5
			},
			expectError: true,
		},
		{
			name: "Invalid height - too high",
			waypoint: WaylinesWaypoint{
				Latitude:  39.9042,
				Longitude: 116.4074,
				Height:    600.0, // > 500
			},
			expectError: true,
		},
		{
			name: "Invalid speed - too low",
			waypoint: WaylinesWaypoint{
				Latitude:  39.9042,
				Longitude: 116.4074,
				Height:    50.0,
				Speed:     0.5, // < 1
			},
			expectError: true,
		},
		{
			name: "Invalid speed - too high",
			waypoint: WaylinesWaypoint{
				Latitude:  39.9042,
				Longitude: 116.4074,
				Height:    50.0,
				Speed:     20.0, // > 15
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			waylines := &Waylines{
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
				Waypoints:               []WaylinesWaypoint{tt.waypoint},
			}

			err := validator.ValidateStruct(waylines)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
