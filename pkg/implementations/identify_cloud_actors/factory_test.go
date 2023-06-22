package identifyCloudActors

import (
	"context"
	"testing"

	terraformValueObjects "github.com/dragondrop-cloud/driftmitigation/implementations/terraform_value_objects"
	"github.com/dragondrop-cloud/driftmitigation/interfaces"
	"github.com/stretchr/testify/assert"
)

func TestCreateNotIsolated(t *testing.T) {
	// Given
	ctx := context.Background()
	config := Config{}
	env := "not_isolated"
	identifyCloudActorsFactory := new(Factory)
	dragonDrop := new(interfaces.DragonDropMock)
	divisionToProvider := make(map[terraformValueObjects.Division]terraformValueObjects.Provider)

	// When
	calculator, err := identifyCloudActorsFactory.Instantiate(ctx, env, dragonDrop, divisionToProvider, config)

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
	divisionToProvider := make(map[terraformValueObjects.Division]terraformValueObjects.Provider)

	// When
	calculator, err := identifyCloudActorsFactory.Instantiate(ctx, env, dragonDrop, divisionToProvider, config)

	// Then
	assert.Nil(t, err)
	assert.NotNil(t, calculator)
}
