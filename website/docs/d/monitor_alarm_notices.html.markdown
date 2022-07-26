---
subcategory: "Monitor"
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
  module     = "monitor"
  pagenumber = 1
  pagesize   = 20
  order      = "DESC"

}
```

## Argument Reference

The following arguments are supported:

* `module` - (Required, String) Module name, fill in 'monitor' here.
* `order` - (Required, String) Sort by update time ASC=forward order DESC=reverse order.
* `pagenumber` - (Required, Int) Page number minimum 1.
* `pagesize` - (Required, Int) Page size 1-200.
* `notices` - (Optional, List) Alarm notification template list.
* `result_output_file` - (Optional, String) Used to store results.

The `notices` object supports the following:

* `is_preset` - (Optional, Int) Whether it is the system default notification template 0=No 1=Yes.
* `notice_language` - (Optional, String) Notification language zh-CN=Chinese en-US=English.
* `notice_type` - (Optional, String) Alarm notification type ALARM=Notification not restored OK=Notification restored ALL.
* `notices_id` - (Optional, String) Alarm notification template ID.
* `notices_name` - (Optional, String) Alarm notification template name.
* `updated_at` - (Optional, String) Last modified time.
* `updated_by` - (Optional, String) Last Modified By.
* `user_notices` - (Optional, List) Alarm notification template list.

The `user_notices` object supports the following:

* `endtime` - (Optional, Int) The number of seconds since the notification start time 00:00:00 (value range 0-86399).
* `receiver_type` - (Optional, String) Recipient Type USER=User GROUP=User Group.
* `start_time` - (Optional, Int) The number of seconds since the notification start time 00:00:00 (value range 0-86399).


