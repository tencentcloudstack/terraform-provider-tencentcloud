---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_trigger_task_run"
sidebar_current: "docs-tencentcloud-datasource-wedata_trigger_task_run"
description: |-
  Use this data source to query detailed information of wedata trigger task run
---

# tencentcloud_wedata_trigger_task_run

Use this data source to query detailed information of wedata trigger task run

## Example Usage

```hcl
data "tencentcloud_wedata_trigger_task_run" "trigger_task_run" {
  project_id        = "1840731342293818368"
  task_execution_id = "20260109165716558"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required, String) Workspace ID.
* `task_execution_id` - (Required, String) Task execution ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `data` - Workflow task execution list pagination result.


