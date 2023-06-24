package terraformWorkspace

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dragondrop-cloud/cloud-concierge/main/internal/interfaces"
)

func TestCreateIsolatedTerraformWorkspace(t *testing.T) {
	// Given
	ctx := context.Background()
	config := TerraformCloudConfig{}
	terraformWorkspaceProvider := "isolated"
	terraformWorkspaceFactory := new(Factory)
	dragonDrop := new(interfaces.DragonDropMock)

	// When
	terraformWorkspace, err := terraformWorkspaceFactory.Instantiate(ctx, terraformWorkspaceProvider, dragonDrop, config)

	// Then
	assert.Nil(t, err)
	assert.NotNil(t, terraformWorkspace)
}
