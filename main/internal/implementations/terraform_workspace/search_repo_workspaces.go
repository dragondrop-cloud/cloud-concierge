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
	workspace, err := extractTFCloudWorkspaceNameIfExists(ctx, filePath, fileContent)
	if err != nil && workspace == "" {
		return "", nil, false
	}

	// checking to see if a major-cloud-provider hosted Terraform backend configuration exists, and if so, extracting the configuration data.
	details, err := extractBackendDetails(ctx, filePath, fileContent, backendType)
	if err != nil {
		return "", nil, false
	}

	return workspace, details, true
}

// TerraformCloudFile is a struct representation of a terraform versions.tf file.
type TerraformCloudFile struct {
	Cloud struct {
		Organization string
		Workspaces   struct {
			Name string   `hcl:"name,optional"`
			Tags []string `hcl:"tags,optional"`
		} `hcl:"workspaces,block"`
	} `hcl:"cloud,block"`
}

// extractTFCloudWorkspaceNameIfExists extracts the workspace name from a terraform versions.tf file if exits.
func extractTFCloudWorkspaceNameIfExists(ctx context.Context, filePath string, fileContent []byte) (string, error) {
	var config TerraformCloudFile
	err := hclsimple.Decode(filePath, fileContent, nil, &config)
	if err != nil {
		return "", err
	}

	return config.Cloud.Workspaces.Name, nil
}

// TerraformBackendS3 is a struct representation of a terraform backend.tf file for aws provider.
type TerraformBackendS3 struct {
	BackendS3 BackendS3 `hcl:"backend,s3"`
}

// BackendS3 is a struct representation of a terraform backend.tf file for s3 backend.
type BackendS3 struct {
	Bucket string `hcl:"bucket,optional"`
	Key    string `hcl:"key,optional"`
	Region string `hcl:"region,optional"`
}

// TerraformBackendAzure is a struct representation of a terraform backend.tf file for azure provider.
type TerraformBackendAzure struct {
	BackendAzure BackendAzure `hcl:"backend,azurerm"`
}

// BackendAzure is a struct representation of a terraform backend.tf file for azure blob storage.
type BackendAzure struct {
	ResourceGroupName  string `hcl:"resource_group_name,optional"`
	StorageAccountName string `hcl:"storage_account_name,optional"`
	ContainerName      string `hcl:"container_name,optional"`
	Key                string `hcl:"key,optional"`
}

// TerraformBackendGoogle is a struct representation of a terraform backend.tf file for google provider.
type TerraformBackendGoogle struct {
	BackendGCS BackendGCS `hcl:"backend,gcs"`
}

// BackendGCS is a struct representation of a terraform backend.tf file for gcs.
type BackendGCS struct {
	Bucket string `hcl:"bucket,optional"`
	Prefix string `hcl:"prefix,optional"`
}

// extractBackendDetails extracts the backend details from a .tf file if it exists.
func extractBackendDetails(ctx context.Context, filePath string, fileContent []byte, backendType string) (interface{}, error) {
	switch backendType {
	case "s3":
		var config TerraformBackendS3
		err := hclsimple.Decode(filePath, fileContent, nil, &config)
		if err != nil {
			return "", err
		}

		return config.BackendS3, nil
	case "azurerm":
		var config TerraformBackendAzure
		err := hclsimple.Decode(filePath, fileContent, nil, &config)
		if err != nil {
			return "", err
		}

		return config.BackendAzure, nil
	case "gcs":
		var config TerraformBackendGoogle
		err := hclsimple.Decode(filePath, fileContent, nil, &config)
		if err != nil {
			return "", err
		}

		return config.BackendGCS, nil
	default:
		return "Not yet supported", nil
	}
}
