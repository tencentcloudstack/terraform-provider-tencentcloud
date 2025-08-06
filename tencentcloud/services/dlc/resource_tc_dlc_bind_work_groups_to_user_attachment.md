Provides a resource to create a DLC bind work groups to user

Example Usage

```hcl
resource "tencentcloud_dlc_bind_work_groups_to_user_attachment" "example" {
  add_info {
    user_id        = "100032772113"
    work_group_ids = [23184, 23181]
  }
}
```

Import

DLC bind work groups to user can be imported using the id, e.g.

```
terraform import tencentcloud_dlc_bind_work_groups_to_user_attachment.example 100032772113
```
