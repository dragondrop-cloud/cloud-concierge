package resourcescalculator

import (
	"context"
	"testing"

	terraformValueObjects "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_value_objects"
	"github.com/stretchr/testify/assert"
)

func TestCreateNotIsolated(t *testing.T) {
	// Given
	ctx := context.Background()
	environment := "not_isolated"
	resourcesCalculatorFactory := new(Factory)
	provider := terraformValueObjects.Provider("provider")

	// When
	calculator, err := resourcesCalculatorFactory.Instantiate(ctx, environment, provider)

	// Then
	assert.Nil(t, err)
	assert.NotNil(t, calculator)
}

func TestCreateIsolatedResourcesCalculator(t *testing.T) {
	// Given
	ctx := context.Background()
	environment := "isolated"
	resourcesCalculatorFactory := new(Factory)
	provider := terraformValueObjects.Provider("")

	// When
	calculator, err := resourcesCalculatorFactory.Instantiate(ctx, environment, provider)

	// Then
	assert.Nil(t, err)
	assert.NotNil(t, calculator)
}
