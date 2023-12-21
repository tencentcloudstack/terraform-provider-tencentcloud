package monitor

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudMonitorAlarmNoticeCallbacks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMonitorAlarmNoticeCallbacksRead,
		Schema: map[string]*schema.Schema{
			"url_notices": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Alarm callback notification.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Callback URL (limited to 256 characters).",
						},
						"is_valid": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Verified 0=No 1=Yes.",
						},
						"validation_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Verification code.",
						},
						"start_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of seconds starting from the day of notification start time.",
						},
						"end_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of seconds from the end of the notification day.",
						},
						"weekday": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
							Computed:    true,
							Description: "Notification period 1-7 represents Monday to Sunday.",
						},
					},
				},
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudMonitorAlarmNoticeCallbacksRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_monitor_alarm_notice_callbacks.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := MonitorService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var urlNotices []*monitor.URLNotice
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMonitorAlarmNoticeCallbacksByFilter(ctx)
		if e != nil {
			return tccommon.RetryError(e)
		}
		urlNotices = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(urlNotices))
	tmpList := make([]map[string]interface{}, 0, len(urlNotices))

	if urlNotices != nil {
		for _, urlNotice := range urlNotices {
			urlNoticeMap := map[string]interface{}{}

			if urlNotice.URL != nil {
				urlNoticeMap["url"] = urlNotice.URL
			}

			if urlNotice.IsValid != nil {
				urlNoticeMap["is_valid"] = urlNotice.IsValid
			}

			if urlNotice.ValidationCode != nil {
				urlNoticeMap["validation_code"] = urlNotice.ValidationCode
			}

			if urlNotice.StartTime != nil {
				urlNoticeMap["start_time"] = urlNotice.StartTime
			}

			if urlNotice.EndTime != nil {
				urlNoticeMap["end_time"] = urlNotice.EndTime
			}

			if urlNotice.Weekday != nil {
				urlNoticeMap["weekday"] = urlNotice.Weekday
			}

			ids = append(ids, *urlNotice.URL)
			tmpList = append(tmpList, urlNoticeMap)
		}

		_ = d.Set("url_notices", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
