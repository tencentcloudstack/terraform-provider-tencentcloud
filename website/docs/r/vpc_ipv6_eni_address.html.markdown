---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc_ipv6_eni_address"
sidebar_current: "docs-tencentcloud-resource-vpc_ipv6_eni_address"
description: |-
  Provides a resource to create a vpc ipv6_eni_address
---

# tencentcloud_vpc_ipv6_eni_address

Provides a resource to create a vpc ipv6_eni_address

## Example Usage

```hcl
resource "tencentcloud_vpc_ipv6_eni_address" "ipv6_eni_address" {
  vpc_id               = "vpc-7w3kgnpl"
  network_interface_id = "eni-pzl7fz37"
  ipv6_addresses {
    address     = "2402:4e00:1019:6a7b:0:994e:7cd7:fb46"
    description = "test123"
  }
}
```

## Argument Reference

The following arguments are supported:

* `network_interface_id` - (Required, String) ENI instance `ID`, in the form of `eni-m6dyj72l`.
* `vpc_id` - (Required, String) VPC `ID`, in the form of `vpc-m6dyj72l`.
* `ipv6_addresses` - (Optional, List) The specified `IPv6` address list, up to 10 can be specified at a time. Combined with the input parameter `Ipv6AddressCount` to calculate the quota. Mandatory one with Ipv6AddressCount.

The `ipv6_addresses` object supports the following:

* `address` - (Required, String) `IPv6` address, in the form of: `3402:4e00:20:100:0:8cd9:2a67:71f3`.
* `address_id` - (Optional, String) `EIP` instance `ID`, such as:`eip-hxlqja90`.
* `description` - (Optional, String) Description.
* `is_wan_ip_blocked` - (Optional, Bool) Whether the public network IP is blocked.
* `primary` - (Optional, Bool) Whether to master `IP`.
* `state` - (Optional, String) `IPv6` address status: `PENDING`: pending, `MIGRATING`: migrating, `DELETING`: deleting, `AVAILABLE`: available.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



