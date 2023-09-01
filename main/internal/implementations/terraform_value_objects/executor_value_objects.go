package terraformvalueobjects

import (
	"fmt"
	"strings"
)

// Credential is a credential which can be used to read resources within a cloud footprint.
// The string is a json structure in json format.
type Credential string

// Division is the name of a division within a cloud provider. For AWS a region, for Azure a resource group, and for GCP
// this is a project name.
type Division string

// Provider is the name of a cloud computing resource provider.
type Provider string

// Version is a Terraform module version string.
type Version string

// Decode provides the object decoding logic for Version, in accordance with the envconfig
// package's requirements.
func (v *Version) Decode(value string) error {
	if string(value[1]) != "." {
		return fmt.Errorf("terraform version should start with 'major version[.]'")
	}

	stringComponents := strings.Split(value, ".")
	versionLength := len(stringComponents)
	if versionLength != 3 {
		return fmt.Errorf("expected three pieces of the version once split by '.', instead got %v", versionLength)
	}

	majorVersion := stringComponents[0]

	if (majorVersion != "0") && (majorVersion != "1") {
		return fmt.Errorf("terraform major version must be either '0' or '1', got %v", majorVersion)
	}

	*v = Version(value)
	return nil
}
