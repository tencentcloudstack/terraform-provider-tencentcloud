---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_downstream_task_instances"
sidebar_current: "docs-tencentcloud-datasource-wedata_downstream_task_instances"
description: |-
  Use this data source to query detailed information of wedata downstream task instances
---

# tencentcloud_wedata_downstream_task_instances

Use this data source to query detailed information of wedata downstream task instances

## Example Usage

```hcl
data "tencentcloud_wedata_downstream_task_instances" "wedata_down_task_instances" {
  project_id   = "1859317240494305280"
  instance_key = "20250731151633120_2025-10-13 17:00:00"
}
```

## Argument Reference

The following arguments are supported:

* `instance_key` - (Required, String) Instance unique identifier.
* `project_id` - (Required, String) Project ID.
* `result_output_file` - (Optional, String) Used to save results.
* `time_zone` - (Optional, String) Time zone timeZone, default UTC+8.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `data` - Direct downstream task instances list.


