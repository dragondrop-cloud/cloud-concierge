package dragonDrop

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateIsolatedDragonDrop(t *testing.T) {
	// Given
	config := HTTPDragonDropClientConfig{}
	dragonDropProtocol := "isolated"
	dragonDropFactory := new(Factory)

	// When
	dragonDrop, err := dragonDropFactory.Instantiate(dragonDropProtocol, config)

	// Then
	assert.Nil(t, err)
	assert.NotNil(t, dragonDrop)
}
