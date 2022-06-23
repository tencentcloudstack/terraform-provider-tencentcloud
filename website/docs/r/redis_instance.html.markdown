---
subcategory: "Redis"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_redis_instance"
sidebar_current: "docs-tencentcloud-resource-redis_instance"
description: |-
  Provides a resource to create a Redis instance and set its attributes.
---

# tencentcloud_redis_instance

Provides a resource to create a Redis instance and set its attributes.

## Example Usage

```hcl
data "tencentcloud_redis_zone_config" "zone" {
}

resource "tencentcloud_redis_instance" "redis_instance_test_2" {
  availability_zone  = data.tencentcloud_redis_zone_config.zone.list[0].zone
  type_id            = data.tencentcloud_redis_zone_config.zone.list[0].type_id
  password           = "test12345789"
  mem_size           = 8192
  redis_shard_num    = data.tencentcloud_redis_zone_config.zone.list[0].redis_shard_nums[0]
  redis_replicas_num = data.tencentcloud_redis_zone_config.zone.list[0].redis_replicas_nums[0]
  name               = "terrform_test"
  port               = 6379
}
```

Using multi replica zone set

```hcl
data "tencentcloud_availability_zones" "az" {

}

variable "redis_replicas_num" {
  default = 3
}

resource "tencentcloud_redis_instance" "red1" {
  availability_zone  = data.tencentcloud_availability_zones.az.zones[0].name
  charge_type        = "POSTPAID"
  mem_size           = 1024
  name               = "test-redis"
  port               = 6379
  project_id         = 0
  redis_replicas_num = var.redis_replicas_num
  redis_shard_num    = 1
  security_groups = [
    "sg-d765yoec",
  ]
  subnet_id = "subnet-ie01x91v"
  type_id   = 6
  vpc_id    = "vpc-k4lrsafc"
  password  = "a12121312334"

  replica_zone_ids = [
    for i in range(var.redis_replicas_num)
  : data.tencentcloud_availability_zones.az.zones[i % length(data.tencentcloud_availability_zones.az.zones)].id]
}
```

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Required, ForceNew) The available zone ID of an instance to be created, please refer to `tencentcloud_redis_zone_config.list`.
* `mem_size` - (Required) The memory volume of an available instance(in MB), please refer to `tencentcloud_redis_zone_config.list[zone].shard_memories`. When redis is standard type, it represents total memory size of the instance; when Redis is cluster type, it represents memory size of per sharding.
* `auto_renew_flag` - (Optional, ForceNew) Auto-renew flag. 0 - default state (manual renewal); 1 - automatic renewal; 2 - explicit no automatic renewal.
* `charge_type` - (Optional, ForceNew) The charge type of instance. Valid values: `PREPAID` and `POSTPAID`. Default value is `POSTPAID`. Note: TencentCloud International only supports `POSTPAID`. Caution that update operation on this field will delete old instances and create new with new charge type.
* `force_delete` - (Optional) Indicate whether to delete Redis instance directly or not. Default is false. If set true, the instance will be deleted instead of staying recycle bin. Note: only works for `PREPAID` instance.
* `name` - (Optional) Instance name.
* `no_auth` - (Optional, ForceNew) Indicates whether the redis instance support no-auth access. NOTE: Only available in private cloud environment.
* `password` - (Optional) Password for a Redis user, which should be 8 to 16 characters. NOTE: Only `no_auth=true` specified can make password empty.
* `port` - (Optional, ForceNew) The port used to access a redis instance. The default value is 6379. And this value can't be changed after creation, or the Redis instance will be recreated.
* `prepaid_period` - (Optional) The tenancy (time unit is month) of the prepaid instance, NOTE: it only works when charge_type is set to `PREPAID`. Valid values are `1`, `2`, `3`, `4`, `5`, `6`, `7`, `8`, `9`, `10`, `11`, `12`, `24`, `36`.
* `project_id` - (Optional) Specifies which project the instance should belong to.
* `redis_replicas_num` - (Optional) The number of instance copies. This is not required for standalone and master slave versions.
* `redis_shard_num` - (Optional) The number of instance shard, default is 1. This is not required for standalone and master slave versions.
* `replica_zone_ids` - (Optional) ID of replica nodes available zone. This is not required for standalone and master slave versions.
* `replicas_read_only` - (Optional, ForceNew) Whether copy read-only is supported, Redis 2.8 Standard Edition and CKV Standard Edition do not support replica read-only, turn on replica read-only, the instance will automatically read and write separate, write requests are routed to the primary node, read requests are routed to the replica node, if you need to open replica read-only, the recommended number of replicas >=2.
* `security_groups` - (Optional, ForceNew) ID of security group. If both vpc_id and subnet_id are not set, this argument should not be set either.
* `subnet_id` - (Optional, ForceNew) Specifies which subnet the instance should belong to.
* `tags` - (Optional) Instance tags.
* `type_id` - (Optional, ForceNew) Instance type. Available values reference data source `tencentcloud_redis_zone_config` or [document](https://intl.cloud.tencent.com/document/product/239/32069).
* `type` - (Optional, ForceNew, **Deprecated**) It has been deprecated from version 1.33.1. Please use 'type_id' instead. Instance type. Available values: `cluster_ckv`,`cluster_redis5.0`,`cluster_redis`,`master_slave_ckv`,`master_slave_redis4.0`,`master_slave_redis5.0`,`master_slave_redis`,`standalone_redis`, specific region support specific types, need to refer data `tencentcloud_redis_zone_config`.
* `vpc_id` - (Optional, ForceNew) ID of the vpc with which the instance is to be associated.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - The time when the instance was created.
* `ip` - IP address of an instance.
* `status` - Current status of an instance, maybe: init, processing, online, isolate and todelete.


## Import

Redis instance can be imported, e.g.

```
$ terraform import tencentcloud_redis_instance.redislab redis-id
```

