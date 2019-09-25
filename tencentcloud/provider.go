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
  tencentcloud_cbs_snapshots
  tencentcloud_cbs_storages
  tencentcloud_ccn_bandwidth_limits
  tencentcloud_ccn_instances
  tencentcloud_clb_instances
  tencentcloud_clb_listeners
  tencentcloud_clb_listener_rules
  tencentcloud_clb_attachments
  tencentcloud_clb_redirections
  tencentcloud_container_cluster_instances
  tencentcloud_container_clusters
  tencentcloud_cos_bucket_object
  tencentcloud_cos_buckets
  tencentcloud_dc_instances
  tencentcloud_dc_gateway_ccn_routes
  tencentcloud_dc_gateway_instances
  tencentcloud_dcx_instances
  tencentcloud_dnats
  tencentcloud_eip
  tencentcloud_eips
  tencentcloud_gaap_certificates
  tencentcloud_gaap_http_domains
  tencentcloud_gaap_http_rules
  tencentcloud_gaap_layer4_listeners
  tencentcloud_gaap_layer7_listeners
  tencentcloud_gaap_proxies
  tencentcloud_gaap_realservers
  tencentcloud_gaap_security_policies
  tencentcloud_gaap_security_rules
  tencentcloud_image
  tencentcloud_instances
  tencentcloud_instance_types
  tencentcloud_key_pairs
  tencentcloud_kubernetes_clusters
  tencentcloud_mongodb_instances
  tencentcloud_mongodb_zone_config
  tencentcloud_mysql_backup_list
  tencentcloud_mysql_instance
  tencentcloud_mysql_parameter_list
  tencentcloud_mysql_zone_config
  tencentcloud_nats
  tencentcloud_nat_gateways
  tencentcloud_placement_groups
  tencentcloud_redis_instances
  tencentcloud_redis_zone_config
  tencentcloud_route_table
  tencentcloud_security_group
  tencentcloud_security_groups
  tencentcloud_ssl_certificates
  tencentcloud_subnet
  tencentcloud_vpc
  tencentcloud_vpc_instances
  tencentcloud_vpc_route_tables
  tencentcloud_vpc_subnets

AS Resources
  tencentcloud_as_scaling_config
  tencentcloud_as_scaling_group
  tencentcloud_as_attachment
  tencentcloud_as_scaling_policy
  tencentcloud_as_schedule
  tencentcloud_as_lifecycle_hook
  tencentcloud_as_notification

CBS Resources
  tencentcloud_cbs_storage
  tencentcloud_cbs_storage_attachment
  tencentcloud_cbs_snapshot
  tencentcloud_cbs_snapshot_policy

CCN Resources
  tencentcloud_ccn
  tencentcloud_ccn_attachment
  tencentcloud_ccn_bandwidth_limit

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

MongoDB Resources
  tencentcloud_mongodb_instance
  tencentcloud_mongodb_sharding_instance

MySQL Resources
  tencentcloud_mysql_instance
  tencentcloud_mysql_readonly_instance
  tencentcloud_mysql_account
  tencentcloud_mysql_account_privilege
  tencentcloud_mysql_backup_policy

Redis Resources
  tencentcloud_redis_instance
  tencentcloud_redis_backup_config

SSL Resources
  tencentcloud_ssl_certificate

VPC Resources
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

*/
package tencentcloud

import (
	"os"

	"github.com/hashicorp/terraform/helper/schema"
)

const (
	PROVIDER_SECRET_ID      = "TENCENTCLOUD_SECRET_ID"
	PROVIDER_SECRET_KEY     = "TENCENTCLOUD_SECRET_KEY"
	PROVIDER_SECURITY_TOKEN = "TENCENTCLOUD_SECURITY_TOKEN"
	PROVIDER_REGION         = "TENCENTCLOUD_REGION"
)

func Provider() *schema.Provider {
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
			"tencentcloud_availability_zones":          dataSourceTencentCloudAvailabilityZones(),
			"tencentcloud_eip":                         dataSourceTencentCloudEip(),
			"tencentcloud_image":                       dataSourceTencentCloudSourceImages(),
			"tencentcloud_instance_types":              dataSourceInstanceTypes(),
			"tencentcloud_vpc":                         dataSourceTencentCloudVpc(),
			"tencentcloud_subnet":                      dataSourceTencentCloudSubnet(),
			"tencentcloud_route_table":                 dataSourceTencentCloudRouteTable(),
			"tencentcloud_security_group":              dataSourceTencentCloudSecurityGroup(),
			"tencentcloud_security_groups":             dataSourceTencentCloudSecurityGroups(),
			"tencentcloud_nats":                        dataSourceTencentCloudNats(),
			"tencentcloud_dnats":                       dataSourceTencentCloudDnats(),
			"tencentcloud_nat_gateways":                dataSourceTencentCloudNatGateways(),
			"tencentcloud_container_clusters":          dataSourceTencentCloudContainerClusters(),
			"tencentcloud_container_cluster_instances": dataSourceTencentCloudContainerClusterInstances(),
			"tencentcloud_mysql_backup_list":           dataSourceTencentMysqlBackupList(),
			"tencentcloud_mysql_zone_config":           dataSourceTencentMysqlZoneConfig(),
			"tencentcloud_mysql_parameter_list":        dataSourceTencentCloudMysqlParameterList(),
			"tencentcloud_mysql_instance":              dataSourceTencentCloudMysqlInstance(),
			"tencentcloud_cos_bucket_object":           dataSourceTencentCloudCosBucketObject(),
			"tencentcloud_cos_buckets":                 dataSourceTencentCloudCosBuckets(),
			"tencentcloud_redis_zone_config":           dataSourceTencentRedisZoneConfig(),
			"tencentcloud_redis_instances":             dataSourceTencentRedisInstances(),
			"tencentcloud_as_scaling_configs":          dataSourceTencentCloudAsScalingConfigs(),
			"tencentcloud_as_scaling_groups":           dataSourceTencentCloudAsScalingGroups(),
			"tencentcloud_as_scaling_policies":         dataSourceTencentCloudAsScalingPolicies(),
			"tencentcloud_vpc_instances":               dataSourceTencentCloudVpcInstances(),
			"tencentcloud_vpc_subnets":                 dataSourceTencentCloudVpcSubnets(),
			"tencentcloud_vpc_route_tables":            dataSourceTencentCloudVpcRouteTables(),
			"tencentcloud_ccn_instances":               dataSourceTencentCloudCcnInstances(),
			"tencentcloud_ccn_bandwidth_limits":        dataSourceTencentCloudCcnBandwidthLimits(),
			"tencentcloud_cbs_storages":                dataSourceTencentCloudCbsStorages(),
			"tencentcloud_cbs_snapshots":               dataSourceTencentCloudCbsSnapshots(),
			"tencentcloud_dc_instances":                dataSourceTencentCloudDcInstances(),
			"tencentcloud_clb_instances":               dataSourceTencentCloudClbInstances(),
			"tencentcloud_clb_listeners":               dataSourceTencentCloudClbListeners(),
			"tencentcloud_clb_listener_rules":          dataSourceTencentCloudClbListenerRules(),
			"tencentcloud_clb_attachments":             dataSourceTencentCloudClbServerAttachments(),
			"tencentcloud_clb_redirections":            dataSourceTencentCloudClbRedirections(),
			"tencentcloud_dcx_instances":               dataSourceTencentCloudDcxInstances(),
			"tencentcloud_mongodb_zone_config":         dataSourceTencentCloudMongodbZoneConfig(),
			"tencentcloud_mongodb_instances":           dataSourceTencentCloudMongodbInstances(),
			"tencentcloud_dc_gateway_instances":        dataSourceTencentCloudDcGatewayInstances(),
			"tencentcloud_dc_gateway_ccn_routes":       dataSourceTencentCloudDcGatewayCCNRoutes(),
			"tencentcloud_kubernetes_clusters":         dataSourceTencentCloudKubernetesClusters(),
			"tencentcloud_gaap_proxies":                dataSourceTencentCloudGaapProxies(),
			"tencentcloud_gaap_realservers":            dataSourceTencentCloudGaapRealservers(),
			"tencentcloud_gaap_layer4_listeners":       dataSourceTencentCloudGaapLayer4Listeners(),
			"tencentcloud_gaap_layer7_listeners":       dataSourceTencentCloudGaapLayer7Listeners(),
			"tencentcloud_gaap_http_domains":           dataSourceTencentCloudGaapHttpDomains(),
			"tencentcloud_gaap_http_rules":             dataSourceTencentCloudGaapHttpRules(),
			"tencentcloud_gaap_security_policies":      dataSourceTencentCloudGaapSecurityPolices(),
			"tencentcloud_gaap_security_rules":         dataSourceTencentCloudGaapSecurityRules(),
			"tencentcloud_gaap_certificates":           dataSourceTencentCloudGaapCertificates(),
			"tencentcloud_ssl_certificates":            dataSourceTencentCloudSslCertificates(),
			"tencentcloud_instances":                   dataSourceTencentCloudInstances(),
			"tencentcloud_placement_groups":            dataSourceTencentCloudPlacementGroups(),
			"tencentcloud_eips":                        dataSourceTencentCloudEips(),
			"tencentcloud_key_pairs":                   dataSourceTencentCloudKeyPairs(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"tencentcloud_alb_server_attachment":      resourceTencentCloudAlbServerAttachment(),
			"tencentcloud_cbs_snapshot":               resourceTencentCloudCbsSnapshot(),
			"tencentcloud_cbs_snapshot_policy":        resourceTencentCloudCbsSnapshotPolicy(),
			"tencentcloud_cbs_storage":                resourceTencentCloudCbsStorage(),
			"tencentcloud_cbs_storage_attachment":     resourceTencentCloudCbsStorageAttachment(),
			"tencentcloud_clb_instance":               resourceTencentCloudClbInstance(),
			"tencentcloud_clb_listener":               resourceTencentCloudClbListener(),
			"tencentcloud_clb_listener_rule":          resourceTencentCloudClbListenerRule(),
			"tencentcloud_clb_attachment":             resourceTencentCloudClbServerAttachment(),
			"tencentcloud_clb_redirection":            resourceTencentCloudClbRedirection(),
			"tencentcloud_container_cluster":          resourceTencentCloudContainerCluster(),
			"tencentcloud_container_cluster_instance": resourceTencentCloudContainerClusterInstance(),
			"tencentcloud_dnat":                       resourceTencentCloudDnat(),
			"tencentcloud_eip":                        resourceTencentCloudEip(),
			"tencentcloud_eip_association":            resourceTencentCloudEipAssociation(),
			"tencentcloud_instance":                   resourceTencentCloudInstance(),
			"tencentcloud_key_pair":                   resourceTencentCloudKeyPair(),
			"tencentcloud_lb":                         resourceTencentCloudLB(),
			"tencentcloud_nat_gateway":                resourceTencentCloudNatGateway(),
			"tencentcloud_route_entry":                resourceTencentCloudRouteEntry(),
			"tencentcloud_route_table_entry":          resourceTencentCloudVpcRouteEntry(),
			"tencentcloud_route_table":                resourceTencentCloudVpcRouteTable(),
			"tencentcloud_security_group":             resourceTencentCloudSecurityGroup(),
			"tencentcloud_security_group_rule":        resourceTencentCloudSecurityGroupRule(),
			"tencentcloud_subnet":                     resourceTencentCloudVpcSubnet(),
			"tencentcloud_vpc":                        resourceTencentCloudVpcInstance(),
			"tencentcloud_mysql_backup_policy":        resourceTencentCloudMysqlBackupPolicy(),
			"tencentcloud_mysql_account":              resourceTencentCloudMysqlAccount(),
			"tencentcloud_mysql_account_privilege":    resourceTencentCloudMysqlAccountPrivilege(),
			"tencentcloud_mysql_instance":             resourceTencentCloudMysqlInstance(),
			"tencentcloud_mysql_readonly_instance":    resourceTencentCloudMysqlReadonlyInstance(),
			"tencentcloud_cos_bucket":                 resourceTencentCloudCosBucket(),
			"tencentcloud_cos_bucket_object":          resourceTencentCloudCosBucketObject(),
			"tencentcloud_redis_instance":             resourceTencentCloudRedisInstance(),
			"tencentcloud_redis_backup_config":        resourceTencentCloudRedisBackupConfig(),
			"tencentcloud_as_scaling_config":          resourceTencentCloudAsScalingConfig(),
			"tencentcloud_as_scaling_group":           resourceTencentCloudAsScalingGroup(),
			"tencentcloud_as_attachment":              resourceTencentCloudAsAttachment(),
			"tencentcloud_as_scaling_policy":          resourceTencentCloudAsScalingPolicy(),
			"tencentcloud_as_schedule":                resourceTencentCloudAsSchedule(),
			"tencentcloud_as_lifecycle_hook":          resourceTencentCloudAsLifecycleHook(),
			"tencentcloud_as_notification":            resourceTencentCloudAsNotification(),
			"tencentcloud_ccn":                        resourceTencentCloudCcn(),
			"tencentcloud_ccn_attachment":             resourceTencentCloudCcnAttachment(),
			"tencentcloud_ccn_bandwidth_limit":        resourceTencentCloudCcnBandwidthLimit(),
			"tencentcloud_dcx":                        resourceTencentCloudDcxInstance(),
			"tencentcloud_mongodb_instance":           resourceTencentCloudMongodbInstance(),
			"tencentcloud_mongodb_sharding_instance":  resourceTencentCloudMongodbShardingInstance(),
			"tencentcloud_dc_gateway":                 resourceTencentCloudDcGatewayInstance(),
			"tencentcloud_dc_gateway_ccn_route":       resourceTencentCloudDcGatewayCcnRouteInstance(),
			"tencentcloud_kubernetes_cluster":         resourceTencentCloudTkeCluster(),
			"tencentcloud_kubernetes_scale_worker":    resourceTencentCloudTkeScaleWorker(),
			"tencentcloud_gaap_proxy":                 resourceTencentCloudGaapProxy(),
			"tencentcloud_gaap_realserver":            resourceTencentCloudGaapRealserver(),
			"tencentcloud_gaap_layer4_listener":       resourceTencentCloudGaapLayer4Listener(),
			"tencentcloud_gaap_layer7_listener":       resourceTencentCloudGaapLayer7Listener(),
			"tencentcloud_gaap_http_domain":           resourceTencentCloudGaapHttpDomain(),
			"tencentcloud_gaap_http_rule":             resourceTencentCloudGaapHttpRule(),
			"tencentcloud_gaap_certificate":           resourceTencentCloudGaapCertificate(),
			"tencentcloud_gaap_security_policy":       resourceTencentCloudGaapSecurityPolicy(),
			"tencentcloud_gaap_security_rule":         resourceTencentCloudGaapSecurityRule(),
			"tencentcloud_ssl_certificate":            resourceTencentCloudSslCertificate(),
			"tencentcloud_security_group_lite_rule":   resourceTencentCloudSecurityGroupLiteRule(),
			"tencentcloud_placement_group":            resourceTencentCloudPlacementGroup(),
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
