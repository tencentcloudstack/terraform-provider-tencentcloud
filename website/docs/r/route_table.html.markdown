---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_route_table"
sidebar_current: "docs-tencentcloud-resource-route_table"
description: |-
  Provides a resource to create a VPC routing table.
---

# tencentcloud_route_table

Provides a resource to create a VPC routing table.

## Example Usage

```hcl
resource "tencentcloud_vpc" "foo" {
  name       = "ci-temp-test"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_route_table" "foo" {
  vpc_id = tencentcloud_vpc.foo.id
  name   = "ci-temp-test-rt"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) The name of routing table.
* `vpc_id` - (Required, String, ForceNew) ID of VPC to which the route table should be associated.
* `tags` - (Optional, Map) The tags of routing table.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Creation time of the routing table.
* `is_default` - Indicates whether it is the default routing table.
* `route_entry_ids` - ID list of the routing entries.
* `subnet_ids` - ID list of the subnets associated with this route table.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateRouteTable](https://cloud.tencent.com/document/api/215/15768) - This API is used to create a route table.
* Create: [ModifyResourceTags](https://cloud.tencent.com/document/api/651/35322) - This API is used to modify all tags associated with a resource.
* Read: [DescribeResourceTagsByResourceIds](https://cloud.tencent.com/document/api/651/35319) - This API is used to query the tag key-value pairs associated with an existing resource.
* Read: [DescribeRouteTables](https://cloud.tencent.com/document/api/215/15763) - This API is used to query route tables.
* Update: [ModifyResourceTags](https://cloud.tencent.com/document/api/651/35322) - This API is used to modify all tags associated with a resource.
* Update: [ModifyRouteTableAttribute](https://cloud.tencent.com/document/api/215/15766) - This API (ModifyRouteTableAttribute) is used to modify the attributes of a route table.
* Delete: [DeleteRouteTable](https://cloud.tencent.com/document/api/215/15771) - This API is used to delete a route table.

## Import

Vpc routetable instance can be imported, e.g.

```
$ terraform import tencentcloud_route_table.test route_table_id
```

