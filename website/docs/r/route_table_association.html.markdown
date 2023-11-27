---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_route_table_association"
sidebar_current: "docs-tencentcloud-resource-route_table_association"
description: |-
  Provides a resource to create a vpc route_table
---

# tencentcloud_route_table_association

Provides a resource to create a vpc route_table

## Example Usage

```hcl
resource "tencentcloud_route_table_association" "route_table_association" {
  route_table_id = "rtb-5toos5sy"
  subnet_id      = "subnet-2y2omd4k"
}
```

## Argument Reference

The following arguments are supported:

* `route_table_id` - (Required, String) The route table instance ID, such as `rtb-azd4dt1c`.
* `subnet_id` - (Required, String, ForceNew) Subnet instance ID, such as `subnet-3x5lf5q0`. This can be queried using the DescribeSubnets API.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

vpc route_table can be imported using the id, e.g.

```
terraform import tencentcloud_route_table_association.route_table_association subnet_id
```

