package terraformImportMigrationGenerator

var ResourceTypeLocations = map[ResourceType]ImportLocationFormat{
	"aws_s3_bucket": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_accessanalyzer_analyzer": {
		StringFormat: "$0",
		Attributes:   []string{"arn"},
	},
	"aws_acm_certificate": {
		StringFormat: "$0",
		Attributes:   []string{"arn"},
	},
	"aws_lb": {
		StringFormat: "$0",
		Attributes:   []string{"arn"},
	},
	"aws_lb_listener": {
		StringFormat: "$0/$1",
		Attributes:   []string{"load_balancer_arn", "port"},
	},
	"aws_lb_listener_rule": {
		StringFormat: "$0/$1/$2",
		Attributes:   []string{"load_balancer_arn", "listener_arn", "arn"},
	},

	"aws_lb_listener_certificate": {
		StringFormat: "$0/$1/$2",
		Attributes:   []string{"load_balancer_arn", "listener_arn", "certificate_arn"},
	},
	"aws_lb_target_group": {
		StringFormat: "$0",
		Attributes:   []string{"arn"},
	},
	"aws_lb_target_group_attachment": {
		StringFormat: "$0/$1",
		Attributes:   []string{"target_group_arn", "target_id"},
	},
	"aws_api_gateway_authorizer": {
		StringFormat: "$0/$1",
		Attributes:   []string{"rest_api_id", "id"},
	},
	"aws_api_gateway_api_key": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_api_gateway_documentation_part": {
		StringFormat: "$0/$1",
		Attributes:   []string{"rest_api_id", "id"},
	},
	"aws_api_gateway_gateway_response": {
		StringFormat: "$0/$1",
		Attributes:   []string{"rest_api_id", "response_type"},
	},
	"aws_api_gateway_integration": {
		StringFormat: "$0/$1",
		Attributes:   []string{"rest_api_id", "resource_id"},
	},
	"aws_api_gateway_integration_response": {
		StringFormat: "$0/$1/$2",
		Attributes:   []string{"rest_api_id", "resource_id", "http_method"},
	},
	"aws_api_gateway_method": {
		StringFormat: "$0/$1/$2",
		Attributes:   []string{"rest_api_id", "resource_id", "http_method"},
	},
	"aws_api_gateway_method_response": {
		StringFormat: "$0/$1/$2/$3",
		Attributes:   []string{"rest_api_id", "resource_id", "http_method", "status_code"},
	},
	"aws_api_gateway_model": {
		StringFormat: "$0/$1",
		Attributes:   []string{"rest_api_id", "name"},
	},
	"aws_api_gateway_resource": {
		StringFormat: "$0/$1",
		Attributes:   []string{"rest_api_id", "path"},
	},
	"aws_api_gateway_rest_api": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_api_gateway_stage": {
		StringFormat: "$0/$1",
		Attributes:   []string{"rest_api_id", "stage_name"},
	},
	"aws_api_gateway_usage_plan": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_api_gateway_vpc_link": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_appsync_graphql_api": {
		StringFormat: "$0",
		Attributes:   []string{"arn"},
	},
	"aws_autoscaling_group": {
		StringFormat: "$0",
		Attributes:   []string{"name"},
	},
	"aws_launch_configuration": {
		StringFormat: "$0",
		Attributes:   []string{"name"},
	},
	"aws_launch_template": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_batch_compute_environment": {
		StringFormat: "$0",
		Attributes:   []string{"arn"},
	},
	"aws_batch_job_definition": {
		StringFormat: "$0",
		Attributes:   []string{"arn"},
	},
	"aws_batch_job_queue": {
		StringFormat: "$0",
		Attributes:   []string{"arn"},
	},
	"aws_budgets_budget": {
		StringFormat: "$0",
		Attributes:   []string{"name"},
	},
	"aws_cloud9_environment_ec2": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_cloudformation_stack": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_cloudformation_stack_set": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_cloudformation_stack_set_instance": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_cloudfront_distribution": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_cloudfront_cache_policy": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_cloudhsm_v2_cluster": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_cloudhsm_v2_hsm": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_cloudtrail": {
		StringFormat: "$0",
		Attributes:   []string{"name"},
	},
	"aws_cloudwatch_dashboard": {
		StringFormat: "$0",
		Attributes:   []string{"dashboard_name"},
	},
	"aws_cloudwatch_event_rule": {
		StringFormat: "$0",
		Attributes:   []string{"name"},
	},
	"aws_cloudwatch_event_target": {
		StringFormat: "$0",
		Attributes:   []string{"rule"},
	},
	"aws_cloudwatch_metric_alarm": {
		StringFormat: "$0",
		Attributes:   []string{"alarm_name"},
	},
	"aws_codebuild_project": {
		StringFormat: "$0",
		Attributes:   []string{"name"},
	},
	"aws_codecommit_repository": {
		StringFormat: "$0",
		Attributes:   []string{"repository_name"},
	},
	"aws_codedeploy_app": {
		StringFormat: "$0",
		Attributes:   []string{"name"},
	},
	"aws_codepipeline": {
		StringFormat: "$0",
		Attributes:   []string{"name"},
	},
	"aws_codepipeline_webhook": {
		StringFormat: "$0",
		Attributes:   []string{"name"},
	},
	"aws_cognito_identity_pool": {
		StringFormat: "$0",
		Attributes:   []string{"identity_pool_name"},
	},
	"aws_cognito_user_pool": {
		StringFormat: "$0",
		Attributes:   []string{"name"},
	},
	"aws_config_config_rule": {
		StringFormat: "$0",
		Attributes:   []string{"name"},
	},
	"aws_config_configuration_recorder": {
		StringFormat: "$0",
		Attributes:   []string{"name"},
	},
	"aws_config_delivery_channel": {
		StringFormat: "$0",
		Attributes:   []string{"name"},
	},
	"aws_customer_gateway": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_datapipeline_pipeline": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_devicefarm_project": {
		StringFormat: "$0",
		Attributes:   []string{"name"},
	},
	"aws_docdb_cluster": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_docdb_cluster_instance": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_docdb_cluster_parameter_group": {
		StringFormat: "$0",
		Attributes:   []string{"name"},
	},
	"aws_docdb_subnet_group": {
		StringFormat: "$0",
		Attributes:   []string{"name"},
	},
	"aws_dynamodb_table": {
		StringFormat: "$0",
		Attributes:   []string{"name"},
	},
	"aws_ebs_volume": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_volume_attachment": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_instance": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_ecr_lifecycle_policy": {
		StringFormat: "$0",
		Attributes:   []string{"repository"},
	},
	"aws_ecr_repository": {
		StringFormat: "$0",
		Attributes:   []string{"name"},
	},
	"aws_ecr_repository_policy": {
		StringFormat: "$0",
		Attributes:   []string{"repository"},
	},
	"aws_ecrpublic_repository": {
		StringFormat: "$0",
		Attributes:   []string{"name"},
	},
	"aws_ecs_cluster": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_ecs_service": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_ecs_task_definition": {
		StringFormat: "$0",
		Attributes:   []string{"arn"},
	},
	"aws_efs_access_point": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_efs_file_system": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_efs_file_system_policy": {
		StringFormat: "$0",
		Attributes:   []string{"file_system_id"},
	},
	"aws_efs_mount_target": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_eip": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_eks_cluster": {
		StringFormat: "$0",
		Attributes:   []string{"name"},
	},
	"aws_eks_node_group": {
		StringFormat: "$0:$1",
		Attributes:   []string{"cluster_name", "node_group_name"},
	},
	"aws_elasticache_cluster": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_elasticache_parameter_group": {
		StringFormat: "$0",
		Attributes:   []string{"name"},
	},
	"aws_elasticache_subnet_group": {
		StringFormat: "$0",
		Attributes:   []string{"name"},
	},
	"aws_elasticache_replication_group": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_elastic_beanstalk_application": {
		StringFormat: "$0",
		Attributes:   []string{"name"},
	},
	"aws_elastic_beanstalk_environment": {
		StringFormat: "$0",
		Attributes:   []string{"name"},
	},
	"aws_elb": {
		StringFormat: "$0",
		Attributes:   []string{"name"},
	},
	"aws_emr_cluster": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_emr_security_configuration": {
		StringFormat: "$0",
		Attributes:   []string{"name"},
	},
	"aws_network_interface": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_elasticsearch_domain": {
		StringFormat: "$0",
		Attributes:   []string{"domain_name"},
	},
	"aws_kinesis_firehose_delivery_stream": {
		StringFormat: "$0",
		Attributes:   []string{"name"},
	},
	"aws_glue_crawler": {
		StringFormat: "$0",
		Attributes:   []string{"name"},
	},
	"aws_glue_catalog_database": {
		StringFormat: "$0",
		Attributes:   []string{"name"},
	},
	"aws_glue_catalog_table": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_glue_job": {
		StringFormat: "$0",
		Attributes:   []string{"name"},
	},
	"aws_glue_trigger": {
		StringFormat: "$0",
		Attributes:   []string{"name"},
	},
	"aws_iam_access_key": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_iam_group": {
		StringFormat: "$0",
		Attributes:   []string{"name"},
	},
	"aws_iam_group_policy": {
		StringFormat: "$0:$1",
		Attributes:   []string{"group", "name"},
	},
	"aws_iam_group_policy_attachment": {
		StringFormat: "$0/$1",
		Attributes:   []string{"group", "policy_arn"},
	},
	"aws_iam_instance_profile": {
		StringFormat: "$0",
		Attributes:   []string{"name"},
	},
	"aws_iam_policy": {
		StringFormat: "$0",
		Attributes:   []string{"arn"},
	},
	"aws_iam_role": {
		StringFormat: "$0",
		Attributes:   []string{"name"},
	},
	"aws_iam_role_policy": {
		StringFormat: "$0/$1",
		Attributes:   []string{"role", "name"},
	},
	"aws_iam_role_policy_attachment": {
		StringFormat: "$0/$1",
		Attributes:   []string{"role", "policy_arn"},
	},
	"aws_iam_user": {
		StringFormat: "$0",
		Attributes:   []string{"name"},
	},
	"aws_iam_user_group_membership": {
		StringFormat: "$0/$1",
		Attributes:   []string{"user", "group"},
	},
	"aws_iam_user_policy": {
		StringFormat: "$0:$1",
		Attributes:   []string{"user", "name"},
	},
	"aws_iam_user_policy_attachment": {
		StringFormat: "$0/$1",
		Attributes:   []string{"user", "policy_arn"},
	},
	"aws_internet_gateway": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_iot_thing": {
		StringFormat: "$0",
		Attributes:   []string{"name"},
	},
	"aws_iot_thing_type": {
		StringFormat: "$0",
		Attributes:   []string{"name"},
	},
	"aws_iot_topic_rule": {
		StringFormat: "$0",
		Attributes:   []string{"name"},
	},
	"aws_iot_role_alias": {
		StringFormat: "$0",
		Attributes:   []string{"name"},
	},
	"aws_kinesis_stream": {
		StringFormat: "$0",
		Attributes:   []string{"name"},
	},
	"aws_kms_key": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_kms_alias": {
		StringFormat: "$0",
		Attributes:   []string{"name"},
	},
	"aws_kms_grant": {
		StringFormat: "$0:$1",
		Attributes:   []string{"key_id", "grant_id"},
	},
	"aws_lambda_event_source_mapping": {
		StringFormat: "$0",
		Attributes:   []string{"uuid"},
	},
	"aws_lambda_function": {
		StringFormat: "$0",
		Attributes:   []string{"function_name"},
	},
	"aws_lambda_function_event_invoke_config": {
		StringFormat: "$0",
		Attributes:   []string{"function_name"},
	},
	"aws_lambda_layer_version": {
		StringFormat: "$0",
		Attributes:   []string{"arn"},
	},
	"aws_lambda_permission": {
		StringFormat: "$0",
		Attributes:   []string{"statement_id"},
	},
	"aws_cloudwatch_log_group": {
		StringFormat: "$0",
		Attributes:   []string{"name"},
	},
	"aws_media_package_channel": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_media_store_container": {
		StringFormat: "$0",
		Attributes:   []string{"name"},
	},
	"aws_medialive_channel": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_medialive_input": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_medialive_input_security_group": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_msk_cluster": {
		StringFormat: "$0",
		Attributes:   []string{"cluster_name"},
	},
	"aws_network_acl": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_nat_gateway": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_opsworks_application": {
		StringFormat: "$0",
		Attributes:   []string{"app_id"},
	},
	"aws_opsworks_custom_layer": {
		StringFormat: "$0",
		Attributes:   []string{"layer_id"},
	},
	"aws_opsworks_instance": {
		StringFormat: "$0",
		Attributes:   []string{"instance_id"},
	},
	"aws_opsworks_java_app_layer": {
		StringFormat: "$0",
		Attributes:   []string{"layer_id"},
	},
	"aws_opsworks_php_app_layer": {
		StringFormat: "$0",
		Attributes:   []string{"layer_id"},
	},
	"aws_opsworks_rds_db_instance": {
		StringFormat: "$0",
		Attributes:   []string{"rds_db_instance_arn"},
	},
	"aws_opsworks_stack": {
		StringFormat: "$0",
		Attributes:   []string{"stack_id"},
	},
	"aws_opsworks_static_web_layer": {
		StringFormat: "$0",
		Attributes:   []string{"layer_id"},
	},
	"aws_opsworks_user_profile": {
		StringFormat: "$0",
		Attributes:   []string{"user_arn"},
	},
	"aws_organizations_account": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_organizations_organization": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_organizations_organizational_unit": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_organizations_policy": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_organizations_policy_attachment": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_qldb_ledger": {
		StringFormat: "$0",
		Attributes:   []string{"name"},
	},
	"aws_db_instance": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_db_proxy": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_db_cluster": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_db_cluster_snapshot": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_db_parameter_group": {
		StringFormat: "$0",
		Attributes:   []string{"name"},
	},
	"aws_db_snapshot": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_db_subnet_group": {
		StringFormat: "$0",
		Attributes:   []string{"name"},
	},
	"aws_db_option_group": {
		StringFormat: "$0",
		Attributes:   []string{"name"},
	},
	"aws_db_event_subscription": {
		StringFormat: "$0",
		Attributes:   []string{"name"},
	},
	"aws_rds_global_cluster": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_resourcegroups_group": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_route53_zone": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_route53_record": {
		StringFormat: "$0/$1",
		Attributes:   []string{"zone_id", "name"},
	},
	"aws_route_table": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_main_route_table_association": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_route_table_association": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_secretsmanager_secret": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_securityhub_account": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_securityhub_member": {
		StringFormat: "$0/$1",
		Attributes:   []string{"account_id", "id"},
	},
	"aws_securityhub_standards_subscription": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_servicecatalog_portfolio": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_ses_configuration_set": {
		StringFormat: "$0",
		Attributes:   []string{"name"},
	},
	"aws_ses_domain_identity": {
		StringFormat: "$0",
		Attributes:   []string{"domain"},
	},
	"aws_ses_email_identity": {
		StringFormat: "$0",
		Attributes:   []string{"email"},
	},
	"aws_ses_receipt_rule": {
		StringFormat: "$0/$1",
		Attributes:   []string{"rule_set_name", "name"},
	},
	"aws_ses_receipt_rule_set": {
		StringFormat: "$0",
		Attributes:   []string{"arn"},
	},
	"aws_ses_template": {
		StringFormat: "$0",
		Attributes:   []string{"arn"},
	},
	"aws_sfn_activity": {
		StringFormat: "$0",
		Attributes:   []string{"arn"},
	},
	"aws_sfn_state_machine": {
		StringFormat: "$0",
		Attributes:   []string{"arn"},
	},
	"aws_security_group": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_security_group_rule": {
		StringFormat: "$0/$1",
		Attributes:   []string{"security_group_id", "id"},
	},
	"aws_sns_topic": {
		StringFormat: "$0",
		Attributes:   []string{"arn"},
	},
	"aws_sns_topic_subscription": {
		StringFormat: "$0/$1",
		Attributes:   []string{"topic_arn", "subscription_arn"},
	},
	"aws_sqs_queue": {
		StringFormat: "$0",
		Attributes:   []string{"arn"},
	},
	"aws_ssm_parameter": {
		StringFormat: "$0",
		Attributes:   []string{"name"},
	},
	"aws_subnet": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_swf_domain": {
		StringFormat: "$0",
		Attributes:   []string{"name"},
	},
	"aws_ec2_transit_gateway_route_table": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_ec2_transit_gateway_vpc_attachment": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_vpc": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_vpc_peering_connection": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_vpn_connection": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_vpn_gateway": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_waf_byte_match_set": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_waf_geo_match_set": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_waf_ipset": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_waf_rate_based_rule": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_waf_regex_match_set": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_waf_regex_pattern_set": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_waf_rule": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_waf_rule_group": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_waf_size_constraint_set": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_waf_sql_injection_match_set": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_waf_web_acl": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_waf_xss_match_set": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_wafregional_byte_match_set": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_wafregional_geo_match_set": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_wafregional_ipset": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_wafregional_rate_based_rule": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_wafregional_regex_match_set": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_wafregional_regex_pattern_set": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_wafregional_rule": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_wafregional_rule_group": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_wafregional_size_constraint_set": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_wafregional_sql_injection_match_set": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_wafregional_web_acl": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_wafregional_xss_match_set": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_wafv2_ip_set": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_wafv2_regex_pattern_set": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_wafv2_rule_group": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_wafv2_web_acl": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
	"aws_wafv2_web_acl_association": {
		StringFormat: "$0/$1/$2",
		Attributes:   []string{"web_acl_arn", "resource_type", "resource_id"},
	},
	"aws_wafv2_web_acl_logging_configuration": {
		StringFormat: "$0/$1",
		Attributes:   []string{"resource_arn", "id"},
	},
	"aws_workspaces_directory": {
		StringFormat: "$0",
		Attributes:   []string{"directory_id"},
	},
	"aws_workspaces_ip_group": {
		StringFormat: "$0",
		Attributes:   []string{"group_id"},
	},
	"aws_workspaces_workspace": {
		StringFormat: "$0",
		Attributes:   []string{"workspace_id"},
	},
	"aws_xray_sampling_rule": {
		StringFormat: "$0",
		Attributes:   []string{"id"},
	},
}
