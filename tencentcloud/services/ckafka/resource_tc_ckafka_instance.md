Use this resource to create ckafka instance.

Example Usage

Basic Instance
```hcl
variable "vpc_id" {
  default = "vpc-68vi2d3h"
}

variable "subnet_id" {
  default = "subnet-ob6clqwk"
}

data "tencentcloud_availability_zones_by_product" "gz" {
  name    = "ap-guangzhou-3"
  product = "ckafka"
}

resource "tencentcloud_ckafka_instance" "kafka_instance_prepaid" {
  instance_name      = "ckafka-instance-prepaid"
  zone_id            = data.tencentcloud_availability_zones_by_product.gz.zones.0.id
  period             = 1
  vpc_id             = var.vpc_id
  subnet_id          = var.subnet_id
  msg_retention_time = 1300
  renew_flag         = 0
  kafka_version      = "2.4.1"
  disk_size          = 200
  disk_type          = "CLOUD_BASIC"
  band_width = 20
  partition = 400

  specifications_type = "standard"
  instance_type       = 2

  config {
    auto_create_topic_enable   = true
    default_num_partitions     = 3
    default_replication_factor = 3
  }

  dynamic_retention_config {
    enable = 1
  }
}

resource "tencentcloud_ckafka_instance" "kafka_instance_postpaid" {
  instance_name      = "ckafka-instance-postpaid"
  zone_id            = data.tencentcloud_availability_zones_by_product.gz.zones.0.id
  vpc_id             = var.vpc_id
  subnet_id          = var.subnet_id
  msg_retention_time = 1300
  kafka_version      = "1.1.1"
  disk_size          = 200
  band_width         = 20
  disk_type          = "CLOUD_BASIC"
  partition          = 400
  charge_type        = "POSTPAID_BY_HOUR"

  config {
    auto_create_topic_enable   = true
    default_num_partitions     = 3
    default_replication_factor = 3
  }

  dynamic_retention_config {
    enable = 1
  }
}
```

Multi zone Instance
```hcl
variable "vpc_id" {
  default = "vpc-68vi2d3h"
}

variable "subnet_id" {
  default = "subnet-ob6clqwk"
}

data "tencentcloud_availability_zones_by_product" "gz3" {
  name    = "ap-guangzhou-3"
  product = "ckafka"
}

data "tencentcloud_availability_zones_by_product" "gz6" {
  name    = "ap-guangzhou-6"
  product = "ckafka"
}

resource "tencentcloud_ckafka_instance" "kafka_instance" {
  instance_name   = "ckafka-instance-maz-tf-test"
  zone_id         = data.tencentcloud_availability_zones_by_product.gz3.zones.0.id
  multi_zone_flag = true
  zone_ids        = [
    data.tencentcloud_availability_zones_by_product.gz3.zones.0.id,
    data.tencentcloud_availability_zones_by_product.gz6.zones.0.id
  ]
  period             = 1
  vpc_id             = var.vpc_id
  subnet_id          = var.subnet_id
  msg_retention_time = 1300
  renew_flag         = 0
  kafka_version      = "1.1.1"
  disk_size          = 500
  disk_type          = "CLOUD_BASIC"

  config {
    auto_create_topic_enable   = true
    default_num_partitions     = 3
    default_replication_factor = 3
  }

  dynamic_retention_config {
    enable = 1
  }
}
```

Import

ckafka instance can be imported using the instance_id, e.g.

```
$ terraform import tencentcloud_ckafka_instance.foo ckafka-f9ife4zz
```