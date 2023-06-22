package terraformerCLI

import (
	terraformValueObjects "github.com/dragondrop-cloud/driftmitigation/implementations/terraform_value_objects"
)

// MultiScanResult maps divisions
type MultiScanResult struct {
	stacks map[terraformValueObjects.Division]terraformValueObjects.Path
}

// Scanner interface allows scanning a single division within a cloud environment at a time.
type Scanner interface {

	// Scan uses the TerraformerCLI interface to scan a given division's cloud environment. Returns
	// the name of the Division scanned, and the Stack of Terraformer output for that division.
	Scan(division terraformValueObjects.Division, credential terraformValueObjects.Credential, options ...string) (terraformValueObjects.Path, error)

	// ScanAll wraps Scan to scan each division for the provider.
	ScanAll(options ...string) (*MultiScanResult, error)
}
