Provides a resource to create a dlc add_users_to_work_group_attachment

Example Usage

```hcl
resource "tencentcloud_dlc_add_users_to_work_group_attachment" "add_users_to_work_group_attachment" {
  add_info {
    work_group_id = 23184
    user_ids = [100032676511]
  }
}
}
```

Import

dlc add_users_to_work_group_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_dlc_add_users_to_work_group_attachment.add_users_to_work_group_attachment add_users_to_work_group_attachment_id
```