package vcs

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dragondrop-cloud/driftmitigation/interfaces"
)

func TestCreateNotSupported(t *testing.T) {
	// Given
	ctx := context.Background()
	config := Config{}
	vcsProvider := "not_supported"
	vcsFactory := new(Factory)
	dragonDrop := new(interfaces.DragonDropMock)

	// When
	vcs, err := vcsFactory.Instantiate(ctx, vcsProvider, dragonDrop, config)

	// Then
	assert.NotNil(t, err)
	assert.Nil(t, vcs)
}

func TestCreateIsolatedVCS(t *testing.T) {
	// Given
	ctx := context.Background()
	config := Config{}
	vcsProvider := "isolated"
	vcsFactory := new(Factory)
	dragonDrop := new(interfaces.DragonDropMock)

	// When
	vcs, err := vcsFactory.Instantiate(ctx, vcsProvider, dragonDrop, config)

	// Then
	assert.Nil(t, err)
	assert.NotNil(t, vcs)
}
