---
subcategory: "TencentDB for Redis(crs)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_redis_connection_config"
sidebar_current: "docs-tencentcloud-resource-redis_connection_config"
description: |-
  Provides a resource to create a redis connection_config
---

# tencentcloud_redis_connection_config

Provides a resource to create a redis connection_config

## Example Usage

### Modify the maximum number of connections and maximum network throughput of an instance

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

resource "tencentcloud_redis_connection_config" "connection_config" {
  instance_id   = "crs-fhm9fnv1"
  client_limit  = "20000"
  add_bandwidth = "30"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) The ID of instance.
* `add_bandwidth` - (Optional, Int) Refers to the additional bandwidth of the instance. When the standard bandwidth does not meet the demand, the user can increase the bandwidth by himself. When the read-only copy is enabled, the total bandwidth of the instance = additional bandwidth * number of fragments + standard bandwidth * number of fragments * Max ([number of read-only replicas, 1] ), the number of shards in the standard architecture = 1, and when read-only replicas are not enabled, the total bandwidth of the instance = additional bandwidth * number of shards + standard bandwidth * number of shards, and the number of shards in the standard architecture = 1.
* `client_limit` - (Optional, Int) The total number of connections per shard.If read-only replicas are not enabled, the lower limit is 10,000 and the upper limit is 40,000.When you enable read-only replicas, the minimum limit is 10,000 and the upper limit is 10,000 * (the number of read replicas +3).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `base_bandwidth` - standard bandwidth. Refers to the bandwidth allocated by the system to each node when an instance is purchased.
* `max_add_bandwidth` - Additional bandwidth is capped.
* `min_add_bandwidth` - Additional bandwidth sets the lower limit.
* `total_bandwidth` - Total bandwidth of the instance = additional bandwidth * number of shards + standard bandwidth * number of shards * (number of primary nodes + number of read-only replica nodes), the number of shards of the standard architecture = 1, in Mb/s.


## Import

Redis connectionConfig can be imported, e.g.

```
$ terraform import tencentcloud_redis_connection_config.connection_config instance_id
```

