Provides a resource to create a cam user_permission_boundary

Example Usage

```hcl
resource "tencentcloud_cam_user_permission_boundary_attachment" "user_permission_boundary" {
  target_uin = 100032767426
  policy_id = 151113272
}
```

Import

cam user_permission_boundary can be imported using the id, e.g.

```
terraform import tencentcloud_cam_user_permission_boundary_attachment.user_permission_boundary user_permission_boundary_id
```