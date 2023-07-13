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
data "tencentcloud_vpc_instances" "foo" {
  name = "Default-VPC"
}
# Create EIP
resource "tencentcloud_eip" "eip_dev_dnat" {
  name = "terraform_nat_test"
}
resource "tencentcloud_eip" "new_eip" {
  name = "terraform_nat_test"
}

resource "tencentcloud_nat_gateway" "my_nat" {
  vpc_id         = data.tencentcloud_vpc_instances.foo.instance_list.0.vpc_id
  name           = "new_name"
  max_concurrent = 10000000
  bandwidth      = 1000
  zone           = "ap-guangzhou-3"

  assigned_eip_set = [
    tencentcloud_eip.eip_dev_dnat.public_ip,
    tencentcloud_eip.new_eip.public_ip,
  ]
  tags = {
    tf = "test"
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

