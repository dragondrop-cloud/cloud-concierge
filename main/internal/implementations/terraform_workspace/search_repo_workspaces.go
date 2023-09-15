package terraformworkspace

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	log "github.com/sirupsen/logrus"

	"github.com/dragondrop-cloud/cloud-concierge/main/internal/interfaces"
)

// outFileCloser closes the outFile and returns an error if one occurred.
func outFileCloser(outFile *os.File) error {
	err := outFile.Close()
	if err != nil {
		return fmt.Errorf("[outFile.Close]%v", err)
	}
	return nil
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
func getAllTFFiles(_ context.Context, directory string) []string {
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

	log.Debugf("[get_all_tf_files][found %d tf files]", len(tfFiles))
	return tfFiles
}

// cleanDirectoryName removes any leading or trailing slashes from a directory name.
func cleanDirectoryName(directory string) string {
	directory = strings.Trim(directory, " ")
	directory = strings.Trim(directory, "/")
	return directory
}

// getWorkspaceByFile searches a given file for a Terraform workspace.
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

	log.Debugf("[get_workspace_by_file][found workspace %s in file %s]", workspace, fileName)
	return workspace, details, true
}

// extractTFCloudWorkspaceNameIfExists extracts the workspace name from a Terraform file if it exists.
func extractTFCloudWorkspaceNameIfExists(_ context.Context, fileContent []byte) (string, error) {
	inputHCLFile, hclDiag := hclwrite.ParseConfig(
		fileContent,
		"placeholder.tf",
		hcl.Pos{Line: 0, Column: 0, Byte: 0},
	)
	if hclDiag.HasErrors() {
		return "", fmt.Errorf("error parsing HCL file: %s", hclDiag.Error())
	}

	// checking to see if a Terraform Cloud workspace configuration exists, and if so, extracting the workspace name.
	terraform := inputHCLFile.Body().FirstMatchingBlock("terraform", nil)
	cloud := terraform.Body().FirstMatchingBlock("cloud", nil)
	workspaces := cloud.Body().FirstMatchingBlock("workspaces", nil)
	workspacesName := workspaces.Body().GetAttribute("name")

	workspaceNameTokens := string(workspacesName.BuildTokens(nil).Bytes())

	re := regexp.MustCompile(`"(.*)"`)
	matches := re.FindStringSubmatch(workspaceNameTokens)

	if len(matches) < 2 {
		return "", fmt.Errorf("error extracting workspace name from file")
	}

	return matches[1], nil
}

// S3BackendBlock is a struct representation of a terraform backend block for s3
type S3BackendBlock struct {
	Bucket string
	Key    string
	Region string
}

// AzureBackendBlock is a struct representation of a terraform backend block for azurerm
type AzureBackendBlock struct {
	ResourceGroupName  string
	StorageAccountName string
	ContainerName      string
	Key                string
}

// GCSBackendBlock parses the backend block for GCS
type GCSBackendBlock struct {
	Bucket string
	Prefix string
}

// extractBackendDetails extracts the backend details from a .tf file if it exists.
func extractBackendDetails(_ context.Context, fileContent []byte, backendType string) (interface{}, error) {
	switch backendType {
	case "s3":
		return extractAttributesFromBackendDetails(fileContent, "s3", []string{"bucket", "key", "region"},
			func(attributesMap map[string]string) interface{} {
				return S3BackendBlock{
					Region: attributesMap["region"],
					Key:    attributesMap["key"],
					Bucket: attributesMap["bucket"],
				}
			})
	case "azurerm":
		return extractAttributesFromBackendDetails(fileContent, "azurerm",
			[]string{"resource_group_name", "storage_account_name", "container_name", "key"},
			func(attributesMap map[string]string) interface{} {
				return AzureBackendBlock{
					ResourceGroupName:  attributesMap["resource_group_name"],
					StorageAccountName: attributesMap["storage_account_name"],
					ContainerName:      attributesMap["container_name"],
					Key:                attributesMap["key"],
				}
			},
		)
	case "gcs":
		return extractAttributesFromBackendDetails(fileContent, "gcs", []string{"bucket", "prefix"},
			func(attributesMap map[string]string) interface{} {
				return GCSBackendBlock{
					Bucket: attributesMap["bucket"],
					Prefix: attributesMap["prefix"],
				}
			})
	default:
		return "Not yet supported", nil
	}
}

// extractAttributesFromBackendDetails extracts the backend details from a .tf file if it exists.
func extractAttributesFromBackendDetails(fileContent []byte, provider string, attributes []string, deserializer func(map[string]string) interface{}) (interface{}, error) {
	inputHCLFile, hclDiag := hclwrite.ParseConfig(
		fileContent,
		"placeholder.tf",
		hcl.Pos{Line: 0, Column: 0, Byte: 0},
	)
	if hclDiag.HasErrors() {
		return "", fmt.Errorf("error parsing HCL file: %s", hclDiag.Error())
	}

	// checking to see if a Terraform Cloud workspace configuration exists, and if so, extracting the workspace name.
	terraform := inputHCLFile.Body().FirstMatchingBlock("terraform", nil)
	backend := terraform.Body().FirstMatchingBlock("backend", []string{provider})

	attributesMap := make(map[string]string)
	for _, attribute := range attributes {
		attributeExpression := backend.Body().GetAttribute(attribute)
		attributeTokenBytes := string(attributeExpression.BuildTokens(nil).Bytes())

		re := regexp.MustCompile(`"(.*)"`)
		matches := re.FindStringSubmatch(attributeTokenBytes)

		if len(matches) < 2 {
			return "", fmt.Errorf("error extracting attribute %s for %s provider", attributes, provider)
		}

		attributesMap[attribute] = matches[1]
	}

	return deserializer(attributesMap), nil
}
