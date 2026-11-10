---
subcategory: "Provider Meta"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cam_update_password_rules"
sidebar_current: "docs-tencentcloud-resource-cam_update_password_rules"
description: |-
  Provides a resource to manage the CAM password rules configuration (UpdatePasswordRules / GetPasswordRules).
---

# tencentcloud_cam_update_password_rules

Provides a resource to manage the CAM password rules configuration (UpdatePasswordRules / GetPasswordRules).

This is a singleton, account-wide configuration resource built on the terraform-plugin-framework implementation. Create and update both call the `UpdatePasswordRules` API; read calls `GetPasswordRules`. Destroying this resource only removes it from state, the cloud-side password rules are left unchanged.

-> **Note:** The CAM password rules are a single account-wide configuration. If the rules look like they have already been changed outside of Terraform, the create emits a warning and applying this configuration will overwrite them; run `terraform import tencentcloud_cam_update_password_rules.example cam-password-rules` and inspect the plan first if the existing values should be kept. Every attribute inside `rules` is optional; a value that is left unset keeps its current cloud-side setting.

-> **Note:** `must_contain` is a literal character set, not an enum: write the characters the password must contain, e.g. `"Aa1!"` (upper, lower, digit, special). The console re-orders and normalizes this value (e.g. `"Aa1!"` is stored and read back as `"1!aA"`), and invalid literals are rejected by the API with `InvalidParameter.PasswordRuleError`.

## Example Usage

```hcl
resource "tencentcloud_cam_update_password_rules" "example" {
  rules {
    minimum_length                   = 8
    must_contain                     = "Aa1!"
    force_password_change            = 0
    reuse_password_limit             = 1
    retry_password_limit             = 10
    password_expiration_invalidation = 2
  }
}
```

## Argument Reference

The following arguments are supported:

* `rules` - (Required, Object) Password rules.

The `rules` object supports the following:

* `force_password_change` - (Optional, Int) Forced password change period.
* `minimum_length` - (Optional, Int) Minimum password length.
* `must_contain` - (Optional, String) Characters that the password must contain.
* `password_expiration_invalidation` - (Optional, Int) Whether an expired password becomes invalid. Valid values: `1` (the sub-user cannot log in and an administrator must reset the password), `2` (the sub-user can log in and is forced to change the password).
* `retry_password_limit` - (Optional, Int) Number of allowed password retries.
* `reuse_password_limit` - (Optional, Int) Number of previous passwords that cannot be reused.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the password rules configuration, always `cam-password-rules` because the CAM password rules are a single account-wide configuration.
* `modifier` - Last modifier of the password rules.
* `update_time` - Last modification time of the password rules.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create:
  - [GetPasswordRules](https://www.tencentcloud.com/document/api/598/81237)
  - [UpdatePasswordRules](https://www.tencentcloud.com/document/api/598/81238)
* Read: [GetPasswordRules](https://www.tencentcloud.com/document/api/598/81237)
* Update:
  - [GetPasswordRules](https://www.tencentcloud.com/document/api/598/81237)
  - [UpdatePasswordRules](https://www.tencentcloud.com/document/api/598/81238)

## Import

CAM password rules can be imported using the constant id `cam-password-rules`, e.g.

```
terraform import tencentcloud_cam_update_password_rules.example cam-password-rules
```

