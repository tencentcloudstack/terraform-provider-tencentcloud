---
subcategory: "TDMQ for RocketMQ(trocket)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tdmq_rocketmq_vip_instance"
sidebar_current: "docs-tencentcloud-resource-tdmq_rocketmq_vip_instance"
description: |-
  Provides a resource to create a tdmq rocketmq_vip_instance
---

# tencentcloud_tdmq_rocketmq_vip_instance

Provides a resource to create a tdmq rocketmq_vip_instance

~> **NOTE:** The instance cannot be downgraded, Include parameters `node_count`, `spec`, `storage_size`.
~> **NOTE:** If `spec` is `rocket-vip-basic-2`, configuration changes are not supported.

## Example Usage

```hcl
# query availability zones
data "tencentcloud_availability_zones" "zones" {}

# create vpc
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
}

# create subnet
resource "tencentcloud_subnet" "subnet" {
  availability_zone = data.tencentcloud_availability_zones.zones.zones.1.name
  name              = "subnet-example"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.0.0/16"
  is_multicast      = false
}

# create rocketmq vip instance
resource "tencentcloud_tdmq_rocketmq_vip_instance" "example" {
  name         = "tx-example"
  spec         = "rocket-vip-basic-2"
  node_count   = 2
  storage_size = 200
  zone_ids = [
    data.tencentcloud_availability_zones.zones.zones.0.id,
    data.tencentcloud_availability_zones.zones.zones.1.id
  ]

  vpc_info {
    vpc_id    = tencentcloud_vpc.vpc.id
    subnet_id = tencentcloud_subnet.subnet.id
  }

  time_span = 1
  ip_rules {
    ip_rule = "0.0.0.0/0"
    allow   = true
    remark  = "remark."
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Instance name.
* `node_count` - (Required, Int) Number of nodes, minimum 2, maximum 20.
* `spec` - (Required, String) Instance specification: Universal type, rocket-vip-basic-0, Basic type: `rocket-vip-basic-1`, Standard type: `rocket-vip-basic-2`, Advanced Type I: `rocket-vip-basic-3`, Advanced Type II: `rocket-vip-basic-4`.
* `storage_size` - (Required, Int) Single node storage space, in GB, minimum 200GB.
* `time_span` - (Required, Int) Purchase period, in months.
* `vpc_info` - (Required, List) VPC information.
* `zone_ids` - (Required, Set: [`String`]) The Zone ID list for node deployment, such as Guangzhou Zone 1, is 100001. For details, please refer to the official website of Tencent Cloud.
* `ip_rules` - (Optional, List) Public IP access control rules.

The `ip_rules` object supports the following:

* `allow` - (Required, Bool) Whether to allow or deny.
* `ip_rule` - (Required, String) IP address block information.
* `remark` - (Required, String) Remark.

The `vpc_info` object supports the following:

* `subnet_id` - (Required, String) Subnet ID.
* `vpc_id` - (Required, String) VPC ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



