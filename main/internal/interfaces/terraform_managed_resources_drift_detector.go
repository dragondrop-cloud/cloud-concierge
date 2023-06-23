package interfaces

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// TerraformManagedResourcesDriftDetector is an interface that represents
// the methods required to detect drift in resources managed by Terraform.
type TerraformManagedResourcesDriftDetector interface {
	Execute(ctx context.Context, workspaceToDirectory map[string]string) (bool, error)
}

// TerraformManagedResourcesDriftDetectorMock is a mock implementation of
// the TerraformManagedResourcesDriftDetector interface, suitable for use in
// unit testing.
type TerraformManagedResourcesDriftDetectorMock struct {
	mock.Mock
}

// Execute is a mock implementation of the Execute method in
// the TerraformManagedResourcesDriftDetector interface. It simulates the
// detection of drift in managed resources by returning an error if specified in testing time.
func (m *TerraformManagedResourcesDriftDetectorMock) Execute(ctx context.Context, workspaceToDirectory map[string]string) (bool, error) {
	args := m.Called(ctx, workspaceToDirectory)
	return args.Bool(0), args.Error(1)
}
