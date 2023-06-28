package terraformWorkspace

import (
	"context"
	"os"
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

func TestExtractBackendDetails(t *testing.T) {
	// Given
	ctx := context.Background()

	filePath := "test.tfstate"
	fileContent := []byte(`
			terraform {
			  backend "gcs" {
				bucket  = "my-gcs-terraform-backend-bucket"
				prefix  = "terraform/state"
			  }
			}
		`)
	err := os.WriteFile(filePath, fileContent, 0600)
	require.NoError(t, err)

	// When
	backendType := "gcs"
	backendDetails, err := extractBackendDetails(ctx, filePath, fileContent, backendType)

	// Then
	require.NoError(t, err, "Unexpected error: %v", err)

	gcsBackend, ok := backendDetails.(BackendGCS)

	assert.True(t, ok, "Failed to cast backendDetails to BackendGCS")
	assert.NotEmpty(t, gcsBackend.Bucket, "Failed to extract backend bucket for GCS")
	assert.NotEmpty(t, gcsBackend.Prefix, "Failed to extract backend prefix for GCS")
}

func TestExtractTFCloudWorkspaceNameIfExists(t *testing.T) {
	// Given
	ctx := context.Background()

	tests := []struct {
		name        string
		filePath    string
		fileContent []byte
		want        string
		wantErr     bool
	}{
		{
			name:        "Valid case",
			filePath:    "valid.tf",
			fileContent: []byte(`cloud { workspaces { name = "test" } }`),
			want:        "test",
			wantErr:     false,
		},
		{
			name:        "Invalid case",
			filePath:    "invalid.tf",
			fileContent: []byte(`invalid`),
			want:        "",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// When
			got, err := extractTFCloudWorkspaceNameIfExists(ctx, tt.filePath, tt.fileContent)

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
