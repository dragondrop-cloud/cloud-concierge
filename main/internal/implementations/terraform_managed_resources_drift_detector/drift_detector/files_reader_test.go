package driftDetector

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseStateFile_Valid(t *testing.T) {
	// Given
	stateFileContent := []byte(`{
		"resources": [
			{
				"module": "module.s3-persistent-storage",
				"type": "aws_instance",
				"name": "example",
				"provider": "provider[\"registry.terraform.io/hashicorp/aws\"]",
				"instances": [
					{
						"attributes": {
							"ami": "ami-0c94855ba95b798c7",
							"instance_type": "t2.micro",
							"tags": {
								"Name": "example-instance"
							}
						}
					}
				]
			}
		]
	}`)

	// When
	stateFile, err := ParseTerraformerStateFile(stateFileContent)

	// Then
	require.NoError(t, err)
	require.Len(t, stateFile.Resources, 1)

	resource := stateFile.Resources[0]
	assert.Equal(t, "module.s3-persistent-storage", resource.Module)
	assert.Equal(t, "aws_instance", resource.Type)
	assert.Equal(t, "example", resource.Name)
	assert.Equal(t, `provider["registry.terraform.io/hashicorp/aws"]`, resource.Provider)
	require.Len(t, resource.Instances, 1)
}

func TestParseStateFile_S3Instance(t *testing.T) {
	// Given
	stateFileContent := []byte(`{
		"resources": [
			{
				"module": "module.s3-persistent-storage",
				"type": "aws_s3_bucket",
				"name": "example",
				"provider": "provider[\"registry.terraform.io/hashicorp/aws\"]",
				"instances": [
					{
						"attributes": {
							"bucket": "my-example-bucket",
							"acl": "private"
						}
					}
				]
			}
		]
	}`)

	// When
	stateFile, err := ParseTerraformerStateFile(stateFileContent)

	// Then
	require.NoError(t, err)
	require.Len(t, stateFile.Resources, 1)

	resource := stateFile.Resources[0]
	assert.Equal(t, "module.s3-persistent-storage", resource.Module)
	assert.Equal(t, "aws_s3_bucket", resource.Type)
	assert.Equal(t, "example", resource.Name)
	assert.Equal(t, `provider["registry.terraform.io/hashicorp/aws"]`, resource.Provider)
	require.Len(t, resource.Instances, 1)
}

func TestParseStateFile_EC2Instance(t *testing.T) {
	// Given
	stateFileContent := []byte(`{
		"resources": [
			{
				"module": "module.s3-persistent-storage",
				"type": "aws_instance",
				"name": "example",
				"provider": "provider[\"registry.terraform.io/hashicorp/aws\"]",
				"instances": [
					{
						"attributes": {
							"ami": "ami-0c94855ba95c12345",
							"instance_type": "t2.micro",
							"tags": {
								"Name": "example-instance"
							}
						}
					}
				]
			}
		]
	}`)

	// When
	stateFile, err := ParseTerraformerStateFile(stateFileContent)

	// Then
	require.NoError(t, err)
	require.Len(t, stateFile.Resources, 1)

	resource := stateFile.Resources[0]
	assert.Equal(t, "module.s3-persistent-storage", resource.Module)
	assert.Equal(t, "aws_instance", resource.Type)
	assert.Equal(t, "example", resource.Name)
	assert.Equal(t, `provider["registry.terraform.io/hashicorp/aws"]`, resource.Provider)
	require.Len(t, resource.Instances, 1)
}

func TestParseStateFile_ResourceWithoutModule(t *testing.T) {
	// Given
	stateFileContent := []byte(`{
		"resources": [
			{
				"type": "aws_instance",
				"name": "example",
				"provider": "provider[\"registry.terraform.io/hashicorp/aws\"]",
				"instances": [
					{
						"attributes": {
							"ami": "ami-0c94855ba95c12345",
							"instance_type": "t2.micro",
							"tags": {
								"Name": "example-instance"
							}
						}
					}
				]
			}
		]
	}`)

	// When
	stateFile, err := ParseTerraformerStateFile(stateFileContent)

	// Then
	require.NoError(t, err)
	require.Len(t, stateFile.Resources, 1)

	resource := stateFile.Resources[0]
	assert.Equal(t, "root", resource.Module)
	assert.Equal(t, "aws_instance", resource.Type)
	assert.Equal(t, "example", resource.Name)
	assert.Equal(t, `provider["registry.terraform.io/hashicorp/aws"]`, resource.Provider)
	require.Len(t, resource.Instances, 1)
}

func TestParseStateFile_MultipleResources(t *testing.T) {
	// Given
	stateFileContent := []byte(`{
		"resources": [
			{
				"module": "module.s3-persistent-storage",
				"type": "aws_instance",
				"name": "example",
				"provider": "provider[\"registry.terraform.io/hashicorp/aws\"]",
				"instances": [
					{
						"attributes": {
							"ami": "ami-0c94855ba95b798c7",
							"instance_type": "t2.micro",
							"tags": {
								"Name": "example-instance"
							}
						}
					}
				]
			},
			{
				"module": "module.s3-persistent-storage",
				"type": "aws_instance",
				"name": "example2",
				"provider": "provider[\"registry.terraform.io/hashicorp/aws\"]",
				"instances": [
					{
						"attributes": {
							"ami": "ami-0c94855ba95b798c7",
							"instance_type": "t2.micro",
							"tags": {
								"Name": "example-instance2"
							}
						}
					}
				]
			},
			{
				"module": "module.s3-persistent-storage",
				"type": "aws_instance",
				"name": "example3",
				"provider": "provider[\"registry.terraform.io/hashicorp/aws\"]",
				"instances": [
					{
						"attributes": {
							"ami": "ami-0c94855ba95b798c7",
							"instance_type": "t2.micro",
							"tags": {
								"Name": "example-instance3"
							}
						}
					}
				]
			}
		]
	}`)

	detector := ManagedResourcesDriftDetector{}

	// When
	stateFile, err := detector.parseRemoteStateFile(stateFileContent)

	// Then
	require.NoError(t, err)
	require.Len(t, stateFile.Resources, 3)

	resource := stateFile.Resources[0]
	assert.Equal(t, "module.s3-persistent-storage", resource.Module)
	assert.Equal(t, "aws_instance", resource.Type)
	assert.Equal(t, "example", resource.Name)
	assert.Equal(t, `provider["registry.terraform.io/hashicorp/aws"]`, resource.Provider)
	require.Len(t, resource.Instances, 1)

	secondResource := stateFile.Resources[1]
	assert.Equal(t, "module.s3-persistent-storage", secondResource.Module)
	assert.Equal(t, "aws_instance", secondResource.Type)
	assert.Equal(t, "example2", secondResource.Name)
	assert.Equal(t, `provider["registry.terraform.io/hashicorp/aws"]`, secondResource.Provider)
	require.Len(t, secondResource.Instances, 1)

	thirdResource := stateFile.Resources[2]
	assert.Equal(t, "module.s3-persistent-storage", thirdResource.Module)
	assert.Equal(t, "aws_instance", thirdResource.Type)
	assert.Equal(t, "example3", thirdResource.Name)
	assert.Equal(t, `provider["registry.terraform.io/hashicorp/aws"]`, thirdResource.Provider)
	require.Len(t, thirdResource.Instances, 1)
}

func TestParseStateFile_MultipleInstancesWithinASingleResource(t *testing.T) {
	// Given
	stateFileContent := []byte(`{
		"resources": [
			{
				"type": "aws_instance",
				"name": "example",
				"provider": "provider[\"registry.terraform.io/hashicorp/aws\"]",
				"instances": [
					{
						"attributes": {
							"ami": "ami-0c94855ba95c12345",
							"instance_type": "t2.micro",
							"tags": {
								"Name": "example-instance"
							}
						}
					},
					{
						"attributes": {
							"ami": "ami-0c94855ba95c12345",
							"instance_type": "t2.micro",
							"tags": {
								"Name": "example-instance1"
							}
						}
					},
					{
						"attributes": {
							"ami": "ami-0c94855ba95c12345",
							"instance_type": "t2.micro",
							"tags": {
								"Name": "example-instance2"
							}
						}
					}
				]
			}
		]
	}`)

	// When
	stateFile, err := ParseTerraformerStateFile(stateFileContent)

	// Then
	require.NoError(t, err)
	require.Len(t, stateFile.Resources, 1)

	resource := stateFile.Resources[0]
	assert.Equal(t, "root", resource.Module)
	assert.Equal(t, "aws_instance", resource.Type)
	assert.Equal(t, "example", resource.Name)
	assert.Equal(t, `provider["registry.terraform.io/hashicorp/aws"]`, resource.Provider)
	require.Len(t, resource.Instances, 3)
}
