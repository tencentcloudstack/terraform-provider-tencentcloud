Provides a resource to manage CAM user.

Example Usage

```hcl
resource "tencentcloud_cam_user" "foo" {
  name                = "tf_cam_user"
  remark              = "tf_user_test"
  console_login       = true
  use_api             = true
  need_reset_password = true
  password            = "Gail@1234"
  phone_num           = "12345678910"
  email               = "hello@test.com"
  country_code        = "86"
  force_delete        = true
  tags = {
    test  = "tf_cam_user",
  }
}
```

Import

CAM user can be imported using the user name, e.g.

```
$ terraform import tencentcloud_cam_user.foo cam-user-test
```