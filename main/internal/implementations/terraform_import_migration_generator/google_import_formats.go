package terraformimportmigrationgenerator

var GoogleResourceTypeLocations = map[ResourceType]ImportLocationFormat{
	"google_storage_bucket": {
		StringFormat: "$0/$1",
		Attributes:   []string{"project", "name"},
	},
	"google_storage_bucket_iam_policy": {
		StringFormat: "b/$0",
		Attributes:   []string{"bucket"},
	},
	"google_compute_address": {
		StringFormat: "projects/$0/regions/$1/addresses/$2",
		Attributes:   []string{"project", "region", "name"},
	},
	"google_compute_disk": {
		StringFormat: "projects/$0/zones/$1/disks/$2",
		Attributes:   []string{"project", "zone", "name"},
	},
	"google_compute_firewall": {
		StringFormat: "projects/$0/global/firewalls/$1",
		Attributes:   []string{"project", "name"},
	},
	"google_compute_image": {
		StringFormat: "projects/$0/global/images/$1",
		Attributes:   []string{"project", "name"},
	},
	"google_compute_instance": {
		StringFormat: "projects/$0/zones/$1/instances/$2",
		Attributes:   []string{"project", "zone", "name"},
	},
	"google_compute_network": {
		StringFormat: "projects/$0/global/networks/$1",
		Attributes:   []string{"project", "name"},
	},
	"google_compute_router": {
		StringFormat: "projects/$0/regions/$1/routers/$2",
		Attributes:   []string{"project", "region", "name"},
	},
	"google_compute_subnetwork": {
		StringFormat: "projects/$0/regions/$1/subnetworks/$2",
		Attributes:   []string{"project", "region", "name"},
	},
	"google_container_cluster": {
		StringFormat: "projects/$0/locations/$1/clusters/$2",
		Attributes:   []string{"project", "location", "name"},
	},
	"google_dns_managed_zone": {
		StringFormat: "projects/$0/managedZones/$1",
		Attributes:   []string{"project", "name"},
	},
	"google_kms_crypto_key": {
		StringFormat: "projects/$0/locations/$1/keyRings/$2/cryptoKeys/$3",
		Attributes:   []string{"project", "location", "key_ring", "name"},
	},
	"google_kms_key_ring": {
		StringFormat: "projects/$0/locations/$1/keyRings/$2",
		Attributes:   []string{"project", "location", "name"},
	},
	"google_pubsub_subscription": {
		StringFormat: "projects/$0/subscriptions/$1",
		Attributes:   []string{"project", "name"},
	},
	"google_pubsub_topic": {
		StringFormat: "projects/$0/topics/$1",
		Attributes:   []string{"project", "name"},
	},
	"google_service_account": {
		StringFormat: "projects/$0/serviceAccounts/$1",
		Attributes:   []string{"project", "email"},
	},
	"google_storage_bucket_object": {
		StringFormat: "$0/$1",
		Attributes:   []string{"project", "name"},
	},
	"google_api_gateway_api": {
		StringFormat: "projects/$0/locations/global/apis/$1",
		Attributes:   []string{"project", "api_id"},
	},
	"google_compute_autoscaler": {
		StringFormat: "projects/$0/zones/$1/autoscalers/$2",
		Attributes:   []string{"project", "zone", "name"},
	},
	"google_compute_backend_bucket": {
		StringFormat: "projects/$0/global/backendBuckets/$1",
		Attributes:   []string{"project", "name"},
	},
	"google_compute_backend_service": {
		StringFormat: "projects/$0/global/backendServices/$1",
		Attributes:   []string{"project", "name"},
	},
	"google_bigquery_dataset": {
		StringFormat: "projects/$0/datasets/$1",
		Attributes:   []string{"project", "dataset_id"},
	},
	"google_bigquery_table": {
		StringFormat: "projects/$0/datasets/$1/tables/$2",
		Attributes:   []string{"project", "dataset_id", "table_id"},
	},
	"google_cloudbuild_trigger": {
		StringFormat: "projects/$0/triggers/$1",
		Attributes:   []string{"project", "trigger_id"},
	},
	"google_cloudfunctions_function": {
		StringFormat: "projects/$0/locations/$1/functions/$2",
		Attributes:   []string{"project", "location", "name"},
	},
	"google_sql_database_instance": {
		StringFormat: "projects/$0/instances/$1",
		Attributes:   []string{"project", "name"},
	},
	"google_sql_database": {
		StringFormat: "projects/$0/instances/$1/databases/$2",
		Attributes:   []string{"project", "instance", "name"},
	},
	"google_dataproc_cluster": {
		StringFormat: "projects/$0/regions/$1/clusters/$2",
		Attributes:   []string{"project", "region", "name"},
	},
	"google_compute_external_vpn_gateway": {
		StringFormat: "projects/$0/global/externalVpnGateways/$1",
		Attributes:   []string{"project", "name"},
	},
	"google_dns_record_set": {
		StringFormat: "projects/$0/managedZones/$1/rrsets/$2/$3",
		Attributes:   []string{"project", "zone", "name", "type"},
	},
	"google_compute_forwarding_rule": {
		StringFormat: "projects/$0/regions/$1/forwardingRules/$2",
		Attributes:   []string{"project", "region", "name"},
	},
	"google_storage_bucket_iam_binding": {
		StringFormat: "b/$0/iam",
		Attributes:   []string{"bucket"},
	},
	"google_storage_bucket_iam_member": {
		StringFormat: "b/$0/iam",
		Attributes:   []string{"bucket"},
	},
	"google_storage_notification": {
		StringFormat: "projects/$0/buckets/$1/notificationConfigs/$2",
		Attributes:   []string{"project", "bucket", "name"},
	},
	"google_container_node_pool": {
		StringFormat: "projects/$0/locations/$1/clusters/$2/nodePools/$3",
		Attributes:   []string{"project", "location", "cluster", "name"},
	},
	"google_compute_global_address": {
		StringFormat: "projects/$0/global/addresses/$1",
		Attributes:   []string{"project", "name"},
	},
	"google_compute_global_forwarding_rule": {
		StringFormat: "projects/$0/global/forwardingRules/$1",
		Attributes:   []string{"project", "name"},
	},
	"google_compute_health_check": {
		StringFormat: "projects/$0/global/healthChecks/$1",
		Attributes:   []string{"project", "name"},
	},
	"google_compute_http_health_check": {
		StringFormat: "projects/$0/global/httpHealthChecks/$1",
		Attributes:   []string{"project", "name"},
	},
	"google_compute_https_health_check": {
		StringFormat: "projects/$0/global/httpsHealthChecks/$1",
		Attributes:   []string{"project", "name"},
	},
	"google_project_iam_custom_role": {
		StringFormat: "projects/$0/roles/$1",
		Attributes:   []string{"project", "name"},
	},
	"google_project_iam_member": {
		StringFormat: "projects/$0/iam",
		Attributes:   []string{"project"},
	},
	"google_compute_instance_group_manager": {
		StringFormat: "projects/$0/zones/$1/instanceGroupManagers/$2",
		Attributes:   []string{"project", "zone", "name"},
	},
	"google_compute_instance_group": {
		StringFormat: "projects/$0/zones/$1/instanceGroups/$2",
		Attributes:   []string{"project", "zone", "name"},
	},
	"google_compute_instance_template": {
		StringFormat: "projects/$0/global/instanceTemplates/$1",
		Attributes:   []string{"project", "name"},
	},
	"google_compute_interconnect_attachment": {
		StringFormat: "projects/$0/regions/$1/interconnectAttachments/$2",
		Attributes:   []string{"project", "region", "name"},
	},
	"google_logging_metric": {
		StringFormat: "projects/$0/metrics/$1",
		Attributes:   []string{"project", "name"},
	},
	"google_redis_instance": {
		StringFormat: "projects/$0/locations/$1/instances/$2",
		Attributes:   []string{"project", "location", "name"},
	},
	"google_monitoring_alert_policy": {
		StringFormat: "projects/$0/alertPolicies/$1",
		Attributes:   []string{"project", "name"},
	},
	"google_monitoring_group": {
		StringFormat: "projects/$0/groups/$1",
		Attributes:   []string{"project", "name"},
	},
	"google_monitoring_notification_channel": {
		StringFormat: "projects/$0/notificationChannels/$1",
		Attributes:   []string{"project", "name"},
	},
	"google_monitoring_uptime_check_config": {
		StringFormat: "projects/$0/uptimeCheckConfigs/$1",
		Attributes:   []string{"project", "name"},
	},
	"google_compute_packet_mirroring": {
		StringFormat: "projects/$0/regions/$1/packetMirrorings/$2",
		Attributes:   []string{"project", "region", "name"},
	},
	"google_compute_node_group": {
		StringFormat: "projects/$0/zones/$1/nodeGroups/$2",
		Attributes:   []string{"project", "zone", "name"},
	},
	"google_compute_node_template": {
		StringFormat: "projects/$0/regions/$1/nodeTemplates/$2",
		Attributes:   []string{"project", "region", "name"},
	},
	"google_project": {
		StringFormat: "projects/$0",
		Attributes:   []string{"project_id"},
	},
	"google_compute_region_autoscaler": {
		StringFormat: "projects/$0/regions/$1/autoscalers/$2",
		Attributes:   []string{"project", "region", "name"},
	},
	"google_compute_region_backend_service": {
		StringFormat: "projects/$0/regions/$1/backendServices/$2",
		Attributes:   []string{"project", "region", "name"},
	},
	"google_compute_region_disk": {
		StringFormat: "projects/$0/regions/$1/disks/$2",
		Attributes:   []string{"project", "region", "name"},
	},
	"google_compute_region_health_check": {
		StringFormat: "projects/$0/regions/$1/healthChecks/$2",
		Attributes:   []string{"project", "region", "name"},
	},
	"google_compute_region_instance_group": {
		StringFormat: "projects/$0/regions/$1/instanceGroups/$2",
		Attributes:   []string{"project", "region", "name"},
	},
	"google_compute_region_ssl_certificate": {
		StringFormat: "projects/$0/regions/$1/sslCertificates/$2",
		Attributes:   []string{"project", "region", "name"},
	},
	"google_compute_region_target_http_proxy": {
		StringFormat: "projects/$0/regions/$1/targetHttpProxies/$2",
		Attributes:   []string{"project", "region", "name"},
	},
	"google_compute_region_target_https_proxy": {
		StringFormat: "projects/$0/regions/$1/targetHttpsProxies/$2",
		Attributes:   []string{"project", "region", "name"},
	},
	"google_compute_region_url_map": {
		StringFormat: "projects/$0/regions/$1/urlMaps/$2",
		Attributes:   []string{"project", "region", "name"},
	},
	"google_compute_reservation": {
		StringFormat: "projects/$0/zones/$1/reservations/$2",
		Attributes:   []string{"project", "zone", "name"},
	},
	"google_compute_resource_policy": {
		StringFormat: "projects/$0/regions/$1/resourcePolicies/$2",
		Attributes:   []string{"project", "region", "name"},
	},
	"google_compute_region_instance_group_manager": {
		StringFormat: "projects/$0/regions/$1/instanceGroupManagers/$2",
		Attributes:   []string{"project", "region", "name"},
	},
	"google_compute_route": {
		StringFormat: "projects/$0/global/routes/$1",
		Attributes:   []string{"project", "name"},
	},
	"google_cloud_scheduler_job": {
		StringFormat: "projects/$0/locations/$1/jobs/$2",
		Attributes:   []string{"project", "location", "name"},
	},
	"google_compute_security_policy": {
		StringFormat: "projects/$0/global/securityPolicies/$1",
		Attributes:   []string{"project", "name"},
	},
	"google_compute_managed_ssl_certificate": {
		StringFormat: "projects/$0/global/sslCertificates/$1",
		Attributes:   []string{"project", "name"},
	},
	"google_compute_ssl_policy": {
		StringFormat: "projects/$0/global/sslPolicies/$1",
		Attributes:   []string{"project", "name"},
	},
	"google_compute_target_http_proxy": {
		StringFormat: "projects/$0/global/targetHttpProxies/$1",
		Attributes:   []string{"project", "name"},
	},
	"google_compute_target_https_proxy": {
		StringFormat: "projects/$0/global/targetHttpsProxies/$1",
		Attributes:   []string{"project", "name"},
	},
	"google_compute_target_instance": {
		StringFormat: "projects/$0/zones/$1/targetInstances/$2",
		Attributes:   []string{"project", "zone", "name"},
	},
	"google_compute_target_pool": {
		StringFormat: "projects/$0/regions/$1/targetPools/$2",
		Attributes:   []string{"project", "region", "name"},
	},
	"google_compute_target_ssl_proxy": {
		StringFormat: "projects/$0/global/targetSslProxies/$1",
		Attributes:   []string{"project", "name"},
	},
	"google_compute_target_tcp_proxy": {
		StringFormat: "projects/$0/global/targetTcpProxies/$1",
		Attributes:   []string{"project", "name"},
	},
	"google_compute_vpn_gateway": {
		StringFormat: "projects/$0/regions/$1/vpnGateways/$2",
		Attributes:   []string{"project", "region", "name"},
	},
	"google_compute_url_map": {
		StringFormat: "projects/$0/global/urlMaps/$1",
		Attributes:   []string{"project", "name"},
	},
	"google_compute_vpn_tunnel": {
		StringFormat: "projects/$0/regions/$1/vpnTunnels/$2",
		Attributes:   []string{"project", "region", "name"},
	},
}
