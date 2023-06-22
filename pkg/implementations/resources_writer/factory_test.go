package resourcesWriter

import (
	"context"
	"testing"

	"github.com/dragondrop-cloud/driftmitigation/hclcreate"
	terraformValueObjects "github.com/dragondrop-cloud/driftmitigation/implementations/terraform_value_objects"
	"github.com/dragondrop-cloud/driftmitigation/interfaces"
	"github.com/stretchr/testify/assert"
)

func TestCreateIsolatedResourcesWriter(t *testing.T) {
	// Given
	ctx := context.Background()
	hclConfig := hclcreate.Config{}
	resourcesWriterProvider := "isolated"
	resourcesWriterFactory := new(Factory)
	vcs := new(interfaces.VCSMock)
	dragonDrop := new(interfaces.DragonDropMock)
	divisionToProvider := make(map[terraformValueObjects.Division]terraformValueObjects.Provider)

	// When
	resourcesWriter, err := resourcesWriterFactory.Instantiate(ctx, resourcesWriterProvider, vcs, dragonDrop, divisionToProvider, hclConfig)

	// Then
	assert.Nil(t, err)
	assert.NotNil(t, resourcesWriter)
}
