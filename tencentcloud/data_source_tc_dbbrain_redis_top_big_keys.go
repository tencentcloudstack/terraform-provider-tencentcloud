/*
Use this data source to query detailed information of dbbrain redis_top_big_keys

Example Usage

```hcl
data "tencentcloud_dbbrain_redis_top_big_keys" "redis_top_big_keys" {
  instance_id = ""
  date = ""
  product = ""
  sort_by = ""
  key_type = ""
    }
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dbbrain "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbbrain/v20210527"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudDbbrainRedisTopBigKeys() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDbbrainRedisTopBigKeysRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance id.",
			},

			"date": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Query date, such as 2021-05-27, the earliest date can be the previous 30 days.",
			},

			"product": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Service product type, supported values include redis - cloud database Redis.",
			},

			"sort_by": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sorting field, the value includes Capacity - memory, ItemCount - number of elements, the default is Capacity.",
			},

			"key_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Key type filter condition, the default is no filter, the value includes string, list, set, hash, sortedset, stream.",
			},

			"top_keys": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "List of top keys.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Key name。.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Key type。.",
						},
						"encoding": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Key encoding method.",
						},
						"expire_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Key expiration timestamp (in milliseconds), 0 means no expiration time is set.",
						},
						"length": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Key memory size, unit Byte.",
						},
						"item_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of elements.",
						},
						"max_element_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Maximum element length.",
						},
					},
				},
			},

			"timestamp": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Collection timestamp (seconds).",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudDbbrainRedisTopBigKeysRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_dbbrain_redis_top_big_keys.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("date"); ok {
		paramMap["Date"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("product"); ok {
		paramMap["Product"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("sort_by"); ok {
		paramMap["SortBy"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("key_type"); ok {
		paramMap["KeyType"] = helper.String(v.(string))
	}

	service := DbbrainService{client: meta.(*TencentCloudClient).apiV3Conn}

	var topKeys []*dbbrain.RedisKeySpaceData

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDbbrainRedisTopBigKeysByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		topKeys = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(topKeys))
	tmpList := make([]map[string]interface{}, 0, len(topKeys))

	if topKeys != nil {
		for _, redisKeySpaceData := range topKeys {
			redisKeySpaceDataMap := map[string]interface{}{}

			if redisKeySpaceData.Key != nil {
				redisKeySpaceDataMap["key"] = redisKeySpaceData.Key
			}

			if redisKeySpaceData.Type != nil {
				redisKeySpaceDataMap["type"] = redisKeySpaceData.Type
			}

			if redisKeySpaceData.Encoding != nil {
				redisKeySpaceDataMap["encoding"] = redisKeySpaceData.Encoding
			}

			if redisKeySpaceData.ExpireTime != nil {
				redisKeySpaceDataMap["expire_time"] = redisKeySpaceData.ExpireTime
			}

			if redisKeySpaceData.Length != nil {
				redisKeySpaceDataMap["length"] = redisKeySpaceData.Length
			}

			if redisKeySpaceData.ItemCount != nil {
				redisKeySpaceDataMap["item_count"] = redisKeySpaceData.ItemCount
			}

			if redisKeySpaceData.MaxElementSize != nil {
				redisKeySpaceDataMap["max_element_size"] = redisKeySpaceData.MaxElementSize
			}

			ids = append(ids, *redisKeySpaceData.InstanceId)
			tmpList = append(tmpList, redisKeySpaceDataMap)
		}

		_ = d.Set("top_keys", tmpList)
	}

	if timestamp != nil {
		_ = d.Set("timestamp", timestamp)
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
