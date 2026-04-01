---
subcategory: "Bastion Host(BH)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_bh_auth_mode_config"
sidebar_current: "docs-tencentcloud-resource-bh_auth_mode_config"
description: |-
  Provides a resource to create a BH auth mode config
---

# tencentcloud_bh_auth_mode_config

Provides a resource to create a BH auth mode config

## Example Usage

```hcl
resource "tencentcloud_bh_auth_mode_config" "example" {
  auth_mode = 2
}
```

## Argument Reference

The following arguments are supported:

* `auth_mode_gm` - (Optional, Int) National secret double factor authentication configuration. Valid values: `0` (disabled), `1` (OTP), `2` (SMS), `3` (USB Key). Note: At least one of AuthMode and AuthModeGM must be passed. AuthModeGM has higher priority than ResourceType.
* `auth_mode` - (Optional, Int) Double factor authentication configuration. Valid values: `0` (disabled), `1` (OTP), `2` (SMS), `3` (USB Key, only valid when ResourceType=1 and AuthModeGM is not passed). Note: At least one of AuthMode and AuthModeGM must be passed.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



