package resourcescalculator

import (
	"context"

	"github.com/dragondrop-cloud/cloud-concierge/main/internal/documentize"
	terraformValueObjects "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_value_objects"
	"github.com/dragondrop-cloud/cloud-concierge/main/internal/interfaces"
)

// Factory is a struct that generates implementations of interfaces.ResourcesCalculator.
type Factory struct{}

// Instantiate returns an implementation of interfaces.ResourcesCalculator depending on the passed
// environment specification.
func (f *Factory) Instantiate(
	ctx context.Context, environment string,
	provider terraformValueObjects.Provider,
	nlpEngine interfaces.NLPEngine,
) (interfaces.ResourcesCalculator, error) {
	switch environment {
	case "isolated":
		return new(IsolatedResourcesCalculator), nil
	default:
		return f.bootstrappedResourceCalculator(ctx, provider, nlpEngine)
	}
}

// bootstrappedResourceCalculator creates a complete implementation of the interfaces.ResourcesCalculator interface with
// configuration specified via environment variables.
func (f *Factory) bootstrappedResourceCalculator(
	ctx context.Context,
	provider terraformValueObjects.Provider,
	nlpEngine interfaces.NLPEngine,
) (interfaces.ResourcesCalculator, error) {
	doc, _ := documentize.NewDocumentize(provider)

	return NewTerraformResourcesCalculator(&doc, nlpEngine), nil
}
