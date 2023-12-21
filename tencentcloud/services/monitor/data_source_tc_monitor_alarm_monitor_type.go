package monitor

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudMonitorAlarmMonitorType() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMonitorAlarmMonitorTypeRead,
		Schema: map[string]*schema.Schema{
			"monitor_types": {
				Computed: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Monitoring type, cloud product monitoring is MT_ QCE.",
			},

			"monitor_type_infos": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Monitoring type details.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Monitoring type ID.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Monitoring type.",
						},
						"sort_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Sort order.",
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

func dataSourceTencentCloudMonitorAlarmMonitorTypeRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_monitor_alarm_monitor_type.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	paramMap := make(map[string]interface{})
	service := MonitorService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var monitor *monitor.DescribeMonitorTypesResponseParams
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMonitorAlarmMonitorTypeByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		monitor = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0)
	if monitor.MonitorTypes != nil {
		_ = d.Set("monitor_types", monitor.MonitorTypes)

		for _, v := range monitor.MonitorTypes {
			ids = append(ids, *v)
		}
	}

	if monitor.MonitorTypeInfos != nil {
		tmpList := make([]map[string]interface{}, 0, len(monitor.MonitorTypeInfos))
		for _, monitorTypeInfo := range monitor.MonitorTypeInfos {
			monitorTypeInfoMap := map[string]interface{}{}

			if monitorTypeInfo.Id != nil {
				monitorTypeInfoMap["id"] = monitorTypeInfo.Id
			}

			if monitorTypeInfo.Name != nil {
				monitorTypeInfoMap["name"] = monitorTypeInfo.Name
			}

			if monitorTypeInfo.SortId != nil {
				monitorTypeInfoMap["sort_id"] = monitorTypeInfo.SortId
			}

			ids = append(ids, *monitorTypeInfo.Id)
			tmpList = append(tmpList, monitorTypeInfoMap)
		}

		_ = d.Set("monitor_type_infos", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}
	return nil
}
