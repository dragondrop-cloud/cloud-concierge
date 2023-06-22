package driftDetector

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	terraformValueObjects "github.com/dragondrop-cloud/driftmitigation/implementations/terraform_value_objects"
)

// ManagedResourcesDriftDetector is a type that identifies resources
// managed by Terraform that have drifted from their expected state.
type ManagedResourcesDriftDetector struct {
	// DivisionToProvider is a mapping between a division and the provider that is responsible
	// for that division.
	divisionToProvider map[terraformValueObjects.Division]terraformValueObjects.Provider `required:"true"`
}

// NewManagedResourcesDriftDetector generated a terraformer instance from ManagedResourcesDriftDetector
func NewManagedResourcesDriftDetector(divisionToProvider map[terraformValueObjects.Division]terraformValueObjects.Provider) *ManagedResourcesDriftDetector {
	return &ManagedResourcesDriftDetector{
		divisionToProvider: divisionToProvider,
	}
}

// Execute initiates the process of detecting drift in managed resources
// by comparing the current state of resources with their expected state.
// It takes a context as input to support cancellation and timeouts.
func (m *ManagedResourcesDriftDetector) Execute(ctx context.Context, workspaceToDirectory map[string]string) (bool, error) {
	remoteStateResources, err := m.loadAllRemoteStateFiles(workspaceToDirectory)
	if err != nil {
		return false, fmt.Errorf("[m.loadAllRemoteStateFiles]%w", err)
	}

	terraformerStateResources, err := m.loadAllTerraformerStateFiles()
	if err != nil {
		return false, fmt.Errorf("[m.loadAllTerraformerStateFiles]%w", err)
	}

	wereDeleted, err := m.identifyAndWriteDeletedResources(terraformerStateResources, remoteStateResources)
	if err != nil {
		return false, fmt.Errorf("[m.identifyAndWriteDeletedResources]%w", err)
	}

	differencesFound, err := m.identifyAndWriteResourcesDifferences(terraformerStateResources, remoteStateResources)
	if err != nil {
		return false, fmt.Errorf("[m.identifyAndWriteResourcesDifferences]%w", err)
	}

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

	return os.WriteFile("mappings/drift-resources-deleted.json", differencesJSON, 0400)
}

// writeDifferences writes within a json file the differences between all the drifted resources to render within the PR
func (m *ManagedResourcesDriftDetector) writeDifferences(differences []AttributeDifference) error {
	differencesJSON, err := json.MarshalIndent(differences, "", "  ")
	if err != nil {
		return fmt.Errorf("[json.MarshalIndent]%w", err)
	}

	return os.WriteFile("mappings/drift-resources-differences.json", differencesJSON, 0400)
}
