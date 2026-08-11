Provides a resource to create a PostgreSQL readonly instance.

Example Usage

Create a postgresql readonly instance

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

# create postgresql primary instance
resource "tencentcloud_postgresql_instance" "example" {
  name              = "example"
  availability_zone = var.availability_zone
  charge_type       = "POSTPAID_BY_HOUR"
  vpc_id            = tencentcloud_vpc.vpc.id
  subnet_id         = tencentcloud_subnet.subnet.id
  db_major_version  = "10"
  root_user         = "root123"
  root_password     = "Root123$"
  charset           = "UTF8"
  project_id        = 0
  memory            = 2
  cpu               = 1
  storage           = 10
}

# create postgresql readonly group
resource "tencentcloud_postgresql_readonly_group" "example" {
  master_db_instance_id       = tencentcloud_postgresql_instance.example.id
  name                        = "tf_ro_group"
  project_id                  = 0
  vpc_id                      = tencentcloud_vpc.vpc.id
  subnet_id                   = tencentcloud_subnet.subnet.id
  replay_lag_eliminate        = 1
  replay_latency_eliminate    = 1
  max_replay_lag              = 100
  max_replay_latency          = 512
  min_delay_eliminate_reserve = 1
}

# create postgresql readonly instance v2
resource "tencentcloud_postgresql_readonly_instance_v2" "example" {
  zone                  = var.availability_zone
  master_db_instance_id = tencentcloud_postgresql_instance.example.id
  spec_code             = "pgro.c2.large.ha"
  storage               = 250
  instance_count        = 1
  period                = 1
  instance_charge_type  = "POSTPAID_BY_HOUR"
  vpc_id                = tencentcloud_vpc.vpc.id
  subnet_id             = tencentcloud_subnet.subnet.id
  read_only_group_id    = tencentcloud_postgresql_readonly_group.example.id
  name                  = "tf_ro_instance_v2"
  auto_renew_flag       = 0
  need_support_ipv6     = 0
  project_id            = 0

  timeouts {
    create = "60m"
    delete = "60m"
  }
}
```

Import

postgresql readonly instance can be imported using the instance id, e.g.

```
$ terraform import tencentcloud_postgresql_readonly_instance_v2.example pgro-gih5m0ke
```
