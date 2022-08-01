/*
Provides a alarm notice resource for monitor.

Example Usage

```hcl
resource "tencentcloud_monitor_alarm_notice" "example" {
  module                = "monitor"
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
			"module": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Module name, fill in 'monitor' here.",
			},
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

						"notice_way": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Notification Channel List EMAIL=Mail SMS=SMS CALL=Telephone WECHAT=WeChat RTX=Enterprise WeChat.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"endtime": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The number of seconds since the notification start time 00:00:00 (value range 0-86399).",
									},
								},
							},
						},
					},
				},
			},

			"notice_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Alarm notification template ID.",
			},

			"request_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique request ID, returned on every request. When locating the problem, you need to provide the RequestId of the request.",
			},

			"notice_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of notification rule IDs.",
				Elem: &schema.Schema{
					Type:        schema.TypeString,
					Description: "ID of the notification rule to be queried.",
				},
			},

			"notices": {
				Type:        schema.TypeList,
				Computed:    true,
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

func resourceTencentMonitorAlarmNoticeCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_alarm_notice.create")()

	var (
		monitorService = MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}
		request        = monitor.NewCreateAlarmNoticeRequest()
		err            error
		//noticeids      []interface{}
	)
	request.Module = helper.String("monitor")
	request.Name = helper.String(d.Get("name").(string))
	request.NoticeType = helper.String(d.Get("notice_type").(string))
	request.NoticeLanguage = helper.String(d.Get("notice_language").(string))

	var noticeid *string
	var requestid *string
	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, err := monitorService.client.UseMonitorClient().CreateAlarmNotice(request)
		if err != nil {
			return retryError(err, InternalError)
		}
		noticeid = response.Response.NoticeId
		requestid = response.Response.RequestId
		return nil
	}); err != nil {
		return err
	}

	d.SetId(*noticeid)
	if err = d.Set("request_id", requestid); err != nil {
		return err
	}

	if err = d.Set("notice_id", noticeid); err != nil {
		return err
	}

	//noticeids = append(noticeids, noticeid)
	//
	//if err = d.Set("notice_ids", noticeids); err != nil {
	//	return err
	//}

	return resourceTencentMonitorAlarmNoticeRead(d, meta)
}

func resourceTencentMonitorAlarmNoticeRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_alarm_notice.read")()
	defer inconsistentCheck(d, meta)()

	var (
		monitorService = MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}
		request        = monitor.NewDescribeAlarmNoticeRequest()
		notice         []interface{}
		err            error
	)

	request.Module = helper.String("monitor")
	noticeId := d.Id()
	request.NoticeId = &noticeId

	var requestid *string
	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, err := monitorService.client.UseMonitorClient().DescribeAlarmNotice(request)
		if err != nil {
			return retryError(err, InternalError)
		}
		noticeItem := response.Response.Notice
		noticesItemMap := map[string]interface{}{
			"notices_id":      &noticeId,
			"notices_name":    noticeItem.Name,
			"updated_at":      noticeItem.UpdatedAt,
			"updated_by":      noticeItem.UpdatedBy,
			"notice_type":     noticeItem.NoticeType,
			"is_preset":       noticeItem.IsPreset,
			"notice_language": noticeItem.NoticeLanguage,
		}

		userNoticesItems := make([]interface{}, 0, 100)
		for _, userNoticesItem := range noticeItem.UserNotices {
			userNoticesItems = append(userNoticesItems, map[string]interface{}{
				"receiver_type": userNoticesItem.ReceiverType,
				"start_time":    userNoticesItem.StartTime,
				"endtime":       userNoticesItem.EndTime,
			})
		}
		noticesItemMap["user_notices"] = userNoticesItems
		notice = append(notice, noticesItemMap)
		requestid = response.Response.RequestId

		return nil
	}); err != nil {
		return err
	}

	d.SetId(noticeId)

	if err = d.Set("notices", notice); err != nil {
		return err
	}
	if err = d.Set("request_id", requestid); err != nil {
		return err
	}

	return nil
}

func resourceTencentMonitorAlarmNoticeUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_alarm_notice.update")()

	var (
		monitorService = MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}
		request        = monitor.NewModifyAlarmNoticeRequest()
		err            error
		//noticeids      []interface{}
	)

	request.Module = helper.String("monitor")
	request.Name = helper.String(d.Get("name").(string))
	request.NoticeType = helper.String(d.Get("notice_type").(string))
	request.NoticeLanguage = helper.String(d.Get("notice_language").(string))
	noticeId := d.Id()
	request.NoticeId = &noticeId

	var requestid *string
	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, err := monitorService.client.UseMonitorClient().ModifyAlarmNotice(request)
		if err != nil {
			return retryError(err, InternalError)
		}
		requestid = response.Response.RequestId

		return nil
	}); err != nil {
		return err
	}

	if err = d.Set("request_id", requestid); err != nil {
		return err
	}

	//noticeids = append(noticeids, &noticeId)
	//
	//if err = d.Set("notice_ids", noticeids); err != nil {
	//	return err
	//}

	return resourceTencentMonitorAlarmNoticeRead(d, meta)
}

func resourceTencentMonitorAlarmNoticeDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_alarm_notice.delete")()

	var (
		monitorService = MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}
		request        = monitor.NewDeleteAlarmNoticesRequest()
		err            error
	)

	request.Module = helper.String("monitor")
	noticeId := d.Id()
	var n = []*string{&noticeId}
	request.NoticeIds = n

	//if v, ok := d.GetOk("notice_ids"); ok {
	//	notice := make([]*string, 0, 10)
	//	for _, item := range v.([]interface{}) {
	//		notice = append(notice, helper.String(item.(string)))
	//	}
	//	request.NoticeIds = append(notice, nil)
	//}

	var requestid *string
	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, err := monitorService.client.UseMonitorClient().DeleteAlarmNotices(request)
		if err != nil {
			return retryError(err, InternalError)
		}
		requestid = response.Response.RequestId

		return nil
	}); err != nil {
		return err
	}

	if err = d.Set("request_id", requestid); err != nil {
		return err
	}
	return nil
}
