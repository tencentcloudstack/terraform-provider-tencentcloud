---
subcategory: "Bastion Host(BH)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dasb_user"
sidebar_current: "docs-tencentcloud-resource-dasb_user"
description: |-
  Provides a resource to create a dasb user
---

# tencentcloud_dasb_user

Provides a resource to create a dasb user

## Example Usage

```hcl
resource "tencentcloud_dasb_user" "example" {
  user_name     = "tf_example"
  real_name     = "terraform"
  phone         = "+86|18345678782"
  email         = "demo@tencent.com"
  validate_from = "2023-09-22T02:00:00+08:00"
  validate_to   = "2023-09-23T03:00:00+08:00"
  department_id = "1.2"
  auth_type     = 0
}
```

## Argument Reference

The following arguments are supported:

* `real_name` - (Required, String) Real name, maximum length 20 characters, cannot contain blank characters.
* `user_name` - (Required, String) Username, 3-20 characters, must start with an English letter and cannot contain characters other than letters, numbers, '.', '_', '-'.
* `auth_type` - (Optional, Int) Authentication method, 0 - local, 1 - LDAP, 2 - OAuth. If not passed, the default is 0.
* `department_id` - (Optional, String) Department ID, such as: 1.2.3.
* `email` - (Optional, String) Email.
* `group_id_set` - (Optional, Set: [`Int`]) The set of user group IDs to which it belongs.
* `phone` - (Optional, String) Fill in the mainland mobile phone number directly. If it is a number from other countries or regions, enter it in the format of country area code|mobile phone number. For example: +852|xxxxxxxx.
* `validate_from` - (Optional, String) User effective time, such as: 2021-09-22T00:00:00+00:00If the effective and expiry time are not filled in, the user will be valid for a long time.
* `validate_time` - (Optional, String) Access time period limit, a string composed of 0 and 1, length 168 (7 * 24), representing the time period the user is allowed to access in a week. The Nth character in the string represents the Nth hour of the week, 0 - means access is not allowed, 1 - means access is allowed.
* `validate_to` - (Optional, String) User expiration time, such as: 2021-09-23T00:00:00+00:00If the effective and expiry time are not filled in, the user will be valid for a long time.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

dasb user can be imported using the id, e.g.

```
terraform import tencentcloud_dasb_user.example 134
```

