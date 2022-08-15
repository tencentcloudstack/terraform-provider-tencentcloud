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

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
							Optional:    true,
							Description: "Alarm notification template ID.",
						},
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
						"notice_language": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Notification language zh-CN=Chinese en-US=English.",
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
				"receiver_type": userNotices.ReceiverType,
				"start_time":    userNotices.StartTime,
				"end_time":      userNotices.EndTime,
				"notice_way":    userNotices.NoticeWay,
			})
		}
		noticesItemMap["user_notices"] = userNoticesItems
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
