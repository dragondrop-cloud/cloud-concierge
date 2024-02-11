package documentize

import (
	"reflect"
	"testing"

	"github.com/Jeffail/gabs/v2"
)

func TestExtractResourceDetailsOne(t *testing.T) {
	gre := googleResourceExtractor{
		currentResourceDetails: &googleResourceDetails{},
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
		t.Errorf("Unexpected error in gabs.ParseJSON(): %v", err)
	}

	_ = gre.ExtractResourceDetails(tfStateParsed, false, 3, 0)

	actualResourceDetails := gre.GetCurrentResourceDetails()

	expectedResourceDetails := &googleResourceDetails{
		terraformModule:     "module.google-backend-api",
		terraformName:       "api_compute",
		terraformType:       "google_cloud_run_service",
		gcpInstanceLocation: "us-east4",
		gcpInstanceName:     "cloud-run-api-dev",
		gcpInstanceProject:  "dragondrop-dev",
	}

	if !reflect.DeepEqual(actualResourceDetails, expectedResourceDetails) {
		t.Errorf("got %v, expected %v", actualResourceDetails, expectedResourceDetails)
	}
}

func TestExtractResourceDetailsTwo(t *testing.T) {
	gre := googleResourceExtractor{
		currentResourceDetails: &googleResourceDetails{},
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
      "mode": "managed",
      "type": "google_cloud_run_service",
      "name": "api_compute",
      "provider": "provider[\"registry.terraform.io/hashicorp/google\"]",
      "instances": [
        {
          "schema_version": 1,
          "attributes": {
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

	_ = gre.ExtractResourceDetails(tfStateParsed, false, 3, 0)

	actualResourceDetails := gre.GetCurrentResourceDetails()

	expectedResourceDetails := &googleResourceDetails{
		terraformModule:     "none",
		terraformName:       "api_compute",
		terraformType:       "google_cloud_run_service",
		gcpInstanceLocation: "global",
		gcpInstanceName:     "none",
		gcpInstanceProject:  "none",
	}
	if !reflect.DeepEqual(actualResourceDetails, expectedResourceDetails) {
		t.Errorf("got %v, expected %v", actualResourceDetails, expectedResourceDetails)
	}
}

func TestResourceDetailsToSentence(t *testing.T) {
	// Base case
	gre := googleResourceExtractor{
		currentResourceDetails: &googleResourceDetails{
			terraformModule:     "module.tfModule",
			terraformName:       "tf-Name",
			terraformType:       "google_cloudbuild_trigger",
			gcpInstanceLocation: "us-east4",
			gcpInstanceName:     "instance_-example-_name",
			gcpInstanceProject:  "example-dev",
		},
		typeToCategory: googleResourceCategories(),
	}

	output := gre.ResourceDetailsToSentence()

	expectedOutput := "terraform name of tf name and type google cloudbuild trigger within module tfmodule " +
		"resource at location us-east4 resource name of instance example name resource project of example dev " +
		"with primary category of ci cd. "

	if output != expectedOutput {
		t.Errorf("got '%v', expected '%v'", output, expectedOutput)
	}

	// Dual category case
	gre = googleResourceExtractor{
		currentResourceDetails: &googleResourceDetails{
			terraformModule:     "module.tfModule",
			terraformName:       "tf-Name",
			terraformType:       "google_bigquery_table",
			gcpInstanceLocation: "us-east4",
			gcpInstanceName:     "instance_-example-_name",
			gcpInstanceProject:  "example-dev",
		},
		typeToCategory: googleResourceCategories(),
	}

	output = gre.ResourceDetailsToSentence()

	expectedOutput = "terraform name of tf name and type google bigquery table within module tfmodule " +
		"resource at location us-east4 resource name of instance example name resource project of example dev " +
		"with primary category of storage and secondary category of analytics. "

	if output != expectedOutput {
		t.Errorf("got '%v', expected '%v'", output, expectedOutput)
	}

	// None case
	gre = googleResourceExtractor{
		currentResourceDetails: &googleResourceDetails{
			terraformModule:     "module.tfModule",
			terraformName:       "tf-Name",
			terraformType:       "tf_example_type",
			gcpInstanceLocation: "global",
			gcpInstanceName:     "none",
			gcpInstanceProject:  "none",
		},
	}

	output = gre.ResourceDetailsToSentence()

	expectedOutput = "terraform name of tf name and type tf example type within module tfmodule " +
		"resource at location global. "

	if output != expectedOutput {
		t.Errorf("got '%v', expected '%v'", output, expectedOutput)
	}
}
