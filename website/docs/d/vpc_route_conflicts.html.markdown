---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc_route_conflicts"
sidebar_current: "docs-tencentcloud-datasource-vpc_route_conflicts"
description: |-
  Use this data source to query detailed information of vpc route_conflicts
---

# tencentcloud_vpc_route_conflicts

Use this data source to query detailed information of vpc route_conflicts

## Example Usage

```hcl
data "tencentcloud_vpc_route_conflicts" "route_conflicts" {
  route_table_id          = "rtb-6xypllqe"
  destination_cidr_blocks = ["172.18.111.0/24"]
}
```

## Argument Reference

The following arguments are supported:

* `destination_cidr_blocks` - (Required, Set: [`String`]) List of conflicting destinations to check for.
* `route_table_id` - (Required, String) Routing table instance ID, for example:rtb-azd4dt1c.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `route_conflict_set` - route conflict list.
  * `conflict_set` - route conflict list.
    * `created_time` - create time.
    * `destination_cidr_block` - Destination Cidr Block, like 112.20.51.0/24.
    * `destination_ipv6_cidr_block` - Destination of Ipv6 Cidr Block.
    * `enabled` - if enabled.
    * `gateway_id` - next hop id.
    * `gateway_type` - next gateway type.
    * `published_to_vbc` - if published To ccn.
    * `route_description` - route description.
    * `route_id` - route id.
    * `route_item_id` - unique policy id.
    * `route_table_id` - route table id.
    * `route_type` - routr type.
  * `destination_cidr_block` - destination cidr block.
  * `route_table_id` - route table id.


