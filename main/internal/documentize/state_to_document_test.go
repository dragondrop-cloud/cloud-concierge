package documentize

import (
	"testing"

	"github.com/Jeffail/gabs/v2"

	terraformValueObjects "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_value_objects"
)

func TestConvertWorkspaceDocumentsToJSON(t *testing.T) {
	d := &documentize{}

	inputMap := map[Workspace][]byte{
		Workspace("example_1"): []byte("Workspace 1 document."),
		Workspace("example_2"): []byte("Workspace 2 document."),
	}

	output, err := d.ConvertWorkspaceDocumentsToJSON(inputMap)
	if err != nil {
		t.Errorf("Unexpected error in d.ConvertWorkspaceDocumentsToJSON: %v", err)
	}

	expectedOutput := `{"example_1":"Workspace 1 document.","example_2":"Workspace 2 document."}`

	if string(output) != expectedOutput {
		t.Errorf("got %v, expected %v", output, expectedOutput)
	}
}

func TestRegexProviderName(t *testing.T) {
	inputProvider := `provider["registry.terraform.io/hashicorp/google"]`

	expectedOutput := "provider[\"registry.terraform.io/hashicorp/google\"]"

	actualOutput, err := regexProviderName(inputProvider)
	if err != nil {
		t.Errorf("Unexpected error in regexProviderName: %v", err)
	}

	if actualOutput != expectedOutput {
		t.Errorf("got %v, expected %v", actualOutput, expectedOutput)
	}

	inputProvider = `module.google-networking.provider["registry.terraform.io/hashicorp/google-beta"]`

	expectedOutput = "provider[\"registry.terraform.io/hashicorp/google-beta\"]"

	actualOutput, err = regexProviderName(inputProvider)
	if err != nil {
		t.Errorf("Unexpected error in regexProviderName: %v", err)
	}

	if actualOutput != expectedOutput {
		t.Errorf("got %v, expected %v", actualOutput, expectedOutput)
	}

	inputProvider = `module.aws-persistent-storage.provider["registry.terraform.io/hashicorp/aws"].us_east_1`

	expectedOutput = "provider[\"registry.terraform.io/hashicorp/aws\"]"

	actualOutput, err = regexProviderName(inputProvider)
	if err != nil {
		t.Errorf("Unexpected error in regexProviderName: %v", err)
	}

	if actualOutput != expectedOutput {
		t.Errorf("got %v, expected %v", actualOutput, expectedOutput)
	}
}

func TestWorkspaceDocFromTFState(t *testing.T) {
	resourceExtractors := map[terraformValueObjects.Provider]ResourceExtractor{
		"google": NewGoogleResourceExtractor(),
	}

	d := &documentize{
		resourceExtractors: resourceExtractors,
	}

	tfStateParsed, err := gabs.ParseJSON([]byte(`{
  "resources": [
    {
      "module": "module.google-backend-api",
      "mode": "data",
      "type": "tfe_outputs"
    },
    {
      "module": "module.google-backend-api",
      "mode": "data",
      "type": "tfe_outputs",
      "name": "networking_output"
    },
    {
      "module": "module.google-backend-api",
      "mode": "data",
      "type": "tfe_outputs",
      "name": "persistent_storage",
      "provider": "provider[\"registry.terraform.io/hashicorp/tfe\"]"
    },
    {
      "module": "module.google-backend-api",
      "mode": "managed",
      "type": "google_cloud_run_service",
      "name": "api_compute",
      "provider": "provider[\"registry.terraform.io/hashicorp/google\"]",
      "instances": [
        {
          "schema_version": 1,
          "attributes": {
            "autogenerate_revision_name": false,
            "location": "us-east4",
            "metadata": [],
            "name": "cloud-run-api-dev",
            "project": "dragondrop-dev",
            "status": [],
            "template": [],
            "timeouts": null,
            "traffic": [
              {
                "latest_revision": true,
                "percent": 100,
                "revision_name": "",
                "tag": "",
                "url": ""
              }
            ]
          }
        }
      ]
    }]
}`))
	if err != nil {
		t.Errorf("Error parsing input state with gabs.ParseJSON(): %v", err)
	}

	outputResourceDetails, err := d.workspaceDocFromTFState(tfStateParsed)
	if err != nil {
		t.Errorf("Unexpected error from d.workspaceDocFromTFState: %v", outputResourceDetails)
	}

	expectedResourceDetails := "terraform name of api compute and type google cloud run" +
		" service within module google backend api resource at location us-east4 resource name of cloud run api dev resource project of dragondrop dev. "

	if outputResourceDetails != expectedResourceDetails {
		t.Errorf("got '%v', expected '%v'", outputResourceDetails, expectedResourceDetails)
	}
}
