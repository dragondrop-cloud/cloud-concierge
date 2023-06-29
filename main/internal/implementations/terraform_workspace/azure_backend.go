package terraformWorkspace

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	"github.com/Azure/azure-storage-blob-go/azblob"

	terraformValueObjects "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_value_objects"
	"github.com/dragondrop-cloud/cloud-concierge/main/internal/interfaces"
)

// AzureBlobBackend is an implementation of the interfaces.TerraformWorkspace interface that uses Azure Blob Storage as the backend.
type AzureBlobBackend struct {
	// config is the configuration for the Azure Blob Storage backend.
	config ContainerBackendConfig

	// dragonDrop is the DragonDrop interface that is used to communicate with the DragonDrop API.
	dragonDrop interfaces.DragonDrop

	// workspaceToBackendDetails is a map of Terraform workspace names to their respective backend details.
	workspaceToBackendDetails map[string]interface{}
}

// FindTerraformWorkspaces returns a map of Terraform workspace names to their respective directories.
func (b *AzureBlobBackend) FindTerraformWorkspaces(ctx context.Context) (map[string]string, error) {
	workspaces, workspaceToBackendDetails, err := findTerraformWorkspaces(ctx, b.dragonDrop, b.config.WorkspaceDirectories, "azurerm")
	if err != nil {
		return nil, err
	}
	b.workspaceToBackendDetails = workspaceToBackendDetails

	return workspaces, err
}

// NewAzurermBlobBackend creates a new AzureBlobBackend instance.
func NewAzurermBlobBackend(ctx context.Context, config ContainerBackendConfig, dragonDrop interfaces.DragonDrop) interfaces.TerraformWorkspace {
	dragonDrop.PostLog(ctx, "Created TFWorkspace client.")

	return &AzureBlobBackend{config: config, dragonDrop: dragonDrop}
}

// DownloadWorkspaceState downloads from the remote Azure Blob Storage backend the latest state file
func (b *AzureBlobBackend) DownloadWorkspaceState(ctx context.Context, workspaceToDirectory map[string]string) error {
	b.dragonDrop.PostLog(ctx, "Beginning download of state files to local memory.")

	for workspaceName := range workspaceToDirectory {
		err := b.getWorkspaceStateByTestingAllAzureCredentials(ctx, workspaceName)
		if err == nil {
			break
		}
	}

	b.dragonDrop.PostLog(ctx, "Done with download of state files to local memory.")
	return nil
}

// getWorkspaceStateByTestingAllAzureCredentials downloads the state file for the given workspace from the Azure Blob Storage backend.
func (b *AzureBlobBackend) getWorkspaceStateByTestingAllAzureCredentials(ctx context.Context, workspaceName string) error {
	for _, credential := range b.config.DivisionCloudCredentials {
		serviceURL, err := b.configureAzureBlobURL(ctx, credential, b.workspaceToBackendDetails[workspaceName].(AzureBackendBlock))
		if err != nil {
			continue
		}

		stateFileName := fmt.Sprintf("%v.json", workspaceName)

		fileOutPath := fmt.Sprintf("state_files/%v", stateFileName)

		outFile, err := os.Create(fileOutPath)
		if err != nil {
			continue
		}

		blobURL := serviceURL.NewContainerURL(b.config.ContainerName).NewBlobURL(stateFileName)

		err = azblob.DownloadBlobToFile(ctx, blobURL, 0, azblob.CountToEnd, outFile, azblob.DownloadFromBlobOptions{})
		if err != nil {
			outFile.Close()
			continue
		}

		outFile.Close()
	}

	return nil
}

// AzureCredentials is a struct that holds the credentials for an Azure Blob Storage backend.
type AzureCredentials struct {
	AzureStorageAccountKey string `json:"azure_storage_account_key"`
}

// configureAzureBlobURL configures the Azure Blob Storage URL.
func (b *AzureBlobBackend) configureAzureBlobURL(ctx context.Context, credential terraformValueObjects.Credential, backendAzure AzureBackendBlock) (azblob.ServiceURL, error) {
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
