package hclcreate

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/Jeffail/gabs/v2"
)

func defineInputGabsContainer() (*gabs.Container, error) {
	byteArray := []byte(`{
  "google-dragondrop-dev": [
    {
      "cost_component": "SQL instance (db-f1-micro, zonal)",
      "is_usage_based": false,
      "monthly_cost": "7.665",
      "monthly_quantity": "730",
      "price": "12",
      "resource_name": "google_sql_database_instance.tfer--outside-of-terraform-control-db",
      "sub_resource_name": "",
      "unit": "hours"
    },
    {
      "cost_component": "Storage (SSD, zonal)",
      "is_usage_based": false,
      "monthly_cost": "1.7",
      "monthly_quantity": "10",
      "price": "12",
      "resource_name": "google_sql_database_instance.tfer--outside-of-terraform-control-db",
      "sub_resource_name": "",
      "unit": "GB"
    },
    {
      "cost_component": "IP address (if unused)",
      "is_usage_based": false,
      "monthly_cost": "7.3",
      "monthly_quantity": "730",
      "price": "12",
      "resource_name": "google_sql_database_instance.tfer--outside-of-terraform-control-db",
      "sub_resource_name": "",
      "unit": "hours"
    },
    {
      "cost_component": "Storage (standard)",
      "is_usage_based": true,
      "monthly_cost": "",
      "monthly_quantity": "",
      "price": "12",
      "resource_name": "google_storage_bucket.tfer--dragondrop-modules",
      "sub_resource_name": "",
      "unit": "GiB"
    },
    {
      "cost_component": "Object gets, retrieve bucket/object metadata (class B)",
      "is_usage_based": true,
      "monthly_cost": "",
      "monthly_quantity": "",
      "price": "12",
      "resource_name": "google_storage_bucket.tfer--dragondrop-modules",
      "sub_resource_name": "",
      "unit": "10k operations"
    },
    {
      "cost_component": "Data transfer to Australia (first 1TB)",
      "is_usage_based": true,
      "monthly_cost": "",
      "monthly_quantity": "",
      "price": "12",
      "resource_name": "google_storage_bucket.tfer--dragondrop-modules",
      "sub_resource_name": "Network egress",
      "unit": "GB"
    },
    {
      "cost_component": "Data transfer to elsewhere",
      "is_usage_based": true,
      "monthly_cost": "",
      "monthly_quantity": "",
      "price": "12",
      "resource_name": "google_storage_bucket.tfer--dragondrop-modules",
      "sub_resource_name": "Network egress",
      "unit": "GB"
    }
  ]
}`)

	output, err := gabs.ParseJSON(byteArray)
	if err != nil {
		return nil, fmt.Errorf("[unexpected error in gabs.ParseJSON]%v", err)
	}

	return output, nil
}

func Test_gabsContainerToAllCostsStruct(t *testing.T) {
	inputContainer, err := defineInputGabsContainer()
	if err != nil {
		t.Errorf("[defineInputGabsContainer]%v", err)
	}

	expectedOutput := allCosts{
		"google-dragondrop-dev": divisionCosts{
			"google_sql_database_instance.outside_of_terraform_control_db": resourceCosts{
				subResources: map[string][]costComponent{},
				costComponents: []costComponent{
					{
						componentName: "SQL instance (db-f1-micro, zonal)",
						isUsageBased:  false,
						monthlyCost:   "7.665",
						price:         "12",
						unit:          "hours",
					},
					{
						componentName: "Storage (SSD, zonal)",
						isUsageBased:  false,
						monthlyCost:   "1.7",
						price:         "12",
						unit:          "GB",
					},
					{
						componentName: "IP address (if unused)",
						isUsageBased:  false,
						monthlyCost:   "7.3",
						price:         "12",
						unit:          "hours",
					},
				},
			},
			"google_storage_bucket.dragondrop_modules": resourceCosts{
				subResources: map[string][]costComponent{
					"Network egress": {
						{
							componentName: "Data transfer to Australia (first 1TB)",
							isUsageBased:  true,
							monthlyCost:   "",
							price:         "12",
							unit:          "GB",
						},
						{
							componentName: "Data transfer to elsewhere",
							isUsageBased:  true,
							monthlyCost:   "",
							price:         "12",
							unit:          "GB",
						},
					},
				},
				costComponents: []costComponent{
					{
						componentName: "Storage (standard)",
						isUsageBased:  true,
						monthlyCost:   "",
						price:         "12",
						unit:          "GiB",
					},
					{
						componentName: "Object gets, retrieve bucket/object metadata (class B)",
						isUsageBased:  true,
						monthlyCost:   "",
						price:         "12",
						unit:          "10k operations",
					},
				},
			},
		},
	}

	output, err := gabsContainerToAllCostsStruct(inputContainer)
	if err != nil {
		t.Errorf("unexpected error in gabsContainerToAllCostsStruct: %v", err)
	}

	if !reflect.DeepEqual(expectedOutput, output) {
		t.Errorf("got:\n%v\nexpected:\n%v\n", output, expectedOutput)
	}
}
