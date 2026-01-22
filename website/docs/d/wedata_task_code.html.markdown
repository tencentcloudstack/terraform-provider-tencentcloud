---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_task_code"
sidebar_current: "docs-tencentcloud-datasource-wedata_task_code"
description: |-
  Use this data source to query detailed information of wedata wedata_task_code
---

# tencentcloud_wedata_task_code

Use this data source to query detailed information of wedata wedata_task_code

## Example Usage

```hcl
data "tencentcloud_wedata_task_code" "wedata_task_code" {
  project_id = "2905622749543821312"
  task_id    = "20251015164958429"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required, String) The project id.
* `task_id` - (Required, String) Task ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `data` - Retrieves the task code result.


