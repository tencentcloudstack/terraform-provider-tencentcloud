---
subcategory: "TencentDB for PostgreSQL(PostgreSQL)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_postgresql_instance"
sidebar_current: "docs-tencentcloud-resource-postgresql_instance"
description: |-
  Use this resource to create postgresql instance.
---

# tencentcloud_postgresql_instance

Use this resource to create postgresql instance.

-> **Note:** To update the charge type, please update the `charge_type` and specify the `period` for the charging period. It only supports updating from `POSTPAID_BY_HOUR` to `PREPAID`, and the `period` field only valid in that upgrading case.

-> **Note:** If no values are set for the parameters: `db_kernel_version`, `db_major_version` and `engine_version`, then `engine_version` is set to `10.4` by default. Suggest using parameter `db_major_version` to create an instance

-> **Note:** If you need to upgrade the database version, Please use data source `tencentcloud_postgresql_db_versions` to obtain the valid version value for `db_kernel_version`, `db_major_version` and `engine_version`. And when modifying, `db_kernel_version`, `db_major_version` and `engine_version` must be set.

-> **Note:** If upgrade `db_kernel_version`, will synchronize the upgrade of the read-only instance version; If upgrade `db_major_version`, cannot have read-only instances.

## Example Usage

### Create a postgresql instance

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
  engine_version    = "10.23"
  root_user         = "root123"
  root_password     = "Root123$"
  charset           = "UTF8"
  project_id        = 0
  cpu               = 1
  memory            = 2
  storage           = 10

  tags = {
    CreateBy = "Terraform"
  }
}
```

### Create a postgresql instance with delete protection

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
  engine_version    = "10.23"
  root_user         = "root123"
  root_password     = "Root123$"
  charset           = "UTF8"
  project_id        = 0
  cpu               = 1
  memory            = 2
  storage           = 10
  delete_protection = true

  tags = {
    CreateBy = "Terraform"
  }
}
```

### Create a multi available zone postgresql instance

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-6"
}

variable "standby_availability_zone" {
  default = "ap-guangzhou-7"
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

  db_node_set {
    role = "Primary"
    zone = var.availability_zone
  }

  db_node_set {
    zone = var.standby_availability_zone
  }

  tags = {
    CreateBy = "Terraform"
  }
}
```

### Create a multi available zone postgresql instance of CDC

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
    CreateBy = "Terraform"
  }
}
```

### Create pgsql with kms key

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-6"
}

resource "tencentcloud_postgresql_instance" "example" {
  name              = "tf_postsql_instance"
  availability_zone = var.availability_zone
  charge_type       = "POSTPAID_BY_HOUR"
  vpc_id            = "vpc-86v957zb"
  subnet_id         = "subnet-enm92y0m"
  db_major_version  = "11"
  engine_version    = "11.12"
  db_kernel_version = "v11.12_r1.3"
  need_support_tde  = 1
  kms_key_id        = "788c606a-c7b7-11ec-82d1-5254001e5c4e"
  kms_region        = "ap-guangzhou"
  root_password     = "Root123$"
  charset           = "LATIN1"
  project_id        = 0
  memory            = 4
  storage           = 100

  backup_plan {
    min_backup_start_time        = "00:10:11"
    max_backup_start_time        = "01:10:11"
    base_backup_retention_period = 7
    backup_period                = ["tuesday", "wednesday"]
  }

  tags = {
    CreateBy = "Terraform"
  }
}
```

### Upgrade kernel version

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-6"
}

resource "tencentcloud_postgresql_instance" "example" {
  name                 = "tf_postsql_instance_update_kernel"
  availability_zone    = var.availability_zone
  charge_type          = "POSTPAID_BY_HOUR"
  vpc_id               = "vpc-86v957zb"
  subnet_id            = "subnet-enm92y0m"
  engine_version       = "13.3"
  db_kernel_version    = "v13.3_r1.4" # eg:from v13.3_r1.1 to v13.3_r1.4
  db_major_version     = "13"
  root_password        = "Root123$"
  charset              = "LATIN1"
  project_id           = 0
  public_access_switch = false
  security_groups      = ["sg-cm7fbbf3"]
  memory               = 4
  storage              = 250

  backup_plan {
    min_backup_start_time        = "01:10:11"
    max_backup_start_time        = "02:10:11"
    base_backup_retention_period = 5
    backup_period                = ["monday", "thursday", "sunday"]
  }

  tags = {
    CreateBy = "Terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Required, String) Availability zone. NOTE: This field could not be modified, please use `db_node_set` instead of modification. The changes on this field will be suppressed when using the `db_node_set`.
* `memory` - (Required, Int) Memory size(in GB). Allowed value must be larger than `memory` that data source `tencentcloud_postgresql_specinfos` provides.
* `name` - (Required, String) Name of the postgresql instance.
* `root_password` - (Required, String) Password of root account. This parameter can be specified when you purchase master instances, but it should be ignored when you purchase read-only instances or disaster recovery instances.
* `storage` - (Required, Int) Volume size(in GB). Allowed value must be a multiple of 10. The storage must be set with the limit of `storage_min` and `storage_max` which data source `tencentcloud_postgresql_specinfos` provides.
* `subnet_id` - (Required, String) ID of subnet.
* `vpc_id` - (Required, String) ID of VPC.
* `auto_renew_flag` - (Optional, Int) Auto renew flag, `1` for enabled. NOTES: Only support prepaid instance.
* `auto_voucher` - (Optional, Int) Whether to use voucher, `1` for enabled.
* `backup_plan` - (Optional, List) Specify DB backup plan.
* `charge_type` - (Optional, String) Pay type of the postgresql instance. Values `POSTPAID_BY_HOUR` (Default), `PREPAID`. It only support to update the type from `POSTPAID_BY_HOUR` to `PREPAID`.
* `charset` - (Optional, String, ForceNew) Charset of the root account. Valid values are `UTF8`,`LATIN1`.
* `cpu` - (Optional, Int) Number of CPU cores. Allowed value must be equal `cpu` that data source `tencentcloud_postgresql_specinfos` provides.
* `db_kernel_version` - (Optional, String) PostgreSQL kernel version number. If it is specified, an instance running kernel DBKernelVersion will be created. It supports updating the minor kernel version immediately.
* `db_major_version` - (Optional, String) PostgreSQL major version number. Valid values: 10, 11, 12, 13, 14, 15, 16. If it is specified, an instance running the latest kernel of PostgreSQL DBMajorVersion will be created.
* `db_major_vesion` - (Optional, String, **Deprecated**) `db_major_vesion` will be deprecated, use `db_major_version` instead. PostgreSQL major version number. Valid values: 10, 11, 12, 13, 14, 15, 16. If it is specified, an instance running the latest kernel of PostgreSQL DBMajorVersion will be created.
* `db_node_set` - (Optional, Set) Specify instance node info for disaster migration.
* `delete_protection` - (Optional, Bool) Whether to enable instance deletion protection. Default: false.
* `engine_version` - (Optional, String) Version of the postgresql database engine. Valid values: `10.4`, `10.17`, `10.23`, `11.8`, `11.12`, `11.22`, `12.4`, `12.7`, `12.18`, `13.3`, `14.2`, `14.11`, `15.1`, `16.0`.
* `kms_cluster_id` - (Optional, String) Specify the cluster served by KMS. If KMSClusterId is blank, use the KMS of the default cluster. If you choose to specify a KMS cluster, you need to pass in KMSClusterId.
* `kms_key_id` - (Optional, String) KeyId of the custom key.
* `kms_region` - (Optional, String) Region of the custom key.
* `max_standby_archive_delay` - (Optional, Int) max_standby_archive_delay applies when WAL data is being read from WAL archive (and is therefore not current). Units are milliseconds if not specified.
* `max_standby_streaming_delay` - (Optional, Int) max_standby_streaming_delay applies when WAL data is being received via streaming replication. Units are milliseconds if not specified.
* `need_support_tde` - (Optional, Int) Whether to support data transparent encryption, 1: yes, 0: no (default).
* `period` - (Optional, Int) Specify Prepaid period in month. Default `1`. Values: `1`, `2`, `3`, `4`, `5`, `6`, `7`, `8`, `9`, `10`, `11`, `12`, `24`, `36`. This field is valid only when creating a `PREPAID` type instance, or updating the charge type from `POSTPAID_BY_HOUR` to `PREPAID`.
* `project_id` - (Optional, Int) Project id, default value is `0`.
* `public_access_switch` - (Optional, Bool) Indicates whether to enable the access to an instance from public network or not.
* `root_user` - (Optional, String) Instance root account name. This parameter is optional, Default value is `root`.
* `security_groups` - (Optional, Set: [`String`]) ID of security group. If both vpc_id and subnet_id are not set, this argument should not be set either.
* `tags` - (Optional, Map) The available tags within this postgresql.
* `voucher_ids` - (Optional, List: [`String`]) Specify Voucher Ids if `auto_voucher` was `1`, only support using 1 vouchers for now.
* `wait_switch` - (Optional, Int) Switch time after instance configurations are modified. `0`: Switch immediately; `2`: Switch during maintenance time window. Default: `0`. Note: This only takes effect when updating the `memory`, `storage`, `cpu`, `db_node_set`, `db_kernel_version` fields.

The `backup_plan` object supports the following:

* `backup_period` - (Optional, List) List of backup period per week, available values: `monday`, `tuesday`, `wednesday`, `thursday`, `friday`, `saturday`, `sunday`. NOTE: At least specify two days.
* `base_backup_retention_period` - (Optional, Int) Specify days of the retention.
* `max_backup_start_time` - (Optional, String) Specify latest backup start time, format `hh:mm:ss`.
* `min_backup_start_time` - (Optional, String) Specify earliest backup start time, format `hh:mm:ss`.
* `monthly_backup_period` - (Optional, List) If it is in monthly dimension, the format is numeric characters, such as ["1","2"].
* `monthly_backup_retention_period` - (Optional, Int) Specify days of the retention.

The `db_node_set` object supports the following:

* `zone` - (Required, String) Indicates the node available zone.
* `dedicated_cluster_id` - (Optional, String) Dedicated cluster ID.
* `role` - (Optional, String) Indicates node type, available values:`Primary`, `Standby`. Default: `Standby`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Create time of the postgresql instance.
* `private_access_ip` - IP for private access.
* `private_access_port` - Port for private access.
* `public_access_host` - Host for public access.
* `public_access_port` - Port for public access.
* `uid` - Uid of the postgresql instance.


## Import

postgresql instance can be imported using the id, e.g.

```
$ terraform import tencentcloud_postgresql_instance.example postgres-cda1iex1
```

