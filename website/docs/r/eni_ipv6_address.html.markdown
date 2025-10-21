---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_eni_ipv6_address"
sidebar_current: "docs-tencentcloud-resource-eni_ipv6_address"
description: |-
  Provides a resource to create a vpc eni ipv6 address
---

# tencentcloud_eni_ipv6_address

Provides a resource to create a vpc eni ipv6 address

## Example Usage

```hcl
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = "ap-guangzhou-6"
  name              = "subnet-example"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.0.0/16"
  is_multicast      = false
}

resource "tencentcloud_vpc_ipv6_cidr_block" "example" {
  vpc_id = tencentcloud_vpc.vpc.id
}

resource "tencentcloud_vpc_ipv6_subnet_cidr_block" "example" {
  vpc_id = tencentcloud_vpc.vpc.id
  ipv6_subnet_cidr_blocks {
    subnet_id       = tencentcloud_subnet.subnet.id
    ipv6_cidr_block = "2402:4e00:1015:7500::/64"
  }
}

resource "tencentcloud_eni" "example" {
  name        = "tf-example"
  vpc_id      = tencentcloud_vpc.vpc.id
  subnet_id   = tencentcloud_subnet.subnet.id
  description = "eni description."
  ipv4_count  = 1

  depends_on = [
    tencentcloud_vpc_ipv6_cidr_block.example,
    tencentcloud_vpc_ipv6_subnet_cidr_block.example
  ]
}

resource "tencentcloud_eni_ipv6_address" "example" {
  network_interface_id = tencentcloud_eni.example.id
  ipv6_address_count   = 1
}
```

## Argument Reference

The following arguments are supported:

* `network_interface_id` - (Required, String, ForceNew) ENI instance `ID`, in the form of `eni-m6dyj72l`.
* `ipv6_address_count` - (Optional, Int, ForceNew) The number of automatically assigned IPv6 addresses and the total number of private IP addresses cannot exceed the quota. This should be combined with the input parameter `ipv6_addresses` for quota calculation. At least one of them, either this or 'Ipv6Addresses', must be provided.
* `ipv6_addresses` - (Optional, Set, ForceNew) The specified `IPv6` address list, up to 10 can be specified at a time. Combined with the input parameter `Ipv6AddressCount` to calculate the quota. Mandatory one with Ipv6AddressCount.

The `ipv6_addresses` object supports the following:

* `address` - (Required, String, ForceNew) `IPv6` address, in the form of: `3402:4e00:20:100:0:8cd9:2a67:71f3`.
* `address_id` - (Optional, String, ForceNew) `EIP` instance `ID`, such as:`eip-hxlqja90`.
* `description` - (Optional, String, ForceNew) Description.
* `is_wan_ip_blocked` - (Optional, Bool, ForceNew) Whether the public network IP is blocked.
* `primary` - (Optional, Bool, ForceNew) Whether to master `IP`.
* `state` - (Optional, String, ForceNew) `IPv6` address status: `PENDING`: pending, `MIGRATING`: migrating, `DELETING`: deleting, `AVAILABLE`: available.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

vpc eni ipv6 address can be imported using the id, e.g.

```
terraform import tencentcloud_eni_ipv6_address.example eni-fxrx5d1d
```

