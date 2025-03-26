Provides a resource to create a Redis instance and set its attributes.

~> **NOTE:** The argument vpc_id and subnet_id is now required because Basic Network Instance is no longer supported.

~> **NOTE:** Both adding and removing replications in one change is supported but not recommend.

Example Usage

Create a base version of redis

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
  name               = "tf-example"
  port               = 6379
  vpc_id             = tencentcloud_vpc.vpc.id
  subnet_id          = tencentcloud_subnet.subnet.id
}
```

Using multi replica zone set

```hcl
variable "redis_replicas_num" {
  default = 3
}

variable "redis_type_id" {
  default = 7
}

data "tencentcloud_availability_zones_by_product" "az" {
  product = "redis"
}

resource "tencentcloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
  name       = "tf_redis_vpc"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = data.tencentcloud_availability_zones_by_product.az.zones[0].name
  name              = "tf_redis_subnet"
  cidr_block        = "10.0.1.0/24"
}

resource "tencentcloud_security_group" "security_group" {
  name = "tf-redis-sg"
}

resource "tencentcloud_security_group_lite_rule" "sg_rule" {
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
  availability_zone  = data.tencentcloud_availability_zones_by_product.az.zones[0].name
  type_id            = var.redis_type_id
  charge_type        = "POSTPAID"
  mem_size           = 1024
  name               = "tf-example"
  port               = 6379
  project_id         = 0
  vpc_id             = tencentcloud_vpc.vpc.id
  subnet_id          = tencentcloud_subnet.subnet.id
  password           = "Password@123"
  security_groups    = [tencentcloud_security_group.security_group.id]
  redis_replicas_num = var.redis_replicas_num
  redis_shard_num    = 1
  replica_zone_ids = [
    for i in range(var.redis_replicas_num)
    : data.tencentcloud_availability_zones_by_product.az.zones[i % length(data.tencentcloud_availability_zones_by_product.az.zones)].id
  ]
}
```

Buy a month of prepaid instances

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

resource "tencentcloud_security_group" "security_group" {
  name = "tf-redis-sg"
}

resource "tencentcloud_security_group_lite_rule" "sg_rule" {
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
  availability_zone  = data.tencentcloud_redis_zone_config.zone.list[0].zone
  type_id            = data.tencentcloud_redis_zone_config.zone.list[0].type_id
  password           = "Password@123"
  mem_size           = 8192
  redis_shard_num    = data.tencentcloud_redis_zone_config.zone.list[0].redis_shard_nums[0]
  redis_replicas_num = data.tencentcloud_redis_zone_config.zone.list[0].redis_replicas_nums[0]
  name               = "tf-example"
  port               = 6379
  vpc_id             = tencentcloud_vpc.vpc.id
  subnet_id          = tencentcloud_subnet.subnet.id
  security_groups    = [tencentcloud_security_group.security_group.id]
  charge_type        = "PREPAID"
  prepaid_period     = 1
}
```

Create a multi-AZ instance

```hcl
data "tencentcloud_redis_zone_config" "zone" {
  type_id = 7
  region = "ap-guangzhou"
}

variable "replica_zone_ids" {
  default = [100004,100006]
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

resource "tencentcloud_security_group_lite_rule" "sg_rule" {
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
  name               = "tf-example"
  port               = 6379
  vpc_id             = tencentcloud_vpc.vpc.id
  subnet_id          = tencentcloud_subnet.subnet.id
  security_groups    = [tencentcloud_security_group.security_group.id]
}
```

Create a CDC scenario instance

```hcl
variable "cdc_id" {
  default = "cluster-262n63e8"
}

variable "cdc_region" {
  default = "ap-guangzhou"
}

data "tencentcloud_redis_clusters" "clusters" {
  dedicated_cluster_id = var.cdc_id
}

output "name" {
  value = data.tencentcloud_redis_clusters.clusters.resources[0].redis_cluster_id
}

data "tencentcloud_redis_zone_config" "zone" {
  type_id = 7
  region  = var.cdc_region
}

data "tencentcloud_cdc_dedicated_clusters" "example" {
  dedicated_cluster_ids = [var.cdc_id]
}

data "tencentcloud_vpc_subnets" "subnets" {
  cdc_id = var.cdc_id
}

resource "tencentcloud_redis_instance" "example" {
  availability_zone  = data.tencentcloud_cdc_dedicated_clusters.example.dedicated_cluster_set[0].zone
  type_id            = data.tencentcloud_redis_zone_config.zone.list[0].type_id
  password           = "Password@123"
  mem_size           = 8192
  redis_shard_num    = data.tencentcloud_redis_zone_config.zone.list[0].redis_shard_nums[0]
  redis_replicas_num = data.tencentcloud_redis_zone_config.zone.list[0].redis_replicas_nums[0]
  name               = "tf-cdc-example"
  port               = 6379
  vpc_id             = data.tencentcloud_vpc_subnets.subnets.instance_list[0].vpc_id
  subnet_id          = data.tencentcloud_vpc_subnets.subnets.instance_list[0].subnet_id
  product_version    = "cdc"
  redis_cluster_id   = data.tencentcloud_redis_clusters.clusters.resources[0].redis_cluster_id
}
```

Import

Redis instance can be imported, e.g.

```
$ terraform import tencentcloud_redis_instance.example crs-iu22tdrf
```