---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_task_instance_executions"
sidebar_current: "docs-tencentcloud-datasource-wedata_task_instance_executions"
description: |-
  Use this data source to query detailed information of wedata task instance executions
---

# tencentcloud_wedata_task_instance_executions

Use this data source to query detailed information of wedata task instance executions

## Example Usage

```hcl
data "tencentcloud_wedata_task_instance_executions" "wedata_task_instance_executions" {
  project_id   = "1859317240494305280"
  instance_key = "20250731151633120_2025-10-13 17:00:00"
}
```

## Argument Reference

The following arguments are supported:

* `instance_key` - (Required, String) Instance unique identifier, can be obtained via ListInstances.
* `project_id` - (Required, String) Project ID to which it belongs.
* `result_output_file` - (Optional, String) Used to save results.
* `time_zone` - (Optional, String) **Time zone** timeZone, the time zone of the input time string, default UTC+8.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `data` - Instance details.


