---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_route_table"
sidebar_current: "docs-tencentcloud-resource-vpc-route-table"
description: |-
  Provides a resource to create a VPC routing table.
---

# tencentcloud_route_table

Provides a resource to create a VPC routing table.

## Example Usage

Basic usage:

```hcl
resource "tencentcloud_route_table" "r" {
  name   = "my test route table"
  vpc_id = "${tencent_vpc.main.id}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name for the Route Table.
* `vpc_id` - (Required, Forces new resource) The VPC ID.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Route Table.
* `name` - The name for Route Table.
* `vpc_id` - The VPC ID.
