Use this resource to create CKafka instance.

Example Usage

Create basic instance(prepaid)

```hcl
data "tencentcloud_availability_zones_by_product" "gz" {
  name    = "ap-guangzhou-6"
  product = "ckafka"
}

# create vpc
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc"
  cidr_block = "10.0.0.0/16"
}

# create vpc subnet
resource "tencentcloud_subnet" "subnet" {
  name              = "subnet"
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = "ap-guangzhou-6"
  cidr_block        = "10.0.20.0/28"
  is_multicast      = false
}

# create ckafka
resource "tencentcloud_ckafka_instance" "example" {
  instance_name       = "tf-example"
  zone_id             = data.tencentcloud_availability_zones_by_product.gz.zones.0.id
  period              = 1
  vpc_id              = tencentcloud_vpc.vpc.id
  subnet_id           = tencentcloud_subnet.subnet.id
  msg_retention_time  = 1300
  renew_flag          = 0
  kafka_version       = "2.8.1"
  disk_size           = 200
  disk_type           = "CLOUD_BASIC"
  band_width          = 40
  partition           = 400
  specifications_type = "profession"
  instance_type       = 1

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

Create multi zone instance(postpaid)

```hcl
data "tencentcloud_availability_zones_by_product" "gz6" {
  name    = "ap-guangzhou-6"
  product = "ckafka"
}

data "tencentcloud_availability_zones_by_product" "gz7" {
  name    = "ap-guangzhou-7"
  product = "ckafka"
}

# create vpc
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc"
  cidr_block = "10.0.0.0/16"
}

# create vpc subnet
resource "tencentcloud_subnet" "subnet" {
  name              = "subnet"
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = "ap-guangzhou-6"
  cidr_block        = "10.0.20.0/28"
  is_multicast      = false
}

# create ckafka
resource "tencentcloud_ckafka_instance" "example" {
  instance_name   = "tf-example"
  zone_id         = data.tencentcloud_availability_zones_by_product.gz6.zones.0.id
  multi_zone_flag = true
  zone_ids = [
    data.tencentcloud_availability_zones_by_product.gz6.zones.0.id,
    data.tencentcloud_availability_zones_by_product.gz7.zones.0.id,
  ]
  renew_flag          = 0
  vpc_id              = tencentcloud_vpc.vpc.id
  subnet_id           = tencentcloud_subnet.subnet.id
  msg_retention_time  = 4320
  kafka_version       = "2.8.1"
  disk_size           = 200
  disk_type           = "CLOUD_BASIC"
  band_width          = 20
  partition           = 400
  specifications_type = "profession"
  charge_type         = "POSTPAID_BY_HOUR"
  instance_type       = 1

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

CKafka instance can be imported using the instanceId, e.g.

```
$ terraform import tencentcloud_ckafka_instance.example ckafka-f9ife4zz
```
