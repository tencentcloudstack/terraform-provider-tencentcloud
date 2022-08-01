---
subcategory: "Monitor"
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
  module          = "monitor"
  name            = "yourname"
  notice_type     = "ALL"
  notice_language = "zh-CN"

}
```

## Argument Reference

The following arguments are supported:

* `module` - (Required, String) Module name, fill in 'monitor' here.
* `name` - (Required, String) Notification template name within 60.
* `notice_ids` - (Optional, List: [`String`]) List of notification rule IDs.
* `notice_language` - (Optional, String) Notification language zh-CN=Chinese en-US=English.
* `notice_type` - (Optional, String) Alarm notification type ALARM=Notification not restored OK=Notification restored ALL.
* `user_notices` - (Optional, List) Alarm notification template list.

The `notice_way` object supports the following:

* `endtime` - (Optional, Int) The number of seconds since the notification start time 00:00:00 (value range 0-86399).

The `user_notices` object supports the following:

* `endtime` - (Optional, Int) The number of seconds since the notification start time 00:00:00 (value range 0-86399).
* `notice_way` - (Optional, List) Notification Channel List EMAIL=Mail SMS=SMS CALL=Telephone WECHAT=WeChat RTX=Enterprise WeChat.
* `receiver_type` - (Optional, String) Recipient Type USER=User GROUP=User Group.
* `start_time` - (Optional, Int) The number of seconds since the notification start time 00:00:00 (value range 0-86399).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `notice_id` - Alarm notification template ID.
* `notices` - Alarm notification template list.
* `request_id` - Unique request ID, returned on every request. When locating the problem, you need to provide the RequestId of the request.


