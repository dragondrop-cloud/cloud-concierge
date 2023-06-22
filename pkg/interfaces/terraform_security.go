package interfaces

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// TerraformSecurity is an interface to execute a scanning with tfsec or mocking the files
type TerraformSecurity interface {
	ExecuteScan(ctx context.Context) error
}

// TerraformSecurityMock is a mock for testing purposes that implements the TerraformSecurity interface
type TerraformSecurityMock struct {
	mock.Mock
}

// ExecuteScan is called from the main job flow to execute the tfsec command and save the output
// to show to the user in the PR
func (m *TerraformSecurityMock) ExecuteScan(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}
