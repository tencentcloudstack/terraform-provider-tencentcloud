---
subcategory: "TencentDB for PostgreSQL(PostgreSQL)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_postgresql_readonly_instance"
sidebar_current: "docs-tencentcloud-resource-postgresql_readonly_instance"
description: |-
  Use this resource to create postgresql readonly instance.
---

# tencentcloud_postgresql_readonly_instance

Use this resource to create postgresql readonly instance.

## Example Usage

### Create postgresql readonly instance

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
  db_major_version  = "10"
  root_user         = "root123"
  root_password     = "Root123$"
  charset           = "UTF8"
  project_id        = 0
  memory            = 2
  cpu               = 1
  storage           = 10

  tags = {
    test = "tf"
  }
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

# create security group
resource "tencentcloud_security_group" "example" {
  name        = "tf-example"
  description = "sg desc."
  project_id  = 0

  tags = {
    "example" = "test"
  }
}

# create postgresql readonly instance
resource "tencentcloud_postgresql_readonly_instance" "example" {
  read_only_group_id    = tencentcloud_postgresql_readonly_group.example.id
  master_db_instance_id = tencentcloud_postgresql_instance.example.id
  zone                  = var.availability_zone
  name                  = "example"
  auto_renew_flag       = 0
  db_version            = "10.23"
  instance_charge_type  = "POSTPAID_BY_HOUR"
  memory                = 4
  cpu                   = 2
  storage               = 250
  vpc_id                = tencentcloud_vpc.vpc.id
  subnet_id             = tencentcloud_subnet.subnet.id
  need_support_ipv6     = 0
  project_id            = 0
  security_groups_ids = [
    tencentcloud_security_group.example.id,
  ]
}
```

### Create postgresql readonly instance of CDC

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-4"
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
  name              = "tf-example"
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

  db_node_set {
    role                 = "Primary"
    zone                 = var.availability_zone
    dedicated_cluster_id = "cluster-262n63e8"
  }

  db_node_set {
    zone                 = var.availability_zone
    dedicated_cluster_id = "cluster-262n63e8"
  }

  tags = {
    CreateBy = "terraform"
  }
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

# create security group
resource "tencentcloud_security_group" "example" {
  name        = "tf-example"
  description = "sg desc."
  project_id  = 0

  tags = {
    CreateBy = "terraform"
  }
}

# create postgresql readonly instance
resource "tencentcloud_postgresql_readonly_instance" "example" {
  read_only_group_id    = tencentcloud_postgresql_readonly_group.example.id
  master_db_instance_id = tencentcloud_postgresql_instance.example.id
  zone                  = var.availability_zone
  name                  = "example"
  auto_renew_flag       = 0
  db_version            = "10.23"
  instance_charge_type  = "POSTPAID_BY_HOUR"
  memory                = 4
  cpu                   = 2
  storage               = 250
  vpc_id                = tencentcloud_vpc.vpc.id
  subnet_id             = tencentcloud_subnet.subnet.id
  need_support_ipv6     = 0
  project_id            = 0
  dedicated_cluster_id  = "cluster-262n63e8"
  security_groups_ids = [
    tencentcloud_security_group.example.id,
  ]
}
```

## Argument Reference

The following arguments are supported:

* `db_version` - (Required, String, ForceNew) PostgreSQL kernel version, which must be the same as that of the primary instance.
* `master_db_instance_id` - (Required, String, ForceNew) ID of the primary instance to which the read-only replica belongs.
* `memory` - (Required, Int) Memory size(in GB). Allowed value must be larger than `memory` that data source `tencentcloud_postgresql_specinfos` provides.
* `name` - (Required, String) Instance name.
* `project_id` - (Required, Int) Project ID.
* `security_groups_ids` - (Required, Set: [`String`]) ID of security group.
* `storage` - (Required, Int) Instance storage capacity in GB.
* `subnet_id` - (Required, String) VPC subnet ID.
* `vpc_id` - (Required, String, ForceNew) VPC ID.
* `zone` - (Required, String, ForceNew) Availability zone ID, which can be obtained through the Zone field in the returned value of the DescribeZones API.
* `auto_renew_flag` - (Optional, Int) Auto renew flag, `1` for enabled. NOTES: Only support prepaid instance.
* `auto_voucher` - (Optional, Int) Whether to use voucher, `1` for enabled.
* `cpu` - (Optional, Int) Number of CPU cores. Allowed value must be equal `cpu` that data source `tencentcloud_postgresql_specinfos` provides.
* `dedicated_cluster_id` - (Optional, String) Dedicated cluster ID.
* `instance_charge_type` - (Optional, String, ForceNew) instance billing mode. Valid values: PREPAID (monthly subscription), POSTPAID_BY_HOUR (pay-as-you-go).
* `need_support_ipv6` - (Optional, Int, ForceNew) Whether to support IPv6 address access. Valid values: 1 (yes), 0 (no).
* `period` - (Optional, Int) Specify Prepaid period in month. Default `1`. Values: `1`, `2`, `3`, `4`, `5`, `6`, `7`, `8`, `9`, `10`, `11`, `12`, `24`, `36`.
* `read_only_group_id` - (Optional, String) RO group ID.
* `tags` - (Optional, Map) Tags.
* `voucher_ids` - (Optional, List: [`String`]) Specify Voucher Ids if `auto_voucher` was `1`, only support using 1 vouchers for now.
* `wait_switch` - (Optional, Int) Switch time after instance configurations are modified. `0`: Switch immediately; `2`: Switch during maintenance time window. Default: `0`. Note: This only takes effect when updating the `memory`, `storage`, `cpu` fields.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Create time of the postgresql instance.
* `instance_id` - The instance ID of this readonly resource.
* `private_access_ip` - IP for private access.
* `private_access_port` - Port for private access.


## Import

postgresql readonly instance can be imported using the id, e.g.

```
$ terraform import tencentcloud_postgresql_readonly_instance.example pgro-gih5m0ke
```

