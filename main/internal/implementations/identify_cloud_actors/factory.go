package identifyCloudActors

import (
	"context"

	terraformValueObjects "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_value_objects"
	"github.com/dragondrop-cloud/cloud-concierge/main/internal/interfaces"
)

// Factory is a struct that generates implementations of interfaces.IdentifyCloudActors.
type Factory struct {
}

// Instantiate returns an implementation of interfaces.IdentifyCloudActors depending on the passed
// environment specification.
func (f *Factory) Instantiate(ctx context.Context, environment string, dragonDrop interfaces.DragonDrop, provider terraformValueObjects.Provider, config Config) (interfaces.IdentifyCloudActors, error) {
	switch environment {
	case "isolated":
		return new(IsolatedIdentifyCloudActors), nil
	default:
		return f.bootstrappedResourceCalculator(dragonDrop, divisionToProvider, config)
	}
}

// bootstrappedResourceCalculator creates a complete implementation of the interfaces.IdentifyCloudActors interface with
// configuration specified via environment variables.
func (f *Factory) bootstrappedResourceCalculator(dragonDrop interfaces.DragonDrop, provider terraformValueObjects.Provider, config Config) (interfaces.IdentifyCloudActors, error) {
	return NewIdentifyCloudActors(config, dragonDrop, divisionToProvider)
}
