package terraformValueObjects

import (
	"fmt"
	"regexp"
	"strings"
)

// Credential is a credential which can be used to read resources within a cloud footprint.
// The string is a json structure in json format.
type Credential string

// Division is the name of a division within a cloud provider. For AWS a region, for Azure a resource group, and for GCP
// this is a project name.
type Division string

// DivisionCloudCredentialDecoder is a type that implements the envconfig.Decoder interface and used
// for decoding the input into a map of strings.
type DivisionCloudCredentialDecoder map[Division]Credential

// Decode provides the object decoding logic for DivisionCloudCredentialDecoder, in accordance with the envconfig
// package's requirements.
func (dcd *DivisionCloudCredentialDecoder) Decode(value string) error {
	divToCredential := map[Division]Credential{}

	// First pulling out each individual division-credential pair
	r, _ := regexp.Compile(`(.*?:{[\S\s]*?}),?$?`)
	allGroups := r.FindAllString(value, -1)

	if len(allGroups) == 0 {
		return fmt.Errorf("no cloud division credentials were found in the expected `division:{}` format")
	}
	//// For each pairing, pull out the respective division name and credential json-string
	rGroup, _ := regexp.Compile(`(.*?):(\{[\S\s]*\})`)
	for _, group := range allGroups {
		division := Division(rGroup.FindStringSubmatch(group)[1])
		credential := Credential(rGroup.FindStringSubmatch(group)[2])
		divToCredential[division] = credential
	}

	*dcd = DivisionCloudCredentialDecoder(divToCredential)
	return nil
}

// Path is the relative file path within the 'current_cloud' directory to the division's output content.
type Path string

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
