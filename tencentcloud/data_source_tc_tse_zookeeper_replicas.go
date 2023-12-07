package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tse/v20201207"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudTseZookeeperReplicas() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTseZookeeperReplicasRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "engine instance ID.",
			},

			"replicas": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Engine instance replica information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "name.",
						},
						"role": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "role.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "status.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Subnet IDNote: This field may return null, indicating that a valid value is not available.",
						},
						"zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Available area IDNote: This field may return null, indicating that a valid value is not available.",
						},
						"zone_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Available area IDNote: This field may return null, indicating that a valid value is not available.",
						},
						"alias_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "aliasNote: This field may return null, indicating that a valid value is not available.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VPC IDNote: This field may return null, indicating that a valid value is not available.",
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudTseZookeeperReplicasRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_tse_zookeeper_replicas.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	instanceId := ""

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	service := TseService{client: meta.(*TencentCloudClient).apiV3Conn}

	var replicas []*tse.ZookeeperReplica

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTseZookeeperReplicasByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		replicas = result
		return nil
	})
	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0, len(replicas))

	if replicas != nil {
		for _, zookeeperReplica := range replicas {
			zookeeperReplicaMap := map[string]interface{}{}

			if zookeeperReplica.Name != nil {
				zookeeperReplicaMap["name"] = zookeeperReplica.Name
			}

			if zookeeperReplica.Role != nil {
				zookeeperReplicaMap["role"] = zookeeperReplica.Role
			}

			if zookeeperReplica.Status != nil {
				zookeeperReplicaMap["status"] = zookeeperReplica.Status
			}

			if zookeeperReplica.SubnetId != nil {
				zookeeperReplicaMap["subnet_id"] = zookeeperReplica.SubnetId
			}

			if zookeeperReplica.Zone != nil {
				zookeeperReplicaMap["zone"] = zookeeperReplica.Zone
			}

			if zookeeperReplica.ZoneId != nil {
				zookeeperReplicaMap["zone_id"] = zookeeperReplica.ZoneId
			}

			if zookeeperReplica.AliasName != nil {
				zookeeperReplicaMap["alias_name"] = zookeeperReplica.AliasName
			}

			if zookeeperReplica.VpcId != nil {
				zookeeperReplicaMap["vpc_id"] = zookeeperReplica.VpcId
			}

			tmpList = append(tmpList, zookeeperReplicaMap)
		}

		_ = d.Set("replicas", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash([]string{instanceId}))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
