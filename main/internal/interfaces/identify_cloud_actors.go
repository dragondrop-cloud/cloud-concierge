package interfaces

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// IdentifyCloudActors is an interface for identifying the cloud actors that have made changes
// to identified drifted resources.
type IdentifyCloudActors interface {
	// Execute creates structured query_param_data mapping new or drifted resources to the cloud actor (service principal or user)
	// responsible for the latest changes for that resource.
	Execute(ctx context.Context) error
}

// IdentifyCloudActorsMock implements the IdentifyCloudActors interface for testing purposes.
type IdentifyCloudActorsMock struct {
	mock.Mock
}

// Execute creates structured query_param_data mapping new or drifted resources to the cloud actor (service principal or user)
// responsible for the latest changes for that resource.
func (m *IdentifyCloudActorsMock) Execute(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}
