---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc_instances"
sidebar_current: "docs-tencentcloud-datasource-vpc_instances"
description: |-
  Use this data source to query vpc instances' information.
---

# tencentcloud_vpc_instances

Use this data source to query vpc instances' information.

## Example Usage

```hcl
resource "tencentcloud_vpc" "foo" {
  name       = "guagua_vpc_instance_test"
  cidr_block = "10.0.0.0/16"
}

data "tencentcloud_vpc_instances" "id_instances" {
  vpc_id = tencentcloud_vpc.foo.id
}

data "tencentcloud_vpc_instances" "name_instances" {
  name = tencentcloud_vpc.foo.name
}
```

## Argument Reference

The following arguments are supported:

* `cidr_block` - (Optional, String) Filter VPC with this CIDR.
* `is_default` - (Optional, Bool) Filter default or no default VPC.
* `name` - (Optional, String) Name of the VPC to be queried.
* `result_output_file` - (Optional, String) Used to save results.
* `tag_key` - (Optional, String) Filter if VPC has this tag.
* `tags` - (Optional, Map) Tags of the VPC to be queried.
* `vpc_id` - (Optional, String) ID of the VPC to be queried.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instance_list` - The information list of the VPC.
  * `cidr_block` - A network address block of a VPC CIDR.
  * `create_time` - Creation time of VPC.
  * `dns_servers` - A list of DNS servers which can be used within the VPC.
  * `is_default` - Indicates whether it is the default VPC for this region.
  * `is_multicast` - Indicates whether VPC multicast is enabled.
  * `name` - Name of the VPC.
  * `subnet_ids` - A ID list of subnets within this VPC.
  * `tags` - Tags of the VPC.
  * `vpc_id` - ID of the VPC.


