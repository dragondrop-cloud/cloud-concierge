package terraformWorkspace

import (
	"context"

	"github.com/dragondrop-cloud/cloud-concierge/main/internal/interfaces"
)

// Factory is a struct that generates implementations of interfaces.TerraformWorkspace.
type Factory struct {
}

// Instantiate returns an implementation of the interfaces.TerraformWorkspace interface depending on the passed
// environment specification.
func (f *Factory) Instantiate(ctx context.Context, environment string, dragonDrop interfaces.DragonDrop, tfConfig TfStackConfig) (interfaces.TerraformWorkspace, error) {
	switch environment {
	case "isolated":
		return new(IsolatedTerraformWorkspace), nil
	default:
		return f.bootstrappedTerraformWorkspace(ctx, dragonDrop, tfConfig)
	}
}

// bootstrappedTerraformWorkspace creates a complete implementation of the interfaces.TerraformWorkspace interface with
// configuration specified via environment variables.
func (f *Factory) bootstrappedTerraformWorkspace(ctx context.Context, dragonDrop interfaces.DragonDrop, tfStack TfStackConfig) (interfaces.TerraformWorkspace, error) {
	switch tfStack.StateBackend {
	case "s3":
		return NewS3Backend(ctx, tfStack, dragonDrop), nil
	case "azurerm":
		return NewAzurermBlobBackend(ctx, tfStack, dragonDrop), nil
	case "gcs":
		return NewGCSBackend(ctx, tfStack, dragonDrop), nil
	default:
		return NewTerraformCloud(ctx, tfStack, dragonDrop), nil
	}
}
