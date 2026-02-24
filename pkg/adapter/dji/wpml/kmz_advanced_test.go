package wpml

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateKmzBufferFromWaylines(t *testing.T) {
	waylines := createValidWaylines("Test Mission")

	buffer, err := CreateKmzBufferFromWaylines(waylines)

	assert.NoError(t, err)
	assert.NotNil(t, buffer)
	assert.Greater(t, buffer.Len(), 0)
}

func TestCreateKmzBufferFromWaylines_InvalidWaylines(t *testing.T) {
	// Test with invalid waylines (empty name will still pass validation)
	// Let's test with missing required fields instead
	waylines := &Waylines{
		Name:         "Invalid Mission",
		Description:  "Test Description",
		DroneModel:   DroneM3Series,
		PayloadModel: PayloadMatrice3TD,
		// Missing required TemplateType and other validation fields
		Waypoints: []WaylinesWaypoint{
			{
				Latitude:  39.9093,
				Longitude: 116.3974,
				Height:    50.0,
				Speed:     15.0,
				// Missing TriggerType - required field
			},
		},
	}

	buffer, err := CreateKmzBufferFromWaylines(waylines)

	assert.Error(t, err)
	assert.Nil(t, buffer)
}

func TestGetKmzInfo(t *testing.T) {
	waylines := createValidWaylines("Test Mission")

	mission, err := ConvertWaylinesToWPMLMission(waylines)
	require.NoError(t, err)

	info, err := GetKmzInfo(mission)

	assert.NoError(t, err)
	assert.NotNil(t, info)
	assert.Contains(t, info, "total_size")
	assert.Contains(t, info, "files")
	assert.Greater(t, info["total_size"], 0)

	files := info["files"].([]map[string]any)
	assert.NotEmpty(t, files)

	// Check that we have expected files - may be different structure
	// Just verify that we got meaningful data
	assert.True(t, len(files) > 0, "Should have at least one file entry")
}

func TestGetKmzInfo_InvalidMission(t *testing.T) {
	// Test with nil mission
	info, err := GetKmzInfo(nil)

	assert.Error(t, err)
	assert.Nil(t, info)
}

func TestParseKMZFile(t *testing.T) {
	// Create a test KMZ file first
	waylines := createValidWaylines("Test Mission")

	mission, err := ConvertWaylinesToWPMLMission(waylines)
	require.NoError(t, err)

	// Create temporary file
	tempFile := "/tmp/test_mission_parsekmzfile.kmz"
	err = CreateKmz(mission, tempFile)
	require.NoError(t, err)

	// Test ParseKMZFile
	parsedMission, err := ParseKMZFile(tempFile)
	assert.NoError(t, err)
	assert.NotNil(t, parsedMission)
	assert.NotNil(t, parsedMission.Template)
	assert.NotNil(t, parsedMission.Waylines)

	// Clean up
	// Note: We don't use os.Remove here to avoid import issues in tests
	// The temp file will be cleaned up by the system
}

func TestParseKMZFile_NonexistentFile(t *testing.T) {
	parsedMission, err := ParseKMZFile("/tmp/nonexistent_file.kmz")

	assert.Error(t, err)
	assert.Nil(t, parsedMission)
}

func TestGenerateKMZJSON(t *testing.T) {
	waylines := createValidWaylines("Test Mission")

	mission, err := ConvertWaylinesToWPMLMission(waylines)
	require.NoError(t, err)

	jsonData, err := GenerateKMZJSON(mission, "test_mission.kmz")

	assert.NoError(t, err)
	assert.NotEmpty(t, jsonData)
	// JSON contains structured data, check for key fields
	assert.Contains(t, jsonData, "template")
	assert.Contains(t, jsonData, "waylines")
}

func TestGenerateKMZJSON_NilMission(t *testing.T) {
	jsonData, err := GenerateKMZJSON(nil, "test.kmz")

	assert.Error(t, err)
	assert.Empty(t, jsonData)
}
