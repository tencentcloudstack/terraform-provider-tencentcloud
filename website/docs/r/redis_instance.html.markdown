---
subcategory: "TencentDB for Redis(crs)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_redis_instance"
sidebar_current: "docs-tencentcloud-resource-redis_instance"
description: |-
  Provides a resource to create a Redis instance and set its attributes.
---

# tencentcloud_redis_instance

Provides a resource to create a Redis instance and set its attributes.

~> **NOTE:** The argument vpc_id and subnet_id is now required because Basic Network Instance is no longer supported.

~> **NOTE:** Both adding and removing replications in one change is supported but not recommend.

## Example Usage

### Create a base version of redis

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
```

### Using multi replica zone set

```hcl
variable "redis_replicas_num" {
  default = 3
}

variable "redis_type_id" {
  default = 7
}

data "tencentcloud_availability_zones_by_product" "az" {
  product = "redis"
}

resource "tencentcloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
  name       = "tf_redis_vpc"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = data.tencentcloud_availability_zones_by_product.az.zones[0].name
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

resource "tencentcloud_redis_instance" "red1" {
  availability_zone  = data.tencentcloud_availability_zones_by_product.az.zones[0].name
  type_id            = var.redis_type_id
  charge_type        = "POSTPAID"
  mem_size           = 1024
  name               = "test-redis"
  port               = 6379
  project_id         = 0
  vpc_id             = tencentcloud_vpc.vpc.id
  subnet_id          = tencentcloud_subnet.subnet.id
  password           = "a12121312334"
  security_groups    = [tencentcloud_security_group.foo.id]
  redis_replicas_num = var.redis_replicas_num
  redis_shard_num    = 1
  replica_zone_ids = [
    for i in range(var.redis_replicas_num)
    : data.tencentcloud_availability_zones_by_product.az.zones[i % length(data.tencentcloud_availability_zones_by_product.az.zones)].id
  ]
}
```

### Buy a month of prepaid instances

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
  availability_zone = data.tencentcloud_redis_zone_config.zone.list[1].zone
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
  security_groups    = [tencentcloud_security_group.foo.id]
  charge_type        = "PREPAID"
  prepaid_period     = 1
}
```

### Create a multi-AZ instance

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
```

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Required, String, ForceNew) The available zone ID of an instance to be created, please refer to `tencentcloud_redis_zone_config.list`.
* `mem_size` - (Required, Int) The memory volume of an available instance(in MB), please refer to `tencentcloud_redis_zone_config.list[zone].shard_memories`. When redis is standard type, it represents total memory size of the instance; when Redis is cluster type, it represents memory size of per sharding.
* `auto_renew_flag` - (Optional, Int, ForceNew) Auto-renew flag. 0 - default state (manual renewal); 1 - automatic renewal; 2 - explicit no automatic renewal.
* `charge_type` - (Optional, String, ForceNew) The charge type of instance. Valid values: `PREPAID` and `POSTPAID`. Default value is `POSTPAID`. Note: TencentCloud International only supports `POSTPAID`. Caution that update operation on this field will delete old instances and create new with new charge type.
* `force_delete` - (Optional, Bool) Indicate whether to delete Redis instance directly or not. Default is false. If set true, the instance will be deleted instead of staying recycle bin. Note: only works for `PREPAID` instance.
* `ip` - (Optional, String) IP address of an instance. When the `operation_network` is `changeVip`, this parameter needs to be configured.
* `name` - (Optional, String) Instance name.
* `no_auth` - (Optional, Bool) Indicates whether the redis instance support no-auth access. NOTE: Only available in private cloud environment.
* `operation_network` - (Optional, String) Refers to the category of the pre-modified network, including: `changeVip`: refers to switching the private network, including its intranet IPv4 address and port; `changeVpc`: refers to switching the subnet to which the private network belongs; `changeBaseToVpc`: refers to switching the basic network to a private network; `changeVPort`: refers to only modifying the instance network port.
* `params_template_id` - (Optional, String) Specify params template id. If not set, will use default template.
* `password` - (Optional, String) Password for a Redis user, which should be 8 to 16 characters. NOTE: Only `no_auth=true` specified can make password empty.
* `port` - (Optional, Int) The port used to access a redis instance. The default value is 6379. When the `operation_network` is `changeVPort` or `changeVip`, this parameter needs to be configured.
* `prepaid_period` - (Optional, Int) The tenancy (time unit is month) of the prepaid instance, NOTE: it only works when charge_type is set to `PREPAID`. Valid values are `1`, `2`, `3`, `4`, `5`, `6`, `7`, `8`, `9`, `10`, `11`, `12`, `24`, `36`.
* `project_id` - (Optional, Int) Specifies which project the instance should belong to.
* `recycle` - (Optional, Int) Original intranet IPv4 address retention time: unit: day, value range: `0`, `1`, `2`, `3`, `7`, `15`.
* `redis_replicas_num` - (Optional, Int) The number of instance copies. This is not required for standalone and master slave versions and must equal to count of `replica_zone_ids`, Non-multi-AZ does not require `replica_zone_ids`.
* `redis_shard_num` - (Optional, Int) The number of instance shard, default is 1. This is not required for standalone and master slave versions.
* `replica_zone_ids` - (Optional, List: [`Int`]) ID of replica nodes available zone. This is not required for standalone and master slave versions. NOTE: Removing some of the same zone of replicas (e.g. removing 100001 of [100001, 100001, 100002]) will pick the first hit to remove.
* `replicas_read_only` - (Optional, Bool) Whether copy read-only is supported, Redis 2.8 Standard Edition and CKV Standard Edition do not support replica read-only, turn on replica read-only, the instance will automatically read and write separate, write requests are routed to the primary node, read requests are routed to the replica node, if you need to open replica read-only, the recommended number of replicas >=2.
* `security_groups` - (Optional, Set: [`String`]) ID of security group. If both vpc_id and subnet_id are not set, this argument should not be set either.
* `subnet_id` - (Optional, String) Specifies which subnet the instance should belong to. When the `operation_network` is `changeVpc` or `changeBaseToVpc`, this parameter needs to be configured.
* `tags` - (Optional, Map) Instance tags.
* `type_id` - (Optional, Int, ForceNew) Instance type. Available values reference data source `tencentcloud_redis_zone_config` or [document](https://intl.cloud.tencent.com/document/product/239/32069), toggle immediately when modified.
* `type` - (Optional, String, ForceNew, **Deprecated**) It has been deprecated from version 1.33.1. Please use 'type_id' instead. Instance type. Available values: `cluster_ckv`,`cluster_redis5.0`,`cluster_redis`,`master_slave_ckv`,`master_slave_redis4.0`,`master_slave_redis5.0`,`master_slave_redis`,`standalone_redis`, specific region support specific types, need to refer data `tencentcloud_redis_zone_config`.
* `vpc_id` - (Optional, String) ID of the vpc with which the instance is to be associated. When the `operation_network` is `changeVpc` or `changeBaseToVpc`, this parameter needs to be configured.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - The time when the instance was created.
* `node_info` - Readonly Primary/Replica nodes.
  * `id` - ID of the master or replica node.
  * `master` - Indicates whether the node is master.
  * `zone_id` - ID of the availability zone of the master or replica node.
* `status` - Current status of an instance, maybe: init, processing, online, isolate and todelete.


## Import

Redis instance can be imported, e.g.

```
$ terraform import tencentcloud_redis_instance.redislab redis-id
```

