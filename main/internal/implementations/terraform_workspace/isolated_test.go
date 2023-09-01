package terraformworkspace

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewIsolatedTerraformWorkspace(t *testing.T) {
	// When
	workspace := NewIsolatedTerraformWorkspace()

	// Then
	assert.NotNil(t, workspace)
}
