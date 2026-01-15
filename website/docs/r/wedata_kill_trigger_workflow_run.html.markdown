---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_kill_trigger_workflow_run"
sidebar_current: "docs-tencentcloud-resource-wedata_kill_trigger_workflow_run"
description: |-
  Provides a resource to kill wedata trigger workflow run
---

# tencentcloud_wedata_kill_trigger_workflow_run

Provides a resource to kill wedata trigger workflow run

~> **NOTE:** Both "all" and "pending" require obtaining the execution_id through the query interface before passing it as a parameter..

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

resource "tencentcloud_wedata_kill_trigger_workflow_run" "kill_all" {
  project_id            = "3108707295180644352"
  workflow_id           = "333368d7-bc8e-4b95-9a66-7a5151063eb2"
  workflow_execution_id = data.tencentcloud_wedata_trigger_workflow_runs.trigger_workflow_runs.data[0].items[0].execution_id
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required, String, ForceNew) Project ID.
* `workflow_execution_id` - (Required, String, ForceNew) Workflow execution ID to stop.
* `workflow_id` - (Required, String, ForceNew) Workflow ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



