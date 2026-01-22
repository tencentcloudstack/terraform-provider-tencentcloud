---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_rerun_trigger_workflow_run_async"
sidebar_current: "docs-tencentcloud-resource-wedata_rerun_trigger_workflow_run_async"
description: |-
  Provides a resource to rerun wedata trigger workflow run asynchronously
---

# tencentcloud_wedata_rerun_trigger_workflow_run_async

Provides a resource to rerun wedata trigger workflow run asynchronously

## Example Usage

```hcl
data "tencentcloud_wedata_trigger_workflow_runs" "trigger_workflow_runs" {
  project_id = "3108707295180644352"
  filters {
    name   = "WorkflowId"
    values = ["333368d7-bc8e-4b95-9a66-7a5151063eb2"]
  }
  order_fields {
    name      = "CreateTime"
    direction = "DESC"
  }
}

resource "tencentcloud_wedata_rerun_trigger_workflow_run_async" "rerun_basic" {
  project_id            = "3108707295180644352"
  workflow_id           = "333368d7-bc8e-4b95-9a66-7a5151063eb2"
  workflow_execution_id = data.tencentcloud_wedata_trigger_workflow_runs.trigger_workflow_runs.data[0].items[0].execution_id
  execute_type          = "1"
}
```

## Argument Reference

The following arguments are supported:

* `execute_type` - (Required, String, ForceNew) Execution type: Normal execution with default parameters: 1, Advanced execution with optional task scope and parameters: 2.
* `project_id` - (Required, String, ForceNew) Project ID.
* `workflow_execution_id` - (Required, String, ForceNew) Workflow execution ID.
* `workflow_id` - (Required, String, ForceNew) Workflow ID.
* `advanced_params` - (Optional, List, ForceNew) Custom execution parameters for advanced execution type.
* `integration_resource_group` - (Optional, String, ForceNew) Specified integration resource group, defaults to the original configured integration resource group if empty.
* `scheduling_resource_group` - (Optional, String, ForceNew) Specified scheduling resource group, defaults to the original configured scheduling resource group if empty.
* `task_ids` - (Optional, Set: [`String`], ForceNew) Set of specific task IDs to run in advanced execution mode.

The `advanced_params` object supports the following:

* `ext_properties` - (Optional, String) Extended properties in JSON format, example: "{}".
* `param_key` - (Optional, String) Parameter name.
* `param_value` - (Optional, String) Parameter value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



