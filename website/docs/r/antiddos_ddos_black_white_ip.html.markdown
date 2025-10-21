---
subcategory: "Anti-DDoS(DayuV2)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_antiddos_ddos_black_white_ip"
sidebar_current: "docs-tencentcloud-resource-antiddos_ddos_black_white_ip"
description: |-
  Provides a resource to create a antiddos ddos black white ip
---

# tencentcloud_antiddos_ddos_black_white_ip

Provides a resource to create a antiddos ddos black white ip

## Example Usage

```hcl
resource "tencentcloud_antiddos_ddos_black_white_ip" "ddos_black_white_ip" {
  instance_id = "bgp-xxxxxx"
  ip          = "1.2.3.5"
  mask        = 0
  type        = "black"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) instance id.
* `ip` - (Required, String) ip list.
* `mask` - (Required, Int) ip mask.
* `type` - (Required, String) ip type, black: black ip list, white: white ip list.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

antiddos ddos_black_white_ip can be imported using the id, e.g.

```
terraform import tencentcloud_antiddos_ddos_black_white_ip.ddos_black_white_ip ${instanceId}#${ip}
```

