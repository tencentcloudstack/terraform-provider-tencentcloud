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

~> **NOTE:** It has been deprecated and replaced by `tencentcloud_route_table_entry`.

## Example Usage

```hcl
resource "tencentcloud_vpc" "vpc" {
  name       = "tf-example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_route_table" "example" {
  name   = "tf-example"
  vpc_id = tencentcloud_vpc.vpc.id
}

resource "tencentcloud_route_entry" "example1" {
  vpc_id         = tencentcloud_vpc.vpc.id
  route_table_id = tencentcloud_route_table.example.id
  cidr_block     = "192.168.0.0/24"
  next_type      = "eip"
  next_hub       = "0"
}

resource "tencentcloud_route_entry" "example2" {
  vpc_id         = tencentcloud_vpc.vpc.id
  route_table_id = tencentcloud_route_table.example.id
  cidr_block     = "10.4.8.0/24"
  next_type      = "instance"
  next_hub       = "10.16.1.7"
}
```

## Argument Reference

The following arguments are supported:

* `cidr_block` - (Required, String, ForceNew) The RouteEntry's target network segment.
* `next_hub` - (Required, String, ForceNew) The route entry's next hub. CVM instance ID or VPC router interface ID.
* `next_type` - (Required, String, ForceNew) The next hop type. Valid values: `public_gateway`,`vpn_gateway`,`sslvpn_gateway`,`dc_gateway`,`peering_connection`,`nat_gateway`,`havip`,`local_gateway`, `intranat`, `user_ccn`, `gwlb_endpoint` and `instance`. `instance` points to CVM Instance.
* `route_table_id` - (Required, String, ForceNew) The ID of the route table.
* `vpc_id` - (Required, String, ForceNew) The VPC ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



