package crs

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudRedisInstanceNodeInfo() *schema.Resource {
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
	defer tccommon.LogElapsed("data_source.tencentcloud_redis_instance_node_info.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var instanceId string
	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	service := RedisService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	var instanceNodeInfo *redis.DescribeInstanceNodeInfoResponseParams
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeRedisInstanceNodeInfoByFilter(ctx, paramMap)
		if e != nil {
			if sdkerr, ok := e.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkerr.Code == "FailedOperation.SystemError" {
					return nil
				}
			}
			return tccommon.RetryError(e)
		}
		instanceNodeInfo = result
		return nil
	})
	if err != nil {
		return err
	}

	if instanceNodeInfo != nil {
		if instanceNodeInfo.ProxyCount != nil {
			_ = d.Set("proxy_count", instanceNodeInfo.ProxyCount)
		}

		if instanceNodeInfo.Proxy != nil {
			tmpList := make([]map[string]interface{}, 0, len(instanceNodeInfo.Proxy))
			for _, proxyNodes := range instanceNodeInfo.Proxy {
				proxyNodesMap := map[string]interface{}{}

				if proxyNodes.NodeId != nil {
					proxyNodesMap["node_id"] = proxyNodes.NodeId
				}

				if proxyNodes.ZoneId != nil {
					proxyNodesMap["zone_id"] = proxyNodes.ZoneId
				}

				tmpList = append(tmpList, proxyNodesMap)
			}

			_ = d.Set("proxy", tmpList)
		}

		if instanceNodeInfo.RedisCount != nil {
			_ = d.Set("redis_count", instanceNodeInfo.RedisCount)
		}

		if instanceNodeInfo.Redis != nil {
			tmpList := make([]map[string]interface{}, 0, len(instanceNodeInfo.Redis))
			for _, redisNodes := range instanceNodeInfo.Redis {
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

				tmpList = append(tmpList, redisNodesMap)
			}

			_ = d.Set("redis", tmpList)
		}
	}

	d.SetId(instanceId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}
	return nil
}
