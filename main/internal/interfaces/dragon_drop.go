package interfaces

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// DragonDrop is the interface for communicating with the external DragonDrop API
type DragonDrop interface {
	// CheckLoggerAndToken Check that logging is successful, and implicitly check that the jobToken env variable is valid.
	CheckLoggerAndToken(ctx context.Context) error

	// InformStarted Informs to DragonDropAPI when job is Starting
	InformStarted(ctx context.Context) error

	// AuthorizeJob Check with DragonDropAPI for valid auth of the current job when run oss
	AuthorizeJob(ctx context.Context) error

	// AuthorizeManagedJob Check with DragonDropAPI for valid auth of the current job, for a dragondop managed job
	AuthorizeManagedJob(ctx context.Context) (string, error)

	// InformRepositoryCloned Informs to DragonDropAPI when job cloned the repository
	InformRepositoryCloned(ctx context.Context) error

	// InformCloudEnvironmentScanned Informs to DragonDropAPI when job scanned the cloud environment
	InformCloudEnvironmentScanned(ctx context.Context) error

	// InformCloudActorIdentification Informs to DragonDropAPI when job is identifying the cloud actors
	InformCloudActorIdentification(ctx context.Context) error

	// InformCostEstimation Informs to DragonDropAPI when job is estimating costs
	InformCostEstimation(ctx context.Context) error

	// InformSecurityScan Informs to DragonDropAPI when job is assessing cloud security
	InformSecurityScan(ctx context.Context) error

	// InformCloudResourcesMappedToStateFile Informs to DragonDropAPI when job is mapped the resources to state file
	InformCloudResourcesMappedToStateFile(ctx context.Context) error

	// InformNoResourcesFound Informs to DragonDropAPI when job is InformNoResourcesFound
	InformNoResourcesFound(ctx context.Context) error

	// InformComplete Informs to DragonDropAPI when job is Complete
	InformComplete(ctx context.Context) error

	// PostLogAlert sends alert log to the dragondrop API.
	PostLogAlert(ctx context.Context, log string)

	// PostLog sends log to the dragondrop API.
	PostLog(ctx context.Context, log string)

	// PutJobPullRequestURL sends the job url to the dragondrop API
	PutJobPullRequestURL(ctx context.Context, prURL string) error
}

// DragonDropMock is a struct that implements the DragonDrop interface solely for the purpose
// of testing with the testify library.
type DragonDropMock struct {
	mock.Mock
}

// CheckLoggerAndToken Check that logging is successful, and implicitly check that the jobToken env variable is valid.
func (m *DragonDropMock) CheckLoggerAndToken(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// InformStarted Informs to DragonDropAPI when job is Starting
func (m *DragonDropMock) InformStarted(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// AuthorizeJob Check with DragonDropAPI for valid auth of the current OSS job
func (m *DragonDropMock) AuthorizeJob(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// AuthorizeManagedJob Check with DragonDropAPI for valid auth of the current managed job
func (m *DragonDropMock) AuthorizeManagedJob(ctx context.Context) (string, error) {
	args := m.Called(ctx)
	return args.String(0), args.Error(1)
}

// InformComplete Informs to DragonDropAPI when job is Complete
func (m *DragonDropMock) InformComplete(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// InformRepositoryCloned Informs to DragonDropAPI when job cloned the repository
func (m *DragonDropMock) InformRepositoryCloned(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// InformCloudEnvironmentScanned Informs to DragonDropAPI when job scanned the cloud environment
func (m *DragonDropMock) InformCloudEnvironmentScanned(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// InformCloudActorIdentification Informs to DragonDropAPI when job is identifying the cloud actors
func (m *DragonDropMock) InformCloudActorIdentification(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// InformCostEstimation Informs to DragonDropAPI when job is estimating costs
func (m *DragonDropMock) InformCostEstimation(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// InformSecurityScan Informs to DragonDropAPI when job is assessing cloud security
func (m *DragonDropMock) InformSecurityScan(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// InformCloudResourcesMappedToStateFile Informs to DragonDropAPI when job is mapped the resources to state file
func (m *DragonDropMock) InformCloudResourcesMappedToStateFile(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// InformNoResourcesFound Informs to DragonDropAPI when job is InformNoResourcesFound
func (m *DragonDropMock) InformNoResourcesFound(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// PostLogAlert sends alert log to the dragondrop API.
func (m *DragonDropMock) PostLogAlert(ctx context.Context, log string) {}

// PostLog sends log to the dragondrop API.
func (m *DragonDropMock) PostLog(ctx context.Context, log string) {}

// PutJobPullRequestURL sends the job url to the dragondrop API
func (m *DragonDropMock) PutJobPullRequestURL(ctx context.Context, prURL string) error {
	args := m.Called(ctx, prURL)
	return args.Error(0)
}
