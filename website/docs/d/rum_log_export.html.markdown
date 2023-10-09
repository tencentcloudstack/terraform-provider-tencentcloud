---
subcategory: "Real User Monitoring(RUM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_rum_log_export"
sidebar_current: "docs-tencentcloud-datasource-rum_log_export"
description: |-
  Use this data source to query detailed information of rum log_export
---

# tencentcloud_rum_log_export

Use this data source to query detailed information of rum log_export

## Example Usage

```hcl
data "tencentcloud_rum_log_export" "log_export" {
  name       = "log"
  start_time = "1692594840000"
  query      = "id:123 AND type: \"log\""
  end_time   = "1692609240000"
  project_id = 1
}
```

## Argument Reference

The following arguments are supported:

* `end_time` - (Required, String) End timestamp, in milliseconds.
* `name` - (Required, String) Export flag name.
* `project_id` - (Required, Int) Project ID.
* `query` - (Required, String) Log Query syntax statement.
* `start_time` - (Required, String) Start timestamp, in milliseconds.
* `fields` - (Optional, Set: [`String`]) Log fields.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `result` - Return result.


