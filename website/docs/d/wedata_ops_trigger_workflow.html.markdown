---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_ops_trigger_workflow"
sidebar_current: "docs-tencentcloud-datasource-wedata_ops_trigger_workflow"
description: |-
  Use this data source to query detailed information of wedata ops trigger workflow.
---

# tencentcloud_wedata_ops_trigger_workflow

Use this data source to query detailed information of wedata ops trigger workflow.

## Example Usage

```hcl
data "tencentcloud_wedata_ops_trigger_workflow" "ops_trigger_workflow" {
  project_id  = "3108707295180644352"
  workflow_id = "b41e8d13-905a-4006-9d05-1fe180338f59"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required, String) Project ID.
* `workflow_id` - (Required, String) Workflow ID.
* `result_output_file` - (Optional, String) Used to save results.
* `workflow_execution_id` - (Optional, String) Workflow execution ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `data` - Workflow task information.


