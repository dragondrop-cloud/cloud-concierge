package resourcescalculator

import (
	"context"

	"github.com/dragondrop-cloud/cloud-concierge/main/internal/interfaces"
	log "github.com/sirupsen/logrus"
)

// IsolatedResourcesCalculator is a struct that implements the ResourcesCalculator interface
// for the purpose of running end to end unit tests.
type IsolatedResourcesCalculator struct {
}

// NewIsolatedResourcesCalculator returns an instance of IsolatedResourcesCalculator
func NewIsolatedResourcesCalculator() interfaces.ResourcesCalculator {
	return &IsolatedResourcesCalculator{}
}

// Execute calculates the association between resources and a state file.
func (c *IsolatedResourcesCalculator) Execute(_ context.Context, _ map[string]string) error {
	log.Debug("Executing resource calculator")
	return nil
}
