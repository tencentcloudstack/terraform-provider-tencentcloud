---
subcategory: "Virtual Private Cloud(VPC)"
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
  max_concurrent   = 1000000
  assigned_eip_set = ["1.1.1.1"]

  tags = {
    test = "tf"
  }
}
```

## Argument Reference

The following arguments are supported:

* `assigned_eip_set` - (Required) EIP IP address set bound to the gateway. The value of at least 1 and at most 10.
* `name` - (Required) Name of the NAT gateway.
* `vpc_id` - (Required, ForceNew) ID of the vpc.
* `bandwidth` - (Optional) The maximum public network output bandwidth of NAT gateway (unit: Mbps). Valid values: 20,50,100,200,500,1000,2000,5000. Default is 100.
* `max_concurrent` - (Optional) The upper limit of concurrent connection of NAT gateway. Valid values: 1000000,3000000,10000000. Default is 1000000.
* `tags` - (Optional) The available tags within this NAT gateway.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `created_time` - Create time of the NAT gateway.


## Import

NAT gateway can be imported using the id, e.g.

```
$ terraform import tencentcloud_nat_gateway.foo nat-1asg3t63
```

