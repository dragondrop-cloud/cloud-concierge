package documentize

import (
	"reflect"
	"testing"

	"github.com/Jeffail/gabs/v2"

	terraformValueObjects "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_value_objects"
)

func TestSelectNewResources(t *testing.T) {
	inputWorkspaceToID := map[Workspace]map[ResourceData]bool{
		"workspace_1": {
			ResourceData{
				id:     "8675309",
				tfType: "example_terraform_resource",
			}: true,
			ResourceData{
				id:     "8675309",
				tfType: "example_terraform_resource_2",
			}: true,
		},
		"workspace_2": {
			ResourceData{
				id:     "8675309asd",
				tfType: "example_terraform_resource",
			}: true,
			ResourceData{
				id:     "8asd",
				tfType: "example_terraform_resource_3",
			}: true,
		},
	}

	inputDivisionToID := map[terraformValueObjects.Division]map[ResourceData]bool{
		"division_1": {
			ResourceData{
				id:     "8675309",
				tfType: "example_terraform_resource",
			}: true,
		},
		"division_2": {
			ResourceData{
				id:     "8675309",
				tfType: "example_terraform_resource_2",
			}: true,
			ResourceData{
				id:     "86asd09asd",
				tfType: "le_terraform_resource",
			}: true,
			ResourceData{
				id:     "asd1118asd",
				tfType: "example_terraform_resource_3",
			}: true,
		},
	}

	expectedOutput := map[terraformValueObjects.Division]map[ResourceData]bool{
		"division_2": {
			ResourceData{
				id:     "86asd09asd",
				tfType: "le_terraform_resource",
			}: true,
			ResourceData{
				id:     "asd1118asd",
				tfType: "example_terraform_resource_3",
			}: true,
		},
	}

	actualOutput := selectNewResources(inputWorkspaceToID, inputDivisionToID)

	if !reflect.DeepEqual(actualOutput, expectedOutput) {
		t.Errorf("got %v, expected %v", actualOutput, expectedOutput)
	}

}

func TestCheckIfResourceIsPresent(t *testing.T) {
	inputTFRResource := ResourceData{
		id:     "8675309",
		tfType: "example_terraform_resource",
	}

	inputTypeToIDPresent := map[Workspace]map[ResourceData]bool{
		"workspace_1": {
			ResourceData{
				id:     "8675309",
				tfType: "example_terraform_resource",
			}: true,
			ResourceData{
				id:     "8675309",
				tfType: "example_terraform_resource_2",
			}: true,
		},
		"workspace_2": {
			ResourceData{
				id:     "8675309asd",
				tfType: "example_terraform_resource",
			}: true,
			ResourceData{
				id:     "8asd",
				tfType: "example_terraform_resource_3",
			}: true,
		},
	}

	inputTypeToIDNotPresent := map[Workspace]map[ResourceData]bool{
		"workspace_1": {
			ResourceData{
				id:     "8675309",
				tfType: "example_terraform_resource_2",
			}: true,
		},
		"workspace_2": {
			ResourceData{
				id:     "8675309asd",
				tfType: "example_terraform_resource",
			}: true,
			ResourceData{
				id:     "8asd",
				tfType: "example_terraform_resource_3",
			}: true,
		},
	}

	resourceBlackList := newResourceBlackList()

	if isValidNewResource(resourceBlackList, inputTFRResource, inputTypeToIDPresent) {
		t.Errorf("got 'true', expected 'false")
	}

	if !isValidNewResource(resourceBlackList, inputTFRResource, inputTypeToIDNotPresent) {
		t.Errorf("got 'false', expected 'true")
	}
}

func TestConvertNewResourcesToJSON(t *testing.T) {
	inputResourceDocMap := map[ResourceName]string{
		"resourceName_1": "doc_1",
		"resourceName_2": "doc_2",
	}

	expectedOutput := `{"resourceName_1":"doc_1","resourceName_2":"doc_2"}`

	d := &documentize{}
	actualOutput, err := d.ConvertNewResourcesToJSON(inputResourceDocMap)

	if err != nil {
		t.Errorf("Unexpected error in ConvertNewResourcesToJSON(): %v", err)
	}

	if string(actualOutput) != expectedOutput {
		t.Errorf("got %v, expected %v", string(actualOutput), expectedOutput)
	}
}

func TestExtractResourceDocument(t *testing.T) {
	tfStateParsed, err := gabs.ParseJSON([]byte(`{
  "resources": [
    {
      "module": "module.google-backend-api",
      "mode": "data",
      "type": "tfe_outputs",
      "name": "iam_secrets_output"
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
          "attributes_flat": {
			"id": "example_id",
            "autogenerate_revision_name": false,
            "metadata": [],
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
		t.Errorf("Unexpected error in gabs.ParseJSON(): %v", err)
	}

	d := &documentize{
		resourceExtractors: map[terraformValueObjects.Provider]ResourceExtractor{
			"google": &googleResourceExtractor{
				currentResourceDetails: &googleResourceDetails{},
			},
		},
	}
	resourceDoc, err := d.extractResourceDocument(tfStateParsed, true, 3, 0)

	expectedOutput := "terraform name of api compute and type" +
		" google cloud run service within module google backend api resource at location global. "

	if err != nil {
		t.Errorf("Unexpected error from d.extractResourceDocument(): %v", err)
	}

	if resourceDoc != expectedOutput {
		t.Errorf("got '%v', expected '%v'", resourceDoc, expectedOutput)
	}
}

func TestExtractResourceIdsFromTerraformerState(t *testing.T) {
	tfStateParsed, err := gabs.ParseJSON([]byte(`{
  "resources": [
    {
      "module": "module.google-backend-api",
      "mode": "data",
      "type": "tfe_outputs",
      "name": "iam_secrets_output"
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
          "attributes_flat": {
			"id": "example_id",
            "autogenerate_revision_name": false,
            "metadata": [],
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
		t.Errorf("Unexpected error in gabs.ParseJSON(): %v", err)
	}

	outputMap, _ := extractResourceIdsFromTerraformerState(tfStateParsed)

	expectedOutputMap := map[ResourceData]bool{
		ResourceData{
			id:     "example_id",
			name:   "api_compute",
			tfType: "google_cloud_run_service",
		}: true,
	}

	if !reflect.DeepEqual(outputMap, expectedOutputMap) {
		t.Errorf("got %v, expected %v", outputMap, expectedOutputMap)
	}
}

func TestExtractResourceIdsFromWorkspaceState(t *testing.T) {
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
			"id": "example_id",
            "autogenerate_revision_name": false,
            "metadata": [],
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
		t.Errorf("Unexpected error in gabs.ParseJSON(): %v", err)
	}

	outputMap, _ := extractResourceIdsFromWorkspaceState(tfStateParsed)

	expectedOutputMap := map[ResourceData]bool{
		ResourceData{
			id:     "example_id",
			tfType: "google_cloud_run_service",
		}: true,
	}

	if !reflect.DeepEqual(outputMap, expectedOutputMap) {
		t.Errorf("got %v, expected %v", outputMap, expectedOutputMap)
	}
}

func TestPullResourceDocumentFromDiv(t *testing.T) {
	tfStateParsed, err := gabs.ParseJSON([]byte(`{
  "resources": [
    {
      "module": "module.google-backend-api",
      "mode": "data",
      "type": "tfe_outputs",
      "name": "iam_secrets_output"
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
          "attributes_flat": {
			"id": "example_id",
            "autogenerate_revision_name": false,
            "metadata": [],
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
          "attributes_flat": {
			"id": "example_id2",
			"location": "us-east4",
            "autogenerate_revision_name": false,
            "metadata": [],
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
		t.Errorf("Unexpected error in gabs.ParseJSON(): %v", err)
	}

	d := &documentize{
		resourceExtractors: map[terraformValueObjects.Provider]ResourceExtractor{
			"google": &googleResourceExtractor{
				currentResourceDetails: &googleResourceDetails{},
			},
		},
	}

	inputDiv := terraformValueObjects.Division("exampleDiv")

	inputResourceData := ResourceData{
		tfType: "google_cloud_run_service",
		id:     "example_id",
		name:   "api_compute",
	}

	resourceName, doc, err := d.pullResourceDocumentFromDiv(tfStateParsed, inputDiv, inputResourceData)

	if err != nil {
		t.Errorf("[d.pullResourceDocumentFromDiv] Unexpected error %v", err)
	}

	expectedResourceName := "exampleDiv.google_cloud_run_service.api_compute"

	if resourceName != ResourceName(expectedResourceName) {
		t.Errorf("resourceName got %v, expected %v", resourceName, expectedResourceName)
	}

	expectedDoc := "terraform name of api compute and type" +
		" google cloud run service within module google backend api resource at location global. "

	if doc != expectedDoc {
		t.Errorf("resource document got %v, expected %v", doc, expectedDoc)
	}

}
