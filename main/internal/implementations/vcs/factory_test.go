package vcs

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateNotSupported(t *testing.T) {
	// Given
	ctx := context.Background()
	config := Config{}
	vcsProvider := "not_supported"
	vcsFactory := new(Factory)

	// When
	vcs, err := vcsFactory.Instantiate(ctx, vcsProvider, config, "")

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

	// When
	vcs, err := vcsFactory.Instantiate(ctx, vcsProvider, config, "github")

	// Then
	assert.Nil(t, err)
	assert.NotNil(t, vcs)
}
