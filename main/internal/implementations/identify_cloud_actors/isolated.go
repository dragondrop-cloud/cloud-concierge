package identifyCloudActors

import (
	"context"

	"github.com/dragondrop-cloud/cloud-concierge/main/internal/interfaces"
	log "github.com/sirupsen/logrus"
)

// IsolatedIdentifyCloudActors is a struct that implements interfaces.IdentifyCloudActors
// for the purpose of running end to end unit tests.
type IsolatedIdentifyCloudActors struct {
}

// NewIsolatedIdentifyCloudActors returns an instance of IdentifyCloudActors
func NewIsolatedIdentifyCloudActors() interfaces.IdentifyCloudActors {
	return &IsolatedIdentifyCloudActors{}
}

// Execute calculates the association between resources and a state file.
func (c *IsolatedIdentifyCloudActors) Execute(ctx context.Context) error {
	log.Debug("Executing identify cloud actors")
	return nil
}
