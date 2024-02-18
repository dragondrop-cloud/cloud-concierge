package resourceswriter

import (
	"context"

	"github.com/dragondrop-cloud/cloud-concierge/main/internal/interfaces"
)

// IsolatedResourcesWriter is a struct that implements the interfaces.ResourcesWriter interface for
// the purpose of end-to-end testing.
type IsolatedResourcesWriter struct{}

// NewIsolatedResourcesWriter returns a new instance of IsolatedResourcesWriter.
func NewIsolatedResourcesWriter() interfaces.ResourcesWriter {
	return &IsolatedResourcesWriter{}
}

// Execute writes new resources to the relevant version control system,
// and returns a pull request url corresponding to the new changes.
func (w *IsolatedResourcesWriter) Execute(_ context.Context, _ bool, _ map[string]string) (string, error) {
	return "", nil
}
