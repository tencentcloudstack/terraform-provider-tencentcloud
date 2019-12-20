/*
The TencentCloud provider is used to interact with many resources supported by [TencentCloud](https://intl.cloud.tencent.com).
The provider needs to be configured with the proper credentials before it can be used.

Use the navigation on the left to read about the available resources.

-> **Note:** From version 1.9.0 (June 18, 2019), the provider start to support Terraform 0.12.x.

Example Usage

```hcl
# Configure the TencentCloud Provider
provider "tencentcloud" {
  secret_id  = "${var.secret_id}"
  secret_key = "${var.secret_key}"
  region     = "${var.region}"
}
```

Resources List

Data Sources
  tencentcloud_as_scaling_configs
  tencentcloud_as_scaling_groups
  tencentcloud_as_scaling_policies
  tencentcloud_availability_zones
  tencentcloud_cam_group_memberships
  tencentcloud_cam_group_policy_attachments
  tencentcloud_cam_groups
  tencentcloud_cam_policies
  tencentcloud_cam_role_policy_attachments
  tencentcloud_cam_roles
  tencentcloud_cam_saml_providers
  tencentcloud_cam_user_policy_attachments
  tencentcloud_cam_users
  tencentcloud_cbs_snapshots
  tencentcloud_cbs_snapshot_policies
  tencentcloud_cbs_storages
  tencentcloud_ccn_bandwidth_limits
  tencentcloud_ccn_instances
  tencentcloud_cfs_access_groups
  tencentcloud_cfs_access_rules
  tencentcloud_cfs_file_systems
  tencentcloud_clb_attachments
  tencentcloud_clb_instances
  tencentcloud_clb_listener_rules
  tencentcloud_clb_listeners
  tencentcloud_clb_redirections
  tencentcloud_container_cluster_instances
  tencentcloud_container_clusters
  tencentcloud_cos_bucket_object
  tencentcloud_cos_buckets
  tencentcloud_dc_gateway_ccn_routes
  tencentcloud_dc_gateway_instances
  tencentcloud_dc_instances
  tencentcloud_dcx_instances
  tencentcloud_dnats
  tencentcloud_eip
  tencentcloud_eips
  tencentcloud_enis
  tencentcloud_gaap_certificates
  tencentcloud_gaap_http_domains
  tencentcloud_gaap_http_rules
  tencentcloud_gaap_layer4_listeners
  tencentcloud_gaap_layer7_listeners
  tencentcloud_gaap_proxies
  tencentcloud_gaap_realservers
  tencentcloud_gaap_security_policies
  tencentcloud_gaap_security_rules
  tencentcloud_ha_vip_eip_attachments
  tencentcloud_ha_vips
  tencentcloud_image
  tencentcloud_images
  tencentcloud_instance_types
  tencentcloud_instances
  tencentcloud_key_pairs
  tencentcloud_kubernetes_clusters
  tencentcloud_mongodb_instances
  tencentcloud_mongodb_zone_config
  tencentcloud_mysql_backup_list
  tencentcloud_mysql_instance
  tencentcloud_mysql_parameter_list
  tencentcloud_mysql_zone_config
  tencentcloud_nat_gateways
  tencentcloud_nats
  tencentcloud_placement_groups
  tencentcloud_redis_instances
  tencentcloud_redis_zone_config
  tencentcloud_reserved_instance_configs
  tencentcloud_reserved_instances
  tencentcloud_route_table
  tencentcloud_scf_functions
  tencentcloud_scf_logs
  tencentcloud_scf_namespaces
  tencentcloud_security_group
  tencentcloud_security_groups
  tencentcloud_ssl_certificates
  tencentcloud_subnet
  tencentcloud_tcaplus_applications
  tencentcloud_tcaplus_idls
  tencentcloud_tcaplus_tables
  tencentcloud_tcaplus_zones
  tencentcloud_vpc
  tencentcloud_vpc_instances
  tencentcloud_vpc_route_tables
  tencentcloud_vpc_subnets
  tencentcloud_vpn_connections
  tencentcloud_vpn_customer_gateways
  tencentcloud_vpn_gateways

AS Resources
  tencentcloud_as_scaling_config
  tencentcloud_as_scaling_group
  tencentcloud_as_attachment
  tencentcloud_as_scaling_policy
  tencentcloud_as_schedule
  tencentcloud_as_lifecycle_hook
  tencentcloud_as_notification

CAM Resources
  tencentcloud_cam_role
  tencentcloud_cam_role_policy_attachment
  tencentcloud_cam_policy
  tencentcloud_cam_user
  tencentcloud_cam_user_policy_attachment
  tencentcloud_cam_group
  tencentcloud_cam_group_policy_attachment
  tencentcloud_cam_group_membership
  tencentcloud_cam_saml_provider

CBS Resources
  tencentcloud_cbs_storage
  tencentcloud_cbs_storage_attachment
  tencentcloud_cbs_snapshot
  tencentcloud_cbs_snapshot_policy
  tencentcloud_cbs_snapshot_policy_attachment

CCN Resources
  tencentcloud_ccn
  tencentcloud_ccn_attachment
  tencentcloud_ccn_bandwidth_limit

CFS Resources
  tencentcloud_cfs_file_system
  tencentcloud_cfs_access_group
  tencentcloud_cfs_access_rule

Container Cluster Resources
  tencentcloud_container_cluster
  tencentcloud_container_cluster_instance

CLB Resources
  tencentcloud_clb_instance
  tencentcloud_clb_listener
  tencentcloud_clb_listener_rule
  tencentcloud_clb_attachment
  tencentcloud_clb_redirection
  tencentcloud_lb
  tencentcloud_alb_server_attachment

COS Resources
  tencentcloud_cos_bucket
  tencentcloud_cos_bucket_object

CVM Resources
  tencentcloud_instance
  tencentcloud_eip
  tencentcloud_eip_association
  tencentcloud_key_pair
  tencentcloud_placement_group
  tencentcloud_reserved_instance

DC Resources
  tencentcloud_dcx

DCG Resources
  tencentcloud_dc_gateway
  tencentcloud_dc_gateway_ccn_route

GAAP Resources
  tencentcloud_gaap_proxy
  tencentcloud_gaap_realserver
  tencentcloud_gaap_layer4_listener
  tencentcloud_gaap_layer7_listener
  tencentcloud_gaap_http_domain
  tencentcloud_gaap_http_rule
  tencentcloud_gaap_certificate
  tencentcloud_gaap_security_policy
  tencentcloud_gaap_security_rule

Kubernetes Resources
  tencentcloud_kubernetes_cluster
  tencentcloud_kubernetes_scale_worker
  tencentcloud_kubernetes_as_scaling_group

MongoDB Resources
  tencentcloud_mongodb_instance
  tencentcloud_mongodb_sharding_instance

MySQL Resources
  tencentcloud_mysql_instance
  tencentcloud_mysql_readonly_instance
  tencentcloud_mysql_account
  tencentcloud_mysql_privilege
  tencentcloud_mysql_account_privilege
  tencentcloud_mysql_backup_policy

Redis Resources
  tencentcloud_redis_instance
  tencentcloud_redis_backup_config

SCF Resources
  tencentcloud_scf_function
  tencentcloud_scf_namespace

SSL Resources
  tencentcloud_ssl_certificate

Tcaplus Resources
  tencentcloud_tcaplus_application
  tencentcloud_tcaplus_zone
  tencentcloud_tcaplus_idl
  tencentcloud_tcaplus_table

VPC Resources
  tencentcloud_eni
  tencentcloud_eni_attachment
  tencentcloud_vpc
  tencentcloud_subnet
  tencentcloud_security_group
  tencentcloud_security_group_rule
  tencentcloud_security_group_lite_rule
  tencentcloud_route_table
  tencentcloud_route_entry
  tencentcloud_route_table_entry
  tencentcloud_dnat
  tencentcloud_nat_gateway
  tencentcloud_ha_vip
  tencentcloud_ha_vip_eip_attachment

VPN Resources
  tencentcloud_vpn_customer_gateway
  tencentcloud_vpn_gateway
  tencentcloud_vpn_connection
*/
package tencentcloud

import (
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

const (
	PROVIDER_SECRET_ID      = "TENCENTCLOUD_SECRET_ID"
	PROVIDER_SECRET_KEY     = "TENCENTCLOUD_SECRET_KEY"
	PROVIDER_SECURITY_TOKEN = "TENCENTCLOUD_SECURITY_TOKEN"
	PROVIDER_REGION         = "TENCENTCLOUD_REGION"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"secret_id": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc(PROVIDER_SECRET_ID, nil),
				Description: "This is the TencentCloud access key. It must be provided, but it can also be sourced from the `TENCENTCLOUD_SECRET_ID` environment variable.",
			},
			"secret_key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc(PROVIDER_SECRET_KEY, nil),
				Description: "This is the TencentCloud secret key. It must be provided, but it can also be sourced from the `TENCENTCLOUD_SECRET_KEY` environment variable.",
				Sensitive:   true,
			},
			"security_token": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(PROVIDER_SECURITY_TOKEN, nil),
				Description: "TencentCloud Security Token of temporary access credentials. It can be sourced from the `TENCENTCLOUD_SECURITY_TOKEN` environment variable. Notice: for supported products, please refer to: [temporary key supported products](https://intl.cloud.tencent.com/document/product/598/10588).",
				Sensitive:   true,
			},
			"region": {
				Type:         schema.TypeString,
				Required:     true,
				DefaultFunc:  schema.EnvDefaultFunc(PROVIDER_REGION, nil),
				Description:  "This is the TencentCloud region. It must be provided, but it can also be sourced from the `TENCENTCLOUD_REGION` environment variables. The default input value is ap-guangzhou.",
				InputDefault: "ap-guangzhou",
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"tencentcloud_availability_zones":           dataSourceTencentCloudAvailabilityZones(),
			"tencentcloud_eip":                          dataSourceTencentCloudEip(),
			"tencentcloud_image":                        dataSourceTencentCloudImage(),
			"tencentcloud_images":                       dataSourceTencentCloudImages(),
			"tencentcloud_instance_types":               dataSourceInstanceTypes(),
			"tencentcloud_vpc":                          dataSourceTencentCloudVpc(),
			"tencentcloud_subnet":                       dataSourceTencentCloudSubnet(),
			"tencentcloud_route_table":                  dataSourceTencentCloudRouteTable(),
			"tencentcloud_security_group":               dataSourceTencentCloudSecurityGroup(),
			"tencentcloud_security_groups":              dataSourceTencentCloudSecurityGroups(),
			"tencentcloud_nats":                         dataSourceTencentCloudNats(),
			"tencentcloud_dnats":                        dataSourceTencentCloudDnats(),
			"tencentcloud_nat_gateways":                 dataSourceTencentCloudNatGateways(),
			"tencentcloud_container_clusters":           dataSourceTencentCloudContainerClusters(),
			"tencentcloud_container_cluster_instances":  dataSourceTencentCloudContainerClusterInstances(),
			"tencentcloud_mysql_backup_list":            dataSourceTencentMysqlBackupList(),
			"tencentcloud_mysql_zone_config":            dataSourceTencentMysqlZoneConfig(),
			"tencentcloud_mysql_parameter_list":         dataSourceTencentCloudMysqlParameterList(),
			"tencentcloud_mysql_instance":               dataSourceTencentCloudMysqlInstance(),
			"tencentcloud_cos_bucket_object":            dataSourceTencentCloudCosBucketObject(),
			"tencentcloud_cos_buckets":                  dataSourceTencentCloudCosBuckets(),
			"tencentcloud_redis_zone_config":            dataSourceTencentRedisZoneConfig(),
			"tencentcloud_redis_instances":              dataSourceTencentRedisInstances(),
			"tencentcloud_as_scaling_configs":           dataSourceTencentCloudAsScalingConfigs(),
			"tencentcloud_as_scaling_groups":            dataSourceTencentCloudAsScalingGroups(),
			"tencentcloud_as_scaling_policies":          dataSourceTencentCloudAsScalingPolicies(),
			"tencentcloud_vpc_instances":                dataSourceTencentCloudVpcInstances(),
			"tencentcloud_vpc_subnets":                  dataSourceTencentCloudVpcSubnets(),
			"tencentcloud_vpc_route_tables":             dataSourceTencentCloudVpcRouteTables(),
			"tencentcloud_ccn_instances":                dataSourceTencentCloudCcnInstances(),
			"tencentcloud_ccn_bandwidth_limits":         dataSourceTencentCloudCcnBandwidthLimits(),
			"tencentcloud_cbs_storages":                 dataSourceTencentCloudCbsStorages(),
			"tencentcloud_cbs_snapshots":                dataSourceTencentCloudCbsSnapshots(),
			"tencentcloud_cbs_snapshot_policies":        dataSourceTencentCloudCbsSnapshotPolicies(),
			"tencentcloud_dc_instances":                 dataSourceTencentCloudDcInstances(),
			"tencentcloud_clb_instances":                dataSourceTencentCloudClbInstances(),
			"tencentcloud_clb_listeners":                dataSourceTencentCloudClbListeners(),
			"tencentcloud_clb_listener_rules":           dataSourceTencentCloudClbListenerRules(),
			"tencentcloud_clb_attachments":              dataSourceTencentCloudClbServerAttachments(),
			"tencentcloud_clb_redirections":             dataSourceTencentCloudClbRedirections(),
			"tencentcloud_dcx_instances":                dataSourceTencentCloudDcxInstances(),
			"tencentcloud_mongodb_zone_config":          dataSourceTencentCloudMongodbZoneConfig(),
			"tencentcloud_mongodb_instances":            dataSourceTencentCloudMongodbInstances(),
			"tencentcloud_dc_gateway_instances":         dataSourceTencentCloudDcGatewayInstances(),
			"tencentcloud_dc_gateway_ccn_routes":        dataSourceTencentCloudDcGatewayCCNRoutes(),
			"tencentcloud_kubernetes_clusters":          dataSourceTencentCloudKubernetesClusters(),
			"tencentcloud_gaap_proxies":                 dataSourceTencentCloudGaapProxies(),
			"tencentcloud_gaap_realservers":             dataSourceTencentCloudGaapRealservers(),
			"tencentcloud_gaap_layer4_listeners":        dataSourceTencentCloudGaapLayer4Listeners(),
			"tencentcloud_gaap_layer7_listeners":        dataSourceTencentCloudGaapLayer7Listeners(),
			"tencentcloud_gaap_http_domains":            dataSourceTencentCloudGaapHttpDomains(),
			"tencentcloud_gaap_http_rules":              dataSourceTencentCloudGaapHttpRules(),
			"tencentcloud_gaap_security_policies":       dataSourceTencentCloudGaapSecurityPolices(),
			"tencentcloud_gaap_security_rules":          dataSourceTencentCloudGaapSecurityRules(),
			"tencentcloud_gaap_certificates":            dataSourceTencentCloudGaapCertificates(),
			"tencentcloud_ssl_certificates":             dataSourceTencentCloudSslCertificates(),
			"tencentcloud_instances":                    dataSourceTencentCloudInstances(),
			"tencentcloud_placement_groups":             dataSourceTencentCloudPlacementGroups(),
			"tencentcloud_eips":                         dataSourceTencentCloudEips(),
			"tencentcloud_key_pairs":                    dataSourceTencentCloudKeyPairs(),
			"tencentcloud_enis":                         dataSourceTencentCloudEnis(),
			"tencentcloud_cam_roles":                    dataSourceTencentCloudCamRoles(),
			"tencentcloud_cam_users":                    dataSourceTencentCloudCamUsers(),
			"tencentcloud_cam_groups":                   dataSourceTencentCloudCamGroups(),
			"tencentcloud_cam_group_memberships":        dataSourceTencentCloudCamGroupMemberships(),
			"tencentcloud_cam_policies":                 dataSourceTencentCloudCamPolicies(),
			"tencentcloud_cam_role_policy_attachments":  dataSourceTencentCloudCamRolePolicyAttachments(),
			"tencentcloud_cam_user_policy_attachments":  dataSourceTencentCloudCamUserPolicyAttachments(),
			"tencentcloud_cam_group_policy_attachments": dataSourceTencentCloudCamGroupPolicyAttachments(),
			"tencentcloud_cam_saml_providers":           dataSourceTencentCloudCamSAMLProviders(),
			"tencentcloud_reserved_instance_configs":    dataSourceTencentCloudReservedInstanceConfigs(),
			"tencentcloud_reserved_instances":           dataSourceTencentCloudReservedInstances(),
			"tencentcloud_cfs_file_systems":             dataSourceTencentCloudCfsFileSystems(),
			"tencentcloud_cfs_access_groups":            dataSourceTencentCloudCfsAccessGroups(),
			"tencentcloud_cfs_access_rules":             dataSourceTencentCloudCfsAccessRules(),
			"tencentcloud_scf_functions":                dataSourceTencentCloudScfFunctions(),
			"tencentcloud_scf_namespaces":               dataSourceTencentCloudScfNamespaces(),
			"tencentcloud_scf_logs":                     dataSourceTencentCloudScfLogs(),
			"tencentcloud_vpn_customer_gateways":        dataSourceTencentCloudVpnCustomerGateways(),
			"tencentcloud_vpn_gateways":                 dataSourceTencentCloudVpnGateways(),
			"tencentcloud_vpn_connections":              dataSourceTencentCloudVpnConnections(),
			"tencentcloud_ha_vips":                      dataSourceTencentCloudHaVips(),
			"tencentcloud_ha_vip_eip_attachments":       dataSourceTencentCloudHaVipEipAttachments(),
			"tencentcloud_tcaplus_applications":         dataSourceTencentCloudTcaplusApplications(),
			"tencentcloud_tcaplus_zones":                dataSourceTencentCloudTcaplusZones(),
			"tencentcloud_tcaplus_tables":               dataSourceTencentCloudTcaplusTables(),
			"tencentcloud_tcaplus_idls":                 dataSourceTencentCloudTcaplusIdls(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"tencentcloud_alb_server_attachment":          resourceTencentCloudAlbServerAttachment(),
			"tencentcloud_cbs_snapshot":                   resourceTencentCloudCbsSnapshot(),
			"tencentcloud_cbs_snapshot_policy":            resourceTencentCloudCbsSnapshotPolicy(),
			"tencentcloud_cbs_storage":                    resourceTencentCloudCbsStorage(),
			"tencentcloud_cbs_storage_attachment":         resourceTencentCloudCbsStorageAttachment(),
			"tencentcloud_cbs_snapshot_policy_attachment": resourceTencentCloudCbsSnapshotPolicyAttachment(),
			"tencentcloud_clb_instance":                   resourceTencentCloudClbInstance(),
			"tencentcloud_clb_listener":                   resourceTencentCloudClbListener(),
			"tencentcloud_clb_listener_rule":              resourceTencentCloudClbListenerRule(),
			"tencentcloud_clb_attachment":                 resourceTencentCloudClbServerAttachment(),
			"tencentcloud_clb_redirection":                resourceTencentCloudClbRedirection(),
			"tencentcloud_container_cluster":              resourceTencentCloudContainerCluster(),
			"tencentcloud_container_cluster_instance":     resourceTencentCloudContainerClusterInstance(),
			"tencentcloud_dnat":                           resourceTencentCloudDnat(),
			"tencentcloud_eip":                            resourceTencentCloudEip(),
			"tencentcloud_eip_association":                resourceTencentCloudEipAssociation(),
			"tencentcloud_instance":                       resourceTencentCloudInstance(),
			"tencentcloud_key_pair":                       resourceTencentCloudKeyPair(),
			"tencentcloud_lb":                             resourceTencentCloudLB(),
			"tencentcloud_nat_gateway":                    resourceTencentCloudNatGateway(),
			"tencentcloud_route_entry":                    resourceTencentCloudRouteEntry(),
			"tencentcloud_route_table_entry":              resourceTencentCloudVpcRouteEntry(),
			"tencentcloud_route_table":                    resourceTencentCloudVpcRouteTable(),
			"tencentcloud_security_group":                 resourceTencentCloudSecurityGroup(),
			"tencentcloud_security_group_rule":            resourceTencentCloudSecurityGroupRule(),
			"tencentcloud_subnet":                         resourceTencentCloudVpcSubnet(),
			"tencentcloud_vpc":                            resourceTencentCloudVpcInstance(),
			"tencentcloud_mysql_backup_policy":            resourceTencentCloudMysqlBackupPolicy(),
			"tencentcloud_mysql_account":                  resourceTencentCloudMysqlAccount(),
			"tencentcloud_mysql_account_privilege":        resourceTencentCloudMysqlAccountPrivilege(),
			"tencentcloud_mysql_privilege":                resourceTencentCloudMysqlPrivilege(),
			"tencentcloud_mysql_instance":                 resourceTencentCloudMysqlInstance(),
			"tencentcloud_mysql_readonly_instance":        resourceTencentCloudMysqlReadonlyInstance(),
			"tencentcloud_cos_bucket":                     resourceTencentCloudCosBucket(),
			"tencentcloud_cos_bucket_object":              resourceTencentCloudCosBucketObject(),
			"tencentcloud_redis_instance":                 resourceTencentCloudRedisInstance(),
			"tencentcloud_redis_backup_config":            resourceTencentCloudRedisBackupConfig(),
			"tencentcloud_as_scaling_config":              resourceTencentCloudAsScalingConfig(),
			"tencentcloud_as_scaling_group":               resourceTencentCloudAsScalingGroup(),
			"tencentcloud_as_attachment":                  resourceTencentCloudAsAttachment(),
			"tencentcloud_as_scaling_policy":              resourceTencentCloudAsScalingPolicy(),
			"tencentcloud_as_schedule":                    resourceTencentCloudAsSchedule(),
			"tencentcloud_as_lifecycle_hook":              resourceTencentCloudAsLifecycleHook(),
			"tencentcloud_as_notification":                resourceTencentCloudAsNotification(),
			"tencentcloud_ccn":                            resourceTencentCloudCcn(),
			"tencentcloud_ccn_attachment":                 resourceTencentCloudCcnAttachment(),
			"tencentcloud_ccn_bandwidth_limit":            resourceTencentCloudCcnBandwidthLimit(),
			"tencentcloud_dcx":                            resourceTencentCloudDcxInstance(),
			"tencentcloud_mongodb_instance":               resourceTencentCloudMongodbInstance(),
			"tencentcloud_mongodb_sharding_instance":      resourceTencentCloudMongodbShardingInstance(),
			"tencentcloud_dc_gateway":                     resourceTencentCloudDcGatewayInstance(),
			"tencentcloud_dc_gateway_ccn_route":           resourceTencentCloudDcGatewayCcnRouteInstance(),
			"tencentcloud_kubernetes_cluster":             resourceTencentCloudTkeCluster(),
			"tencentcloud_kubernetes_as_scaling_group":    ResourceTencentCloudKubernetesAsScalingGroup(),
			"tencentcloud_kubernetes_scale_worker":        resourceTencentCloudTkeScaleWorker(),
			"tencentcloud_gaap_proxy":                     resourceTencentCloudGaapProxy(),
			"tencentcloud_gaap_realserver":                resourceTencentCloudGaapRealserver(),
			"tencentcloud_gaap_layer4_listener":           resourceTencentCloudGaapLayer4Listener(),
			"tencentcloud_gaap_layer7_listener":           resourceTencentCloudGaapLayer7Listener(),
			"tencentcloud_gaap_http_domain":               resourceTencentCloudGaapHttpDomain(),
			"tencentcloud_gaap_http_rule":                 resourceTencentCloudGaapHttpRule(),
			"tencentcloud_gaap_certificate":               resourceTencentCloudGaapCertificate(),
			"tencentcloud_gaap_security_policy":           resourceTencentCloudGaapSecurityPolicy(),
			"tencentcloud_gaap_security_rule":             resourceTencentCloudGaapSecurityRule(),
			"tencentcloud_ssl_certificate":                resourceTencentCloudSslCertificate(),
			"tencentcloud_security_group_lite_rule":       resourceTencentCloudSecurityGroupLiteRule(),
			"tencentcloud_placement_group":                resourceTencentCloudPlacementGroup(),
			"tencentcloud_eni":                            resourceTencentCloudEni(),
			"tencentcloud_eni_attachment":                 resourceTencentCloudEniAttachment(),
			"tencentcloud_cam_role":                       resourceTencentCloudCamRole(),
			"tencentcloud_cam_user":                       resourceTencentCloudCamUser(),
			"tencentcloud_cam_policy":                     resourceTencentCloudCamPolicy(),
			"tencentcloud_cam_role_policy_attachment":     resourceTencentCloudCamRolePolicyAttachment(),
			"tencentcloud_cam_user_policy_attachment":     resourceTencentCloudCamUserPolicyAttachment(),
			"tencentcloud_cam_group_policy_attachment":    resourceTencentCloudCamGroupPolicyAttachment(),
			"tencentcloud_cam_group":                      resourceTencentCloudCamGroup(),
			"tencentcloud_cam_group_membership":           resourceTencentCloudCamGroupMembership(),
			"tencentcloud_cam_saml_provider":              resourceTencentCloudCamSAMLProvider(),
			"tencentcloud_reserved_instance":              resourceTencentCloudReservedInstance(),
			"tencentcloud_cfs_file_system":                resourceTencentCloudCfsFileSystem(),
			"tencentcloud_cfs_access_group":               resourceTencentCloudCfsAccessGroup(),
			"tencentcloud_cfs_access_rule":                resourceTencentCloudCfsAccessRule(),
			"tencentcloud_scf_function":                   resourceTencentCloudScfFunction(),
			"tencentcloud_scf_namespace":                  resourceTencentCloudScfNamespace(),
			"tencentcloud_vpn_customer_gateway":           resourceTencentCloudVpnCustomerGateway(),
			"tencentcloud_vpn_gateway":                    resourceTencentCloudVpnGateway(),
			"tencentcloud_vpn_connection":                 resourceTencentCloudVpnConnection(),
			"tencentcloud_ha_vip":                         resourceTencentCloudHaVip(),
			"tencentcloud_ha_vip_eip_attachment":          resourceTencentCloudHaVipEipAttachment(),
			"tencentcloud_tcaplus_application":            resourceTencentCloudTcaplusApplication(),
			"tencentcloud_tcaplus_zone":                   resourceTencentCloudTcaplusZone(),
			"tencentcloud_tcaplus_idl":                    resourceTencentCloudTcaplusIdl(),
			"tencentcloud_tcaplus_table":                  resourceTencentCloudTcaplusTable(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	secretId, ok := d.GetOk("secret_id")
	if !ok {
		secretId = os.Getenv(PROVIDER_SECRET_ID)
	}
	secretKey, ok := d.GetOk("secret_key")
	if !ok {
		secretKey = os.Getenv(PROVIDER_SECRET_KEY)
	}
	securityToken, ok := d.GetOk("security_token")
	if !ok {
		securityToken = os.Getenv(PROVIDER_SECURITY_TOKEN)
	}
	region, ok := d.GetOk("region")
	if !ok {
		region = os.Getenv(PROVIDER_REGION)
	}
	config := Config{
		SecretId:      secretId.(string),
		SecretKey:     secretKey.(string),
		SecurityToken: securityToken.(string),
		Region:        region.(string),
	}
	return config.Client()
}
