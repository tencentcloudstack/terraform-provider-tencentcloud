---
subcategory: "VPC"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc_route_tables"
sidebar_current: "docs-tencentcloud-datasource-vpc_route_tables"
description: |-
  Use this data source to query vpc route tables information.
---

# tencentcloud_vpc_route_tables

Use this data source to query vpc route tables information.

## Example Usage

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-3"
}

resource "tencentcloud_vpc" "foo" {
  name       = "guagua-ci-temp-test"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_route_table" "route_table" {
  vpc_id = tencentcloud_vpc.foo.id
  name   = "ci-temp-test-rt"

  tags = {
    "test" = "test"
  }
}

data "tencentcloud_vpc_route_tables" "id_instances" {
  route_table_id = tencentcloud_route_table.route_table.id
}

data "tencentcloud_vpc_route_tables" "name_instances" {
  name = tencentcloud_route_table.route_table.name
}

data "tencentcloud_vpc_route_tables" "tags_instances" {
  tags = tencentcloud_route_table.route_table.tags
}
```

## Argument Reference

The following arguments are supported:

* `association_main` - (Optional) Filter the main routing table.
* `name` - (Optional) Name of the routing table to be queried.
* `result_output_file` - (Optional) Used to save results.
* `route_table_id` - (Optional) ID of the routing table to be queried.
* `tag_key` - (Optional) Filter if routing table has this tag.
* `tags` - (Optional) Tags of the routing table to be queried.
* `vpc_id` - (Optional) ID of the VPC to be queried.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instance_list` - The information list of the VPC route table.
  * `create_time` - Creation time of the routing table.
  * `is_default` - Indicates whether it is the default routing table.
  * `name` - Name of the routing table.
  * `route_entry_infos` - Detailed information of each entry of the route table.
    * `description` - Description information user defined for a route table rule.
    * `destination_cidr_block` - The destination address block.
    * `next_hub` - ID of next-hop gateway. Note: when 'next_type' is EIP, GatewayId will fix the value '0'.
    * `next_type` - Type of next-hop, and available values include CVM, VPN, DIRECTCONNECT, PEERCONNECTION, SSLVPN, NAT, NORMAL_CVM, EIP and CCN.
    * `route_entry_id` - ID of a route table entry.
  * `route_table_id` - ID of the routing table.
  * `subnet_ids` - List of subnet IDs bound to the route table.
  * `tags` - Tags of the routing table.
  * `vpc_id` - ID of the VPC.


