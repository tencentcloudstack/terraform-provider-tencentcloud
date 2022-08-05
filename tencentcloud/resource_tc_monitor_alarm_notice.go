/*
Provides a alarm notice resource for monitor.

Example Usage

```hcl
resource "tencentcloud_monitor_alarm_notice" "example" {
  name                  = "yourname"
  notice_type           = "ALL"
  notice_language       = "zh-CN"

}

```

*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func resourceTencentCloudMonitorAlarmNotice() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentMonitorAlarmNoticeCreate,
		Read:   resourceTencentMonitorAlarmNoticeRead,
		Update: resourceTencentMonitorAlarmNoticeUpdate,
		Delete: resourceTencentMonitorAlarmNoticeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Notification template name within 60.",
			},
			"notice_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Alarm notification type ALARM=Notification not restored OK=Notification restored ALL.",
			},
			"notice_language": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Notification language zh-CN=Chinese en-US=English.",
			},

			"notice_ids": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Receive group list.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"alarm_notice": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Alarm notification template list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
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
									"end_time": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The number of seconds since the notification start time 00:00:00 (value range 0-86399).",
									},
									"notice_way": {
										Type:        schema.TypeSet,
										Optional:    true,
										Description: "Notification Channel List EMAIL=Mail SMS=SMS CALL=Telephone WECHAT=WeChat RTX=Enterprise WeChat.",
										Elem:        &schema.Schema{Type: schema.TypeString},
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
						"policy_ids": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "List of alarm policy IDs bound to the alarm notification template.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func resourceTencentMonitorAlarmNoticeCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_alarm_notice.create")()

	var (
		monitorService = MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}
		request        = monitor.NewCreateAlarmNoticeRequest()
	)
	request.Module = helper.String("monitor")
	request.Name = helper.String(d.Get("name").(string))
	request.NoticeType = helper.String(d.Get("notice_type").(string))
	request.NoticeLanguage = helper.String(d.Get("notice_language").(string))

	var noticeId *string
	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, err := monitorService.client.UseMonitorClient().CreateAlarmNotice(request)
		if err != nil {
			return retryError(err, InternalError)
		}
		noticeId = response.Response.NoticeId
		return nil
	}); err != nil {
		return err
	}

	d.SetId(*noticeId)

	return resourceTencentMonitorAlarmNoticeRead(d, meta)
}

func resourceTencentMonitorAlarmNoticeRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_alarm_notice.read")()
	defer inconsistentCheck(d, meta)()

	var (
		monitorService = MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}
		err            error
		alarmNotices   []interface{}
		alarmNotice    []*monitor.AlarmNotice
	)

	alarmNoticeMap := make(map[string]interface{})
	alarmNoticeMap["order"] = helper.String("ASC")
	var tmpAlarmNotice = []*string{helper.String(d.Id())}
	alarmNoticeMap["noticeArr"] = tmpAlarmNotice

	alarmNotice, err = monitorService.DescribeAlarmNoticeById(nil, alarmNoticeMap)
	if err != nil {
		return err
	}
	for _, noticesItem := range alarmNotice {
		noticesItemMap := map[string]interface{}{
			"name":        noticesItem.Name,
			"updated_at":  noticesItem.UpdatedAt,
			"updated_by":  noticesItem.UpdatedBy,
			"notice_type": noticesItem.NoticeType,
			"is_preset":   noticesItem.IsPreset,
			"policy_ids":  noticesItem.PolicyIds,
		}

		userNoticesItems := make([]interface{}, 0, 100)
		for _, userNotices := range noticesItem.UserNotices {
			userNoticesItems = append(userNoticesItems, map[string]interface{}{
				"receiver_type": userNotices.ReceiverType,
				"start_time":    userNotices.StartTime,
				"end_time":      userNotices.EndTime,
				"notice_way":    userNotices.NoticeWay,
			})
		}
		noticesItemMap["user_notices"] = userNoticesItems
		alarmNotices = append(alarmNotices, noticesItemMap)
	}

	if err = d.Set("alarm_notice", alarmNotices); err != nil {
		return err
	}

	return nil
}

func resourceTencentMonitorAlarmNoticeUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_alarm_notice.update")()

	var (
		monitorService = MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}
		request        = monitor.NewModifyAlarmNoticeRequest()
	)

	request.Module = helper.String("monitor")
	request.Name = helper.String(d.Get("name").(string))
	request.NoticeType = helper.String(d.Get("notice_type").(string))
	request.NoticeLanguage = helper.String(d.Get("notice_language").(string))
	request.NoticeId = helper.String(d.Id())

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, err := monitorService.client.UseMonitorClient().ModifyAlarmNotice(request)
		if err != nil {
			return retryError(err, InternalError)
		}
		return nil
	}); err != nil {
		return err
	}

	return resourceTencentMonitorAlarmNoticeRead(d, meta)
}

func resourceTencentMonitorAlarmNoticeDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_alarm_notice.delete")()

	var (
		monitorService = MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}
		request        = monitor.NewDeleteAlarmNoticesRequest()
	)

	request.Module = helper.String("monitor")
	noticeId := d.Id()
	var n = []*string{&noticeId}
	request.NoticeIds = n

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, err := monitorService.client.UseMonitorClient().DeleteAlarmNotices(request)
		if err != nil {
			return retryError(err, InternalError)
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}
