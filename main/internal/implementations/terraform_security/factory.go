package terraformSecurity

import (
	"context"

	terraformValueObjects "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_value_objects"
	"github.com/dragondrop-cloud/cloud-concierge/main/internal/interfaces"
)

// Factory is a struct that generates implementations of interfaces.TerraformSecurity.
type Factory struct {
}

// Instantiate returns an implementation of interfaces.TerraformSecurity depending on the passed
// environment specification.
func (f *Factory) Instantiate(ctx context.Context, environment string, divisionToProvider map[terraformValueObjects.Division]terraformValueObjects.Provider) (interfaces.TerraformSecurity, error) {
	switch environment {
	case "isolated":
		return NewIsolatedTerraformSecurity(), nil
	default:
		return f.bootstrappedTerraformSecurity(divisionToProvider)
	}
}

// bootstrappedTerraformSecurity creates a complete implementation of the interfaces.TerraformSecurity interface with
// configuration specified via environment variables.
func (f *Factory) bootstrappedTerraformSecurity(divisionToProvider map[terraformValueObjects.Division]terraformValueObjects.Provider) (interfaces.TerraformSecurity, error) {
	return NewTFSec(divisionToProvider), nil
}
