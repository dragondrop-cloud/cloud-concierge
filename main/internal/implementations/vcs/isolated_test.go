package vcs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateIsolatedVCS_Success(t *testing.T) {
	// Given
	initialCloneTimes := uint(0)

	// When
	isolatedVCS := NewIsolatedVCS(initialCloneTimes)

	// Then
	assert.NotNil(t, isolatedVCS)
	assert.Equal(t, uint(0), isolatedVCS.(*IsolatedVCS).CloneTimes)
}

func TestIsolatedVCS_CloneSuccess(t *testing.T) {
	// Given
	isolatedVCS := NewIsolatedVCS(uint(0))

	// When
	err := isolatedVCS.Clone()

	// Then
	assert.Nil(t, err)
	assert.NotNil(t, isolatedVCS)
	assert.Equal(t, uint(1), isolatedVCS.(*IsolatedVCS).CloneTimes)
}
