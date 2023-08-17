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
data "tencentcloud_availability_zones" "zones" {}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = data.tencentcloud_availability_zones.zones.zones.0.name
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
    ipv6_cidr_block = tencentcloud_vpc_ipv6_cidr_block.example.ipv6_cidr_block
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

