package terraformimportmigrationgenerator

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

	// Division is the name of a cloud division.
	Division terraformValueObjects.Division `required:"true"`
}

// TerraformImportMigrationGenerator is a struct that implements the interfaces.TerraformImportMigrationGenerator interface.
type TerraformImportMigrationGenerator struct {
	// dragonDrop is needed to inform cloud resources mapped to state file
	dragonDrop interfaces.DragonDrop

	// provider is the name of a cloud provider
	provider terraformValueObjects.Provider `required:"true"`

	// config contains the variables that determine the specific behavior of the TerraformImportMigrationGenerator struct.
	config Config
}

// NewTerraformImportMigrationGenerator creates and returns a new instance of TerraformImportMigrationGenerator
func NewTerraformImportMigrationGenerator(ctx context.Context, config Config, dragonDrop interfaces.DragonDrop, provider terraformValueObjects.Provider) interfaces.TerraformImportMigrationGenerator {
	dragonDrop.PostLog(ctx, "Created TFImport client.")

	return &TerraformImportMigrationGenerator{config: config, dragonDrop: dragonDrop, provider: provider}
}

// Execute generates terraform state migration statements for identified resources.
func (i *TerraformImportMigrationGenerator) Execute(ctx context.Context) error {
	i.dragonDrop.PostLog(ctx, "Beginning to map resources to import location.")

	resourceImports, err := i.GenericResourcesToImportLocation(i.provider)
	if err != nil {
		return fmt.Errorf("[terraform_import_migration_generator][error in GenericResourcesToImportLocation]%w", err)
	}

	err = i.dragonDrop.InformCloudResourcesMappedToStateFile(ctx)
	if err != nil {
		return fmt.Errorf("[terraform_import_migration_generator][error informing resources mapped to import location]%w", err)
	}

	resourceImportMapJSON, err := i.convertProviderToResourceImportMapToJSON(resourceImports)
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

	_ = os.MkdirAll("outputs", 0660)
	err = os.WriteFile("outputs/resources-to-import-location.json", []byte(resourceImportMapJSON), 0400)
	if err != nil {
		return fmt.Errorf("[map_resources][os.WriteFile(resources-to-import-location.json]%w", err)
	}

	return nil
}

// convertProviderToResourceImportMapToJSON converts the importMap to a json formatted string.
func (i *TerraformImportMigrationGenerator) convertProviderToResourceImportMapToJSON(importMap terraformValueObjects.ResourceImportMap) (string, error) {
	jsonObj := gabs.New()

	for resourceName, importLocation := range importMap {
		_, err := jsonObj.Set(
			importLocation,
			string(resourceName),
		)
		if err != nil {
			return "", fmt.Errorf("[convert_provider_to_resource_import][jsonObj.Set(%v, %v, %v, %v)]", importLocation, i.provider, i.config.Division, resourceName)
		}
	}

	return jsonObj.String(), nil
}
