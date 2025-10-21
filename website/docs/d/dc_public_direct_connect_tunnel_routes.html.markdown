---
subcategory: "Direct Connect(DC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dc_public_direct_connect_tunnel_routes"
sidebar_current: "docs-tencentcloud-datasource-dc_public_direct_connect_tunnel_routes"
description: |-
  Use this data source to query detailed information of dc public_direct_connect_tunnel_routes
---

# tencentcloud_dc_public_direct_connect_tunnel_routes

Use this data source to query detailed information of dc public_direct_connect_tunnel_routes

## Example Usage

```hcl
data "tencentcloud_dc_public_direct_connect_tunnel_routes" "public_direct_connect_tunnel_routes" {
  direct_connect_tunnel_id = "dcx-4z49tnws"
}
```

## Argument Reference

The following arguments are supported:

* `direct_connect_tunnel_id` - (Required, String) direct connect tunnel id.
* `filters` - (Optional, List) filter condition: route-type: route type, value: BGP/STATIC route-subnet: route cidr, value such as: 192.68.1.0/24.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `name` - (Required, String) Fields to be filtered.
* `values` - (Required, Set) filter value of the field.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `routes` - Internet tunnel route list.
  * `as_path` - ASPath info.
  * `destination_cidr_block` - Network CIDR.
  * `next_hop` - Route next hop ip.
  * `route_id` - direct connect tunnel route id.
  * `route_type` - Route type: BGP/STATIC route.
  * `status` - ENABLE: routing is enabled, DISABLE: routing is disabled.


