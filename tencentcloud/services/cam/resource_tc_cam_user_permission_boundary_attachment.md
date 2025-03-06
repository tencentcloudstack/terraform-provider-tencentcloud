Provides a resource to create a CAM user permission boundary

Example Usage

```hcl
resource "tencentcloud_cam_user_permission_boundary_attachment" "example" {
  target_uin = 100037718101
  policy_id  = 234290251
}
```

Import

CAM user permission boundary can be imported using the id, e.g.

```
terraform import tencentcloud_cam_user_permission_boundary_attachment.example 100037718101#234290251
```