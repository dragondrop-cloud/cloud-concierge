package terraformimportmigrationgenerator

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	terraformValueObjects "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_value_objects"
)

func TestCreateIsolatedTerraformImportMigrationGenerator(t *testing.T) {
	// Given
	ctx := context.Background()
	config := Config{}
	terraformImporterProvider := "isolated"
	terraformImporterFactory := new(Factory)
	provider := terraformValueObjects.Provider("")

	// When
	terraformImporter, err := terraformImporterFactory.Instantiate(ctx, terraformImporterProvider, provider, config)

	// Then
	assert.Nil(t, err)
	assert.NotNil(t, terraformImporter)
}
