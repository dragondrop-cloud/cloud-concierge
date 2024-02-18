package terraformworkspace

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateIsolatedTerraformWorkspace(t *testing.T) {
	// Given
	ctx := context.Background()
	config := TfStackConfig{}
	terraformWorkspaceProvider := "isolated"
	terraformWorkspaceFactory := new(Factory)

	// When
	terraformWorkspace, err := terraformWorkspaceFactory.Instantiate(ctx, terraformWorkspaceProvider, config)

	// Then
	assert.Nil(t, err)
	assert.NotNil(t, terraformWorkspace)
}
