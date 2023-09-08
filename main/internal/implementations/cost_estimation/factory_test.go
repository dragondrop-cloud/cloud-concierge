package costestimation

import (
	"testing"

	"github.com/stretchr/testify/assert"

	terraformValueObjects "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_value_objects"
)

func TestCreateNotSupported(t *testing.T) {
	// Given
	config := CostEstimatorConfig{}
	provider := terraformValueObjects.Provider("")
	costEstimatorFactory := new(Factory)

	// When
	costEstimator, err := costEstimatorFactory.Instantiate("", provider, config)

	// Then
	assert.Nil(t, err)
	assert.NotNil(t, costEstimator)
}

func TestCreateIsolatedDragonDrop(t *testing.T) {
	// Given
	config := CostEstimatorConfig{}
	costEstimatorProtocol := "isolated"
	costEstimatorFactory := new(Factory)
	provider := terraformValueObjects.Provider("")

	// When
	costEstimator, err := costEstimatorFactory.Instantiate(costEstimatorProtocol, provider, config)

	// Then
	assert.Nil(t, err)
	assert.NotNil(t, costEstimator)
}
