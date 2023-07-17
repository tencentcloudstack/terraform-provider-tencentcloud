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

### Create a NAT gateway.

```hcl
resource "tencentcloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
  name       = "tf_nat_gateway_vpc"
}

resource "tencentcloud_eip" "eip_example1" {
  name = "tf_nat_gateway_eip1"
}

resource "tencentcloud_eip" "eip_example2" {
  name = "tf_nat_gateway_eip2"
}

resource "tencentcloud_nat_gateway" "example" {
  name           = "tf_example_nat_gateway"
  vpc_id         = tencentcloud_vpc.vpc.id
  bandwidth      = 100
  max_concurrent = 1000000
  assigned_eip_set = [
    tencentcloud_eip.eip_example1.public_ip,
    tencentcloud_eip.eip_example2.public_ip,
  ]
  tags = {
    tf_tag_key = "tf_tag_value"
  }
}
```

## Argument Reference

The following arguments are supported:

* `assigned_eip_set` - (Required, Set: [`String`]) EIP IP address set bound to the gateway. The value of at least 1 and at most 10.
* `name` - (Required, String) Name of the NAT gateway.
* `vpc_id` - (Required, String, ForceNew) ID of the vpc.
* `bandwidth` - (Optional, Int) The maximum public network output bandwidth of NAT gateway (unit: Mbps). Valid values: `20`, `50`, `100`, `200`, `500`, `1000`, `2000`, `5000`. Default is 100.
* `max_concurrent` - (Optional, Int) The upper limit of concurrent connection of NAT gateway. Valid values: `1000000`, `3000000`, `10000000`. Default is `1000000`.
* `tags` - (Optional, Map) The available tags within this NAT gateway.
* `zone` - (Optional, String) The availability zone, such as `ap-guangzhou-3`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `created_time` - Create time of the NAT gateway.


## Import

NAT gateway can be imported using the id, e.g.

```
$ terraform import tencentcloud_nat_gateway.foo nat-1asg3t63
```

