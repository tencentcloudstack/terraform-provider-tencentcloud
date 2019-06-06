---
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
resource "tencentcloud_redis_instance" "redis_instance_test"{
	availability_zone="ap-hongkong-3"
	type="master_slave_redis"
	password="test12345789"
	mem_size=8192
	name="terrform_test"
	port=6379
}
```

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Required, ForceNew) The available zone ID of an instance to be created., refer to tencentcloud_redis_zone_config.list
* `mem_size` - (Required) The memory volume of an available instance(in MB), refer to tencentcloud_redis_zone_config.list[zone].mem_sizes
* `password` - (Required) Password for a Redis user，which should be 8 to 16 characters.
* `name` - (Optional) Instance name.
* `port` - (Optional, ForceNew) The port used to access a redis instance. The default value is 6379. And this value can't be changed after creation, or the Redis instance will be recreated.
* `project_id` - (Optional) Specifies which project the instance should belong to.
* `security_groups` - (Optional, ForceNew) ID of security group. If both vpc_id and subnet_id are not set, this argument should not be set either. 
* `subnet_id` - (Optional, ForceNew) Specifies which subnet the instance should belong to.
* `type` - (Optional, ForceNew) Instance type. Available values: master_slave_redis.
* `vpc_id` - (Optional, ForceNew) ID of the vpc with which the instance is to be associated.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `create_time` -  The time when the instance was created.
* `ip` - IP address of an instance.
* `status` - Current status of an instance，maybe: init, processing, online, isolate and todelete.


## Import

Redis instance can be imported, e.g.

```hcl
$ terraform import tencentcloud_redis_instance.redislab redis-id
```

