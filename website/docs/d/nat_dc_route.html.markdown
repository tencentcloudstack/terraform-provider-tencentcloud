---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_nat_dc_route"
sidebar_current: "docs-tencentcloud-datasource-nat_dc_route"
description: |-
  Use this data source to query detailed information of vpc nat_dc_route
---

# tencentcloud_nat_dc_route

Use this data source to query detailed information of vpc nat_dc_route

## Example Usage

```hcl
data "tencentcloud_nat_dc_route" "nat_dc_route" {
  nat_gateway_id = "nat-gnxkey2e"
  vpc_id         = "vpc-pyyv5k3v"
}
```

## Argument Reference

The following arguments are supported:

* `nat_gateway_id` - (Required, String) Unique identifier of Nat Gateway.
* `vpc_id` - (Required, String) Unique identifier of Vpc.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `nat_direct_connect_gateway_route_set` - Data of route.
  * `create_time` - Create time of route.
  * `destination_cidr_block` - IPv4 CIDR of subnet.
  * `gateway_id` - Id of next-hop gateway.
  * `gateway_type` - Type of next-hop gateway, valid values: DIRECTCONNECT.
  * `update_time` - Update time of route.


