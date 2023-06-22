package costEstimation

import (
	"testing"

	"github.com/stretchr/testify/assert"

	terraformValueObjects "github.com/dragondrop-cloud/driftmitigation/implementations/terraform_value_objects"
)

func TestCreateNotSupported(t *testing.T) {
	// Given
	config := CostEstimatorConfig{}
	divisionToProvider := make(map[terraformValueObjects.Division]terraformValueObjects.Provider)
	costEstimatorFactory := new(Factory)

	// When
	costEstimator, err := costEstimatorFactory.Instantiate("", divisionToProvider, config)

	// Then
	assert.Nil(t, err)
	assert.NotNil(t, costEstimator)
}

func TestCreateIsolatedDragonDrop(t *testing.T) {
	// Given
	config := CostEstimatorConfig{}
	costEstimatorProtocol := "isolated"
	costEstimatorFactory := new(Factory)
	divisionToProvider := make(map[terraformValueObjects.Division]terraformValueObjects.Provider)

	// When
	costEstimator, err := costEstimatorFactory.Instantiate(costEstimatorProtocol, divisionToProvider, config)

	// Then
	assert.Nil(t, err)
	assert.NotNil(t, costEstimator)
}
