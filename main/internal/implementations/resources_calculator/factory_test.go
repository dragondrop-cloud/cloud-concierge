package resourcesCalculator

import (
	"context"
	"testing"

	terraformValueObjects "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_value_objects"
	"github.com/stretchr/testify/assert"

	"github.com/dragondrop-cloud/cloud-concierge/main/internal/interfaces"
)

func TestCreateNotIsolated(t *testing.T) {
	// Given
	ctx := context.Background()
	environment := "not_isolated"
	resourcesCalculatorFactory := new(Factory)
	dragonDrop := new(interfaces.DragonDropMock)
	provider := terraformValueObjects.Provider("provider")

	// When
	calculator, err := resourcesCalculatorFactory.Instantiate(ctx, environment, dragonDrop, provider)

	// Then
	assert.Nil(t, err)
	assert.NotNil(t, calculator)
}

func TestCreateIsolatedResourcesCalculator(t *testing.T) {
	// Given
	ctx := context.Background()
	environment := "isolated"
	resourcesCalculatorFactory := new(Factory)
	dragonDrop := new(interfaces.DragonDropMock)
	provider := terraformValueObjects.Provider("")

	// When
	calculator, err := resourcesCalculatorFactory.Instantiate(ctx, environment, dragonDrop, provider)

	// Then
	assert.Nil(t, err)
	assert.NotNil(t, calculator)
}
