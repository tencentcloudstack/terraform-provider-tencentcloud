package sqlserver

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudSqlserverDescHaLog() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSqlserverDescHaLogRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},
			"start_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Start time (yyyy-MM-dd HH:mm:ss).",
			},
			"end_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "End time (yyyy-MM-dd HH:mm:ss).",
			},
			"switch_type": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Switching mode 0-system automatically switches, 1-manual switch, if not filled in, all will be checked by default.",
			},
			"switch_log": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Master/Slave switching log.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"event_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Switch event ID Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"switch_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Switching mode 0-system automatic switching, 1-manual switching Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Switch start time Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Switch end time Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"reason": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Machine failure causes automatic switching Note: This field may return null, indicating that no valid value can be obtained.",
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

func dataSourceTencentCloudSqlserverDescHaLogRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_sqlserver_desc_ha_log.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = SqlserverService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		switchLogs []*sqlserver.SwitchLog
		instanceId string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("start_time"); ok {
		paramMap["StartTime"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_time"); ok {
		paramMap["EndTime"] = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("switch_type"); ok {
		paramMap["SwitchType"] = helper.IntUint64(v.(int))
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeSqlserverDescHaLogByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		switchLogs = result
		return nil
	})

	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0, len(switchLogs))

	if switchLogs != nil {
		for _, switchLog := range switchLogs {
			switchLogMap := map[string]interface{}{}

			if switchLog.EventId != nil {
				switchLogMap["event_id"] = switchLog.EventId
			}

			if switchLog.SwitchType != nil {
				switchLogMap["switch_type"] = switchLog.SwitchType
			}

			if switchLog.StartTime != nil {
				switchLogMap["start_time"] = switchLog.StartTime
			}

			if switchLog.EndTime != nil {
				switchLogMap["end_time"] = switchLog.EndTime
			}

			if switchLog.Reason != nil {
				switchLogMap["reason"] = switchLog.Reason
			}

			tmpList = append(tmpList, switchLogMap)
		}

		_ = d.Set("switch_log", tmpList)
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
