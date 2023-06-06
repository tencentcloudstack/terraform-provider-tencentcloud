---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc_ipv6_subnet_cidr_block"
sidebar_current: "docs-tencentcloud-resource-vpc_ipv6_subnet_cidr_block"
description: |-
  Provides a resource to create a vpc ipv6_subnet_cidr_block
---

# tencentcloud_vpc_ipv6_subnet_cidr_block

Provides a resource to create a vpc ipv6_subnet_cidr_block

## Example Usage

```hcl
resource "tencentcloud_vpc_ipv6_subnet_cidr_block" "ipv6_subnet_cidr_block" {
  vpc_id = "vpc-7w3kgnpl"
  ipv6_subnet_cidr_blocks {
    subnet_id       = "subnet-plg028y8"
    ipv6_cidr_block = "2402:4e00:1019:6a7b::/64"
  }
}
```

## Argument Reference

The following arguments are supported:

* `ipv6_subnet_cidr_blocks` - (Required, List, ForceNew) Allocate a list of `IPv6` subnets.
* `vpc_id` - (Required, String, ForceNew) The private network `ID` where the subnet is located. Such as:`vpc-f49l6u0z`.

The `ipv6_subnet_cidr_blocks` object supports the following:

* `ipv6_cidr_block` - (Required, String) `IPv6` subnet segment. Such as: `3402:4e00:20:1001::/64`.
* `subnet_id` - (Required, String) Subnet instance `ID`. Such as:`subnet-pxir56ns`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

vpc ipv6_subnet_cidr_block can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_ipv6_subnet_cidr_block.ipv6_subnet_cidr_block ipv6_subnet_cidr_block_id
```

