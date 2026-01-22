Provides a resource to create a Clickhouse instance.

Example Usage

Create POSTPAID instance

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-6"
}

data "tencentcloud_clickhouse_spec" "spec" {
  zone       = var.availability_zone
  pay_mode   = "POSTPAID_BY_HOUR"
  is_elastic = false
}

locals {
  data_spec              = [for i in data.tencentcloud_clickhouse_spec.spec.data_spec : i if lookup(i, "cpu") == 4 && lookup(i, "mem") == 16]
  data_spec_name_4c16m   = local.data_spec.0.name
  common_spec            = [for i in data.tencentcloud_clickhouse_spec.spec.common_spec : i if lookup(i, "cpu") == 4 && lookup(i, "mem") == 16]
  common_spec_name_4c16m = local.common_spec.0.name
}

resource "tencentcloud_vpc" "vpc" {
  name       = "cdwch-vpc"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  name              = "cdwch-subnet"
  cidr_block        = "10.0.0.0/16"
  availability_zone = var.availability_zone
  is_multicast      = false
}

resource "tencentcloud_clickhouse_instance" "example" {
  instance_name       = "tf-example"
  charge_type         = "POSTPAID_BY_HOUR"
  zone                = var.availability_zone
  ha_flag             = true
  vpc_id              = tencentcloud_vpc.vpc.id
  subnet_id           = tencentcloud_subnet.subnet.id
  product_version     = "21.8.12.29"
  ck_default_user_pwd = "Password@123"
  data_spec {
    spec_name = local.data_spec_name_4c16m
    count     = 2
    disk_size = 300
  }
  
  common_spec {
    spec_name = local.common_spec_name_4c16m
    count     = 3
    disk_size = 300
  }
}
```

Create PREPAID instance

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-6"
}

data "tencentcloud_clickhouse_spec" "spec" {
  zone       = var.availability_zone
  pay_mode   = "POSTPAID_BY_HOUR"
  is_elastic = false
}

locals {
  data_spec              = [for i in data.tencentcloud_clickhouse_spec.spec.data_spec : i if lookup(i, "cpu") == 4 && lookup(i, "mem") == 16]
  data_spec_name_4c16m   = local.data_spec.0.name
  common_spec            = [for i in data.tencentcloud_clickhouse_spec.spec.common_spec : i if lookup(i, "cpu") == 4 && lookup(i, "mem") == 16]
  common_spec_name_4c16m = local.common_spec.0.name
}

resource "tencentcloud_vpc" "vpc" {
  name       = "cdwch-vpc"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  name              = "cdwch-subnet"
  cidr_block        = "10.0.0.0/16"
  availability_zone = var.availability_zone
  is_multicast      = false
}

resource "tencentcloud_clickhouse_instance" "example" {
  instance_name       = "tf-example"
  charge_type         = "PREPAID"
  renew_flag          = 1
  time_span           = 1
  zone                = var.availability_zone
  ha_flag             = true
  vpc_id              = tencentcloud_vpc.vpc.id
  subnet_id           = tencentcloud_subnet.subnet.id
  product_version     = "21.8.12.29"
  ck_default_user_pwd = "Password@123"
  data_spec {
    spec_name = local.data_spec_name_4c16m
    count     = 2
    disk_size = 300
  }

  common_spec {
    spec_name = local.common_spec_name_4c16m
    count     = 3
    disk_size = 300
  }
}
```

Import

Clickhouse instance can be imported using the id, e.g.

```
$ terraform import tencentcloud_clickhouse_instance.example cdwch-4l6mm8p7
```
