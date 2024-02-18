package terraformercli

import (
	terraformValueObjects "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_value_objects"
)

// Scanner interface allows scanning a single division within a cloud environment at a time.
type Scanner interface {
	// Scan uses the TerraformerCLI interface to scan a given division's cloud environment
	Scan(division terraformValueObjects.Division, credential terraformValueObjects.Credential, options ...string) error
}
