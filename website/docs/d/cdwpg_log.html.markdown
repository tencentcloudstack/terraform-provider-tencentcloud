---
subcategory: "CDWPG"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cdwpg_log"
sidebar_current: "docs-tencentcloud-datasource-cdwpg_log"
description: |-
  Use this data source to query detailed information of cdwpg cdwpg_log
---

# tencentcloud_cdwpg_log

Use this data source to query detailed information of cdwpg cdwpg_log

## Example Usage

```hcl
data "tencentcloud_cdwpg_log" "cdwpg_log" {
  instance_id = "cdwpg-gexy9tue"
  start_time  = "2025-03-21 00:00:00"
  end_time    = "2025-03-21 23:59:59"
}
```

## Argument Reference

The following arguments are supported:

* `end_time` - (Required, String) End time.
* `instance_id` - (Required, String) Instance id.
* `start_time` - (Required, String) Start time.
* `database` - (Optional, String) Database.
* `duration` - (Optional, Float64) Filter duration.
* `order_by_type` - (Optional, String) Ascending/Descending.
* `order_by` - (Optional, String) Sort by.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `error_log_details` - Error log details.
* `slow_log_details` - Slow sql log details.


