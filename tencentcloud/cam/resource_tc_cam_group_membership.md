Provides a resource to create a CAM group membership.

Example Usage

```hcl
variable "cam_group_basic" {
  default = "keep-cam-group"
}

data "tencentcloud_cam_groups" "groups" {
  name   = var.cam_group_basic
}

resource "tencentcloud_cam_user" "foo" {
  name                = "tf_cam_user"
  remark              = "tf_user_remark"
  console_login       = true
  use_api             = true
  need_reset_password = true
  password            = "Gail@1234"
  phone_num           = "12345678910"
  country_code        = "86"
  email               = "1234@qq.com"
  force_delete        = true
}

resource "tencentcloud_cam_group_membership" "group_membership_basic" {
  group_id = data.tencentcloud_cam_groups.groups.group_list.0.group_id
  user_names = [tencentcloud_cam_user.foo.id]
}

```

Import

CAM group membership can be imported using the id, e.g.

```
$ terraform import tencentcloud_cam_group_membership.foo 12515263
```