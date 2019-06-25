---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc_subnets"
sidebar_current: "docs-tencentcloud-datasource-vpc_subnets"
description: |-
  Use this data source to query vpc subnets information.
---

# tencentcloud_vpc_subnets

Use this data source to query vpc subnets information.

## Example Usage

```hcl
variable "availability_zone" {
	default = "ap-guangzhou-3"
}

resource "tencentcloud_vpc" "foo" {
    name="guagua_vpc_instance_test"
    cidr_block="10.0.0.0/16"
}
resource "tencentcloud_subnet" "subnet" {
	availability_zone="${var.availability_zone}"
	name="guagua_vpc_subnet_test"
	vpc_id="${tencentcloud_vpc.foo.id}"
	cidr_block="10.0.20.0/28"
	is_multicast=false
}
data "tencentcloud_vpc_subnets" "id_instances" {
	subnet_id="${tencentcloud_subnet.subnet.id}"
}
data "tencentcloud_vpc_subnets" "name_instances" {
	name="${tencentcloud_subnet.subnet.name}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional, ForceNew) Name of the subnet to be queried.
* `result_output_file` - (Optional, ForceNew) Used to save results.
* `subnet_id` - (Optional, ForceNew) ID of the subnet to be queried.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instance_list` - List of subnets.
  * `availability_zone` - The availability zone of the subnet.
  * `available_ip_count` - The number of available IPs.
  * `cidr_block` - A network address block of the subnet.
  * `create_time` - Creation time of the subnet resource.
  * `is_default` - Indicates whether it is the default subnet of the VPC for this region.
  * `is_multicast` - Indicates whether multicast is enabled.
  * `name` - Name of the subnet.
  * `route_table_id` - ID of the routing table.
  * `subnet_id` - ID of the subnet.
  * `vpc_id` - ID of the VPC.


