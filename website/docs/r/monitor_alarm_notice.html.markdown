---
subcategory: "Cloud Monitor(Monitor)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_alarm_notice"
sidebar_current: "docs-tencentcloud-resource-monitor_alarm_notice"
description: |-
  Provides a alarm notice resource for monitor.
---

# tencentcloud_monitor_alarm_notice

Provides a alarm notice resource for monitor.

## Example Usage

```hcl
resource "tencentcloud_monitor_alarm_notice" "example" {
  name            = "test_alarm_notice_1"
  notice_type     = "ALL"
  notice_language = "zh-CN"

  user_notices {
    receiver_type            = "USER"
    start_time               = 0
    end_time                 = 1
    notice_way               = ["SMS", "EMAIL"]
    user_ids                 = [10001]
    group_ids                = []
    phone_order              = [10001]
    phone_circle_times       = 2
    phone_circle_interval    = 50
    phone_inner_interval     = 60
    need_phone_arrive_notice = 1
    phone_call_type          = "CIRCLE"
    weekday                  = [1, 2, 3, 4, 5, 6, 7]
  }

  url_notices {
    url        = "https://www.mytest.com/validate"
    end_time   = 0
    start_time = 1
    weekday    = [1, 2, 3, 4, 5, 6, 7]
  }

}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Notification template name within 60.
* `notice_language` - (Required, String) Notification language zh-CN=Chinese en-US=English.
* `notice_type` - (Required, String) Alarm notification type ALARM=Notification not restored OK=Notification restored ALL.
* `cls_notices` - (Optional, List) A maximum of one alarm notification can be pushed to the CLS service.
* `url_notices` - (Optional, List) The maximum number of callback notifications is 3.
* `user_notices` - (Optional, List) Alarm notification template list.(At most five).

The `cls_notices` object supports the following:

* `log_set_id` - (Required, String) Log collection Id.
* `region` - (Required, String) Regional.
* `topic_id` - (Required, String) Theme Id.
* `enable` - (Optional, Int) Start-stop status, can not be transmitted, default enabled. 0= Disabled, 1= enabled.

The `url_notices` object supports the following:

* `url` - (Required, String) Callback URL (limited to 256 characters).
* `end_time` - (Optional, Int) Notification End Time Seconds at the start of a day.
* `start_time` - (Optional, Int) Notification Start Time Number of seconds at the start of a day.
* `weekday` - (Optional, Set) Notification period 1-7 indicates Monday to Sunday.

The `user_notices` object supports the following:

* `end_time` - (Required, Int) The number of seconds since the notification end time 00:00:00 (value range 0-86399).
* `notice_way` - (Required, Set) Notification Channel List EMAIL=Mail SMS=SMS CALL=Telephone WECHAT=WeChat RTX=Enterprise WeChat.
* `receiver_type` - (Required, String) Recipient Type USER=User GROUP=User Group.
* `start_time` - (Required, Int) The number of seconds since the notification start time 00:00:00 (value range 0-86399).
* `group_ids` - (Optional, Set) User group ID list.
* `need_phone_arrive_notice` - (Optional, Int) Contact notification required 0= No 1= Yes.
* `phone_call_type` - (Optional, String) Call type SYNC= Simultaneous call CIRCLE= Round call If this parameter is not specified, the default value is round call.
* `phone_circle_interval` - (Optional, Int) Number of seconds between polls (value range: 60-900).
* `phone_circle_times` - (Optional, Int) Number of telephone polls (value range: 1-5).
* `phone_inner_interval` - (Optional, Int) Number of seconds between calls in a polling session (value range: 60-900).
* `phone_order` - (Optional, Set) Telephone polling list.
* `user_ids` - (Optional, Set) User UID List.
* `weekday` - (Optional, Set) Notification period 1-7 indicates Monday to Sunday.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `is_preset` - Whether it is the system default notification template 0=No 1=Yes.
* `policy_ids` - List of alarm policy IDs bound to the alarm notification template.
* `updated_at` - Last modified time.
* `updated_by` - Last Modified By.


## Import

Monitor Alarm Notice can be imported, e.g.

```
$ terraform import tencentcloud_monitor_alarm_notice.import-test noticeId
```

