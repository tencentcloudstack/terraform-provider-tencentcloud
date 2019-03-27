package tencentcloud

import (
	"os"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

const (
	PROVIDER_SECRET_ID  = "TENCENTCLOUD_SECRET_ID"
	PROVIDER_SECRET_KEY = "TENCENTCLOUD_SECRET_KEY"
	PROVIDER_REGION     = "TENCENTCLOUD_REGION"
	PROVIDER_DEBUG      = "TENCENTCLOUD_DEBUG"
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
			"tencentcloud_nats":                        dataSourceTencentCloudNats(),
			"tencentcloud_container_clusters":          dataSourceTencentCloudContainerClusters(),
			"tencentcloud_container_cluster_instances": dataSourceTencentCloudContainerClusterInstances(),
			"tencentcloud_mysql_backup_list":           dataSourceTencentMysqlBackupList(),
			"tencentcloud_mysql_zone_config":           dataSourceTencentMysqlZoneConfig(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"tencentcloud_alb_server_attachment":      resourceTencentCloudAlbServerAttachment(),
			"tencentcloud_cbs_snapshot":               resourceTencentCloudCbsSnapshot(),
			"tencentcloud_cbs_storage":                resourceTencentCloudCbsStorage(),
			"tencentcloud_cbs_storage_attachment":     resourceTencentCloudCbsStorageAttachment(),
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
			"tencentcloud_route_table":                resourceTencentCloudRouteTable(),
			"tencentcloud_security_group":             resourceTencentCloudSecurityGroup(),
			"tencentcloud_security_group_rule":        resourceTencentCloudSecurityGroupRule(),
			"tencentcloud_subnet":                     resourceTencentCloudSubnet(),
			"tencentcloud_vpc":                        resourceTencentCloudVpc(),
			"tencentcloud_mysql_backup_policy":        resourceTencentCloudMysqlBackupPolicy(),
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
	if strings.TrimSpace(os.Getenv(PROVIDER_DEBUG)) != "" {
		InitLogConfig(true)
	} else {
		InitLogConfig(true)
	}
	return config.Client()
}
