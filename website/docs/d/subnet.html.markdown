---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_subnet"
sidebar_current: "docs-tencentcloud-datasource-subnet-x"
description: |-
  Provides details about a specific VPC subnet.
---

# tencentcloud_subnet

`tencentcloud_subnet` provides details about a specific VPC subnet.

This resource can prove useful when a module accepts a subnet id as an input variable and needs to, for example, determine the id of the VPC that the subnet belongs to.

~> **NOTE:** It has been deprecated and replaced by tencentcloud_vpc_subnets.

## Example Usage

The following example shows how one might accept a subnet id as a variable and use this data source to obtain the data necessary to create a security group that allows connections from hosts in that subnet.

```hcl
variable "subnet_id" {}
variable "vpc_id" {}

data "tencentcloud_subnet" "selected" {
  vpc_id    = "${var.vpc_id}"
  subnet_id = "${var.subnet_id}"
}

resource "tencentcloud_security_group" "default" {
  name        = "test subnet data"
  description = "test subnet data description"
}

resource "tencentcloud_security_group_rule" "subnet" {
  security_group_id = "${tencentcloud_security_group.default.id}"
  type              = "ingress"
  cidr_ip           = "${data.tencentcloud_subnet.selected.cidr_block}"
  ip_protocol       = "tcp"
  port_range        = "80,8080"
  policy            = "accept"
}
```

## Argument Reference

The arguments of this data source act as filters for querying the available subnets in the current region. The given filters must match exactly one subnet whose data will be exported as attributes.

* `vpc_id` - (Required) The VPC ID.
* `subnet_id` - (Required) The ID of the Subnet.

## Attributes Reference

The following attributes are exported:

* `name` - The name for the Subnet.
* `cidr_block` - The CIDR block of the Subnet.
* `availability_zone`- The AZ for the subnet.
* `route_table_id` - The Route Table ID.
