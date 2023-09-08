package interfaces

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// ResourcesWriter is an interface for writing Terraform resource information to
// a version control system.
type ResourcesWriter interface {

	// Execute writes new resources to the relevant version control system,
	// and returns a pull request url corresponding to the new changes.
	Execute(ctx context.Context, jobName string, createDummyFile bool, workspaceToDirectory map[string]string) (string, error)
}

// ResourcesWriterMock implements the ResourcesWriter interface for testing purposes.
type ResourcesWriterMock struct {
	mock.Mock
}

// Execute writes new resources to the relevant version control system,
// and returns a pull request url corresponding to the new changes.
func (m *ResourcesWriterMock) Execute(_ context.Context, _ string, _ bool, _ map[string]string) (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}
