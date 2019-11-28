---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cam_user"
sidebar_current: "docs-tencentcloud-resource-cam_user"
description: |-
  Provides a resource to create a CAM user.
---

# tencentcloud_cam_user

Provides a resource to create a CAM user.

## Example Usage

```hcl
resource "tencentcloud_cam_user" "foo" {
  name                = "cam-user-test"
  remark              = "test"
  console_login       = true
  use_api             = true
  need_reset_password = true
  password            = "Gail@1234"
  phone_num           = "12345678910"
  country_code        = "86"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, ForceNew) Name of CAM user.
* `console_login` - (Optional) Indicade whether the CAM user can login or not.
* `country_code` - (Optional) Country code of the phone num, like '86'.
* `email` - (Optional) Email of the CAM user.
* `need_reset_password` - (Optional) Indicate whether the CAM user will reset the password the next time he/her logs in.
* `password` - (Optional) The password of the CAM user. The password should be set with 8 characters or more and contains uppercase small letters, numbers, and special characters. Only valid when console_login set true. If not set and the value of console_login is true, a random password is automatically generated.
* `phone_num` - (Optional) Phone num of the CAM user.
* `remark` - (Optional) Remark of the CAM user.
* `use_api` - (Optional) Indicate whether to generate a secret key or not.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `secret_id` - Secret Id of the CAM user.
* `secret_key` - Secret key of the CAM user.
* `uid` - Id of the CAM user.
* `uin` - Uin of the CAM User.


## Import

CAM user can be imported using the id, e.g.

```
$ terraform import tencentcloud_cam_user.foo cam-user-test
```

