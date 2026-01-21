---
subcategory: "Cwp"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cwp_auto_open_proversion_config"
sidebar_current: "docs-tencentcloud-resource-cwp_auto_open_proversion_config"
description: |-
  Provides a resource to create a CWP auto open proversion config
---

# tencentcloud_cwp_auto_open_proversion_config

Provides a resource to create a CWP auto open proversion config

## Example Usage

```hcl
resource "tencentcloud_cwp_auto_open_proversion_config" "example" {
  status                       = "OPEN"
  protect_type                 = "FLAGSHIP_PREPAY"
  auto_repurchase_renew_switch = 1
  auto_repurchase_switch       = 1
  repurchase_renew_switch      = 0
  auto_bind_rasp_switch        = 1
  auto_open_rasp_switch        = 1
  auto_downgrade_switch        = 0
}
```

## Argument Reference

The following arguments are supported:

* `status` - (Required, String) Set the auto-activation status.
<li>CLOSE: off</li>
<li>OPEN: on</li>.
* `auto_bind_rasp_switch` - (Optional, Int) Newly added machines will be automatically bound to Rasp. 0: Disabled, 1: Enabled.
* `auto_downgrade_switch` - (Optional, Int) Automatic scaling switch: 0 for off, 1 for on.
* `auto_open_rasp_switch` - (Optional, Int) Newly added machines will have automatic Raspberry Pi protection enabled by default. (0: Disabled, 1: Enabled).
* `auto_repurchase_renew_switch` - (Optional, Int) Auto-renewal or not for auto-purchased orders, 0 by default, 0 for OFF, 1 for ON.
* `auto_repurchase_switch` - (Optional, Int) Automatic purchase/expansion authorization switch, 1 by default, 0 for OFF, 1 for ON.
* `protect_type` - (Optional, String) Enhanced Protection Mode PROVERSION_POSTPAY Professional Edition - Pay-as-you-go PROVERSION_PREPAY Professional Edition - Annual/Monthly Subscription FLAGSHIP_PREPAY Flagship Edition - Annual/Monthly Subscription.
* `repurchase_renew_switch` - (Optional, Int) Whether the manually purchased order is automatically renewed (defaults to 0). 0 - off; 1 -on.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

CWP auto open proversion config can be imported using the customId(like uuid or base64 string), e.g.

```
terraform import tencentcloud_cwp_auto_open_proversion_config.example 9n0PnuKR3sEmY8COuKVb7g==
```

