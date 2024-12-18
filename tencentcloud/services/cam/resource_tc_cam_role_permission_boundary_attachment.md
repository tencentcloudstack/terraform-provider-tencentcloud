Provides a resource to create a CAM role permission boundary attachment

Example Usage

Use role_name

```hcl
resource "tencentcloud_cam_role_permission_boundary_attachment" "example" {
  policy_id = 1
  role_name = "tf-example"
}
```

Use role_id

```hcl
resource "tencentcloud_cam_role_permission_boundary_attachment" "example" {
  policy_id = 1
  role_id   = "4611686018441060141"
}
```

Use all

```hcl
resource "tencentcloud_cam_role_permission_boundary_attachment" "example" {
  policy_id = 1
  role_name = "tf-example"
  role_id   = "4611686018441060141"
}
```

Import

CAM role permission boundary attachment can be imported using the id, e.g.

```
terraform import tencentcloud_cam_role_permission_boundary_attachment.example 1##tf-example

terraform import tencentcloud_cam_role_permission_boundary_attachment.example 1#4611686018441060141#

terraform import tencentcloud_cam_role_permission_boundary_attachment.example 1#4611686018441060141#tf-example
```