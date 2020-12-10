---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_route_entry"
sidebar_current: "docs-tencentcloud-resource-route_entry"
description: |-
  Provides a resource to create a routing entry in a VPC routing table.
---

# tencentcloud_route_entry

Provides a resource to create a routing entry in a VPC routing table.

~> **NOTE:** It has been deprecated and replaced by tencentcloud_route_table_entry.

## Example Usage

```hcl
resource "tencentcloud_vpc" "main" {
  name       = "Used to test the routing entry"
  cidr_block = "10.4.0.0/16"
}

resource "tencentcloud_route_table" "r" {
  name   = "Used to test the routing entry"
  vpc_id = tencentcloud_vpc.main.id
}

resource "tencentcloud_route_entry" "rtb_entry_instance" {
  vpc_id         = tencentcloud_route_table.main.vpc_id
  route_table_id = tencentcloud_route_table.r.id
  cidr_block     = "10.4.8.0/24"
  next_type      = "instance"
  next_hub       = "10.16.1.7"
}

resource "tencentcloud_route_entry" "rtb_entry_instance" {
  vpc_id         = tencentcloud_route_table.main.vpc_id
  route_table_id = tencentcloud_route_table.r.id
  cidr_block     = "10.4.5.0/24"
  next_type      = "vpn_gateway"
  next_hub       = "vpngw-db52irtl"
}
```

## Argument Reference

The following arguments are supported:

* `cidr_block` - (Required, ForceNew) The RouteEntry's target network segment.
* `next_hub` - (Required, ForceNew) The route entry's next hub. CVM instance ID or VPC router interface ID.
* `next_type` - (Required, ForceNew) The next hop type. Valid values: `public_gateway`,`vpn_gateway`,`sslvpn_gateway`,`dc_gateway`,`peering_connection`,`nat_gateway` and `instance`. `instance` points to CVM Instance.
* `route_table_id` - (Required, ForceNew) The ID of the route table.
* `vpc_id` - (Required, ForceNew) The VPC ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



