package wpml

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewDocument(t *testing.T) {
	doc := NewDocument()

	assert.NotNil(t, doc)
	assert.Equal(t, "http://www.opengis.net/kml/2.2", doc.XMLNS)
	assert.Equal(t, "http://www.dji.com/wpmz/1.0.6", doc.WPMLNS)
	assert.NotZero(t, doc.Document.CreateTime)
	assert.NotZero(t, doc.Document.UpdateTime)
}

func TestDocument_SetAuthor(t *testing.T) {
	doc := NewDocument()
	author := "Test Author"

	doc.SetAuthor(author)

	assert.Equal(t, author, doc.Document.Author)
}

func TestDocument_UpdateTimestamp(t *testing.T) {
	doc := NewDocument()
	originalTime := doc.Document.UpdateTime

	// Sleep to ensure time difference
	time.Sleep(1 * time.Millisecond)

	doc.UpdateTimestamp()

	assert.Greater(t, doc.Document.UpdateTime, originalTime)
}
