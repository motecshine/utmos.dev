package wpml

import (
	"bytes"
	"encoding/xml"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewXMLSerializer(t *testing.T) {
	serializer := NewXMLSerializer(true)
	require.NotNil(t, serializer)
	assert.True(t, serializer.indent)

	serializer = NewXMLSerializer(false)
	require.NotNil(t, serializer)
	assert.False(t, serializer.indent)
}

func TestXMLSerializer_Marshal(t *testing.T) {
	serializer := NewXMLSerializer(true)

	doc := &Document{
		XMLNS:  "http://www.opengis.net/kml/2.2",
		WPMLNS: "http://www.dji.com/wpmz/1.0.6",
		Document: DocumentContent{
			Author:     "Test Author",
			CreateTime: 1640995200000,
			UpdateTime: 1640995200000,
			MissionConfig: &MissionConfig{
				FlyToWaylineMode:        FlightModeSafely,
				FinishAction:            FinishActionGoHome,
				ExitOnRCLost:            RCLostActionGoContinue,
				TakeOffSecurityHeight:   30.0,
				GlobalTransitionalSpeed: 5.0,
				DroneInfo: DroneInfo{
					DroneEnumValue:    91,
					DroneSubEnumValue: 0,
				},
				PayloadInfo: PayloadInfo{
					PayloadEnumValue:     81,
					PayloadPositionIndex: 0,
				},
			},
		},
	}

	xmlData, err := serializer.Marshal(doc)
	require.NoError(t, err)
	require.NotEmpty(t, xmlData)

	// Verify basic XML structure
	xmlString := string(xmlData)
	assert.Contains(t, xmlString, `<?xml version="1.0" encoding="UTF-8"?>`)
	assert.Contains(t, xmlString, `xmlns="http://www.opengis.net/kml/2.2"`)
	assert.Contains(t, xmlString, `xmlns:wpml="http://www.dji.com/wpmz/1.0.6"`)
	assert.Contains(t, xmlString, `<wpml:author>Test Author</wpml:author>`)

	// Verify it's valid XML
	var parsedDoc Document
	err = xml.Unmarshal(xmlData, &parsedDoc)
	assert.NoError(t, err)
}

func TestXMLSerializer_Unmarshal(t *testing.T) {
	serializer := NewXMLSerializer(false)

	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<kml xmlns="http://www.opengis.net/kml/2.2" xmlns:wpml="http://www.dji.com/wpmz/1.0.6">
	<Document>
		<wpml:author>Test Author</wpml:author>
		<wpml:createTime>1640995200000</wpml:createTime>
		<wpml:updateTime>1640995200000</wpml:updateTime>
		<wpml:missionConfig>
			<wpml:flyToWaylineMode>safely</wpml:flyToWaylineMode>
			<wpml:finishAction>goHome</wpml:finishAction>
			<wpml:exitOnRCLost>goContinue</wpml:exitOnRCLost>
			<wpml:takeOffSecurityHeight>30</wpml:takeOffSecurityHeight>
			<wpml:globalTransitionalSpeed>5</wpml:globalTransitionalSpeed>
			<wpml:droneInfo>
				<wpml:droneEnumValue>91</wpml:droneEnumValue>
				<wpml:droneSubEnumValue>0</wpml:droneSubEnumValue>
			</wpml:droneInfo>
			<wpml:payloadInfo>
				<wpml:payloadEnumValue>81</wpml:payloadEnumValue>
				<wpml:payloadPositionIndex>0</wpml:payloadPositionIndex>
			</wpml:payloadInfo>
		</wpml:missionConfig>
	</Document>
</kml>`

	var doc Document
	err := serializer.Unmarshal([]byte(xmlData), &doc)
	require.NoError(t, err)

	assert.Equal(t, "Test Author", doc.Document.Author)
	assert.Equal(t, int64(1640995200000), doc.Document.CreateTime)
	assert.Equal(t, FlightModeSafely, doc.Document.MissionConfig.FlyToWaylineMode)
	assert.Equal(t, FinishActionGoHome, doc.Document.MissionConfig.FinishAction)
	assert.Equal(t, 91, doc.Document.MissionConfig.DroneInfo.DroneEnumValue)
	assert.Equal(t, 81, doc.Document.MissionConfig.PayloadInfo.PayloadEnumValue)
}

func TestXMLSerializer_MarshalToWriter(t *testing.T) {
	serializer := NewXMLSerializer(true)

	doc := &Document{
		XMLNS:  "http://www.opengis.net/kml/2.2",
		WPMLNS: "http://www.dji.com/wpmz/1.0.6",
		Document: DocumentContent{
			Author: "Writer Test",
		},
	}

	var buf bytes.Buffer
	err := serializer.MarshalToWriter(doc, &buf)
	require.NoError(t, err)

	xmlString := buf.String()
	assert.Contains(t, xmlString, `<?xml version="1.0" encoding="UTF-8"?>`)
	assert.Contains(t, xmlString, `<wpml:author>Writer Test</wpml:author>`)
}

func TestXMLSerializer_UnmarshalFromReader(t *testing.T) {
	serializer := NewXMLSerializer(false)

	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<kml xmlns="http://www.opengis.net/kml/2.2" xmlns:wpml="http://www.dji.com/wpmz/1.0.6">
	<Document>
		<wpml:author>Reader Test</wpml:author>
		<wpml:createTime>1640995200000</wpml:createTime>
	</Document>
</kml>`

	reader := strings.NewReader(xmlData)
	var doc Document
	err := serializer.UnmarshalFromReader(reader, &doc)
	require.NoError(t, err)

	assert.Equal(t, "Reader Test", doc.Document.Author)
	assert.Equal(t, int64(1640995200000), doc.Document.CreateTime)
}

func TestXMLSerializer_InvalidXML(t *testing.T) {
	serializer := NewXMLSerializer(false)

	invalidXML := `<?xml version="1.0" encoding="UTF-8"?>
<kml>
	<Document>
		<unclosed_tag>
	</Document>
</kml>`

	var doc Document
	err := serializer.Unmarshal([]byte(invalidXML), &doc)
	assert.Error(t, err)
}

func TestXMLSerializer_EmptyInput(t *testing.T) {
	serializer := NewXMLSerializer(false)

	var doc Document
	err := serializer.Unmarshal([]byte(""), &doc)
	assert.Error(t, err)
}

func TestXMLFormatting_Indented(t *testing.T) {
	serializer := NewXMLSerializer(true)

	doc := &Document{
		XMLNS:  "http://www.opengis.net/kml/2.2",
		WPMLNS: "http://www.dji.com/wpmz/1.0.6",
		Document: DocumentContent{
			Author: "Formatting Test",
		},
	}

	xmlData, err := serializer.Marshal(doc)
	require.NoError(t, err)

	xmlString := string(xmlData)

	// Check for proper XML formatting
	assert.True(t, strings.Contains(xmlString, "\n"), "XML should be formatted with newlines")
	assert.True(t, strings.Contains(xmlString, "  "), "XML should be formatted with indentation")

	// Should start with XML declaration
	assert.True(t, strings.HasPrefix(xmlString, `<?xml version="1.0" encoding="UTF-8"?>`))
}

func TestXMLFormatting_Compact(t *testing.T) {
	serializer := NewXMLSerializer(false)

	doc := &Document{
		XMLNS:  "http://www.opengis.net/kml/2.2",
		WPMLNS: "http://www.dji.com/wpmz/1.0.6",
		Document: DocumentContent{
			Author: "Compact Test",
		},
	}

	xmlData, err := serializer.Marshal(doc)
	require.NoError(t, err)

	xmlString := string(xmlData)

	// Should have fewer newlines and no indentation
	lines := strings.Split(xmlString, "\n")
	assert.LessOrEqual(t, len(lines), 3, "Compact XML should have fewer lines")
}

func TestRoundTripSerialization(t *testing.T) {
	serializer := NewXMLSerializer(true)

	// Create a document
	original := &Document{
		XMLNS:  "http://www.opengis.net/kml/2.2",
		WPMLNS: "http://www.dji.com/wpmz/1.0.6",
		Document: DocumentContent{
			Author:     "Roundtrip Test",
			CreateTime: 1640995200000,
			UpdateTime: 1640995200000,
		},
	}

	// Marshal to XML
	xmlData, err := serializer.Marshal(original)
	require.NoError(t, err)

	// Unmarshal back
	var parsed Document
	err = serializer.Unmarshal(xmlData, &parsed)
	require.NoError(t, err)

	// Compare key fields
	assert.Equal(t, original.Document.Author, parsed.Document.Author)
	assert.Equal(t, original.Document.CreateTime, parsed.Document.CreateTime)
	assert.Equal(t, original.XMLNS, parsed.XMLNS)
	assert.Equal(t, original.WPMLNS, parsed.WPMLNS)
}
