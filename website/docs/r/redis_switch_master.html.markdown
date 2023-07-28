---
subcategory: "TencentDB for Redis(crs)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_redis_switch_master"
sidebar_current: "docs-tencentcloud-resource-redis_switch_master"
description: |-
  Provides a resource to create a redis switch_master
---

# tencentcloud_redis_switch_master

Provides a resource to create a redis switch_master

## Example Usage

```hcl
data "tencentcloud_redis_zone_config" "zone" {
  type_id = 7
  region  = "ap-guangzhou"
}

variable "replica_zone_ids" {
  default = [100004, 100006]
}

resource "tencentcloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
  name       = "tf_redis_vpc"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = data.tencentcloud_redis_zone_config.zone.list[2].zone
  name              = "tf_redis_subnet"
  cidr_block        = "10.0.1.0/24"
}

resource "tencentcloud_security_group" "foo" {
  name = "tf-redis-sg"
}

resource "tencentcloud_security_group_lite_rule" "foo" {
  security_group_id = tencentcloud_security_group.foo.id

  ingress = [
    "ACCEPT#192.168.1.0/24#80#TCP",
    "DROP#8.8.8.8#80,90#UDP",
    "DROP#0.0.0.0/0#80-90#TCP",
  ]

  egress = [
    "ACCEPT#192.168.0.0/16#ALL#TCP",
    "ACCEPT#10.0.0.0/8#ALL#ICMP",
    "DROP#0.0.0.0/0#ALL#ALL",
  ]
}

resource "tencentcloud_redis_instance" "foo" {
  availability_zone  = data.tencentcloud_redis_zone_config.zone.list[2].zone
  type_id            = data.tencentcloud_redis_zone_config.zone.list[2].type_id
  password           = "test12345789"
  mem_size           = 8192
  redis_shard_num    = data.tencentcloud_redis_zone_config.zone.list[2].redis_shard_nums[0]
  redis_replicas_num = 2
  replica_zone_ids   = var.replica_zone_ids
  name               = "terrform_test"
  port               = 6379
  vpc_id             = tencentcloud_vpc.vpc.id
  subnet_id          = tencentcloud_subnet.subnet.id
  security_groups    = [tencentcloud_security_group.foo.id]
}

data "tencentcloud_redis_instance_zone_info" "foo" {
  instance_id = tencentcloud_redis_instance.foo.id
}

resource "tencentcloud_redis_switch_master" "switch_master" {
  instance_id = tencentcloud_redis_instance.foo.id
  group_id    = data.tencentcloud_redis_instance_zone_info.foo.replica_groups[1].group_id
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) The ID of instance.
* `group_id` - (Optional, Int) Replication group ID, required for multi-AZ instances.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



