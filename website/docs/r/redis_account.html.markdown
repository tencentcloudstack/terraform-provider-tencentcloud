---
subcategory: "TencentDB for Redis(crs)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_redis_account"
sidebar_current: "docs-tencentcloud-resource-redis_account"
description: |-
  Provides a resource to create a redis account
---

# tencentcloud_redis_account

Provides a resource to create a redis account

## Example Usage

### Create an account with read and write permissions

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

resource "tencentcloud_redis_instance" "foo" {
  availability_zone  = data.tencentcloud_redis_zone_config.zone.list[1].zone
  type_id            = data.tencentcloud_redis_zone_config.zone.list[1].type_id
  password           = "test12345789"
  mem_size           = 8192
  redis_shard_num    = data.tencentcloud_redis_zone_config.zone.list[1].redis_shard_nums[0]
  redis_replicas_num = data.tencentcloud_redis_zone_config.zone.list[1].redis_replicas_nums[0]
  name               = "terrform_test"
  port               = 6379
  vpc_id             = tencentcloud_vpc.vpc.id
  subnet_id          = tencentcloud_subnet.subnet.id
}

resource "tencentcloud_redis_account" "account" {
  instance_id      = tencentcloud_redis_instance.foo.id
  account_name     = "account_test"
  account_password = "test1234"
  remark           = "master"
  readonly_policy  = ["master"]
  privilege        = "rw"
}
```

### Create an account with read-only permissions

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

resource "tencentcloud_redis_instance" "foo" {
  availability_zone  = data.tencentcloud_redis_zone_config.zone.list[1].zone
  type_id            = data.tencentcloud_redis_zone_config.zone.list[1].type_id
  password           = "test12345789"
  mem_size           = 8192
  redis_shard_num    = data.tencentcloud_redis_zone_config.zone.list[1].redis_shard_nums[0]
  redis_replicas_num = data.tencentcloud_redis_zone_config.zone.list[1].redis_replicas_nums[0]
  name               = "terrform_test"
  port               = 6379
  vpc_id             = tencentcloud_vpc.vpc.id
  subnet_id          = tencentcloud_subnet.subnet.id
}

resource "tencentcloud_redis_account" "account" {
  instance_id      = tencentcloud_redis_instance.foo.id
  account_name     = "account_test"
  account_password = "test1234"
  remark           = "master"
  readonly_policy  = ["master"]
  privilege        = "r"
}
```

## Argument Reference

The following arguments are supported:

* `account_name` - (Required, String) The account name.
* `account_password` - (Required, String) 1: Length 8-30 digits, it is recommended to use a password of more than 12 digits; 2: Cannot start with `/`; 3: Include at least two items: a.Lowercase letters `a-z`; b.Uppercase letters `A-Z` c.Numbers `0-9`;  d.`()`~!@#$%^&*-+=_|{}[]:;<>,.?/`.
* `instance_id` - (Required, String) The ID of instance.
* `privilege` - (Required, String) Read and write policy: Enter R and RW to indicate read-only, read-write, cannot be empty when modifying operations.
* `readonly_policy` - (Required, Set: [`String`]) Routing policy: Enter master or replication, which indicates the master node or slave node, cannot be empty when modifying operations.
* `remark` - (Optional, String) Remark.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

redis account can be imported using the id, e.g.

```
terraform import tencentcloud_redis_account.account crs-xxxxxx#account_test
```

