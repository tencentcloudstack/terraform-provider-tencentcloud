Provides a resource to create a DLC add users to work group attachment

Example Usage

```hcl
resource "tencentcloud_dlc_add_users_to_work_group_attachment" "example" {
  add_info {
    work_group_id = 70220
    user_ids      = ["100032717595", "100030773831"]
  }
}
```

Import

DLC add users to work group attachment can be imported using the id, e.g.

```
terraform import tencentcloud_dlc_add_users_to_work_group_attachment.example '70220#100032717595|100030773831'
```
