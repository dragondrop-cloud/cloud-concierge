package hclcreate

import (
	"reflect"
	"strings"
	"testing"

	"github.com/Jeffail/gabs/v2"
	terraformValueObjects "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_value_objects"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

func TestSplitResourceIdentifier(t *testing.T) {
	h := hclCreate{}

	expectedOutput := ResourceIdentifier{
		division:     "div",
		resourceType: "tf_type",
		resourceName: "tf_name",
	}

	output := h.splitResourceIdentifier("div.tf_type.tf_name")

	if !reflect.DeepEqual(expectedOutput, output) {
		t.Errorf("got %v, expected %v", output, expectedOutput)
	}
}

func TestExtractResourceBlockDefinition(t *testing.T) {
	h := hclCreate{}

	hclBytes := []byte(`resource "google_storage_bucket" "tfer--dragondrop-current-cloud-definitions-dev" {
  location                    = "US"
  name                        = "dragondrop-current-cloud-definitions-dev"
  project                     = "dragondrop-dev"
  uniform_bucket_level_access = "false"
}

resource "google_storage_bucket_acl" "tfer--dragondrop-modules" {
  default_event_based_hold = "false"
  force_destroy            = "false"
}

resource "google_storage_bucket_acl" "tfer--the-one" {
  default_event_based_hold = "false"
  force_destroy            = "false"
}
`)

	inputHCLFile, hclDiag := hclwrite.ParseConfig(
		hclBytes,
		"testName",
		hcl.Pos{Line: 0, Column: 0, Byte: 0},
	)

	if hclDiag != nil {
		t.Errorf("[hclwrite.ParseConfig] Unexpected error: %v", hclDiag)
	}

	expectedOutput := `resource "google_storage_bucket_acl" "the_one" {
  default_event_based_hold = "false"
  force_destroy            = "false"
}`

	output, err := h.extractResourceBlockDefinition(
		inputHCLFile,
		ConvertTerraformerResourceName("the-one"),
		ResourceIdentifier{resourceType: "google_storage_bucket_acl", resourceName: "tfer--the-one"},
	)

	if err != nil {
		t.Errorf("[h.extractResourceBlockDefinition] Unexpected error: %v", err)
	}

	if strings.Replace(expectedOutput, "\n", "", -1) != strings.Replace(createOutputString(output), "\n", "", -1) {
		t.Errorf("\ngot:\n%v\n\nexpected:\n%v", createOutputString(output), expectedOutput)
	}
}

func TestWriteBlockToWorkspaceHCL(t *testing.T) {
	h := hclCreate{}

	f := hclwrite.NewEmptyFile()

	newBlock := hclwrite.NewBlock("resource", []string{"tf_type", "tf_resource"})

	cloudActorTokens := hclwrite.Tokens{
		&hclwrite.Token{
			Type:         hclsyntax.TokenComment,
			Bytes:        []byte("cloudActorActionStatement"),
			SpacesBefore: 0,
		},
	}

	costEstimateTokens := hclwrite.Tokens{
		&hclwrite.Token{
			Type:         hclsyntax.TokenComment,
			Bytes:        []byte("costEstimateStatement"),
			SpacesBefore: 0,
		},
	}

	outputFile := h.writeBlockToWorkspaceHCL(
		hclwrite.NewEmptyFile(),
		cloudActorTokens,
		costEstimateTokens,
		newBlock,
	)

	f.Body().AppendUnstructuredTokens(costEstimateTokens)
	f.Body().AppendUnstructuredTokens(cloudActorTokens)
	f.Body().AppendNewline()
	f.Body().AppendBlock(newBlock)
	f.Body().AppendNewline()

	expectedOutput := string(f.Bytes())

	if expectedOutput != string(outputFile.Bytes()) {
		t.Errorf("got:\n%v\n\nexpected:\n%v", string(outputFile.Bytes()), expectedOutput)
	}
}

// createOutputString is a helper function used in testing for use in converting an *hclwrite.Block
// object to a string.
func createOutputString(block *hclwrite.Block) string {
	f := hclwrite.NewEmptyFile()
	rootBody := f.Body()

	rootBody.AppendBlock(block)

	return string(f.Bytes())
}

func TestCloudActionsToResourceActionMap(t *testing.T) {
	// Given
	h := hclCreate{}
	rawCloudActions := `{
		"google_storage_bucket.testing-out-this-bucket": {
			"creation": {
				"actor":"example@dragondrop.cloud",
				"timestamp":"2023-02-25"
			},
			"modified":{
				"actor":"example@dragondrop.cloud",
				"timestamp":"2023-03-08"
			}
		}
	}
`
	parsedJSON, err := gabs.ParseJSON([]byte(rawCloudActions))
	if err != nil {
		t.Errorf("[gabs.ParseJSON unexpected error]%v", err)
	}

	// When
	output, err := h.cloudActionsToResourceActionMap(parsedJSON)
	if err != nil {
		t.Errorf("[h.subsetCloudActionsToCurrentDivision]%v", err)
	}

	// Then
	expectedOutput := terraformValueObjects.ResourceActionMap{
		"google_storage_bucket.testing-out-this-bucket": {
			Creator: terraformValueObjects.CloudActorTimeStamp{
				Actor:     "example@dragondrop.cloud",
				Timestamp: "2023-02-25",
			},
			Modifier: terraformValueObjects.CloudActorTimeStamp{
				Actor:     "example@dragondrop.cloud",
				Timestamp: "2023-03-08",
			},
		},
	}

	if !reflect.DeepEqual(output, expectedOutput) {
		t.Errorf("got:\n%v\nexpected:\n%v", output, expectedOutput)
	}
}

func TestSubsetCloudActionsToCurrentDivisionsModifierOnly(t *testing.T) {
	// Given
	h := hclCreate{}

	rawCloudActions := `{
		"google_storage_bucket.testing-out-this-bucket": {
			"modified":{
				"actor":"example@dragondrop.cloud",
				"timestamp":"2025-03-08"
			}
		}
	}
`
	parsedJSON, err := gabs.ParseJSON([]byte(rawCloudActions))
	if err != nil {
		t.Errorf("[gabs.ParseJSON unexpected error]%v", err)
	}

	// When
	output, err := h.cloudActionsToResourceActionMap(parsedJSON)
	if err != nil {
		t.Errorf("[h.subsetCloudActionsToCurrentDivision]%v", err)
	}

	// Then
	expectedOutput := terraformValueObjects.ResourceActionMap{
		"google_storage_bucket.testing-out-this-bucket": {
			Modifier: terraformValueObjects.CloudActorTimeStamp{
				Actor:     "example@dragondrop.cloud",
				Timestamp: "2025-03-08",
			},
		},
	}

	if !reflect.DeepEqual(output, expectedOutput) {
		t.Errorf("got:\n%v\nexpected:\n%v", output, expectedOutput)
	}
}

func TestSubsetCloudActionsToCurrentDivisionsCreationOnly(t *testing.T) {
	// Given
	h := hclCreate{}
	rawCloudActions := `{
		"google_storage_bucket.testing-out-this-bucket": {
			"creation": {
				"actor":"example@dragondrop.cloud",
				"timestamp":"2023-02-25"
			}
		}
	}
`
	parsedJSON, err := gabs.ParseJSON([]byte(rawCloudActions))
	if err != nil {
		t.Errorf("[gabs.ParseJSON unexpected error]%v", err)
	}

	// When
	output, err := h.cloudActionsToResourceActionMap(parsedJSON)
	if err != nil {
		t.Errorf("[h.subsetCloudActionsToCurrentDivision]%v", err)
	}

	// Then
	expectedOutput := terraformValueObjects.ResourceActionMap{
		"google_storage_bucket.testing-out-this-bucket": {
			Creator: terraformValueObjects.CloudActorTimeStamp{
				Actor:     "example@dragondrop.cloud",
				Timestamp: "2023-02-25",
			},
		},
	}

	if !reflect.DeepEqual(output, expectedOutput) {
		t.Errorf("got:\n%v\nexpected:\n%v", output, expectedOutput)
	}
}

func TestSubsetCloudActionsToCurrentDivisionsEmpty(t *testing.T) {
	// Given
	h := hclCreate{}

	rawCloudActions := `{}`
	parsedJSON, err := gabs.ParseJSON([]byte(rawCloudActions))
	if err != nil {
		t.Errorf("[gabs.ParseJSON unexpected error]%v", err)
	}

	// When
	output, err := h.cloudActionsToResourceActionMap(parsedJSON)
	if err != nil {
		t.Errorf("[h.subsetCloudActionsToCurrentDivision]%v", err)
	}

	// Then
	expectedOutput := terraformValueObjects.ResourceActionMap{}

	if !reflect.DeepEqual(output, expectedOutput) {
		t.Errorf("got:\n%v\nexpected:\n%v", output, expectedOutput)
	}
}

func TestGenerateHCLActorsCommentCompleteData(t *testing.T) {
	// Given
	h := hclCreate{}
	inputResourceType := "google_storage_bucket"
	inputResourceName := "testing-out-this-bucket"

	inputResourceToCloudActions := terraformValueObjects.ResourceActionMap{
		"google_storage_bucket.testing-out-this-bucket": {
			Creator: terraformValueObjects.CloudActorTimeStamp{
				Actor:     "example@dragondrop.cloud",
				Timestamp: "2023-02-25",
			},
			Modifier: terraformValueObjects.CloudActorTimeStamp{
				Actor:     "example@dragondrop.cloud",
				Timestamp: "2023-03-08",
			},
		},
	}

	// When
	output := h.generateHCLCloudActorsComment(
		inputResourceType,
		inputResourceName,
		inputResourceToCloudActions,
	)

	// Then
	expectedOutputString := "\n# Created at 2023-02-25 by example@dragondrop.cloud\n# Last Modified at 2023-03-08 by example@dragondrop.cloud"
	if string(output.Bytes()) != expectedOutputString {
		t.Errorf("got:\n%v\nexpected:\n%v", string(output.Bytes()), expectedOutputString)
	}
}

func TestGenerateHCLActorsCommentModifierOnly(t *testing.T) {
	// Given
	h := hclCreate{}
	inputResourceType := "google_storage_bucket"
	inputResourceName := "testing-out-this-bucket"

	inputResourceToCloudActions := terraformValueObjects.ResourceActionMap{
		"google_storage_bucket.testing-out-this-bucket": {
			Modifier: terraformValueObjects.CloudActorTimeStamp{
				Actor:     "example@dragondrop.cloud",
				Timestamp: "2023-03-08",
			},
		},
	}

	// When
	output := h.generateHCLCloudActorsComment(
		inputResourceType,
		inputResourceName,
		inputResourceToCloudActions,
	)

	// Then
	expectedOutputString := "\n# Last Modified at 2023-03-08 by example@dragondrop.cloud"
	if string(output.Bytes()) != expectedOutputString {
		t.Errorf("got:\n%v\nexpected:\n%v", string(output.Bytes()), expectedOutputString)
	}
}

func TestGenerateHCLActorsCommentCreatorOnly(t *testing.T) {
	// Given
	h := hclCreate{}
	inputResourceType := "google_storage_bucket"
	inputResourceName := "testing-out-this-bucket"

	inputResourceToCloudActions := terraformValueObjects.ResourceActionMap{
		"google_storage_bucket.testing-out-this-bucket": {
			Creator: terraformValueObjects.CloudActorTimeStamp{
				Actor:     "example@dragondrop.cloud",
				Timestamp: "2023-02-25",
			},
		},
	}

	// When
	output := h.generateHCLCloudActorsComment(
		inputResourceType,
		inputResourceName,
		inputResourceToCloudActions,
	)

	// Then
	expectedOutputString := "\n# Created at 2023-02-25 by example@dragondrop.cloud"
	if string(output.Bytes()) != expectedOutputString {
		t.Errorf("got:\n%v\nexpected:\n%v", string(output.Bytes()), expectedOutputString)
	}
}

func TestGenerateHCLActorsCommentNoInput(t *testing.T) {
	// Given
	h := hclCreate{}
	inputResourceType := "google_storage_bucket"
	inputResourceName := "testing-out-this-bucket"

	inputResourceToCloudActions := terraformValueObjects.ResourceActionMap{}

	// When
	output := h.generateHCLCloudActorsComment(
		inputResourceType,
		inputResourceName,
		inputResourceToCloudActions,
	)

	// Then
	expectedOutputString := ""
	if string(output.Bytes()) != expectedOutputString {
		t.Errorf("got:\n%v\nexpected:\n%v", string(output.Bytes()), expectedOutputString)
	}
}
