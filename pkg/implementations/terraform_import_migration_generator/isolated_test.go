package terraformImportMigrationGenerator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewIsolatedTerraformImportMigrationGenerator(t *testing.T) {
	// When
	terraformImport := NewIsolatedTerraformImportMigrationGenerator()

	// Then
	assert.NotNil(t, terraformImport)
}
