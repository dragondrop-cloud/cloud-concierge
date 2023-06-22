package driftDetector

import (
	"fmt"
	"strings"
)

// ResourceIDCalculator determines the unique ID of a terraform resource based on attributes,
// cloud provider and resource type.
func ResourceIDCalculator(attributesFlat map[string]string, cloudProvider string, resourceType string) (string, error) {
	switch cloudProvider {
	case "google":
		selfLink, ok := attributesFlat["self_link"]
		if ok {
			outputString := parseIDFromField(selfLink)
			if outputString != "" {
				return outputString, nil
			}
		}

		id := attributesFlat["id"]
		outputString := parseIDFromField(id)
		if outputString != "" {
			return outputString, nil
		}

		// At this point there is no clear id, so we construct one from custom data loaded here
		googleToTFQueryName := NewGoogleTfToQueryName()
		if stringToFill, ok := googleToTFQueryName[TerraformResourceType(resourceType)]; ok {
			return fmt.Sprintf(string(stringToFill), attributesFlat["id"]), nil
		}
		return attributesFlat["id"], nil
	default:
		return attributesFlat["id"], nil
	}
}

// parseIDFromField extracts the google cloud-resource identifying string from a field. This
// string starts with either "projects" or "namespaces" and includes at least on "/".
func parseIDFromField(field string) string {
	startingIndex := -1
	if (strings.Contains(field, "projects") || strings.Contains(field, "namespaces")) && strings.Contains(field, "/") {
		startingIndex = strings.Index(field, "projects")
		if startingIndex == -1 {
			startingIndex = strings.Index(field, "namespaces")
		}
	}
	if startingIndex == -1 {
		return ""
	}

	return field[startingIndex:]
}

// TerraformResourceType is a Terraform resource type represented by a string.
type TerraformResourceType string

// LoggingResourceName is the value which can be used as a ResourceName filter in a GCP log query.
type LoggingResourceName string

// TerraformResourceToLoggingName is a map between a Terraform resource type and the LoggingResourceName
// for that resource type.
type TerraformResourceToLoggingName map[TerraformResourceType]LoggingResourceName

// NewGoogleTfToQueryName returns a new instance of TerraformResourceToGoogleQueryParams.
func NewGoogleTfToQueryName() TerraformResourceToLoggingName {
	return TerraformResourceToLoggingName{
		"google_storage_bucket": "projects/_/buckets/%v",
	}
}
