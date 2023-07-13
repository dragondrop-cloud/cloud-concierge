package terraformWorkspace

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/Jeffail/gabs/v2"
	log "github.com/sirupsen/logrus"

	"github.com/dragondrop-cloud/cloud-concierge/main/internal/interfaces"
)

type WorkspaceDirectoriesDecoder []string

func (d *WorkspaceDirectoriesDecoder) Decode(value string) error {
	value = strings.Trim(value, "[]")
	arrayValue := strings.Split(value, ",")

	*d = make([]string, 0)
	for _, directory := range arrayValue {
		*d = append(*d, strings.Trim(directory, "\""))
	}
	return nil
}

// TerraformCloudConfig is a struct containing the variables that determine the specific
// behavior of the TerraformCloud struct.
type TerraformCloudConfig struct {
	// StateBackend is the name of the backend used for storing State.
	StateBackend string `required:"true"`

	// TerraformCloudOrganization is the name of the organization within TerraformCloudFile Cloud
	TerraformCloudOrganization string `required:"true"`

	// TerraformCloudToken is the auth token to access TerraformCloudFile Cloud programmatically.
	TerraformCloudToken string `required:"true"`

	// WorkspaceDirectories is a slice of directories that contains terraform workspaces within the user repo.
	WorkspaceDirectories WorkspaceDirectoriesDecoder `required:"true"`
}

// TerraformCloud is a struct that implements the interfaces.TerraformWorkspace interface.
type TerraformCloud struct {
	// httpClient is the client used to send http requests
	httpClient http.Client

	// config contains the variables that determine the specific behavior of the TerraformCloud struct
	config TerraformCloudConfig

	// dragonDrop is an implementation of the interfaces.dragonDrop interface for communicating with the
	// dragondrop API.
	dragonDrop interfaces.DragonDrop
}

// NewTerraformCloud creates a new instance of the TerraformCloud struct.
func NewTerraformCloud(ctx context.Context, config TerraformCloudConfig, dragonDrop interfaces.DragonDrop) interfaces.TerraformWorkspace {
	dragonDrop.PostLog(ctx, "Created TFWorkspace client.")
	return &TerraformCloud{config: config, dragonDrop: dragonDrop}
}

func (c *TerraformCloud) FindTerraformWorkspaces(ctx context.Context) (map[string]string, error) {
	workspaceToDirectory := make(map[string]string)

	c.dragonDrop.PostLog(ctx, "Searching for terraform workspaces names.")
	for _, directory := range c.config.WorkspaceDirectories {
		workspace, err := c.searchDirectoryForWorkspaceName(ctx, directory)
		if err != nil {
			return nil, fmt.Errorf("[found_terraform_workspaces][error searching directory %s]%w", directory, err)
		}

		workspaceToDirectory[workspace] = directory
	}
	return workspaceToDirectory, nil
}

// searchDirectoryForWorkspaceName searches a directory for a terraform workspace name.
func (c *TerraformCloud) searchDirectoryForWorkspaceName(ctx context.Context, directory string) (string, error) {
	directory = cleanDirectoryName(directory)
	tfFiles := []string{"versions.tf", "main.tf"}
	tfFiles = append(tfFiles, getAllTFFiles(ctx, directory)...)

	for _, tfFile := range tfFiles {
		workspace, _, found := getWorkspaceByFile(ctx, directory, tfFile, "terraform")
		if found {
			log.Debug(fmt.Sprintf("[search_directory_for_workspace_name][found workspace %s]", workspace))
			return workspace, nil
		}
	}

	return "", fmt.Errorf("[search_directory_for_workspace_name][error searching directory %s]", directory)
}

// DownloadWorkspaceState downloads from the remote TerraformCloudFile backend the latest state file
// for each "workspace".
func (c *TerraformCloud) DownloadWorkspaceState(ctx context.Context, WorkspaceToDirectory map[string]string) error {
	c.dragonDrop.PostLog(ctx, "Beginning download of state files to local memory.")

	for workspaceName := range WorkspaceToDirectory {
		err := c.getWorkspaceState(ctx, workspaceName)

		if err != nil {
			return fmt.Errorf("[download_workspace_state][error getting state for %s]%w", workspaceName, err)
		}
	}

	c.dragonDrop.PostLog(ctx, "Done with download of state files to local memory.")
	return nil
}

// getWorkspaceState downloads from the remote TerraformCloudFile backend a single "workspace"'s latest
// state file.
func (c *TerraformCloud) getWorkspaceState(ctx context.Context, workspaceName string) error {
	workspaceID, err := c.getWorkspaceID(ctx, workspaceName)
	if err != nil {
		return err
	}

	requestName := "getWorkspaceStateByTestingAllS3Credentials"
	requestPath := fmt.Sprintf("https://app.terraform.io/api/v2/workspaces/%v/current-state-version", workspaceID)

	request, err := c.buildTFCloudHTTPRequest(ctx, requestName, "GET", requestPath)

	if err != nil {
		return fmt.Errorf("[get_workspace_state][error creating request %s]%w", requestName, err)
	}

	jsonResponseBytes, err := c.terraformCloudRequest(request, requestName)

	if err != nil {
		return fmt.Errorf("[get_workspace_state][error executing terraform cloud request]%w", err)
	}

	rawStateURL, err := c.extractRawStateURL(jsonResponseBytes)

	if err != nil {
		return fmt.Errorf("[get_workspace_state][error extracting raw state url]%w", err)
	}

	jsonResponseBytes, err = c.getRawTerraformStateFile(rawStateURL)

	if err != nil {
		return fmt.Errorf("[get_workspace_state][error getting raw terraform state file]%w", err)
	}

	_ = os.MkdirAll("state_files", 0660)
	fileOutPath := fmt.Sprintf("state_files/%v.json", workspaceName)

	err = os.WriteFile(fileOutPath, jsonResponseBytes, 0400)
	if err != nil {
		return fmt.Errorf("[get_workspace_state][error saving state file to memory]%w", err)
	}

	return nil
}

// getWorkspaceID calls the TerraformCloudFile Cloud API and gets the workspace ID for the
// relevant workspace name in the relevant organization.
func (c *TerraformCloud) getWorkspaceID(ctx context.Context, workspaceName string) (string, error) {
	requestName := "getWorkspaceID"
	requestPath := fmt.Sprintf("https://app.terraform.io/api/v2/organizations/%v/workspaces/%v", c.config.TerraformCloudOrganization, workspaceName)

	request, err := c.buildTFCloudHTTPRequest(ctx, requestName, "GET", requestPath)

	if err != nil {
		return "", fmt.Errorf("[get_workspace_id][error building terraform cloud request %s]%w", requestName, err)
	}

	jsonResponseBytes, err := c.terraformCloudRequest(request, requestName)

	if err != nil {
		return "", err
	}

	return c.extractWorkspaceID(jsonResponseBytes)
}

// buildTFCloudHTTPRequest structures a request to the TerraformCloudFile Cloud api.
func (c *TerraformCloud) buildTFCloudHTTPRequest(ctx context.Context, requestName string, method string, requestPath string) (*http.Request, error) {
	request, err := http.NewRequestWithContext(ctx, method, requestPath, nil)
	if err != nil {
		return nil, fmt.Errorf("[%v] error in http request instantiation: %v", requestName, err)
	}

	request.Header = http.Header{
		"Authorization": {"Bearer " + c.config.TerraformCloudToken},
		"Content-Type":  {"application/vnd.api+json"},
	}

	return request, nil
}

// terraformCloudRequest build, executes, and processes an API call to the TerraformCloudFile Cloud API.
func (c *TerraformCloud) terraformCloudRequest(request *http.Request, requestName string) ([]byte, error) {

	response, err := c.httpClient.Do(request)

	if err != nil {
		return nil, fmt.Errorf("[terraform_cloud_request][error in http GET request to TerraformCloudFile cloud: %s]%w", requestName, err)
	}

	defer response.Body.Close()
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("[terraform_cloud_request][request %s was unsuccessful, with the server returning: %d]", requestName, response.StatusCode)
	}

	// Read in response body to bytes array.
	outputBytes, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, fmt.Errorf("[error in reading response into bytes array in request: %s]%w", requestName, err)
	}

	return outputBytes, nil
}

// extractRawStateURL is a helper function that uses the gabs library to pull out the raw state URL
// from a TerraformCloudFile Cloud API response of the current workspace.
func (c *TerraformCloud) extractRawStateURL(jsonBytes []byte) (string, error) {
	jsonParsed, err := gabs.ParseJSON(jsonBytes)
	if err != nil {
		return "", fmt.Errorf("[extract_raw_state_url][error in parsing bytes array to json via 'gabs']%w", err)
	}

	value, ok := jsonParsed.Path("data.attributes.hosted-state-download-url").Data().(string)
	if !ok {
		return "", fmt.Errorf("[extract_raw_state_url][unable to find hosted-state-download-url]%w", err)
	}

	return value, nil
}

// extractWorkspaceID is a helper function that uses the gabs library to pull out the workspace ID
// from a TerraformCloudFile Cloud API response.
func (c *TerraformCloud) extractWorkspaceID(jsonBytes []byte) (string, error) {
	jsonParsed, err := gabs.ParseJSON(jsonBytes)
	if err != nil {
		return "", fmt.Errorf("[extract_workspace_id][error in parsing bytes array to json via 'gabs']%w", err)
	}

	value, ok := jsonParsed.Path("data.id").Data().(string)
	if !ok {
		return "", fmt.Errorf("[extract_workspace_id][unable to find workspace id]%w", err)
	}

	return value, nil
}

// getRawTerraformStateFile gets the raw terraform state file contents from the passed url.
func (c *TerraformCloud) getRawTerraformStateFile(rawStateURL string) ([]byte, error) {
	resp, err := http.Get(rawStateURL) //nolint

	if err != nil {
		return nil, fmt.Errorf("[get_raw_terraform_state_file][error executing http get request]%w", err)
	}
	defer resp.Body.Close()

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("[get_raw_terraform_state_file][was unsuccessful, with the server returning: %d]", resp.StatusCode)
	}

	// Read in response body to bytes array.
	outputBytes, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("[get_raw_terraform_state_file][error in reading response into bytes array]%w", err)
	}

	return outputBytes, nil
}
