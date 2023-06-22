package resourcesCalculator

import (
	"context"
	"testing"

	terraformValueObjects "github.com/dragondrop-cloud/driftmitigation/implementations/terraform_value_objects"
	"github.com/stretchr/testify/assert"

	"github.com/dragondrop-cloud/driftmitigation/interfaces"
)

func TestCreateNotIsolated(t *testing.T) {
	// Given
	ctx := context.Background()
	provider := "not_isolated"
	resourcesCalculatorFactory := new(Factory)
	dragonDrop := new(interfaces.DragonDropMock)
	divisionToProvider := make(map[terraformValueObjects.Division]terraformValueObjects.Provider)

	// When
	calculator, err := resourcesCalculatorFactory.Instantiate(ctx, provider, dragonDrop, divisionToProvider)

	// Then
	assert.Nil(t, err)
	assert.NotNil(t, calculator)
}

func TestCreateIsolatedResourcesCalculator(t *testing.T) {
	// Given
	ctx := context.Background()
	provider := "isolated"
	resourcesCalculatorFactory := new(Factory)
	dragonDrop := new(interfaces.DragonDropMock)
	divisionToProvider := make(map[terraformValueObjects.Division]terraformValueObjects.Provider)

	// When
	calculator, err := resourcesCalculatorFactory.Instantiate(ctx, provider, dragonDrop, divisionToProvider)

	// Then
	assert.Nil(t, err)
	assert.NotNil(t, calculator)
}
