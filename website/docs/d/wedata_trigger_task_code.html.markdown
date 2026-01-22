---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_trigger_task_code"
sidebar_current: "docs-tencentcloud-datasource-wedata_trigger_task_code"
description: |-
  Use this data source to query detailed information of wedata trigger task code
---

# tencentcloud_wedata_trigger_task_code

Use this data source to query detailed information of wedata trigger task code

## Example Usage

```hcl
data "tencentcloud_wedata_trigger_task_code" "wedata_trigger_task_code" {
  project_id = "3108707295180644352"
  task_id    = "20250109174507653"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required, String) project id.
* `task_id` - (Required, String) task Id.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `data` - Get task code results.


