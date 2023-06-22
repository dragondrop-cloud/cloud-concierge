package interfaces

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// TerraformImportMigrationGenerator is an interface for generating terraform state migration import
// statements.
type TerraformImportMigrationGenerator interface {

	// Execute generates terraform state migration statements for identified resources.
	Execute(ctx context.Context) error
}

// TerraformImportMigrationGeneratorMock implements the TerraformImportMigrationGenerator interface for unit testing purposes.
type TerraformImportMigrationGeneratorMock struct {
	mock.Mock
}

// Execute generates terraform state migration statements for identified resources.
func (m *TerraformImportMigrationGeneratorMock) Execute(ctx context.Context) error {
	args := m.Called()
	return args.Error(0)
}
