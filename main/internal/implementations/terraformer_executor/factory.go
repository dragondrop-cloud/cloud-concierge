package terraformerExecutor

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/dragondrop-cloud/cloud-concierge/main/internal/hclcreate"
	terraformValueObjects "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_value_objects"
	terraformerCli "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraformer_executor/terraformer_cli"
	"github.com/dragondrop-cloud/cloud-concierge/main/internal/interfaces"
)

// Factory is a struct that generates implementations of interfaces.TerraformerExecutor
type Factory struct {
}

// Instantiate returns an implementation of interfaces.TerraformerExecutor depending on the passed
// environment specification.
func (f *Factory) Instantiate(ctx context.Context, environment string, dragonDrop interfaces.DragonDrop, divisionToProvider map[terraformValueObjects.Division]terraformValueObjects.Provider, hclConfig hclcreate.Config, executorConfig terraformerCli.TerraformerExecutorConfig, cliConfig terraformerCli.Config) (interfaces.TerraformerExecutor, error) {
	switch environment {
	case "isolated":
		return new(IsolatedTerraformerExecutor), nil
	default:
		return f.bootstrappedTerraformerExecutor(ctx, dragonDrop, divisionToProvider, hclConfig, executorConfig, cliConfig)
	}
}

// bootstrappedTerraformerExecutor creates a complete implementation of the interfaces.TerraformerExecutor interface with
// configuration specified via environment variables.
func (f *Factory) bootstrappedTerraformerExecutor(ctx context.Context, dragonDrop interfaces.DragonDrop, divisionToProvider map[terraformValueObjects.Division]terraformValueObjects.Provider, hclConfig hclcreate.Config, executorConfig terraformerCli.TerraformerExecutorConfig, cliConfig terraformerCli.Config) (interfaces.TerraformerExecutor, error) {
	hclCreate, err := hclcreate.NewHCLCreate(hclConfig, divisionToProvider)
	if err != nil {
		log.Errorf("[cannot instantiate hclCreate config]%s", err.Error())
		return nil, fmt.Errorf("[cannot instantiate hclCreate config]%w", err)
	}

	return terraformerCli.NewTerraformerExecutor(ctx, hclCreate, dragonDrop, executorConfig, cliConfig, divisionToProvider)
}
