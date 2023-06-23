package hclcreate

import (
	"reflect"
	"testing"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

func Test_GenerateImportBlockFile(t *testing.T) {
	// Given
	h := hclCreate{}

	inputResourceToImportLoc := ResourceImportsByDivision{
		"dev-division": {
			"resource_type_1.resource_name_1": {
				TerraformConfigLocation: "resource_type_1.resource_name_1",
				RemoteCloudReference:    "remote/cloud/reference",
			},
		},
		"prod-division": {
			"resource_type_2.resource_name_2": {
				TerraformConfigLocation: "resource_type_2.resource_name_2",
				RemoteCloudReference:    "remote/cloud/reference/2",
			},
		},
	}

	inputResourceToWorkspace := NewResourceToWorkspace{
		"dev-division.resource_type_1.resource_name_1":  "my-dev-workspace",
		"prod-division.resource_type_2.resource_name_2": "my-prod-workspace",
	}

	inputWorkspace := "my-dev-workspace"

	expectedOutputFile := hclwrite.NewEmptyFile()
	expectedOutputBody := expectedOutputFile.Body()

	importBlock := expectedOutputBody.AppendNewBlock("import", nil)
	importBlock.Body().SetAttributeValue("to", cty.StringVal("resource_type_1.resource_name_1"))
	importBlock.Body().SetAttributeValue("id", cty.StringVal("remote/cloud/reference"))

	expectedOutput := string(expectedOutputFile.Bytes())

	// When
	hclFile, err := h.generateImportBlockFile(
		inputWorkspace,
		inputResourceToImportLoc,
		inputResourceToWorkspace,
	)

	if err != nil {
		t.Errorf("unexpected error in h.generateImportBlockFile: %v", err)
	}

	// Then
	if string(hclFile) != expectedOutput {
		t.Errorf("expected:\n%v\ngot:\n%v", expectedOutput, hclFile)
	}
}

func Test_SetOfWorkspacesWithMigrationsStruct(t *testing.T) {
	// Given
	h := hclCreate{}

	inputResourceToWorkspace := NewResourceToWorkspace{
		"dev-division.resource_type_1.resource_name_1": "my-dev-workspace",
	}

	expectedOutput := map[string]bool{
		"my-dev-workspace": true,
	}

	// When
	output := h.setOfWorkspacesWithMigrationsStruct(inputResourceToWorkspace)

	// Then
	if !reflect.DeepEqual(output, expectedOutput) {
		t.Errorf("expected:\n%v\ngot:\n%v", expectedOutput, output)
	}

}
