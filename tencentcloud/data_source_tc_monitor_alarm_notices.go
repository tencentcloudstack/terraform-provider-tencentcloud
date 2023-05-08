/*
Use this data source to Interlude notification list.

Example Usage

```hcl
data "tencentcloud_monitor_alarm_notices" "notices" {
    order = "DESC"
    owner_uid = 1
    name = ""
    receiver_type = ""
    user_ids = []
    group_ids = []
    notice_ids = []
}
```

*/
package tencentcloud

import (
	"context"
	"crypto/md5"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentMonitorAlarmNotices() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentMonitorAlarmNoticesRead,
		Schema: map[string]*schema.Schema{
			"order": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "ASC",
				Description: "Sort by update time ASC=forward order DESC=reverse order.",
			},
			"owner_uid": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The primary account uid is used to create a preset notification.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Alarm notification template name Used for fuzzy search.",
			},
			"receiver_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "To filter alarm notification templates according to recipients, you need to select the notification user type. USER=user GROUP=user group Leave blank = not filter by recipient.",
			},
			"user_ids": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "List of recipients.",
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			"group_ids": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Receive group list.",
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			"notice_ids": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Receive group list.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to store results.",
			},

			"alarm_notice": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Alarm notification template list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Alarm notification template ID.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Alarm notification template name.",
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Last modified time.",
						},
						"updated_by": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Last Modified By.",
						},
						"notice_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Alarm notification type ALARM=Notification not restored OK=Notification restored ALL.",
						},
						"user_notices": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Alarm notification template list.(At most five).",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"receiver_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Recipient Type USER=User GROUP=User Group.",
									},
									"start_time": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The number of seconds since the notification start time 00:00:00 (value range 0-86399).",
									},
									"end_time": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The number of seconds since the notification end time 00:00:00 (value range 0-86399).",
									},
									"notice_way": {
										Type:        schema.TypeSet,
										Computed:    true,
										Description: "Notification Channel List EMAIL=Mail SMS=SMS CALL=Telephone WECHAT=WeChat RTX=Enterprise WeChat.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"user_ids": {
										Type:        schema.TypeSet,
										Computed:    true,
										Description: "User UID List.",
										Elem:        &schema.Schema{Type: schema.TypeInt},
									},
									"group_ids": {
										Type:        schema.TypeSet,
										Computed:    true,
										Description: "User group ID list.",
										Elem:        &schema.Schema{Type: schema.TypeInt},
									},
									"phone_order": {
										Type:        schema.TypeSet,
										Computed:    true,
										Description: "Telephone polling list.",
										Elem:        &schema.Schema{Type: schema.TypeInt},
									},
									"phone_circle_times": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Number of telephone polls (value range: 1-5).",
									},
									"phone_inner_interval": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Number of seconds between calls in a polling session (value range: 60-900).",
									},
									"phone_circle_interval": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Number of seconds between polls (value range: 60-900).",
									},
									"need_phone_arrive_notice": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Contact notification required 0= No 1= Yes.",
									},
									"phone_call_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Call type SYNC= Simultaneous call CIRCLE= Round call If this parameter is not specified, the default value is round call.",
									},
									"weekday": {
										Type:        schema.TypeSet,
										Computed:    true,
										Description: "Notification period 1-7 indicates Monday to Sunday.",
										Elem:        &schema.Schema{Type: schema.TypeInt},
									},
								},
							},
						},
						"url_notices": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The maximum number of callback notifications is 3.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"url": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Callback URL (limited to 256 characters).",
									},
									"start_time": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Notification Start Time Number of seconds at the start of a day.",
									},
									"end_time": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Notification End Time Seconds at the start of a day.",
									},
									"weekday": {
										Type:        schema.TypeSet,
										Computed:    true,
										Description: "Notification period 1-7 indicates Monday to Sunday.",
										Elem:        &schema.Schema{Type: schema.TypeInt},
									},
								},
							},
						},
						"cls_notices": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "A maximum of one alarm notification can be pushed to the CLS service.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"region": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Regional.",
									},
									"log_set_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Log collection Id.",
									},
									"topic_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Theme Id.",
									},
									"enable": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Start-stop status, can not be transmitted, default enabled. 0= Disabled, 1= enabled.",
									},
								},
							},
						},
						"is_preset": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether it is the system default notification template 0=No 1=Yes.",
						},
						"notice_language": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Notification language zh-CN=Chinese en-US=English.",
						},
						"policy_ids": {
							Type:        schema.TypeSet,
							Computed:    true,
							Description: "List of alarm policy IDs bound to the alarm notification template.",
							Elem:        &schema.Schema{Type: schema.TypeString},
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
		err            error
		alarmNotices   []interface{}
		alarmNotice    []*monitor.AlarmNotice
	)

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	alarmNoticeMap := make(map[string]interface{})
	alarmNoticeMap["order"] = helper.String(d.Get("order").(string))

	if v, ok := d.GetOk("owner_uid"); ok {
		alarmNoticeMap["ownerUid"] = helper.IntInt64(v.(int))
	}
	if v, ok := d.GetOk("name"); ok {
		alarmNoticeMap["name"] = helper.String(v.(string))
	}
	if v, ok := d.GetOk("receiver_type"); ok {
		alarmNoticeMap["receiverType"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("user_ids"); ok {
		userIds := v.(*schema.Set).List()
		userIdsArr := make([]*int64, 0, len(userIds))
		for _, userId := range userIds {
			userIdsArr = append(userIdsArr, helper.Int64(userId.(int64)))
		}
		alarmNoticeMap["userIdArr"] = userIdsArr
	}

	if v, ok := d.GetOk("group_ids"); ok {
		groupIds := v.(*schema.Set).List()
		groupIdsArr := make([]*int64, 0, len(groupIds))
		for _, groupId := range groupIds {
			groupIdsArr = append(groupIdsArr, helper.Int64(groupId.(int64)))
		}
		alarmNoticeMap["groupArr"] = groupIdsArr
	}

	if v, ok := d.GetOk("notice_ids"); ok {
		noticeIds := v.(*schema.Set).List()
		noticeIdsArr := make([]*string, 0, len(noticeIds))
		for _, noticeId := range noticeIds {
			noticeIdsArr = append(noticeIdsArr, helper.String(noticeId.(string)))
		}
		alarmNoticeMap["noticeArr"] = noticeIdsArr
	}

	alarmNotice, err = monitorService.DescribeAlarmNoticeById(ctx, alarmNoticeMap)
	if err != nil {
		return err
	}
	for _, noticesItem := range alarmNotice {
		noticesItemMap := map[string]interface{}{
			"id":              noticesItem.Id,
			"name":            noticesItem.Name,
			"updated_at":      noticesItem.UpdatedAt,
			"updated_by":      noticesItem.UpdatedBy,
			"notice_type":     noticesItem.NoticeType,
			"is_preset":       noticesItem.IsPreset,
			"notice_language": noticesItem.NoticeLanguage,
			"policy_ids":      noticesItem.PolicyIds,
		}

		userNoticesItems := make([]interface{}, 0, 100)
		for _, userNotices := range noticesItem.UserNotices {
			userNoticesItems = append(userNoticesItems, map[string]interface{}{
				"receiver_type":            userNotices.ReceiverType,
				"start_time":               userNotices.StartTime,
				"end_time":                 userNotices.EndTime,
				"notice_way":               userNotices.NoticeWay,
				"user_ids":                 userNotices.UserIds,
				"group_ids":                userNotices.GroupIds,
				"phone_order":              userNotices.PhoneOrder,
				"phone_circle_times":       userNotices.PhoneCircleTimes,
				"phone_inner_interval":     userNotices.PhoneInnerInterval,
				"phone_circle_interval":    userNotices.PhoneCircleInterval,
				"need_phone_arrive_notice": userNotices.NeedPhoneArriveNotice,
				"phone_call_type":          userNotices.PhoneCallType,
				"weekday":                  userNotices.Weekday,
			})
		}

		urlNoticesItems := make([]interface{}, 0, 100)
		for _, urlNotice := range noticesItem.URLNotices {
			urlNoticesItems = append(urlNoticesItems, map[string]interface{}{
				"url":        urlNotice.URL,
				"start_time": urlNotice.StartTime,
				"end_time":   urlNotice.EndTime,
				"weekday":    urlNotice.Weekday,
			})
		}

		clsNoticesItems := make([]interface{}, 0, 100)
		for _, clsNotice := range noticesItem.CLSNotices {
			clsNoticesItems = append(clsNoticesItems, map[string]interface{}{
				"region":     clsNotice.Region,
				"log_set_id": clsNotice.LogSetId,
				"topic_id":   clsNotice.TopicId,
				"enable":     clsNotice.Enable,
			})
		}
		noticesItemMap["user_notices"] = userNoticesItems
		noticesItemMap["url_notices"] = urlNoticesItems
		noticesItemMap["cls_notices"] = clsNoticesItems
		alarmNotices = append(alarmNotices, noticesItemMap)
	}

	md := md5.New()
	id := fmt.Sprintf("%x", md.Sum(nil))
	d.SetId(id)

	if err = d.Set("alarm_notice", alarmNotices); err != nil {
		return err
	}
	if output, ok := d.GetOk("result_output_file"); ok {
		return writeToFile(output.(string), alarmNotices)
	}
	return nil
}
