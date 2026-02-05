package wpml

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWPMLMission_AddResource(t *testing.T) {
	mission := &WPMLMission{}
	filename := "test.jpg"
	data := []byte("test image data")

	mission.AddResource(filename, data)

	assert.NotNil(t, mission.Resources)
	assert.Contains(t, mission.Resources, filename)
	assert.Equal(t, data, mission.Resources[filename])
}

func TestWPMLMission_AddResource_Multiple(t *testing.T) {
	mission := &WPMLMission{}

	mission.AddResource("file1.jpg", []byte("data1"))
	mission.AddResource("file2.jpg", []byte("data2"))

	assert.Len(t, mission.Resources, 2)
	assert.Equal(t, []byte("data1"), mission.Resources["file1.jpg"])
	assert.Equal(t, []byte("data2"), mission.Resources["file2.jpg"])
}
