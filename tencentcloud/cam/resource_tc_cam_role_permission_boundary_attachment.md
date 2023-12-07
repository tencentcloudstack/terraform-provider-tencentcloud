Provides a resource to create a cam role_permission_boundary_attachment

Example Usage

```hcl
resource "tencentcloud_cam_role_permission_boundary_attachment" "role_permission_boundary_attachment" {
  policy_id = 1
  role_name = "test-cam-tag"
}
```

Import

cam role_permission_boundary_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_cam_role_permission_boundary_attachment.role_permission_boundary_attachment role_permission_boundary_attachment_id
```