---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_workflow_max_permission"
sidebar_current: "docs-tencentcloud-datasource-wedata_workflow_max_permission"
description: |-
  Use this data source to query detailed information of WeData workflow max permission
---

# tencentcloud_wedata_workflow_max_permission

Use this data source to query detailed information of WeData workflow max permission

## Example Usage

```hcl
data "tencentcloud_wedata_workflow_max_permission" "example" {
  project_id  = "3108707295180644352"
  entity_id   = "53e78f97-f145-11f0-ba36-b8cef6a5af5c"
  entity_type = "folder"
}
```

## Argument Reference

The following arguments are supported:

* `entity_id` - (Required, String) Authorization entity ID.
* `entity_type` - (Required, String) Authorization entity type, folder/workflow.
* `project_id` - (Required, String) Project ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `data` - Current user's recursive maximum permission type for entity resources.
  * `permission_type` - Authorization permission type (CAN_VIEW/CAN_RUN/CAN_EDIT/CAN_MANAGE, currently only supports CAN_MANAGE).


