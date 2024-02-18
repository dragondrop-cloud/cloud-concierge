package resourcescalculator

import (
	"context"
	"testing"

	nlpenginerequestor "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/nlp_engine_requester"
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
	nlpEngine, _ := (&nlpenginerequestor.Factory{}).Instantiate(nlpenginerequestor.HTTPNLPEngineClientConfig{})
	calculator, err := resourcesCalculatorFactory.Instantiate(ctx, environment, provider, nlpEngine)

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
	nlpEngine, _ := (&nlpenginerequestor.Factory{}).Instantiate(nlpenginerequestor.HTTPNLPEngineClientConfig{})
	calculator, err := resourcesCalculatorFactory.Instantiate(ctx, environment, provider, nlpEngine)

	// Then
	assert.Nil(t, err)
	assert.NotNil(t, calculator)
}
