package interfaces

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// TerraformerExecutor is an interface for getting the current state of an external cloud
// by programmatically running terraformer commands.
type TerraformerExecutor interface {
	// Execute runs the workflow needed to capture the current state of an
	// external cloud environment via the terraformer package.
	Execute(ctx context.Context) error
}

// TerraformerExecutorMock is a struct that implements the TerraformerExecutor interface for
// unit testing purposes.
type TerraformerExecutorMock struct {
	mock.Mock
}

// Execute runs the workflow needed to capture the current state of an
// external cloud environment via the terraformer package.
func (m *TerraformerExecutorMock) Execute(_ context.Context) error {
	args := m.Called()
	return args.Error(0)
}
