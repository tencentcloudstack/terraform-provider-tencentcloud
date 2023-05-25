---
subcategory: "Direct Connect(DC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dcx_extra_config"
sidebar_current: "docs-tencentcloud-resource-dcx_extra_config"
description: |-
  Provides a resource to create a dc dcx_extra_config
---

# tencentcloud_dcx_extra_config

Provides a resource to create a dc dcx_extra_config

## Example Usage

```hcl
resource "tencentcloud_dcx_extra_config" "dcx_extra_config" {
  direct_connect_tunnel_id = "dcx-4z49tnws"
  vlan                     = 123
  bgp_peer {
    asn      = 65101
    auth_key = "test123"

  }
  route_filter_prefixes {
    cidr = "192.168.0.0/24"
  }
  tencent_address        = "192.168.1.1"
  tencent_backup_address = "192.168.1.2"
  customer_address       = "192.168.1.4"
  bandwidth              = 10
  enable_bgp_community   = false
  bfd_enable             = 0
  nqa_enable             = 1
  bfd_info {
    probe_failed_times = 3
    interval           = 100

  }
  nqa_info {
    probe_failed_times = 3
    interval           = 100
    destination_ip     = "192.168.2.2"

  }
  ipv6_enable  = 0
  jumbo_enable = 0
}
```

## Argument Reference

The following arguments are supported:

* `direct_connect_tunnel_id` - (Required, String) direct connect tunnel id.
* `bandwidth` - (Optional, Int) direct connect tunnel bandwidth.
* `bfd_enable` - (Optional, Int) be enabled BFD.
* `bfd_info` - (Optional, List) BFD config info.
* `bgp_peer` - (Optional, List) idc BGP, Asn, AuthKey.
* `customer_address` - (Optional, String) direct connect tunnel user idc connect ip.
* `enable_bgp_community` - (Optional, Bool) BGP community attribute.
* `ipv6_enable` - (Optional, Int) 0: disable IPv61: enable IPv6.
* `jumbo_enable` - (Optional, Int) direct connect tunnel support jumbo frame1: enable direct connect tunnel jumbo frame0: disable direct connect tunnel jumbo frame.
* `nqa_enable` - (Optional, Int) be enabled NQA.
* `nqa_info` - (Optional, List) NQA config info.
* `route_filter_prefixes` - (Optional, List) user filter network prefixes.
* `tencent_address` - (Optional, String) direct connect tunnel tencent cloud connect ip.
* `tencent_backup_address` - (Optional, String) direct connect tunnel tencent cloud backup connect ip.
* `vlan` - (Optional, Int) direct connect tunnel vlan id.

The `bfd_info` object supports the following:

* `interval` - (Optional, Int) detect interval.
* `probe_failed_times` - (Optional, Int) detect times.

The `bgp_peer` object supports the following:

* `asn` - (Optional, Int) user idc BGP Asn.
* `auth_key` - (Optional, String) user bgp key.

The `nqa_info` object supports the following:

* `destination_ip` - (Optional, String) detect ip.
* `interval` - (Optional, Int) detect interval.
* `probe_failed_times` - (Optional, Int) detect times.

The `route_filter_prefixes` object supports the following:

* `cidr` - (Optional, String) user network prefixes.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

dc dcx_extra_config can be imported using the id, e.g.

```
terraform import tencentcloud_dcx_extra_config.dcx_extra_config dcx_id
```

