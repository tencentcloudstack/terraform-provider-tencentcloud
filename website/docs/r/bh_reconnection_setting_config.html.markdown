---
subcategory: "Bastion Host(BH)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_bh_reconnection_setting_config"
sidebar_current: "docs-tencentcloud-resource-bh_reconnection_setting_config"
description: |-
  Provides a resource to create a BH reconnection setting config
---

# tencentcloud_bh_reconnection_setting_config

Provides a resource to create a BH reconnection setting config

## Example Usage

```hcl
resource "tencentcloud_bh_reconnection_setting_config" "example" {
  reconnection_max_count = 5
  enable                 = true
}
```

## Argument Reference

The following arguments are supported:

* `enable` - (Required, Bool) true: limit reconnection count, false: do not limit reconnection count.
* `reconnection_max_count` - (Required, Int) Retry count, value range: 0-20.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

BH reconnection setting config can be imported using the customId(like uuid or base64 string), e.g.

```
terraform import tencentcloud_bh_reconnection_setting_config.example gO1Ew6OEgLcQun164XiWmw==
```

