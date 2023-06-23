package terraformImportMigrationGenerator

import (
	"testing"

	"github.com/Jeffail/gabs/v2"
	"github.com/stretchr/testify/assert"

	terraformValueObjects "github.com/dragondrop-cloud/cloud-concierge/main/internal/implementations/terraform_value_objects"
)

func TestGetResourceLocationFormatted_GCP_StorageBucket(t *testing.T) {
	// Given
	provider := terraformValueObjects.Provider("google")
	resourceType := ResourceType("google_storage_bucket")
	resourcesJSON := []byte(`{
		"name": "tfer--dragondrop-example-2",
		"instances": [
			{
				"attributes_flat": {
					"id": "dragondrop-example-2",
					"arn": "arn:aws:s3:::dragondrop-example-2",
					"project": "example-project",
					"name": "dragondrop-example-2"
				}
			}
		]
	}`)
	resourcesParsed, err := gabs.ParseJSON(resourcesJSON)
	assert.Nil(t, err)

	// When
	resourceFormatted, err := GetRemoteCloudReference(resourcesParsed, provider, resourceType)

	// Then
	assert.Nil(t, err)
	assert.Equal(t, "example-project/dragondrop-example-2", resourceFormatted)
}

func TestGetResourceLocationFormatted_GCP_Resources(t *testing.T) {
	type args struct {
		provider      string
		resourceType  string
		resourcesJSON []byte
	}
	tests := []struct {
		name                     string
		args                     args
		expectedResourceLocation string
	}{
		{
			name: "api-gateway-api",
			args: args{
				provider:     "google",
				resourceType: "google_api_gateway_api",
				resourcesJSON: []byte(`{
					"name": "tfer--dragondrop-api-gateway-api",
					"type": "google_api_gateway_api",
					"instances": [
						{
							"attributes_flat": {
								"id": "dragondrop-api-gateway-api",
								"api_id": "example-api-id",
								"project": "example-project",
								"name": "dragondrop-api-gateway-api"
							}
						}
					]
				}`),
			},
			expectedResourceLocation: "projects/example-project/locations/global/apis/example-api-id",
		},
		{
			name: "google_storage_bucket",
			args: args{
				provider:     "google",
				resourceType: "google_storage_bucket",
				resourcesJSON: []byte(`{
					"name": "tfer--dragondrop-google_storage_bucket",
					"type": "google_storage_bucket",
					"instances": [
						{
							"attributes_flat": {
								"id": "dragondrop-google_storage_bucket",
								"api_id": "example-api-id",
								"project": "example-project",
								"name": "dragondrop-google_storage_bucket",
								"region": "example-project"
							}
						}
					]
				}`),
			},
			expectedResourceLocation: "example-project/dragondrop-google_storage_bucket",
		},
		{
			name: "google_compute_address",
			args: args{
				provider:     "google",
				resourceType: "google_compute_address",
				resourcesJSON: []byte(`{
					"name": "tfer--dragondrop-google_compute_address",
					"type": "google_compute_address",
					"instances": [
						{
							"attributes_flat": {
								"id": "dragondrop-google_compute_address",
								"api_id": "example-api-id",
								"project": "example-project",
								"name": "dragondrop-google_compute_address",
								"region": "example-region"
							}
						}
					]
				}`),
			},
			expectedResourceLocation: "projects/example-project/regions/example-region/addresses/dragondrop-google_compute_address",
		},
		{
			name: "google_compute_disk",
			args: args{
				provider:     "google",
				resourceType: "google_compute_disk",
				resourcesJSON: []byte(`{
					"name": "tfer--dragondrop-google_compute_disk",
					"type": "google_compute_disk",
					"instances": [
						{
							"attributes_flat": {
								"id": "dragondrop-google_compute_disk",
								"api_id": "example-api-id",
								"project": "example-project",
								"name": "dragondrop-google_compute_disk",
								"zone": "example-zone"
							}
						}
					]
				}`),
			},
			expectedResourceLocation: "projects/example-project/zones/example-zone/disks/dragondrop-google_compute_disk",
		},
		{
			name: "google_compute_firewall",
			args: args{
				provider:     "google",
				resourceType: "google_compute_firewall",
				resourcesJSON: []byte(`{
					"name": "tfer--dragondrop-google_compute_firewall",
					"type": "google_compute_firewall",
					"instances": [
						{
							"attributes_flat": {
								"id": "dragondrop-google_compute_firewall",
								"api_id": "example-api-id",
								"project": "example-project",
								"name": "dragondrop-google_compute_firewall"
							}
						}
					]
				}`),
			},
			expectedResourceLocation: "projects/example-project/global/firewalls/dragondrop-google_compute_firewall",
		},
		{
			name: "google_compute_image",
			args: args{
				provider:     "google",
				resourceType: "google_compute_image",
				resourcesJSON: []byte(`{
					"name": "tfer--dragondrop-google_compute_image",
					"type": "google_compute_image",
					"instances": [
						{
							"attributes_flat": {
								"id": "dragondrop-google_compute_image",
								"api_id": "example-api-id",
								"project": "example-project",
								"name": "dragondrop-google_compute_image"
							}
						}
					]
				}`),
			},
			expectedResourceLocation: "projects/example-project/global/images/dragondrop-google_compute_image",
		},
		{
			name: "google_compute_instance",
			args: args{
				provider:     "google",
				resourceType: "google_compute_instance",
				resourcesJSON: []byte(`{
					"name": "tfer--dragondrop-google_compute_instance",
					"type": "google_compute_instance",
					"instances": [
						{
							"attributes_flat": {
								"id": "dragondrop-google_compute_instance",
								"api_id": "example-api-id",
								"project": "example-project",
								"name": "dragondrop-google_compute_instance",
								"zone": "example-zone"
							}
						}
					]
				}`),
			},
			expectedResourceLocation: "projects/example-project/zones/example-zone/instances/dragondrop-google_compute_instance",
		},
		{
			name: "google_compute_network",
			args: args{
				provider:     "google",
				resourceType: "google_compute_network",
				resourcesJSON: []byte(`{
					"name": "tfer--dragondrop-google_compute_network",
					"type": "google_compute_network",
					"instances": [
						{
							"attributes_flat": {
								"id": "dragondrop-google_compute_network",
								"api_id": "example-api-id",
								"project": "example-project",
								"name": "dragondrop-google_compute_network"
							}
						}
					]
				}`),
			},
			expectedResourceLocation: "projects/example-project/global/networks/dragondrop-google_compute_network",
		},
		{
			name: "google_compute_router",
			args: args{
				provider:     "google",
				resourceType: "google_compute_router",
				resourcesJSON: []byte(`{
					"name": "tfer--dragondrop-google_compute_router",
					"type": "google_compute_router",
					"instances": [
						{
							"attributes_flat": {
								"id": "dragondrop-google_compute_router",
								"api_id": "example-api-id",
								"project": "example-project",
								"name": "dragondrop-google_compute_router",
								"region": "example-region"
							}
						}
					]
				}`),
			},
			expectedResourceLocation: "projects/example-project/regions/example-region/routers/dragondrop-google_compute_router",
		},
		{
			name: "google_compute_subnetwork",
			args: args{
				provider:     "google",
				resourceType: "google_compute_subnetwork",
				resourcesJSON: []byte(`{
					"name": "tfer--dragondrop-google_compute_subnetwork",
					"type": "google_compute_subnetwork",
					"instances": [
						{
							"attributes_flat": {
								"id": "dragondrop-google_compute_subnetwork",
								"api_id": "example-api-id",
								"project": "example-project",
								"name": "dragondrop-google_compute_subnetwork",
								"region": "example-region"
							}
						}
					]
				}`),
			},
			expectedResourceLocation: "projects/example-project/regions/example-region/subnetworks/dragondrop-google_compute_subnetwork",
		},
		{
			name: "google_container_cluster",
			args: args{
				provider:     "google",
				resourceType: "google_container_cluster",
				resourcesJSON: []byte(`{
					"name": "tfer--dragondrop-google_container_cluster",
					"type": "google_container_cluster",
					"instances": [
						{
							"attributes_flat": {
								"id": "dragondrop-google_container_cluster",
								"api_id": "example-api-id",
								"project": "example-project",
								"name": "dragondrop-google_container_cluster",
								"location": "example-location"
							}
						}
					]
				}`),
			},
			expectedResourceLocation: "projects/example-project/locations/example-location/clusters/dragondrop-google_container_cluster",
		},
		{
			name: "google_dns_managed_zone",
			args: args{
				provider:     "google",
				resourceType: "google_dns_managed_zone",
				resourcesJSON: []byte(`{
					"name": "tfer--dragondrop-google_dns_managed_zone",
					"type": "google_dns_managed_zone",
					"instances": [
						{
							"attributes_flat": {
								"id": "dragondrop-google_dns_managed_zone",
								"api_id": "example-api-id",
								"project": "example-project",
								"name": "dragondrop-google_dns_managed_zone"
							}
						}
					]
				}`),
			},
			expectedResourceLocation: "projects/example-project/managedZones/dragondrop-google_dns_managed_zone",
		},
		{
			name: "google_kms_crypto_key",
			args: args{
				provider:     "google",
				resourceType: "google_kms_crypto_key",
				resourcesJSON: []byte(`{
					"name": "tfer--dragondrop-google_kms_crypto_key",
					"type": "google_kms_crypto_key",
					"instances": [
						{
							"attributes_flat": {
								"id": "dragondrop-google_kms_crypto_key",
								"api_id": "example-api-id",
								"project": "example-project",
								"name": "dragondrop-google_kms_crypto_key",
								"key_ring": "example-key-ring",
								"location": "example-location"
							}
						}
					]
				}`),
			},
			expectedResourceLocation: "projects/example-project/locations/example-location/keyRings/example-key-ring/cryptoKeys/dragondrop-google_kms_crypto_key",
		},
		{
			name: "google_kms_key_ring",
			args: args{
				provider:     "google",
				resourceType: "google_kms_key_ring",
				resourcesJSON: []byte(`{
					"name": "tfer--dragondrop-google_kms_key_ring",
					"type": "google_kms_key_ring",
					"instances": [
						{
							"attributes_flat": {
								"id": "dragondrop-google_kms_key_ring",
								"api_id": "example-api-id",
								"project": "example-project",
								"name": "dragondrop-google_kms_key_ring",
								"location": "example-location"
							}
						}
					]
				}`),
			},
			expectedResourceLocation: "projects/example-project/locations/example-location/keyRings/dragondrop-google_kms_key_ring",
		},
		{
			name: "google_pubsub_subscription",
			args: args{
				provider:     "google",
				resourceType: "google_pubsub_subscription",
				resourcesJSON: []byte(`{
					"name": "tfer--dragondrop-google_pubsub_subscription",
					"type": "google_pubsub_subscription",
					"instances": [
						{
							"attributes_flat": {
								"id": "dragondrop-google_pubsub_subscription",
								"api_id": "example-api-id",
								"project": "example-project",
								"name": "dragondrop-google_pubsub_subscription"
							}
						}
					]
				}`),
			},
			expectedResourceLocation: "projects/example-project/subscriptions/dragondrop-google_pubsub_subscription",
		},
		{
			name: "google_pubsub_topic",
			args: args{
				provider:     "google",
				resourceType: "google_pubsub_topic",
				resourcesJSON: []byte(`{
					"name": "tfer--dragondrop-google_pubsub_topic",
					"type": "google_pubsub_topic",
					"instances": [
						{
							"attributes_flat": {
								"id": "dragondrop-google_pubsub_topic",
								"api_id": "example-api-id",
								"project": "example-project",
								"name": "dragondrop-google_pubsub_topic"
							}
						}
					]
				}`),
			},
			expectedResourceLocation: "projects/example-project/topics/dragondrop-google_pubsub_topic",
		},
		{
			name: "google_compute_autoscaler",
			args: args{
				provider:     "google",
				resourceType: "google_compute_autoscaler",
				resourcesJSON: []byte(`{
				"name": "tfer--dragondrop-google_compute_autoscaler",
				"type": "google_compute_autoscaler",
				"instances": [
					{
						"attributes_flat": {
							"id": "dragondrop-google_compute_autoscaler",
							"project": "example-project",
							"name": "dragondrop-google_compute_autoscaler",
							"zone": "example-zone"
						}
					}
				]
			}`),
			},
			expectedResourceLocation: "projects/example-project/zones/example-zone/autoscalers/dragondrop-google_compute_autoscaler",
		},
		{
			name: "google_compute_backend_bucket",
			args: args{
				provider:     "google",
				resourceType: "google_compute_backend_bucket",
				resourcesJSON: []byte(`{
				"name": "tfer--dragondrop-google_compute_backend_bucket",
				"type": "google_compute_backend_bucket",
				"instances": [
					{
						"attributes_flat": {
							"id": "dragondrop-google_compute_backend_bucket",
							"project": "example-project",
							"name": "dragondrop-google_compute_backend_bucket"
						}
					}
				]
			}`),
			},
			expectedResourceLocation: "projects/example-project/global/backendBuckets/dragondrop-google_compute_backend_bucket",
		},
		{
			name: "google_compute_backend_service",
			args: args{
				provider:     "google",
				resourceType: "google_compute_backend_service",
				resourcesJSON: []byte(`{
				"name": "tfer--dragondrop-google_compute_backend_service",
				"type": "google_compute_backend_service",
				"instances": [
					{
						"attributes_flat": {
							"id": "dragondrop-google_compute_backend_service",
							"project": "example-project",
							"name": "dragondrop-google_compute_backend_service"
						}
					}
				]
			}`),
			},
			expectedResourceLocation: "projects/example-project/global/backendServices/dragondrop-google_compute_backend_service",
		},
		{
			name: "google_bigquery_dataset",
			args: args{
				provider:     "google",
				resourceType: "google_bigquery_dataset",
				resourcesJSON: []byte(`{
				"name": "tfer--dragondrop-google_bigquery_dataset",
				"type": "google_bigquery_dataset",
				"instances": [
					{
						"attributes_flat": {
							"id": "dragondrop-google_bigquery_dataset",
							"project": "example-project",
							"dataset_id": "dragondrop-google_bigquery_dataset"
						}
					}
				]
			}`),
			},
			expectedResourceLocation: "projects/example-project/datasets/dragondrop-google_bigquery_dataset",
		},
		{
			name: "google_bigquery_table",
			args: args{
				provider:     "google",
				resourceType: "google_bigquery_table",
				resourcesJSON: []byte(`{
				"name": "tfer--dragondrop-google_bigquery_table",
				"type": "google_bigquery_table",
				"instances": [
					{
						"attributes_flat": {
							"id": "dragondrop-google_bigquery_table",
							"project": "example-project",
							"dataset_id": "example-dataset",
							"table_id": "example-table"
						}
					}
				]
			}`),
			},
			expectedResourceLocation: "projects/example-project/datasets/example-dataset/tables/example-table",
		},
		{
			name: "google_cloudbuild_trigger",
			args: args{
				provider:     "google",
				resourceType: "google_cloudbuild_trigger",
				resourcesJSON: []byte(`{
				"name": "tfer--dragondrop-google_cloudbuild_trigger",
				"type": "google_cloudbuild_trigger",
				"instances": [
					{
						"attributes_flat": {
							"id": "dragondrop-google_cloudbuild_trigger",
							"project": "example-project",
							"trigger_id": "example-trigger-id"
						}
					}
				]
			}`),
			},
			expectedResourceLocation: "projects/example-project/triggers/example-trigger-id",
		},
		{
			name: "google_cloudfunctions_function",
			args: args{
				provider:     "google",
				resourceType: "google_cloudfunctions_function",
				resourcesJSON: []byte(`{
				"name": "tfer--dragondrop-google_cloudfunctions_function",
				"type": "google_cloudfunctions_function",
				"instances": [
					{
						"attributes_flat": {
							"id": "dragondrop-google_cloudfunctions_function",
							"project": "example-project",
							"location": "example-location",
							"name": "example-function-name"
						}
					}
				]
			}`),
			},
			expectedResourceLocation: "projects/example-project/locations/example-location/functions/example-function-name",
		},
		{
			name: "google_sql_database_instance",
			args: args{
				provider:     "google",
				resourceType: "google_sql_database_instance",
				resourcesJSON: []byte(`{
				"name": "tfer--dragondrop-google_sql_database_instance",
				"type": "google_sql_database_instance",
				"instances": [
					{
						"attributes_flat": {
							"id": "dragondrop-google_sql_database_instance",
							"project": "example-project",
							"name": "example-instance-name"
						}
					}
				]
			}`),
			},
			expectedResourceLocation: "projects/example-project/instances/example-instance-name",
		},
		{
			name: "google_sql_database",
			args: args{
				provider:     "google",
				resourceType: "google_sql_database",
				resourcesJSON: []byte(`{
				"name": "tfer--dragondrop-google_sql_database",
				"type": "google_sql_database",
				"instances": [
					{
						"attributes_flat": {
							"id": "dragondrop-google_sql_database",
							"api_id": "example-api-id",
							"project": "example-project",
							"instance": "example-instance",
							"name": "dragondrop-google_sql_database"
						}
					}
				]
			}`),
			},
			expectedResourceLocation: "projects/example-project/instances/example-instance/databases/dragondrop-google_sql_database",
		},
		{
			name: "google_dataproc_cluster",
			args: args{
				provider:     "google",
				resourceType: "google_dataproc_cluster",
				resourcesJSON: []byte(`{
				"name": "tfer--dragondrop-google_dataproc_cluster",
				"type": "google_dataproc_cluster",
				"instances": [
					{
						"attributes_flat": {
							"id": "dragondrop-google_dataproc_cluster",
							"api_id": "example-api-id",
							"project": "example-project",
							"name": "dragondrop-google_dataproc_cluster",
							"region": "example-region"
						}
					}
				]
			}`),
			},
			expectedResourceLocation: "projects/example-project/regions/example-region/clusters/dragondrop-google_dataproc_cluster",
		},
		{
			name: "google_compute_external_vpn_gateway",
			args: args{
				provider:     "google",
				resourceType: "google_compute_external_vpn_gateway",
				resourcesJSON: []byte(`{
				"name": "tfer--dragondrop-google_compute_external_vpn_gateway",
				"type": "google_compute_external_vpn_gateway",
				"instances": [
					{
						"attributes_flat": {
							"id": "dragondrop-google_compute_external_vpn_gateway",
							"api_id": "example-api-id",
							"project": "example-project",
							"name": "dragondrop-google_compute_external_vpn_gateway"
						}
					}
				]
			}`),
			},
			expectedResourceLocation: "projects/example-project/global/externalVpnGateways/dragondrop-google_compute_external_vpn_gateway",
		},
		{
			name: "google_dns_record_set",
			args: args{
				provider:     "google",
				resourceType: "google_dns_record_set",
				resourcesJSON: []byte(`{
				"name": "tfer--dragondrop-google_dns_record_set",
				"type": "google_dns_record_set",
				"instances": [
					{
						"attributes_flat": {
							"id": "dragondrop-google_dns_record_set",
							"api_id": "example-api-id",
							"project": "example-project",
							"zone": "example-zone",
							"name": "dragondrop-google_dns_record_set",
							"type": "example-type"
						}
					}
				]
			}`),
			},
			expectedResourceLocation: "projects/example-project/managedZones/example-zone/rrsets/dragondrop-google_dns_record_set/example-type",
		},
		{
			name: "google_compute_forwarding_rule",
			args: args{
				provider:     "google",
				resourceType: "google_compute_forwarding_rule",
				resourcesJSON: []byte(`{
				"name": "tfer--dragondrop-google_compute_forwarding_rule",
				"type": "google_compute_forwarding_rule",
				"instances": [
					{
						"attributes_flat": {
							"id": "dragondrop-google_compute_forwarding_rule",
							"api_id": "example-api-id",
							"project": "example-project",
							"name": "dragondrop-google_compute_forwarding_rule",
							"region": "example-region"
						}
					}
				]
			}`),
			},
			expectedResourceLocation: "projects/example-project/regions/example-region/forwardingRules/dragondrop-google_compute_forwarding_rule",
		},
		{
			name: "google_storage_bucket_iam_binding",
			args: args{
				provider:     "google",
				resourceType: "google_storage_bucket_iam_binding",
				resourcesJSON: []byte(`{
				"name": "tfer--dragondrop-google_storage_bucket_iam_binding",
				"type": "google_storage_bucket_iam_binding",
				"instances": [
					{
						"attributes_flat": {
							"id": "dragondrop-google_storage_bucket_iam_binding",
							"api_id": "example-api-id",
							"bucket": "example-bucket"
						}
					}
				]
			}`),
			},
			expectedResourceLocation: "b/example-bucket/iam",
		},
		{
			name: "google_storage_bucket_iam_member",
			args: args{
				provider:     "google",
				resourceType: "google_storage_bucket_iam_member",
				resourcesJSON: []byte(`{
				"name": "tfer--dragondrop-google_storage_bucket_iam_member",
				"type": "google_storage_bucket_iam_member",
				"instances": [
					{
						"attributes_flat": {
							"id": "dragondrop-google_storage_bucket_iam_member",
							"api_id": "example-api-id",
							"bucket": "example-bucket"
						}
					}
				]
			}`),
			},
			expectedResourceLocation: "b/example-bucket/iam",
		},
		{
			name: "google_compute_global_address",
			args: args{
				provider:     "google",
				resourceType: "google_compute_global_address",
				resourcesJSON: []byte(`{
				"name": "tfer--dragondrop-google_compute_global_address",
				"type": "google_compute_global_address",
				"instances": [
					{
						"attributes_flat": {
							"id": "dragondrop-google_compute_global_address",
							"api_id": "example-api-id",
							"project": "example-project",
							"name": "dragondrop-google_compute_global_address"
						}
					}
				]
			}`),
			},
			expectedResourceLocation: "projects/example-project/global/addresses/dragondrop-google_compute_global_address",
		},
		{
			name: "google_compute_global_forwarding_rule",
			args: args{
				provider:     "google",
				resourceType: "google_compute_global_forwarding_rule",
				resourcesJSON: []byte(`{
				"name": "tfer--dragondrop-google_compute_global_forwarding_rule",
				"type": "google_compute_global_forwarding_rule",
				"instances": [
					{
						"attributes_flat": {
							"id": "dragondrop-google_compute_global_forwarding_rule",
							"api_id": "example-api-id",
							"project": "example-project",
							"name": "dragondrop-google_compute_global_forwarding_rule"
						}
					}
				]
			}`),
			},
			expectedResourceLocation: "projects/example-project/global/forwardingRules/dragondrop-google_compute_global_forwarding_rule",
		},
		{
			name: "google_container_node_pool",
			args: args{
				provider:     "google",
				resourceType: "google_container_node_pool",
				resourcesJSON: []byte(`{
				"name": "tfer--dragondrop-google_container_node_pool",
				"type": "google_container_node_pool",
				"instances": [
					{
						"attributes_flat": {
							"id": "dragondrop-google_container_node_pool",
							"api_id": "example-api-id",
							"project": "example-project",
							"name": "dragondrop-google_container_node_pool",
							"location": "example-location",
							"cluster": "example-cluster"
						}
					}
				]
			}`),
			},
			expectedResourceLocation: "projects/example-project/locations/example-location/clusters/example-cluster/nodePools/dragondrop-google_container_node_pool",
		},
		{
			name: "google_storage_notification",
			args: args{
				provider:     "google",
				resourceType: "google_storage_notification",
				resourcesJSON: []byte(`{
				"name": "tfer--dragondrop-google_storage_notification",
				"type": "google_storage_notification",
				"instances": [
					{
						"attributes_flat": {
							"id": "dragondrop-google_storage_notification",
							"api_id": "example-api-id",
							"project": "example-project",
							"bucket": "example-bucket",
							"name": "dragondrop-google_storage_notification"
						}
					}
				]
			}`),
			},
			expectedResourceLocation: "projects/example-project/buckets/example-bucket/notificationConfigs/dragondrop-google_storage_notification",
		},
		{
			name: "google_compute_global_forwarding_rule",
			args: args{
				provider:     "google",
				resourceType: "google_compute_global_forwarding_rule",
				resourcesJSON: []byte(`{
				"name": "tfer--dragondrop-google_compute_global_forwarding_rule",
				"type": "google_compute_global_forwarding_rule",
				"instances": [
					{
						"attributes_flat": {
							"id": "dragondrop-google_compute_global_forwarding_rule",
							"api_id": "example-api-id",
							"project": "example-project",
							"name": "dragondrop-google_compute_global_forwarding_rule"
						}
					}
				]
			}`),
			},
			expectedResourceLocation: "projects/example-project/global/forwardingRules/dragondrop-google_compute_global_forwarding_rule",
		},
		{
			name: "google_compute_health_check",
			args: args{
				provider:     "google",
				resourceType: "google_compute_health_check",
				resourcesJSON: []byte(`{
				"name": "tfer--dragondrop-google_compute_health_check",
				"type": "google_compute_health_check",
				"instances": [
					{
						"attributes_flat": {
							"id": "dragondrop-google_compute_health_check",
							"api_id": "example-api-id",
							"project": "example-project",
							"name": "dragondrop-google_compute_health_check"
						}
					}
				]
			}`),
			},
			expectedResourceLocation: "projects/example-project/global/healthChecks/dragondrop-google_compute_health_check",
		},
		{
			name: "google_compute_http_health_check",
			args: args{
				provider:     "google",
				resourceType: "google_compute_http_health_check",
				resourcesJSON: []byte(`{
				"name": "tfer--dragondrop-google_compute_http_health_check",
				"type": "google_compute_http_health_check",
				"instances": [
					{
						"attributes_flat": {
							"id": "dragondrop-google_compute_http_health_check",
							"api_id": "example-api-id",
							"project": "example-project",
							"name": "dragondrop-google_compute_http_health_check"
						}
					}
				]
			}`),
			},
			expectedResourceLocation: "projects/example-project/global/httpHealthChecks/dragondrop-google_compute_http_health_check",
		},
		{
			name: "google_compute_https_health_check",
			args: args{
				provider:     "google",
				resourceType: "google_compute_https_health_check",
				resourcesJSON: []byte(`{
				"name": "tfer--dragondrop-google_compute_https_health_check",
				"type": "google_compute_https_health_check",
				"instances": [
					{
						"attributes_flat": {
							"id": "dragondrop-google_compute_https_health_check",
							"api_id": "example-api-id",
							"project": "example-project",
							"name": "dragondrop-google_compute_https_health_check"
						}
					}
				]
			}`),
			},
			expectedResourceLocation: "projects/example-project/global/httpsHealthChecks/dragondrop-google_compute_https_health_check",
		},
		{
			name: "google_project_iam_custom_role",
			args: args{
				provider:     "google",
				resourceType: "google_project_iam_custom_role",
				resourcesJSON: []byte(`{
				"name": "tfer--dragondrop-google_project_iam_custom_role",
				"type": "google_project_iam_custom_role",
				"instances": [
					{
						"attributes_flat": {
							"id": "dragondrop-google_project_iam_custom_role",
							"api_id": "example-api-id",
							"project": "example-project",
							"name": "dragondrop-google_project_iam_custom_role"
						}
					}
				]
			}`),
			},
			expectedResourceLocation: "projects/example-project/roles/dragondrop-google_project_iam_custom_role",
		},
		{
			name: "google_project_iam_member",
			args: args{
				provider:     "google",
				resourceType: "google_project_iam_member",
				resourcesJSON: []byte(`{
				"name": "tfer--dragondrop-google_project_iam_member",
				"type": "google_project_iam_member",
				"instances": [
					{
						"attributes_flat": {
							"id": "dragondrop-google_project_iam_member",
							"api_id": "example-api-id",
							"project": "example-project"
						}
					}
				]
			}`),
			},
			expectedResourceLocation: "projects/example-project/iam",
		},
		{
			name: "google_compute_instance_group_manager",
			args: args{
				provider:     "google",
				resourceType: "google_compute_instance_group_manager",
				resourcesJSON: []byte(`{
				"name": "tfer--dragondrop-google_compute_instance_group_manager",
				"type": "google_compute_instance_group_manager",
				"instances": [
					{
						"attributes_flat": {
							"id": "dragondrop-google_compute_instance_group_manager",
							"api_id": "example-api-id",
							"project": "example-project",
							"zone": "example-zone",
							"name": "dragondrop-google_compute_instance_group_manager"
						}
					}
				]
			}`),
			},
			expectedResourceLocation: "projects/example-project/zones/example-zone/instanceGroupManagers/dragondrop-google_compute_instance_group_manager",
		},
		{
			name: "google_compute_instance_group",
			args: args{
				provider:     "google",
				resourceType: "google_compute_instance_group",
				resourcesJSON: []byte(`{
				"name": "tfer--dragondrop-google_compute_instance_group",
				"type": "google_compute_instance_group",
				"instances": [
					{
						"attributes_flat": {
							"id": "dragondrop-google_compute_instance_group",
							"api_id": "example-api-id",
							"project": "example-project",
							"zone": "example-zone",
							"name": "dragondrop-google_compute_instance_group"
						}
					}
				]
			}`),
			},
			expectedResourceLocation: "projects/example-project/zones/example-zone/instanceGroups/dragondrop-google_compute_instance_group",
		},
		{
			name: "google_compute_instance_template",
			args: args{
				provider:     "google",
				resourceType: "google_compute_instance_template",
				resourcesJSON: []byte(`{
                "name": "tfer--dragondrop-google_compute_instance_template",
                "type": "google_compute_instance_template",
                "instances": [
                    {
                        "attributes_flat": {
                            "id": "dragondrop-google_compute_instance_template",
                            "api_id": "example-api-id",
                            "project": "example-project",
                            "name": "dragondrop-google_compute_instance_template"
                        }
                    }
                ]
            }`),
			},
			expectedResourceLocation: "projects/example-project/global/instanceTemplates/dragondrop-google_compute_instance_template",
		},
		{
			name: "google_compute_interconnect_attachment",
			args: args{
				provider:     "google",
				resourceType: "google_compute_interconnect_attachment",
				resourcesJSON: []byte(`{
                "name": "tfer--dragondrop-google_compute_interconnect_attachment",
                "type": "google_compute_interconnect_attachment",
                "instances": [
                    {
                        "attributes_flat": {
                            "id": "dragondrop-google_compute_interconnect_attachment",
                            "api_id": "example-api-id",
                            "project": "example-project",
                            "region": "example-region",
                            "name": "dragondrop-google_compute_interconnect_attachment"
                        }
                    }
                ]
            }`),
			},
			expectedResourceLocation: "projects/example-project/regions/example-region/interconnectAttachments/dragondrop-google_compute_interconnect_attachment",
		},
		{
			name: "google_logging_metric",
			args: args{
				provider:     "google",
				resourceType: "google_logging_metric",
				resourcesJSON: []byte(`{
                "name": "tfer--dragondrop-google_logging_metric",
                "type": "google_logging_metric",
                "instances": [
                    {
                        "attributes_flat": {
                            "id": "dragondrop-google_logging_metric",
                            "api_id": "example-api-id",
                            "project": "example-project",
                            "name": "dragondrop-google_logging_metric"
                        }
                    }
                ]
            }`),
			},
			expectedResourceLocation: "projects/example-project/metrics/dragondrop-google_logging_metric",
		},
		{
			name: "google_redis_instance",
			args: args{
				provider:     "google",
				resourceType: "google_redis_instance",
				resourcesJSON: []byte(`{
                "name": "tfer--dragondrop-google_redis_instance",
                "type": "google_redis_instance",
                "instances": [
                    {
                        "attributes_flat": {
                            "id": "dragondrop-google_redis_instance",
                            "api_id": "example-api-id",
                            "project": "example-project",
                            "location": "example-location",
                            "name": "dragondrop-google_redis_instance"
                        }
                    }
                ]
            }`),
			},
			expectedResourceLocation: "projects/example-project/locations/example-location/instances/dragondrop-google_redis_instance",
		},
		{
			name: "google_monitoring_alert_policy",
			args: args{
				provider:     "google",
				resourceType: "google_monitoring_alert_policy",
				resourcesJSON: []byte(`{
				"name": "tfer--dragondrop-google_monitoring_alert_policy",
				"type": "google_monitoring_alert_policy",
				"instances": [
					{
						"attributes_flat": {
							"id": "dragondrop-google_monitoring_alert_policy",
							"api_id": "example-api-id",
							"project": "example-project",
							"name": "dragondrop-google_monitoring_alert_policy"
						}
					}
				]
			}`),
			},
			expectedResourceLocation: "projects/example-project/alertPolicies/dragondrop-google_monitoring_alert_policy",
		},
		{
			name: "google_monitoring_group",
			args: args{
				provider:     "google",
				resourceType: "google_monitoring_group",
				resourcesJSON: []byte(`{
				"name": "tfer--example-google_monitoring_group",
				"type": "google_monitoring_group",
				"instances": [
					{
						"attributes_flat": {
							"id": "example-google_monitoring_group",
							"project": "example-project",
							"name": "example-google_monitoring_group"
						}
					}
				]
			}`),
			},
			expectedResourceLocation: "projects/example-project/groups/example-google_monitoring_group",
		},
		{
			name: "google_monitoring_notification_channel",
			args: args{
				provider:     "google",
				resourceType: "google_monitoring_notification_channel",
				resourcesJSON: []byte(`{
				"name": "tfer--example-google_monitoring_notification_channel",
				"type": "google_monitoring_notification_channel",
				"instances": [
					{
						"attributes_flat": {
							"id": "example-google_monitoring_notification_channel",
							"project": "example-project",
							"name": "example-google_monitoring_notification_channel"
						}
					}
				]
			}`),
			},
			expectedResourceLocation: "projects/example-project/notificationChannels/example-google_monitoring_notification_channel",
		},
		{
			name: "google_monitoring_uptime_check_config",
			args: args{
				provider:     "google",
				resourceType: "google_monitoring_uptime_check_config",
				resourcesJSON: []byte(`{
				"name": "tfer--example-google_monitoring_uptime_check_config",
				"type": "google_monitoring_uptime_check_config",
				"instances": [
					{
						"attributes_flat": {
							"id": "example-google_monitoring_uptime_check_config",
							"project": "example-project",
							"name": "example-google_monitoring_uptime_check_config"
						}
					}
				]
			}`),
			},
			expectedResourceLocation: "projects/example-project/uptimeCheckConfigs/example-google_monitoring_uptime_check_config",
		},
		{
			name: "google_compute_packet_mirroring",
			args: args{
				provider:     "google",
				resourceType: "google_compute_packet_mirroring",
				resourcesJSON: []byte(`{
				"name": "tfer--example-google_compute_packet_mirroring",
				"type": "google_compute_packet_mirroring",
				"instances": [
					{
						"attributes_flat": {
							"id": "example-google_compute_packet_mirroring",
							"project": "example-project",
							"name": "example-google_compute_packet_mirroring",
							"region": "example-region"
						}
					}
				]
			}`),
			},
			expectedResourceLocation: "projects/example-project/regions/example-region/packetMirrorings/example-google_compute_packet_mirroring",
		},
		{
			name: "google_compute_node_group",
			args: args{
				provider:     "google",
				resourceType: "google_compute_node_group",
				resourcesJSON: []byte(`{
				"name": "tfer--example-google_compute_node_group",
				"type": "google_compute_node_group",
				"instances": [
					{
						"attributes_flat": {
							"id": "example-google_compute_node_group",
							"project": "example-project",
							"name": "example-google_compute_node_group",
							"zone": "example-zone"
						}
					}
				]
			}`),
			},
			expectedResourceLocation: "projects/example-project/zones/example-zone/nodeGroups/example-google_compute_node_group",
		},
		{
			name: "google_compute_node_template",
			args: args{
				provider:     "google",
				resourceType: "google_compute_node_template",
				resourcesJSON: []byte(`{
				"name": "tfer--example-google_compute_node_template",
				"type": "google_compute_node_template",
				"instances": [
					{
						"attributes_flat": {
							"id": "example-google_compute_node_template",
							"project": "example-project",
							"name": "example-google_compute_node_template",
							"region": "example-region"
						}
					}
				]
			}`),
			},
			expectedResourceLocation: "projects/example-project/regions/example-region/nodeTemplates/example-google_compute_node_template",
		},
		{
			name: "google_project",
			args: args{
				provider:     "google",
				resourceType: "google_project",
				resourcesJSON: []byte(`{
				"name": "tfer--example-google_project",
				"type": "google_project",
				"instances": [
					{
						"attributes_flat": {
							"id": "example-google_project",
							"project_id": "example-project"
						}
					}
				]
			}`),
			},
			expectedResourceLocation: "projects/example-project",
		},
		{
			name: "google_compute_region_autoscaler",
			args: args{
				provider:     "google",
				resourceType: "google_compute_region_autoscaler",
				resourcesJSON: []byte(`{
				"name": "tfer--example-google_compute_region_autoscaler",
				"type": "google_compute_region_autoscaler",
				"instances": [
					{
						"attributes_flat": {
							"id": "example-google_compute_region_autoscaler",
							"project": "example-project",
							"name": "example-google_compute_region_autoscaler",
							"region": "example-region"
						}
					}
				]
			}`),
			},
			expectedResourceLocation: "projects/example-project/regions/example-region/autoscalers/example-google_compute_region_autoscaler",
		},
		{
			name: "google_compute_region_backend_service",
			args: args{
				provider:     "google",
				resourceType: "google_compute_region_backend_service",
				resourcesJSON: []byte(`{
				"name": "tfer--example-google_compute_region_backend_service",
				"type": "google_compute_region_backend_service",
				"instances": [
					{
						"attributes_flat": {
							"id": "example-google_compute_region_backend_service",
							"project": "example-project",
							"name": "example-google_compute_region_backend_service",
							"region": "example-region"
						}
					}
				]
			}`),
			},
			expectedResourceLocation: "projects/example-project/regions/example-region/backendServices/example-google_compute_region_backend_service",
		},
		{
			name: "google_compute_region_disk",
			args: args{
				provider:     "google",
				resourceType: "google_compute_region_disk",
				resourcesJSON: []byte(`{
				"name": "tfer--example-google_compute_region_disk",
				"type": "google_compute_region_disk",
				"instances": [
					{
						"attributes_flat": {
							"id": "example-google_compute_region_disk",
							"project": "example-project",
							"name": "example-google_compute_region_disk",
							"region": "example-region"
						}
					}
				]
			}`),
			},
			expectedResourceLocation: "projects/example-project/regions/example-region/disks/example-google_compute_region_disk",
		},
		{
			name: "google_compute_region_health_check",
			args: args{
				provider:     "google",
				resourceType: "google_compute_region_health_check",
				resourcesJSON: []byte(`{
				"name": "tfer--example-google_compute_region_health_check",
				"type": "google_compute_region_health_check",
				"instances": [
					{
						"attributes_flat": {
							"id": "example-google_compute_region_health_check",
							"project": "example-project",
							"name": "example-google_compute_region_health_check",
							"region": "example-region"
						}
					}
				]
			}`),
			},
			expectedResourceLocation: "projects/example-project/regions/example-region/healthChecks/example-google_compute_region_health_check",
		},
		{
			name: "google_compute_region_instance_group",
			args: args{
				provider:     "google",
				resourceType: "google_compute_region_instance_group",
				resourcesJSON: []byte(`{
				"name": "tfer--example-google_compute_region_instance_group",
				"type": "google_compute_region_instance_group",
				"instances": [
					{
						"attributes_flat": {
							"id": "example-google_compute_region_instance_group",
							"project": "example-project",
							"name": "example-google_compute_region_instance_group",
							"region": "example-region"
						}
					}
				]
			}`),
			},
			expectedResourceLocation: "projects/example-project/regions/example-region/instanceGroups/example-google_compute_region_instance_group",
		},
		{
			name: "google_compute_region_ssl_certificate",
			args: args{
				provider:     "google",
				resourceType: "google_compute_region_ssl_certificate",
				resourcesJSON: []byte(`{
				"name": "tfer--example-google_compute_region_ssl_certificate",
				"type": "google_compute_region_ssl_certificate",
				"instances": [
					{
						"attributes_flat": {
							"id": "example-google_compute_region_ssl_certificate",
							"project": "example-project",
							"name": "example-google_compute_region_ssl_certificate",
							"region": "example-region"
						}
					}
				]
			}`),
			},
			expectedResourceLocation: "projects/example-project/regions/example-region/sslCertificates/example-google_compute_region_ssl_certificate",
		},
		{
			name: "google_compute_region_target_http_proxy",
			args: args{
				provider:     "google",
				resourceType: "google_compute_region_target_http_proxy",
				resourcesJSON: []byte(`{
				"name": "tfer--example-google_compute_region_target_http_proxy",
				"type": "google_compute_region_target_http_proxy",
				"instances": [
					{
						"attributes_flat": {
							"id": "example-google_compute_region_target_http_proxy",
							"project": "example-project",
							"name": "example-google_compute_region_target_http_proxy",
							"region": "example-region"
						}
					}
				]
			}`),
			},
			expectedResourceLocation: "projects/example-project/regions/example-region/targetHttpProxies/example-google_compute_region_target_http_proxy",
		},
		{
			name: "google_compute_region_target_https_proxy",
			args: args{
				provider:     "google",
				resourceType: "google_compute_region_target_https_proxy",
				resourcesJSON: []byte(`{
				"name": "tfer--example-google_compute_region_target_https_proxy",
				"type": "google_compute_region_target_https_proxy",
				"instances": [
					{
						"attributes_flat": {
							"id": "example-google_compute_region_target_https_proxy",
							"project": "example-project",
							"name": "example-google_compute_region_target_https_proxy",
							"region": "example-region"
						}
					}
				]
			}`),
			},
			expectedResourceLocation: "projects/example-project/regions/example-region/targetHttpsProxies/example-google_compute_region_target_https_proxy",
		},
		{
			name: "google_compute_region_url_map",
			args: args{
				provider:     "google",
				resourceType: "google_compute_region_url_map",
				resourcesJSON: []byte(`{
				"name": "tfer--example-google_compute_region_url_map",
				"type": "google_compute_region_url_map",
				"instances": [
					{
						"attributes_flat": {
							"project": "example-project",
							"region": "example-region",
							"name": "example-url-map"
						}
					}
				]
			}`),
			},
			expectedResourceLocation: "projects/example-project/regions/example-region/urlMaps/example-url-map",
		},
		{
			name: "google_compute_reservation",
			args: args{
				provider:     "google",
				resourceType: "google_compute_reservation",
				resourcesJSON: []byte(`{
				"name": "tfer--example-google_compute_reservation",
				"type": "google_compute_reservation",
				"instances": [
					{
						"attributes_flat": {
							"project": "example-project",
							"zone": "example-zone",
							"name": "example-reservation"
						}
					}
				]
			}`),
			},
			expectedResourceLocation: "projects/example-project/zones/example-zone/reservations/example-reservation",
		},
		{
			name: "google_compute_resource_policy",
			args: args{
				provider:     "google",
				resourceType: "google_compute_resource_policy",
				resourcesJSON: []byte(`{
				"name": "tfer--example-google_compute_resource_policy",
				"type": "google_compute_resource_policy",
				"instances": [
					{
						"attributes_flat": {
							"project": "example-project",
							"region": "example-region",
							"name": "example-resource-policy"
						}
					}
				]
			}`),
			},
			expectedResourceLocation: "projects/example-project/regions/example-region/resourcePolicies/example-resource-policy",
		},
		{
			name: "google_compute_region_instance_group_manager",
			args: args{
				provider:     "google",
				resourceType: "google_compute_region_instance_group_manager",
				resourcesJSON: []byte(`{
				"name": "tfer--example-google_compute_region_instance_group_manager",
				"type": "google_compute_region_instance_group_manager",
				"instances": [
					{
						"attributes_flat": {
							"project": "example-project",
							"region": "example-region",
							"name": "example-instance-group-manager"
						}
					}
				]
			}`),
			},
			expectedResourceLocation: "projects/example-project/regions/example-region/instanceGroupManagers/example-instance-group-manager",
		},
		{
			name: "google_compute_route",
			args: args{
				provider:     "google",
				resourceType: "google_compute_route",
				resourcesJSON: []byte(`{
				"name": "tfer--example-google_compute_route",
				"type": "google_compute_route",
				"instances": [
					{
						"attributes_flat": {
							"project": "example-project",
							"name": "example-route"
						}
					}
				]
			}`),
			},
			expectedResourceLocation: "projects/example-project/global/routes/example-route",
		},
		{
			name: "google_cloud_scheduler_job",
			args: args{
				provider:     "google",
				resourceType: "google_cloud_scheduler_job",
				resourcesJSON: []byte(`{
				"name": "tfer--example-google_cloud_scheduler_job",
				"type": "google_cloud_scheduler_job",
				"instances": [
					{
						"attributes_flat": {
							"project": "example-project",
							"location": "example-location",
							"name": "example-job"
						}
					}
				]
			}`),
			},
			expectedResourceLocation: "projects/example-project/locations/example-location/jobs/example-job",
		},
		{
			name: "google_compute_security_policy",
			args: args{
				provider:     "google",
				resourceType: "google_compute_security_policy",
				resourcesJSON: []byte(`{
				"name": "tfer--example-google_compute_security_policy",
				"type": "google_compute_security_policy",
				"instances": [
					{
						"attributes_flat": {
							"project": "example-project",
							"name": "example-security-policy"
						}
					}
				]
			}`),
			},
			expectedResourceLocation: "projects/example-project/global/securityPolicies/example-security-policy",
		},
		{
			name: "google_compute_managed_ssl_certificate",
			args: args{
				provider:     "google",
				resourceType: "google_compute_managed_ssl_certificate",
				resourcesJSON: []byte(`{
				"name": "tfer--example-google_compute_managed_ssl_certificate",
				"type": "google_compute_managed_ssl_certificate",
				"instances": [
					{
						"attributes_flat": {
							"project": "example-project",
							"name": "example-ssl-certificate"
						}
					}
				]
			}`),
			},
			expectedResourceLocation: "projects/example-project/global/sslCertificates/example-ssl-certificate",
		},
		{
			name: "google_compute_managed_ssl_certificate",
			args: args{
				provider:     "google",
				resourceType: "google_compute_managed_ssl_certificate",
				resourcesJSON: []byte(`{
				"name": "tfer--example-google_compute_managed_ssl_certificate",
				"type": "google_compute_managed_ssl_certificate",
				"instances": [
					{
						"attributes_flat": {
							"project": "example-project",
							"name": "example-ssl-certificate"
						}
					}
				]
			}`),
			},
			expectedResourceLocation: "projects/example-project/global/sslCertificates/example-ssl-certificate",
		},
		{
			name: "google_compute_ssl_policy",
			args: args{
				provider:     "google",
				resourceType: "google_compute_ssl_policy",
				resourcesJSON: []byte(`{
				"name": "tfer--example-google_compute_ssl_policy",
				"type": "google_compute_ssl_policy",
				"instances": [
					{
						"attributes_flat": {
							"project": "example-project",
							"name": "example-ssl-policy"
						}
					}
				]
			}`),
			},
			expectedResourceLocation: "projects/example-project/global/sslPolicies/example-ssl-policy",
		},
		{
			name: "google_compute_target_http_proxy",
			args: args{
				provider:     "google",
				resourceType: "google_compute_target_http_proxy",
				resourcesJSON: []byte(`{
				"name": "tfer--example-google_compute_target_http_proxy",
				"type": "google_compute_target_http_proxy",
				"instances": [
					{
						"attributes_flat": {
							"project": "example-project",
							"name": "example-target-http-proxy"
						}
					}
				]
			}`),
			},
			expectedResourceLocation: "projects/example-project/global/targetHttpProxies/example-target-http-proxy",
		},
		{
			name: "google_compute_target_https_proxy",
			args: args{
				provider:     "google",
				resourceType: "google_compute_target_https_proxy",
				resourcesJSON: []byte(`{
				"name": "tfer--example-google_compute_target_https_proxy",
				"type": "google_compute_target_https_proxy",
				"instances": [
					{
						"attributes_flat": {
							"project": "example-project",
							"name": "example-target-https-proxy"
						}
					}
				]
			}`),
			},
			expectedResourceLocation: "projects/example-project/global/targetHttpsProxies/example-target-https-proxy",
		},
		{
			name: "google_compute_target_instance",
			args: args{
				provider:     "google",
				resourceType: "google_compute_target_instance",
				resourcesJSON: []byte(`{
				"name": "tfer--example-google_compute_target_instance",
				"type": "google_compute_target_instance",
				"instances": [
					{
						"attributes_flat": {
							"project": "example-project",
							"zone": "example-zone",
							"name": "example-target-instance"
						}
					}
				]
			}`),
			},
			expectedResourceLocation: "projects/example-project/zones/example-zone/targetInstances/example-target-instance",
		},
		{
			name: "google_compute_target_pool",
			args: args{
				provider:     "google",
				resourceType: "google_compute_target_pool",
				resourcesJSON: []byte(`{
				"name": "tfer--example-google_compute_target_pool",
				"type": "google_compute_target_pool",
				"instances": [
					{
						"attributes_flat": {
							"project": "example-project",
							"region": "example-region",
							"name": "example-target-pool"
						}
					}
				]
			}`),
			},
			expectedResourceLocation: "projects/example-project/regions/example-region/targetPools/example-target-pool",
		},
		{
			name: "google_compute_target_ssl_proxy",
			args: args{
				provider:     "google",
				resourceType: "google_compute_target_ssl_proxy",
				resourcesJSON: []byte(`{
				"name": "tfer--example-google_compute_target_ssl_proxy",
				"type": "google_compute_target_ssl_proxy",
				"instances": [
					{
						"attributes_flat": {
							"project": "example-project",
							"name": "example-target-ssl-proxy"
						}
					}
				]
			}`),
			},
			expectedResourceLocation: "projects/example-project/global/targetSslProxies/example-target-ssl-proxy",
		},
		{
			name: "google_compute_target_tcp_proxy",
			args: args{
				provider:     "google",
				resourceType: "google_compute_target_tcp_proxy",
				resourcesJSON: []byte(`{
				"name": "tfer--example-google_compute_target_tcp_proxy",
				"type": "google_compute_target_tcp_proxy",
				"instances": [
					{
						"attributes_flat": {
							"project": "example-project",
							"name": "example-target-tcp-proxy"
						}
					}
				]
			}`),
			},
			expectedResourceLocation: "projects/example-project/global/targetTcpProxies/example-target-tcp-proxy",
		},
		{
			name: "google_compute_vpn_gateway",
			args: args{
				provider:     "google",
				resourceType: "google_compute_vpn_gateway",
				resourcesJSON: []byte(`{
				"name": "tfer--example-google_compute_vpn_gateway",
				"type": "google_compute_vpn_gateway",
				"instances": [
					{
						"attributes_flat": {
							"project": "example-project",
							"region": "example-region",
							"name": "example-vpn-gateway"
						}
					}
				]
			}`),
			},
			expectedResourceLocation: "projects/example-project/regions/example-region/vpnGateways/example-vpn-gateway",
		},
		{
			name: "google_compute_url_map",
			args: args{
				provider:     "google",
				resourceType: "google_compute_url_map",
				resourcesJSON: []byte(`{
				"name": "tfer--example-google_compute_url_map",
				"type": "google_compute_url_map",
				"instances": [
					{
						"attributes_flat": {
							"project": "example-project",
							"name": "example-url-map"
						}
					}
				]
			}`),
			},
			expectedResourceLocation: "projects/example-project/global/urlMaps/example-url-map",
		},
		{
			name: "google_compute_vpn_tunnel",
			args: args{
				provider:     "google",
				resourceType: "google_compute_vpn_tunnel",
				resourcesJSON: []byte(`{
				"name": "tfer--example-google_compute_vpn_tunnel",
				"type": "google_compute_vpn_tunnel",
				"instances": [
					{
						"attributes_flat": {
							"project": "example-project",
							"region": "example-region",
							"name": "example-vpn-tunnel"
						}
					}
				]
			}`),
			},
			expectedResourceLocation: "projects/example-project/regions/example-region/vpnTunnels/example-vpn-tunnel",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			resourcesParsed, err := gabs.ParseJSON(tt.args.resourcesJSON)
			assert.Nil(t, err)

			// When
			resourceFormatted, err := GetRemoteCloudReference(resourcesParsed, terraformValueObjects.Provider(tt.args.provider), ResourceType(tt.args.resourceType))

			// Then
			assert.Nil(t, err)
			assert.Equal(t, tt.expectedResourceLocation, resourceFormatted)
		})
	}
}
