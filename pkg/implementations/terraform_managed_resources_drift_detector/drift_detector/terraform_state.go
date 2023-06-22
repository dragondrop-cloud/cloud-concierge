package driftDetector

import (
	"fmt"
	"os"
)

// TerraformStateFile represents the structure of a Terraform state file from terraform cloud.
type TerraformStateFile struct {
	Resources []*Resource `json:"resources"`
}

// TerraformStateResourceIDToData is a map between a resource's unique id and the terraform state files unique
// resource data.
type TerraformStateResourceIDToData map[string]TerraformStateUniqueResourceData

// TerraformStateUniqueResourceData is a struct for storing the data definition of a single Terraform State resource instance.
type TerraformStateUniqueResourceData struct {
	StateFile  string
	Module     string
	Type       string
	Name       string
	Provider   string
	Attributes map[string]interface{}
}

// Resource represents a Terraform resource within a state file.
type Resource struct {
	Mode      string             `json:"mode"`
	Module    string             `json:"module"`
	Type      string             `json:"type"`
	Name      string             `json:"name"`
	Provider  string             `json:"provider"`
	Instances []ResourceInstance `json:"instances"`
}

// ResourceInstance represents a Terraform resource instance within a state file.
type ResourceInstance struct {
	SchemaVersion int                    `json:"schema_version"`
	Attributes    map[string]interface{} `json:"attributes"`
}

// loadAllRemoteStateFiles loads from memory the remote state files and aggregates data.
func (m *ManagedResourcesDriftDetector) loadAllRemoteStateFiles(workspaceToDirectory map[string]string) (TerraformStateResourceIDToData, error) {
	fileNames := make([]string, 0)
	for workspaceName := range workspaceToDirectory {
		fileNames = append(fileNames, workspaceName)
	}

	resources := TerraformStateResourceIDToData{}

	for _, remoteStateFile := range fileNames {
		fileContent, err := os.ReadFile(fmt.Sprintf("state_files/%v.json", remoteStateFile))
		if err != nil {
			return nil, fmt.Errorf("failed to read state file %s: %v", remoteStateFile, err)
		}

		file, err := m.parseRemoteStateFile(fileContent)
		if err != nil {
			return nil, err
		}

		resourcesFromStateFile := m.terraformStateExtractUniqueResourceIDToData(remoteStateFile, file)

		for resourceID, resourceData := range resourcesFromStateFile {
			resources[resourceID] = resourceData
		}
	}

	return resources, nil
}

// terraformStateExtractUniqueResourceIDToData reformats resource data to pull out the attribute "id" as the unique
// resource identifier.
func (m *ManagedResourcesDriftDetector) terraformStateExtractUniqueResourceIDToData(stateFileName string, stateFile TerraformStateFile) TerraformStateResourceIDToData {
	outputIDToData := TerraformStateResourceIDToData{}

	for _, resource := range stateFile.Resources {
		for _, instance := range resource.Instances {
			id := fmt.Sprintf("%v.%v", resource.Type, instance.Attributes["id"])
			outputIDToData[id] = TerraformStateUniqueResourceData{
				StateFile:  stateFileName,
				Module:     resource.Module,
				Type:       resource.Type,
				Name:       resource.Name,
				Provider:   resource.Provider,
				Attributes: instance.Attributes,
			}
		}
	}
	return outputIDToData
}
