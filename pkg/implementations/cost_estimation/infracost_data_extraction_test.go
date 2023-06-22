package costEstimation

import (
	"reflect"
	"strings"
	"testing"
)

func generateInputJSON() []byte {
	return []byte(`{
  "currency": "USD",
  "projects": [
    {
      "name": "./current_cloud/google-dragondrop-dev/",
      "metadata": {
        "path": "./current_cloud/google-dragondrop-dev/",
        "type": "terraform_dir"
      },
      "breakdown": {
        "resources": [
          {
            "name": "google_sql_database_instance.tfer--web-app-instance-dev",
            "metadata": {
              "calls": [
                {
                  "blockName": "google_sql_database_instance.tfer--web-app-instance-dev",
                  "filename": "current_cloud/google-dragondrop-dev/resources.tf"
                }
              ],
              "filename": "current_cloud/google-dragondrop-dev/resources.tf"
            },
            "hourlyCost": "0.022828767123287671",
            "monthlyCost": "16.665",
            "costComponents": [
              {
                "name": "SQL instance (db-f1-micro, zonal)",
                "unit": "hours",
                "hourlyQuantity": "1",
                "monthlyQuantity": "730",
                "price": "0.0105",
                "hourlyCost": "0.0105",
                "monthlyCost": "7.665"
              },
              {
                "name": "Storage (SSD, zonal)",
                "unit": "GB",
                "hourlyQuantity": "0.0136986301369863",
                "monthlyQuantity": "10",
                "price": "0.17",
                "hourlyCost": "0.002328767123287671",
                "monthlyCost": "1.7"
              },
              {
                "name": "Backups",
                "unit": "GB",
                "hourlyQuantity": null,
                "monthlyQuantity": null,
                "price": "0.08",
                "hourlyCost": null,
                "monthlyCost": null
              },
              {
                "name": "IP address (if unused)",
                "unit": "hours",
                "hourlyQuantity": "1",
                "monthlyQuantity": "730",
                "price": "0.01",
                "hourlyCost": "0.01",
                "monthlyCost": "7.3"
              }
            ]
          },
          {
            "name": "google_storage_bucket.tfer--dragondrop-migrations-history-dev",
            "metadata": {
              "calls": [
                {
                  "blockName": "google_storage_bucket.tfer--dragondrop-migrations-history-dev",
                  "filename": "current_cloud/google-dragondrop-dev/resources.tf"
                }
              ],
              "filename": "current_cloud/google-dragondrop-dev/resources.tf"
            },
            "hourlyCost": null,
            "monthlyCost": null,
            "costComponents": [
              {
                "name": "Storage (standard)",
                "unit": "GiB",
                "hourlyQuantity": null,
                "monthlyQuantity": null,
                "price": "0.023",
                "hourlyCost": null,
                "monthlyCost": null
              }
            ],
            "subresources": [
              {
                "name": "Network egress",
                "metadata": {},
                "hourlyCost": null,
                "monthlyCost": null,
                "costComponents": [
                  {
                    "name": "Data transfer in same continent",
                    "unit": "GB",
                    "hourlyQuantity": null,
                    "monthlyQuantity": null,
                    "price": "0.02",
                    "hourlyCost": null,
                    "monthlyCost": null
                  },
                  {
                    "name": "Data transfer to worldwide excluding Asia, Australia (first 1TB)",
                    "unit": "GB",
                    "hourlyQuantity": null,
                    "monthlyQuantity": null,
                    "price": "0.12",
                    "hourlyCost": null,
                    "monthlyCost": null
                  }
                ]
              }
            ]
          }
        ],
        "totalHourlyCost": "0.022828767123287671",
        "totalMonthlyCost": "16.665"
      }
    }
  ],
  "totalHourlyCost": "0.022828767123287671",
  "totalMonthlyCost": "16.665",
  "pastTotalHourlyCost": "0.022828767123287671",
  "pastTotalMonthlyCost": "16.665",
  "diffTotalHourlyCost": "0",
  "diffTotalMonthlyCost": "0",
  "timeGenerated": "2023-04-08T22:44:32.5033959Z",
  "summary": {
    "totalDetectedResources": 18,
    "totalSupportedResources": 4,
    "totalUnsupportedResources": 0,
    "totalUsageBasedResources": 4,
    "totalNoPriceResources": 14,
    "unsupportedResourceCounts": {},
    "noPriceResourceCounts": {
      "google_sql_database": 2,
      "google_storage_bucket_acl": 3,
      "google_storage_bucket_iam_binding": 3,
      "google_storage_bucket_iam_policy": 3,
      "google_storage_default_object_acl": 3
    }
  }
}`)
}

func generateInfracostResourceData() []InfracostResourceData {
	return []InfracostResourceData{
		{
			resourceID:            "google_sql_database_instance.tfer--web-app-instance-dev",
			monthlyCost:           "16.665",
			isPrimarilyUsageBased: false,
			costComponentList: []CostComponent{
				{
					name:            "SQL instance (db-f1-micro, zonal)",
					unit:            "hours",
					unitPrice:       "0.0105",
					monthlyQuantity: "730",
					monthlyCost:     "7.665",
				},
				{
					name:            "Storage (SSD, zonal)",
					unit:            "GB",
					unitPrice:       "0.17",
					monthlyQuantity: "10",
					monthlyCost:     "1.7",
				},
				{
					name:            "Backups",
					unit:            "GB",
					unitPrice:       "0.08",
					monthlyQuantity: "",
					monthlyCost:     "",
				},
				{
					name:            "IP address (if unused)",
					unit:            "hours",
					unitPrice:       "0.01",
					monthlyQuantity: "730",
					monthlyCost:     "7.3",
				},
			},
			subResources: []SubResource{},
		},
		{
			resourceID:            "google_storage_bucket.tfer--dragondrop-migrations-history-dev",
			monthlyCost:           "",
			isPrimarilyUsageBased: true,
			costComponentList: []CostComponent{
				{
					name:            "Storage (standard)",
					unit:            "GiB",
					unitPrice:       "0.023",
					monthlyQuantity: "",
					monthlyCost:     "",
				},
			},
			subResources: []SubResource{
				{
					name:                  "Network egress",
					monthlyCost:           "",
					isPrimarilyUsageBased: true,
					costComponentList: []CostComponent{
						{
							name:            "Data transfer in same continent",
							unit:            "GB",
							unitPrice:       "0.02",
							monthlyQuantity: "",
							monthlyCost:     "",
						},
						{
							name:            "Data transfer to worldwide excluding Asia, Australia (first 1TB)",
							unit:            "GB",
							unitPrice:       "0.12",
							monthlyQuantity: "",
							monthlyCost:     "",
						},
					},
				},
			},
		},
	}
}

func TestCostEstimator_ParseJSONGABSContainerToStruct(t *testing.T) {
	ce := CostEstimator{}

	inputJSONBytes := generateInputJSON()

	expectedOutput := generateInfracostResourceData()

	output, err := ce.ParseJSONToStruct(inputJSONBytes)
	if err != nil {
		t.Errorf("unexpected error in ce.ParseJSONGABSContainerTOStruct: %v", err)
	}

	if !reflect.DeepEqual(expectedOutput, output) {
		t.Errorf("got:\n%v\nexpected:\n%v", output, expectedOutput)
	}
}

func TestCostEstimator_StructToJSONString(t *testing.T) {
	ce := CostEstimator{}

	inputDataStructList := generateInfracostResourceData()

	expectedOutput := `[{"cost_component":"SQL instance (db-f1-micro, zonal)","is_usage_based":false,"monthly_cost":"7.665","monthly_quantity":"730","price":"0.0105","resource_name":"google_sql_database_instance.tfer--web-app-instance-dev","sub_resource_name":"","unit":"hours"},{"cost_component":"Storage (SSD, zonal)","is_usage_based":false,"monthly_cost":"1.7","monthly_quantity":"10","price":"0.17","resource_name":"google_sql_database_instance.tfer--web-app-instance-dev","sub_resource_name":"","unit":"GB"},{"cost_component":"Backups","is_usage_based":false,"monthly_cost":"","monthly_quantity":"","price":"0.08","resource_name":"google_sql_database_instance.tfer--web-app-instance-dev","sub_resource_n
ame":"","unit":"GB"},{"cost_component":"IP address (if unused)","is_usage_based":false,"monthly_cost":"7.3","monthly_quantity":"730","price":"0.01","resource_name":"google_sql_database_instance.tfer--web-app-instance-dev","sub_reso
urce_name":"","unit":"hours"},{"cost_component":"Storage (standard)","is_usage_based":true,"monthly_cost":"","monthly_quantity":"","price":"0.023","resource_name":"google_storage_bucket.tfer--dragondrop-migrations-history-dev","sub
_resource_name":"","unit":"GiB"},{"cost_component":"Data transfer in same continent","is_usage_based":true,"monthly_cost":"","monthly_quantity":"","price":"0.02","resource_name":"google_storage_bucket.tfer--dragondrop-migrations-hi
story-dev","sub_resource_name":"Network egress","unit":"GB"},{"cost_component":"Data transfer to worldwide excluding Asia, Australia (first 1TB)","is_usage_based":true,"monthly_cost":"","monthly_quantity":"","price":"0.12","resource_name":"google_storage_bucket.tfer--dragondrop-migrations-history-dev","sub_resource_name":"Network egress","unit":"GB"}]`

	expectedOutput = strings.Replace(expectedOutput, "\n", "", -1)

	output, err := ce.StructToJSONString(inputDataStructList)
	if err != nil {
		t.Errorf("unexpected error in ce.StructToJSONString: %v", err)
	}

	if !reflect.DeepEqual(expectedOutput, output) {
		t.Errorf("got:\n%v\nexpected:\n%v", output, expectedOutput)
	}
}
