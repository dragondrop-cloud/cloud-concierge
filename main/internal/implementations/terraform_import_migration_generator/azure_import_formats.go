package terraformImportMigrationGenerator

var AzureResourceTypeLocations = map[ResourceType]ImportLocationFormat{
	"azurerm_analysis_services_server": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.AnalysisServices/servers/$2",
		Attributes:   []string{"subscription_id", "resource_group_name", "name"},
	},
	"azurerm_app_service": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.Web/sites/$2",
		Attributes:   []string{"subscription_id", "resource_group_name", "name"},
	},
	"azurerm_application_gateway": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.Network/applicationGateways/$2",
		Attributes:   []string{"subscription_id", "resource_group_name", "name"},
	},
	"azurerm_container_group": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.ContainerInstance/containerGroups/$2",
		Attributes:   []string{"subscription_id", "resource_group_name", "name"},
	},
	"azurerm_container_registry": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.ContainerRegistry/registries/$2",
		Attributes:   []string{"subscription_id", "resource_group_name", "name"},
	},
	"azurerm_container_registry_webhook": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.ContainerRegistry/registries/$2/webhooks/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "registry_name", "name"},
	},
	"azurerm_cosmosdb_account": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.DocumentDB/databaseAccounts/$2",
		Attributes:   []string{"subscription_id", "resource_group_name", "name"},
	},
	"azurerm_cosmosdb_sql_container": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.DocumentDB/databaseAccounts/$2/sqlDatabases/$3/containers/$4",
		Attributes:   []string{"subscription_id", "resource_group_name", "account_name", "database_name", "name"},
	},
	"azurerm_cosmosdb_sql_database": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.DocumentDB/databaseAccounts/$2/sqlDatabases/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "account_name", "name"},
	},
	"azurerm_cosmosdb_table": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.DocumentDB/databaseAccounts/$2/tables/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "account_name", "name"},
	},
	"azurerm_mariadb_configuration": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.DBforMariaDB/servers/$2/configurations/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "server_name", "name"},
	},
	"azurerm_mariadb_database": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.DBforMariaDB/servers/$2/databases/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "server_name", "name"},
	},
	"azurerm_mariadb_firewall_rule": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.DBforMariaDB/servers/$2/firewallRules/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "server_name", "name"},
	},
	"azurerm_mariadb_server": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.DBforMariaDB/servers/$2",
		Attributes:   []string{"subscription_id", "resource_group_name", "name"},
	},
	"azurerm_mariadb_virtual_network_rule": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.DBforMariaDB/servers/$2/virtualNetworkRules/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "server_name", "name"},
	},
	"azurerm_mysql_configuration": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.DBforMySQL/servers/$2/configurations/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "server_name", "name"},
	},
	"azurerm_mysql_database": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.DBforMySQL/servers/$2/databases/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "server_name", "name"},
	},
	"azurerm_mysql_firewall_rule": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.DBforMySQL/servers/$2/firewallRules/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "server_name", "name"},
	},
	"azurerm_mysql_server": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.DBforMySQL/servers/$2",
		Attributes:   []string{"subscription_id", "resource_group_name", "name"},
	},
	"azurerm_mysql_virtual_network_rule": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.DBforMySQL/servers/$2/virtualNetworkRules/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "server_name", "name"},
	},
	"azurerm_postgresql_configuration": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.DBforPostgreSQL/servers/$2/configurations/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "server_name", "name"},
	},
	"azurerm_postgresql_database": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.DBforPostgreSQL/servers/$2/databases/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "server_name", "name"},
	},
	"azurerm_postgresql_firewall_rule": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.DBforPostgreSQL/servers/$2/firewallRules/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "server_name", "name"},
	},
	"azurerm_postgresql_server": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.DBforPostgreSQL/servers/$2",
		Attributes:   []string{"subscription_id", "resource_group_name", "name"},
	},
	"azurerm_postgresql_virtual_network_rule": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.DBforPostgreSQL/servers/$2/virtualNetworkRules/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "server_name", "name"},
	},
	"azurerm_sql_database": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.Sql/servers/$2/databases/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "server_name", "name"},
	},
	"azurerm_sql_active_directory_administrator": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.Sql/servers/$2/administrators/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "server_name", "name"},
	},
	"azurerm_sql_elasticpool": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.Sql/servers/$2/elasticPools/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "server_name", "name"},
	},
	"azurerm_sql_failover_group": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.Sql/servers/$2/failoverGroups/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "server_name", "name"},
	},
	"azurerm_sql_firewall_rule": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.Sql/servers/$2/firewallRules/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "server_name", "name"},
	},
	"azurerm_sql_server": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.Sql/servers/$2",
		Attributes:   []string{"subscription_id", "resource_group_name", "name"},
	},
	"azurerm_sql_virtual_network_rule": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.Sql/servers/$2/virtualNetworkRules/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "server_name", "name"},
	},
	"azurerm_databricks_workspace": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.Databricks/workspaces/$2",
		Attributes:   []string{"subscription_id", "resource_group_name", "name"},
	},
	"azurerm_data_factory": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.DataFactory/factories/$2",
		Attributes:   []string{"subscription_id", "resource_group_name", "name"},
	},
	"azurerm_data_factory_pipeline": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.DataFactory/factories/$2/pipelines/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "factory_name", "name"},
	},
	"azurerm_data_factory_data_flow": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.DataFactory/factories/$2/dataflows/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "factory_name", "name"},
	},
	"azurerm_data_factory_dataset_azure_blob": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.DataFactory/factories/$2/datasets/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "factory_name", "name"},
	},
	"azurerm_data_factory_dataset_binary": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.DataFactory/factories/$2/datasets/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "factory_name", "name"},
	},
	"azurerm_data_factory_dataset_cosmosdb_sqlapi": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.DataFactory/factories/$2/datasets/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "factory_name", "name"},
	},
	"azurerm_data_factory_custom_dataset": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.DataFactory/factories/$2/datasets/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "factory_name", "name"},
	},
	"azurerm_data_factory_dataset_delimited_text": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.DataFactory/factories/$2/datasets/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "factory_name", "name"},
	},
	"azurerm_data_factory_dataset_http": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.DataFactory/factories/$2/datasets/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "factory_name", "name"},
	},
	"azurerm_data_factory_dataset_json": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.DataFactory/factories/$2/datasets/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "factory_name", "name"},
	},
	"azurerm_data_factory_dataset_mysql": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.DataFactory/factories/$2/datasets/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "factory_name", "name"},
	},
	"azurerm_data_factory_dataset_parquet": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.DataFactory/factories/$2/datasets/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "factory_name", "name"},
	},
	"azurerm_data_factory_dataset_postgresql": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.DataFactory/factories/$2/datasets/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "factory_name", "name"},
	},
	"azurerm_data_factory_dataset_snowflake": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.DataFactory/factories/$2/datasets/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "factory_name", "name"},
	},
	"azurerm_data_factory_dataset_sql_server_table": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.DataFactory/factories/$2/datasets/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "factory_name", "name"},
	},
	"azurerm_data_factory_integration_runtime_azure": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.DataFactory/factories/$2/integrationRuntimes/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "factory_name", "name"},
	},
	"azurerm_data_factory_integration_runtime_managed": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.DataFactory/factories/$2/integrationRuntimes/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "factory_name", "name"},
	},
	"azurerm_data_factory_integration_runtime_azure_ssis": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.DataFactory/factories/$2/integrationRuntimes/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "factory_name", "name"},
	},
	"azurerm_data_factory_integration_runtime_self_hosted": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.DataFactory/factories/$2/integrationRuntimes/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "factory_name", "name"},
	},
	"azurerm_data_factory_linked_service_azure_blob_storage": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.DataFactory/factories/$2/linkedservices/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "factory_name", "name"},
	},
	"azurerm_data_factory_linked_service_azure_databricks": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.DataFactory/factories/$2/linkedservices/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "factory_name", "name"},
	},
	"azurerm_data_factory_linked_service_azure_file_storage": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.DataFactory/factories/$2/linkedservices/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "factory_name", "name"},
	},
	"azurerm_data_factory_linked_service_azure_function": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.DataFactory/factories/$2/linkedservices/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "factory_name", "name"},
	},
	"azurerm_data_factory_linked_service_azure_search": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.DataFactory/factories/$2/linkedservices/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "factory_name", "name"},
	},
	"azurerm_data_factory_linked_service_azure_sql_database": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.DataFactory/factories/$2/linkedservices/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "factory_name", "name"},
	},
	"azurerm_data_factory_linked_service_azure_table_storage": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.DataFactory/factories/$2/linkedservices/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "factory_name", "name"},
	},
	"azurerm_data_factory_linked_service_cosmosdb": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.DataFactory/factories/$2/linkedservices/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "factory_name", "name"},
	},
	"azurerm_data_factory_linked_custom_service": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.DataFactory/factories/$2/linkedservices/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "factory_name", "name"},
	},
	"azurerm_data_factory_linked_service_data_lake_storage_gen2": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.DataFactory/factories/$2/linkedservices/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "factory_name", "name"},
	},
	"azurerm_data_factory_linked_service_key_vault": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.DataFactory/factories/$2/linkedservices/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "factory_name", "name"},
	},
	"azurerm_data_factory_linked_service_kusto": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.DataFactory/factories/$2/linkedservices/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "factory_name", "name"},
	},
	"azurerm_data_factory_linked_service_mysql": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.DataFactory/factories/$2/linkedservices/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "factory_name", "name"},
	},
	"azurerm_data_factory_linked_service_odata": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.DataFactory/factories/$2/linkedservices/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "factory_name", "name"},
	},
	"azurerm_data_factory_linked_service_postgresql": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.DataFactory/factories/$2/linkedservices/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "factory_name", "name"},
	},
	"azurerm_data_factory_linked_service_sftp": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.DataFactory/factories/$2/linkedservices/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "factory_name", "name"},
	},
	"azurerm_data_factory_linked_service_snowflake": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.DataFactory/factories/$2/linkedservices/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "factory_name", "name"},
	},
	"azurerm_data_factory_linked_service_sql_server": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.DataFactory/factories/$2/linkedservices/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "factory_name", "name"},
	},
	"azurerm_data_factory_linked_service_synapse": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.DataFactory/factories/$2/linkedservices/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "factory_name", "name"},
	},
	"azurerm_data_factory_linked_service_web": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.DataFactory/factories/$2/linkedservices/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "factory_name", "name"},
	},
	"azurerm_data_factory_trigger_blob_event": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.DataFactory/factories/$2/triggers/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "factory_name", "name"},
	},
	"azurerm_data_factory_trigger_schedule": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.DataFactory/factories/$2/triggers/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "factory_name", "name"},
	},
	"azurerm_data_factory_trigger_tumbling_window": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.DataFactory/factories/$2/triggers/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "factory_name", "name"},
	},
	"azurerm_managed_disk": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.Compute/disks/$2",
		Attributes:   []string{"subscription_id", "resource_group_name", "name"},
	},
	"azurerm_dns_a_record": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.Network/dnszones/$2/A/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "zone_name", "name"},
	},
	"azurerm_dns_aaaa_record": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.Network/dnszones/$2/AAAA/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "zone_name", "name"},
	},
	"azurerm_dns_caa_record": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.Network/dnszones/$2/CAA/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "zone_name", "name"},
	},
	"azurerm_dns_cname_record": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.Network/dnszones/$2/CNAME/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "zone_name", "name"},
	},
	"azurerm_dns_mx_record": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.Network/dnszones/$2/MX/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "zone_name", "name"},
	},
	"azurerm_dns_ns_record": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.Network/dnszones/$2/NS/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "zone_name", "name"},
	},
	"azurerm_dns_ptr_record": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.Network/dnszones/$2/PTR/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "zone_name", "name"},
	},
	"azurerm_dns_srv_record": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.Network/dnszones/$2/SRV/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "zone_name", "name"},
	},
	"azurerm_dns_txt_record": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.Network/dnszones/$2/TXT/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "zone_name", "name"},
	},
	"azurerm_dns_zone": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.Network/dnszones/$2",
		Attributes:   []string{"subscription_id", "resource_group_name", "name"},
	},
	"azurerm_lb": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.Network/loadBalancers/$2",
		Attributes:   []string{"subscription_id", "resource_group_name", "name"},
	},
	"azurerm_lb_backend_address_pool": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.Network/loadBalancers/$2/backendAddressPools/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "loadbalancer_name", "name"},
	},
	"azurerm_lb_nat_rule": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.Network/loadBalancers/$2/inboundNatRules/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "loadbalancer_name", "name"},
	},
	"azurerm_lb_probe": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.Network/loadBalancers/$2/probes/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "loadbalancer_name", "name"},
	},
	"azurerm_eventhub_namespace": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.EventHub/namespaces/$2",
		Attributes:   []string{"subscription_id", "resource_group_name", "name"},
	},
	"azurerm_eventhub": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.EventHub/namespaces/$2/eventhubs/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "namespace_name", "name"},
	},
	"azurerm_eventhub_consumer_group": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.EventHub/namespaces/$2/eventhubs/$3/consumergroups/$4",
		Attributes:   []string{"subscription_id", "resource_group_name", "namespace_name", "eventhub_name", "name"},
	},
	"azurerm_eventhub_namespace_authorization_rule": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.EventHub/namespaces/$2/AuthorizationRules/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "namespace_name", "name"},
	},
	"azurerm_network_interface": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.Network/networkInterfaces/$2",
		Attributes:   []string{"subscription_id", "resource_group_name", "name"},
	},
	"azurerm_network_security_group": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.Network/networkSecurityGroups/$2",
		Attributes:   []string{"subscription_id", "resource_group_name", "name"},
	},
	"azurerm_network_security_rule": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.Network/networkSecurityGroups/$2/securityRules/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "network_security_group_name", "name"},
	},
	"azurerm_network_watcher": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.Network/networkWatchers/$2",
		Attributes:   []string{"subscription_id", "resource_group_name", "name"},
	},
	"azurerm_network_watcher_flow_log": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.Network/networkWatchers/$2/flowLogs/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "network_watcher_name", "name"},
	},
	"azurerm_network_packet_capture": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.Network/networkWatchers/$2/packetCaptures/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "network_watcher_name", "name"},
	},
	"azurerm_private_dns_a_record": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.Network/privateDnsZones/$2/A/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "zone_name", "name"},
	},
	"azurerm_private_dns_aaaa_record": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.Network/privateDnsZones/$2/AAAA/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "zone_name", "name"},
	},
	"azurerm_private_dns_cname_record": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.Network/privateDnsZones/$2/CNAME/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "zone_name", "name"},
	},
	"azurerm_private_dns_mx_record": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.Network/privateDnsZones/$2/MX/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "zone_name", "name"},
	},
	"azurerm_private_dns_ptr_record": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.Network/privateDnsZones/$2/PTR/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "zone_name", "name"},
	},
	"azurerm_private_dns_srv_record": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.Network/privateDnsZones/$2/SRV/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "zone_name", "name"},
	},
	"azurerm_private_dns_txt_record": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.Network/privateDnsZones/$2/TXT/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "zone_name", "name"},
	},
	"azurerm_private_dns_zone": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.Network/privateDnsZones/$2",
		Attributes:   []string{"subscription_id", "resource_group_name", "name"},
	},
	"azurerm_private_dns_zone_virtual_network_link": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.Network/privateDnsZones/$2/virtualNetworkLinks/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "private_dns_zone_name", "name"},
	},
	"azurerm_private_endpoint": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.Network/privateEndpoints/$2",
		Attributes:   []string{"subscription_id", "resource_group_name", "name"},
	},
	"azurerm_private_link_service": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.Network/privateLinkServices/$2",
		Attributes:   []string{"subscription_id", "resource_group_name", "name"},
	},
	"azurerm_public_ip": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.Network/publicIPAddresses/$2",
		Attributes:   []string{"subscription_id", "resource_group_name", "name"},
	},
	"azurerm_public_ip_prefix": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.Network/publicIPPrefixes/$2",
		Attributes:   []string{"subscription_id", "resource_group_name", "name"},
	},
	"azurerm_redis_cache": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.Cache/Redis/$2",
		Attributes:   []string{"subscription_id", "resource_group_name", "name"},
	},
	"azurerm_purview_account": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.Purview/accounts/$2",
		Attributes:   []string{"subscription_id", "resource_group_name", "name"},
	},
	"azurerm_resource_group": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1",
		Attributes:   []string{"subscription_id", "name"},
	},
	"azurerm_management_lock": {
		StringFormat: "/subscriptions/$0/providers/Microsoft.Authorization/locks/$1",
		Attributes:   []string{"subscription_id", "name"},
	},
	"azurerm_route_table": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.Network/routeTables/$2",
		Attributes:   []string{"subscription_id", "resource_group_name", "name"},
	},
	"azurerm_route": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.Network/routeTables/$2/routes/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "route_table_name", "name"},
	},
	"azurerm_route_filter": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.Network/routeFilters/$2",
		Attributes:   []string{"subscription_id", "resource_group_name", "name"},
	},
	"azurerm_virtual_machine_scale_set": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.Compute/virtualMachineScaleSets/$2",
		Attributes:   []string{"subscription_id", "resource_group_name", "name"},
	},
	"azurerm_security_center_contact": {
		StringFormat: "/subscriptions/$0/providers/Microsoft.Security/securityContacts/default1",
		Attributes:   []string{"subscription_id"},
	},
	"azurerm_security_center_subscription_pricing": {
		StringFormat: "/subscriptions/$0/providers/Microsoft.Security/pricings/$1",
		Attributes:   []string{"subscription_id", "pricing_name"},
	},
	"azurerm_storage_account": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.Storage/storageAccounts/$2",
		Attributes:   []string{"subscription_id", "resource_group_name", "name"},
	},
	"azurerm_storage_blob": {
		StringFormat: "https://$0.blob.core.windows.net/$1/$2",
		Attributes:   []string{"storage_account_name", "container_name", "name"},
	},
	"azurerm_storage_container": {
		StringFormat: "https://$0.blob.core.windows.net/$1",
		Attributes:   []string{"storage_account_name", "name"},
	},
	"azurerm_synapse_workspace": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.Synapse/workspaces/$2",
		Attributes:   []string{"subscription_id", "resource_group_name", "name"},
	},
	"azurerm_synapse_sql_pool": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.Synapse/workspaces/$2/sqlPools/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "workspace_name", "name"},
	},
	"azurerm_synapse_spark_pool": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.Synapse/workspaces/$2/bigDataPools/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "workspace_name", "name"},
	},
	"azurerm_synapse_firewall_rule": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.Synapse/workspaces/$2/firewallRules/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "workspace_name", "name"},
	},
	"azurerm_synapse_managed_private_endpoint": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.Synapse/workspaces/$2/managedPrivateEndpoints/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "workspace_name", "name"},
	},
	"azurerm_synapse_private_link_hub": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.Synapse/privateLinkHubs/$2",
		Attributes:   []string{"subscription_id", "resource_group_name", "name"},
	},
	"azurerm_ssh_public_key": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.Compute/sshPublicKeys/$2",
		Attributes:   []string{"subscription_id", "resource_group_name", "name"},
	},
	"azurerm_virtual_machine": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.Compute/virtualMachines/$2",
		Attributes:   []string{"subscription_id", "resource_group_name", "name"},
	},
	"azurerm_virtual_network": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.Network/virtualNetworks/$2",
		Attributes:   []string{"subscription_id", "resource_group_name", "name"},
	},
	"azurerm_subnet": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.Network/virtualNetworks/$2/subnets/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "virtual_network_name", "name"},
	},
	"azurerm_subnet_service_endpoint_storage_policy": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.Network/serviceEndpointPolicies/$2",
		Attributes:   []string{"subscription_id", "resource_group_name", "name"},
	},
	"azurerm_subnet_nat_gateway_association": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.Network/virtualNetworks/$2/subnets/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "virtual_network_name", "subnet_name"},
	},
	"azurerm_subnet_route_table_association": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.Network/virtualNetworks/$2/subnets/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "virtual_network_name", "subnet_name"},
	},
	"azurerm_subnet_network_security_group_association": {
		StringFormat: "/subscriptions/$0/resourceGroups/$1/providers/Microsoft.Network/virtualNetworks/$2/subnets/$3",
		Attributes:   []string{"subscription_id", "resource_group_name", "virtual_network_name", "subnet_name"},
	},
}
