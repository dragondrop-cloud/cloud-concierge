package terraformworkspace

import (
	"context"
	"fmt"
	"io"
	"os"

	"cloud.google.com/go/storage"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/option"

	"github.com/dragondrop-cloud/cloud-concierge/main/internal/interfaces"
)

// GCSBackend is an implementation of the interfaces.TerraformWorkspace interface that uses GCS as a backend.
type GCSBackend struct {
	// config is the configuration for the Azure Blob Storage backend.
	config TfStackConfig

	// dragonDrop is the DragonDrop interface that is used to communicate with the DragonDrop API.
	dragonDrop interfaces.DragonDrop

	// workspaceToBackendDetails is a map of Terraform workspace names to their respective backend details.
	workspaceToBackendDetails map[string]interface{}
}

// FindTerraformWorkspaces returns a map of Terraform workspace names to their respective directories.
func (b *GCSBackend) FindTerraformWorkspaces(ctx context.Context) (map[string]string, error) {
	logrus.Debugf("[GCS Terraform workspace] Finding Terraform workspaces in %v", b.config.WorkspaceDirectories)

	workspaces, workspaceToBackendDetails, err := findTerraformWorkspaces(ctx, b.dragonDrop, b.config.WorkspaceDirectories, "gcs")
	if err != nil {
		return nil, err
	}
	b.workspaceToBackendDetails = workspaceToBackendDetails

	return workspaces, err
}

// NewGCSBackend creates a new GCSBackend instance.
func NewGCSBackend(ctx context.Context, config TfStackConfig, dragonDrop interfaces.DragonDrop) interfaces.TerraformWorkspace {
	dragonDrop.PostLog(ctx, "Created TFWorkspace client.")

	return &GCSBackend{config: config, dragonDrop: dragonDrop}
}

// DownloadWorkspaceState downloads from the remote Azure Blob Storage backend the latest state file
func (b *GCSBackend) DownloadWorkspaceState(ctx context.Context, workspaceToDirectory map[string]string) error {
	logrus.Debugf("[GCS Terraform workspace] Downloading workspace state files for %v", workspaceToDirectory)
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

	stateFileName := fmt.Sprintf("%v.json", workspaceName)

	fileOutPath := fmt.Sprintf("state_files/%v", stateFileName)

	outFile, err := os.Create(fileOutPath)
	if err != nil {
		return fmt.Errorf("[os.Create] %v", err)
	}

	client, err := storage.NewClient(ctx, option.WithCredentialsJSON([]byte(b.config.CloudCredential)))
	if err != nil {
		return fmt.Errorf("[storage.NewClient] %v", err)
	}

	gcsBackendDetails := b.workspaceToBackendDetails[workspaceName].(GCSBackendBlock)
	bucket := client.Bucket(gcsBackendDetails.Bucket)
	rc, err := bucket.Object(stateFileName).NewReader(ctx)
	if err != nil {
		err = outFileCloser(outFile)
		if err != nil {
			return fmt.Errorf("[bucket.Object().NewReader][outFileCloser]%v", err)
		}
		return fmt.Errorf("[bucket.Object().NewReader] %v", err)
	}
	defer rc.Close()

	if _, err = io.Copy(outFile, rc); err != nil {
		err = outFileCloser(outFile)
		if err != nil {
			return fmt.Errorf("[io.Copy][outFileCloser]%v", err)
		}
	}

	err = outFileCloser(outFile)
	if err != nil {
		return fmt.Errorf("[outFileCloser] %v", err)
	}

	return nil
}
