package resourcesCalculator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewIsolatedResourcesCalculator(t *testing.T) {
	// When
	calculator := NewIsolatedResourcesCalculator()

	// Then
	assert.NotNil(t, calculator)
}
