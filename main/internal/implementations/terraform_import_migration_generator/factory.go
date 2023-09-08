package terraformimportmigrationgenerator

import (
	"context"

	terraformValueObjects "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_value_objects"
	"github.com/dragondrop-cloud/cloud-concierge/main/internal/interfaces"
)

// Factory is a struct that generates implementations of the interfaces.TerraformImportMigrationGenerator interface.
type Factory struct {
}

// Instantiate returns an implementation of the interfaces.TerraformImportMigrationGenerator interface depending on the passed
// environment specification.
func (f *Factory) Instantiate(ctx context.Context, environment string, dragonDrop interfaces.DragonDrop, provider terraformValueObjects.Provider, config Config) (interfaces.TerraformImportMigrationGenerator, error) {
	switch environment {
	case "isolated":
		return new(IsolatedTerraformImportMigrationGenerator), nil
	default:
		return f.bootstrappedTerraformImportMigrationGenerator(ctx, dragonDrop, provider, config)
	}
}

// bootstrappedTerraformImportMigrationGenerator creates a complete implementation of the TerraformImportMigrationGenerator interface with
// configuration specified via environment variables.
func (f *Factory) bootstrappedTerraformImportMigrationGenerator(ctx context.Context, dragonDrop interfaces.DragonDrop, provider terraformValueObjects.Provider, config Config) (interfaces.TerraformImportMigrationGenerator, error) {
	return NewTerraformImportMigrationGenerator(ctx, config, dragonDrop, provider), nil
}
