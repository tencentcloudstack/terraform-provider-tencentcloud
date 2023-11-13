/*
Use this data source to query detailed information of redis instance_node_info

Example Usage

```hcl
data "tencentcloud_redis_instance_node_info" "instance_node_info" {
  instance_id = "crs-c1nl9rpv"
        }
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudRedisInstanceNodeInfo() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudRedisInstanceNodeInfoRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The ID of instance.",
			},

			"proxy_count": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Number of proxy nodes.",
			},

			"proxy": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Proxy node information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Node ID.",
						},
						"zone_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Zone ID.",
						},
					},
				},
			},

			"redis_count": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Number of redis nodes.",
			},

			"redis": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Redis node information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Node ID.",
						},
						"node_role": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Node role.",
						},
						"cluster_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Shard ID.",
						},
						"zone_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Zone ID.",
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

func dataSourceTencentCloudRedisInstanceNodeInfoRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_redis_instance_node_info.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	service := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeRedisInstanceNodeInfoByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		proxyCount = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(proxyCount))
	if proxyCount != nil {
		_ = d.Set("proxy_count", proxyCount)
	}

	if proxy != nil {
		for _, proxyNodes := range proxy {
			proxyNodesMap := map[string]interface{}{}

			if proxyNodes.NodeId != nil {
				proxyNodesMap["node_id"] = proxyNodes.NodeId
			}

			if proxyNodes.ZoneId != nil {
				proxyNodesMap["zone_id"] = proxyNodes.ZoneId
			}

			ids = append(ids, *proxyNodes.InstanceId)
			tmpList = append(tmpList, proxyNodesMap)
		}

		_ = d.Set("proxy", tmpList)
	}

	if redisCount != nil {
		_ = d.Set("redis_count", redisCount)
	}

	if redis != nil {
		for _, redisNodes := range redis {
			redisNodesMap := map[string]interface{}{}

			if redisNodes.NodeId != nil {
				redisNodesMap["node_id"] = redisNodes.NodeId
			}

			if redisNodes.NodeRole != nil {
				redisNodesMap["node_role"] = redisNodes.NodeRole
			}

			if redisNodes.ClusterId != nil {
				redisNodesMap["cluster_id"] = redisNodes.ClusterId
			}

			if redisNodes.ZoneId != nil {
				redisNodesMap["zone_id"] = redisNodes.ZoneId
			}

			ids = append(ids, *redisNodes.InstanceId)
			tmpList = append(tmpList, redisNodesMap)
		}

		_ = d.Set("redis", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string)); e != nil {
			return e
		}
	}
	return nil
}
