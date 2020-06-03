---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_redis_instances"
sidebar_current: "docs-tencentcloud-datasource-redis_instances"
description: |-
  Use this data source to query the detail information of redis instance.
---

# tencentcloud_redis_instances

Use this data source to query the detail information of redis instance.

## Example Usage

```hcl
data "tencentcloud_redis_instances" "redislab" {
  zone               = "ap-hongkong-1"
  search_key         = "myredis"
  project_id         = 0
  limit              = 20
  result_output_file = "/tmp/redis_instances"
}
```

## Argument Reference

The following arguments are supported:

* `limit` - (Optional) The number limitation of results for a query.
* `project_id` - (Optional) ID of the project to which redis instance belongs.
* `result_output_file` - (Optional) Used to save results.
* `search_key` - (Optional) Key words used to match the results, and the key words can be: instance ID, instance name and IP address.
* `tags` - (Optional) Tags of redis instance.
* `zone` - (Optional) ID of an available zone.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instance_list` - A list of redis instance. Each element contains the following attributes:
  * `charge_type` - The charge type of instance. Valid values are `POSTPAID` and `PREPAID`. Default value is `POSTPAID`.
  * `create_time` - The time when the instance is created.
  * `ip` - IP address of an instance.
  * `mem_size` - Memory size in MB.
  * `name` - Name of a redis instance.
  * `port` - The port used to access a redis instance.
  * `project_id` - ID of the project to which a redis instance belongs.
  * `redis_id` - ID of a redis instance.
  * `redis_replicas_num` - The number of instance copies.
  * `redis_shard_num` - The number of instance shard.
  * `status` - Current status of an instance, maybe: init, processing, online, isolate and todelete.
  * `subnet_id` - ID of the vpc subnet.
  * `tags` - Tags of an instance.
  * `type_id` - Instance type. Refer to `data.tencentcloud_redis_zone_config.list.type_id` get available values.
  * `type` - (**Deprecated**) It has been deprecated from version 1.33.1. Please use 'type_id' instead. Instance type. Available values: master_slave_redis, master_slave_ckv, cluster_ckv, cluster_redis and standalone_redis.
  * `vpc_id` - ID of the vpc with which the instance is associated.
  * `zone` - Available zone to which a redis instance belongs.


