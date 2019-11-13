---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_nat_gateway"
sidebar_current: "docs-tencentcloud-resource-nat_gateway"
description: |-
  Provides a resource to create a NAT gateway.
---

# tencentcloud_nat_gateway

Provides a resource to create a NAT gateway.

## Example Usage

```hcl
resource "tencentcloud_nat_gateway" "foo" {
  name             = "test_nat_gateway"
  vpc_id           = "vpc-4xxr2cy7"
  bandwidth        = 100
  max_connection   = 1000000
  assigned_eip_set = ["eip-da12w5re5"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the NAT gateway.
* `vpc_id` - (Required, ForceNew) Id of the VPC.
* `assigned_eip_set` - (Optional) EIP arrays bound to the gateway. The value of at least 1.
* `bandwidth` - (Optional) The maximum public network output bandwidth of nat gateway (unit: Mbps), the available values include: 20,50,100,200,500,1000,2000,5000. Default is 100.
* `max_concurrent` - (Optional) The upper limit of concurrent connection of nat gateway, the available values include: 1000000,3000000,10000000. Default is 1000000.


## Import

NAT gateway can be imported using the id, e.g.

```
$ terraform import tencentcloud_nat_gateway.foo nat-1asg3t63
```

