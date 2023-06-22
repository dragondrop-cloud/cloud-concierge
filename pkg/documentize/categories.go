package documentize

// awsResourceCategories stores categories for different terraform resources for aws.
// Possible AWS Categories: compute, security, serverless, storage, database, application integration,
// networking, operations, ci cd, containers, tools, analytics, and artificial intelligence
func awsResourceCategories() TypeToCategory {
	return TypeToCategory{
		"aws_accessanalyzer_analyzer": ResourceCategory{
			primaryCat: "security",
		},
		"aws_acm_certificate": ResourceCategory{
			primaryCat: "networking",
		},
		"aws_lb": ResourceCategory{
			primaryCat: "networking",
		},
		"aws_lb_listener": ResourceCategory{
			primaryCat: "networking",
		},
		"aws_lb_listener_rule": ResourceCategory{
			primaryCat: "networking",
		},
		"aws_lb_listener_certificate": ResourceCategory{
			primaryCat: "networking",
		},
		"aws_lb_target_group": ResourceCategory{
			primaryCat: "networking",
		},
		"aws_lb_target_group_attachment": ResourceCategory{
			primaryCat: "networking",
		},
		"aws_api_gateway_authorizer": ResourceCategory{
			primaryCat:   "networking",
			secondaryCat: "security",
		},
		"aws_api_gateway_api_key": ResourceCategory{
			primaryCat: "networking",
		},
		"aws_api_gateway_documentation_part": ResourceCategory{
			primaryCat: "networking",
		},
		"aws_api_gateway_integration": ResourceCategory{
			primaryCat: "networking",
		},
		"aws_api_gateway_integration_response": ResourceCategory{
			primaryCat: "networking",
		},
		"aws_api_gateway_method": ResourceCategory{
			primaryCat: "networking",
		},
		"aws_api_gateway_method_response": ResourceCategory{
			primaryCat: "networking",
		},
		"aws_api_gateway_model": ResourceCategory{
			primaryCat: "networking",
		},
		"aws_api_gateway_resource": ResourceCategory{
			primaryCat: "networking",
		},
		"aws_api_gateway_rest_api": ResourceCategory{
			primaryCat: "networking",
		},
		"aws_api_gateway_stage": ResourceCategory{
			primaryCat: "networking",
		},
		"aws_api_gateway_usage_plan": ResourceCategory{
			primaryCat: "networking",
		},
		"aws_api_gateway_vpc_link": ResourceCategory{
			primaryCat: "networking",
		},
		"aws_appsync_graphql_api": ResourceCategory{
			primaryCat:   "compute",
			secondaryCat: "serverless",
		},
		"aws_autoscaling_group": ResourceCategory{
			primaryCat: "compute",
		},
		"aws_launch_configuration": ResourceCategory{
			primaryCat: "compute",
		},
		"aws_launch_template": ResourceCategory{
			primaryCat: "compute",
		},
		"aws_batch_compute_environment": ResourceCategory{
			primaryCat: "compute",
		},
		"aws_batch_job_definition": ResourceCategory{
			primaryCat: "compute",
		},
		"aws_batch_job_queue": ResourceCategory{
			primaryCat: "compute",
		},
		"aws_budgets_budget": ResourceCategory{
			primaryCat: "operations",
		},
		"aws_cloud9_environment_ec2": ResourceCategory{
			primaryCat:   "operations",
			secondaryCat: "compute",
		},
		"aws_cloud_formation_stack": ResourceCategory{
			primaryCat:   "tools",
			secondaryCat: "ci cd",
		},
		"aws_cloud_formation_stack_set": ResourceCategory{
			primaryCat:   "tools",
			secondaryCat: "ci cd",
		},
		"aws_cloud_formation_stack_set_instance": ResourceCategory{
			primaryCat:   "tools",
			secondaryCat: "ci cd",
		},
		"aws_cloudfront_distribution": ResourceCategory{
			primaryCat: "networking",
		},
		"aws_cloudfront_cache_policy": ResourceCategory{
			primaryCat: "networking",
		},
		"aws_cloudhsm_v2_cluster": ResourceCategory{
			primaryCat: "security",
		},
		"aws_cloudhsm_v2_hsm": ResourceCategory{
			primaryCat: "security",
		},
		"aws_cloudtrail": ResourceCategory{
			primaryCat:   "operations",
			secondaryCat: "tools",
		},
		"aws_cloudwatch_dashboard": ResourceCategory{
			primaryCat:   "operations",
			secondaryCat: "tools",
		},
		"aws_cloudwatch_event_rule": ResourceCategory{
			primaryCat:   "operations",
			secondaryCat: "tools",
		},
		"aws_cloudwatch_event_target": ResourceCategory{
			primaryCat:   "operations",
			secondaryCat: "tools",
		},
		"aws_cloudwatch_metric_alarm": ResourceCategory{
			primaryCat:   "operations",
			secondaryCat: "tools",
		},
		"aws_codebuild_project": ResourceCategory{
			primaryCat: "ci cd",
		},
		"aws_codecommit_repository": ResourceCategory{
			primaryCat:   "operations",
			secondaryCat: "tools",
		},
		"aws_codedeploy_app": ResourceCategory{
			primaryCat: "ci cd",
		},
		"aws_codepipeline": ResourceCategory{
			primaryCat: "ci cd",
		},
		"aws_codepipeline_webhook": ResourceCategory{
			primaryCat: "ci cd",
		},
		"aws_cognito_identity_pool": ResourceCategory{
			primaryCat: "security",
		},
		"aws_cognito_user_pool": ResourceCategory{
			primaryCat: "security",
		},
		"aws_config_rule": ResourceCategory{
			primaryCat:   "security",
			secondaryCat: "operations",
		},
		"aws_config_configuration_recorder": ResourceCategory{
			primaryCat:   "security",
			secondaryCat: "operations",
		},
		"aws_config_delivery_channel": ResourceCategory{
			primaryCat:   "security",
			secondaryCat: "operations",
		},
		"aws_customer_gateway": ResourceCategory{
			primaryCat: "networking",
		},
		"aws_datapipeline_pipeline": ResourceCategory{
			primaryCat:   "analytics",
			secondaryCat: "application integration",
		},
		"aws_devicefarm_project": ResourceCategory{
			primaryCat: "tools",
		},
		"aws_docdb_cluster": ResourceCategory{
			primaryCat: "database",
		},
		"aws_docdb_cluster_instance": ResourceCategory{
			primaryCat: "database",
		},
		"aws_docdb_cluster_parameter_group": ResourceCategory{
			primaryCat: "database",
		},
		"aws_docdb_subnet_group": ResourceCategory{
			primaryCat: "database",
		},
		"aws_dynamodb_table": ResourceCategory{
			primaryCat: "database",
		},
		"aws_ebs_volume": ResourceCategory{
			primaryCat: "storage",
		},
		"aws_volume_attachment": ResourceCategory{
			primaryCat: "storage",
		},
		"aws_instance": ResourceCategory{
			primaryCat: "compute",
		},
		"aws_ecr_lifecycle_policy": ResourceCategory{
			primaryCat:   "containers",
			secondaryCat: "storage",
		},
		"aws_ecr_repository": ResourceCategory{
			primaryCat:   "containers",
			secondaryCat: "storage",
		},
		"aws_ecr_repository_policy": ResourceCategory{
			primaryCat:   "containers",
			secondaryCat: "storage",
		},
		"aws_ecrpublic_repository": ResourceCategory{
			primaryCat:   "containers",
			secondaryCat: "storage",
		},
		"aws_ecs_cluster": ResourceCategory{
			primaryCat: "compute",
		},
		"aws_ecs_service": ResourceCategory{
			primaryCat: "compute",
		},
		"aws_ecs_task_definition": ResourceCategory{
			primaryCat: "compute",
		},
		"aws_efs_access_point": ResourceCategory{
			primaryCat:   "storage",
			secondaryCat: "serverless",
		},
		"aws_efs_file_system": ResourceCategory{
			primaryCat:   "storage",
			secondaryCat: "serverless",
		},
		"aws_efs_file_system_policy": ResourceCategory{
			primaryCat:   "storage",
			secondaryCat: "serverless",
		},
		"aws_efs_mount_target": ResourceCategory{
			primaryCat:   "storage",
			secondaryCat: "serverless",
		},
		"aws_eip": ResourceCategory{
			primaryCat: "networking",
		},
		"aws_eks_cluster": ResourceCategory{
			primaryCat:   "compute",
			secondaryCat: "containers",
		},
		"aws_eks_node_group": ResourceCategory{
			primaryCat:   "compute",
			secondaryCat: "containers",
		},
		"aws_elasticache_cluster": ResourceCategory{
			primaryCat: "database",
		},
		"aws_elasticache_parameter_group": ResourceCategory{
			primaryCat: "database",
		},
		"aws_elasticache_subnet_group": ResourceCategory{
			primaryCat: "database",
		},
		"aws_elasticache_replication_group": ResourceCategory{
			primaryCat: "database",
		},
		"aws_elastic_beanstalk_application": ResourceCategory{
			primaryCat:   "application integration",
			secondaryCat: "tools",
		},
		"aws_elastic_beanstalk_environment": ResourceCategory{
			primaryCat:   "application integration",
			secondaryCat: "tools",
		},
		"aws_elb": ResourceCategory{
			primaryCat: "networking",
		},
		"aws_emr_cluster": ResourceCategory{
			primaryCat:   "compute",
			secondaryCat: "analytics",
		},
		"aws_emr_security_configuration": ResourceCategory{
			primaryCat:   "compute",
			secondaryCat: "analytics",
		},
		"aws_network_interface": ResourceCategory{
			primaryCat: "networking",
		},
		"aws_elasticsearch_domain": ResourceCategory{
			primaryCat:   "database",
			secondaryCat: "analytics",
		},
		"aws_kinesis_firehose_delivery_stream": ResourceCategory{
			primaryCat:   "application integration",
			secondaryCat: "analytics",
		},
		"aws_glue_crawler": ResourceCategory{
			primaryCat:   "analytics",
			secondaryCat: "application integration",
		},
		"aws_glue_catalog_database": ResourceCategory{
			primaryCat:   "analytics",
			secondaryCat: "application integration",
		},
		"aws_glue_catalog_table": ResourceCategory{
			primaryCat:   "analytics",
			secondaryCat: "application integration",
		},
		"aws_glue_job": ResourceCategory{
			primaryCat:   "analytics",
			secondaryCat: "application integration",
		},
		"aws_glue_trigger": ResourceCategory{
			primaryCat:   "analytics",
			secondaryCat: "application integration",
		},
		"aws_iam_access_key": ResourceCategory{
			primaryCat:   "security",
			secondaryCat: "operations",
		},
		"aws_iam_group": ResourceCategory{
			primaryCat:   "security",
			secondaryCat: "operations",
		},
		"aws_iam_group_policy": ResourceCategory{
			primaryCat:   "security",
			secondaryCat: "operations",
		},
		"aws_iam_group_policy_attachment": ResourceCategory{
			primaryCat:   "security",
			secondaryCat: "operations",
		},
		"aws_iam_instance_profile": ResourceCategory{
			primaryCat:   "security",
			secondaryCat: "operations",
		},
		"aws_iam_policy": ResourceCategory{
			primaryCat:   "security",
			secondaryCat: "operations",
		},
		"aws_iam_role": ResourceCategory{
			primaryCat:   "security",
			secondaryCat: "operations",
		},
		"aws_iam_role_policy": ResourceCategory{
			primaryCat:   "security",
			secondaryCat: "operations",
		},
		"aws_iam_role_policy_attachment": ResourceCategory{
			primaryCat:   "security",
			secondaryCat: "operations",
		},
		"aws_iam_user": ResourceCategory{
			primaryCat:   "security",
			secondaryCat: "operations",
		},
		"aws_iam_user_group_membership": ResourceCategory{
			primaryCat:   "security",
			secondaryCat: "operations",
		},
		"aws_iam_user_policy": ResourceCategory{
			primaryCat:   "security",
			secondaryCat: "operations",
		},
		"aws_iam_user_policy_attachment": ResourceCategory{
			primaryCat:   "security",
			secondaryCat: "operations",
		},
		"aws_internet_gateway": ResourceCategory{
			primaryCat: "networking",
		},
		"aws_iot_thing": ResourceCategory{
			primaryCat:   "iot",
			secondaryCat: "application integration",
		},
		"aws_iot_thing_type": ResourceCategory{
			primaryCat:   "iot",
			secondaryCat: "application integration",
		},
		"aws_iot_topic_rule": ResourceCategory{
			primaryCat:   "iot",
			secondaryCat: "application integration",
		},
		"aws_iot_role_alias": ResourceCategory{
			primaryCat:   "iot",
			secondaryCat: "application integration",
		},
		"aws_kinesis_stream": ResourceCategory{
			primaryCat:   "data",
			secondaryCat: "application integration",
		},
		"aws_kms_key": ResourceCategory{
			primaryCat: "security",
		},
		"aws_kms_alias": ResourceCategory{
			primaryCat: "security",
		},
		"aws_kms_grant": ResourceCategory{
			primaryCat: "security",
		},
		"aws_lambda_event_source_mapping": ResourceCategory{
			primaryCat:   "serverless",
			secondaryCat: "compute",
		},
		"aws_lambda_function": ResourceCategory{
			primaryCat:   "serverless",
			secondaryCat: "compute",
		},
		"aws_lambda_function_event_invoke_config": ResourceCategory{
			primaryCat:   "serverless",
			secondaryCat: "compute",
		},
		"aws_lambda_layer_version": ResourceCategory{
			primaryCat:   "serverless",
			secondaryCat: "compute",
		},
		"aws_lambda_permission": ResourceCategory{
			primaryCat:   "serverless",
			secondaryCat: "compute",
		},
		"aws_cloudwatch_log_group": ResourceCategory{
			primaryCat: "operations",
		},
		"aws_media_package_channel": ResourceCategory{
			primaryCat: "tools",
		},
		"aws_media_store_container": ResourceCategory{
			primaryCat: "storage",
		},
		"aws_medialive_channel": ResourceCategory{
			primaryCat: "tools",
		},
		"aws_medialive_input": ResourceCategory{
			primaryCat: "tools",
		},
		"aws_medialive_input_security_group": ResourceCategory{
			primaryCat: "tools",
		},
		"aws_msk_cluster": ResourceCategory{
			primaryCat:   "compute",
			secondaryCat: "analytics",
		},
		"aws_network_acl": ResourceCategory{
			primaryCat:   "security",
			secondaryCat: "networking",
		},
		"aws_nat_gateway": ResourceCategory{
			primaryCat: "networking",
		},
		"aws_opsworks_application": ResourceCategory{
			primaryCat:   "ci cd",
			secondaryCat: "application integration",
		},
		"aws_opsworks_custom_layer": ResourceCategory{
			primaryCat:   "ci cd",
			secondaryCat: "application integration",
		},
		"aws_opsworks_instance": ResourceCategory{
			primaryCat:   "ci cd",
			secondaryCat: "application integration",
		},
		"aws_opsworks_java_app_layer": ResourceCategory{
			primaryCat:   "ci cd",
			secondaryCat: "application integration",
		},
		"aws_opsworks_php_app_layer": ResourceCategory{
			primaryCat:   "ci cd",
			secondaryCat: "application integration",
		},
		"aws_opsworks_rds_db_instance": ResourceCategory{
			primaryCat:   "ci cd",
			secondaryCat: "application integration",
		},
		"aws_opsworks_stack": ResourceCategory{
			primaryCat:   "ci cd",
			secondaryCat: "application integration",
		},
		"aws_opsworks_static_web_layer": ResourceCategory{
			primaryCat:   "ci cd",
			secondaryCat: "application integration",
		},
		"aws_opsworks_user_profile": ResourceCategory{
			primaryCat:   "ci cd",
			secondaryCat: "application integration",
		},
		"aws_organizations_account": ResourceCategory{
			primaryCat:   "security",
			secondaryCat: "operations",
		},
		"aws_organizations_organization": ResourceCategory{
			primaryCat:   "security",
			secondaryCat: "operations",
		},
		"aws_organizations_organizational_unit": ResourceCategory{
			primaryCat:   "security",
			secondaryCat: "operations",
		},
		"aws_organizations_policy": ResourceCategory{
			primaryCat:   "security",
			secondaryCat: "operations",
		},
		"aws_organizations_policy_attachment": ResourceCategory{
			primaryCat:   "security",
			secondaryCat: "operations",
		},
		"aws_qldb_ledger": ResourceCategory{
			primaryCat: "database",
		},
		"aws_db_instance": ResourceCategory{
			primaryCat:   "database",
			secondaryCat: "compute",
		},
		"aws_db_proxy": ResourceCategory{
			primaryCat: "database",
		},
		"aws_db_cluster": ResourceCategory{
			primaryCat:   "database",
			secondaryCat: "compute",
		},
		"aws_db_cluster_snapshot": ResourceCategory{
			primaryCat: "database",
		},
		"aws_db_parameter_group": ResourceCategory{
			primaryCat: "database",
		},
		"aws_db_snapshot": ResourceCategory{
			primaryCat:   "database",
			secondaryCat: "storage",
		},
		"aws_db_subnet_group": ResourceCategory{
			primaryCat:   "database",
			secondaryCat: "network",
		},
		"aws_db_option_group": ResourceCategory{
			primaryCat: "database",
		},
		"aws_db_event_subscription": ResourceCategory{
			primaryCat: "database",
		},
		"aws_rds_global_cluster": ResourceCategory{
			primaryCat: "database",
		},
		"aws_route53_zone": ResourceCategory{
			primaryCat: "networking",
		},
		"aws_route53_record": ResourceCategory{
			primaryCat: "networking",
		},
		"aws_route_table": ResourceCategory{
			primaryCat: "networking",
		},
		"aws_main_route_table_association": ResourceCategory{
			primaryCat: "networking",
		},
		"aws_route_table_association": ResourceCategory{
			primaryCat: "networking",
		},
		"aws_s3_bucket": ResourceCategory{
			primaryCat: "storage",
		},
		"aws_secretsmanager_secret": ResourceCategory{
			primaryCat: "security",
		},
		"aws_securityhub_account": ResourceCategory{
			primaryCat: "security",
		},
		"aws_securityhub_member": ResourceCategory{
			primaryCat: "security",
		},
		"aws_securityhub_standards_subscription": ResourceCategory{
			primaryCat: "security",
		},
		"aws_servicecatalog_portfolio": ResourceCategory{
			primaryCat:   "operations",
			secondaryCat: "tools",
		},
		"aws_ses_configuration_set": ResourceCategory{
			primaryCat:   "application integration",
			secondaryCat: "tools",
		},
		"aws_ses_domain_identity": ResourceCategory{
			primaryCat:   "application integration",
			secondaryCat: "tools",
		},
		"aws_ses_email_identity": ResourceCategory{
			primaryCat:   "application integration",
			secondaryCat: "tools",
		},
		"aws_ses_receipt_rule": ResourceCategory{
			primaryCat:   "application integration",
			secondaryCat: "tools",
		},
		"aws_ses_receipt_rule_set": ResourceCategory{
			primaryCat:   "application integration",
			secondaryCat: "tools",
		},
		"aws_ses_template": ResourceCategory{
			primaryCat:   "application integration",
			secondaryCat: "tools",
		},
		"aws_sfn_activity": ResourceCategory{
			primaryCat:   "operations",
			secondaryCat: "tools",
		},
		"aws_sfn_state_machine": ResourceCategory{
			primaryCat:   "operations",
			secondaryCat: "tools",
		},
		"aws_security_group": ResourceCategory{
			primaryCat: "security",
		},
		"aws_security_group_rule": ResourceCategory{
			primaryCat: "security",
		},
		"aws_sns_topic": ResourceCategory{
			primaryCat: "application integration",
		},
		"aws_sns_topic_subscription": ResourceCategory{
			primaryCat: "application integration",
		},
		"aws_sqs_queue": ResourceCategory{
			primaryCat: "application integration",
		},
		"aws_ssm_parameter": ResourceCategory{
			primaryCat:   "operations",
			secondaryCat: "security",
		},
		"aws_subnet": ResourceCategory{
			primaryCat: "networking",
		},
		"aws_swf_domain": ResourceCategory{
			primaryCat:   "operations",
			secondaryCat: "tools",
		},
		"aws_ec2_transit_gateway_route_table": ResourceCategory{
			primaryCat: "networking",
		},
		"aws_ec2_transit_gateway_vpc_attachment": ResourceCategory{
			primaryCat: "networking",
		},
		"aws_vpc": ResourceCategory{
			primaryCat: "networking",
		},
		"aws_vpc_peering_connection": ResourceCategory{
			primaryCat: "networking",
		},
		"aws_vpn_connection": ResourceCategory{
			primaryCat: "networking",
		},
		"aws_vpn_gateway": ResourceCategory{
			primaryCat: "networking",
		},
		"aws_wafregional_byte_match_set": ResourceCategory{
			primaryCat:   "networking",
			secondaryCat: "security",
		},
		"aws_wafregional_geo_match_set": ResourceCategory{
			primaryCat:   "networking",
			secondaryCat: "security",
		},
		"aws_wafregional_ipset": ResourceCategory{
			primaryCat:   "networking",
			secondaryCat: "security",
		},
		"aws_wafregional_rate_based_rule": ResourceCategory{
			primaryCat:   "networking",
			secondaryCat: "security",
		},
		"aws_wafregional_regex_match_set": ResourceCategory{
			primaryCat:   "networking",
			secondaryCat: "security",
		},
		"aws_wafregional_regex_pattern_set": ResourceCategory{
			primaryCat:   "networking",
			secondaryCat: "security",
		},
		"aws_wafregional_rule": ResourceCategory{
			primaryCat:   "networking",
			secondaryCat: "security",
		},
		"aws_wafregional_rule_group": ResourceCategory{
			primaryCat:   "networking",
			secondaryCat: "security",
		},
		"aws_wafregional_size_constraint_set": ResourceCategory{
			primaryCat:   "networking",
			secondaryCat: "security",
		},
		"aws_wafregional_sql_injection_match_set": ResourceCategory{
			primaryCat:   "networking",
			secondaryCat: "security",
		},
		"aws_wafregional_web_acl": ResourceCategory{
			primaryCat:   "networking",
			secondaryCat: "security",
		},
		"aws_wafregional_xss_match_set": ResourceCategory{
			primaryCat:   "networking",
			secondaryCat: "security",
		},
		"aws_wafv2_ip_set": ResourceCategory{
			primaryCat:   "networking",
			secondaryCat: "security",
		},
		"aws_wafv2_regex_pattern_set": ResourceCategory{
			primaryCat:   "networking",
			secondaryCat: "security",
		},
		"aws_wafv2_rule_group": ResourceCategory{
			primaryCat:   "networking",
			secondaryCat: "security",
		},
		"aws_wafv2_web_acl": ResourceCategory{
			primaryCat:   "networking",
			secondaryCat: "security",
		},
		"aws_wafv2_web_acl_association": ResourceCategory{
			primaryCat:   "networking",
			secondaryCat: "security",
		},
		"aws_wafv2_web_acl_logging_configuration": ResourceCategory{
			primaryCat:   "networking",
			secondaryCat: "security",
		},
		"aws_workspaces_directory": ResourceCategory{
			primaryCat:   "compute",
			secondaryCat: "tools",
		},
		"aws_workspaces_ip_group": ResourceCategory{
			primaryCat:   "compute",
			secondaryCat: "tools",
		},
		"aws_workspaces_workspace": ResourceCategory{
			primaryCat:   "compute",
			secondaryCat: "tools",
		},
		"aws_xray_sampling_rule": ResourceCategory{
			primaryCat:   "security",
			secondaryCat: "networking",
		},
	}
}

// googleResourceCategories stores categories for different terraform resources for gcp.
// Possible GCP Categories: compute, security, serverless, storage, database, application integration,
// networking, operations, ci cd, containers, tools, analytics, and artificial intelligence
func googleResourceCategories() TypeToCategory {
	return TypeToCategory{
		"google_compute_address": ResourceCategory{
			primaryCat: "compute",
		},
		"google_compute_autoscaler": ResourceCategory{
			primaryCat: "compute",
		},
		"google_compute_backend_bucket": ResourceCategory{
			primaryCat:   "compute",
			secondaryCat: "storage",
		},
		"google_compute_backend_service": ResourceCategory{
			primaryCat: "compute",
		},
		"google_bigquery_dataset": ResourceCategory{
			primaryCat:   "storage",
			secondaryCat: "analytics",
		},
		"google_bigquery_table": ResourceCategory{
			primaryCat:   "storage",
			secondaryCat: "analytics",
		},
		"google_cloudbuild_trigger": ResourceCategory{
			primaryCat: "ci cd",
		},
		"google_cloudfunctions_function": ResourceCategory{
			primaryCat:   "compute",
			secondaryCat: "serverless",
		},
		"google_sql_database_instance": ResourceCategory{
			primaryCat:   "storage",
			secondaryCat: "database",
		},
		"google_sql_database": ResourceCategory{
			primaryCat:   "storage",
			secondaryCat: "database",
		},
		"google_dataproc_cluster": ResourceCategory{
			primaryCat:   "analytics",
			secondaryCat: "compute",
		},
		"google_compute_disk": ResourceCategory{
			primaryCat:   "compute",
			secondaryCat: "storage",
		},
		"google_compute_external_vpn_gateway": ResourceCategory{
			primaryCat:   "compute",
			secondaryCat: "networking",
		},
		"google_dns_managed_zone": ResourceCategory{
			primaryCat: "networking",
		},
		"google_dns_record_set": ResourceCategory{
			primaryCat: "networking",
		},
		"google_compute_firewall": ResourceCategory{
			primaryCat:   "compute",
			secondaryCat: "networking",
		},
		"google_compute_forwarding_rule": ResourceCategory{
			primaryCat:   "compute",
			secondaryCat: "networking",
		},
		"google_storage_bucket": ResourceCategory{
			primaryCat: "storage",
		},
		"google_storage_bucket_acl": ResourceCategory{
			primaryCat:   "storage",
			secondaryCat: "security",
		},
		"google_storage_default_object_acl": ResourceCategory{
			primaryCat:   "storage",
			secondaryCat: "security",
		},
		"google_storage_bucket_iam_binding": ResourceCategory{
			primaryCat:   "storage",
			secondaryCat: "security",
		},
		"google_storage_bucket_iam_member": ResourceCategory{
			primaryCat:   "storage",
			secondaryCat: "security",
		},
		"google_storage_bucket_iam_policy": ResourceCategory{
			primaryCat:   "storage",
			secondaryCat: "security",
		},
		"google_storage_notification": ResourceCategory{
			primaryCat:   "storage",
			secondaryCat: "application integration",
		},
		"google_container_cluster": ResourceCategory{
			primaryCat:   "compute",
			secondaryCat: "containers",
		},
		"google_container_node_pool": ResourceCategory{
			primaryCat:   "compute",
			secondaryCat: "containers",
		},
		"google_compute_global_address": ResourceCategory{
			primaryCat:   "compute",
			secondaryCat: "networking",
		},
		"google_compute_global_forwarding_rule": ResourceCategory{
			primaryCat:   "compute",
			secondaryCat: "networking",
		},
		"google_compute_health_check": ResourceCategory{
			primaryCat: "compute",
		},
		"google_compute_http_health_check": ResourceCategory{
			primaryCat: "compute",
		},
		"google_compute_https_health_check": ResourceCategory{
			primaryCat: "compute",
		},
		"google_project_iam_custom_role": ResourceCategory{
			primaryCat: "security",
		},
		"google_project_iam_member": ResourceCategory{
			primaryCat: "security",
		},
		"google_service_account": ResourceCategory{
			primaryCat: "security",
		},
		"google_compute_image": ResourceCategory{
			primaryCat:   "compute",
			secondaryCat: "containers",
		},
		"google_compute_instance_group_manager": ResourceCategory{
			primaryCat:   "compute",
			secondaryCat: "operations",
		},
		"google_compute_instance_group": ResourceCategory{
			primaryCat:   "compute",
			secondaryCat: "operations",
		},
		"google_compute_instance_template": ResourceCategory{
			primaryCat:   "compute",
			secondaryCat: "operations",
		},
		"google_compute_instance": ResourceCategory{
			primaryCat: "compute",
		},
		"google_compute_interconnect_attachment": ResourceCategory{
			primaryCat:   "compute",
			secondaryCat: "networking",
		},
		"google_kms_key_ring": ResourceCategory{
			primaryCat: "security",
		},
		"google_kms_crypto_key": ResourceCategory{
			primaryCat: "security",
		},
		"google_logging_metric": ResourceCategory{
			primaryCat: "operations",
		},
		"google_redis_instance": ResourceCategory{
			primaryCat: "storage",
		},
		"google_monitoring_alert_policy": ResourceCategory{
			primaryCat: "operations",
		},
		"google_monitoring_group": ResourceCategory{
			primaryCat: "operations",
		},
		"google_monitoring_notification_channel": ResourceCategory{
			primaryCat: "operations",
		},
		"google_monitoring_uptime_check_config": ResourceCategory{
			primaryCat: "operations",
		},
		"google_compute_network": ResourceCategory{
			primaryCat: "networking",
		},
		"google_compute_packet_mirroring": ResourceCategory{
			primaryCat:   "compute",
			secondaryCat: "operations",
		},
		"google_compute_node_group": ResourceCategory{
			primaryCat:   "compute",
			secondaryCat: "operations",
		},
		"google_compute_node_template": ResourceCategory{
			primaryCat:   "compute",
			secondaryCat: "operations",
		},
		"google_pubsub_subscription": ResourceCategory{
			primaryCat: "application integration",
		},
		"google_pubsub_topic": ResourceCategory{
			primaryCat: "application integration",
		},
		"google_compute_region_autoscaler": ResourceCategory{
			primaryCat:   "compute",
			secondaryCat: "operations",
		},
		"google_compute_region_backend_service": ResourceCategory{
			primaryCat: "compute",
		},
		"google_compute_region_disk": ResourceCategory{
			primaryCat:   "compute",
			secondaryCat: "storage",
		},
		"google_compute_region_health_check": ResourceCategory{
			primaryCat:   "compute",
			secondaryCat: "operations",
		},
		"google_compute_region_instance_group": ResourceCategory{
			primaryCat:   "compute",
			secondaryCat: "operations",
		},
		"google_compute_region_target_http_proxy": ResourceCategory{
			primaryCat: "networking",
		},
		"google_compute_region_target_https_proxy": ResourceCategory{
			primaryCat: "networking",
		},
		"google_compute_region_ssl_certificate": ResourceCategory{
			primaryCat:   "networking",
			secondaryCat: "compute",
		},
		"google_compute_region_url_map": ResourceCategory{
			primaryCat: "networking",
		},
		"google_compute_reservation": ResourceCategory{
			primaryCat: "compute",
		},
		"google_compute_resource_policy": ResourceCategory{
			primaryCat:   "compute",
			secondaryCat: "security",
		},
		"google_compute_region_instance_group_manager": ResourceCategory{
			primaryCat:   "compute",
			secondaryCat: "operations",
		},
		"google_compute_router": ResourceCategory{
			primaryCat: "networking",
		},
		"google_compute_route": ResourceCategory{
			primaryCat:   "networking",
			secondaryCat: "compute",
		},
		"google_cloud_scheduler_job": ResourceCategory{
			primaryCat: "application integration",
		},
		"google_compute_security_policy": ResourceCategory{
			primaryCat:   "security",
			secondaryCat: "compute",
		},
		"google_compute_managed_ssl_certificate": ResourceCategory{
			primaryCat: "networking",
		},
		"google_compute_ssl_policy": ResourceCategory{
			primaryCat: "networking",
		},
		"google_compute_subnetwork": ResourceCategory{
			primaryCat: "networking",
		},
		"google_compute_target_http_proxy": ResourceCategory{
			primaryCat:   "networking",
			secondaryCat: "compute",
		},
		"google_compute_target_https_proxy": ResourceCategory{
			primaryCat:   "networking",
			secondaryCat: "compute",
		},
		"google_compute_target_instance": ResourceCategory{
			primaryCat: "compute",
		},
		"google_compute_target_pool": ResourceCategory{
			primaryCat:   "compute",
			secondaryCat: "operations",
		},
		"google_compute_target_ssl_proxy": ResourceCategory{
			primaryCat: "networking",
		},
		"google_compute_target_tcp_proxy": ResourceCategory{
			primaryCat: "networking",
		},
		"google_compute_vpn_gateway": ResourceCategory{
			primaryCat: "networking",
		},
		"google_compute_url_map": ResourceCategory{
			primaryCat: "networking",
		},
		"google_compute_vpn_tunnel": ResourceCategory{
			primaryCat: "networking",
		},
	}
}

func azureResourceCategories() TypeToCategory {
	return TypeToCategory{
		"azurerm_analysis_services_server": ResourceCategory{
			primaryCat: "analytics",
		},
		"azurerm_app_service": ResourceCategory{
			primaryCat: "compute",
		},
		"azurerm_application_gateway": ResourceCategory{
			primaryCat:   "networking",
			secondaryCat: "security",
		},
		"azurerm_container_group": ResourceCategory{
			primaryCat: "compute",
		},
		"azurerm_container_registry": ResourceCategory{
			primaryCat: "compute",
		},
		"azurerm_container_registry_webhook": ResourceCategory{
			primaryCat: "compute",
		},
		"azurerm_cosmosdb_account": ResourceCategory{
			primaryCat: "database",
		},
		"azurerm_cosmosdb_sql_container": ResourceCategory{
			primaryCat: "database",
		},
		"azurerm_cosmosdb_sql_database": ResourceCategory{
			primaryCat: "database",
		},
		"azurerm_cosmosdb_table": ResourceCategory{
			primaryCat: "database",
		},
		"azurerm_mariadb_configuration": ResourceCategory{
			primaryCat: "database",
		},
		"azurerm_mariadb_database": ResourceCategory{
			primaryCat: "database",
		},
		"azurerm_mariadb_firewall_rule": ResourceCategory{
			primaryCat:   "database",
			secondaryCat: "security",
		},
		"azurerm_mariadb_server": ResourceCategory{
			primaryCat: "database",
		},
		"azurerm_mariadb_virtual_network_rule": ResourceCategory{
			primaryCat:   "database",
			secondaryCat: "networking",
		},
		"azurerm_mysql_configuration": ResourceCategory{
			primaryCat: "database",
		},
		"azurerm_mysql_database": ResourceCategory{
			primaryCat: "database",
		},
		"azurerm_mysql_firewall_rule": ResourceCategory{
			primaryCat:   "database",
			secondaryCat: "security",
		},
		"azurerm_mysql_server": ResourceCategory{
			primaryCat: "database",
		},
		"azurerm_mysql_virtual_network_rule": ResourceCategory{
			primaryCat:   "database",
			secondaryCat: "networking",
		},
		"azurerm_postgresql_configuration": ResourceCategory{
			primaryCat: "database",
		},
		"azurerm_postgresql_database": ResourceCategory{
			primaryCat: "database",
		},
		"azurerm_postgresql_firewall_rule": ResourceCategory{
			primaryCat:   "database",
			secondaryCat: "security",
		},
		"azurerm_postgresql_server": ResourceCategory{
			primaryCat: "database",
		},
		"azurerm_postgresql_virtual_network_rule": ResourceCategory{
			primaryCat:   "database",
			secondaryCat: "networking",
		},
		"azurerm_sql_database": ResourceCategory{
			primaryCat: "database",
		},
		"azurerm_sql_active_directory_administrator": ResourceCategory{
			primaryCat: "database",
		},
		"azurerm_sql_elasticpool": ResourceCategory{
			primaryCat: "database",
		},
		"azurerm_sql_failover_group": ResourceCategory{
			primaryCat: "database",
		},
		"azurerm_sql_firewall_rule": ResourceCategory{
			primaryCat:   "database",
			secondaryCat: "security",
		},
		"azurerm_sql_server": ResourceCategory{
			primaryCat: "database",
		},
		"azurerm_sql_virtual_network_rule": ResourceCategory{
			primaryCat:   "database",
			secondaryCat: "networking",
		},
		"azurerm_databricks_workspace": ResourceCategory{
			primaryCat: "analytics",
		},
		"azurerm_data_factory": ResourceCategory{
			primaryCat: "data_pipeline",
		},
		"azurerm_data_factory_pipeline": ResourceCategory{
			primaryCat: "data_pipeline",
		},
		"azurerm_data_factory_data_flow": ResourceCategory{
			primaryCat: "data_pipeline",
		},
		"azurerm_data_factory_dataset_azure_blob": ResourceCategory{
			primaryCat: "data_pipeline",
		},
		"azurerm_data_factory_dataset_binary": ResourceCategory{
			primaryCat: "data_pipeline",
		},
		"azurerm_data_factory_dataset_cosmosdb_sqlapi": ResourceCategory{
			primaryCat: "data_pipeline",
		},
		"azurerm_data_factory_custom_dataset": ResourceCategory{
			primaryCat: "data_pipeline",
		},
		"azurerm_data_factory_dataset_delimited_text": ResourceCategory{
			primaryCat: "data_pipeline",
		},
		"azurerm_data_factory_dataset_http": ResourceCategory{
			primaryCat: "data_pipeline",
		},
		"azurerm_data_factory_dataset_json": ResourceCategory{
			primaryCat: "data_pipeline",
		},
		"azurerm_data_factory_dataset_mysql": ResourceCategory{
			primaryCat: "data_pipeline",
		},
		"azurerm_data_factory_dataset_parquet": ResourceCategory{
			primaryCat: "data_pipeline",
		},
		"azurerm_data_factory_dataset_postgresql": ResourceCategory{
			primaryCat: "data_pipeline",
		},
		"azurerm_data_factory_dataset_snowflake": ResourceCategory{
			primaryCat: "data_pipeline",
		},
		"azurerm_data_factory_dataset_sql_server_table": ResourceCategory{
			primaryCat: "data_pipeline",
		},
		"azurerm_data_factory_integration_runtime_azure": ResourceCategory{
			primaryCat: "data_pipeline",
		},
		"azurerm_data_factory_integration_runtime_managed": ResourceCategory{
			primaryCat: "data_pipeline",
		},
		"azurerm_data_factory_integration_runtime_azure_ssis": ResourceCategory{
			primaryCat: "data_pipeline",
		},
		"azurerm_data_factory_integration_runtime_self_hosted": ResourceCategory{
			primaryCat: "data_pipeline",
		},
		"azurerm_data_factory_linked_service_azure_blob_storage": ResourceCategory{
			primaryCat: "data_pipeline",
		},
		"azurerm_data_factory_linked_service_azure_databricks": ResourceCategory{
			primaryCat: "data_pipeline",
		},
		"azurerm_data_factory_linked_service_azure_file_storage": ResourceCategory{
			primaryCat: "data_pipeline",
		},
		"azurerm_data_factory_linked_service_azure_function": ResourceCategory{
			primaryCat: "data_pipeline",
		},
		"azurerm_data_factory_linked_service_azure_search": ResourceCategory{
			primaryCat: "data_pipeline",
		},
		"azurerm_data_factory_linked_service_azure_sql_database": ResourceCategory{
			primaryCat: "data_pipeline",
		},
		"azurerm_data_factory_linked_service_azure_table_storage": ResourceCategory{
			primaryCat: "data_pipeline",
		},
		"azurerm_data_factory_linked_service_cosmosdb": ResourceCategory{
			primaryCat: "data_pipeline",
		},
		"azurerm_data_factory_linked_custom_service": ResourceCategory{
			primaryCat: "data_pipeline",
		},
		"azurerm_data_factory_linked_service_data_lake_storage_gen2": ResourceCategory{
			primaryCat: "data_pipeline",
		},
		"azurerm_data_factory_linked_service_key_vault": ResourceCategory{
			primaryCat: "data_pipeline",
		},
		"azurerm_data_factory_linked_service_kusto": ResourceCategory{
			primaryCat: "data_pipeline",
		},
		"azurerm_data_factory_linked_service_mysql": ResourceCategory{
			primaryCat: "data_pipeline",
		},
		"azurerm_data_factory_linked_service_odata": ResourceCategory{
			primaryCat: "data_pipeline",
		},
		"azurerm_data_factory_linked_service_postgresql": ResourceCategory{
			primaryCat: "data_pipeline",
		},
		"azurerm_data_factory_linked_service_sftp": ResourceCategory{
			primaryCat: "data_pipeline",
		},
		"azurerm_data_factory_linked_service_snowflake": ResourceCategory{
			primaryCat: "data_pipeline",
		},
		"azurerm_data_factory_linked_service_sql_server": ResourceCategory{
			primaryCat: "data_pipeline",
		},
		"azurerm_data_factory_linked_service_synapse": ResourceCategory{
			primaryCat: "data_pipeline",
		},
		"azurerm_data_factory_linked_service_web": ResourceCategory{
			primaryCat: "data_pipeline",
		},
		"azurerm_data_factory_trigger_blob_event": ResourceCategory{
			primaryCat: "data_pipeline",
		},
		"azurerm_data_factory_trigger_schedule": ResourceCategory{
			primaryCat: "data_pipeline",
		},
		"azurerm_data_factory_trigger_tumbling_window": ResourceCategory{
			primaryCat: "data_pipeline",
		},
		"azurerm_managed_disk": ResourceCategory{
			primaryCat: "storage",
		},
		"azurerm_dns_a_record": ResourceCategory{
			primaryCat: "networking",
		},
		"azurerm_dns_aaaa_record": ResourceCategory{
			primaryCat: "networking",
		},
		"azurerm_dns_caa_record": ResourceCategory{
			primaryCat: "networking",
		},
		"azurerm_dns_cname_record": ResourceCategory{
			primaryCat: "networking",
		},
		"azurerm_dns_mx_record": ResourceCategory{
			primaryCat: "networking",
		},
		"azurerm_dns_ns_record": ResourceCategory{
			primaryCat: "networking",
		},
		"azurerm_dns_ptr_record": ResourceCategory{
			primaryCat: "networking",
		},
		"azurerm_dns_srv_record": ResourceCategory{
			primaryCat: "networking",
		},
		"azurerm_dns_txt_record": ResourceCategory{
			primaryCat: "networking",
		},
		"azurerm_dns_zone": ResourceCategory{
			primaryCat: "networking",
		},
		"azurerm_lb": ResourceCategory{
			primaryCat: "networking",
		},
		"azurerm_lb_backend_address_pool": ResourceCategory{
			primaryCat: "networking",
		},
		"azurerm_lb_nat_rule": ResourceCategory{
			primaryCat: "networking",
		},
		"azurerm_lb_probe": ResourceCategory{
			primaryCat: "networking",
		},
		"azurerm_eventhub_namespace": ResourceCategory{
			primaryCat: "messaging",
		},
		"azurerm_eventhub": ResourceCategory{
			primaryCat: "messaging",
		},
		"azurerm_eventhub_consumer_group": ResourceCategory{
			primaryCat: "messaging",
		},
		"azurerm_eventhub_namespace_authorization_rule": ResourceCategory{
			primaryCat: "messaging",
		},
		"azurerm_network_interface": ResourceCategory{
			primaryCat: "networking",
		},
		"azurerm_network_security_group": ResourceCategory{
			primaryCat: "networking",
		},
		"azurerm_network_security_rule": ResourceCategory{
			primaryCat: "networking",
		},
		"azurerm_network_watcher": ResourceCategory{
			primaryCat: "networking",
		},
		"azurerm_network_watcher_flow_log": ResourceCategory{
			primaryCat: "networking",
		},
		"azurerm_network_packet_capture": ResourceCategory{
			primaryCat: "networking",
		},
		"azurerm_private_dns_a_record": ResourceCategory{
			primaryCat: "networking",
		},
		"azurerm_private_dns_aaaa_record": ResourceCategory{
			primaryCat: "networking",
		},
		"azurerm_private_dns_cname_record": ResourceCategory{
			primaryCat: "networking",
		},
		"azurerm_private_dns_mx_record": ResourceCategory{
			primaryCat: "networking",
		},
		"azurerm_private_dns_ptr_record": ResourceCategory{
			primaryCat: "networking",
		},
		"azurerm_private_dns_srv_record": ResourceCategory{
			primaryCat: "networking",
		},
		"azurerm_private_dns_txt_record": ResourceCategory{
			primaryCat: "networking",
		},
		"azurerm_private_dns_zone": ResourceCategory{
			primaryCat: "networking",
		},
		"azurerm_private_dns_zone_virtual_network_link": ResourceCategory{
			primaryCat: "networking",
		},
		"azurerm_private_endpoint": ResourceCategory{
			primaryCat: "networking",
		},
		"azurerm_private_link_service": ResourceCategory{
			primaryCat: "networking",
		},
		"azurerm_public_ip": ResourceCategory{
			primaryCat: "networking",
		},
		"azurerm_public_ip_prefix": ResourceCategory{
			primaryCat: "networking",
		},
		"azurerm_redis_cache": ResourceCategory{
			primaryCat: "database",
		},
		"azurerm_purview_account": ResourceCategory{
			primaryCat: "management",
		},
		"azurerm_resource_group": ResourceCategory{
			primaryCat: "management",
		},
		"azurerm_management_lock": ResourceCategory{
			primaryCat: "management",
		},
		"azurerm_route_table": ResourceCategory{
			primaryCat: "networking",
		},
		"azurerm_route": ResourceCategory{
			primaryCat: "networking",
		},
		"azurerm_route_filter": ResourceCategory{
			primaryCat: "networking",
		},
		"azurerm_virtual_machine_scale_set": ResourceCategory{
			primaryCat: "compute",
		},
		"azurerm_security_center_contact": ResourceCategory{
			primaryCat: "security",
		},
		"azurerm_security_center_subscription_pricing": ResourceCategory{
			primaryCat: "security",
		},
		"azurerm_storage_account": ResourceCategory{
			primaryCat: "storage",
		},
		"azurerm_storage_blob": ResourceCategory{
			primaryCat: "storage",
		},
		"azurerm_storage_container": ResourceCategory{
			primaryCat: "storage",
		},
		"azurerm_synapse_workspace": ResourceCategory{
			primaryCat: "analytics",
		},
		"azurerm_synapse_sql_pool": ResourceCategory{
			primaryCat: "analytics",
		},
		"azurerm_synapse_spark_pool": ResourceCategory{
			primaryCat: "analytics",
		},
		"azurerm_synapse_firewall_rule": ResourceCategory{
			primaryCat: "security",
		},
		"azurerm_synapse_managed_private_endpoint": ResourceCategory{
			primaryCat: "networking",
		},
		"azurerm_synapse_private_link_hub": ResourceCategory{
			primaryCat: "networking",
		},
		"azurerm_ssh_public_key": ResourceCategory{
			primaryCat: "security",
		},
		"azurerm_virtual_machine": ResourceCategory{
			primaryCat: "compute",
		},
		"azurerm_virtual_network": ResourceCategory{
			primaryCat: "networking",
		},
		"azurerm_subnet": ResourceCategory{
			primaryCat: "networking",
		},
		"azurerm_subnet_service_endpoint_storage_policy": ResourceCategory{
			primaryCat: "networking",
		},
		"azurerm_subnet_nat_gateway_association": ResourceCategory{
			primaryCat: "networking",
		},
		"azurerm_subnet_route_table_association": ResourceCategory{
			primaryCat: "networking",
		},
		"azurerm_subnet_network_security_group_association": ResourceCategory{
			primaryCat: "security",
		},
	}
}
