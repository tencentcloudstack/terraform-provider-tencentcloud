Provides a resource to create a CAM group membership.

Example Usage

```hcl
// create cam user
resource "tencentcloud_cam_user" "example" {
  name                = "tf-example"
  remark              = "remark."
  console_login       = true
  use_api             = true
  need_reset_password = true
  password            = "Password@123"
  phone_num           = "18611111111"
  email               = "example@tencent.com"
  country_code        = "86"
  force_delete        = true
  tags = {
    createBy  = "Terraform"
  }
}

// create cam group
resource "tencentcloud_cam_group" "example" {
  name   = "tf-example"
  remark = "remark."
}

// create cam group membership
resource "tencentcloud_cam_group_membership" "example" {
  group_id = tencentcloud_cam_group.example.id
  user_names = [tencentcloud_cam_user.example.id]
}
```

Import

CAM group membership can be imported using the id, e.g.

```
$ terraform import tencentcloud_cam_group_membership.example 353251
```