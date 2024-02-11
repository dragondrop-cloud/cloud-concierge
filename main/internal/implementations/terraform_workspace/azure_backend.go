package terraformworkspace

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/sirupsen/logrus"

	terraformValueObjects "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_value_objects"
	"github.com/dragondrop-cloud/cloud-concierge/main/internal/interfaces"
)

// AzureBlobBackend is an implementation of the interfaces.TerraformWorkspace interface that uses Azure Blob Storage as the backend.
type AzureBlobBackend struct {
	// config is the configuration for the Azure Blob Storage backend.
	config TfStackConfig

	// workspaceToBackendDetails is a map of Terraform workspace names to their respective backend details.
	workspaceToBackendDetails map[string]interface{}
}

// FindTerraformWorkspaces returns a map of Terraform workspace names to their respective directories.
func (b *AzureBlobBackend) FindTerraformWorkspaces(ctx context.Context) (map[string]string, error) {
	logrus.Debugf("[Azure Terraform workspace] Finding Terraform workspaces in %v", b.config.WorkspaceDirectories)
	workspaces, workspaceToBackendDetails, err := findTerraformWorkspaces(ctx, b.config.WorkspaceDirectories, "azurerm")
	if err != nil {
		return nil, err
	}
	b.workspaceToBackendDetails = workspaceToBackendDetails

	return workspaces, err
}

// NewAzurermBlobBackend creates a new AzureBlobBackend instance.
func NewAzurermBlobBackend(ctx context.Context, config TfStackConfig) interfaces.TerraformWorkspace {
	return &AzureBlobBackend{config: config}
}

// DownloadWorkspaceState downloads from the remote Azure Blob Storage backend the latest state file
func (b *AzureBlobBackend) DownloadWorkspaceState(ctx context.Context, workspaceToDirectory map[string]string) error {
	logrus.Debugf("[Azure Terraform workspace] Downloading workspace state files for %v", workspaceToDirectory)

	for workspaceName := range workspaceToDirectory {
		err := b.getWorkspaceStateFromAzureCredentials(ctx, workspaceName)
		if err != nil {
			return err
		}
	}

	return nil
}

// getWorkspaceStateFromAzureCredentials downloads the state file for the given workspace from the Azure Blob Storage backend.
func (b *AzureBlobBackend) getWorkspaceStateFromAzureCredentials(ctx context.Context, workspaceName string) error {
	serviceURL, err := b.configureAzureBlobURL(ctx, b.config.CloudCredential, b.workspaceToBackendDetails[workspaceName].(AzureBackendBlock))
	if err != nil {
		return fmt.Errorf("[b.configureAzureBlobURL]%v", err)
	}

	stateFileName := fmt.Sprintf("%v.json", workspaceName)

	fileOutPath := fmt.Sprintf("state_files/%v", stateFileName)

	outFile, err := os.Create(fileOutPath)
	if err != nil {
		return fmt.Errorf("[os.Create]%v", err)
	}

	azureBackendDetails := b.workspaceToBackendDetails[workspaceName].(AzureBackendBlock)
	blobURL := serviceURL.NewContainerURL(azureBackendDetails.ContainerName).NewBlobURL(stateFileName)

	err = azblob.DownloadBlobToFile(ctx, blobURL, 0, azblob.CountToEnd, outFile, azblob.DownloadFromBlobOptions{})
	if err != nil {
		err = outFileCloser(outFile)
		if err != nil {
			return fmt.Errorf("[azblob.DownloadBlobToFile][outFileCloser]%v", err)
		}
		return fmt.Errorf("[azblob.DownloadBlobToFile]%v", err)
	}

	return outFileCloser(outFile)
}

// AzureCredentials is a struct that holds the credentials for an Azure Blob Storage backend.
type AzureCredentials struct {
	AzureStorageAccountKey string `json:"azure_storage_account_key"`
}

// configureAzureBlobURL configures the Azure Blob Storage URL.
func (b *AzureBlobBackend) configureAzureBlobURL(_ context.Context, credential terraformValueObjects.Credential, backendAzure AzureBackendBlock) (azblob.ServiceURL, error) {
	azureCredentials := new(AzureCredentials)
	err := json.Unmarshal([]byte(credential), &azureCredentials)
	if err != nil {
		return azblob.ServiceURL{}, err
	}

	sharedCredential, err := azblob.NewSharedKeyCredential(backendAzure.StorageAccountName, azureCredentials.AzureStorageAccountKey)
	if err != nil {
		return azblob.ServiceURL{}, err
	}

	p := azblob.NewPipeline(sharedCredential, azblob.PipelineOptions{})
	URL, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net", backendAzure.StorageAccountName))
	return azblob.NewServiceURL(*URL, p), nil
}
