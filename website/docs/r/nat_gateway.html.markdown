---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_nat_gateway"
sidebar_current: "docs-tencentcloud-resource-vpc-nat-gateway"
description: |-
  Provides a resource to create a VPC NAT Gateway.
---

# tencentcloud_nat_gateway

Provides a resource to create a VPC NAT Gateway.

## Example Usage

Basic usage:

```hcl
resource "tencentcloud_vpc" "main" {
  name       = "terraform test"
  cidr_block = "10.6.0.0/16"
}

# Create EIP
resource "tencentcloud_eip" "eip_dev_dnat" {
  name = "terraform_test"
}
resource "tencentcloud_eip" "eip_test_dnat" {
  name = "terraform_test"
}

# Create NAT Gateway
resource "tencentcloud_nat_gateway" "my_nat" {
  vpc_id           = "${tencentcloud_vpc.main.id}"
  name             = "terraform test"
  max_concurrent   = 3000000
  bandwidth        = 500
  assigned_eip_set = [
    "${tencentcloud_eip.eip_dev_dnat.public_ip}",
    "${tencentcloud_eip.eip_test_dnat.public_ip}",
  ]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name for the NAT Gateway.
* `vpc_id` - (Required, Forces new resource) The VPC ID.
* `max_concurrent` - (Required) The upper limit of concurrent connection of NAT gateway, for example: 1000000, 3000000, 10000000. To learn more, please refer to [Virtual Private Cloud Gateway Description](https://intl.cloud.tencent.com/doc/product/215/1682).
* `bandwidth` - (Required) The maximum public network output bandwidth of the gateway (unit: Mbps), for example: 10, 20, 50, 100, 200, 500, 1000, 2000, 5000. For more information, please refer to [Virtual Private Cloud Gateway Description](https://intl.cloud.tencent.com/doc/product/215/1682).
* `assigned_eip_set` - (Required) Elastic IP arrays bound to the gateway, For more information on elastic IP, please refer to [Elastic IP](eip.html).

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the NAT Gateway.
* `name` - The name of the NAT Gateway.
* `max_concurrent` - The upper limit of concurrent connection of NAT gateway.
* `bandwidth` - The maximum public network output bandwidth of the gateway (unit: Mbps).
* `assigned_eip_set` - Elastic IP arrays bound to the gateway
