Provide a resource to create a SSM secret.

Example Usage

Create user defined secret

```hcl
resource "tencentcloud_ssm_secret" "example" {
  secret_name             = "tf-example"
  description             = "desc."
  is_enabled              = true
  recovery_window_in_days = 0

  tags = {
    createBy = "terraform"
  }
}
```

Create redis secret

```hcl
data "tencentcloud_redis_zone_config" "zone" {
  type_id = 8
}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = data.tencentcloud_redis_zone_config.zone.list[3].zone
  name              = "subnet-example"
  cidr_block        = "10.0.0.0/16"
}

resource "tencentcloud_redis_instance" "example" {
  availability_zone  = data.tencentcloud_redis_zone_config.zone.list[3].zone
  type_id            = data.tencentcloud_redis_zone_config.zone.list[3].type_id
  password           = "Qwer@234"
  mem_size           = data.tencentcloud_redis_zone_config.zone.list[3].mem_sizes[0]
  redis_shard_num    = data.tencentcloud_redis_zone_config.zone.list[3].redis_shard_nums[0]
  redis_replicas_num = data.tencentcloud_redis_zone_config.zone.list[3].redis_replicas_nums[0]
  name               = "tf_example"
  port               = 6379
  vpc_id             = tencentcloud_vpc.vpc.id
  subnet_id          = tencentcloud_subnet.subnet.id
}

resource "tencentcloud_ssm_secret" "example" {
  secret_name       = "tf-example"
  description       = "redis desc."
  is_enabled        = true
  secret_type       = 4
  additional_config = jsonencode(
    {
      "Region" : "ap-guangzhou"
      "Privilege" : "r",
      "InstanceId" : tencentcloud_redis_instance.example.id
      "ReadonlyPolicy" : ["master"],
      "Remark" : "for tf test"
    }
  )
  tags = {
    createdBy = "terraform"
  }
  recovery_window_in_days = 0
}
```

Import

SSM secret can be imported using the secretName, e.g.
```
$ terraform import tencentcloud_ssm_secret.foo test
```