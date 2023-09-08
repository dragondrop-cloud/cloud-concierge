package identifycloudactors

import (
	"context"
	"testing"

	terraformValueObjects "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_value_objects"
	"github.com/dragondrop-cloud/cloud-concierge/main/internal/interfaces"
	"github.com/stretchr/testify/assert"
)

func TestCreateNotIsolated(t *testing.T) {
	// Given
	ctx := context.Background()
	config := Config{}
	env := "not_isolated"
	identifyCloudActorsFactory := new(Factory)
	dragonDrop := new(interfaces.DragonDropMock)
	provider := terraformValueObjects.Provider("gcp")

	// When
	calculator, err := identifyCloudActorsFactory.Instantiate(ctx, env, dragonDrop, provider, config)

	// Then
	assert.Nil(t, err)
	assert.NotNil(t, calculator)
}

func TestCreateIsolatedResourcesCalculator(t *testing.T) {
	// Given
	ctx := context.Background()
	config := Config{}
	env := "isolated"
	identifyCloudActorsFactory := new(Factory)
	dragonDrop := new(interfaces.DragonDropMock)
	provider := terraformValueObjects.Provider("aws")

	// When
	calculator, err := identifyCloudActorsFactory.Instantiate(ctx, env, dragonDrop, provider, config)

	// Then
	assert.Nil(t, err)
	assert.NotNil(t, calculator)
}
