---
subcategory: "TencentDB for Redis(crs)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_redis_instance_password_policy_config"
sidebar_current: "docs-tencentcloud-resource-redis_instance_password_policy_config"
description: |-
  Provides a resource to manage redis instance password policy config
---

# tencentcloud_redis_instance_password_policy_config

Provides a resource to manage redis instance password policy config

## Example Usage

### Manage the password complexity policy of a redis instance

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

resource "tencentcloud_redis_instance" "example" {
  availability_zone  = data.tencentcloud_redis_zone_config.zone.list[0].zone
  type_id            = data.tencentcloud_redis_zone_config.zone.list[0].type_id
  password           = "Password@123"
  mem_size           = 8192
  redis_shard_num    = data.tencentcloud_redis_zone_config.zone.list[0].redis_shard_nums[0]
  redis_replicas_num = data.tencentcloud_redis_zone_config.zone.list[0].redis_replicas_nums[0]
  name               = "tf_example"
  port               = 6379
  vpc_id             = tencentcloud_vpc.vpc.id
  subnet_id          = tencentcloud_subnet.subnet.id
}

resource "tencentcloud_redis_instance_password_policy_config" "example" {
  instance_id       = tencentcloud_redis_instance.example.id
  enabled           = true
  min_letter_count  = 1
  min_digit_count   = 1
  min_special_count = 1
  min_length        = 8
}
```

## Argument Reference

The following arguments are supported:

* `enabled` - (Required, Bool) Whether to enable the instance-level password complexity policy. true: enable; false: disable.
* `instance_id` - (Required, String, ForceNew) The ID of redis instance.
* `min_digit_count` - (Optional, Int) The minimum number of digit characters. Value range: [1,16].
* `min_length` - (Optional, Int) The minimum total length of the password. Value range: [8,64].
* `min_letter_count` - (Optional, Int) The minimum number of letters (uppercase and lowercase). Value range: [1,16].
* `min_special_count` - (Optional, Int) The minimum number of special characters. Value range: [1,16].

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

redis instance password policy config can be imported using the instance_id, e.g.

```
$ terraform import tencentcloud_redis_instance_password_policy_config.example crs-cqdfdzvt
```

