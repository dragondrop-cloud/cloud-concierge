package documentize

import (
	"reflect"
	"testing"

	"github.com/Jeffail/gabs/v2"
	"github.com/stretchr/testify/require"
)

func TestAzureResourceExtractor_ExtractResourceDetails_NotFlat(t *testing.T) {
	are := azureResourceExtractor{
		currentResourceDetails: &azureResourceDetails{},
	}

	tfStateParsed, err := gabs.ParseJSON([]byte(`{
		"resources": [
			{
			  "mode": "managed",
			  "type": "azurerm_resource_group",
			  "name": "drifted-resource-group",
			  "provider": "provider[\"registry.terraform.io/hashicorp/azurerm\"]",
			  "instances": [
				{
				  "schema_version": 0,
				  "attributes": {
					"id": "/subscriptions/123/resourceGroups/drifted-resource-group",
					"location": "eastus2",
					"name": "drifted-resource-group",
					"tags": {
						"Hello": "Drift"
					}
				  },
				  "sensitive_attributes": []
				}
			  ]
			}
		  ]
	}`))

	if err != nil {
		t.Errorf("Unexpected error in gabs.ParseJSON(): %v", err)
	}

	err = are.ExtractResourceDetails(tfStateParsed, false, 0, 0)
	require.NoError(t, err)

	actualResourceDetails := are.GetCurrentResourceDetails()

	expectedResourceDetails := &azureResourceDetails{
		terraformModule:       "none",
		terraformName:         "drifted-resource-group",
		terraformType:         "azurerm_resource_group",
		azureInstanceLocation: "eastus2",
		azureInstanceName:     "drifted-resource-group",
		azureInstanceTags: map[string]string{
			"Hello": "Drift",
		},
	}

	if !reflect.DeepEqual(actualResourceDetails, expectedResourceDetails) {
		t.Errorf("got:\n%v\nexpected:\n%v", actualResourceDetails, expectedResourceDetails)
	}
}

func TestAzureResourceExtractor_ExtractResourceDetails(t *testing.T) {
	are := azureResourceExtractor{
		currentResourceDetails: &azureResourceDetails{},
	}

	tfStateParsed, err := gabs.ParseJSON([]byte(`{
		"resources": [
			{
			  "mode": "managed",
			  "type": "azurerm_resource_group",
			  "name": "tfer--drifted-resource-group",
			  "provider": "provider[\"registry.terraform.io/hashicorp/azurerm\"]",
			  "instances": [
				{
				  "schema_version": 0,
				  "attributes_flat": {
					"id": "/subscriptions/123/resourceGroups/drifted-resource-group",
					"location": "eastus2",
					"name": "drifted-resource-group",
					"tags.%": "1",
					"tags.Hello": "Drift"
				  },
				  "sensitive_attributes": []
				}
			  ]
			}
		  ]
	}`))

	if err != nil {
		t.Errorf("Unexpected error in gabs.ParseJSON(): %v", err)
	}

	err = are.ExtractResourceDetails(tfStateParsed, true, 0, 0)
	require.NoError(t, err)

	actualResourceDetails := are.GetCurrentResourceDetails()

	expectedResourceDetails := &azureResourceDetails{
		terraformModule:       "none",
		terraformName:         "tfer--drifted-resource-group",
		terraformType:         "azurerm_resource_group",
		azureInstanceLocation: "eastus2",
		azureInstanceName:     "tfer--drifted-resource-group",
		azureInstanceTags: map[string]string{
			"Hello": "Drift",
		},
	}

	if !reflect.DeepEqual(actualResourceDetails, expectedResourceDetails) {
		t.Errorf("got:\n%v\nexpected:\n%v", actualResourceDetails, expectedResourceDetails)
	}
}

func TestResourceDetailsToSentence_Azure(t *testing.T) {
	// Base case
	are := azureResourceExtractor{
		currentResourceDetails: &azureResourceDetails{
			terraformModule:       "module.example-module",
			terraformName:         "drifted-resource-group",
			terraformType:         "azurerm_resource_group",
			azureInstanceLocation: "eastus2",
			azureInstanceName:     "drifted-resource-group",
			azureInstanceTags: map[string]string{
				"Hello": "Drift",
			},
		},
		typeToCategory: azureResourceCategories(),
	}

	output := are.ResourceDetailsToSentence()

	expectedOutput := "terraform name of drifted resource group and type azurerm resource group within module example module " +
		"resource at location eastus2 resource name of drifted resource group " +
		"with tag key of Hello and value of Drift with primary category of management."

	require.Equal(t, expectedOutput, output)
}

func TestResourceDetailsToSentence_DualCategory_Azure(t *testing.T) {
	are := azureResourceExtractor{
		currentResourceDetails: &azureResourceDetails{
			terraformModule:       "module.example_module",
			terraformName:         "drifted-resource-group",
			terraformType:         "azurerm_sql_virtual_network_rule",
			azureInstanceLocation: "eastus2",
			azureInstanceName:     "drifted-resource-group",
			azureInstanceTags: map[string]string{
				"Hello": "Drift",
			},
		},
		typeToCategory: azureResourceCategories(),
	}

	output := are.ResourceDetailsToSentence()

	expectedOutput := "terraform name of drifted resource group and type azurerm sql virtual network rule within module example module " +
		"resource at location eastus2 resource name of drifted resource group " +
		"with tag key of Hello and value of Drift with primary category of database and secondary category of networking."

	require.Equal(t, expectedOutput, output)
}

func TestResourceDetailsToSentence_NoneCase_Azure(t *testing.T) {
	are := azureResourceExtractor{
		currentResourceDetails: &azureResourceDetails{
			terraformModule:       "module.example_module",
			terraformName:         "drifted-resource-group",
			terraformType:         "azurerm_resource_group",
			azureInstanceLocation: "global",
			azureInstanceName:     "none",
			azureInstanceTags:     map[string]string{},
		},
		typeToCategory: azureResourceCategories(),
	}

	output := are.ResourceDetailsToSentence()

	expectedOutput := "terraform name of drifted resource group and type azurerm resource group within module example module " +
		"resource at location global with primary category of management."

	require.Equal(t, expectedOutput, output)
}
