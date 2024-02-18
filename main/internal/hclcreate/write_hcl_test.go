package hclcreate

import (
	"strconv"
	"testing"
)

func TestCreateMainTF(t *testing.T) {
	inputProvidersMap := map[string]string{"google": "~>4.27.0", "tfe": "~>0.33.0"}

	expectedOutput := "terraform {\n  required_version = \"~>1.2.4\"\n\n  required_providers {" +
		"\n    google = {\n      source  = \"hashicorp/google\"\n      version = \"~>4.27.0\"\n    }\n\n    tfe = {\n" +
		"      source  = \"hashicorp/tfe\"\n      version = \"~>0.33.0\"\n    }\n\n  }\n}\n"

	expectedOutputTwo := "terraform {\n  required_version = \"~>1.2.4\"\n\n  required_providers {" +
		"\n    tfe = {\n      source  = \"hashicorp/tfe\"\n      version = \"~>0.33.0\"\n    }\n\n" +
		"    google = {\n      source  = \"hashicorp/google\"\n      version = \"~>4.27.0\"\n    }\n\n  }\n}\n"

	hclCreate, _ := NewHCLCreate(
		Config{
			TerraformVersion: "~>1.2.4",
		},
		"",
	)
	f, err := hclCreate.CreateMainTF(inputProvidersMap)
	if err != nil {
		t.Errorf("unexpected error in createMainTF: %v", err)
	}

	fString := string(f)

	// Due to random return over map iteration, checks for different orders of provider specifications
	if (fString != expectedOutput) && (fString != expectedOutputTwo) {
		t.Errorf(
			"got:\n%s\n\n expected:\n%v\n\nor expected:\n%v",
			strconv.Quote(fString),
			strconv.Quote(expectedOutput),
			strconv.Quote(expectedOutputTwo),
		)
	}
}
