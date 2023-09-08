package terraformimportmigrationgenerator

import (
	"fmt"
	"strings"

	"github.com/Jeffail/gabs/v2"

	terraformValueObjects "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_value_objects"
)

// ResourceType is the type of the terraform cloud resource
type ResourceType string

// ImportLocationFormat is the format which is applied to the import statement
type ImportLocationFormat struct {
	StringFormat string
	Attributes   []string
}

// providerResourceLocationFormats are all the format rules to apply to the terraform import migration statement
var providerResourceLocationFormats = map[terraformValueObjects.Provider]map[ResourceType]ImportLocationFormat{
	"aws":    ResourceTypeLocations,
	"google": GoogleResourceTypeLocations,
}

// GetRemoteCloudReference extracts the formatted string from the resources json
func GetRemoteCloudReference(resource *gabs.Container, provider terraformValueObjects.Provider, resourceType ResourceType) (string, error) {
	format := providerResourceLocationFormats[provider][resourceType]
	formattedString := format.StringFormat

	for i, attribute := range format.Attributes {
		value := resource.Path(fmt.Sprintf("instances.0.attributes_flat.%s", attribute)).Data().(string)
		formattedString = strings.Replace(formattedString, fmt.Sprintf("$%d", i), value, -1)
	}

	return formattedString, nil
}
