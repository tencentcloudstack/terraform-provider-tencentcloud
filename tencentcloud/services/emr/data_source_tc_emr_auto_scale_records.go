package emr

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	emr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/emr/v20190103"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudEmrAutoScaleRecords() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudEmrAutoScaleRecordsRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "EMR cluster ID.",
			},

			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Record filtering parameters, currently only `StartTime`, `EndTime` and `StrategyName` are supported. `StartTime` and `EndTime` support the time format of 2006-01-02 15:04:05 or 2006/01/02 15:04:05.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Key. Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Value. Note: This field may return null, indicating that no valid value can be obtained.",
						},
					},
				},
			},

			"record_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Record list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"strategy_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Rule name of expanding and shrinking capacity.",
						},
						"scale_action": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "`SCALE_OUT` and `SCALE_IN` respectively represent expanding and shrinking capacity.",
						},
						"action_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "`SUCCESS`, `FAILED`, `PART_SUCCESS`, `IN_PROCESS`.",
						},
						"action_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Process Trigger Time.",
						},
						"scale_info": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Scalability-related Description.",
						},
						"expect_scale_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Effective only when ScaleAction is SCALE_OUT.",
						},
						"end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Process End Time.",
						},
						"strategy_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Strategy Type, 1 for Load scaling, 2 for Time scaling.",
						},
						"spec_info": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specification information used when expanding capacity.",
						},
						"compensate_flag": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Compensation and expansion, 0 represents no start, 1 represents start. Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"compensate_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Compensation Times Note: This field may return null, indicating that no valid value can be obtained.",
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

func dataSourceTencentCloudEmrAutoScaleRecordsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_emr_auto_scale_records.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var instanceId string

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*emr.KeyValue, 0, len(filtersSet))

		for _, item := range filtersSet {
			keyValue := emr.KeyValue{}
			keyValueMap := item.(map[string]interface{})

			if v, ok := keyValueMap["key"]; ok {
				keyValue.Key = helper.String(v.(string))
			}
			if v, ok := keyValueMap["value"]; ok {
				keyValue.Value = helper.String(v.(string))
			}
			tmpSet = append(tmpSet, &keyValue)
		}
		paramMap["Filters"] = tmpSet
	}

	service := EMRService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var recordList []*emr.AutoScaleRecord
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeEmrAutoScaleRecordsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		recordList = result
		return nil
	})
	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0, len(recordList))

	if recordList != nil {
		for _, autoScaleRecord := range recordList {
			autoScaleRecordMap := map[string]interface{}{}

			if autoScaleRecord.StrategyName != nil {
				autoScaleRecordMap["strategy_name"] = autoScaleRecord.StrategyName
			}

			if autoScaleRecord.ScaleAction != nil {
				autoScaleRecordMap["scale_action"] = autoScaleRecord.ScaleAction
			}

			if autoScaleRecord.ActionStatus != nil {
				autoScaleRecordMap["action_status"] = autoScaleRecord.ActionStatus
			}

			if autoScaleRecord.ActionTime != nil {
				autoScaleRecordMap["action_time"] = autoScaleRecord.ActionTime
			}

			if autoScaleRecord.ScaleInfo != nil {
				autoScaleRecordMap["scale_info"] = autoScaleRecord.ScaleInfo
			}

			if autoScaleRecord.ExpectScaleNum != nil {
				autoScaleRecordMap["expect_scale_num"] = autoScaleRecord.ExpectScaleNum
			}

			if autoScaleRecord.EndTime != nil {
				autoScaleRecordMap["end_time"] = autoScaleRecord.EndTime
			}

			if autoScaleRecord.StrategyType != nil {
				autoScaleRecordMap["strategy_type"] = autoScaleRecord.StrategyType
			}

			if autoScaleRecord.SpecInfo != nil {
				autoScaleRecordMap["spec_info"] = autoScaleRecord.SpecInfo
			}

			if autoScaleRecord.CompensateFlag != nil {
				autoScaleRecordMap["compensate_flag"] = autoScaleRecord.CompensateFlag
			}

			if autoScaleRecord.CompensateCount != nil {
				autoScaleRecordMap["compensate_count"] = autoScaleRecord.CompensateCount
			}
			tmpList = append(tmpList, autoScaleRecordMap)
		}

		_ = d.Set("record_list", tmpList)
	}

	d.SetId(instanceId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
