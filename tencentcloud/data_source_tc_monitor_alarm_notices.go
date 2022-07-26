/*
Use this data source to Interlude notification list.

Example Usage

```hcl
data "tencentcloud_monitor_alarm_notices" "notices" {
  module     = "monitor"
  pagenumber = 1
  pagesize   = 20
  order      = "DESC"

}
```

*/
package tencentcloud

import (
	"crypto/md5"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func dataSourceTencentMonitorAlarmNotices() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentMonitorAlarmNoticesRead,
		Schema: map[string]*schema.Schema{
			"module": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Module name, fill in 'monitor' here.",
			},
			"pagenumber": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Page number minimum 1.",
			},
			"pagesize": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Page size 1-200.",
			},
			"order": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Sort by update time ASC=forward order DESC=reverse order.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to store results.",
			},

			"notices": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Alarm notification template list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"notices_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Alarm notification template ID.",
						},
						"notices_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Alarm notification template name.",
						},
						"updated_at": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Last modified time.",
						},
						"updated_by": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Last Modified By.",
						},
						"notice_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Alarm notification type ALARM=Notification not restored OK=Notification restored ALL.",
						},
						"user_notices": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Alarm notification template list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"receiver_type": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Recipient Type USER=User GROUP=User Group.",
									},
									"start_time": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The number of seconds since the notification start time 00:00:00 (value range 0-86399).",
									},
									"endtime": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The number of seconds since the notification start time 00:00:00 (value range 0-86399).",
									},
								},
							},
						},
						"is_preset": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     1,
							Description: "Whether it is the system default notification template 0=No 1=Yes.",
						},
						"notice_language": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Notification language zh-CN=Chinese en-US=English.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentMonitorAlarmNoticesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_monitor_alarm_notices.read")()

	var (
		monitorService = MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}
		request        = monitor.NewDescribeAlarmNoticesRequest()
		response       *monitor.DescribeAlarmNoticesResponse
		err            error
		notices        []interface{}
	)
	request.Module = helper.String(d.Get("module").(string))
	request.PageNumber = helper.IntInt64(d.Get("pagenumber").(int))
	request.PageSize = helper.IntInt64(d.Get("pagesize").(int))
	request.Order = helper.String(d.Get("order").(string))

	if err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		if response, err = monitorService.client.UseMonitorClient().DescribeAlarmNotices(request); err != nil {
			return retryError(err, InternalError)
		}
		return nil
	}); err != nil {
		return err
	}

	for _, noticesItem := range response.Response.Notices {
		noticesItemMap := map[string]interface{}{
			"notices_id":      noticesItem.Id,
			"notices_name":    noticesItem.Name,
			"updated_at":      noticesItem.UpdatedAt,
			"updated_by":      noticesItem.UpdatedBy,
			"notice_type":     noticesItem.NoticeType,
			"is_preset":       noticesItem.IsPreset,
			"notice_language": noticesItem.NoticeLanguage,
		}

		user_noticesItems := make([]interface{}, 0, 100)
		for _, user_noticesItem := range noticesItem.UserNotices {
			user_noticesItems = append(user_noticesItems, map[string]interface{}{
				"receiver_type": user_noticesItem.ReceiverType,
				"start_time":    user_noticesItem.StartTime,
				"endtime":       user_noticesItem.EndTime,
			})
		}
		noticesItemMap["user_notices"] = user_noticesItems
		notices = append(notices, noticesItemMap)
	}

	md := md5.New()
	_, _ = md.Write([]byte(request.ToJsonString()))
	id := fmt.Sprintf("%x", md.Sum(nil))
	d.SetId(id)

	if err = d.Set("notices", notices); err != nil {
		return err
	}
	if output, ok := d.GetOk("result_output_file"); ok {
		return writeToFile(output.(string), notices)
	}
	return nil
}
