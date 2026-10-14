---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc_notify_routes"
sidebar_current: "docs-tencentcloud-resource-vpc_notify_routes"
description: |-
  Provides a resource to create a vpc notify_routes
---

# tencentcloud_vpc_notify_routes

Provides a resource to create a vpc notify_routes

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

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `published_to_vbc` - If published to vbc.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [NotifyRoutes](https://cloud.tencent.com/document/api/215/52618) - This API is used to publish a route to CCN. This can also be done by clicking "Publish to CCN" in the operation column on the page of route table list.
* Read: [DescribeRouteTables](https://cloud.tencent.com/document/api/215/15763) - This API is used to query route tables.
* Delete: [WithdrawNotifyRoutes](https://cloud.tencent.com/document/api/215/52617) - This API is used to withdraw a route from CCN.

## Import

vpc notify_routes can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_notify_routes.notify_routes route_table_id#route_item_id
```

