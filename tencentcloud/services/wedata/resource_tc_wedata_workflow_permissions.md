Provides a resource to create a WeData workflow permissions

Example Usage

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

Import

WeData workflow permissions can be imported using the projectId#entityId#entityType, e.g.

```
terraform import tencentcloud_wedata_workflow_permissions.example 3108707295180644352#53e78f97-f145-11f0-ba36-b8cef6a5af5c#folder
```
