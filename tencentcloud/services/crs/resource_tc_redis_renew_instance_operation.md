Provides a resource to create a redis renew_instance_operation

Example Usage

Renew Subscription Instances

```hcl
data "tencentcloud_redis_zone_config" "zone" {
  type_id = 7
}

data "tencentcloud_vpc" "vpc" {
  name = "Default-VPC"
}

data "tencentcloud_vpc_subnets" "subnet" {
  vpc_id = data.tencentcloud_vpc.vpc.id
  availability_zone = data.tencentcloud_redis_zone_config.zone.list[1].zone
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
  vpc_id             = data.tencentcloud_vpc.vpc.id
  subnet_id          = data.tencentcloud_vpc_subnets.subnet.instance_list[0].subnet_id
  security_groups    = [tencentcloud_security_group.foo.id]
  charge_type        = "PREPAID"
  prepaid_period     = 1
}

resource "tencentcloud_redis_renew_instance_operation" "foo" {
  instance_id     = tencentcloud_redis_instance.foo.id
  period          = 1
  modify_pay_mode = "prepaid"
}
```