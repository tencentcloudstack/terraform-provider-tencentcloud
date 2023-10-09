---
subcategory: "Real User Monitoring(RUM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_rum_rum_log_stats_log_list"
sidebar_current: "docs-tencentcloud-datasource-rum_rum_log_stats_log_list"
description: |-
  Use this data source to query detailed information of rum rum_log_stats_log_list
---

# tencentcloud_rum_rum_log_stats_log_list

Use this data source to query detailed information of rum rum_log_stats_log_list

## Example Usage

```hcl
data "tencentcloud_rum_rum_log_stats_log_list" "rum_log_stats_log_list" {
  start_time = 1625444040
  query      = "id:123 AND type:\"log\""
  end_time   = 1625454840
  project_id = 1
}
```

## Argument Reference

The following arguments are supported:

* `end_time` - (Required, String) End time but is represented using a timestamp in seconds.
* `project_id` - (Required, Int) Project ID.
* `query` - (Required, String) Log Query syntax statement.
* `start_time` - (Required, String) Start time but is represented using a timestamp in seconds.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `result` - Return value.


