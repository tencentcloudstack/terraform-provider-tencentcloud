---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_ops_workflow"
sidebar_current: "docs-tencentcloud-datasource-wedata_ops_workflow"
description: |-
  Use this data source to query detailed information of wedata ops workflow
---

# tencentcloud_wedata_ops_workflow

Use this data source to query detailed information of wedata ops workflow

## Example Usage

```hcl
data "tencentcloud_wedata_ops_workflow" "wedata_ops_workflow" {
  project_id  = "2905622749543821312"
  workflow_id = "f328ab83-62e1-4b0a-9a18-a79b42722792"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required, String) Project ID.
* `workflow_id` - (Required, String) Workflow ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `data` - Workflow scheduling details.


