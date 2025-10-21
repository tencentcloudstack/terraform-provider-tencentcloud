---
subcategory: "Cloud Monitor(Monitor)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_alarm_notices"
sidebar_current: "docs-tencentcloud-datasource-monitor_alarm_notices"
description: |-
  Use this data source to Interlude notification list.
---

# tencentcloud_monitor_alarm_notices

Use this data source to Interlude notification list.

## Example Usage

```hcl
data "tencentcloud_monitor_alarm_notices" "notices" {
  order         = "DESC"
  owner_uid     = 1
  name          = ""
  receiver_type = ""
  user_ids      = []
  group_ids     = []
  notice_ids    = []
}
```

## Argument Reference

The following arguments are supported:

* `group_ids` - (Optional, Set: [`Int`]) Receive group list.
* `name` - (Optional, String) Alarm notification template name Used for fuzzy search.
* `notice_ids` - (Optional, Set: [`String`]) Receive group list.
* `order` - (Optional, String) Sort by update time ASC=forward order DESC=reverse order.
* `owner_uid` - (Optional, Int) The primary account uid is used to create a preset notification.
* `receiver_type` - (Optional, String) To filter alarm notification templates according to recipients, you need to select the notification user type. USER=user GROUP=user group Leave blank = not filter by recipient.
* `result_output_file` - (Optional, String) Used to store results.
* `user_ids` - (Optional, Set: [`Int`]) List of recipients.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `alarm_notice` - Alarm notification template list.
  * `amp_consumer_id` - AMP consumer ID.
  * `cls_notices` - A maximum of one alarm notification can be pushed to the CLS service.
    * `enable` - Start-stop status, can not be transmitted, default enabled. 0= Disabled, 1= enabled.
    * `log_set_id` - Log collection Id.
    * `region` - Regional.
    * `topic_id` - Theme Id.
  * `id` - Alarm notification template ID.
  * `is_preset` - Whether it is the system default notification template 0=No 1=Yes.
  * `name` - Alarm notification template name.
  * `notice_language` - Notification language zh-CN=Chinese en-US=English.
  * `notice_type` - Alarm notification type ALARM=Notification not restored OK=Notification restored ALL.
  * `policy_ids` - List of alarm policy IDs bound to the alarm notification template.
  * `updated_at` - Last modified time.
  * `updated_by` - Last Modified By.
  * `url_notices` - The maximum number of callback notifications is 3.
    * `end_time` - Notification End Time Seconds at the start of a day.
    * `start_time` - Notification Start Time Number of seconds at the start of a day.
    * `url` - Callback URL (limited to 256 characters).
    * `weekday` - Notification period 1-7 indicates Monday to Sunday.
  * `user_notices` - Alarm notification template list.(At most five).
    * `end_time` - The number of seconds since the notification end time 00:00:00 (value range 0-86399).
    * `group_ids` - User group ID list.
    * `need_phone_arrive_notice` - Contact notification required 0= No 1= Yes.
    * `notice_way` - Notification Channel List EMAIL=Mail SMS=SMS CALL=Telephone WECHAT=WeChat RTX=Enterprise WeChat.
    * `phone_call_type` - Call type SYNC= Simultaneous call CIRCLE= Round call If this parameter is not specified, the default value is round call.
    * `phone_circle_interval` - Number of seconds between polls (value range: 60-900).
    * `phone_circle_times` - Number of telephone polls (value range: 1-5).
    * `phone_inner_interval` - Number of seconds between calls in a polling session (value range: 60-900).
    * `phone_order` - Telephone polling list.
    * `receiver_type` - Recipient Type USER=User GROUP=User Group.
    * `start_time` - The number of seconds since the notification start time 00:00:00 (value range 0-86399).
    * `user_ids` - User UID List.
    * `weekday` - Notification period 1-7 indicates Monday to Sunday.


