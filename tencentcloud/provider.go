/*
The TencentCloud provider is used to interact with many resources supported by TencentCloud. The provider needs to be configured with the proper credentials before it can be used.

The TencentCloud provider is used to interact with the many resources supported by [TencentCloud](https://intl.cloud.tencent.com).
The provider needs to be configured with the proper credentials before it can be used.

Use the navigation to the left to read about the available resources.

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
  tencentcloud_container_cluster_instances
  tencentcloud_container_clusters
  tencentcloud_cos_bucket_object
  tencentcloud_cos_buckets
  tencentcloud_dc_instances
  tencentcloud_dc_gateway_ccn_routes
  tencentcloud_dc_gateway_instances
  tencentcloud_dcx_instances
  tencentcloud_clb_instances
  tencentcloud_clb_listeners
  tencentcloud_clb_listener_rules
  tencentcloud_clb_attachments
  tencentcloud_clb_redirections
  tencentcloud_eip
  tencentcloud_image
  tencentcloud_instance_types
  tencentcloud_kubernetes_clusters
  tencentcloud_mongodb_instances
  tencentcloud_mongodb_zone_config
  tencentcloud_mysql_backup_list
  tencentcloud_mysql_instance
  tencentcloud_mysql_parameter_list
  tencentcloud_mysql_zone_config
  tencentcloud_nats
  tencentcloud_redis_instances
  tencentcloud_redis_zone_config
  tencentcloud_route_table
  tencentcloud_security_group
  tencentcloud_security_groups
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

DC Resources
  tencentcloud_dcx

DCG Resources
  tencentcloud_dc_gateway
  tencentcloud_dc_gateway_ccn_route

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

VPC Resources
  tencentcloud_vpc
  tencentcloud_subnet
  tencentcloud_security_group
  tencentcloud_security_group_rule
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
	PROVIDER_SECRET_ID  = "TENCENTCLOUD_SECRET_ID"
	PROVIDER_SECRET_KEY = "TENCENTCLOUD_SECRET_KEY"
	PROVIDER_REGION     = "TENCENTCLOUD_REGION"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"secret_id": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc(PROVIDER_SECRET_ID, nil),
				Description: "Secret ID of Tencent Cloud",
			},
			"secret_key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc(PROVIDER_SECRET_KEY, nil),
				Description: "Secret key of Tencent Cloud",
				Sensitive:   true,
			},
			"region": {
				Type:         schema.TypeString,
				Required:     true,
				DefaultFunc:  schema.EnvDefaultFunc(PROVIDER_REGION, nil),
				Description:  "Region of Tencent Cloud",
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
	region, ok := d.GetOk("region")
	if !ok {
		region = os.Getenv(PROVIDER_REGION)
	}
	config := Config{
		SecretId:  secretId.(string),
		SecretKey: secretKey.(string),
		Region:    region.(string),
	}
	return config.Client()
}
