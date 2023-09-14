package driftdetector

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"

	terraformValueObjects "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_value_objects"
)

// ManagedResourceDriftDetectorConfig is a type that contains configuration
type ManagedResourceDriftDetectorConfig struct {
	// ResourcesWhiteList represents the list of resource names that will be exclusively considered for inclusion in the import statement.
	ResourcesWhiteList terraformValueObjects.ResourceNameList

	// ResourcesBlackList represents the list of resource names that will be excluded from consideration for inclusion in the import statement.
	ResourcesBlackList terraformValueObjects.ResourceNameList
}

// ManagedResourcesDriftDetector is a type that identifies resources
// managed by Terraform that have drifted from their expected state.
type ManagedResourcesDriftDetector struct {
	// config is the configuration for the ManagedResourcesDriftDetector
	config ManagedResourceDriftDetectorConfig
}

// NewManagedResourcesDriftDetector generated a terraformer instance from ManagedResourcesDriftDetector
func NewManagedResourcesDriftDetector(config ManagedResourceDriftDetectorConfig) *ManagedResourcesDriftDetector {
	return &ManagedResourcesDriftDetector{
		config: config,
	}
}

// Execute initiates the process of detecting drift in managed resources
// by comparing the current state of resources with their expected state.
// It takes a context as input to support cancellation and timeouts.
func (m *ManagedResourcesDriftDetector) Execute(_ context.Context, workspaceToDirectory map[string]string) (bool, error) {
	logrus.Debugf("[drift_detector] workspaceToDirectory: %v", workspaceToDirectory)

	remoteStateResources, err := m.loadAllRemoteStateFiles(workspaceToDirectory)
	if err != nil {
		return false, fmt.Errorf("[m.loadAllRemoteStateFiles]%w", err)
	}
	logrus.Debugf("[drift_detector] remoteStateResources: %v", remoteStateResources)

	terraformerStateResources, err := m.loadAllTerraformerStateFiles()
	if err != nil {
		return false, fmt.Errorf("[m.loadAllTerraformerStateFiles]%w", err)
	}
	logrus.Debugf("[drift_detector] terraformerStateResources: %v", terraformerStateResources)

	wereDeleted, err := m.identifyAndWriteDeletedResources(terraformerStateResources, remoteStateResources)
	if err != nil {
		return false, fmt.Errorf("[m.identifyAndWriteDeletedResources]%w", err)
	}
	logrus.Debugf("[drift_detector] wereDeleted: %v", wereDeleted)

	differencesFound, err := m.identifyAndWriteResourcesDifferences(terraformerStateResources, remoteStateResources)
	if err != nil {
		return false, fmt.Errorf("[m.identifyAndWriteResourcesDifferences]%w", err)
	}
	logrus.Debugf("[drift_detector] differencesFound: %v", differencesFound)

	return wereDeleted || differencesFound, nil
}

// identifyAndWriteResourcesDifferences found the resources differences and writes in the mapping file
func (m *ManagedResourcesDriftDetector) identifyAndWriteResourcesDifferences(terraformerResources TerraformerResourceIDToData, terraformResources TerraformStateResourceIDToData) (bool, error) {
	differences, err := m.identifyResourceDifferences(terraformerResources, terraformResources)
	if err != nil {
		return false, fmt.Errorf("[m.identifyResourceDifferences]%w", err)
	}

	err = m.writeDifferences(differences)
	if err != nil {
		return false, fmt.Errorf("[m.writeDifferences]%w", err)
	}

	return len(differences) > 0, nil
}

// identifyAndWriteDeletedResources found the resources differences and writes in the mapping file
func (m *ManagedResourcesDriftDetector) identifyAndWriteDeletedResources(terraformerResources TerraformerResourceIDToData, terraformResources TerraformStateResourceIDToData) (bool, error) {
	deleted, err := m.identifyDeletedResources(terraformerResources, terraformResources)
	if err != nil {
		return false, fmt.Errorf("[m.identifyDeletedResources]%w", err)
	}

	err = m.writeDeletedResources(deleted)
	if err != nil {
		return false, fmt.Errorf("[m.writeDeletedResources]%w", err)
	}

	return len(deleted) > 0, nil
}

// writeDeletedResources writes within a json file all the deleted resources to render within the PR
func (m *ManagedResourcesDriftDetector) writeDeletedResources(deleted []DeletedResource) error {
	differencesJSON, err := json.MarshalIndent(deleted, "", "  ")
	if err != nil {
		return fmt.Errorf("[json.MarshalIndent]%w", err)
	}

	return os.WriteFile("outputs/drift-resources-deleted.json", differencesJSON, 0400)
}

// writeDifferences writes within a json file the differences between all the drifted resources to render within the PR
func (m *ManagedResourcesDriftDetector) writeDifferences(differences []AttributeDifference) error {
	differencesJSON, err := json.MarshalIndent(differences, "", "  ")
	if err != nil {
		return fmt.Errorf("[json.MarshalIndent]%w", err)
	}

	return os.WriteFile("outputs/drift-resources-differences.json", differencesJSON, 0400)
}
