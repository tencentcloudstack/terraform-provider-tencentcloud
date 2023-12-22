package monitor

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func ResourceTencentCloudMonitorAlarmNotice() *schema.Resource {
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
				Required:    true,
				Description: "Alarm notification type ALARM=Notification not restored OK=Notification restored ALL.",
			},
			"notice_language": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Notification language zh-CN=Chinese en-US=English.",
			},
			"user_notices": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Alarm notification template list.(At most five).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"receiver_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Recipient Type USER=User GROUP=User Group.",
						},
						"start_time": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The number of seconds since the notification start time 00:00:00 (value range 0-86399).",
						},
						"end_time": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The number of seconds since the notification end time 00:00:00 (value range 0-86399).",
						},
						"notice_way": {
							Type:        schema.TypeSet,
							Required:    true,
							Description: "Notification Channel List EMAIL=Mail SMS=SMS CALL=Telephone WECHAT=WeChat RTX=Enterprise WeChat.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"user_ids": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "User UID List.",
							Elem:        &schema.Schema{Type: schema.TypeInt},
						},
						"group_ids": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "User group ID list.",
							Elem:        &schema.Schema{Type: schema.TypeInt},
						},
						"phone_order": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "Telephone polling list.",
							Elem:        &schema.Schema{Type: schema.TypeInt},
						},
						"phone_circle_times": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Number of telephone polls (value range: 1-5).",
						},
						"phone_inner_interval": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Number of seconds between calls in a polling session (value range: 60-900).",
						},
						"phone_circle_interval": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Number of seconds between polls (value range: 60-900).",
						},
						"need_phone_arrive_notice": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Contact notification required 0= No 1= Yes.",
						},
						"phone_call_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Call type SYNC= Simultaneous call CIRCLE= Round call If this parameter is not specified, the default value is round call.",
						},
						"weekday": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "Notification period 1-7 indicates Monday to Sunday.",
							Elem:        &schema.Schema{Type: schema.TypeInt},
						},
					},
				},
			},
			"url_notices": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The maximum number of callback notifications is 3.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"url": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Callback URL (limited to 256 characters).",
						},
						"is_valid": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "If passed verification `0` is no, `1` is yes. Default `0`.",
						},
						"validation_code": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Verification code.",
						},
						"start_time": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Notification Start Time Number of seconds at the start of a day.",
						},
						"end_time": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Notification End Time Seconds at the start of a day.",
						},
						"weekday": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "Notification period 1-7 indicates Monday to Sunday.",
							Elem:        &schema.Schema{Type: schema.TypeInt},
						},
					},
				},
			},
			"cls_notices": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A maximum of one alarm notification can be pushed to the CLS service.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Regional.",
						},
						"log_set_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Log collection Id.",
						},
						"topic_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Theme Id.",
						},
						"enable": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     1,
							Description: "Start-stop status, can not be transmitted, default enabled. 0= Disabled, 1= enabled.",
						},
					},
				},
			},
			"updated_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last Modified By.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last modified time.",
			},
			"is_preset": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Whether it is the system default notification template 0=No 1=Yes.",
			},
			"amp_consumer_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Amp consumer ID.",
			},
			"policy_ids": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "List of alarm policy IDs bound to the alarm notification template.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceTencentMonitorAlarmNoticeCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_alarm_notice.create")()

	var (
		monitorService = MonitorService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request        = monitor.NewCreateAlarmNoticeRequest()
	)
	request.Module = helper.String("monitor")
	request.Name = helper.String(d.Get("name").(string))
	request.NoticeType = helper.String(d.Get("notice_type").(string))
	request.NoticeLanguage = helper.String(d.Get("notice_language").(string))

	if v, ok := d.GetOk("user_notices"); ok {
		userNotices := make([]*monitor.UserNotice, 0, 10)
		for _, item := range v.([]interface{}) {
			m := item.(map[string]interface{})
			userNotice := monitor.UserNotice{}
			userNotice.ReceiverType = helper.String(m["receiver_type"].(string))
			userNotice.StartTime = helper.IntInt64(m["start_time"].(int))
			userNotice.EndTime = helper.IntInt64(m["end_time"].(int))

			if v, ok := m["notice_way"]; ok {
				noticeWay := v.(*schema.Set).List()
				noticeWayArr := make([]*string, 0, len(noticeWay))
				for _, noticeId := range noticeWay {
					noticeWayArr = append(noticeWayArr, helper.String(noticeId.(string)))
				}
				userNotice.NoticeWay = noticeWayArr
			}

			if v, ok := m["user_ids"]; ok {
				userIds := v.(*schema.Set).List()
				userIdsArr := make([]*int64, 0, len(userIds))
				for _, userId := range userIds {
					userIdsArr = append(userIdsArr, helper.IntInt64(userId.(int)))
				}
				userNotice.UserIds = userIdsArr
			}

			if v, ok := m["group_ids"]; ok {
				groupIds := v.(*schema.Set).List()
				groupIdsArr := make([]*int64, 0, len(groupIds))
				for _, groupId := range groupIds {
					groupIdsArr = append(groupIdsArr, helper.IntInt64(groupId.(int)))
				}
				userNotice.GroupIds = groupIdsArr
			}

			if v, ok := m["phone_order"]; ok {
				phoneOrder := v.(*schema.Set).List()
				phoneOrderArr := make([]*int64, 0, len(phoneOrder))
				for _, phone := range phoneOrder {
					phoneOrderArr = append(phoneOrderArr, helper.IntInt64(phone.(int)))
				}
				userNotice.PhoneOrder = phoneOrderArr
			}

			if m["phone_circle_times"] != nil {
				userNotice.PhoneCircleTimes = helper.IntInt64(m["phone_circle_times"].(int))
			}

			if m["phone_inner_interval"] != nil {
				userNotice.PhoneInnerInterval = helper.IntInt64(m["phone_inner_interval"].(int))
			}

			if m["phone_circle_interval"] != nil {
				userNotice.PhoneCircleInterval = helper.IntInt64(m["phone_circle_interval"].(int))
			}

			if m["need_phone_arrive_notice"] != nil {
				userNotice.NeedPhoneArriveNotice = helper.IntInt64(m["need_phone_arrive_notice"].(int))
			}

			if m["phone_call_type"] != nil {
				userNotice.PhoneCallType = helper.String(m["phone_call_type"].(string))
			}

			if v, ok := m["weekday"]; ok {
				weekday := v.(*schema.Set).List()
				weekdayArr := make([]*int64, 0, len(weekday))
				for _, week := range weekday {
					weekdayArr = append(weekdayArr, helper.IntInt64(week.(int)))
				}
				userNotice.Weekday = weekdayArr
			}
			userNotices = append(userNotices, &userNotice)
		}
		request.UserNotices = userNotices
	}

	if v, ok := d.GetOk("url_notices"); ok {
		urlNotices := make([]*monitor.URLNotice, 0, 10)
		for _, item := range v.([]interface{}) {
			m := item.(map[string]interface{})
			urlNotice := monitor.URLNotice{}
			urlNotice.URL = helper.String(m["url"].(string))

			if m["is_valid"] != nil {
				urlNotice.IsValid = helper.IntInt64(m["is_valid"].(int))
			}

			if m["validation_code"] != "" {
				urlNotice.ValidationCode = helper.String(m["validation_code"].(string))
			}

			if m["start_time"] != nil {
				urlNotice.StartTime = helper.IntInt64(m["start_time"].(int))
			}

			if m["end_time"] != nil {
				urlNotice.EndTime = helper.IntInt64(m["end_time"].(int))
			}

			if v, ok := m["weekday"]; ok {
				weekday := v.(*schema.Set).List()
				weekdayArr := make([]*int64, 0, len(weekday))
				for _, week := range weekday {
					weekdayArr = append(weekdayArr, helper.IntInt64(week.(int)))
				}
				urlNotice.Weekday = weekdayArr
			}
			urlNotices = append(urlNotices, &urlNotice)
		}
		request.URLNotices = urlNotices
	}

	if v, ok := d.GetOk("cls_notices"); ok {
		clsNotices := make([]*monitor.CLSNotice, 0, 10)
		for _, item := range v.([]interface{}) {
			m := item.(map[string]interface{})
			clsNotice := monitor.CLSNotice{}
			clsNotice.Region = helper.String(m["region"].(string))
			clsNotice.LogSetId = helper.String(m["log_set_id"].(string))
			clsNotice.TopicId = helper.String(m["topic_id"].(string))

			if m["enable"] != nil {
				clsNotice.Enable = helper.IntInt64(m["enable"].(int))
			}
			clsNotices = append(clsNotices, &clsNotice)
		}
		request.CLSNotices = clsNotices
	}

	var noticeId *string
	if err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, err := monitorService.client.UseMonitorClient().CreateAlarmNotice(request)
		if err != nil {
			return tccommon.RetryError(err, tccommon.InternalError)
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
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_alarm_notice.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var (
		monitorService = MonitorService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		err            error
		alarmNotice    []*monitor.AlarmNotice
	)

	alarmNoticeMap := make(map[string]interface{})
	alarmNoticeMap["order"] = helper.String("ASC")
	var tmpAlarmNotice = []*string{helper.String(d.Id())}
	alarmNoticeMap["noticeArr"] = tmpAlarmNotice

	alarmNotice, err = monitorService.DescribeAlarmNoticeById(ctx, alarmNoticeMap)
	if err != nil {
		return err
	}
	for _, noticesItem := range alarmNotice {
		if err = d.Set("name", noticesItem.Name); err != nil {
			return err
		}
		if err = d.Set("notice_type", noticesItem.NoticeType); err != nil {
			return err
		}
		if err = d.Set("notice_language", noticesItem.NoticeLanguage); err != nil {
			return err
		}
		if err = d.Set("updated_by", noticesItem.UpdatedBy); err != nil {
			return err
		}
		if err = d.Set("updated_at", noticesItem.UpdatedAt); err != nil {
			return err
		}
		if err = d.Set("is_preset", noticesItem.IsPreset); err != nil {
			return err
		}
		if err = d.Set("policy_ids", noticesItem.PolicyIds); err != nil {
			return err
		}
		if err = d.Set("amp_consumer_id", noticesItem.AMPConsumerId); err != nil {
			return err
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
				"url":             urlNotice.URL,
				"is_valid":        urlNotice.IsValid,
				"validation_code": urlNotice.ValidationCode,
				"start_time":      urlNotice.StartTime,
				"end_time":        urlNotice.EndTime,
				"weekday":         urlNotice.Weekday,
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

		if err = d.Set("user_notices", userNoticesItems); err != nil {
			return err
		}
		if err = d.Set("url_notices", urlNoticesItems); err != nil {
			return err
		}
		if err = d.Set("cls_notices", clsNoticesItems); err != nil {
			return err
		}

	}

	return nil
}

func resourceTencentMonitorAlarmNoticeUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_alarm_notice.update")()

	var (
		monitorService = MonitorService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request        = monitor.NewModifyAlarmNoticeRequest()
	)

	request.Module = helper.String("monitor")
	request.Name = helper.String(d.Get("name").(string))
	request.NoticeType = helper.String(d.Get("notice_type").(string))
	request.NoticeLanguage = helper.String(d.Get("notice_language").(string))
	request.NoticeId = helper.String(d.Id())

	if v, ok := d.GetOk("user_notices"); ok {
		userNotices := make([]*monitor.UserNotice, 0, 10)
		for _, item := range v.([]interface{}) {
			m := item.(map[string]interface{})
			userNotice := monitor.UserNotice{}
			userNotice.ReceiverType = helper.String(m["receiver_type"].(string))
			userNotice.StartTime = helper.IntInt64(m["start_time"].(int))
			userNotice.EndTime = helper.IntInt64(m["end_time"].(int))

			if v, ok := m["notice_way"]; ok {
				noticeWay := v.(*schema.Set).List()
				noticeWayArr := make([]*string, 0, len(noticeWay))
				for _, noticeId := range noticeWay {
					noticeWayArr = append(noticeWayArr, helper.String(noticeId.(string)))
				}
				userNotice.NoticeWay = noticeWayArr
			}

			if v, ok := m["user_ids"]; ok {
				userIds := v.(*schema.Set).List()
				userIdsArr := make([]*int64, 0, len(userIds))
				for _, userId := range userIds {
					userIdsArr = append(userIdsArr, helper.IntInt64(userId.(int)))
				}
				userNotice.UserIds = userIdsArr
			}

			if v, ok := m["group_ids"]; ok {
				groupIds := v.(*schema.Set).List()
				groupIdsArr := make([]*int64, 0, len(groupIds))
				for _, groupId := range groupIds {
					groupIdsArr = append(groupIdsArr, helper.IntInt64(groupId.(int)))
				}
				userNotice.GroupIds = groupIdsArr
			}

			if v, ok := m["phone_order"]; ok {
				phoneOrder := v.(*schema.Set).List()
				phoneOrderArr := make([]*int64, 0, len(phoneOrder))
				for _, phone := range phoneOrder {
					phoneOrderArr = append(phoneOrderArr, helper.IntInt64(phone.(int)))
				}
				userNotice.PhoneOrder = phoneOrderArr
			}

			if m["phone_circle_times"] != nil {
				userNotice.PhoneCircleTimes = helper.IntInt64(m["phone_circle_times"].(int))
			}

			if m["phone_inner_interval"] != nil {
				userNotice.PhoneInnerInterval = helper.IntInt64(m["phone_inner_interval"].(int))
			}

			if m["phone_circle_interval"] != nil {
				userNotice.PhoneCircleInterval = helper.IntInt64(m["phone_circle_interval"].(int))
			}

			if m["need_phone_arrive_notice"] != nil {
				userNotice.NeedPhoneArriveNotice = helper.IntInt64(m["need_phone_arrive_notice"].(int))
			}

			if m["phone_call_type"] != nil {
				userNotice.PhoneCallType = helper.String(m["phone_call_type"].(string))
			}

			if v, ok := m["weekday"]; ok {
				weekday := v.(*schema.Set).List()
				weekdayArr := make([]*int64, 0, len(weekday))
				for _, week := range weekday {
					weekdayArr = append(weekdayArr, helper.IntInt64(week.(int)))
				}
				userNotice.Weekday = weekdayArr
			}
			userNotices = append(userNotices, &userNotice)
		}
		request.UserNotices = userNotices
	}

	if v, ok := d.GetOk("url_notices"); ok {
		urlNotices := make([]*monitor.URLNotice, 0, 10)
		for _, item := range v.([]interface{}) {
			m := item.(map[string]interface{})
			urlNotice := monitor.URLNotice{}
			urlNotice.URL = helper.String(m["url"].(string))

			if m["is_valid"] != nil {
				urlNotice.IsValid = helper.IntInt64(m["is_valid"].(int))
			}

			if m["validation_code"] != "" {
				urlNotice.ValidationCode = helper.String(m["validation_code"].(string))
			}

			if m["start_time"] != nil {
				urlNotice.StartTime = helper.IntInt64(m["start_time"].(int))
			}

			if m["end_time"] != nil {
				urlNotice.EndTime = helper.IntInt64(m["end_time"].(int))
			}

			if v, ok := m["weekday"]; ok {
				weekday := v.(*schema.Set).List()
				weekdayArr := make([]*int64, 0, len(weekday))
				for _, week := range weekday {
					weekdayArr = append(weekdayArr, helper.IntInt64(week.(int)))
				}
				urlNotice.Weekday = weekdayArr
			}
			urlNotices = append(urlNotices, &urlNotice)
		}
		request.URLNotices = urlNotices
	}

	if v, ok := d.GetOk("cls_notices"); ok {
		clsNotices := make([]*monitor.CLSNotice, 0, 10)
		for _, item := range v.([]interface{}) {
			m := item.(map[string]interface{})
			clsNotice := monitor.CLSNotice{}
			clsNotice.Region = helper.String(m["region"].(string))
			clsNotice.LogSetId = helper.String(m["log_set_id"].(string))
			clsNotice.TopicId = helper.String(m["topic_id"].(string))

			if m["enable"] != nil {
				clsNotice.Enable = helper.IntInt64(m["enable"].(int))
			}
			clsNotices = append(clsNotices, &clsNotice)
		}
		request.CLSNotices = clsNotices
	}

	if err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, err := monitorService.client.UseMonitorClient().ModifyAlarmNotice(request)
		if err != nil {
			return tccommon.RetryError(err, tccommon.InternalError)
		}
		return nil
	}); err != nil {
		return err
	}

	d.SetId(d.Id())

	return resourceTencentMonitorAlarmNoticeRead(d, meta)
}

func resourceTencentMonitorAlarmNoticeDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_alarm_notice.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var (
		monitorService = MonitorService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	err := monitorService.DeleteMonitorAlarmNoticeById(ctx, d.Id())
	if err != nil {
		return err
	}
	return nil
}
