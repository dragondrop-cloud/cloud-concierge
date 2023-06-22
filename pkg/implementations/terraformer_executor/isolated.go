package terraformerExecutor

import (
	"context"

	"github.com/dragondrop-cloud/driftmitigation/interfaces"
)

// IsolatedTerraformerExecutor  is a struct that implements the interfaces.TerraformerExecutor interface for
// the purpose of end-to-end testing.
type IsolatedTerraformerExecutor struct {
}

// NewIsolatedTerraformerExecutor returns a new instance of IsolatedTerraformerExecutor
func NewIsolatedTerraformerExecutor() interfaces.TerraformerExecutor {
	return &IsolatedTerraformerExecutor{}
}

// Execute runs the workflow needed to capture the current state of an
// external cloud environment via the terraformer package.
func (v *IsolatedTerraformerExecutor) Execute(ctx context.Context) error {
	return nil
}
