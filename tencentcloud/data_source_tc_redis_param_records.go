/*
Use this data source to query detailed information of redis param_records

Example Usage

```hcl
data "tencentcloud_redis_param_records" "param_records" {
  instance_id = "crs-c1nl9rpv"
  limit = &lt;nil&gt;
  offset = &lt;nil&gt;
  }
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudRedisParamRecords() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudRedisParamRecordsRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The ID of instance.",
			},

			"limit": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Page size.",
			},

			"offset": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Page offset.",
			},

			"instance_param_history": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "The parameter name.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"param_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The parameter name.",
						},
						"pre_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Modify the previous value.",
						},
						"new_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The modified value.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Parameter status:1: parameter configuration modification.2: The parameter configuration is modified successfully.3: Parameter configuration modification failed.",
						},
						"modify_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Modification time.",
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

func dataSourceTencentCloudRedisParamRecordsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_redis_param_records.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("limit"); v != nil {
		paramMap["Limit"] = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOk("offset"); v != nil {
		paramMap["Offset"] = helper.IntUint64(v.(int))
	}

	service := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}

	var instanceParamHistory []*redis.InstanceParamHistory

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeRedisParamRecordsByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		instanceParamHistory = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(instanceParamHistory))
	tmpList := make([]map[string]interface{}, 0, len(instanceParamHistory))

	if instanceParamHistory != nil {
		for _, instanceParamHistory := range instanceParamHistory {
			instanceParamHistoryMap := map[string]interface{}{}

			if instanceParamHistory.ParamName != nil {
				instanceParamHistoryMap["param_name"] = instanceParamHistory.ParamName
			}

			if instanceParamHistory.PreValue != nil {
				instanceParamHistoryMap["pre_value"] = instanceParamHistory.PreValue
			}

			if instanceParamHistory.NewValue != nil {
				instanceParamHistoryMap["new_value"] = instanceParamHistory.NewValue
			}

			if instanceParamHistory.Status != nil {
				instanceParamHistoryMap["status"] = instanceParamHistory.Status
			}

			if instanceParamHistory.ModifyTime != nil {
				instanceParamHistoryMap["modify_time"] = instanceParamHistory.ModifyTime
			}

			ids = append(ids, *instanceParamHistory.InstanceId)
			tmpList = append(tmpList, instanceParamHistoryMap)
		}

		_ = d.Set("instance_param_history", tmpList)
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
