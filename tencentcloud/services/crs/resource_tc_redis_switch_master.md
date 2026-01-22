Provides a resource to create a redis switch master

Example Usage

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

resource "tencentcloud_security_group" "security_group" {
  name = "tf-redis-sg"
}

resource "tencentcloud_security_group_lite_rule" "example" {
  security_group_id = tencentcloud_security_group.security_group.id

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

resource "tencentcloud_redis_instance" "example" {
  availability_zone  = data.tencentcloud_redis_zone_config.zone.list[2].zone
  type_id            = data.tencentcloud_redis_zone_config.zone.list[2].type_id
  password           = "Password@123"
  mem_size           = 8192
  redis_shard_num    = data.tencentcloud_redis_zone_config.zone.list[2].redis_shard_nums[0]
  redis_replicas_num = 2
  replica_zone_ids   = var.replica_zone_ids
  name               = "tf_example"
  port               = 6379
  vpc_id             = tencentcloud_vpc.vpc.id
  subnet_id          = tencentcloud_subnet.subnet.id
  security_groups    = [tencentcloud_security_group.security_group.id]
}

data "tencentcloud_redis_instance_zone_info" "example" {
  instance_id = tencentcloud_redis_instance.example.id
}

resource "tencentcloud_redis_switch_master" "example" {
  instance_id = tencentcloud_redis_instance.example.id
  group_id    = data.tencentcloud_redis_instance_zone_info.example.replica_groups[1].group_id
}
```