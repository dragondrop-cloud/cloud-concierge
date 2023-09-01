package terraformimportmigrationgenerator

import (
	"testing"

	terraformValueObjects "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_value_objects"
)

func TestConvertProviderToResourceImportMapToJSON(t *testing.T) {
	// Given
	inputResourceImportMap := terraformValueObjects.ResourceImportMap{
		"resourceName": terraformValueObjects.ImportMigration{
			TerraformConfigLocation: "config_location",
			RemoteCloudReference:    "remote.reference",
		},
	}

	i := TerraformImportMigrationGenerator{}

	// Then
	output, err := i.convertProviderToResourceImportMapToJSON(inputResourceImportMap)
	if err != nil {
		t.Errorf("unexpected error in i.convertProviderToResourceImportMapToJSON: %v", err)
	}

	// When
	expectedOutput := `{"resourceName":{"TerraformConfigLocation":"config_location","RemoteCloudReference":"remote.reference"}}`
	if output != expectedOutput {
		t.Errorf("got:\n%v\nexpected:%v", output, expectedOutput)
	}
}
