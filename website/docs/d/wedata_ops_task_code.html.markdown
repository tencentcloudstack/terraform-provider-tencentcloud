---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_ops_task_code"
sidebar_current: "docs-tencentcloud-datasource-wedata_ops_task_code"
description: |-
  Use this data source to query detailed information of wedata ops task code
---

# tencentcloud_wedata_ops_task_code

Use this data source to query detailed information of wedata ops task code

## Example Usage

```hcl
data "tencentcloud_wedata_ops_task_code" "wedata_ops_task_code" {
  project_id = "1859317240494305280"
  task_id    = "20230901114849281"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required, String) Project id.
* `task_id` - (Required, String) Task ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `data` - Retrieves the task code result.


