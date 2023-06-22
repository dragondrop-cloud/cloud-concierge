package identifyCloudActors

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/Jeffail/gabs/v2"
	"github.com/dragondrop-cloud/driftmitigation/hclcreate"
	resourcesCalculator "github.com/dragondrop-cloud/driftmitigation/implementations/resources_calculator"
	driftDetector "github.com/dragondrop-cloud/driftmitigation/implementations/terraform_managed_resources_drift_detector/drift_detector"
	terraformValueObjects "github.com/dragondrop-cloud/driftmitigation/implementations/terraform_value_objects"
)

// GoogleLogQuerier implements the LogQuerier interface for Google Cloud.
type GoogleLogQuerier struct {

	// authToken is an auth token for the Google Cloud REST API for a particular GCP project
	authToken string

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
}

// NewGoogleLogQuerier instantiates a new instance of GoogleLogQuerier
func NewGoogleLogQuerier(divisionToCredentials terraformValueObjects.DivisionCloudCredentialDecoder) (LogQuerier, error) {
	return &GoogleLogQuerier{
		divisionToCredentials: divisionToCredentials,
	}, nil
}

// loadUpstreamDataToGoogleLogQuerier loads all data needed for querying logs from upstream
// saved data sources.
func (glc *GoogleLogQuerier) loadUpstreamDataToGoogleLogQuerier() error {
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

	glc.divisionToUniqueManagedDriftedResources = divToUniqueDriftedResources
	glc.divisionToNewResources = divToNewResources
	glc.managedDriftAttributeDifferences = attributeDifferences
	return nil
}

// QueryForAllResources coordinates calls of QueryForResourcesInDivision for all
// divisions from which drifted resources have been identified.
func (glc *GoogleLogQuerier) QueryForAllResources(ctx context.Context) (terraformValueObjects.DivisionResourceActions, error) {
	divisionToResourceActions := terraformValueObjects.DivisionResourceActions{}

	err := glc.loadUpstreamDataToGoogleLogQuerier()
	if err != nil {
		return divisionToResourceActions, fmt.Errorf("[glc.loadUpstreamDataToGoogleLogQuerier]%v", err)
	}

	for division := range glc.divisionToCredentials {
		divisionResourceActions, err := glc.QueryForResourcesInDivision(ctx, division)
		if err != nil {
			return divisionToResourceActions, fmt.Errorf("[glc.QueryForResourcesInDivision]%v", err)
		}
		divisionToResourceActions[division] = divisionResourceActions
	}
	return divisionToResourceActions, nil
}

// QueryForResourcesInDivision coordinates calls of QuerySingleResource for all resources within a division.
func (glc *GoogleLogQuerier) QueryForResourcesInDivision(ctx context.Context, division terraformValueObjects.Division) (map[terraformValueObjects.ResourceName]terraformValueObjects.ResourceActions, error) {
	divisionResourceActions := map[terraformValueObjects.ResourceName]terraformValueObjects.ResourceActions{}

	err := glc.gcloudAuthTokenFromServiceAccount(division)
	if err != nil {
		return divisionResourceActions, fmt.Errorf("[glc.gcloudAuthTokenFromServiceAccount]%v", err)
	}
	dragondropDivision := terraformValueObjects.Division("google-" + string(division))

	// Calculating cloud actors for managed resource drift
	currentUniqueDriftedResources, ok := glc.divisionToUniqueManagedDriftedResources[dragondropDivision]
	if ok {
		for _, driftedResource := range currentUniqueDriftedResources {
			resourceActions, err := glc.adminLogSearch(ctx, division, driftedResource.InstanceID, false)
			if err != nil {
				return divisionResourceActions, fmt.Errorf("[glc.QuerySingleResource]%v", err)
			}

			currentResourceName := uniqueDriftedResourceToName(driftedResource)
			divisionResourceActions[currentResourceName] = resourceActions
		}

		glc.UpdateManagedDriftAttributeDifferences(divisionResourceActions)

		// Overwrite the drift-resources-differences.json file with the new data.
		managedAttributeDifferencesBytes, err := json.MarshalIndent(glc.managedDriftAttributeDifferences, "", "  ")
		if err != nil {
			return divisionResourceActions, fmt.Errorf("[json.MarshalIndent]%v", err)
		}

		err = os.WriteFile("mappings/drift-resources-differences.json", managedAttributeDifferencesBytes, 0400)
		if err != nil {
			return divisionResourceActions, fmt.Errorf("[os.WriteFile]%v", err)
		}
	}

	// Calculating cloud actors for new resource drift
	currentNewResources, ok := glc.divisionToNewResources[dragondropDivision]
	if ok {
		for id, resource := range currentNewResources {
			resourceActions, err := glc.adminLogSearch(ctx, division, string(id), true)
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
func (glc *GoogleLogQuerier) UpdateManagedDriftAttributeDifferences(
	divisionResourceActions map[terraformValueObjects.ResourceName]terraformValueObjects.ResourceActions,
) {
	newAttributeDifferences := []driftDetector.AttributeDifference{}

	for _, attributeDifference := range glc.managedDriftAttributeDifferences {
		currentDifferenceResourceName := attributeDifferenceToResourceName(attributeDifference)

		if _, ok := divisionResourceActions[currentDifferenceResourceName]; ok {
			resourceAction := divisionResourceActions[currentDifferenceResourceName]

			attributeDifference.RecentActor = resourceAction.Modifier.Actor
			attributeDifference.RecentActionTimestamp = resourceAction.Modifier.Timestamp
		}
		newAttributeDifferences = append(newAttributeDifferences, attributeDifference)
	}

	glc.managedDriftAttributeDifferences = newAttributeDifferences
}

// adminLogSearch pulls logs for a single resource from the cloud provider.
func (glc *GoogleLogQuerier) adminLogSearch(
	ctx context.Context, division terraformValueObjects.Division, resourceID string, isNewToTerraform bool,
) (terraformValueObjects.ResourceActions, error) {
	result, err := glc.queryGCPAPI(ctx, division, resourceID)
	if err != nil {
		return terraformValueObjects.ResourceActions{}, fmt.Errorf("[glc.queryGCPAPI]%w", err)
	}

	resourceActions, err := glc.ExtractDataFromResourceResult(result, "", isNewToTerraform)
	if err != nil {
		return terraformValueObjects.ResourceActions{}, fmt.Errorf("[glc.ExtractDataFromResourceResult]%w", err)
	}
	return resourceActions, nil
}

// GCPAdminLogPostBody contains the fields needed for the body of a post request to the GCP api
// for getting admin action log data.
type GCPAdminLogPostBody struct {

	// ResourceNames are the names of one or more parent resources from which to retrieve log entries.
	// For our use case, each value will always take the form of "projects/[PROJECT_ID]"
	ResourceNames []string `json:"resourceNames"`

	// Filter is the filter of the resource specified within resourceNames.
	Filter string `json:"filter"`

	// OrderBy is the timeline order of returned results.
	OrderBy string `json:"orderBy"`

	// PageSize is the number of records to return.
	PageSize int `json:"pageSize"`
}

// queryGCPAPI sends a REST API POST request to the Google Cloud endpoint corresponding to admin
// log querying.
func (glc *GoogleLogQuerier) queryGCPAPI(ctx context.Context, division terraformValueObjects.Division, resourceID string) ([]byte, error) {
	logFilterString := glc.generateLogFilter(division, resourceID)
	jsonBody, err := json.Marshal(&GCPAdminLogPostBody{
		ResourceNames: []string{fmt.Sprintf("projects/%v", division)},
		Filter:        logFilterString,
		OrderBy:       "timestamp desc",
		PageSize:      1000,
	})

	if err != nil {
		return []byte{}, fmt.Errorf("[glc.queryGCPAPI][error in json marshal]%v", err)
	}

	request, err := glc.newRequest(
		ctx,
		"Query logs",
		glc.authToken,
		"https://logging.googleapis.com/v2/entries:list",
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return []byte{}, fmt.Errorf("[glc.newRequest]%w", err)
	}

	response, err := glc.httpClient.Do(request)
	if err != nil {
		return []byte{}, fmt.Errorf("[glc.httpClient.Do][error in executing request]%w", err)
	}

	if response.StatusCode != 200 {
		return []byte{}, fmt.Errorf("[glc.queryGCPAPI POST request][was unsuccessful, with the server returning: %v]", response.StatusCode)
	}

	// Read in response body to bytes array.
	outputBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("[io.ReadAll][error in reading response into bytes array]%w", err)
	}
	err = response.Body.Close()
	if err != nil {
		return []byte{}, fmt.Errorf("[response.Body.Close()][error in closing response]%w", err)
	}

	return outputBytes, nil
}

// generateLogFilter generates a string formatted for filtering admin query logs within the GCP API.
func (glc *GoogleLogQuerier) generateLogFilter(division terraformValueObjects.Division, resourceID string) string {
	logNameFilter := fmt.Sprintf("logName=projects/%v", division) + "/logs/cloudaudit.googleapis.com%2Factivity"

	resourceTypeFilter := fmt.Sprintf("protoPayload.resourceName=%v", resourceID)

	combinedFilter := fmt.Sprintf("%v AND %v", logNameFilter, resourceTypeFilter)

	return combinedFilter
}

// newRequest creates a new http request.
func (glc *GoogleLogQuerier) newRequest(ctx context.Context, requestName string, authToken string, requestPath string, body io.Reader) (*http.Request, error) {
	request, err := http.NewRequestWithContext(ctx, "POST", requestPath, body)
	if err != nil {
		return nil, fmt.Errorf("[glc.newRequest][error in http request instantiation with name: %s, err: %v]", requestName, err)
	}

	request.Header = http.Header{
		"Authorization": {fmt.Sprintf("Bearer %v", authToken)},
		"Content-Type":  {"application/json"},
	}

	return request, nil
}

// Entries is a struct representing the entries field of a GCP logging query response.
type Entries struct {
	Entries []Entry `json:"entries"`
}

// Entry is a struct representing a single entry in a GCP logging query response.
type Entry struct {
	ProtoPayload     ProtoPayload `json:"protoPayload"`
	ReceiveTimestamp string       `json:"receiveTimestamp"`
}

// ProtoPayload is a struct representing the protoPayload field of a GCP logging query response.
type ProtoPayload struct {
	MethodName         string             `json:"methodName"`
	AuthenticationInfo AuthenticationInfo `json:"authenticationInfo"`
}

// AuthenticationInfo is a struct representing the authenticationInfo field of the ProtoPayload component
// of a GCP logging query response.
type AuthenticationInfo struct {
	PrincipalEmail string `json:"principalEmail"`
}

// ExtractDataFromResourceResult parses the log response from the provider API
// and extracts needed data (namely who made the most recent relevant change to the resource).
func (glc *GoogleLogQuerier) ExtractDataFromResourceResult(resourceResult []byte, resourceType string, isNewToTerraform bool) (terraformValueObjects.ResourceActions, error) {
	resourceActions := terraformValueObjects.ResourceActions{}
	var entries Entries
	if err := json.Unmarshal(resourceResult, &entries); err != nil {
		return resourceActions, fmt.Errorf("failed to parse resource result: %v", err)
	}

	// Algorithm: Iterate through each response entry.
	// If "Create" is in the action, stop.
	// If "Modify", capture the first one.
	// ---- If isNewToTerraform continue until we hit "create", otherwise return immediately.
	isModifyIdentified := false

	for _, entry := range entries.Entries {
		classification := determineActionClass(entry.ProtoPayload.MethodName)

		switch classification {
		case "creation":
			resourceActions.Creator = terraformValueObjects.CloudActorTimeStamp{
				Actor:     terraformValueObjects.CloudActor(entry.ProtoPayload.AuthenticationInfo.PrincipalEmail),
				Timestamp: terraformValueObjects.Timestamp(entry.ReceiveTimestamp[:10]),
			}
			return resourceActions, nil
		case "modification":
			if !isModifyIdentified {
				isModifyIdentified = true
				resourceActions.Modifier = terraformValueObjects.CloudActorTimeStamp{
					Actor:     terraformValueObjects.CloudActor(entry.ProtoPayload.AuthenticationInfo.PrincipalEmail),
					Timestamp: terraformValueObjects.Timestamp(entry.ReceiveTimestamp[:10]),
				}
				if !isNewToTerraform {
					return resourceActions, nil
				}
			}
		}
	}
	return resourceActions, nil
}

// gcloudAuthTokenFromServiceAccount gets an authentication token for REST API requests from the
// passed service account keys.
func (glc *GoogleLogQuerier) gcloudAuthTokenFromServiceAccount(division terraformValueObjects.Division) error {
	account, err := glc.parseGCPServiceAccountEmailAddress(division)
	if err != nil {
		return fmt.Errorf("[gcloud_authentication][error parsing service account email address]%w", err)
	}

	// Authenticate gcloud for current division
	keyFilePath := fmt.Sprintf("--key-file=current_cloud/credentials/google-%s.json", division)
	authArgs := []string{"auth", "activate-service-account", string(account), keyFilePath}

	_, err = executeCommand("gcloud", authArgs...)
	if err != nil {
		return fmt.Errorf("[gcloud_authentication][gcloud auth activate-service-account, failed to authenticate]%w", err)
	}

	printAccessTokenArgs := []string{"auth", "print-access-token", string(account)}
	token, err := executeCommand("gcloud", printAccessTokenArgs...)
	if err != nil {
		return fmt.Errorf("[executeCommand][gcloud auth print-access-token]%w", err)
	}

	glc.authToken = strings.Replace(token, "\n", "", -1)

	return nil
}

// parseGCPServiceAccountEmailAddress pulls out the service account email address from the service account
// key file.
func (glc *GoogleLogQuerier) parseGCPServiceAccountEmailAddress(division terraformValueObjects.Division) (terraformValueObjects.Account, error) {
	credentialString := glc.divisionToCredentials[division]
	serviceAccountParsed, err := gabs.ParseJSON([]byte(credentialString))
	if err != nil {
		return "", fmt.Errorf("[parse_gcp_service_account_email_address][error parsing JSON with gabs.ParseJSON]%w", err)
	}

	account, ok := serviceAccountParsed.Path("client_email").Data().(string)
	if !ok {
		return "", fmt.Errorf("[parse_gcp_service_account_email_address][client_email not found within service account key]")
	}

	return terraformValueObjects.Account(account), nil
}
