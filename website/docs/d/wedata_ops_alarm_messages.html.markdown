---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_ops_alarm_messages"
sidebar_current: "docs-tencentcloud-datasource-wedata_ops_alarm_messages"
description: |-
  Use this data source to query detailed information of wedata ops alarm messages
---

# tencentcloud_wedata_ops_alarm_messages

Use this data source to query detailed information of wedata ops alarm messages

## Example Usage

```hcl
data "tencentcloud_wedata_ops_alarm_messages" "wedata_ops_alarm_messages" {
  project_id  = "1859317240494305280"
  start_time  = "2025-10-14 21:09:26"
  end_time    = "2025-10-14 21:10:26"
  alarm_level = 1
  time_zone   = "UTC+8"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required, String) Project id.
* `alarm_level` - (Optional, Int) Alarm level.
* `alarm_recipient_id` - (Optional, String) Alert recipient Id.
* `end_time` - (Optional, String) Specifies the Alarm end time in the format yyyy-MM-dd HH:MM:ss.
* `result_output_file` - (Optional, String) Used to save results.
* `start_time` - (Optional, String) Starting Alarm time. format: yyyy-MM-dd HH:MM:ss.
* `time_zone` - (Optional, String) For incoming and returned filter time zone, default UTC+8.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `data` - Alarm information list.


