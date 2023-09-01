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
func (d *IsolatedDragonDrop) CheckLoggerAndToken(ctx context.Context) error {
	return nil
}

// InformStarted Informs to DragonDropAPI when job is Starting
func (d *IsolatedDragonDrop) InformStarted(ctx context.Context) error {
	return nil
}

// AuthorizeJob Check with DragonDropAPI for valid auth of the current OSS job
func (d *IsolatedDragonDrop) AuthorizeJob(ctx context.Context) error {
	log.Debug("Authorizing job")
	return nil
}

// AuthorizeManagedJob Check with DragonDropAPI for valid auth of the current managed job
func (d *IsolatedDragonDrop) AuthorizeManagedJob(ctx context.Context) (string, error) {
	log.Debug("Authorizing job")
	return "", nil
}

// InformComplete Informs to DragonDropAPI when job is Complete
func (d *IsolatedDragonDrop) InformComplete(ctx context.Context) error {
	return nil
}

// InformRepositoryCloned Informs to DragonDropAPI when job cloned the repository
func (d *IsolatedDragonDrop) InformRepositoryCloned(ctx context.Context) error {
	return nil
}

// InformCloudEnvironmentScanned Informs to DragonDropAPI when job scanned the cloud environment
func (d *IsolatedDragonDrop) InformCloudEnvironmentScanned(ctx context.Context) error {
	return nil
}

// InformCloudActorIdentification Informs to DragonDropAPI when job is identifying the cloud actors
func (d *IsolatedDragonDrop) InformCloudActorIdentification(ctx context.Context) error {
	return nil
}

// InformCostEstimation Informs to DragonDropAPI when job is estimating costs
func (d *IsolatedDragonDrop) InformCostEstimation(ctx context.Context) error {
	return nil
}

// InformSecurityScan Informs to DragonDropAPI when job is assessing cloud security
func (d *IsolatedDragonDrop) InformSecurityScan(ctx context.Context) error {
	return nil
}

// InformCloudResourcesMappedToStateFile Informs to DragonDropAPI when job is mapped the resources to state file
func (d *IsolatedDragonDrop) InformCloudResourcesMappedToStateFile(ctx context.Context) error {
	return nil
}

// InformNoResourcesFound Informs to DragonDropAPI when job is InformNoResourcesFound
func (d *IsolatedDragonDrop) InformNoResourcesFound(ctx context.Context) error {
	return nil
}

// PostLogAlert sends alert log to the dragondrop API.
func (d *IsolatedDragonDrop) PostLogAlert(ctx context.Context, log string) {}

// PostLog sends log to the dragondrop API.
func (d *IsolatedDragonDrop) PostLog(ctx context.Context, log string) {}

// PutJobPullRequestURL sends the job url to the dragondrop API
func (d *IsolatedDragonDrop) PutJobPullRequestURL(ctx context.Context, prURL string) error {
	return nil
}

func (d *IsolatedDragonDrop) SendCloudPerchData(ctx context.Context) error {
	return nil
}
