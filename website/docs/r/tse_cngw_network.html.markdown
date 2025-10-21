---
subcategory: "Tencent Cloud Service Engine(TSE)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tse_cngw_network"
sidebar_current: "docs-tencentcloud-resource-tse_cngw_network"
description: |-
  Provides a resource to create a tse cngw_network
---

# tencentcloud_tse_cngw_network

Provides a resource to create a tse cngw_network

## Example Usage

```hcl
resource "tencentcloud_tse_cngw_network" "cngw_network" {
  gateway_id                 = "gateway-cf8c99c3"
  group_id                   = "group-a160d123"
  internet_address_version   = "IPV4"
  internet_pay_mode          = "BANDWIDTH"
  description                = "des-test1"
  internet_max_bandwidth_out = 1
  master_zone_id             = "ap-guangzhou-3"
  multi_zone_flag            = true
  sla_type                   = "clb.c2.medium"
  slave_zone_id              = "ap-guangzhou-4"
}
```

## Argument Reference

The following arguments are supported:

* `gateway_id` - (Required, String, ForceNew) gateway ID.
* `group_id` - (Required, String, ForceNew) gateway group ID.
* `description` - (Optional, String) description of clb.
* `internet_address_version` - (Optional, String) internet type. Reference value:`IPV4` (default value), `IPV6`.
* `internet_max_bandwidth_out` - (Optional, Int) public network bandwidth.
* `internet_pay_mode` - (Optional, String) trade type of internet. Reference value:`BANDWIDTH` (default value), `TRAFFIC`.
* `master_zone_id` - (Optional, String) primary availability zone.
* `multi_zone_flag` - (Optional, Bool) Whether load balancing has multiple availability zones.
* `sla_type` - (Optional, String) specification type of clb. Default `shared` type when this parameter is empty, Note: input `shared` is not supported when creating. Reference value:`clb.c2.medium`, `clb.c3.small`, `clb.c3.medium`, `clb.c4.small`, `clb.c4.medium`, `clb.c4.large`, `clb.c4.xlarge`.
* `slave_zone_id` - (Optional, String) alternate availability zone.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `network_id` - network id.
* `vip` - clb vip.


