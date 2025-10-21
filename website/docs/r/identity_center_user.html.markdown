---
subcategory: "Tencent Cloud Organization (TCO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_identity_center_user"
sidebar_current: "docs-tencentcloud-resource-identity_center_user"
description: |-
  Provides a resource to create an identity center user
---

# tencentcloud_identity_center_user

Provides a resource to create an identity center user

## Example Usage

```hcl
resource "tencentcloud_identity_center_user" "example" {
  zone_id     = "z-1os7c9tyugct"
  user_name   = "tf-example"
  description = "desc."
}
```

### Or

```hcl
resource "tencentcloud_identity_center_user" "example" {
  zone_id      = "z-1os7c9tyugct"
  user_name    = "tf-example"
  description  = "desc."
  first_name   = "FirstName"
  last_name    = "LastName"
  display_name = "DisplayName"
  email        = "example@tencent.com"
  user_status  = "Enabled"
}
```

## Argument Reference

The following arguments are supported:

* `user_name` - (Required, String) User name. It must be unique in space. Modifications are not supported. Format: Contains numbers, English letters and special symbols(`+`, `=`, `,`, `.`, `@`, `-`, `_`). Length: Maximum 64 characters.
* `zone_id` - (Required, String) Zone id.
* `description` - (Optional, String) User's description. Length: Maximum 1024 characters.
* `display_name` - (Optional, String) The display name of the user. Length: Maximum 256 characters.
* `email` - (Optional, String) The user's email address. Must be unique within the catalog. Length: Maximum 128 characters.
* `first_name` - (Optional, String) The user's last name. Length: Maximum 64 characters.
* `last_name` - (Optional, String) The user's name. Length: Maximum 64 characters.
* `user_status` - (Optional, String) The status of the user. Value: Enabled (default): Enabled. Disabled: Disabled.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Create time.
* `update_time` - Update time.
* `user_id` - User id.
* `user_type` - User type.


## Import

organization identity center user can be imported using the id, e.g.

```
terraform import tencentcloud_identity_center_user.example z-1os7c9tyugct#u-rdvm4xdqi8pr
```

