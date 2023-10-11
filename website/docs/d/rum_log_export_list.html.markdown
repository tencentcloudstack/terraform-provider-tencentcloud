---
subcategory: "Real User Monitoring(RUM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_rum_log_export_list"
sidebar_current: "docs-tencentcloud-datasource-rum_log_export_list"
description: |-
  Use this data source to query detailed information of rum log_export_list
---

# tencentcloud_rum_log_export_list

Use this data source to query detailed information of rum log_export_list

## Example Usage

```hcl
data "tencentcloud_rum_log_export_list" "log_export_list" {
  project_id = 1
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required, Int) Project ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `result` - Return result.


