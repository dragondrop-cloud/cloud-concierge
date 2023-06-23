package interfaces

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// CostEstimation is an interface for estimating the cost of Terraform resources.
type CostEstimation interface {

	// Execute creates structured cost estimation data for the current identified/scanned
	// cloud resources.
	Execute(ctx context.Context) error
}

// CostEstimationMock implements the CostEstimation interface for testing purposes.
type CostEstimationMock struct {
	mock.Mock
}

// Execute creates structured cost estimation data for the current identified/scanned
// cloud resources.
func (m *CostEstimationMock) Execute(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}
