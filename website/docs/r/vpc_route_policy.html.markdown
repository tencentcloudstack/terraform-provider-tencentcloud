---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc_route_policy"
sidebar_current: "docs-tencentcloud-resource-vpc_route_policy"
description: |-
  Provides a resource to create a VPC route policy
---

# tencentcloud_vpc_route_policy

Provides a resource to create a VPC route policy

## Example Usage

```hcl
resource "tencentcloud_vpc_route_policy" "example" {
  route_policy_name        = "tf-example"
  route_policy_description = "remark."
}
```

## Argument Reference

The following arguments are supported:

* `route_policy_description` - (Required, String) Routing policy description.
* `route_policy_name` - (Required, String) Specifies the routing strategy name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `route_policy_id` - Route policy ID.


## Import

VPC route policy can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_route_policy.example rrp-lpv8rjp8
```

