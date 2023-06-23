package resourcesWriter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewIsolatedResourcesWriter(t *testing.T) {
	// When
	writer := NewIsolatedResourcesWriter()

	// Then
	assert.NotNil(t, writer)
}
