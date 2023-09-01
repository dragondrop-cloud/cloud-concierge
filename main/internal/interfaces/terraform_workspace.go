package interfaces

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// TerraformWorkspace is an interface for pulling in data from remote Terraform state files.
type TerraformWorkspace interface {
	FindTerraformWorkspaces(ctx context.Context) (map[string]string, error)

	// DownloadWorkspaceState downloads from the remote Terraform backend the latest state file
	// for each "workspace".
	DownloadWorkspaceState(ctx context.Context, WorkspaceToDirectory map[string]string) error
}

// TerraformWorkspaceMock implements the TerraformWorkspace interface for unit testing purposes.
type TerraformWorkspaceMock struct {
	mock.Mock
}

// FindTerraformWorkspaces returns a map of Terraform workspace names to their respective directories.
func (m *TerraformWorkspaceMock) FindTerraformWorkspaces(ctx context.Context) (map[string]string, error) {
	args := m.Called(ctx)
	return args.Get(0).(map[string]string), args.Error(1)
}

// DownloadWorkspaceState downloads from the remote Terraform backend the latest state file
// for each "workspace".
func (m *TerraformWorkspaceMock) DownloadWorkspaceState(_ context.Context, _ map[string]string) error {
	args := m.Called()
	return args.Error(0)
}
