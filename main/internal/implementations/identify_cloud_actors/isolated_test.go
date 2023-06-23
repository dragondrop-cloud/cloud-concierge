package identifyCloudActors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewIsolatedResourcesCalculator(t *testing.T) {
	// When
	identifier := NewIsolatedIdentifyCloudActors()

	// Then
	assert.NotNil(t, identifier)
}
