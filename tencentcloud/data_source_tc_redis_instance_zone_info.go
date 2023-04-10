/*
Use this data source to query detailed information of redis instance_zone_info

Example Usage

```hcl
data "tencentcloud_redis_instance_zone_info" "instance_zone_info" {
  instance_id = "crs-c1nl9rpv"
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudRedisInstanceZoneInfo() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudRedisInstanceZoneInfoRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The ID of instance.",
			},

			"replica_groups": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "List of instance node groups.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Node group ID.",
						},
						"group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Node group Name.",
						},
						"zone_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "he availability zone ID of the node, such as ap-guangzhou-1.",
						},
						"role": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The node group type, master is the primary node, and replica is the replica node.",
						},
						"redis_nodes": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Node group node list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"keys": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The number of node keys.",
									},
									"slot": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Node slot distribution.",
									},
									"node_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Node ID.",
									},
									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Node status.",
									},
									"role": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Node role.",
									},
								},
							},
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

func dataSourceTencentCloudRedisInstanceZoneInfoRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_redis_instance_zone_info.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	var instanceId string

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	service := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}

	var replicaGroups []*redis.ReplicaGroup

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeRedisInstanceZoneInfoByFilter(ctx, paramMap)
		if e != nil {
			if ee, ok := e.(*sdkErrors.TencentCloudSDKError); ok {
				if ee.Code == "FailedOperation.UnSupportError" {
					return resource.NonRetryableError(e)
				}
			}
			return retryError(e)
		}
		replicaGroups = result
		return nil
	})
	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0, len(replicaGroups))

	if replicaGroups != nil {
		for _, replicaGroup := range replicaGroups {
			replicaGroupMap := map[string]interface{}{}

			if replicaGroup.GroupId != nil {
				replicaGroupMap["group_id"] = replicaGroup.GroupId
			}

			if replicaGroup.GroupName != nil {
				replicaGroupMap["group_name"] = replicaGroup.GroupName
			}

			if replicaGroup.ZoneId != nil {
				replicaGroupMap["zone_id"] = replicaGroup.ZoneId
			}

			if replicaGroup.Role != nil {
				replicaGroupMap["role"] = replicaGroup.Role
			}

			if replicaGroup.RedisNodes != nil {
				redisNodesList := []interface{}{}
				for _, redisNodes := range replicaGroup.RedisNodes {
					redisNodesMap := map[string]interface{}{}

					if redisNodes.Keys != nil {
						redisNodesMap["keys"] = redisNodes.Keys
					}

					if redisNodes.Slot != nil {
						redisNodesMap["slot"] = redisNodes.Slot
					}

					if redisNodes.NodeId != nil {
						redisNodesMap["node_id"] = redisNodes.NodeId
					}

					if redisNodes.Status != nil {
						redisNodesMap["status"] = redisNodes.Status
					}

					if redisNodes.Role != nil {
						redisNodesMap["role"] = redisNodes.Role
					}

					redisNodesList = append(redisNodesList, redisNodesMap)
				}

				replicaGroupMap["redis_nodes"] = redisNodesList
			}

			tmpList = append(tmpList, replicaGroupMap)
		}

		_ = d.Set("replica_groups", tmpList)
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
