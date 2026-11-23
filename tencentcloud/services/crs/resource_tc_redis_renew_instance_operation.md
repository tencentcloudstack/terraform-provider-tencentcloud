Provides a resource to create a redis renew instance operation

Example Usage

Renew Subscription Instances

```hcl
data "tencentcloud_redis_zone_config" "zone" {
  type_id = 7
}

data "tencentcloud_vpc_instances" "vpc" {
  name = "Default-VPC"
}

data "tencentcloud_vpc_subnets" "subnet" {
  vpc_id = data.tencentcloud_vpc_instances.vpc.instance_list.0.vpc_id
  availability_zone = data.tencentcloud_redis_zone_config.zone.list[1].zone
}

resource "tencentcloud_security_group" "security_group" {
  name = "tf-redis-sg"
}

resource "tencentcloud_security_group_rule_set" "example" {
  security_group_id = tencentcloud_security_group.security_group.id

  ingress {
    action     = "ACCEPT"
    cidr_block = "192.168.1.0/24"
    protocol   = "TCP"
    port       = "80"
  }

  ingress {
    action     = "DROP"
    cidr_block = "8.8.8.8"
    protocol   = "UDP"
    port       = "80,90"
  }

  ingress {
    action     = "DROP"
    cidr_block = "0.0.0.0/0"
    protocol   = "TCP"
    port       = "80-90"
  }

  egress {
    action     = "ACCEPT"
    cidr_block = "192.168.0.0/16"
    protocol   = "TCP"
    port       = "ALL"
  }

  egress {
    action     = "ACCEPT"
    cidr_block = "10.0.0.0/8"
    protocol   = "ICMP"
    port       = "ALL"
  }

  egress {
    action     = "DROP"
    cidr_block = "0.0.0.0/0"
    protocol   = "ALL"
    port       = "ALL"
  }
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
  vpc_id             = data.tencentcloud_vpc_instances.vpc.instance_list.0.vpc_id
  subnet_id          = data.tencentcloud_vpc_subnets.subnet.instance_list[0].subnet_id
  security_groups    = [tencentcloud_security_group_rule_set.example.id]
  charge_type        = "PREPAID"
  prepaid_period     = 1
}

resource "tencentcloud_redis_renew_instance_operation" "example" {
  instance_id     = tencentcloud_redis_instance.example.id
  period          = 1
  modify_pay_mode = "prepaid"
}
```