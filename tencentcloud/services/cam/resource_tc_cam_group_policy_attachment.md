Provides a resource to create a CAM group policy attachment.

Example Usage

```hcl
variable "cam_policy_basic" {
  default = "keep-cam-policy"
}

variable "cam_group_basic" {
  default = "keep-cam-group"
}

data "tencentcloud_cam_groups" "groups" {
  name = var.cam_group_basic
}

data "tencentcloud_cam_policies" "policy" {
  name = var.cam_policy_basic
}

resource "tencentcloud_cam_group_policy_attachment" "group_policy_attachment_basic" {
  group_id  = data.tencentcloud_cam_groups.groups.group_list.0.group_id
  policy_id = data.tencentcloud_cam_policies.policy.policy_list.0.policy_id
}
```

Import

CAM group policy attachment can be imported using the id, e.g.

```
$ terraform import tencentcloud_cam_group_policy_attachment.foo 12515263#26800353
```