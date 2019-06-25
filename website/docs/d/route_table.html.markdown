---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_route_table"
sidebar_current: "docs-tencentcloud-datasource-route-table"
description: |-
  Provides details about a specific Route Table.
---

# tencentcloud_route_table

`tencentcloud_route_table` provides details about a specific Route Table.

This resource can prove useful when a module accepts a Subnet id as an input variable and needs to, for example, add a route in the Route Table.

~> **NOTE:** It has been deprecated and replaced by tencentcloud_vpc_route_tables.

## Example Usage

The following example shows how one might accept a vpc id as a variable and use this data source to obtain the data necessary to create a route.

```hcl
variable "route_table_id" {}

data "tencentcloud_route_table" "selected" {
  route_table_id = "${var.route_table_id}"
}

resource "tencentcloud_route_entry" "rtb_entry_instance" {
  vpc_id         = "{data.tencentcloud_route_table.selected.vpc_id}"
  route_table_id = "${var.route_table_id}"
  cidr_block     = "10.4.8.0/24"
  next_type      = "instance"
  next_hub       = "10.16.1.7"
}
```

## Argument Reference

The arguments of this data source act as filters for querying the available
Route Table in the current region. The given filters must match exactly one
Route Table whose data will be exported as attributes.

* `route_table_id` - (Required) The Route Table ID.

## Attributes Reference

* `name` - The name for Route Table.
* `vpc_id` - The VPC ID.
* `routes` - routes are also exported with the following attributes, when there are relevants: Each route supports the following:
  * `cidr_block` - The RouteEntry's target network segment.
  * `next_type` - The `next_hub` type.
  * `next_hub` - The RouteEntry's next hub.
  * `description` - The RouteEntry's description.
* `subnet_num` - Number of associated subnets.
* `create_time` - Creation time of routing table, for example: 2018-01-22 17:50:21.
