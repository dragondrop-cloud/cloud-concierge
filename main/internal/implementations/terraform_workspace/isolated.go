package terraformworkspace

import (
	"context"

	log "github.com/sirupsen/logrus"
)

// IsolatedTerraformWorkspace is a struct that implements the interfaces.TerraformWorkspace interface for
// the purpose of end-to-end testing.
type IsolatedTerraformWorkspace struct{}

// NewIsolatedTerraformWorkspace returns a new instance of IsolatedTerraformWorkspace
func NewIsolatedTerraformWorkspace() *IsolatedTerraformWorkspace {
	return &IsolatedTerraformWorkspace{}
}

// DownloadWorkspaceState downloads from the remote TerraformCloudFile backend the latest state file
// for each "workspace".
func (v *IsolatedTerraformWorkspace) DownloadWorkspaceState(_ context.Context, _ map[string]string) error {
	log.Debug("Downloading terraform workspace state")
	return nil
}

func (v *IsolatedTerraformWorkspace) FindTerraformWorkspaces(_ context.Context) (map[string]string, error) {
	log.Debug("finding terraform workspaces")
	return make(map[string]string), nil
}
