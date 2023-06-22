package terraformImportMigrationGenerator

import (
	"testing"

	terraformValueObjects "github.com/dragondrop-cloud/driftmitigation/implementations/terraform_value_objects"
)

// TODO: This needs implementation!
func TestConvertProviderToResourceImportMapToJSON(t *testing.T) {
	// Given
	inputResourceImportMap := terraformValueObjects.ProviderToResourceImportMap{
		"google": {
			"example-division": {
				"resourceName": terraformValueObjects.ImportMigration{
					TerraformConfigLocation: "config_location",
					RemoteCloudReference:    "remote.reference",
				},
			},
		},
	}
	i := TerraformImportMigrationGenerator{}

	// Then
	output, err := i.convertProviderToResourceImportMapToJSON(inputResourceImportMap)
	if err != nil {
		t.Errorf("unexpected error in i.convertProviderToResourceImportMapToJSON: %v", err)
	}

	// When
	expectedOutput := `{"google-example-division":{"resourceName":{"TerraformConfigLocation":"config_location","RemoteCloudReference":"remote.reference"}}}`
	if output != expectedOutput {
		t.Errorf("got:\n%v\nexpected:%v", output, expectedOutput)
	}
}
