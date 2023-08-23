package terraformWorkspace

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	terraformValueObjects "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_value_objects"
	"github.com/dragondrop-cloud/cloud-concierge/main/internal/interfaces"
)

type S3Backend struct {
	// config contains the variables that determine the specific behavior of the ContainerBackendConfig struct
	config TfStackConfig

	// dragonDrop is an implementation of the interfaces.dragonDrop interface for communicating with the
	// dragondrop API.
	dragonDrop interfaces.DragonDrop

	// s3Client is the client used to send s3 requests
	s3Client *s3.S3

	// workspaceToBackendDetails is a map of Terraform workspace names to their respective backend details.
	workspaceToBackendDetails map[string]interface{}
}

// NewS3Backend creates a new instance of the TerraformCloud struct.
func NewS3Backend(ctx context.Context, config TfStackConfig, dragonDrop interfaces.DragonDrop) interfaces.TerraformWorkspace {
	dragonDrop.PostLog(ctx, "Created TFWorkspace client.")

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(config.Region),
	}))
	s3Client := s3.New(sess)

	return &S3Backend{config: config, dragonDrop: dragonDrop, s3Client: s3Client}
}

// FindTerraformWorkspaces returns a map of TerraformCloudFile workspace names to their respective directories.
func (s *S3Backend) FindTerraformWorkspaces(ctx context.Context) (map[string]string, error) {
	workspaces, workspaceToBackendDetails, err := findTerraformWorkspaces(ctx, s.dragonDrop, s.config.WorkspaceDirectories, "s3")
	if err != nil {
		return nil, err
	}
	s.workspaceToBackendDetails = workspaceToBackendDetails

	return workspaces, err
}

// DownloadWorkspaceState downloads from the remote S3 backend the latest state file
// for each "workspace".
func (s *S3Backend) DownloadWorkspaceState(ctx context.Context, workspaceToDirectory map[string]string) error {
	s.dragonDrop.PostLog(ctx, "Beginning download of state files to local memory.")

	for workspaceName := range workspaceToDirectory {
		err := s.getWorkspaceStateByTestingAllS3Credentials(ctx, workspaceName)
		if err == nil {
			break
		}
	}

	s.dragonDrop.PostLog(ctx, "Done with download of state files to local memory.")
	return nil
}

// AWSCredentials is a struct that contains the AWS credentials for a specific workspace.
type AWSCredentials struct {
	AWSAccessKeyID     string `json:"awsAccessKeyID"`
	AWSSecretKeyAccess string `json:"awsSecretAccessKey"`
	Token              string `json:"token"`
}

// configureS3Client configures the S3 client to use the correct credentials that have read-access for the specified storage bucket.
func (s *S3Backend) configureS3Client(credential terraformValueObjects.Credential) error {
	awsCredentials := new(AWSCredentials)
	err := json.Unmarshal([]byte(credential), &awsCredentials)
	if err != nil {
		return err
	}

	staticCredentials := credentials.NewStaticCredentials(awsCredentials.AWSAccessKeyID, awsCredentials.AWSSecretKeyAccess, awsCredentials.Token)
	_, err = staticCredentials.Get()
	if err != nil {
		return err
	}

	cfg := aws.NewConfig().WithRegion(s.config.Region).WithCredentials(staticCredentials)
	newSession, err := session.NewSession(cfg)
	if err != nil {
		return err
	}

	s.s3Client = s3.New(newSession)
	return nil
}

// getWorkspaceStateByTestingAllS3Credentials downloads from the remote S3 backend a single "workspace"'s latest
// state file testing all the s3 credentials.
func (s *S3Backend) getWorkspaceStateByTestingAllS3Credentials(ctx context.Context, workspaceName string) error {
	err := s.configureS3Client(s.config.CloudCredential)
	if err != nil {
		return fmt.Errorf("[s.configureS3Client]%w", err)
	}

	stateFileName := fmt.Sprintf("%v.json", workspaceName)

	fileOutPath := fmt.Sprintf("state_files/%v", stateFileName)

	outFile, err := os.Create(fileOutPath)
	if err != nil {
		return fmt.Errorf("[get_workspace_state][error creating file]%w", err)
	}

	s3BackendDetails := s.workspaceToBackendDetails[workspaceName].(S3BackendBlock)
	downloadInput := &s3.GetObjectInput{
		Bucket: aws.String(s3BackendDetails.Bucket),
		Key:    aws.String(stateFileName),
	}

	_, err = s.s3Client.GetObject(downloadInput)
	if err != nil {
		err = outFileCloser(outFile)
		if err != nil {
			return fmt.Errorf("[s.s3Client.GetObject][outFileCloser]%v", err)
		}
		return fmt.Errorf("[s.s3Client.GetObject]%w", err)
	}

	err = outFileCloser(outFile)
	if err != nil {
		return fmt.Errorf("[outFileCloser]%v", err)
	}

	return nil
}
