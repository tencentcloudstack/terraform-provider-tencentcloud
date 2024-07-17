---
subcategory: "Cloud Connect Network(CCN)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ccn_routes"
sidebar_current: "docs-tencentcloud-datasource-ccn_routes"
description: |-
  Use this data source to query detailed information of CCN routes.
---

# tencentcloud_ccn_routes

Use this data source to query detailed information of CCN routes.

## Example Usage

### Query CCN instance all routes

```hcl
data "tencentcloud_ccn_routes" "routes" {
  ccn_id = "ccn-gr7nynbd"
}
```

### Query CCN instance routes by filter

```hcl
data "tencentcloud_ccn_routes" "routes" {
  ccn_id = "ccn-gr7nynbd"
  filters {
    name   = "route-table-id"
    values = ["ccnrtb-jpf7bzn3"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `ccn_id` - (Required, String) ID of the CCN to be queried.
* `filters` - (Optional, List) Filter conditions.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `name` - (Required, String) Field to be filtered. Support `route-id`, `cidr-block`, `instance-type`, `instance-region`, `instance-id`, `route-table-id`.
* `values` - (Required, Set) Filter value of the field.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `route_list` - CCN route list.
  * `destination_cidr_block` - Destination.
  * `enabled` - Is routing enabled.
  * `extra_state` - Extension status of routing.
  * `instance_extra_name` - Next hop extension name (associated instance extension name).
  * `instance_id` - Next jump (associated instance ID).
  * `instance_name` - Next jump (associated instance name).
  * `instance_region` - Next jump (associated instance region).
  * `instance_type` - Next hop type (associated instance type), all types: VPC, DIRECTCONNECT.
  * `instance_uin` - The UIN (root account) to which the associated instance belongs.
  * `is_bgp` - Is it dynamic routing.
  * `route_id` - route ID.
  * `route_priority` - Routing priority.
  * `update_time` - update time.


