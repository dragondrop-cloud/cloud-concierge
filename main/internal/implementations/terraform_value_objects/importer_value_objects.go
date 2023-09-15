package terraformvalueobjects

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

// ResourceName is the name of a cloud computing resource
type ResourceName string

// ResourceNameList is a list of ResourceNames
type ResourceNameList []ResourceName

// Decode allows ResourceNameList to be decoded by the ENVConfig function.
func (rnl *ResourceNameList) Decode(value string) error {
	logrus.Debugf("Decoding terraform resource name list %v", value)

	if value == "None" {
		*rnl = nil
		return nil
	}

	newValue := strings.TrimSuffix(value, "]")
	newValue = strings.TrimPrefix(newValue, "[")
	if newValue == value {
		return fmt.Errorf("expected the values formatted in a list bookended by '[' and ']' chars")
	}

	stringSlice := strings.Split(newValue, ",")

	resourceNameSlice := ResourceNameList{}
	for _, val := range stringSlice {
		val = strings.Trim(val, `"`)
		resourceNameSlice = append(resourceNameSlice, ResourceName(val))
	}
	*rnl = resourceNameSlice

	return nil
}

// RemoteCloudReference is the identifying string where a resource is located within a remote cloud
// environment for use within a `terraform import` statement
type RemoteCloudReference string

// TerraformConfigLocation is the location where a resource can be identified within Terraform by the
// {TerraformResourceType}.{TerraformResourceName} syntax.
type TerraformConfigLocation string

// ImportMigration are the full args for terraform import statement
type ImportMigration struct {
	TerraformConfigLocation TerraformConfigLocation
	RemoteCloudReference    RemoteCloudReference
}

// ResourceImportMap is the import migration statement for the resource name
type ResourceImportMap map[ResourceName]ImportMigration

// ProviderToResourceImportMap is a map of structure {Provider: {Division: ResourceImportMap}}
type ProviderToResourceImportMap map[Provider]map[Division]ResourceImportMap

// Account is a string that represents an AWS cloud account.
type Account string
