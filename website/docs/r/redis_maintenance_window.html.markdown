---
subcategory: "TencentDB for Redis(crs)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_redis_maintenance_window"
sidebar_current: "docs-tencentcloud-resource-redis_maintenance_window"
description: |-
  Provides a resource to create a redis maintenance_window
---

# tencentcloud_redis_maintenance_window

Provides a resource to create a redis maintenance_window

## Example Usage

```hcl
data "tencentcloud_redis_zone_config" "zone" {
  type_id = 7
}

resource "tencentcloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
  name       = "tf_redis_vpc"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = data.tencentcloud_redis_zone_config.zone.list[0].zone
  name              = "tf_redis_subnet"
  cidr_block        = "10.0.1.0/24"
}

resource "tencentcloud_redis_instance" "foo" {
  availability_zone  = data.tencentcloud_redis_zone_config.zone.list[0].zone
  type_id            = data.tencentcloud_redis_zone_config.zone.list[0].type_id
  password           = "test12345789"
  mem_size           = 8192
  redis_shard_num    = data.tencentcloud_redis_zone_config.zone.list[0].redis_shard_nums[0]
  redis_replicas_num = data.tencentcloud_redis_zone_config.zone.list[0].redis_replicas_nums[0]
  name               = "terrform_test"
  port               = 6379
  vpc_id             = tencentcloud_vpc.vpc.id
  subnet_id          = tencentcloud_subnet.subnet.id
}

resource "tencentcloud_redis_maintenance_window" "foo" {
  instance_id = tencentcloud_redis_instance.foo.id
  start_time  = "17:00"
  end_time    = "19:00"
}
```

## Argument Reference

The following arguments are supported:

* `end_time` - (Required, String) The end time of the maintenance window, e.g. 19:00.
* `instance_id` - (Required, String) The ID of instance.
* `start_time` - (Required, String) Maintenance window start time, e.g. 17:00.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

redis maintenance_window can be imported using the id, e.g.

```
terraform import tencentcloud_redis_maintenance_window.foo instance_id
```

