package wpml

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateKmz(t *testing.T) {
	// Create a test mission
	waylines := &Waylines{
		Name:                    "Test KMZ Mission",
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
	}

	mission, err := ConvertWaylinesToWPMLMission(waylines)
	require.NoError(t, err)

	// Create temp directory for test
	tmpDir, err := os.MkdirTemp("", "wpml_test_*")
	require.NoError(t, err)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	kmzPath := filepath.Join(tmpDir, "test_mission.kmz")

	// Test CreateKmz
	err = CreateKmz(mission, kmzPath)
	require.NoError(t, err)

	// Verify file exists
	_, err = os.Stat(kmzPath)
	require.NoError(t, err, "KMZ file should exist")

	// Verify file is not empty
	info, err := os.Stat(kmzPath)
	require.NoError(t, err)
	assert.Greater(t, info.Size(), int64(0), "KMZ file should not be empty")
}

func TestCreateKmzBuffer(t *testing.T) {
	waylines := &Waylines{
		Name:                    "Buffer Test Mission",
		DroneModel:              DroneM3DSeries,
		PayloadModel:            PayloadMatrice3TD,
		TemplateType:            TemplateTypeWaypoint,
		GlobalHeight:            80.0,
		GlobalSpeed:             8.0,
		ClimbMode:               "vertical",
		SafeHeight:              30.0,
		GlobalRTHHeight:         120.0,
		AircraftYawMode:         "followWayline",
		GimbalPitchMode:         "usePointSetting",
		GlobalTransitionalSpeed: 8.0,
		Waypoints: []WaylinesWaypoint{
			{
				Latitude:    40.7589,
				Longitude:   -73.9851,
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
	}

	mission, err := ConvertWaylinesToWPMLMission(waylines)
	require.NoError(t, err)

	// Test CreateKmzBuffer
	buffer, err := CreateKmzBuffer(mission)
	require.NoError(t, err)
	require.NotNil(t, buffer)

	// Buffer should not be empty
	assert.Greater(t, buffer.Len(), 0, "Buffer should not be empty")

	// Buffer should contain ZIP data
	bufferBytes := buffer.Bytes()
	assert.True(t, len(bufferBytes) > 0)

	// ZIP files start with "PK"
	assert.Equal(t, byte('P'), bufferBytes[0])
	assert.Equal(t, byte('K'), bufferBytes[1])
}

func TestParseKMZBuffer(t *testing.T) {
	// First create a KMZ buffer
	waylines := &Waylines{
		Name:                    "Parse Test Mission",
		DroneModel:              DroneM3DSeries,
		PayloadModel:            PayloadMatrice3TD,
		TemplateType:            TemplateTypeWaypoint,
		GlobalHeight:            60.0,
		GlobalSpeed:             6.0,
		FinishAction:            FinishActionGoHome,
		ClimbMode:               "vertical",
		SafeHeight:              30.0,
		GlobalRTHHeight:         100.0,
		AircraftYawMode:         "followWayline",
		GimbalPitchMode:         "usePointSetting",
		GlobalTransitionalSpeed: 6.0,
		Waypoints: []WaylinesWaypoint{
			{
				Latitude:    35.6762,
				Longitude:   139.6503,
				Height:      60.0,
				Speed:       6.0,
				TriggerType: "reachPoint",
			},
			{
				Latitude:    35.6772,
				Longitude:   139.6513,
				Height:      70.0,
				Speed:       7.0,
				TriggerType: "reachPoint",
			},
		},
	}

	originalMission, err := ConvertWaylinesToWPMLMission(waylines)
	require.NoError(t, err)

	buffer, err := CreateKmzBuffer(originalMission)
	require.NoError(t, err)

	// Test ParseKMZBuffer
	parsedMission, err := ParseKMZBuffer(buffer.Bytes())
	require.NoError(t, err)
	require.NotNil(t, parsedMission)

	// Verify basic properties
	assert.Equal(t, originalMission.Template.Document.Author, parsedMission.Template.Document.Author)
	assert.NotNil(t, parsedMission.Template.Document.MissionConfig)
	assert.NotNil(t, parsedMission.Template)
	assert.NotNil(t, parsedMission.Waylines)
}

func TestCreateKmzErrors(t *testing.T) {
	tests := []struct {
		name        string
		mission     *WPMLMission
		kmzPath     string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Nil mission",
			mission:     nil,
			kmzPath:     "/tmp/test.kmz",
			expectError: true,
		},
		{
			name: "Invalid path",
			mission: &WPMLMission{
				Template: &TemplateDocument{},
			},
			kmzPath:     "/invalid/path/that/doesnt/exist/and/cannot/be/created/test.kmz",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CreateKmz(tt.mission, tt.kmzPath)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
				// Clean up
				_ = os.Remove(tt.kmzPath)
			}
		})
	}
}

func TestKmzWithComplexWaylines(t *testing.T) {
	waylines := &Waylines{
		Name:                     "Complex Test Mission",
		Description:              "A complex mission with multiple features",
		DroneModel:               DroneM3DSeries,
		PayloadModel:             PayloadMatrice3TD,
		TemplateType:             TemplateTypeWaypoint,
		GlobalHeight:             100.0,
		GlobalSpeed:              10.0,
		HeightType:               HeightModeRelativeToStartPoint,
		FinishAction:             FinishActionAutoLand,
		SafeHeight:               50.0,
		GlobalRTHHeight:          120.0,
		AircraftYawMode:          "followWayline",
		GimbalPitchMode:          "usePointSetting",
		GlobalTransitionalSpeed:  8.0,
		TakeOffRefPointLatitude:  39.9042,
		TakeOffRefPointLongitude: 116.4074,
		TakeOffRefPointHeight:    10.0,
		GlobalWaypointTurnMode:   "coordinateTurn",
		PhotoSettings:            []string{"wide", "zoom"},
		UseLowLightSmart:         true,
		ClimbMode:                "vertical",
		Waypoints: []WaylinesWaypoint{
			{
				Latitude:         39.9042,
				Longitude:        116.4074,
				Height:           100.0,
				Speed:            10.0,
				TriggerType:      "reachPoint",
				WaypointTurnMode: "coordinateTurn",
				Actions: []ActionRequest{
					{
						Type:   ActionTypeTakePhoto,
						Action: &TakePhotoAction{PayloadPositionIndex: PayloadPosition0},
					},
					{
						Type: ActionTypeGimbalRotate,
						Action: &GimbalRotateAction{
							PayloadPositionIndex:    PayloadPosition0,
							GimbalHeadingYawBase:    "north",
							GimbalRotateMode:        "absoluteAngle",
							GimbalPitchRotateEnable: true,
							GimbalPitchRotateAngle:  -45.0,
							GimbalYawRotateEnable:   true,
							GimbalYawRotateAngle:    0.0,
							GimbalRotateTimeEnable:  true,
							GimbalRotateTime:        2.0,
						},
					},
				},
			},
			{
				Latitude:         39.9052,
				Longitude:        116.4084,
				Height:           120.0,
				Speed:            12.0,
				TriggerType:      "passPoint",
				WaypointTurnMode: "toPointAndStopWithContinuityCurvature",
				Actions: []ActionRequest{
					{
						Type: ActionTypeHover,
						Action: &HoverAction{
							HoverTime: 5.0,
						},
					},
					{
						Type: ActionTypeRotateYaw,
						Action: &RotateYawAction{
							AircraftHeading: 90.0,
						},
					},
				},
			},
		},
	}

	mission, err := ConvertWaylinesToWPMLMission(waylines)
	require.NoError(t, err)

	// Create KMZ buffer
	buffer, err := CreateKmzBuffer(mission)
	require.NoError(t, err)
	require.NotNil(t, buffer)

	// Verify buffer contains expected content
	assert.Greater(t, buffer.Len(), 1000, "Complex mission should generate substantial KMZ")

	// Test parsing back
	parsedMission, err := ParseKMZBuffer(buffer.Bytes())
	require.NoError(t, err)
	require.NotNil(t, parsedMission)

	// Verify mission config
	assert.NotNil(t, parsedMission.Template.Document.MissionConfig)
	assert.Equal(t, FinishActionAutoLand, parsedMission.Template.Document.MissionConfig.FinishAction)

	// Verify waylines structure
	assert.NotEmpty(t, parsedMission.Waylines.Document.Folders)
	waylineFolder := parsedMission.Waylines.Document.Folders[0]
	assert.NotEmpty(t, waylineFolder.Placemarks)

	// Verify we have waypoints with actions
	assert.Greater(t, len(waylineFolder.Placemarks), 0)
	firstPlacemark := waylineFolder.Placemarks[0]
	assert.NotNil(t, firstPlacemark)
}
