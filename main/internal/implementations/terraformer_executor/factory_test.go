package terraformerexecutor

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dragondrop-cloud/cloud-concierge/main/internal/hclcreate"
	terraformValueObjects "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_value_objects"
	terraformerCli "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraformer_executor/terraformer_cli"
)

func TestCreateIsolatedTerraformerExecutor(t *testing.T) {
	// Given
	ctx := context.Background()
	hclConfig := hclcreate.Config{}
	executorConfig := terraformerCli.TerraformerExecutorConfig{}
	cliConfig := terraformerCli.Config{}
	terraformerExecutorProvider := "isolated"
	terraformerExecutorFactory := new(Factory)
	provider := terraformValueObjects.Provider("")

	// When
	terraformerExecutor, err := terraformerExecutorFactory.Instantiate(ctx, terraformerExecutorProvider, provider, hclConfig, executorConfig, cliConfig)

	// Then
	assert.Nil(t, err)
	assert.NotNil(t, terraformerExecutor)
}
