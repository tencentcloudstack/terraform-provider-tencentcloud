---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc_private_nat_gateway"
sidebar_current: "docs-tencentcloud-resource-vpc_private_nat_gateway"
description: |-
  Provides a resource to create a vpc private nat gateway
---

# tencentcloud_vpc_private_nat_gateway

Provides a resource to create a vpc private nat gateway

## Example Usage

```hcl
resource "tencentcloud_vpc_private_nat_gateway" "private_nat_gateway" {
  nat_gateway_name = "xxx"
  vpc_id           = "xxx"
}
```

## Argument Reference

The following arguments are supported:

* `nat_gateway_name` - (Required, String) Private network gateway name.
* `ccn_id` - (Optional, String) Cloud Connect Network type The Cloud Connect Network instance ID required to be bound to the private network NAT gateway.
* `cross_domain` - (Optional, Bool) Cross-domain parameters. Cross-domain binding of VPCs is supported only when the value is True.
* `vpc_id` - (Optional, String) Private Cloud instance ID. This parameter is required when creating a VPC type private network NAT gateway or a private network NAT gateway of private network gateway.
* `vpc_type` - (Optional, Bool) VPC type private network NAT gateway. Only when the value is True will a VPC type private network NAT gateway be created.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

vpc private_nat_gateway can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_private_nat_gateway.private_nat_gateway private_nat_gateway_id
```

