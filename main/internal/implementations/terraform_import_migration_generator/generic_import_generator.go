package terraformImportMigrationGenerator

import (
	"fmt"
	"os"

	"github.com/Jeffail/gabs/v2"
	log "github.com/sirupsen/logrus"

	terraformValueObjects "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_value_objects"
)

// GenericResourcesToImportLocation creates a map between any cloud provider resources Terraform definition location and the corresponding Import location reference.
func GenericResourcesToImportLocation(division terraformValueObjects.Division, provider terraformValueObjects.Provider) (map[terraformValueObjects.Division]terraformValueObjects.ResourceImportMap, error) {
	stateFileContent, err := readTerraformerStateFile(division, provider)
	if err != nil {
		return nil, err
	}

	return mapResourcesToImportLocation(division, provider, stateFileContent)
}

// readTerraformerResourcesFile reads the terraformer resource.tf file
func readTerraformerResourcesFile(division terraformValueObjects.Division, provider terraformValueObjects.Provider) ([]byte, error) {
	return readTerraformerFileByName(division, provider, "resources.tf")
}

// readTerraformerStateFile reads the terraformer terraform.tfstate file
func readTerraformerStateFile(division terraformValueObjects.Division, provider terraformValueObjects.Provider) ([]byte, error) {
	return readTerraformerFileByName(division, provider, "terraform.tfstate")
}

// readTerraformerFileByName reads files from /current cloud directory specifying its name
func readTerraformerFileByName(division terraformValueObjects.Division, provider terraformValueObjects.Provider, fileName string) ([]byte, error) {
	fileContent, err := os.ReadFile(fmt.Sprintf("%s-%s/%s", provider, division, fileName))
	if err != nil {
		return nil, fmt.Errorf("[generic_import_migration_generator][error reading terraformer %s file]%w", fileName, err)
	}

	return fileContent, nil
}

// mapResourcesToImportLocation maps the resources locations using the terraform.tfstate file
func mapResourcesToImportLocation(division terraformValueObjects.Division, provider terraformValueObjects.Provider, stateFileContent []byte) (map[terraformValueObjects.Division]terraformValueObjects.ResourceImportMap, error) {
	resourceImportMap := terraformValueObjects.ResourceImportMap{}

	stateFileJSON, err := gabs.ParseJSON(stateFileContent)
	if err != nil {
		return nil, fmt.Errorf("[generic_import_migration_generator][error parsing state file content]%w", err)
	}

	for _, resource := range stateFileJSON.S("resources").Children() {
		if resource.Exists("name") {
			resourceName := resource.S("name").Data().(string)
			resourceType := resource.S("type").Data().(string)
			fmt.Printf(
				"Calculating the import statement for resource type: %v, resource name: %v\n",
				resourceType,
				resourceName,
			)

			terraformConfigLocation, err := getTerraformConfigLocation(resourceType, resourceName)
			if err != nil {
				return nil, fmt.Errorf("[generic_import_migration_generator][error obtaining terraform config location]%w", err)
			}

			remoteCloudReference, err := GetRemoteCloudReference(resource, provider, ResourceType(resourceType))
			if err != nil {
				return nil, fmt.Errorf("[generic_import_migration_generator][error obtaining remote cloud reference]%w", err)
			}

			resourceImportMap[terraformValueObjects.ResourceName(terraformConfigLocation)] = terraformValueObjects.ImportMigration{
				TerraformConfigLocation: terraformConfigLocation,
				RemoteCloudReference:    terraformValueObjects.RemoteCloudReference(remoteCloudReference),
			}
		} else {
			log.Warnf("Resource doesn't have name: %v", resource)
		}
	}

	divisionToResourceImportMap := map[terraformValueObjects.Division]terraformValueObjects.ResourceImportMap{
		division: resourceImportMap,
	}
	return divisionToResourceImportMap, nil
}

// getTerraformConfigLocation gets the resources location with the specific format applied
func getTerraformConfigLocation(resourceType string, resourceName string) (terraformValueObjects.TerraformConfigLocation, error) {
	importLocation := fmt.Sprintf("%s.%s", resourceType, resourceName)
	return terraformValueObjects.TerraformConfigLocation(importLocation), nil
}
