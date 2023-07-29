package identifyCloudActors

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

	// QueryForAllResources coordinates calls of QueryForResourcesInDivision for all
	// divisions from which drifted resources have been identified.
	QueryForAllResources(ctx context.Context) (terraformValueObjects.DivisionResourceActions, error)

	// QueryForResourcesInDivision coordinates calls of QuerySingleResource for all
	// resources within a division.
	QueryForResourcesInDivision(ctx context.Context, division terraformValueObjects.Division) (map[terraformValueObjects.ResourceName]terraformValueObjects.ResourceActions, error)

	// ExtractDataFromResourceResult parses the log response from the provider API
	// and extracts needed data (namely who made the most recent relevant change to the resource).
	ExtractDataFromResourceResult(resourceResult []byte, resourceType string, isNewToTerraform bool) (terraformValueObjects.ResourceActions, error)
}

type DivisionToUniqueDriftedResources map[terraformValueObjects.Division]map[terraformValueObjects.ResourceName]UniqueDriftedResource

// UniqueDriftedResource is a type that represents a cloud resource that has drifted from its expected state.
// without any information on individual attributes that have drifted.
type UniqueDriftedResource struct {
	CloudDivision string
	InstanceID    string
	Region        string
	ResourceType  string
	ResourceName  string
	StateFileName driftDetector.StateFileName
}

// NewProviderToLogQuerierMap returns a map between cloud providers and an instantiated LogQuerier
// implementation for that provider.
func NewProviderToLogQuerierMap(globalConfig Config, provider terraformValueObjects.Provider) (map[terraformValueObjects.Provider]LogQuerier, error) {
	providerToQuerier := map[terraformValueObjects.Provider]LogQuerier{}

	gcpDivCredentials := filterDivisionCloudCredentialsForProvider("google", divisionToProvider, globalConfig)
	if len(gcpDivCredentials) > 0 {
		googleLogQuerier, err := NewGoogleLogQuerier(gcpDivCredentials)
		if err != nil {
			return nil, fmt.Errorf("[NewGoogleLogQuerier]%v", err)
		}
		providerToQuerier["google"] = googleLogQuerier
	}

	awsDivCredentials := filterDivisionCloudCredentialsForProvider("aws", divisionToProvider, globalConfig)
	if len(awsDivCredentials) > 0 {
		awsLogQuerier, err := NewAWSLogQuerier(awsDivCredentials)
		if err != nil {
			return nil, fmt.Errorf("[NewAWSLogQuerier]%v", err)
		}
		providerToQuerier["aws"] = awsLogQuerier
	}

	return providerToQuerier, nil
}

// filterDivisionCloudCredentialsForProvider is the subset of cloud credentials for a particular cloud provider.
func filterDivisionCloudCredentialsForProvider(providerName terraformValueObjects.Provider, provider terraformValueObjects.Provider, globalConfig Config) terraformValueObjects.DivisionCloudCredentialDecoder {
	filteredDivToCredential := terraformValueObjects.DivisionCloudCredentialDecoder{}

	for division, provider := range divisionToProvider {
		if provider == providerName {
			filteredDivToCredential[division] = globalConfig.DivisionCloudCredentials[division]
		}
	}
	return filteredDivToCredential
}

// determineActionClass determines the classification of an input method, which is either a resource
// "modification", "creation", "deletion", or "not_classified".
func determineActionClass(value string) string {
	value = strings.ToLower(value)
	if strings.Contains(value, "update") || strings.Contains(value, "replace") || strings.Contains(value, "modify") {
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
	if _, err := os.Stat("mappings/drift-resources-differences.json"); errors.Is(err, os.ErrNotExist) {
		return resourceDifferences, nil
	}

	fileContent, err := os.ReadFile("mappings/drift-resources-differences.json")
	if err != nil {
		return nil, fmt.Errorf("[os.ReadFile]%v", err)
	}

	err = json.Unmarshal(fileContent, &resourceDifferences)
	if err != nil {
		return nil, fmt.Errorf("[json.Unmarshal]%v", err)
	}

	return resourceDifferences, nil
}

// loadDivisionToNewResources loads the division-to-new-resources file as a
// resourcesCalculator.DivisionToNewResources struct.
func loadDivisionToNewResources() (resourcesCalculator.DivisionToNewResources, error) {
	newResources := resourcesCalculator.DivisionToNewResources{}
	if _, err := os.Stat("mappings/division-to-new-resources.json"); errors.Is(err, os.ErrNotExist) {
		return newResources, nil
	}

	fileContent, err := os.ReadFile("mappings/division-to-new-resources.json")
	if err != nil {
		return newResources, fmt.Errorf("[os.ReadFile]%v", err)
	}

	err = json.Unmarshal(fileContent, &newResources)
	if err != nil {
		return newResources, fmt.Errorf("[json.Unmarshal]%v", err)
	}
	return newResources, nil
}

// createDivisionUniqueDriftedResources converts a slice of AttributeDifference into
// DivisionToUniqueDriftedResources.
func createDivisionUniqueDriftedResources(differences []driftDetector.AttributeDifference) (DivisionToUniqueDriftedResources, error) {
	output := DivisionToUniqueDriftedResources{}

	for _, dif := range differences {
		division := terraformValueObjects.Division(dif.CloudDivision)

		currentDivision := map[terraformValueObjects.ResourceName]UniqueDriftedResource{}
		if _, ok := output[division]; ok {
			currentDivision = output[division]
		}

		// create a unique resource name
		uniqueResourceName := terraformValueObjects.ResourceName(
			strings.Join([]string{string(dif.StateFileName), dif.ResourceType, dif.ResourceName, dif.InstanceID}, "."),
		)
		if _, ok := currentDivision[uniqueResourceName]; !ok {
			currentDivision[uniqueResourceName] = UniqueDriftedResource{
				CloudDivision: dif.CloudDivision,
				InstanceID:    dif.InstanceID,
				Region:        dif.InstanceRegion,
				ResourceType:  dif.ResourceType,
				ResourceName:  dif.ResourceName,
				StateFileName: dif.StateFileName,
			}
			output[division] = currentDivision
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
