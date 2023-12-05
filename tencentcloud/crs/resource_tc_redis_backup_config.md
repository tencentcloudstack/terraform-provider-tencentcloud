Use this resource to create a backup config of redis.

Example Usage

Set configuration for automatic backups

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

resource "tencentcloud_redis_backup_config" "foo" {
  redis_id      = tencentcloud_redis_instance.foo.id
  backup_time   = "04:00-05:00"
  backup_period = ["Monday"]
}
```

Import

Redis  backup config can be imported, e.g.

```
$ terraform import tencentcloud_redis_backup_config.foo redis-id
```