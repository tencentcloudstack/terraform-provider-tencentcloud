---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc_notify_routes"
sidebar_current: "docs-tencentcloud-resource-vpc_notify_routes"
description: |-
  Provides a resource to create a VPC notify routes
---

# tencentcloud_vpc_notify_routes

Provides a resource to create a VPC notify routes

## Example Usage

```hcl
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_route_table" "route_table" {
  vpc_id = tencentcloud_vpc.vpc.id
  name   = "tf-example"
}

resource "tencentcloud_vpc_notify_routes" "example" {
  route_table_id = tencentcloud_route_table.route_table.id
  route_item_ids = ["rti-i8bap903"]
}
```

## Argument Reference

The following arguments are supported:

* `route_item_ids` - (Required, Set: [`String`], ForceNew) The unique ID of the routing policy.
* `route_table_id` - (Required, String, ForceNew) The unique ID of the routing table.
* `expected_published_status` - (Optional, Bool, ForceNew) Set the desired publication status: true: published; false: not published.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `published_to_vbc` - Whether to publish policies to vbc.


## Import

VPC notify routes can be imported using the routeTableId#routeItemId, e.g.

```
terraform import tencentcloud_vpc_notify_routes.example route_table_id#route_item_id
```

