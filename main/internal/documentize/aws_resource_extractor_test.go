package documentize

import (
	"reflect"
	"testing"

	"github.com/Jeffail/gabs/v2"
)

func TestExtractResourceDetailsOne_AWS(t *testing.T) {
	are := awsResourceExtractor{
		currentResourceDetails: &awsResourceDetails{},
	}

	tfStateParsed, err := gabs.ParseJSON([]byte(`{
  "version": 4,
  "terraform_version": "1.2.9",
  "serial": 26,
  "outputs": {
  },
  "resources": [
	{
      "module": "module.aws-backend-api",
      "mode": "data",
      "type": "tfe_outputs"
    },
    {
      "module": "module.dragondrop_compute.module.ecs_fargate_task",
      "mode": "managed",
      "type": "aws_ecs_cluster",
      "name": "fargate_cluster",
      "provider": "provider[\"registry.terraform.io/hashicorp/aws\"]",
      "instances": [
        {
          "schema_version": 0,
          "attributes": {
            "arn": "arn:aws:ecs:us-east-1:682649898103:cluster/dragondrop-ecs-fargate-cluster",
            "configuration": [],
            "id": "arn:aws:ecs:us-east-1:682649898103:cluster/dragondrop-ecs-fargate-cluster",
            "name": "dragondrop-ecs-fargate-cluster",
            "tags": {
              "env": "dev",
              "origin": "dragondrop-compute-module"
            },
            "tags_all": {
              "env": "dev",
              "origin": "dragondrop-compute-module"
            }
          },
          "sensitive_attributes": [],
          "private": "bnVsbA=="
        }
      ]
    }
  ]
}`))
	if err != nil {
		t.Errorf("Unexpected error in gabs.ParseJSON(): %v", err)
	}

	_ = are.ExtractResourceDetails(tfStateParsed, false, 1, 0)

	actualResourceDetails := are.GetCurrentResourceDetails()

	expectedResourceDetails := &awsResourceDetails{
		terraformModule:      "module.dragondrop_compute.module.ecs_fargate_task",
		terraformName:        "fargate_cluster",
		terraformType:        "aws_ecs_cluster",
		awsInstanceRegion:    "us-east-1",
		awsInstanceName:      "dragondrop-ecs-fargate-cluster",
		awsInstanceAccountID: "682649898103",
		awsInstanceTags: map[string]string{
			"env":    "dev",
			"origin": "dragondrop-compute-module",
		},
	}

	if !reflect.DeepEqual(actualResourceDetails, expectedResourceDetails) {
		t.Errorf("got:\n%v\nexpected:\n%v", actualResourceDetails, expectedResourceDetails)
	}
}

func TestExtractResourceDetailsTwo_AWS(t *testing.T) {
	are := awsResourceExtractor{
		currentResourceDetails: &awsResourceDetails{},
	}

	tfStateParsed, err := gabs.ParseJSON([]byte(`{
  "version": 4,
  "terraform_version": "1.2.9",
  "serial": 26,
  "outputs": {
  },
  "resources": [
	{
      "module": "module.aws-backend-api",
      "mode": "data",
      "type": "tfe_outputs"
    },
    {
      "mode": "managed",
      "type": "aws_ecs_cluster",
      "name": "fargate_cluster",
      "provider": "provider[\"registry.terraform.io/hashicorp/aws\"]",
      "instances": [
        {
          "schema_version": 0,
          "attributes": {
            "configuration": [],
            "id": "arn:aws:ecs:us-east-1:682649898103:cluster/dragondrop-ecs-fargate-cluster",
            "tags": {
              "env": "dev",
              "origin": "dragondrop-compute-module"
            },
            "tags_all": {}
          },
          "sensitive_attributes": [],
          "private": "bnVsbA=="
        }
      ]
    }
  ]
}`))
	if err != nil {
		t.Errorf("Unexpected error in gabs.ParseJSON(): %v", err)
	}

	_ = are.ExtractResourceDetails(tfStateParsed, false, 1, 0)

	actualResourceDetails := are.GetCurrentResourceDetails()

	expectedResourceDetails := &awsResourceDetails{
		terraformModule:      "none",
		terraformName:        "fargate_cluster",
		terraformType:        "aws_ecs_cluster",
		awsInstanceRegion:    "none",
		awsInstanceName:      "none",
		awsInstanceAccountID: "none",
		awsInstanceTags:      map[string]string{},
	}
	if !reflect.DeepEqual(actualResourceDetails, expectedResourceDetails) {
		t.Errorf("got:\n%v\nexpected:\n%v", actualResourceDetails, expectedResourceDetails)
	}
}

func TestResourceDetailsToSentence_AWS(t *testing.T) {
	// Base case
	are := awsResourceExtractor{
		currentResourceDetails: &awsResourceDetails{
			terraformModule:      "module.dragondrop_compute.module.ecs_fargate_task",
			terraformName:        "fargate_cluster",
			terraformType:        "aws_ecs_cluster",
			awsInstanceRegion:    "us-east-1",
			awsInstanceName:      "instance-example-name",
			awsInstanceAccountID: "68263",
			awsInstanceTags: map[string]string{
				"env":    "dev",
				"origin": "dragondrop-compute-module",
			},
		},
		typeToCategory: awsResourceCategories(),
	}

	output := are.ResourceDetailsToSentence()

	expectedOutput := "terraform name of fargate cluster and type aws ecs cluster within module dragondrop compute ecs fargate task " +
		"resource at location us-east-1 resource name of instance example name resource account of 68263 " +
		"with tag key of env and value of dev with tag key of origin and value of dragondrop compute module " +
		"with primary category of compute. "

	expectedOutputTwo := "terraform name of fargate cluster and type aws ecs cluster within module dragondrop compute ecs fargate task " +
		"resource at location us-east-1 resource name of instance example name resource account of 68263 " +
		"with tag key of origin and value of dragondrop compute module with tag key of env and value of dev " +
		"with primary category of compute. "

	if output != expectedOutput && output != expectedOutputTwo {
		t.Errorf("got:\n'%v'\nexpected one of:\n'%v'\nor:\n'%v'", output, expectedOutput, expectedOutputTwo)
	}

	// Dual category case
	are = awsResourceExtractor{
		currentResourceDetails: &awsResourceDetails{
			terraformModule:      "module.dragondrop.module.storage",
			terraformName:        "fargate_cluster",
			terraformType:        "aws_ecr_repository",
			awsInstanceRegion:    "us-east-1",
			awsInstanceName:      "instance-example-name",
			awsInstanceAccountID: "68263",
			awsInstanceTags: map[string]string{
				"env":    "dev",
				"origin": "dragondrop-compute-module",
			},
		},
		typeToCategory: awsResourceCategories(),
	}

	output = are.ResourceDetailsToSentence()

	expectedOutput = "terraform name of fargate cluster and type aws ecr repository within module dragondrop storage " +
		"resource at location us-east-1 resource name of instance example name resource account of 68263 " +
		"with tag key of env and value of dev with tag key of origin and value of dragondrop compute module " +
		"with primary category of containers and secondary category of storage. "

	expectedOutputTwo = "terraform name of fargate cluster and type aws ecr repository within module " +
		"dragondrop storage resource at location us-east-1 resource name of instance example name resource " +
		"account of 68263 with tag key of origin and value of dragondrop compute module with tag key " +
		"of env and value of dev with primary category of containers and secondary category of storage. "

	if output != expectedOutput && output != expectedOutputTwo {
		t.Errorf("got:\n'%v'\nexpected one of:\n'%v'\nor:\n'%v'", output, expectedOutput, expectedOutputTwo)
	}

	// None case
	are = awsResourceExtractor{
		currentResourceDetails: &awsResourceDetails{
			terraformModule:      "example",
			terraformName:        "fargate_cluster",
			terraformType:        "aws_ecr_xyz",
			awsInstanceRegion:    "none",
			awsInstanceName:      "none",
			awsInstanceAccountID: "none",
			awsInstanceTags:      map[string]string{},
		},
		typeToCategory: awsResourceCategories(),
	}

	output = are.ResourceDetailsToSentence()

	expectedOutput = "terraform name of fargate cluster and type aws ecr xyz within module example " +
		"resource at location none. "

	if output != expectedOutput {
		t.Errorf("got:\n'%v'\nexpected:\n'%v'", output, expectedOutput)
	}
}
