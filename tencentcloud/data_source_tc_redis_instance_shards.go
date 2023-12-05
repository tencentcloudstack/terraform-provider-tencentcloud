package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudRedisInstanceShards() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudRedisInstanceShardsRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The ID of instance.",
			},

			"filter_slave": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Whether to filter out slave information.",
			},

			"instance_shards": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Instance shard list information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"shard_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Shard node name.",
						},
						"shard_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Shard node ID.",
						},
						"role": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "role.",
						},
						"keys": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of keys.",
						},
						"slots": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Slot information.",
						},
						"storage": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Used capacity.",
						},
						"storage_slope": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Capacity tilt.",
						},
						"runid": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The node ID of the instance runtime.",
						},
						"connected": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Service status: 0-down;1-on.",
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

func dataSourceTencentCloudRedisInstanceShardsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_redis_instance_shards.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	var instanceId string

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("filter_slave"); v != nil {
		paramMap["FilterSlave"] = helper.Bool(v.(bool))
	}

	service := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}

	var instanceShards []*redis.InstanceClusterShard

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeRedisInstanceShardsByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		instanceShards = result
		return nil
	})
	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0, len(instanceShards))
	if instanceShards != nil {
		for _, instanceClusterShard := range instanceShards {
			instanceClusterShardMap := map[string]interface{}{}

			if instanceClusterShard.ShardName != nil {
				instanceClusterShardMap["shard_name"] = instanceClusterShard.ShardName
			}

			if instanceClusterShard.ShardId != nil {
				instanceClusterShardMap["shard_id"] = instanceClusterShard.ShardId
			}

			if instanceClusterShard.Role != nil {
				instanceClusterShardMap["role"] = instanceClusterShard.Role
			}

			if instanceClusterShard.Keys != nil {
				instanceClusterShardMap["keys"] = instanceClusterShard.Keys
			}

			if instanceClusterShard.Slots != nil {
				instanceClusterShardMap["slots"] = instanceClusterShard.Slots
			}

			if instanceClusterShard.Storage != nil {
				instanceClusterShardMap["storage"] = instanceClusterShard.Storage
			}

			if instanceClusterShard.StorageSlope != nil {
				instanceClusterShardMap["storage_slope"] = instanceClusterShard.StorageSlope
			}

			if instanceClusterShard.Runid != nil {
				instanceClusterShardMap["runid"] = instanceClusterShard.Runid
			}

			if instanceClusterShard.Connected != nil {
				instanceClusterShardMap["connected"] = instanceClusterShard.Connected
			}

			tmpList = append(tmpList, instanceClusterShardMap)
		}

		_ = d.Set("instance_shards", tmpList)
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
