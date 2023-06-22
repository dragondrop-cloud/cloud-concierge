package terraformImportMigrationGenerator

import (
	"context"
	"testing"

	"github.com/dragondrop-cloud/driftmitigation/interfaces"
	"github.com/stretchr/testify/assert"

	terraformValueObjects "github.com/dragondrop-cloud/driftmitigation/implementations/terraform_value_objects"
)

func TestCreateIsolatedTerraformImportMigrationGenerator(t *testing.T) {
	// Given
	ctx := context.Background()
	config := Config{}
	terraformImporterProvider := "isolated"
	terraformImporterFactory := new(Factory)
	dragonDrop := new(interfaces.DragonDropMock)
	divisionToProvider := make(map[terraformValueObjects.Division]terraformValueObjects.Provider)

	// When
	terraformImporter, err := terraformImporterFactory.Instantiate(ctx, terraformImporterProvider, dragonDrop, divisionToProvider, config)

	// Then
	assert.Nil(t, err)
	assert.NotNil(t, terraformImporter)
}
