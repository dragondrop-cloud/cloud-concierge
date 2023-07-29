package terraformImportMigrationGenerator

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Jeffail/gabs/v2"
	terraformValueObjects "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_value_objects"
	"github.com/dragondrop-cloud/cloud-concierge/main/internal/interfaces"
)

// Config is a struct for variables that determine the specific behavior of the TerraformImportMigrationGenerator struct.
type Config struct {
	// CloudCredential is a cloud credential with read-only access to a cloud division and, if applicable, access to read Terraform state files.
	CloudCredential terraformValueObjects.Credential `required:"true"`
}

// TerraformImportMigrationGenerator is a struct that implements the interfaces.TerraformImportMigrationGenerator interface.
type TerraformImportMigrationGenerator struct {
	// dragonDrop is needed to inform cloud resources mapped to state file
	dragonDrop interfaces.DragonDrop

	// DivisionToProvider is a map between the string representing a Division and the corresponding
	// cloud Provider (aws, azurerm, google, etc.).
	// For AWS, an account is the Division, for GCP a project name is the Division,
	// and for azurerm a resource group is a Division.
	divisionToProvider map[terraformValueObjects.Division]terraformValueObjects.Provider `required:"true"`

	// config contains the variables that determine the specific behavior of the TerraformImportMigrationGenerator struct.
	config Config
}

// NewTerraformImportMigrationGenerator creates and returns a new instance of TerraformImportMigrationGenerator
func NewTerraformImportMigrationGenerator(ctx context.Context, config Config, dragonDrop interfaces.DragonDrop, divisionToProvider map[terraformValueObjects.Division]terraformValueObjects.Provider) interfaces.TerraformImportMigrationGenerator {
	dragonDrop.PostLog(ctx, "Created TFImport client.")

	return &TerraformImportMigrationGenerator{config: config, dragonDrop: dragonDrop, divisionToProvider: divisionToProvider}
}

// Execute generates terraform state migration statements for identified resources.
func (i *TerraformImportMigrationGenerator) Execute(ctx context.Context) error {
	i.dragonDrop.PostLog(ctx, "Beginning to map resources to import location.")

	providerToResourceImportMap, err := i.mapResourcesToImportLocation()
	if err != nil {
		return fmt.Errorf("[terraform_import_migration_generator][error mapping resources to import location]%w", err)
	}

	err = i.dragonDrop.InformCloudResourcesMappedToStateFile(ctx)
	if err != nil {
		return fmt.Errorf("[terraform_import_migration_generator][error informing resources mapped to import location]%w", err)
	}

	resourceImportMapJSON, err := i.convertProviderToResourceImportMapToJSON(providerToResourceImportMap)
	if err != nil {
		return fmt.Errorf("[terraform_import_migration_generator][error converting Provider to resource import]%w", err)
	}

	err = i.writeResourcesMap(resourceImportMapJSON)
	if err != nil {
		return fmt.Errorf("[terraform_import_migration_generator][error mapping resources]%w", err)
	}

	i.dragonDrop.PostLog(ctx, "Generated map of remote resources to import location.")
	return nil
}

// writeResourcesMap writes a json file for the resource Import Map.
func (i *TerraformImportMigrationGenerator) writeResourcesMap(resourceImportMapJSON string) error {
	err := os.Chdir("/")
	if err != nil {
		return fmt.Errorf("[map_resources][os.Chdir(/)]%w", err)
	}

	// Check to see if running locally and move into the docker volume directory if necessary
	files, err := ioutil.ReadDir("./")
	if err != nil {
		return fmt.Errorf("[map_resources][ioutil.ReadDir('./')]%w", err)
	}

	for _, f := range files {
		fmt.Println(f.Name())
		if f.Name() == "main" {
			err = os.Chdir("/main")
			if err != nil {
				return fmt.Errorf("[map_resources][os.Chdir(/main)]%w", err)
			}
			break
		}
	}

	_ = os.MkdirAll("mappings", 0660)
	err = os.WriteFile("mappings/resources-to-import-location.json", []byte(resourceImportMapJSON), 0400)
	if err != nil {
		return fmt.Errorf("[map_resources][os.WriteFile(resources-to-import-location.json]%w", err)
	}

	return nil
}

// convertProviderToResourceImportMapToJSON converts the importMap to a json formatted string.
func (i *TerraformImportMigrationGenerator) convertProviderToResourceImportMapToJSON(importMap terraformValueObjects.ProviderToResourceImportMap) (string, error) {
	jsonObj := gabs.New()

	for provider, divisionToResourceImportMap := range importMap {
		for division, resourceImportMap := range divisionToResourceImportMap {
			for resourceName, importLocation := range resourceImportMap {
				_, err := jsonObj.Set(
					importLocation,
					fmt.Sprintf("%v-%v", string(provider), string(division)),
					string(resourceName),
				)
				if err != nil {
					return "", fmt.Errorf("[convert_provider_to_resource_import][jsonObj.Set(%v, %v, %v, %v)]", importLocation, provider, division, resourceName)
				}
			}
		}
	}

	return jsonObj.String(), nil
}

// TODO: Should try to write unit tests for this helper function if possible.
// mapResourcesToImportLocation maps cloud resources to the appropriate Terraform import migration statement.
func (i *TerraformImportMigrationGenerator) mapResourcesToImportLocation() (terraformValueObjects.ProviderToResourceImportMap, error) {
	providerToResourceImportMap := terraformValueObjects.ProviderToResourceImportMap{}

	for division, provider := range i.divisionToProvider {
		divisionResourceImports, err := GenericResourcesToImportLocation(division, provider)
		if err != nil {
			return nil, fmt.Errorf("[map_resources_to_import_location][error in GenericResourcesToImportLocation]%w", err)
		}

		_, okay := providerToResourceImportMap[provider]
		if !okay {
			providerToResourceImportMap[provider] = map[terraformValueObjects.Division]terraformValueObjects.ResourceImportMap{}
		}

		providerToResourceImportMap[provider][division] = divisionResourceImports[division]
	}

	return providerToResourceImportMap, nil
}
