---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc_subnets"
sidebar_current: "docs-tencentcloud-datasource-vpc_subnets"
description: |-
  Use this data source to query vpc subnets information.
---

# tencentcloud_vpc_subnets

Use this data source to query vpc subnets information.

## Example Usage

### Create subnet resource

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-3"
}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = var.availability_zone
  name              = "subnet1"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.20.0/28"
  is_multicast      = false

  tags = {
    "test" = "test"
  }
}

resource "tencentcloud_subnet" "subnetCDC" {
  vpc_id            = tencentcloud_vpc.vpc.id
  name              = "subnet2"
  cidr_block        = "10.0.0.0/16"
  cdc_id            = "cluster-lchwgxhs"
  availability_zone = data.tencentcloud_availability_zones.zones.zones.0.name
  is_multicast      = false
}
```

### Query all subnets

```hcl
data "tencentcloud_vpc_subnets" "subnets" {}
```

### Query subnets by filter

```hcl
data "tencentcloud_vpc_subnets" "subnets" {
  vpc_id = tencentcloud_vpc.vpc.id
}

data "tencentcloud_vpc_subnets" "subnets" {
  subnet_id = tencentcloud_subnet.subnet.id
}

data "tencentcloud_vpc_subnets" "subnets" {
  name = tencentcloud_subnet.subnet.name
}

data "tencentcloud_vpc_subnets" "subnets" {
  tags = tencentcloud_subnet.subnet.tags
}

data "tencentcloud_vpc_subnets" "subnets" {
  cdc_id = tencentcloud_subnet.subnetCDC.cdc_id
}
```

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Optional, String) Zone of the subnet to be queried.
* `cdc_id` - (Optional, String) ID of CDC instance.
* `cidr_block` - (Optional, String) Filter subnet with this CIDR.
* `is_default` - (Optional, Bool) Filter default or no default subnets.
* `is_remote_vpc_snat` - (Optional, Bool) Filter the VPC SNAT address pool subnet.
* `name` - (Optional, String) Name of the subnet to be queried.
* `result_output_file` - (Optional, String) Used to save results.
* `subnet_id` - (Optional, String) ID of the subnet to be queried.
* `tag_key` - (Optional, String) Filter if subnet has this tag.
* `tags` - (Optional, Map) Tags of the subnet to be queried.
* `vpc_id` - (Optional, String) ID of the VPC to be queried.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instance_list` - List of subnets.
  * `availability_zone` - The availability zone of the subnet.
  * `available_ip_count` - The number of available IPs.
  * `cdc_id` - ID of CDC instance.
  * `cidr_block` - A network address block of the subnet.
  * `create_time` - Creation time of the subnet resource.
  * `is_default` - Indicates whether it is the default subnet of the VPC for this region.
  * `is_multicast` - Indicates whether multicast is enabled.
  * `name` - Name of the subnet.
  * `route_table_id` - ID of the routing table.
  * `subnet_id` - ID of the subnet.
  * `tags` - Tags of the subnet resource.
  * `vpc_id` - ID of the VPC.


