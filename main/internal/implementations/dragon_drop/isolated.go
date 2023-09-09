package dragondrop

import (
	"context"

	"github.com/dragondrop-cloud/cloud-concierge/main/internal/interfaces"
	log "github.com/sirupsen/logrus"
)

// IsolatedDragonDrop is a struct that implements the DragonDrop interface for the purpose
// of end-to-end testing.
type IsolatedDragonDrop struct {
}

// NewIsolatedDragonDrop creates an instance of IsolatedDragonDrop.
func NewIsolatedDragonDrop() interfaces.DragonDrop {
	return &IsolatedDragonDrop{}
}

// CheckLoggerAndToken Check that logging is successful, and implicitly check that the jobToken env variable is valid.
func (d *IsolatedDragonDrop) CheckLoggerAndToken(_ context.Context) error {
	return nil
}

// InformStarted Informs to DragonDropAPI when job is Starting
func (d *IsolatedDragonDrop) InformStarted(_ context.Context) error {
	return nil
}

// AuthorizeJob Check with DragonDropAPI for valid auth of the current OSS job
func (d *IsolatedDragonDrop) AuthorizeJob(_ context.Context) error {
	log.Debug("Authorizing job")
	return nil
}

// AuthorizeManagedJob Check with DragonDropAPI for valid auth of the current managed job
func (d *IsolatedDragonDrop) AuthorizeManagedJob(_ context.Context) (string, error) {
	log.Debug("Authorizing job")
	return "", nil
}

// InformComplete Informs to DragonDropAPI when job is Complete
func (d *IsolatedDragonDrop) InformComplete(_ context.Context) error {
	return nil
}

// InformRepositoryCloned Informs to DragonDropAPI when job cloned the repository
func (d *IsolatedDragonDrop) InformRepositoryCloned(_ context.Context) error {
	return nil
}

// InformCloudEnvironmentScanned Informs to DragonDropAPI when job scanned the cloud environment
func (d *IsolatedDragonDrop) InformCloudEnvironmentScanned(_ context.Context) error {
	return nil
}

// InformCloudActorIdentification Informs to DragonDropAPI when job is identifying the cloud actors
func (d *IsolatedDragonDrop) InformCloudActorIdentification(_ context.Context) error {
	return nil
}

// InformCostEstimation Informs to DragonDropAPI when job is estimating costs
func (d *IsolatedDragonDrop) InformCostEstimation(_ context.Context) error {
	return nil
}

// InformSecurityScan Informs to DragonDropAPI when job is assessing cloud security
func (d *IsolatedDragonDrop) InformSecurityScan(_ context.Context) error {
	return nil
}

// InformCloudResourcesMappedToStateFile Informs to DragonDropAPI when job is mapped the resources to state file
func (d *IsolatedDragonDrop) InformCloudResourcesMappedToStateFile(_ context.Context) error {
	return nil
}

// InformNoResourcesFound Informs to DragonDropAPI when job is InformNoResourcesFound
func (d *IsolatedDragonDrop) InformNoResourcesFound(_ context.Context) error {
	return nil
}

// PostLogAlert sends alert log to the dragondrop API.
func (d *IsolatedDragonDrop) PostLogAlert(_ context.Context, _ string) {}

// PostLog sends log to the dragondrop API.
func (d *IsolatedDragonDrop) PostLog(_ context.Context, _ string) {}

// PutJobPullRequestURL sends the job url to the dragondrop API
func (d *IsolatedDragonDrop) PutJobPullRequestURL(_ context.Context, _ string) error {
	return nil
}

// PostNLPEngine posts to the NLPEngine for calculating a mapping between uncontrolled cloud resources and
// the appropriate state file.
func (d *IsolatedDragonDrop) PostNLPEngine(_ context.Context) error {
	return nil
}

// SendCloudPerchData posts anonymized cloud footprint visualization data for managed cloud-concierge jobs.
func (d *IsolatedDragonDrop) SendCloudPerchData(_ context.Context) error {
	return nil
}
