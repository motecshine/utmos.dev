package wpml

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarshalTemplate(t *testing.T) {
	waylines := createValidWaylines("Template Marshal Test")

	mission, err := ConvertWaylinesToWPMLMission(waylines)
	require.NoError(t, err)

	// Test MarshalTemplate
	xmlBytes, err := MarshalTemplate(mission.Template)

	assert.NoError(t, err)
	assert.NotEmpty(t, xmlBytes)
	assert.Contains(t, string(xmlBytes), "wpml:templateType")
	assert.Contains(t, string(xmlBytes), "Placemark")
}

func TestMarshalWaylines(t *testing.T) {
	waylines := createValidWaylines("Waylines Marshal Test")

	mission, err := ConvertWaylinesToWPMLMission(waylines)
	require.NoError(t, err)

	// Test MarshalWaylines
	xmlBytes, err := MarshalWaylines(mission.Waylines)

	assert.NoError(t, err)
	assert.NotEmpty(t, xmlBytes)
	assert.Contains(t, string(xmlBytes), "wpml:waylineId")
	assert.Contains(t, string(xmlBytes), "Placemark")
}

func TestUnmarshalTemplate(t *testing.T) {
	waylines := createValidWaylines("Unmarshal Template Test")

	mission, err := ConvertWaylinesToWPMLMission(waylines)
	require.NoError(t, err)

	xmlBytes, err := MarshalTemplate(mission.Template)
	require.NoError(t, err)

	// Test UnmarshalTemplate
	template, err := UnmarshalTemplate(xmlBytes)

	assert.NoError(t, err)
	assert.NotNil(t, template)
	// Verify some content was preserved
	assert.Equal(t, "http://www.opengis.net/kml/2.2", template.XMLNS)
}

func TestUnmarshalTemplate_InvalidXML(t *testing.T) {
	invalidXML := []byte("<invalid><xml>unclosed")

	template, err := UnmarshalTemplate(invalidXML)

	assert.Error(t, err)
	assert.Nil(t, template)
}

func TestUnmarshalWaylines(t *testing.T) {
	waylines := createValidWaylines("Unmarshal Waylines Test")

	mission, err := ConvertWaylinesToWPMLMission(waylines)
	require.NoError(t, err)

	xmlBytes, err := MarshalWaylines(mission.Waylines)
	require.NoError(t, err)

	// Test UnmarshalWaylines
	waylinesDoc, err := UnmarshalWaylines(xmlBytes)

	assert.NoError(t, err)
	assert.NotNil(t, waylinesDoc)
	// Verify some content was preserved
	assert.Equal(t, "http://www.opengis.net/kml/2.2", waylinesDoc.XMLNS)
}

func TestUnmarshalWaylines_InvalidXML(t *testing.T) {
	invalidXML := []byte("<invalid><xml>unclosed")

	waylinesDoc, err := UnmarshalWaylines(invalidXML)

	assert.Error(t, err)
	assert.Nil(t, waylinesDoc)
}
