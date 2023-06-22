package terraformWorkspace

import (
	"context"

	"github.com/dragondrop-cloud/driftmitigation/interfaces"
)

// Factory is a struct that generates implementations of interfaces.TerraformWorkspace.
type Factory struct {
}

// Instantiate returns an implementation of the interfaces.TerraformWorkspace interface depending on the passed
// environment specification.
func (f *Factory) Instantiate(ctx context.Context, environment string, dragonDrop interfaces.DragonDrop, config TerraformCloudConfig) (interfaces.TerraformWorkspace, error) {
	switch environment {
	case "isolated":
		return new(IsolatedTerraformWorkspace), nil
	default:
		return f.bootstrappedTerraformWorkspace(ctx, dragonDrop, config)
	}
}

// bootstrappedTerraformWorkspace creates a complete implementation of the interfaces.TerraformWorkspace interface with
// configuration specified via environment variables.
func (f *Factory) bootstrappedTerraformWorkspace(ctx context.Context, dragonDrop interfaces.DragonDrop, config TerraformCloudConfig) (interfaces.TerraformWorkspace, error) {
	return NewTerraformCloud(ctx, config, dragonDrop), nil
}
