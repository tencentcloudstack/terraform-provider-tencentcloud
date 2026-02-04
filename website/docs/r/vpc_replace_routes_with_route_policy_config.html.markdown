---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc_replace_routes_with_route_policy_config"
sidebar_current: "docs-tencentcloud-resource-vpc_replace_routes_with_route_policy_config"
description: |-
  Provides a resource to create a VPC replace routes with route policy config
---

# tencentcloud_vpc_replace_routes_with_route_policy_config

Provides a resource to create a VPC replace routes with route policy config

## Example Usage

```hcl
resource "tencentcloud_vpc_replace_routes_with_route_policy_config" "example" {
  route_table_id = "rtb-olsbhnyc"
  routes {
    route_item_id      = "rti-araogi5t"
    force_match_policy = true
  }

  routes {
    route_item_id      = "rti-kiyt72op"
    force_match_policy = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `route_table_id` - (Required, String, ForceNew) Route Table Instance ID.
* `routes` - (Required, Set) Routing policy object. requires specifying the unique ID of routing policy (RouteItemId).

The `routes` object supports the following:

* `force_match_policy` - (Optional, Bool) Match the route reception policy tag.
* `route_item_id` - (Optional, String) Route unique policy ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



