package terraformWorkspace

import (
	"context"
	"fmt"

	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"

	"github.com/dragondrop-cloud/cloud-concierge/main/internal/interfaces"
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
	var backendConfig ContainerBackendConfig
	err := envconfig.Process("DRAGONDROP", &backendConfig)
	if err != nil {
		log.Errorf("[cannot create terraform workspace config]%s", err.Error())
		return nil, fmt.Errorf("[cannot create terraform workspace config]%w", err)
	}

	switch config.StateBackend {
	case "s3":
		return NewS3Backend(ctx, backendConfig, dragonDrop), nil
	case "azurerm":
		return NewAzurermBlobBackend(ctx, backendConfig, dragonDrop), nil
	case "gcs":
		return NewGCSBackend(ctx, backendConfig, dragonDrop), nil
	default:
		return NewTerraformCloud(ctx, config, dragonDrop), nil
	}
}
