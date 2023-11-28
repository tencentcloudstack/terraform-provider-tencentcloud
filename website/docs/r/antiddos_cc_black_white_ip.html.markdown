---
subcategory: "Anti-DDoS(DayuV2)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_antiddos_cc_black_white_ip"
sidebar_current: "docs-tencentcloud-resource-antiddos_cc_black_white_ip"
description: |-
  Provides a resource to create a antiddos cc black white ip
---

# tencentcloud_antiddos_cc_black_white_ip

Provides a resource to create a antiddos cc black white ip

## Example Usage

```hcl
resource "tencentcloud_antiddos_cc_black_white_ip" "cc_black_white_ip" {
  instance_id = "bgpip-xxxxxx"
  black_white_ip {
    ip   = "1.2.3.5"
    mask = 0

  }
  type     = "black"
  ip       = "xxx.xxx.xxx.xxx"
  domain   = "t.baidu.com"
  protocol = "http"
}
```

## Argument Reference

The following arguments are supported:

* `black_white_ip` - (Required, List, ForceNew) Black white ip.
* `domain` - (Required, String, ForceNew) domain.
* `instance_id` - (Required, String, ForceNew) instance id.
* `ip` - (Required, String, ForceNew) ip address.
* `protocol` - (Required, String, ForceNew) protocol.
* `type` - (Required, String, ForceNew) IP type, value [black(blacklist IP), white(whitelist IP)].

The `black_white_ip` object supports the following:

* `ip` - (Required, String) ip address.
* `mask` - (Required, Int) ip mask.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

antiddos cc_black_white_ip can be imported using the id, e.g.

```
terraform import tencentcloud_antiddos_cc_black_white_ip.cc_black_white_ip ${instanceId}#${policyId}#${instanceIp}#${domain}#${protocol}
```

