Provides a resource to manage CAM user.

Example Usage

```hcl
resource "tencentcloud_cam_user" "example" {
  name                = "tf-example"
  remark              = "Remark."
  console_login       = true
  use_api             = true
  need_reset_password = true
  password            = "Password@123"
  phone_num           = "189********"
  email               = "example@qq.com"
  country_code        = "86"
  force_delete        = true
  tags = {
    createBy = "Terraform"
  }
}
```

Import

CAM user can be imported using the user name, e.g.

```
$ terraform import tencentcloud_cam_user.example tf-example
```