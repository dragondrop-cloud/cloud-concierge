package costEstimation

import (
	terraformValueObjects "github.com/dragondrop-cloud/driftmitigation/implementations/terraform_value_objects"
	"github.com/dragondrop-cloud/driftmitigation/interfaces"
)

// Factory is a struct for creating different implementations of interfaces.CostEstimation.
type Factory struct {
}

// Instantiate creates an implementation of interfaces.CostEstimation.
func (f *Factory) Instantiate(environment string, divisionToProvider map[terraformValueObjects.Division]terraformValueObjects.Provider, config CostEstimatorConfig) (interfaces.CostEstimation, error) {
	switch environment {
	case "isolated":
		return new(IsolatedCostEstimator), nil
	default:
		return f.bootstrappedCostEstimator(divisionToProvider, config)
	}
}

// bootstrappedCostEstimator instantiates an instance of CostEstimator with the proper environment
// variables read in.
func (f *Factory) bootstrappedCostEstimator(divisionToProvider map[terraformValueObjects.Division]terraformValueObjects.Provider, config CostEstimatorConfig) (interfaces.CostEstimation, error) {
	return NewCostEstimator(config, divisionToProvider), nil
}
