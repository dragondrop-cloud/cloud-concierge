package terraformerexecutor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewIsolatedTerraformerExecutor(t *testing.T) {
	// When
	executor := NewIsolatedTerraformerExecutor()

	// Then
	assert.NotNil(t, executor)
}
