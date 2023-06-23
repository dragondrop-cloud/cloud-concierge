package vcs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractOrgAndRepoName(t *testing.T) {
	// Given
	github := &GitHub{}
	input := "https://github.com/dragondrop-cloud-org/dragondrop-cloud-repo1.git"

	// When
	org, repo, err := github.extractOrgAndRepoName(input)

	// Then
	assert.Nil(t, err)
	assert.Equal(t, "dragondrop-cloud-org", org)
	assert.Equal(t, "dragondrop-cloud-repo1", repo)
}
