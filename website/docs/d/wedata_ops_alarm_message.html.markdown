---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_ops_alarm_message"
sidebar_current: "docs-tencentcloud-datasource-wedata_ops_alarm_message"
description: |-
  Use this data source to query detailed information of wedata ops alarm message
---

# tencentcloud_wedata_ops_alarm_message

Use this data source to query detailed information of wedata ops alarm message

## Example Usage

```hcl
data "tencentcloud_wedata_ops_alarm_message" "wedata_ops_alarm_message" {
  project_id       = "1859317240494305280"
  alarm_message_id = 263840
}
```

## Argument Reference

The following arguments are supported:

* `alarm_message_id` - (Required, String) Alarm message Id.
* `project_id` - (Required, String) Project ID.
* `result_output_file` - (Optional, String) Used to save results.
* `time_zone` - (Optional, String) Specifies the time zone of the return date. default UTC+8.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `data` - Alarm information.


