package identifyCloudActors

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/dragondrop-cloud/driftmitigation/hclcreate"
	queryParamData "github.com/dragondrop-cloud/driftmitigation/implementations/identify_cloud_actors/query_param_data"
	resourcesCalculator "github.com/dragondrop-cloud/driftmitigation/implementations/resources_calculator"
	driftDetector "github.com/dragondrop-cloud/driftmitigation/implementations/terraform_managed_resources_drift_detector/drift_detector"
	terraformValueObjects "github.com/dragondrop-cloud/driftmitigation/implementations/terraform_value_objects"
)

// AWSLogQuerier implements the LogQuerier interface for AWS.
type AWSLogQuerier struct {

	// divisionToCredentials is a map between a division and request cloud credentials.
	divisionToCredentials terraformValueObjects.DivisionCloudCredentialDecoder `required:"true"`

	// divisionToNewResources is a map between a division and a list of new resource objects.
	divisionToNewResources resourcesCalculator.DivisionToNewResources

	// divisionToUniqueManagedDriftedResources is a map between a division and a list of unique drifted resource objects.
	divisionToUniqueManagedDriftedResources DivisionToUniqueDriftedResources

	// httpClient is a http client shared across all http requests within this package.
	httpClient http.Client

	// managedDriftAttributeDifferences is a list of all attribute differences.
	managedDriftAttributeDifferences []driftDetector.AttributeDifference

	// resourceToCloudTrailType is a map between a Terraform resource type and the corresponding Cloud Trail event type.
	resourceToCloudTrailType queryParamData.AWSResourceToCloudTrailResource
}

// AWSEnvironment is a struct defining the credential values needed for authenticating with an AWS account.
type AWSEnvironment struct {
	AWSAccessKeyID     string `json:"awsAccessKeyID"`
	AWSSecretKeyAccess string `json:"awsSecretAccessKey"`
}

// CloudTrailEvents is a struct containing all the data returned from the AWS CLI command
// `aws cloudtrail lookup-events`.
type CloudTrailEvents struct {
	Events []CloudTrailEvent `json:"Events"`
}

// CloudTrailEvent is a struct for a single event within return data from the AWS CLI
// command `aws cloudtrail lookup-events`.
type CloudTrailEvent struct {
	EventID              string  `json:"EventId"`
	EventName            string  `json:"EventName"`
	EventTimeUnformatted float64 `json:"EventTime"`
	EventTime            string
	UserName             string               `json:"Username"`
	Resources            []CloudTrailResource `json:"Resources"`
}

// CloudTrailResource is a struct for a resource identity within a CloudTrailEvent.
type CloudTrailResource struct {
	ResourceType string `json:"ResourceType"`
	ResourceName string `json:"ResourceName"`
}

// NewAWSLogQuerier instantiates a new instance of GoogleLogQuerier
func NewAWSLogQuerier(
	divisionToCredentials terraformValueObjects.DivisionCloudCredentialDecoder,
) (LogQuerier, error) {
	return &AWSLogQuerier{
		divisionToCredentials:    divisionToCredentials,
		resourceToCloudTrailType: queryParamData.NewAWSResourceToCloudTrailLookup(),
	}, nil
}

// loadUpstreamDataToAWSLogQuerier loads all data needed for querying logs from upstream
// saved data sources.
func (alc *AWSLogQuerier) loadUpstreamDataToAWSLogQuerier() error {
	attributeDifferences, err := loadDriftResourcesDifferences()
	if err != nil {
		return fmt.Errorf("[loadDriftResourcesDifferences]%v", err)
	}

	divToNewResources, err := loadDivisionToNewResources()
	if err != nil {
		return fmt.Errorf("[loadDivisionToNewResources]%v", err)
	}
	divToUniqueDriftedResources, err := createDivisionUniqueDriftedResources(attributeDifferences)
	if err != nil {
		return fmt.Errorf("[createDivisionUniqueDriftedResources]%v", err)
	}

	alc.divisionToUniqueManagedDriftedResources = divToUniqueDriftedResources
	alc.divisionToNewResources = divToNewResources
	alc.managedDriftAttributeDifferences = attributeDifferences
	return nil
}

// QueryForAllResources coordinates calls of QueryForResourcesInDivision for all
// divisions from which drifted resources have been identified.
func (alc *AWSLogQuerier) QueryForAllResources(ctx context.Context) (terraformValueObjects.DivisionResourceActions, error) {
	divisionToResourceActions := terraformValueObjects.DivisionResourceActions{}

	err := alc.loadUpstreamDataToAWSLogQuerier()
	if err != nil {
		return divisionToResourceActions, fmt.Errorf("[alc.loadUpstreamDataToAWSLogQuerier]%v", err)
	}

	for division := range alc.divisionToCredentials {
		fmt.Printf("Pulling cloud actor actions for the AWS account represented by %v\n", division)
		divisionAllResourceActions, err := alc.QueryForResourcesInDivision(ctx, division)
		if err != nil {
			return divisionToResourceActions, fmt.Errorf("[alc.QueryForResourcesInDivision]%v", err)
		}
		divisionToResourceActions[division] = divisionAllResourceActions
	}
	return divisionToResourceActions, nil
}

// QueryForResourcesInDivision coordinates calls of cloudTrailEventHistorySearch for all
// resources within a division - both managed drift and resources outside of Terraform control.
func (alc *AWSLogQuerier) QueryForResourcesInDivision(ctx context.Context, division terraformValueObjects.Division) (map[terraformValueObjects.ResourceName]terraformValueObjects.ResourceActions, error) {
	divisionResourceActions := map[terraformValueObjects.ResourceName]terraformValueObjects.ResourceActions{}
	credential := alc.divisionToCredentials[division]
	err := alc.setAWSCredentials(credential)
	if err != nil {
		return nil, fmt.Errorf("[alc.setAWSCredentials]%v", err)
	}

	division = terraformValueObjects.Division("aws-" + string(division))

	// Calculating cloud actors for managed resource drift
	currentUniqueDriftedResources, ok := alc.divisionToUniqueManagedDriftedResources[division]
	if ok {
		for _, driftedResource := range currentUniqueDriftedResources {
			resourceActions, err := alc.cloudTrailEventHistorySearch(ctx, driftedResource.ResourceType, driftedResource.InstanceID, driftedResource.Region, false)
			if err != nil {
				return divisionResourceActions, fmt.Errorf("[alc.cloudTrailEventHistory]%v", err)
			}

			currentResourceName := uniqueDriftedResourceToName(driftedResource)
			divisionResourceActions[currentResourceName] = resourceActions
		}

		alc.UpdateManagedDriftAttributeDifferences(divisionResourceActions)

		// Overwrite the drift-resources-differences.json file with the new data.
		managedAttributeDifferencesBytes, err := json.MarshalIndent(alc.managedDriftAttributeDifferences, "", "  ")
		if err != nil {
			return divisionResourceActions, fmt.Errorf("[json.MarshalIndent]%v", err)
		}

		err = os.WriteFile("mappings/drift-resources-differences.json", managedAttributeDifferencesBytes, 0400)
		if err != nil {
			return divisionResourceActions, fmt.Errorf("[os.WriteFile]%v", err)
		}
	}

	// Calculating cloud actors for new resource drift
	currentNewResources, ok := alc.divisionToNewResources[division]
	if ok {
		for id, resource := range currentNewResources {
			resourceActions, err := alc.cloudTrailEventHistorySearch(ctx, resource.ResourceType, string(id), resource.Region, true)
			if err != nil {
				return divisionResourceActions, fmt.Errorf("[alc.cloudTrailEventHistory]%v", err)
			}

			currentResourceName := terraformValueObjects.ResourceName(
				resource.ResourceType + "." + hclcreate.ConvertTerraformerResourceName(resource.ResourceTerraformerName),
			)
			divisionResourceActions[currentResourceName] = resourceActions
		}
	}

	return divisionResourceActions, nil
}

// UpdateManagedDriftAttributeDifferences updates the RecentActor and RecentActionTimestamp fields
// for each struct within the alc.managedDriftAttributeDifferences slice.
func (alc *AWSLogQuerier) UpdateManagedDriftAttributeDifferences(
	divisionResourceActions map[terraformValueObjects.ResourceName]terraformValueObjects.ResourceActions,
) {
	newAttributeDifferences := []driftDetector.AttributeDifference{}

	for _, attributeDifference := range alc.managedDriftAttributeDifferences {
		currentDifferenceResourceName := attributeDifferenceToResourceName(attributeDifference)

		if _, ok := divisionResourceActions[currentDifferenceResourceName]; ok {
			resourceAction := divisionResourceActions[currentDifferenceResourceName]

			attributeDifference.RecentActor = resourceAction.Modifier.Actor
			attributeDifference.RecentActionTimestamp = resourceAction.Modifier.Timestamp
		}
		newAttributeDifferences = append(newAttributeDifferences, attributeDifference)
	}

	alc.managedDriftAttributeDifferences = newAttributeDifferences
}

// TODO: At some point results may be further back in time, a subsequent improvement would be to look further back if resource details are not found.
// cloudTrailEventHistorySearch runs AWS CLI commands to pull data on who modified and created the cloud resource in question.
func (alc *AWSLogQuerier) cloudTrailEventHistorySearch(ctx context.Context, resourceType string, resourceID string, resourceRegion string, isNewToTerraform bool) (terraformValueObjects.ResourceActions, error) {
	lookupAttributeString := fmt.Sprintf("AttributeKey=ResourceName,AttributeValue=%v", resourceID)
	cloudTrailCommand := []string{"cloudtrail", "lookup-events", "--max-results", "50", "--output", "json", "--region", resourceRegion, "--lookup-attributes", lookupAttributeString}

	result, err := executeCommandReturnStdOut("aws", cloudTrailCommand...)
	if err != nil {
		return terraformValueObjects.ResourceActions{}, fmt.Errorf("[executeCommandReturnStdOut]%v", err)
	}

	return alc.ExtractDataFromResourceResult([]byte(result), resourceType, isNewToTerraform)
}

// ExtractDataFromResourceResult parses the log response from the provider API
// and extracts needed data (namely who made the most recent relevant change to the resource).
func (alc *AWSLogQuerier) ExtractDataFromResourceResult(resourceResult []byte, resourceType string, isNewToTerraform bool) (terraformValueObjects.ResourceActions, error) {
	resourceActions := terraformValueObjects.ResourceActions{}

	var cloudTrailEvents CloudTrailEvents
	if err := json.Unmarshal(resourceResult, &cloudTrailEvents); err != nil {
		return resourceActions, fmt.Errorf("failed to parse resource results to cloudTrailEvents struct: %v", err)
	}
	resourceType = string(alc.resourceToCloudTrailType[resourceType])

	isComplete := false
	isModificationIdentified := false
	i := 0

	for !isComplete {
		event := cloudTrailEvents.Events[i]
		classification := determineActionClass(event.EventName)
		event.EventTime = decimalToFormattedTimestamp(event.EventTimeUnformatted)

		// check to ensure that ResourceType is present within one of the event.Resources elements
		// if not, move on to the next event
		isValidResourceType := false
		for _, resource := range event.Resources {
			if resource.ResourceType == resourceType {
				isValidResourceType = true
				break
			}
		}
		if !isValidResourceType {
			i++
			continue
		}

		switch classification {
		case "creation":
			if isNewToTerraform {
				resourceActions.Creator = terraformValueObjects.CloudActorTimeStamp{
					Actor:     terraformValueObjects.CloudActor(event.UserName),
					Timestamp: terraformValueObjects.Timestamp(event.EventTime),
				}
				return resourceActions, nil
			}
		case "modification":
			if !isModificationIdentified {
				isModificationIdentified = true
				resourceActions.Modifier = terraformValueObjects.CloudActorTimeStamp{
					Actor:     terraformValueObjects.CloudActor(event.UserName),
					Timestamp: terraformValueObjects.Timestamp(event.EventTime),
				}
				if !isNewToTerraform {
					return resourceActions, nil
				}
			}
		}
		i++
	}

	return resourceActions, nil
}

// decimalToFormattedTimestamp is a function to convert decimal to a unix timestamp
func decimalToFormattedTimestamp(decimal float64) string {
	sec, dec := math.Modf(decimal)
	return time.Unix(int64(sec), int64(dec*(1e9))).Format("2006-01-02")
}

// executeCommandReturnStdOut wraps os.exec.Command with capturing of std output and errors.
// It also returns the command results.
func executeCommandReturnStdOut(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)

	// Setting up logging objects
	var out bytes.Buffer
	cmd.Stdout = &out

	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil {
		return "", fmt.Errorf("%v\n\n%v", err, stderr.String()+out.String())
	}
	return out.String(), nil
}

// setAWSCredentials loads and sets as environment variables AWS credentials for a given AWS account.
func (alc *AWSLogQuerier) setAWSCredentials(credential terraformValueObjects.Credential) error {
	envVars := new(AWSEnvironment)
	err := json.Unmarshal([]byte(credential), &envVars)
	if err != nil {
		return fmt.Errorf("[json.Unmarshal] %w", err)
	}

	err = os.Setenv("AWS_ACCESS_KEY_ID", envVars.AWSAccessKeyID)
	if err != nil {
		return fmt.Errorf("[os.Setenv]%w", err)
	}

	err = os.Setenv("AWS_SECRET_ACCESS_KEY", envVars.AWSSecretKeyAccess)
	if err != nil {
		return fmt.Errorf("[os.Setenv]%w", err)
	}

	return nil
}
