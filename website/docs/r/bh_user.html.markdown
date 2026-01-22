---
subcategory: "Bastion Host(BH)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_bh_user"
sidebar_current: "docs-tencentcloud-resource-bh_user"
description: |-
  Provides a resource to create a BH user
---

# tencentcloud_bh_user

Provides a resource to create a BH user

## Example Usage

```hcl
resource "tencentcloud_bh_user" "example" {
  user_name = "tf-example"
  real_name = "Terraform"
  phone     = "+86|18991162528"
  email     = "demo@tencent.com"
  auth_type = 0
}
```

## Argument Reference

The following arguments are supported:

* `real_name` - (Required, String) User's real name, maximum length 20 characters, cannot contain whitespace characters.
* `user_name` - (Required, String, ForceNew) Username, 3-20 characters, must start with an English letter and cannot contain characters other than `letters`, `numbers`, `.`, `_`, `-`.
* `auth_type` - (Optional, Int) Authentication method, 0 - local, 1 - LDAP, 2 - OAuth. Default is 0 if not provided.
* `department_id` - (Optional, String) Department ID to which the user belongs, e.g.: "1.2.3".
* `email` - (Optional, String) Email address. At least one of phone and email parameters must be provided.
* `group_id_set` - (Optional, Set: [`Int`]) User group ID set to which the user belongs.
* `phone` - (Optional, String) Input in the format of "country code|phone number", e.g.: "+86|xxxxxxxx". At least one of phone and email parameters must be provided.
* `validate_from` - (Optional, String) User effective time, e.g.: "2021-09-22T00:00:00+00:00". If effective and expiration times are not filled, the user will be valid permanently.
* `validate_time` - (Optional, String) Access time restriction, a string composed of 0 and 1 with length 168 (7 * 24), representing the time slots allowed for the user in a week. The Nth character in the string represents the Nth hour in the week, 0 - not allowed to access, 1 - allowed to access.
* `validate_to` - (Optional, String) User expiration time, e.g.: "2021-09-23T00:00:00+00:00". If effective and expiration times are not filled, the user will be valid permanently.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `user_id` - User ID.


## Import

BH user can be imported using the id, e.g.

```
terraform import tencentcloud_bh_user.example 2322
```

