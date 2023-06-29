package terraformWorkspace

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/hcl/v2/hclsimple"
	log "github.com/sirupsen/logrus"

	terraformValueObjects "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_value_objects"
	"github.com/dragondrop-cloud/cloud-concierge/main/internal/interfaces"
)

type ContainerBackendConfig struct {
	// AWSRegion is the region of the AWS account that contains the S3 bucket.
	AWSRegion string

	// WorkspaceDirectories is a slice of directories that contains terraform workspaces within the user repo.
	WorkspaceDirectories WorkspaceDirectoriesDecoder `required:"true"`

	// ContainerName is the name of the S3 bucket that contains the state files.
	ContainerName string `required:"true"`

	// DivisionCloudCredentials is a map between a Division and request cloud credentials.
	DivisionCloudCredentials terraformValueObjects.DivisionCloudCredentialDecoder `required:"true"`
}

// findTerraformWorkspaces searches a repo for terraform workspaces.
func findTerraformWorkspaces(ctx context.Context, dragonDrop interfaces.DragonDrop, workspaceDirectories []string, backendType string) (map[string]string, map[string]interface{}, error) {
	workspaceToDirectory := make(map[string]string)
	workspaceToBackendDetails := make(map[string]interface{})

	dragonDrop.PostLog(ctx, "Searching for terraform workspaces names.")
	for _, directory := range workspaceDirectories {
		workspace, backendDetails, err := searchDirectoryForWorkspaceName(ctx, directory, backendType)
		if err != nil {
			return nil, nil, fmt.Errorf("[found_terraform_workspaces][error searching directory %s]%w", directory, err)
		}

		workspaceToDirectory[workspace] = directory
		workspaceToBackendDetails[workspace] = backendDetails
	}
	return workspaceToDirectory, workspaceToBackendDetails, nil
}

// searchDirectoryForWorkspaceName searches a directory for a terraform workspace name.
func searchDirectoryForWorkspaceName(ctx context.Context, directory string, backendType string) (string, interface{}, error) {
	directory = cleanDirectoryName(directory)
	tfFiles := []string{"versions.tf", "main.tf"}
	tfFiles = append(tfFiles, getAllTFFiles(ctx, directory)...)

	for _, tfFile := range tfFiles {
		workspace, backendDetails, found := getWorkspaceByFile(ctx, directory, tfFile, backendType)
		if found {
			log.Debug(fmt.Sprintf("[search_directory_for_workspace_name][found workspace %s]", workspace))
			return workspace, backendDetails, nil
		}
	}

	return "", nil, fmt.Errorf("[search_directory_for_workspace_name][error searching directory %s]", directory)
}

// getAllTFFiles searches a directory for all terraform files.
func getAllTFFiles(ctx context.Context, directory string) []string {
	files, err := os.ReadDir(fmt.Sprintf("repo/%s", directory))
	if err != nil {
		return make([]string, 0)
	}

	tfFiles := make([]string, 0)
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if strings.HasSuffix(file.Name(), ".tf") {
			tfFiles = append(tfFiles, file.Name())
		}
	}

	return tfFiles
}

// cleanDirectoryName removes any leading or trailing slashes from a directory name.
func cleanDirectoryName(directory string) string {
	directory = strings.Trim(directory, " ")
	directory = strings.Trim(directory, "/")
	return directory
}

// getWorkspaceByFile searches a given file for a terraform workspace.
func getWorkspaceByFile(ctx context.Context, directory string, fileName string, backendType string) (string, interface{}, bool) {
	filePath := fmt.Sprintf("repo/%s/%s", directory, fileName)
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return "", nil, false
	}

	// checking to see if a Terraform Cloud workspace configuration exists, and if so, extracting the workspace data.
	workspace, err := extractTFCloudWorkspaceNameIfExists(ctx, fileContent)
	if err != nil && workspace == "" {
		return "", nil, false
	}

	// checking to see if a major-cloud-provider hosted Terraform backend configuration exists, and if so, extracting the configuration data.
	details, err := extractBackendDetails(ctx, fileContent, backendType)
	if err != nil {
		return "", nil, false
	}

	return workspace, details, true
}

// TerraformCloudFile is a struct representation of a terraform block
type TerraformCloudFile struct {
	Cloud CloudBlock `hcl:"terraform,block"`
}

type CloudBlock struct {
	Workspace WorkspaceBlock `hcl:"cloud,block"`
}

type WorkspaceBlock struct {
	Details WorkspaceDetails `hcl:"workspaces,block"`
}

type WorkspaceDetails struct {
	Name string `hcl:"name,attr"`
}

// extractTFCloudWorkspaceNameIfExists extracts the workspace name from a Terraform versions.tf file if exits.
func extractTFCloudWorkspaceNameIfExists(ctx context.Context, fileContent []byte) (string, error) {
	var config TerraformCloudFile
	err := hclsimple.Decode("placeholder.hcl", fileContent, nil, &config)
	if err != nil {
		return "", err
	}

	return config.Cloud.Workspace.Details.Name, nil
}

// S3 Related backend configuration parsing
// S3TerraformBackend is a struct representation of a terraform backend file for s3
type S3TerraformBackend struct {
	TerraformBlock S3TerraformBlock `hcl:"terraform,block"`
}

// S3TerraformBlock is a struct representation of a terraform block for s3
type S3TerraformBlock struct {
	Backend S3BackendBlock `hcl:"backend,block"`
}

// S3BackendBlock is a struct representation of a terraform backend block for s3
type S3BackendBlock struct {
	Name   string `hcl:"name,label"`
	Bucket string `hcl:"bucket,attr"`
	Key    string `hcl:"key,attr"`
	Region string `hcl:"region,attr"`
}

// Azurerm Related backend configuration parsing
// AzurermTerraformBackend is a struct representation of a terraform backend file for azurerm
type AzurermTerraformBackend struct {
	TerraformBlock AzurermTerraformBlock `hcl:"terraform,block"`
}

// AzurermTerraformBlock is a struct representation of a terraform block for azurerm
type AzurermTerraformBlock struct {
	Backend AzureBackendBlock `hcl:"backend,block"`
}

// AzureBackendBlock is a struct representation of a terraform backend block for azurerm
type AzureBackendBlock struct {
	Name               string `hcl:"name,label"`
	ResourceGroupName  string `hcl:"resource_group_name,attr"`
	StorageAccountName string `hcl:"storage_account_name,attr"`
	ContainerName      string `hcl:"container_name,attr"`
	Key                string `hcl:"key,attr"`
}

// GCS Related backend configuration parsing
// GCSTerraformBackend is a struct representation of a terraform backend file for gcs.
type GCSTerraformBackend struct {
	TerraformBlock GCSTerraformBlock `hcl:"terraform,block"`
}

// GCSTerraformBlock parses the "terraform" block for GCS
type GCSTerraformBlock struct {
	Backend GCSBackendBlock `hcl:"backend,block"`
}

// GCSBackendBlock parses the backend block for GCS
type GCSBackendBlock struct {
	Name   string `hcl:"name,label"`
	Bucket string `hcl:"bucket,attr"`
	Prefix string `hcl:"prefix,attr"`
}

// extractBackendDetails extracts the backend details from a .tf file if it exists.
func extractBackendDetails(ctx context.Context, fileContent []byte, backendType string) (interface{}, error) {
	filePath := "backend.hcl"

	switch backendType {
	case "s3":
		var config S3TerraformBackend
		err := hclsimple.Decode(filePath, fileContent, nil, &config)
		if err != nil {
			return "", err
		}

		return config.TerraformBlock.Backend, nil
	case "azurerm":
		var config AzurermTerraformBackend
		err := hclsimple.Decode(filePath, fileContent, nil, &config)
		if err != nil {
			return "", err
		}

		return config.TerraformBlock.Backend, nil
	case "gcs":
		var config GCSTerraformBackend
		err := hclsimple.Decode(filePath, fileContent, nil, &config)
		if err != nil {
			return "", err
		}

		return config.TerraformBlock.Backend, nil
	default:
		return "Not yet supported", nil
	}
}
