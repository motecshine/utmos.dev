package wpml

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatXML(t *testing.T) {
	// Create a simple valid XML string
	xmlString := `<kml xmlns="http://www.opengis.net/kml/2.2"><Document><name>Test</name></Document></kml>`

	formatted, err := FormatXML(xmlString)

	assert.NoError(t, err)
	assert.Contains(t, formatted, "Test")
	assert.Contains(t, formatted, "kml")
	// Should be formatted with indentation
	assert.Contains(t, formatted, "\n")
}

func TestFormatXML_InvalidXML(t *testing.T) {
	// Test with invalid XML
	invalidXML := "<invalid><unclosed>"

	_, err := FormatXML(invalidXML)
	assert.Error(t, err)
}

func TestValidateXML(t *testing.T) {
	// Test with valid XML
	validXML := []byte(`<kml xmlns="http://www.opengis.net/kml/2.2"><Document><name>Test</name></Document></kml>`)

	err := ValidateXML(validXML)
	assert.NoError(t, err)
}

func TestValidateXML_Invalid(t *testing.T) {
	// Test with invalid XML
	invalidXML := []byte("<invalid><unclosed>")

	err := ValidateXML(invalidXML)
	assert.Error(t, err)
}

func TestGetXMLElements(t *testing.T) {
	// Create XML with multiple elements
	xmlString := `<kml xmlns="http://www.opengis.net/kml/2.2">
		<Document>
			<Placemark><name>Point1</name></Placemark>
			<Placemark><name>Point2</name></Placemark>
		</Document>
	</kml>`

	elements, err := GetXMLElements([]byte(xmlString), "Placemark")

	assert.NoError(t, err)
	assert.Len(t, elements, 2)
	assert.Contains(t, elements[0], "Point1")
	assert.Contains(t, elements[1], "Point2")
}

func TestGetXMLElements_NotFound(t *testing.T) {
	xmlString := `<kml xmlns="http://www.opengis.net/kml/2.2"><Document><name>Test</name></Document></kml>`

	elements, err := GetXMLElements([]byte(xmlString), "NonExistent")

	assert.NoError(t, err)
	assert.Empty(t, elements)
}

func TestGetXMLElements_InvalidXML(t *testing.T) {
	invalidXML := []byte("<invalid><unclosed>")

	elements, err := GetXMLElements(invalidXML, "Placemark")

	assert.Error(t, err)
	assert.Empty(t, elements)
}
