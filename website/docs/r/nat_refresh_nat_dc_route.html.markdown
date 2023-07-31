---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_nat_refresh_nat_dc_route"
sidebar_current: "docs-tencentcloud-resource-nat_refresh_nat_dc_route"
description: |-
  Provides a resource to create a vpc refresh_nat_dc_route
---

# tencentcloud_nat_refresh_nat_dc_route

Provides a resource to create a vpc refresh_nat_dc_route

## Example Usage

### is True

```hcl
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_eip" "eip_example" {
  name = "eip_example"
}

resource "tencentcloud_nat_gateway" "nat" {
  vpc_id         = tencentcloud_vpc.vpc.id
  name           = "tf_example_nat_gateway"
  max_concurrent = 3000000
  bandwidth      = 500

  assigned_eip_set = [
    tencentcloud_eip.eip_example.public_ip,
  ]
}

resource "tencentcloud_nat_refresh_nat_dc_route" "refresh_nat_dc_route" {
  nat_gateway_id = tencentcloud_nat_gateway.nat.id
  vpc_id         = tencentcloud_vpc.vpc.id
  dry_run        = true
}
```

### is False

```hcl
resource "tencentcloud_nat_refresh_nat_dc_route" "refresh_nat_dc_route" {
  nat_gateway_id = tencentcloud_nat_gateway.nat.id
  vpc_id         = tencentcloud_vpc.vpc.id
  dry_run        = false
}
```

## Argument Reference

The following arguments are supported:

* `dry_run` - (Required, Bool, ForceNew) Whether to pre-refresh, valid values: True:yes, False:no.
* `nat_gateway_id` - (Required, String, ForceNew) Unique identifier of Nat Gateway.
* `vpc_id` - (Required, String, ForceNew) Unique identifier of Vpc.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

vpc refresh_nat_dc_route can be imported using the id, e.g.

```
terraform import tencentcloud_nat_refresh_nat_dc_route.refresh_nat_dc_route vpc_id#nat_gateway_id
```

