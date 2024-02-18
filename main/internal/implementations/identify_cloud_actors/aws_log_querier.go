package identifycloudactors

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudtrail"
	"github.com/dragondrop-cloud/cloud-concierge/main/internal/hclcreate"
	queryParamData "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/identify_cloud_actors/query_param_data"
	resourcesCalculator "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/resources_calculator"
	driftDetector "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_managed_resources_drift_detector/drift_detector"
	terraformValueObjects "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_value_objects"
	log "github.com/sirupsen/logrus"
)

var ErrNoCloudTrailEvents = errors.New("no events found")

// AWSLogQuerier implements the LogQuerier interface for AWS.
type AWSLogQuerier struct {
	// cloudCredential is a map between a division and request cloud credentials.
	cloudCredential terraformValueObjects.Credential `required:"true"`

	// awsSession from the AWS Go SDK.
	awsSession *session.Session

	// division is the division that the AWSLogQuerier is querying for.
	division terraformValueObjects.Division

	// newResources is a map between a division and a list of new resource objects.
	newResources resourcesCalculator.NewResourceMap

	// uniqueManagedDriftedResources is a map between a division and a list of unique drifted resource objects.
	uniqueManagedDriftedResources UniqueDriftedResources

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

// NewAWSLogQuerier instantiates a new instance of GoogleLogQuerier
func NewAWSLogQuerier(
	config Config,
) (LogQuerier, error) {
	return &AWSLogQuerier{
		cloudCredential:          config.CloudCredential,
		division:                 config.Division,
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

	divToNewResources, err := loadNewResources()
	if err != nil {
		return fmt.Errorf("[loadDivisionToNewResources]%v", err)
	}
	divToUniqueDriftedResources, err := createUniqueDriftedResources(attributeDifferences)
	if err != nil {
		return fmt.Errorf("[createDivisionUniqueDriftedResources]%v", err)
	}

	alc.uniqueManagedDriftedResources = divToUniqueDriftedResources
	alc.newResources = divToNewResources
	alc.managedDriftAttributeDifferences = attributeDifferences
	return nil
}

// QueryForAllResources coordinates calls of QueryForResourcesInDivision for all
// divisions from which drifted resources have been identified.
func (alc *AWSLogQuerier) QueryForAllResources(ctx context.Context) (terraformValueObjects.ResourceActionMap, error) {
	resourceActions := terraformValueObjects.ResourceActionMap{}

	err := alc.loadUpstreamDataToAWSLogQuerier()
	if err != nil {
		return resourceActions, fmt.Errorf("[alc.loadUpstreamDataToAWSLogQuerier]%v", err)
	}

	err = alc.setAWSCredentials(alc.cloudCredential)
	if err != nil {
		return nil, fmt.Errorf("[alc.setAWSCredentials]%v", err)
	}

	alc.awsSession = session.Must(session.NewSession())

	// Calculating cloud actors for managed resource drift
	for _, driftedResource := range alc.uniqueManagedDriftedResources {
		currentActions, err := alc.cloudTrailEventHistorySearch(ctx, driftedResource.ResourceType, driftedResource.InstanceID, driftedResource.Region, false)
		if err != nil {
			if err != ErrNoCloudTrailEvents {
				return nil, fmt.Errorf("[alc.cloudTrailEventHistorySearch]%v", err)
			}
			log.Errorf("[no cloud trail events found for resource %v]", driftedResource)
			continue
		}

		currentResourceName := uniqueDriftedResourceToName(driftedResource)
		resourceActions[currentResourceName] = &currentActions
	}

	alc.UpdateManagedDriftAttributeDifferences(resourceActions)

	// Overwrite the drift-resources-differences.json file with the new data.
	managedAttributeDifferencesBytes, err := json.MarshalIndent(alc.managedDriftAttributeDifferences, "", "  ")
	if err != nil {
		return resourceActions, fmt.Errorf("[json.MarshalIndent]%v", err)
	}

	err = os.WriteFile("outputs/drift-resources-differences.json", managedAttributeDifferencesBytes, 0o400)
	if err != nil {
		return resourceActions, fmt.Errorf("[os.WriteFile]%v", err)
	}
	log.Debugf("[aws_log_querier][QueryForAllResources] Wrote drift-resources-differences.json file")

	// Calculating cloud actors for new resource drift
	for id, resource := range alc.newResources {
		currentActions, err := alc.cloudTrailEventHistorySearch(ctx, resource.ResourceType, string(id), resource.Region, true)
		if err != nil {
			if err != ErrNoCloudTrailEvents {
				return nil, fmt.Errorf("[alc.cloudTrailEventHistorySearch]%v", err)
			}
			log.Errorf("[no cloud trail events found for resource %v]", resource)
			continue
		}

		currentResourceName := terraformValueObjects.ResourceName(
			resource.ResourceType + "." + hclcreate.ConvertTerraformerResourceName(resource.ResourceTerraformerName),
		)
		resourceActions[currentResourceName] = &currentActions
	}

	return resourceActions, nil
}

// UpdateManagedDriftAttributeDifferences updates the RecentActor and RecentActionTimestamp fields
// for each struct within the alc.managedDriftAttributeDifferences slice.
func (alc *AWSLogQuerier) UpdateManagedDriftAttributeDifferences(
	resourceActions terraformValueObjects.ResourceActionMap,
) {
	newAttributeDifferences := []driftDetector.AttributeDifference{}

	for _, attributeDifference := range alc.managedDriftAttributeDifferences {
		currentDifferenceResourceName := attributeDifferenceToResourceName(attributeDifference)

		if _, ok := resourceActions[currentDifferenceResourceName]; ok {
			resourceAction := resourceActions[currentDifferenceResourceName]

			attributeDifference.RecentActor = resourceAction.Modifier.Actor
			attributeDifference.RecentActionTimestamp = resourceAction.Modifier.Timestamp
		}
		newAttributeDifferences = append(newAttributeDifferences, attributeDifference)
	}

	alc.managedDriftAttributeDifferences = newAttributeDifferences
}

// cloudTrailEventHistorySearch runs AWS CLI commands to pull data on who modified and created the cloud resource in question.
func (alc *AWSLogQuerier) cloudTrailEventHistorySearch(_ context.Context, resourceType string, resourceID string, resourceRegion string, isNewToTerraform bool) (terraformValueObjects.ResourceActions, error) {
	svc := cloudtrail.New(alc.awsSession, aws.NewConfig().WithRegion(resourceRegion))

	maxResults := int64(50)

	result, err := svc.LookupEvents(&cloudtrail.LookupEventsInput{
		MaxResults: &maxResults,
		LookupAttributes: []*cloudtrail.LookupAttribute{
			{
				AttributeKey:   aws.String("ResourceName"),
				AttributeValue: aws.String(resourceID),
			},
		},
	})
	if err != nil {
		return terraformValueObjects.ResourceActions{}, fmt.Errorf("[alc.cloudTrailClient.LookupEvents]%v", err)
	}

	return alc.ExtractDataFromResourceResult(result.Events, resourceType, isNewToTerraform)
}

// ExtractDataFromResourceResult parses the log response from the provider API
// and extracts needed data (namely who made the most recent relevant change to the resource).
func (alc *AWSLogQuerier) ExtractDataFromResourceResult(resourceResult []*cloudtrail.Event, resourceType string, isNewToTerraform bool) (terraformValueObjects.ResourceActions, error) {
	resourceActions := terraformValueObjects.ResourceActions{}

	if len(resourceResult) == 0 {
		return resourceActions, ErrNoCloudTrailEvents
	}

	resourceType = string(alc.resourceToCloudTrailType[resourceType])

	eventsLength := len(resourceResult)
	isComplete := false
	isModificationIdentified := false
	i := 0

	for !isComplete {
		event := resourceResult[i]
		classification := determineActionClass(*event.EventName)
		cleanedTime := event.EventTime.Format("2006-01-02")

		// check to ensure that ResourceType is present within one of the event.Resources elements
		// if not, move on to the next event
		isValidResourceType := false
		for _, resource := range event.Resources {
			if *resource.ResourceType == resourceType || *resource.ResourceType == "" {
				isValidResourceType = true
				break
			}
		}
		if !isValidResourceType {
			if i+1 >= eventsLength {
				isComplete = true
			}
			i++
			continue
		}

		switch classification {
		case "creation":
			if isNewToTerraform {
				resourceActions.Creator = &terraformValueObjects.CloudActorTimeStamp{
					Actor:     terraformValueObjects.CloudActor(*event.Username),
					Timestamp: terraformValueObjects.Timestamp(cleanedTime),
				}
				return resourceActions, nil
			}
		case "modification":
			if !isModificationIdentified {
				isModificationIdentified = true
				resourceActions.Modifier = &terraformValueObjects.CloudActorTimeStamp{
					Actor:     terraformValueObjects.CloudActor(*event.Username),
					Timestamp: terraformValueObjects.Timestamp(cleanedTime),
				}
				if !isNewToTerraform {
					return resourceActions, nil
				}
			}
		}

		if i+1 >= eventsLength {
			isComplete = true
		}
		i++
	}

	return resourceActions, nil
}

// setAWSCredentials loads and sets as environment variables AWS credentials for a given AWS account.
func (alc *AWSLogQuerier) setAWSCredentials(credential terraformValueObjects.Credential) error {
	envVars := new(AWSEnvironment)
	err := json.Unmarshal([]byte(credential), &envVars)
	if err != nil {
		return fmt.Errorf("[json.Unmarshal] %w", err)
	}
	log.Debugf("setting AWS credentials for account: %+v", envVars)

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
