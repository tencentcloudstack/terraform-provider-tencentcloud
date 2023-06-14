/*
Use this data source to query detailed information of dbbrain redis_top_key_prefix_list

Example Usage

```hcl
data "tencentcloud_dbbrain_redis_top_key_prefix_list" "redis_top_key_prefix_list" {
	instance_id = local.redis_id
	date        = "%s"
	product     = "redis"
}
```
*/
package tencentcloud

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dbbrain "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbbrain/v20210527"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudDbbrainRedisTopKeyPrefixList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDbbrainRedisTopKeyPrefixListRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "instance id.",
			},

			"date": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Query date, such as 2021-05-27, the earliest date can be the previous 30 days.",
			},

			"product": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Service product type, supported values include `redis` - cloud database Redis.",
			},

			"items": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "list of top key prefixes.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ave_element_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Average element length.",
						},
						"length": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total occupied memory (Byte).",
						},
						"key_pre_index": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "key prefix.",
						},
						"item_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "number of elements.",
						},
						"count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of keys.",
						},
						"max_element_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Maximum element length.",
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

func dataSourceTencentCloudDbbrainRedisTopKeyPrefixListRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_dbbrain_redis_top_key_prefix_list.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	var (
		instanceId string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("date"); ok {
		paramMap["Date"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("product"); ok {
		paramMap["Product"] = helper.String(v.(string))
	}

	service := DbbrainService{client: meta.(*TencentCloudClient).apiV3Conn}

	var items []*dbbrain.RedisPreKeySpaceData

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDbbrainRedisTopKeyPrefixListByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		items = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(items))
	tmpList := make([]map[string]interface{}, 0, len(items))

	if items != nil {
		for _, redisPreKeySpaceData := range items {
			redisPreKeySpaceDataMap := map[string]interface{}{}

			if redisPreKeySpaceData.AveElementSize != nil {
				redisPreKeySpaceDataMap["ave_element_size"] = redisPreKeySpaceData.AveElementSize
			}

			if redisPreKeySpaceData.Length != nil {
				redisPreKeySpaceDataMap["length"] = redisPreKeySpaceData.Length
			}

			if redisPreKeySpaceData.KeyPreIndex != nil {
				redisPreKeySpaceDataMap["key_pre_index"] = redisPreKeySpaceData.KeyPreIndex
			}

			if redisPreKeySpaceData.ItemCount != nil {
				redisPreKeySpaceDataMap["item_count"] = redisPreKeySpaceData.ItemCount
			}

			if redisPreKeySpaceData.Count != nil {
				redisPreKeySpaceDataMap["count"] = redisPreKeySpaceData.Count
			}

			if redisPreKeySpaceData.MaxElementSize != nil {
				redisPreKeySpaceDataMap["max_element_size"] = redisPreKeySpaceData.MaxElementSize
			}

			ids = append(ids, strings.Join([]string{instanceId, *redisPreKeySpaceData.KeyPreIndex}, FILED_SP))
			tmpList = append(tmpList, redisPreKeySpaceDataMap)
		}

		_ = d.Set("items", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
