---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_trigger_workflow_runs"
sidebar_current: "docs-tencentcloud-datasource-wedata_trigger_workflow_runs"
description: |-
  Use this data source to query detailed information of wedata trigger workflow runs.
---

# tencentcloud_wedata_trigger_workflow_runs

Use this data source to query detailed information of wedata trigger workflow runs.

## Example Usage

```hcl
data "tencentcloud_wedata_trigger_workflow_runs" "trigger_workflow_runs" {
  project_id = "1840731342293643264"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required, String) Project ID.
* `filters` - (Optional, List) Filter parameters. Workflow name or ID query name: `Keyword`; workflow ID query name: `WorkflowId`; folder query name: `FolderId`; owner query name: `InChargeUin`; workflow execution ID: `ExecutionId`.
* `order_fields` - (Optional, List) Sort fields. Sort field names include, for example, start time: `CreateTime`; end time: `EndTime`.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `name` - (Optional, String) Filter field name.
* `values` - (Optional, Set) List of filter values.

The `order_fields` object supports the following:

* `direction` - (Required, String) Sort direction: ASC|DESC.
* `name` - (Required, String) Sort field name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `data` - Workflow run query results.


