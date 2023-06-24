package terraformImportMigrationGenerator

import (
	"context"

	"github.com/dragondrop-cloud/cloud-concierge/main/internal/interfaces"
)

// IsolatedTerraformImportMigrationGenerator is a struct that implements the interfaces.TerraformImportMigrationGenerator interface for
// the purpose of end-to-end testing.
type IsolatedTerraformImportMigrationGenerator struct {
}

// NewIsolatedTerraformImportMigrationGenerator returns a new instance of IsolatedTerraformImportMigrationGenerator.
func NewIsolatedTerraformImportMigrationGenerator() interfaces.TerraformImportMigrationGenerator {
	return &IsolatedTerraformImportMigrationGenerator{}
}

// Execute generates terraform state migration statements for identified resources.
func (i *IsolatedTerraformImportMigrationGenerator) Execute(ctx context.Context) error {
	return nil
}
