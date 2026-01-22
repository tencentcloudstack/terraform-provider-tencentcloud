---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_workflow_permissions"
sidebar_current: "docs-tencentcloud-resource-wedata_workflow_permissions"
description: |-
  Provides a resource to create a WeData workflow permissions
---

# tencentcloud_wedata_workflow_permissions

Provides a resource to create a WeData workflow permissions

## Example Usage

```hcl
resource "tencentcloud_wedata_workflow_permissions" "example" {
  project_id  = "3108707295180644352"
  entity_id   = "53e78f97-f145-11f0-ba36-b8cef6a5af5c"
  entity_type = "folder"
  permission_list {
    permission_target_type = "user"
    permission_target_id   = "100028448903"
    permission_type_list   = ["CAN_MANAGE"]
  }

  permission_list {
    permission_target_type = "role"
    permission_target_id   = "308335260676890624"
    permission_type_list   = ["CAN_MANAGE"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `entity_id` - (Required, String, ForceNew) Authorization entity ID.
* `entity_type` - (Required, String, ForceNew) Authorization entity type, folder/workflow.
* `permission_list` - (Required, Set) Authorization information array.
* `project_id` - (Required, String, ForceNew) Project ID.

The `permission_list` object supports the following:

* `permission_target_id` - (Required, String) Authorization target ID array (userId/roleId).
* `permission_target_type` - (Required, String) Authorization target type (user: user, role: role).
* `permission_type_list` - (Required, Set) Authorization permission type array (CAN_VIEW/CAN_RUN/CAN_EDIT/CAN_MANAGE, currently only supports CAN_MANAGE).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

WeData workflow permissions can be imported using the projectId#entityId#entityType, e.g.

```
terraform import tencentcloud_wedata_workflow_permissions.example 3108707295180644352#53e78f97-f145-11f0-ba36-b8cef6a5af5c#folder
```

