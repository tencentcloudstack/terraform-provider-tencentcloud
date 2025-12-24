---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc_route_policy_association"
sidebar_current: "docs-tencentcloud-resource-vpc_route_policy_association"
description: |-
  Provides a resource to create a VPC route policy association
---

# tencentcloud_vpc_route_policy_association

Provides a resource to create a VPC route policy association

## Example Usage

```hcl
resource "tencentcloud_vpc_route_policy_association" "example" {
  route_policy_id = "rrp-7dnu4yoi"
  route_table_id  = "rtb-389phpuq"
  priority        = 10
}
```

## Argument Reference

The following arguments are supported:

* `priority` - (Required, Int) Priority.
* `route_policy_id` - (Required, String, ForceNew) Specifies the unique ID of the route reception policy.
* `route_table_id` - (Required, String, ForceNew) Unique route table ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

VPC route policy association can be imported using the routePolicyId#routeTableId, e.g.

```
terraform import tencentcloud_vpc_route_policy_association.example rrp-7dnu4yoi#rtb-389phpuq
```

