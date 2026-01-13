---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_ops_trigger_workflows"
sidebar_current: "docs-tencentcloud-datasource-wedata_ops_trigger_workflows"
description: |-
  Use this data source to query detailed information of wedata ops trigger workflows.
---

# tencentcloud_wedata_ops_trigger_workflows

Use this data source to query detailed information of wedata ops trigger workflows.

## Example Usage

```hcl
data "tencentcloud_wedata_ops_trigger_workflows" "ops_trigger_workflows" {
  project_id = "1840731342293643264"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required, String) Project ID.
* `filters` - (Optional, List) Filter parameters. Workflow name or ID query name: `Keyword`; workflow ID query name: `WorkflowId`; folder query name: `FolderId`; owner query name: `InChargeUin`.
* `order_fields` - (Optional, List) Sort fields. Sort field names include, for example, task count: TaskCount.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `name` - (Optional, String) Filter field name.
* `values` - (Optional, Set) List of filter values.

The `order_fields` object supports the following:

* `direction` - (Required, String) Sort direction: ASC|DESC.
* `name` - (Required, String) Sort field name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `data` - Workflow query results.


