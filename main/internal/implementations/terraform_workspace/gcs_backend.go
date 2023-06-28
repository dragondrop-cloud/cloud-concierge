package terraformWorkspace

import (
	"context"
	"fmt"
	"io"
	"os"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"

	"github.com/dragondrop-cloud/cloud-concierge/main/internal/interfaces"
)

// GCSBackend is an implementation of the interfaces.TerraformWorkspace interface that uses GCS as a backend.
type GCSBackend struct {
	// config is the configuration for the Azure Blob Storage backend.
	config ContainerBackendConfig

	// dragonDrop is the DragonDrop interface that is used to communicate with the DragonDrop API.
	dragonDrop interfaces.DragonDrop

	// workspaceToBackendDetails is a map of Terraform workspace names to their respective backend details.
	workspaceToBackendDetails map[string]interface{}
}

// FindTerraformWorkspaces returns a map of Terraform workspace names to their respective directories.
func (b *GCSBackend) FindTerraformWorkspaces(ctx context.Context) (map[string]string, error) {
	workspaces, workspaceToBackendDetails, err := findTerraformWorkspaces(ctx, b.dragonDrop, b.config.WorkspaceDirectories, "gcs")
	if err != nil {
		return nil, err
	}
	b.workspaceToBackendDetails = workspaceToBackendDetails

	return workspaces, err
}

// NewGCSBackend creates a new GCSBackend instance.
func NewGCSBackend(ctx context.Context, config ContainerBackendConfig, dragonDrop interfaces.DragonDrop) interfaces.TerraformWorkspace {
	dragonDrop.PostLog(ctx, "Created TFWorkspace client.")

	return &GCSBackend{config: config, dragonDrop: dragonDrop}
}

// DownloadWorkspaceState downloads from the remote Azure Blob Storage backend the latest state file
func (b *GCSBackend) DownloadWorkspaceState(ctx context.Context, workspaceToDirectory map[string]string) error {
	b.dragonDrop.PostLog(ctx, "Beginning download of state files to local memory.")

	for workspaceName := range workspaceToDirectory {
		err := b.getWorkspaceStateByTestingAllGoogleCredentials(ctx, workspaceName)
		if err == nil {
			break
		}
	}

	b.dragonDrop.PostLog(ctx, "Done with download of state files to local memory.")
	return nil
}

// getWorkspaceStateByTestingAllGoogleCredentials attempts to download the state file for the given workspace using all
func (b *GCSBackend) getWorkspaceStateByTestingAllGoogleCredentials(ctx context.Context, workspaceName string) error {
	for _, credential := range b.config.DivisionCloudCredentials {
		stateFileName := fmt.Sprintf("%v.json", workspaceName)

		fileOutPath := fmt.Sprintf("state_files/%v", stateFileName)

		outFile, err := os.Create(fileOutPath)
		if err != nil {
			continue
		}

		client, err := storage.NewClient(ctx, option.WithCredentialsJSON([]byte(credential)))
		if err != nil {
			continue
		}

		bucket := client.Bucket(b.config.ContainerName)
		rc, err := bucket.Object(stateFileName).NewReader(ctx)
		if err != nil {
			outFile.Close()
			continue
		}
		defer rc.Close()

		if _, err = io.Copy(outFile, rc); err != nil {
			outFile.Close()
			continue
		}

		outFile.Close()
	}

	return nil
}
