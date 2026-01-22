---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_trigger_workflow_run"
sidebar_current: "docs-tencentcloud-datasource-wedata_trigger_workflow_run"
description: |-
  Use this data source to query detailed information of wedata trigger_workflow_run
---

# tencentcloud_wedata_trigger_workflow_run

Use this data source to query detailed information of wedata trigger_workflow_run

## Example Usage

```hcl
data "tencentcloud_wedata_trigger_workflow_run" "trigger_workflow_run" {
  project_id            = "1840731342293087232"
  workflow_execution_id = "82c15b04-a6ef-4075-bed2-d20d23457297"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required, String) Project ID.
* `workflow_execution_id` - (Required, String) Workflow execution ID.
* `filters` - (Optional, List) Filter conditions.
* `order_fields` - (Optional, List) Sort conditions.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `name` - (Optional, String) Filter field name.
* `values` - (Optional, Set) Filter value list.

The `order_fields` object supports the following:

* `direction` - (Required, String) Sort direction: `ASC`, `DESC`.
* `name` - (Required, String) Sort field name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `data` - Workflow task information.


