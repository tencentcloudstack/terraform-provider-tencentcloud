---
subcategory: "Bastion Host(BH)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_bh_asset_sync_flag_config"
sidebar_current: "docs-tencentcloud-resource-bh_asset_sync_flag_config"
description: |-
  Provides a resource to create a BH asset sync flag config
---

# tencentcloud_bh_asset_sync_flag_config

Provides a resource to create a BH asset sync flag config

## Example Usage

```hcl
resource "tencentcloud_bh_asset_sync_flag_config" "example" {
  auto_sync = true
}
```

## Argument Reference

The following arguments are supported:

* `auto_sync` - (Required, Bool) Whether to enable asset auto-sync, false - disabled, true - enabled.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `role_granted` - Whether the role has been authorized, false - not authorized, true - authorized.


## Import

BH asset sync flag config can be imported using the customId(like uuid or base64 string), e.g.

```
terraform import tencentcloud_bh_asset_sync_flag_config.example yci5a1o76a5HzCaqJM2bQA==
```

