package terraformImportMigrationGenerator

import (
	"testing"

	"github.com/stretchr/testify/assert"

	terraformValueObjects "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_value_objects"
)

func TestAWSImportMigrationGenerator_mapResourcesToImportLocation(t *testing.T) {
	// Given
	stateFileContent := []byte(`{
		"version": 4,
		"terraform_version": "1.2.6",
		"serial": 2,
		"resources": [
			{
				"name": "tfer--dragondrop-example-2",
				"type": "aws_s3_bucket",
				"instances": [
					{
						"attributes_flat": {
							"id": "dragondrop-example-2",
					  		"arn": "arn:aws:s3:::dragondrop-example-2"
						}
					}
				]
			},
			{
				"name": "tfer--dragondrop-example",
				"type": "aws_s3_bucket",
				"instances": [
					{
						"attributes_flat": {
							"id": "dragondrop-example",
					  		"arn": "arn:aws:s3:::dragondrop-example"
						}
					}
				]
			}
		]
	}`)
	provider := terraformValueObjects.Provider("aws")

	i := TerraformImportMigrationGenerator{}

	// When
	divisionToResourceImportMap, err := i.mapResourcesToImportLocation(provider, stateFileContent)

	// Then
	assert.Nil(t, err)

	expectedDivisionToResourceImportMap := terraformValueObjects.ResourceImportMap{
		terraformValueObjects.ResourceName("aws_s3_bucket.tfer--dragondrop-example"): terraformValueObjects.ImportMigration{
			TerraformConfigLocation: terraformValueObjects.TerraformConfigLocation("aws_s3_bucket.tfer--dragondrop-example"),
			RemoteCloudReference:    terraformValueObjects.RemoteCloudReference("dragondrop-example"),
		},
		terraformValueObjects.ResourceName("aws_s3_bucket.tfer--dragondrop-example-2"): terraformValueObjects.ImportMigration{
			TerraformConfigLocation: terraformValueObjects.TerraformConfigLocation("aws_s3_bucket.tfer--dragondrop-example-2"),
			RemoteCloudReference:    terraformValueObjects.RemoteCloudReference("dragondrop-example-2"),
		},
	}
	assert.Equal(t, expectedDivisionToResourceImportMap, divisionToResourceImportMap)
}

func TestImportMigrationGenerator_mapResourcesToImportLocation_GCP_DifferentResources(t *testing.T) {
	// Given
	stateFileContent := []byte(`{
		"version": 4,
		"terraform_version": "1.2.6",
		"serial": 2,
		"resources": [
			{
				"name": "tfer--dragondrop-example",
				"type": "google_storage_bucket",
				"instances": [
					{
						"attributes_flat": {
							"id": "tfer--dragondrop-example",
							"project": "example-project",
							"name": "dragondrop-example-2"
						}
					}
				]
			},
			{
				"name": "tfer--dragondrop-iam-example",
				"type": "google_storage_bucket_iam_policy",
				"instances": [
					{
						"attributes_flat": {
							"id": "tfer--dragondrop-iam-example",
							"project": "example-project",
							"bucket": "dragondrop-example"
						}
					}
				]
			}
		]
	}`)
	provider := terraformValueObjects.Provider("google")
	i := TerraformImportMigrationGenerator{}

	// When
	divisionToResourceImportMap, err := i.mapResourcesToImportLocation(provider, stateFileContent)

	// Then
	assert.Nil(t, err)

	expectedDivisionToResourceImportMap := terraformValueObjects.ResourceImportMap{
		terraformValueObjects.ResourceName("google_storage_bucket.tfer--dragondrop-example"): terraformValueObjects.ImportMigration{
			TerraformConfigLocation: terraformValueObjects.TerraformConfigLocation("google_storage_bucket.tfer--dragondrop-example"),
			RemoteCloudReference:    terraformValueObjects.RemoteCloudReference("example-project/dragondrop-example-2"),
		},
		terraformValueObjects.ResourceName("google_storage_bucket_iam_policy.tfer--dragondrop-iam-example"): terraformValueObjects.ImportMigration{
			TerraformConfigLocation: terraformValueObjects.TerraformConfigLocation("google_storage_bucket_iam_policy.tfer--dragondrop-iam-example"),
			RemoteCloudReference:    terraformValueObjects.RemoteCloudReference(`b/dragondrop-example`),
		},
	}
	assert.Equal(t, expectedDivisionToResourceImportMap, divisionToResourceImportMap)
}
