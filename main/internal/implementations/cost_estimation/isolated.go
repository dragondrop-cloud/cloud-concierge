package costEstimation

import (
	"context"

	"github.com/dragondrop-cloud/cloud-concierge/main/internal/interfaces"
)

// IsolatedCostEstimator is a struct that implements interfaces.CostEstimation for the purpose
// of end-to-end testing.
type IsolatedCostEstimator struct {
}

// NewIsolatedCostEstimator creates an instance of IsolatedCostEstimator
func NewIsolatedCostEstimator() interfaces.CostEstimation {
	return &IsolatedCostEstimator{}
}

// Execute creates structured cost estimation data for the current identified/scanned
// cloud resources.
func (ice *IsolatedCostEstimator) Execute(ctx context.Context) error {
	return nil
}
