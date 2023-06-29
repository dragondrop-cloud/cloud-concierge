package terraformWorkspace

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTerraformCloud_cleanDirectoryName(t *testing.T) {
	type args struct {
		directory string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "happy path",
			args: args{
				directory: "google-backend-api-dev",
			},
			want: "google-backend-api-dev",
		},
		{
			name: "happy path with trailing slash",
			args: args{
				directory: "google-backend-api-dev/",
			},
			want: "google-backend-api-dev",
		},
		{
			name: "happy path with trailing slash and whitespace",
			args: args{
				directory: "google-backend-api-dev/ ",
			},
			want: "google-backend-api-dev",
		},
		{
			name: "happy path with trailing at beginning and end",
			args: args{
				directory: "/google-backend-api-dev/",
			},
			want: "google-backend-api-dev",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, cleanDirectoryName(tt.args.directory), "cleanDirectoryName(%v)", tt.args.directory)
		})
	}
}

func TestExtractBackendDetails_GCS(t *testing.T) {
	// Given
	ctx := context.Background()

	fileContent := []byte(`
			terraform {
			  backend "gcs" {
				bucket  = "my-gcs-terraform-backend-bucket"
				prefix  = "terraform/state"
			  }
			}
		`)

	// When
	backendType := "gcs"
	backend, err := extractBackendDetails(ctx, fileContent, backendType)

	// Then
	require.NoError(t, err, "Unexpected error: %v", err)

	gcsBackend, ok := backend.(GCSBackendBlock)

	assert.True(t, ok, "Failed to cast backendDetails to BackendGCS")
	assert.NotEmpty(t, gcsBackend.Bucket, "Failed to extract backend bucket for GCS")
	assert.NotEmpty(t, gcsBackend.Prefix, "Failed to extract backend prefix for GCS")

	require.Equal(t, "my-gcs-terraform-backend-bucket", gcsBackend.Bucket, "Failed to extract backend bucket for GCS")
	require.Equal(t, "terraform/state", gcsBackend.Prefix, "Failed to extract backend prefix for GCS")
}

func TestExtractBackendDetails_S3(t *testing.T) {
	// Given
	ctx := context.Background()

	fileContent := []byte(`
			terraform {
			  backend "s3" {
				bucket = "state-management-bucket"
				key    = "files/test.tfstate"
				region = "us-east-1"
			  }
			}
		`)

	// When
	backendType := "s3"
	backend, err := extractBackendDetails(ctx, fileContent, backendType)

	// Then
	require.NoError(t, err, "Unexpected error: %v", err)

	s3BackendBlock, ok := backend.(S3BackendBlock)

	assert.True(t, ok, "Failed to cast backendDetails to BackendGCS")
	assert.NotEmpty(t, s3BackendBlock.Bucket, "Failed to extract backend bucket for S3")
	assert.NotEmpty(t, s3BackendBlock.Key, "Failed to extract backend key for S3")
	assert.NotEmpty(t, s3BackendBlock.Region, "Failed to extract backend region for S3")

	require.Equal(t, "state-management-bucket", s3BackendBlock.Bucket)
	require.Equal(t, "files/test.tfstate", s3BackendBlock.Key)
	require.Equal(t, "us-east-1", s3BackendBlock.Region)
}

func TestExtractBackendDetails_Azurerm(t *testing.T) {
	// Given
	ctx := context.Background()

	fileContent := []byte(`
			terraform {
			  backend "azurerm" {
				resource_group_name  = "dragondrop-dev"
				storage_account_name = "dragondropstoragedev"
				container_name       = "state-files-container"
				key                  = "test-state.terraform.tfstate"
			  }
			}
		`)

	// When
	backendType := "azurerm"
	backend, err := extractBackendDetails(ctx, fileContent, backendType)

	// Then
	require.NoError(t, err, "Unexpected error: %v", err)

	azureBackendBlock, ok := backend.(AzureBackendBlock)

	assert.True(t, ok, "Failed to cast backendDetails to BackendGCS")
	assert.NotEmpty(t, azureBackendBlock.ResourceGroupName, "Failed to extract backend resource group name for azurerm")
	assert.NotEmpty(t, azureBackendBlock.StorageAccountName, "Failed to extract backend storge account name for azurerm")
	assert.NotEmpty(t, azureBackendBlock.ContainerName, "Failed to extract backend container name for azurerm")
	assert.NotEmpty(t, azureBackendBlock.Key, "Failed to extract backend key for azurerm")

	require.Equal(t, "dragondrop-dev", azureBackendBlock.ResourceGroupName, "Failed to extract backend resource group name for azurerm")
	require.Equal(t, "dragondropstoragedev", azureBackendBlock.StorageAccountName, "Failed to extract backend storge account name for azurerm")
	require.Equal(t, "state-files-container", azureBackendBlock.ContainerName, "Failed to extract backend container name for azurerm")
	require.Equal(t, "test-state.terraform.tfstate", azureBackendBlock.Key, "Failed to extract backend key for azurerm")
}

func TestExtractTFCloudWorkspaceNameIfExists(t *testing.T) {
	// Given
	ctx := context.Background()

	tests := []struct {
		name        string
		fileContent []byte
		want        string
		wantErr     bool
	}{
		{
			name: "Valid case",
			fileContent: []byte(`
terraform {
	cloud {
		workspaces {
			name = "test" 
		}
	}
}
`),
			want:    "test",
			wantErr: false,
		},
		{
			name: "Valid case with full block specification",
			fileContent: []byte(`
terraform {
  required_version = "~> 1.2.6"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "4.59.0"
    }
  }

  cloud {
    organization = "my-org"

    workspaces {
      name = "workspace"
    }
  }
}`),
			want:    "workspace",
			wantErr: false,
		},
		{
			name: "Valid case with full file specification",
			fileContent: []byte(`
provider "aws" {
  region = var.region
}

terraform {
  required_version = "~> 1.2.6"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "4.59.0"
    }
  }

  cloud {
    organization = "my-org"

    workspaces {
      name = "workspace"
    }
  }
}

provider "tfe" {
  token = var.tfe_token
}`),
			want:    "workspace",
			wantErr: false,
		},
		{
			name:        "Invalid case",
			fileContent: []byte(`invalid`),
			want:        "",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// When
			got, err := extractTFCloudWorkspaceNameIfExists(ctx, tt.fileContent)

			// Then
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
