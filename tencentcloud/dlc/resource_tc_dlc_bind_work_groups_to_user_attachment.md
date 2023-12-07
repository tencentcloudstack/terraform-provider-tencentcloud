Provides a resource to create a dlc bind_work_groups_to_user

Example Usage

```hcl
resource "tencentcloud_dlc_bind_work_groups_to_user_attachment" "bind_work_groups_to_user" {
  add_info {
    user_id = "100032772113"
    work_group_ids = [23184,23181]
  }
}
```

Import

dlc bind_work_groups_to_user can be imported using the id, e.g.

```
terraform import tencentcloud_dlc_bind_work_groups_to_user_attachment.bind_work_groups_to_user bind_work_groups_to_user_id
```