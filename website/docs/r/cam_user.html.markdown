---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cam_user"
sidebar_current: "docs-tencentcloud-resource-cam_user"
description: |-
  Provides a resource to manage CAM user.
---

# tencentcloud_cam_user

Provides a resource to manage CAM user.

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
  email               = "hello@test.com"
  country_code        = "86"
  force_delete        = true
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, ForceNew) Name of the CAM user.
* `console_login` - (Optional) Indicate whether the CAM user can login to the web console or not.
* `country_code` - (Optional) Country code of the phone number, for example: '86'.
* `email` - (Optional) Email of the CAM user.
* `force_delete` - (Optional) Indicate whether to force deletes the CAM user. If set false, the API secret key will be checked and failed when exists; otherwise the user will be deleted directly. Default is false.
* `need_reset_password` - (Optional) Indicate whether the CAM user need to reset the password when first logins.
* `password` - (Optional) The password of the CAM user. Password should be at least 8 characters and no more than 32 characters, includes uppercase letters, lowercase letters, numbers and special characters. Only required when `console_login` is true. If not set, a random password will be automatically generated.
* `phone_num` - (Optional) Phone number of the CAM user.
* `remark` - (Optional) Remark of the CAM user.
* `use_api` - (Optional) Indicate whether to generate the API secret key or not.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `secret_id` - Secret ID of the CAM user.
* `secret_key` - Secret key of the CAM user.
* `uid` - ID of the CAM user.
* `uin` - Uin of the CAM User.


## Import

CAM user can be imported using the user name, e.g.

```
$ terraform import tencentcloud_cam_user.foo cam-user-test
```

