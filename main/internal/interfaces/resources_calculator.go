package interfaces

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// ResourcesCalculator is the interface for determining the workspace/state file
// where different cloud resources should live.
type ResourcesCalculator interface {

	// Execute calculates the association between resources and a state file.
	Execute(ctx context.Context, workspaceToDirectory map[string]string) error
}

// ResourcesCalculatorMock implements the ResourcesCalculator interface for testing purposes.
type ResourcesCalculatorMock struct {
	mock.Mock
}

// Execute calculates the association between resources and a state file.
func (m *ResourcesCalculatorMock) Execute(ctx context.Context, workspaceToDirectory map[string]string) error {
	args := m.Called()
	return args.Error(0)
}
