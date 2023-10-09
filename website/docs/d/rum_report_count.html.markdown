---
subcategory: "Real User Monitoring(RUM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_rum_report_count"
sidebar_current: "docs-tencentcloud-datasource-rum_report_count"
description: |-
  Use this data source to query detailed information of rum report_count
---

# tencentcloud_rum_report_count

Use this data source to query detailed information of rum report_count

## Example Usage

```hcl
data "tencentcloud_rum_report_count" "report_count" {
  start_time  = 1625444040
  end_time    = 1625454840
  project_id  = 1
  report_type = "log"
}
```

## Argument Reference

The following arguments are supported:

* `end_time` - (Required, Int) End time but is represented using a timestamp in seconds.
* `project_id` - (Required, Int) Project ID.
* `start_time` - (Required, Int) Start time but is represented using a timestamp in seconds.
* `instance_id` - (Optional, String) Instance ID.
* `report_type` - (Optional, String) Report type, empty is meaning all type count. `log`:log report count, `pv`:pv report count, `event`:event report count, `speed`:speed report count, `performance`:performance report count, `custom`:custom report count, `webvitals`:webvitals report count, `miniProgramData`:miniProgramData report count.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `result` - Return value.


