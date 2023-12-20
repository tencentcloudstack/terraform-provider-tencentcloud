package crs

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudRedisRecordsParam() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudRedisParamRecordsRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The ID of instance.",
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
	defer tccommon.LogElapsed("data_source.tencentcloud_redis_param_records.read")()
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

	var instanceParamHistory []*redis.InstanceParamHistory

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeRedisParamRecordsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		instanceParamHistory = result
		return nil
	})
	if err != nil {
		return err
	}

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

			tmpList = append(tmpList, instanceParamHistoryMap)
		}

		_ = d.Set("instance_param_history", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash([]string{instanceId}))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
