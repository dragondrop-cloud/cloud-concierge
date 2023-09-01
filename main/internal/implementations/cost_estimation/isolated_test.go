package costestimation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewIsolatedCostEstimator(t *testing.T) {
	// When
	dragonDrop := NewIsolatedCostEstimator()

	// Then
	assert.NotNil(t, dragonDrop)
}
