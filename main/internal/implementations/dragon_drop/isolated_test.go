package dragondrop

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewIsolatedDragonDrop(t *testing.T) {
	// When
	dragonDrop := NewIsolatedDragonDrop()

	// Then
	assert.NotNil(t, dragonDrop)
}
