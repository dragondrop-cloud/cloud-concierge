package identifycloudactors

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	resourcesCalculator "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/resources_calculator"
	driftDetector "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_managed_resources_drift_detector/drift_detector"
	terraformValueObjects "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_value_objects"
)

// LogQuerier is an interface for querying information from a single cloud providers
// administrative logs.
type LogQuerier interface {
	// QueryForAllResources coordinates API calls that receive data on user actions on individual resources.
	QueryForAllResources(ctx context.Context) (terraformValueObjects.ResourceActionMap, error)
}

type UniqueDriftedResources map[terraformValueObjects.ResourceName]UniqueDriftedResource

// UniqueDriftedResource is a type that represents a cloud resource that has drifted from its expected state.
// without any information on individual attributes that have drifted.
type UniqueDriftedResource struct {
	InstanceID    string
	Region        string
	ResourceType  string
	ResourceName  string
	StateFileName driftDetector.StateFileName
}

// NewLogQuerier returns an instantiated LogQuerier implementation for the specified provider.
func NewLogQuerier(globalConfig Config, provider terraformValueObjects.Provider) (LogQuerier, error) {
	switch provider {
	case "google":
		googleLogQuerier, err := NewGoogleLogQuerier(globalConfig)
		if err != nil {
			return nil, fmt.Errorf("[NewGoogleLogQuerier]%v", err)
		}
		return googleLogQuerier, nil
	case "aws":
		awsLogQuerier, err := NewAWSLogQuerier(globalConfig)
		if err != nil {
			return nil, fmt.Errorf("[NewAWSLogQuerier]%v", err)
		}
		return awsLogQuerier, nil
	default:
		fmt.Printf("provider %s not supported for log querying", provider)
		return nil, nil
	}
}

// determineActionClass determines the classification of an input method, which is either a resource
// "modification", "creation", "deletion", or "not_classified".
func determineActionClass(value string) string {
	value = strings.ToLower(value)
	if strings.Contains(value, "update") || strings.Contains(value, "replace") || strings.Contains(value, "modify") || strings.Contains(value, "createtags") {
		return "modification"
	}

	if strings.Contains(value, "delete") || strings.Contains(value, "deletion") {
		return "deletion"
	}

	if strings.Contains(value, "create") {
		return "creation"
	}

	return "not_classified"
}

// loadDriftResourcesDifferences loads the drift-resources-differences file as a slice
// of driftDetector.AttributeDifference.
func loadDriftResourcesDifferences() ([]driftDetector.AttributeDifference, error) {
	var resourceDifferences []driftDetector.AttributeDifference
	if _, err := os.Stat("outputs/drift-resources-differences.json"); errors.Is(err, os.ErrNotExist) {
		return resourceDifferences, nil
	}

	fileContent, err := os.ReadFile("outputs/drift-resources-differences.json")
	if err != nil {
		return nil, fmt.Errorf("[os.ReadFile]%v", err)
	}

	err = json.Unmarshal(fileContent, &resourceDifferences)
	if err != nil {
		return nil, fmt.Errorf("[json.Unmarshal]%v", err)
	}

	return resourceDifferences, nil
}

// loadNewResources loads the new-resources file as a resourcesCalculator.NewResourceMap struct.
func loadNewResources() (resourcesCalculator.NewResourceMap, error) {
	newResources := resourcesCalculator.NewResourceMap{}
	if _, err := os.Stat("outputs/new-resources.json"); errors.Is(err, os.ErrNotExist) {
		return newResources, nil
	}

	fileContent, err := os.ReadFile("outputs/new-resources.json")
	if err != nil {
		return newResources, fmt.Errorf("[os.ReadFile]%v", err)
	}

	err = json.Unmarshal(fileContent, &newResources)
	if err != nil {
		return newResources, fmt.Errorf("[json.Unmarshal]%v", err)
	}
	return newResources, nil
}

// createUniqueDriftedResources converts a slice of AttributeDifference into UniqueDriftedResources.
func createUniqueDriftedResources(differences []driftDetector.AttributeDifference) (UniqueDriftedResources, error) {
	output := UniqueDriftedResources{}

	for _, dif := range differences {
		// create a unique resource name
		uniqueResourceName := terraformValueObjects.ResourceName(
			strings.Join([]string{string(dif.StateFileName), dif.ResourceType, dif.ResourceName, dif.InstanceID}, "."),
		)
		if _, ok := output[uniqueResourceName]; !ok {
			output[uniqueResourceName] = UniqueDriftedResource{
				InstanceID:    dif.InstanceID,
				Region:        dif.InstanceRegion,
				ResourceType:  dif.ResourceType,
				ResourceName:  dif.ResourceName,
				StateFileName: dif.StateFileName,
			}
		}
	}

	return output, nil
}

// uniqueDriftedResourceToName converts a UniqueDriftedResource into a unique name
func uniqueDriftedResourceToName(resource UniqueDriftedResource) terraformValueObjects.ResourceName {
	return terraformValueObjects.ResourceName(
		strings.Join(
			[]string{string(resource.StateFileName), resource.ResourceType, resource.ResourceName, resource.InstanceID},
			".",
		),
	)
}

// attributeDifferenceToResourceName converts a attributeDifferenceToResourceName into a unique name
func attributeDifferenceToResourceName(resource driftDetector.AttributeDifference) terraformValueObjects.ResourceName {
	return terraformValueObjects.ResourceName(
		strings.Join(
			[]string{string(resource.StateFileName), resource.ResourceType, resource.ResourceName, resource.InstanceID},
			".",
		),
	)
}
