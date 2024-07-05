Use this resource to create postgresql readonly group.

Example Usage

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-3"
}

# create vpc
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc"
  cidr_block = "10.0.0.0/16"
}

# create vpc subnet
resource "tencentcloud_subnet" "subnet" {
  availability_zone = var.availability_zone
  name              = "subnet"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.20.0/28"
  is_multicast      = false
}

# create postgresql
resource "tencentcloud_postgresql_instance" "example" {
  name              = "example"
  availability_zone = var.availability_zone
  charge_type       = "POSTPAID_BY_HOUR"
  vpc_id            = tencentcloud_vpc.vpc.id
  subnet_id         = tencentcloud_subnet.subnet.id
  engine_version    = "10.4"
  root_user         = "root123"
  root_password     = "Root123$"
  charset           = "UTF8"
  project_id        = 0
  memory            = 4
  cpu               = 2
  storage           = 50

  tags = {
    test = "tf"
  }
}

# create security group
resource "tencentcloud_security_group" "example" {
  name        = "tf-example"
  description = "sg desc."
  project_id  = 0

  tags = {
    "example" = "test"
  }
}

resource "tencentcloud_postgresql_readonly_group" "example" {
  master_db_instance_id       = tencentcloud_postgresql_instance.example.id
  name                        = "tf_ro_group"
  project_id                  = 0
  vpc_id                      = tencentcloud_vpc.vpc.id
  subnet_id                   = tencentcloud_subnet.subnet.id
  security_groups_ids         = [tencentcloud_security_group.example.id]
  replay_lag_eliminate        = 1
  replay_latency_eliminate    = 1
  max_replay_lag              = 100
  max_replay_latency          = 512
  min_delay_eliminate_reserve = 1
}
```

Import

postgresql readonly group can be imported, e.g.

```
$ terraform import tencentcloud_postgresql_readonly_group.example pgrogrp-lckioi2a
```
