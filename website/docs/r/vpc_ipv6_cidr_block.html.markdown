---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc_ipv6_cidr_block"
sidebar_current: "docs-tencentcloud-resource-vpc_ipv6_cidr_block"
description: |-
  Provides a resource to create a VPC ipv6 cidr block
---

# tencentcloud_vpc_ipv6_cidr_block

Provides a resource to create a VPC ipv6 cidr block

## Example Usage

```hcl
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_vpc_ipv6_cidr_block" "example" {
  vpc_id = tencentcloud_vpc.vpc.id
}
```

### Or

```hcl
resource "tencentcloud_vpc_ipv6_cidr_block" "example" {
  vpc_id       = tencentcloud_vpc.vpc.id
  address_type = "ULA"
}
```

## Argument Reference

The following arguments are supported:

* `vpc_id` - (Required, String, ForceNew) `VPC` instance `ID`, in the form of `vpc-f49l6u0z`.
* `address_type` - (Optional, String, ForceNew) Apply for the type of IPv6 Cidr, GUA (Global Unicast Address), ULA (Unique Local Address).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `ipv6_cidr_block_set` - Ipv6 cidr block set.
  * `address_type` - Apply for the type of IPv6 Cidr, GUA (Global Unicast Address), ULA (Unique Local Address).
  * `ipv6_cidr_block` - Ipv6 cidr block.
  * `isp_type` - Range of network operator types: 'BGP' - default, 'CMCC' - China Mobile, 'CTCC' - China Telecom, 'CUCC' - China Joint Debugging.
* `ipv6_cidr_block` - Ipv6 cidr block.


## Import

vpc ipv6_cidr_block can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_ipv6_cidr_block.example vpc-826mi3hd
```

