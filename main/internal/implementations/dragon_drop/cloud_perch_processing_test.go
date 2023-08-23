package dragonDrop

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/stretchr/testify/require"
)

func writeFile(fileName string, content []byte, directoryName string) error {
	err := os.MkdirAll(directoryName, os.ModePerm)
	if err != nil {
		return err
	}
	f, err := os.Create(fmt.Sprintf("%s/%v", directoryName, fileName))
	if err != nil {
		return fmt.Errorf("[error writing output file]%w", err)
	}

	defer f.Close()
	_, err = f.Write(content)
	if err != nil {
		return fmt.Errorf("[error writing output file]%w", err)
	}

	return nil
}

func writeOutputFile(filename string, content []byte) error {
	return writeFile(filename, content, "outputs")
}

func writeCurrentCloudFile(filename string, content []byte) error {
	return writeFile(filename, content, "current_cloud")
}

func cleanMockedDirectories(t *testing.T) {
	cmd := exec.Command("rm", "-rf", "outputs")
	err := cmd.Run()
	if err != nil {
		t.Logf("Error cleaning mocked directories outputs: %v", err)
	}

	cmd = exec.Command("rm", "-rf", "current_cloud")
	err = cmd.Run()
	if err != nil {
		t.Logf("Error cleaning mocked directories current_cloud: %v", err)
	}

	cmd = exec.Command("rm", "-rf", "repo")
	err = cmd.Run()
	if err != nil {
		t.Logf("Error cleaning mocked directories repo: %v", err)
	}
}

func TestHTTPDragonDropClient_getResourceInventoryData(t *testing.T) {
	defer cleanMockedDirectories(t)

	// Given
	ctx := context.Background()
	client := HTTPDragonDropClient{}

	err := writeOutputFile("new-resources.json", []byte(`
		{
 			"resource_id_1": {
				"ResourceType": "resource_type_1",
				"ResourceTerraformerName": "resource_terraformer_name_1",
				"Region": "region_1"
			},
			"resource_id_2": {
				"ResourceType": "resource_type_1",
				"ResourceTerraformerName": "resource_terraformer_name_2",
				"Region": "region_1"
			}
		}
	`))
	require.NoError(t, err)

	err = writeOutputFile("drift-resources-differences.json", []byte(`
		[
 			{
				"RecentActor": "root",
				"RecentActionTimestamp": "2023-08-09",
				"AttributeName": "access_logs.s3.enabled",
				"TerraformValue": "false",
				"CloudValue": "true",
				"InstanceID": "instance_id_2",
				"InstanceRegion": "us-east-1",
				"StateFileName": "aws-networking-dev",
				"ModuleName": "root",
				"ResourceType": "aws_lb",
				"ResourceName": "example_lb"
			},
			{
				"RecentActor": "root",
				"RecentActionTimestamp": "2023-08-09",
				"AttributeName": "access_logs.s3.enabled",
				"TerraformValue": "false",
				"CloudValue": "true",
				"InstanceID": "i-instance_id_1",
				"InstanceRegion": "us-east-1",
				"StateFileName": "aws-networking-dev",
				"ModuleName": "root",
				"ResourceType": "aws_lb",
				"ResourceName": "example_lb_2"
			}
		]
	`))
	require.NoError(t, err)

	// When
	resourceInventory, newResources, err := client.getResourceInventoryData(ctx)

	// Then
	require.NoError(t, err)
	require.Equal(t, 2, resourceInventory.DriftedResources)
	require.Equal(t, 2, resourceInventory.ResourcesOutsideTerraformControl)
	require.NotNil(t, newResources["resource_id_1"])
	require.NotNil(t, newResources["resource_id_2"])

}

func TestHTTPDragonDropClient_getCloudSecurityData(t *testing.T) {
	// Given
	defer cleanMockedDirectories(t)

	ctx := context.Background()
	client := HTTPDragonDropClient{}

	err := writeOutputFile("security-scan.json", []byte(`
	{
	  "results": [
			{
			  "id": "arn:aws:elasticloadbalancing:us-east-1:682649898103:loadbalancer/app/tf-managed-demo-alb/4c89e21113613302",
			  "rule_id": "AVD-AWS-0053",
			  "description": "Load balancer is exposed publicly.",
			  "severity": "HIGH"
			},
			{
			  "id": "arn:aws:elasticloadbalancing:us-east-1:682649898103:loadbalancer/app/tf-managed-demo-alb/4c89e21113613302",
			  "rule_id": "AVD-AWS-0053",
			  "description": "Load balancer is exposed publicly.",
			  "severity": "CRITICAL"
			},
			{
			  "id": "arn:aws:elasticloadbalancing:us-east-1:682649898103:loadbalancer/app/tf-managed-demo-alb/4c89e21113613302",
			  "rule_id": "AVD-AWS-0053",
			  "description": "Load balancer is exposed publicly.",
			  "severity": "MEDIUM"
			},
			{
			  "id": "arn:aws:elasticloadbalancing:us-east-1:682649898103:loadbalancer/app/tf-managed-demo-alb/4c89e21113613302",
			  "rule_id": "AVD-AWS-0053",
			  "description": "Load balancer is exposed publicly.",
			  "severity": "LOW"
			}
		]
   }`))
	require.NoError(t, err)

	// When
	cloudSecurity, err := client.getCloudSecurityData(ctx)

	// Then
	require.NoError(t, err)
	require.Equal(t, 1, cloudSecurity.SecurityRiskCritical)
	require.Equal(t, 1, cloudSecurity.SecurityRiskHigh)
	require.Equal(t, 1, cloudSecurity.SecurityRiskMedium)
	require.Equal(t, 1, cloudSecurity.SecurityRiskLow)
}

func TestHTTPDragonDropClient_getCloudCostsData(t *testing.T) {
	// Given
	defer cleanMockedDirectories(t)

	ctx := context.Background()
	client := HTTPDragonDropClient{}

	newResources := map[string]interface{}{
		"resource_type_1.resource_terraformer_name_1": interface{}(
			map[string]interface{}{
				"ResourceType":            "resource_type_1",
				"ResourceTerraformerName": "resource_terraformer_name_1",
				"Region":                  "region_1",
			},
		),
		"resource_type_2.resource_terraformer_name_2": interface{}(
			map[string]interface{}{
				"ResourceType":            "resource_type_2",
				"ResourceTerraformerName": "resource_terraformer_name_2",
				"Region":                  "region_1",
			},
		),
	}

	err := writeOutputFile("cost-estimates.json", []byte(`
	[
		{
			"cost_component": "Application load balancer",
			"is_usage_based": false,
			"monthly_cost": "16.425",
			"monthly_quantity": "730",
			"price": "0.0225",
			"resource_name": "resource_type_1.resource_terraformer_name_1",
			"sub_resource_name": "",
			"unit": "hours"
		},
		{
			"cost_component": "Application load balancer",
			"is_usage_based": false,
			"monthly_cost": "",
			"monthly_quantity": "",
			"price": "5.84",
			"resource_name": "resource_type_2.resource_terraformer_name_2",
			"sub_resource_name": "",
			"unit": "hours"
		},
		{
			"cost_component": "Application load balancer",
			"is_usage_based": false,
			"monthly_cost": "112.746",
			"monthly_quantity": "730",
			"price": "0.1542",
			"resource_name": "resource_type_3.resource_terraformer_name_3",
			"sub_resource_name": "",
			"unit": "hours"
		},
		{
			"cost_component": "Application load balancer",
			"is_usage_based": false,
			"monthly_cost": "",
			"monthly_quantity": "",
			"price": "10.8",
			"resource_name": "resource_type_4.resource_terraformer_name_4",
			"sub_resource_name": "",
			"unit": "hours"
		}
	]`))
	require.NoError(t, err)

	// When
	cloudCosts, err := client.getCloudCostsData(ctx, newResources)

	// Then
	require.NoError(t, err)
	require.Equal(t, 22.265, cloudCosts.CostsOutsideOfTerraform)
	require.Equal(t, 123.546, cloudCosts.CostsTerraformControlled)
}

func Test_formatResources(t *testing.T) {
	// Given
	defer cleanMockedDirectories(t)

	newResources := map[string]interface{}{
		"instance_id_1": interface{}(
			map[string]interface{}{
				"ResourceType":            "resource_type_1",
				"ResourceTerraformerName": "resource_terraformer_name_1",
				"Region":                  "region_1",
			},
		),
		"instance_id_2": interface{}(
			map[string]interface{}{
				"ResourceType":            "resource_type_2",
				"ResourceTerraformerName": "resource_terraformer_name_2",
				"Region":                  "region_1",
			},
		),
	}

	// When
	resourcesFormatted := formatResources(newResources)

	// Then
	expectedResourcesFormatted := map[string]interface{}{
		"resource_type_1.resource_terraformer_name_1": interface{}(
			map[string]interface{}{
				"ResourceType":            "resource_type_1",
				"ResourceTerraformerName": "resource_terraformer_name_1",
				"Region":                  "region_1",
			},
		),
		"resource_type_2.resource_terraformer_name_2": interface{}(
			map[string]interface{}{
				"ResourceType":            "resource_type_2",
				"ResourceTerraformerName": "resource_terraformer_name_2",
				"Region":                  "region_1",
			},
		),
	}
	require.Equal(t, expectedResourcesFormatted, resourcesFormatted)
}

func TestHTTPDragonDropClient_getTerraformFootprintData_main_tf_no_modules_aws(t *testing.T) {
	// Given
	defer cleanMockedDirectories(t)

	ctx := context.Background()
	client := HTTPDragonDropClient{}

	err := writeCurrentCloudFile("main.tf", []byte(`
		terraform {
		  required_version = "1.5.1"

		  required_providers {
			aws = {
			  source  = "hashicorp/aws"
			  version = "~>4.59.0"
			}

		  }
		}
	`))
	require.NoError(t, err)

	// When
	terraformFootprint, err := client.getTerraformFootprintData(ctx)

	// Then
	require.NoError(t, err)
	require.Equal(t, `{"hashicorp/aws":{"4.59.0":1}}`, terraformFootprint.VersionsTFProviders)
	require.Equal(t, `{"1.5.1":1}`, terraformFootprint.VersionsTF)
	require.Equal(t, "{}", terraformFootprint.VersionsTFModules)
}

func TestHTTPDragonDropClient_getTerraformFootprintData_main_tf_no_modules_azure(t *testing.T) {
	// Given
	defer cleanMockedDirectories(t)

	ctx := context.Background()
	client := HTTPDragonDropClient{}

	err := writeCurrentCloudFile("main.tf", []byte(`
		terraform {
		  required_version = "1.5.1"

		  required_providers {
			azurerm = {
			  source  = "hashicorp/azurerm"
			  version = "~>3.30.0"
			}
		  }
		}
	`))
	require.NoError(t, err)

	// When
	terraformFootprint, err := client.getTerraformFootprintData(ctx)

	// Then
	require.NoError(t, err)
	require.Equal(t, `{"hashicorp/azurerm":{"3.30.0":1}}`, terraformFootprint.VersionsTFProviders)
	require.Equal(t, `{"1.5.1":1}`, terraformFootprint.VersionsTF)
	require.Equal(t, "{}", terraformFootprint.VersionsTFModules)
}

func TestHTTPDragonDropClient_getTerraformFootprintData_main_tf_no_modules_google(t *testing.T) {
	// Given
	defer cleanMockedDirectories(t)

	ctx := context.Background()
	client := HTTPDragonDropClient{}

	err := writeCurrentCloudFile("main.tf", []byte(`
		terraform {
		  required_version = "1.5.1"

		  required_providers {
			google = {
			  source  = "hashicorp/google"
			  version = "~>4.77.0"
			}
		  }
		}
	`))
	require.NoError(t, err)

	// When
	terraformFootprint, err := client.getTerraformFootprintData(ctx)

	// Then
	require.NoError(t, err)
	require.Equal(t, `{"hashicorp/google":{"4.77.0":1}}`, terraformFootprint.VersionsTFProviders)
	require.Equal(t, `{"1.5.1":1}`, terraformFootprint.VersionsTF)
	require.Equal(t, "{}", terraformFootprint.VersionsTFModules)
}

func TestHTTPDragonDropClient_getTerraformFootprintData_versions_tf_no_modules(t *testing.T) {
	// Given
	defer cleanMockedDirectories(t)

	ctx := context.Background()
	client := HTTPDragonDropClient{}

	err := writeCurrentCloudFile("versions.tf", []byte(`
		terraform {
		  required_version = "1.5.1"

		  required_providers {
			google = {
			  source  = "hashicorp/google"
			  version = "~>4.77.0"
			}
		  }
		}
	`))
	require.NoError(t, err)

	// When
	terraformFootprint, err := client.getTerraformFootprintData(ctx)

	// Then
	require.NoError(t, err)
	require.Equal(t, `{"hashicorp/google":{"4.77.0":1}}`, terraformFootprint.VersionsTFProviders)
	require.Equal(t, `{"1.5.1":1}`, terraformFootprint.VersionsTF)
	require.Equal(t, "{}", terraformFootprint.VersionsTFModules)
}

func TestHTTPDragonDropClient_getTerraformFootprintData_another_tf_file_no_modules(t *testing.T) {
	// Given
	defer cleanMockedDirectories(t)

	ctx := context.Background()
	client := HTTPDragonDropClient{
		config: HTTPDragonDropClientConfig{
			WorkspaceDirectories: []string{"workspace1"},
		},
	}

	err := writeFile("another.tf", []byte(`
		terraform {
		  required_version = "1.5.1"

		  required_providers {
			google = {
			  source  = "hashicorp/google"
			  version = "~>4.77.0"
			}
		  }
		}
	`), "repo/workspace1")
	require.NoError(t, err)

	// When
	terraformFootprint, err := client.getTerraformFootprintData(ctx)

	// Then
	require.NoError(t, err)
	require.Equal(t, `{"hashicorp/google":{"4.77.0":1}}`, terraformFootprint.VersionsTFProviders)
	require.Equal(t, `{"1.5.1":1}`, terraformFootprint.VersionsTF)
	require.Equal(t, "{}", terraformFootprint.VersionsTFModules)
}

func TestHTTPDragonDropClient_getTerraformFootprintData_another_tf_file_one_module(t *testing.T) {
	// Given
	defer cleanMockedDirectories(t)

	ctx := context.Background()
	client := HTTPDragonDropClient{
		config: HTTPDragonDropClientConfig{
			WorkspaceDirectories: []string{"workspace1"},
		},
	}

	err := writeFile("another.tf", []byte(`
		terraform {
		  required_version = "1.5.1"

		  required_providers {
			google = {
			  source  = "hashicorp/google"
			  version = "~>4.77.0"
			}
		  }
		}
	`), "repo/workspace1")
	require.NoError(t, err)

	err = writeFile("module_file.tf", []byte(`
		module "servers" {
			source = "./app-cluster"
			version = "0.0.5"

			servers = 5
		}
	`), "repo/workspace1")
	require.NoError(t, err)

	// When
	terraformFootprint, err := client.getTerraformFootprintData(ctx)

	// Then
	require.NoError(t, err)
	require.Equal(t, `{"hashicorp/google":{"4.77.0":1}}`, terraformFootprint.VersionsTFProviders)
	require.Equal(t, `{"1.5.1":1}`, terraformFootprint.VersionsTF)
	require.Equal(t, `{"./app-cluster":{"0.0.5":1}}`, terraformFootprint.VersionsTFModules)
}

func TestHTTPDragonDropClient_getTerraformFootprintData_another_tf_file_two_modules_in_a_file(t *testing.T) {
	// Given
	defer cleanMockedDirectories(t)

	ctx := context.Background()
	client := HTTPDragonDropClient{
		config: HTTPDragonDropClientConfig{
			WorkspaceDirectories: []string{"workspace1"},
		},
	}

	err := writeFile("another.tf", []byte(`
		terraform {
		  required_version = "1.5.1"

		  required_providers {
			google = {
			  source  = "hashicorp/google"
			  version = "~>4.77.0"
			}
		  }
		}
	`), "repo/workspace1")
	require.NoError(t, err)

	err = writeFile("module_file.tf", []byte(`
		module "servers" {
			source = "./app-cluster"
			version = "0.0.5"

			servers = 5
		}

		module "servers2" {
			source = "./app-cluster"
			version = "0.0.6"

			servers = 5
		}
	`), "repo/workspace1")
	require.NoError(t, err)

	// When
	terraformFootprint, err := client.getTerraformFootprintData(ctx)

	// Then
	require.NoError(t, err)
	require.Equal(t, `{"hashicorp/google":{"4.77.0":1}}`, terraformFootprint.VersionsTFProviders)
	require.Equal(t, `{"1.5.1":1}`, terraformFootprint.VersionsTF)
	require.Equal(t, `{"./app-cluster":{"0.0.5":1,"0.0.6":1}}`, terraformFootprint.VersionsTFModules)
}

func TestHTTPDragonDropClient_getTerraformFootprintData_another_tf_file_multiple_modules_multiple_files(t *testing.T) {
	// Given
	defer cleanMockedDirectories(t)

	ctx := context.Background()
	client := HTTPDragonDropClient{
		config: HTTPDragonDropClientConfig{
			WorkspaceDirectories: []string{"workspace1"},
		},
	}

	err := writeFile("another.tf", []byte(`
		terraform {
		  required_version = "1.5.1"

		  required_providers {
			google = {
			  source  = "hashicorp/google"
			  version = "~>4.77.0"
			}
		  }
		}
	`), "repo/workspace1")
	require.NoError(t, err)

	err = writeFile("module_file1.tf", []byte(`
		module "servers" {
			source = "./app-cluster"
			version = "0.0.5"

			servers = 5
		}

		module "servers2" {
			source = "./app-cluster"
			version = "0.0.6"

			servers = 5
		}
	`), "repo/workspace1")
	require.NoError(t, err)

	err = writeFile("module_file2.tf", []byte(`
		module "server4" {
			source = "./app-cluster"
			version = "0.0.7"

			servers = 5
		}

		module "servers3" {
			source = "./app-cluster"
			version = "0.0.8"

			servers = 5
		}
	`), "repo/workspace1")
	require.NoError(t, err)

	// When
	terraformFootprint, err := client.getTerraformFootprintData(ctx)

	// Then
	require.NoError(t, err)
	require.Equal(t, `{"hashicorp/google":{"4.77.0":1}}`, terraformFootprint.VersionsTFProviders)
	require.Equal(t, `{"1.5.1":1}`, terraformFootprint.VersionsTF)
	require.Equal(t, `{"./app-cluster":{"0.0.5":1,"0.0.6":1,"0.0.7":1,"0.0.8":1}}`, terraformFootprint.VersionsTFModules)
}

func TestHTTPDragonDropClient_getTerraformFootprintData_another_tf_file_sum_module_versions(t *testing.T) {
	// Given
	defer cleanMockedDirectories(t)

	ctx := context.Background()
	client := HTTPDragonDropClient{
		config: HTTPDragonDropClientConfig{
			WorkspaceDirectories: []string{"workspace1"},
		},
	}

	err := writeFile("another.tf", []byte(`
		terraform {
		  required_version = "1.5.1"

		  required_providers {
			google = {
			  source  = "hashicorp/google"
			  version = "~>4.77.0"
			}
		  }
		}
	`), "repo/workspace1")
	require.NoError(t, err)

	err = writeFile("module_file1.tf", []byte(`
		module "servers" {
			source = "./app-cluster"
			version = "0.0.5"

			servers = 5
		}

		module "servers2" {
			source = "./app-cluster"
			version = "0.0.6"

			servers = 5
		}
	`), "repo/workspace1")
	require.NoError(t, err)

	err = writeFile("module_file2.tf", []byte(`
		module "server4" {
			source = "./app-cluster"
			version = "0.0.5"

			servers = 5
		}

		module "servers3" {
			source = "./app-cluster"
			version = "0.0.6"

			servers = 5
		}
	`), "repo/workspace1")
	require.NoError(t, err)

	// When
	terraformFootprint, err := client.getTerraformFootprintData(ctx)

	// Then
	require.NoError(t, err)
	require.Equal(t, `{"hashicorp/google":{"4.77.0":1}}`, terraformFootprint.VersionsTFProviders)
	require.Equal(t, `{"1.5.1":1}`, terraformFootprint.VersionsTF)
	require.Equal(t, `{"./app-cluster":{"0.0.5":2,"0.0.6":2}}`, terraformFootprint.VersionsTFModules)
}

func TestHTTPDragonDropClient_getTerraformFootprintData_another_tf_file_more_than_one_source(t *testing.T) {
	// Given
	defer cleanMockedDirectories(t)

	ctx := context.Background()
	client := HTTPDragonDropClient{
		config: HTTPDragonDropClientConfig{
			WorkspaceDirectories: []string{"workspace1"},
		},
	}

	err := writeFile("another.tf", []byte(`
		terraform {
		  required_version = "1.5.1"

		  required_providers {
			google = {
			  source  = "hashicorp/google"
			  version = "~>4.77.0"
			}
		  }
		}
	`), "repo/workspace1")
	require.NoError(t, err)

	err = writeFile("module_file1.tf", []byte(`
		module "servers" {
			source = "./app-cluster"
			version = "0.0.5"

			servers = 5
		}

		module "servers2" {
			source = "./app-cluster2"
			version = "0.0.6"

			servers = 5
		}
	`), "repo/workspace1")
	require.NoError(t, err)

	err = writeFile("module_file2.tf", []byte(`
		module "server4" {
			source = "./app-cluster"
			version = "0.0.5"

			servers = 5
		}

		module "servers3" {
			source = "./app-cluster2"
			version = "0.0.6"

			servers = 5
		}
	`), "repo/workspace1")
	require.NoError(t, err)

	// When
	terraformFootprint, err := client.getTerraformFootprintData(ctx)

	// Then
	require.NoError(t, err)
	require.Equal(t, `{"hashicorp/google":{"4.77.0":1}}`, terraformFootprint.VersionsTFProviders)
	require.Equal(t, `{"1.5.1":1}`, terraformFootprint.VersionsTF)
	require.Equal(t, `{"./app-cluster":{"0.0.5":2},"./app-cluster2":{"0.0.6":2}}`, terraformFootprint.VersionsTFModules)
}

func TestHTTPDragonDropClient_getTerraformFootprintData_another_tf_file_more_than_one_source_multiple_versions(t *testing.T) {
	// Given
	defer cleanMockedDirectories(t)

	ctx := context.Background()
	client := HTTPDragonDropClient{
		config: HTTPDragonDropClientConfig{
			WorkspaceDirectories: []string{"workspace1"},
		},
	}

	err := writeFile("another.tf", []byte(`
		terraform {
		  required_version = "1.5.1"

		  required_providers {
			google = {
			  source  = "hashicorp/google"
			  version = "~>4.77.0"
			}
		  }
		}
	`), "repo/workspace1")
	require.NoError(t, err)

	err = writeFile("module_file1.tf", []byte(`
		module "servers" {
			source = "./app-cluster"
			version = "0.0.5"

			servers = 5
		}

		module "servers2" {
			source = "./app-cluster"
			version = "0.0.6"

			servers = 5
		}

		module "servers2" {
			source = "./app-cluster2"
			version = "0.0.5"

			servers = 5
		}

		module "servers2" {
			source = "./app-cluster2"
			version = "0.0.6"

			servers = 5
		}
	`), "repo/workspace1")
	require.NoError(t, err)

	err = writeFile("module_file2.tf", []byte(`
		module "server4" {
			source = "./app-cluster"
			version = "0.0.5"

			servers = 5
		}

		module "servers3" {
			source = "./app-cluster"
			version = "0.0.6"

			servers = 5
		}

		module "servers2" {
			source = "./app-cluster2"
			version = "0.0.5"

			servers = 5
		}

		module "servers2" {
			source = "./app-cluster2"
			version = "0.0.6"

			servers = 5
		}
	`), "repo/workspace1")
	require.NoError(t, err)

	// When
	terraformFootprint, err := client.getTerraformFootprintData(ctx)

	// Then
	require.NoError(t, err)
	require.Equal(t, `{"hashicorp/google":{"4.77.0":1}}`, terraformFootprint.VersionsTFProviders)
	require.Equal(t, `{"1.5.1":1}`, terraformFootprint.VersionsTF)
	require.Equal(t, `{"./app-cluster":{"0.0.5":2,"0.0.6":2},"./app-cluster2":{"0.0.5":2,"0.0.6":2}}`, terraformFootprint.VersionsTFModules)
}

func Test_getVersionFromProviderAttribute_nil_attribute(t *testing.T) {
	// When
	version := getVersionFromProviderAttribute(nil)

	// Then
	require.Equal(t, "", version)
}

func Test_getVersionFromProviderAttribute_tilde_greater_than(t *testing.T) {
	// Given
	terraformBlock := []byte(`
	required_providers {
			google = {
			  source  = "hashicorp/google"
			  version = "~>4.77.0"
			}
		}
	`)
	terraformFile, _ := hclwrite.ParseConfig(
		terraformBlock,
		"placeholder.tf",
		hcl.Pos{Line: 0, Column: 0, Byte: 0},
	)
	requiredVersionAttribute := terraformFile.Body().FirstMatchingBlock("required_providers", nil)
	providerAttribute := requiredVersionAttribute.Body().GetAttribute("google")

	// When
	version := getVersionFromProviderAttribute(providerAttribute)

	// Then
	require.Equal(t, "4.77.0", version)
}

func Test_getVersionFromProviderAttribute_equals(t *testing.T) {
	// Given
	terraformBlock := []byte(`
	required_providers {
			google = {
			  source  = "hashicorp/google"
			  version = "=4.77.0"
			}
		}
	`)
	terraformFile, _ := hclwrite.ParseConfig(
		terraformBlock,
		"placeholder.tf",
		hcl.Pos{Line: 0, Column: 0, Byte: 0},
	)
	requiredVersionAttribute := terraformFile.Body().FirstMatchingBlock("required_providers", nil)
	providerAttribute := requiredVersionAttribute.Body().GetAttribute("google")

	// When
	version := getVersionFromProviderAttribute(providerAttribute)

	// Then
	require.Equal(t, "4.77.0", version)
}

func Test_getVersionFromProviderAttribute_greater_or_equals_than(t *testing.T) {
	// Given
	terraformBlock := []byte(`
	required_providers {
			google = {
			  source  = "hashicorp/google"
			  version = ">=4.77.0"
			}
		}
	`)
	terraformFile, _ := hclwrite.ParseConfig(
		terraformBlock,
		"placeholder.tf",
		hcl.Pos{Line: 0, Column: 0, Byte: 0},
	)
	requiredVersionAttribute := terraformFile.Body().FirstMatchingBlock("required_providers", nil)
	providerAttribute := requiredVersionAttribute.Body().GetAttribute("google")

	// When
	version := getVersionFromProviderAttribute(providerAttribute)

	// Then
	require.Equal(t, "4.77.0", version)
}

func Test_concatenateVersions(t *testing.T) {
	// Given
	versions1 := ModulesVersions{
		"source1": {
			"1.0.0": 1,
			"1.0.1": 2,
			"1.0.2": 3,
		},
		"source2": {
			"1.0.0": 1,
			"1.0.1": 2,
			"1.0.2": 3,
		},
	}
	versions2 := ModulesVersions{
		"source1": {
			"1.0.1": 2,
			"1.0.2": 3,
			"1.0.3": 4,
		},
	}

	// When
	versions := concatenateVersions(versions1, versions2)

	// Then
	require.Equal(t, ModulesVersions{
		"source1": {
			"1.0.0": 1,
			"1.0.1": 4,
			"1.0.2": 6,
			"1.0.3": 4,
		},
		"source2": {
			"1.0.0": 1,
			"1.0.1": 2,
			"1.0.2": 3,
		},
	}, versions)
}

func Test_getAttributeValue(t *testing.T) {
	// Given
	sourceAttribute := []byte(`source = "hashicorp/google"`)
	versionAttribute := []byte(`version = ">=4.77.0"`)

	// When
	sourceValue, err := getAttributeValue(sourceAttribute)
	require.NoError(t, err)

	versionValue, err := getAttributeValue(versionAttribute)
	require.NoError(t, err)

	// Then
	require.Equal(t, "hashicorp/google", sourceValue)
	require.Equal(t, ">=4.77.0", versionValue)
}

func Test_getUniqueDriftedResourceCount(t *testing.T) {
	// Given
	terraformState := []byte(`[
 {
   "InstanceID": "arn:aws:elasticloadbalancing:us-east-1:6"
 },
 {
   "InstanceID": "arn:aws:elasticloadbalancing:us-east-1:6"
 },
 {
   "InstanceID": "arn:aws:elasticloadbalancing:us-east-1:7"
 },
 {
   "InstanceID": "arn:aws:elasticloadbalancing:us-east-1:7"
 },
 {
   "InstanceID": "arn:aws:elasticloadbalancing:us-east-1:6"
 }
]`)
	var data []interface{}
	err := json.Unmarshal(terraformState, &data)
	if err != nil {
		t.Errorf("unexpected error unmarshalling value: %s", err)
	}

	// When
	count := getUniqueDriftedResourceCount(data)

	// Then
	if count != 2 {
		t.Errorf("expected 2, got %d", count)
	}
}
