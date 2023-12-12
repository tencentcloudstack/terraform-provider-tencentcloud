Provides a resource to create a CAM role policy attachment.

Example Usage

```hcl
variable "cam_policy_basic" {
  default = "keep-cam-policy"
}

variable "cam_role_basic" {
  default = "keep-cam-role"
}

data "tencentcloud_cam_policies" "policy" {
  name        = var.cam_policy_basic
}

data "tencentcloud_cam_roles" "roles" {
  name        = var.cam_role_basic
}

resource "tencentcloud_cam_role_policy_attachment_by_name" "role_policy_attachment_basic" {
  role_name   = var.cam_role_basic
  policy_name = var.cam_policy_basic
}
```

Import

CAM role policy attachment can be imported using the id, e.g.

```
$ terraform import tencentcloud_cam_role_policy_attachment_by_name.foo ${role_name}#${policy_name}
```