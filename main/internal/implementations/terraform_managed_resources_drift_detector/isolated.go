package terraformManagedResourcesDriftDetector

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/dragondrop-cloud/cloud-concierge/main/internal/interfaces"
)

// IsolatedTerraformResourcesManagedDriftDetector is a struct that implements the TerraformManagedResourcesDriftDetector interface
// for the purpose of running end to end unit tests.
type IsolatedTerraformResourcesManagedDriftDetector struct {
}

// NewIsolatedDriftDetector returns an instance of IsolatedTerraformResourcesManagedDriftDetector
func NewIsolatedDriftDetector() interfaces.TerraformManagedResourcesDriftDetector {
	return &IsolatedTerraformResourcesManagedDriftDetector{}
}

// Execute detects drift in Terraform-managed resources.
func (d *IsolatedTerraformResourcesManagedDriftDetector) Execute(ctx context.Context, workspaceToDirectory map[string]string) (bool, error) {
	log.Debug("Executing drift detector")
	return false, nil
}
