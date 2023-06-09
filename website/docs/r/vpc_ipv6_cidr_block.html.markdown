---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc_ipv6_cidr_block"
sidebar_current: "docs-tencentcloud-resource-vpc_ipv6_cidr_block"
description: |-
  Provides a resource to create a vpc ipv6_cidr_block
---

# tencentcloud_vpc_ipv6_cidr_block

Provides a resource to create a vpc ipv6_cidr_block

## Example Usage

```hcl
resource "tencentcloud_vpc" "cidr-block" {
  name         = "ipv6-cidr-block-for-test"
  cidr_block   = "10.0.0.0/16"
  is_multicast = false
}

resource "tencentcloud_vpc_ipv6_cidr_block" "ipv6_cidr_block" {
  vpc_id = tencentcloud_vpc.cidr-block.id
}
```

## Argument Reference

The following arguments are supported:

* `vpc_id` - (Required, String, ForceNew) `VPC` instance `ID`, in the form of `vpc-f49l6u0z`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `ipv6_cidr_block` - ipv6 cidr block.


## Import

vpc ipv6_cidr_block can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_ipv6_cidr_block.ipv6_cidr_block vpc_id
```

