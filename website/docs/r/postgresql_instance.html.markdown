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

## Example Usage

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-1"
}

# create vpc
resource "tencentcloud_vpc" "vpc" {
  name       = "guagua_vpc_instance_test"
  cidr_block = "10.0.0.0/16"
}

# create vpc subnet
resource "tencentcloud_subnet" "subnet" {
  availability_zone = var.availability_zone
  name              = "guagua_vpc_subnet_test"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.20.0/28"
  is_multicast      = false
}

# create postgresql
resource "tencentcloud_postgresql_instance" "foo" {
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
  memory            = 2
  storage           = 10

  tags = {
    test = "tf"
  }
}
```

Create a multi available zone bucket

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-6"
}

variable "standby_availability_zone" {
  default = "ap-guangzhou-7"
}

# create vpc
resource "tencentcloud_vpc" "vpc" {
  name       = "guagua_vpc_instance_test"
  cidr_block = "10.0.0.0/16"
}

# create vpc subnet
resource "tencentcloud_subnet" "subnet" {
  availability_zone = var.availability_zone
  name              = "guagua_vpc_subnet_test"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.20.0/28"
  is_multicast      = false
}

# create postgresql
resource "tencentcloud_postgresql_instance" "foo" {
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
  memory            = 2
  storage           = 10

  db_node_set {
    role = "Primary"
    zone = var.availability_zone
  }
  db_node_set {
    zone = var.standby_availability_zone
  }

  tags = {
    test = "tf"
  }
}
```

create pgsql with kms key

```hcl
resource "tencentcloud_postgresql_instance" "pg" {
  name              = "tf_postsql_instance"
  availability_zone = "ap-guangzhou-6"
  charge_type       = "POSTPAID_BY_HOUR"
  vpc_id            = "vpc-86v957zb"
  subnet_id         = "subnet-enm92y0m"
  engine_version    = "11.12"
  #  db_major_vesion   = "11"
  db_kernel_version = "v11.12_r1.3"
  need_support_tde  = 1
  kms_key_id        = "788c606a-c7b7-11ec-82d1-5254001e5c4e"
  kms_region        = "ap-guangzhou"
  root_password     = "xxxxxxxxxx"
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
    tf = "test"
  }
}
```

upgrade kernel version

```hcl
resource "tencentcloud_postgresql_instance" "test" {
  name                 = "tf_postsql_instance_update"
  availability_zone    = data.tencentcloud_availability_zones_by_product.zone.zones[5].name
  charge_type          = "POSTPAID_BY_HOUR"
  vpc_id               = local.vpc_id
  subnet_id            = local.subnet_id
  engine_version       = "13.3"
  root_password        = "*"
  charset              = "LATIN1"
  project_id           = 0
  public_access_switch = false
  security_groups      = [local.sg_id]
  memory               = 4
  storage              = 250
  backup_plan {
    min_backup_start_time        = "01:10:11"
    max_backup_start_time        = "02:10:11"
    base_backup_retention_period = 5
    backup_period                = ["monday", "thursday", "sunday"]
  }

  db_kernel_version = "v13.3_r1.4" # eg:from v13.3_r1.1 to v13.3_r1.4

  tags = {
    tf = "teest"
  }
}
```

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Required, String, ForceNew) Availability zone. NOTE: If value modified but included in `db_node_set`, the diff will be suppressed.
* `memory` - (Required, Int) Memory size(in GB). Allowed value must be larger than `memory` that data source `tencentcloud_postgresql_specinfos` provides.
* `name` - (Required, String) Name of the postgresql instance.
* `root_password` - (Required, String) Password of root account. This parameter can be specified when you purchase master instances, but it should be ignored when you purchase read-only instances or disaster recovery instances.
* `storage` - (Required, Int) Volume size(in GB). Allowed value must be a multiple of 10. The storage must be set with the limit of `storage_min` and `storage_max` which data source `tencentcloud_postgresql_specinfos` provides.
* `subnet_id` - (Required, String) ID of subnet.
* `vpc_id` - (Required, String) ID of VPC.
* `auto_renew_flag` - (Optional, Int) Auto renew flag, `1` for enabled. NOTES: Only support prepaid instance.
* `auto_voucher` - (Optional, Int) Whether to use voucher, `1` for enabled.
* `backup_plan` - (Optional, List) Specify DB backup plan.
* `charge_type` - (Optional, String) Pay type of the postgresql instance. Values `POSTPAID_BY_HOUR` (Default), `PREPAID`. It support to update the type from `POSTPAID_BY_HOUR` to `PREPAID`.
* `charset` - (Optional, String, ForceNew) Charset of the root account. Valid values are `UTF8`,`LATIN1`.
* `db_kernel_version` - (Optional, String) PostgreSQL kernel version number. If it is specified, an instance running kernel DBKernelVersion will be created. It supports updating the minor kernel version immediately.
* `db_major_version` - (Optional, String) PostgreSQL major version number. Valid values: 10, 11, 12, 13. If it is specified, an instance running the latest kernel of PostgreSQL DBMajorVersion will be created.
* `db_major_vesion` - (Optional, String, **Deprecated**) `db_major_vesion` will be deprecated, use `db_major_version` instead. PostgreSQL major version number. Valid values: 10, 11, 12, 13. If it is specified, an instance running the latest kernel of PostgreSQL DBMajorVersion will be created.
* `db_node_set` - (Optional, Set) Specify instance node info for disaster migration.
* `engine_version` - (Optional, String, ForceNew) Version of the postgresql database engine. Valid values: `10.4`, `11.8`, `12.4`.
* `kms_key_id` - (Optional, String) KeyId of the custom key.
* `kms_region` - (Optional, String) Region of the custom key.
* `max_standby_archive_delay` - (Optional, Int) max_standby_archive_delay applies when WAL data is being read from WAL archive (and is therefore not current). Units are milliseconds if not specified.
* `max_standby_streaming_delay` - (Optional, Int) max_standby_streaming_delay applies when WAL data is being received via streaming replication. Units are milliseconds if not specified.
* `need_support_tde` - (Optional, Int) Whether to support data transparent encryption, 1: yes, 0: no (default).
* `period` - (Optional, Int) Specify Prepaid period in month. Default `1`. Values: `1`, `2`, `3`, `4`, `5`, `6`, `7`, `8`, `9`, `10`, `11`, `12`, `24`, `36`.
* `project_id` - (Optional, Int) Project id, default value is `0`.
* `public_access_switch` - (Optional, Bool) Indicates whether to enable the access to an instance from public network or not.
* `root_user` - (Optional, String, ForceNew) Instance root account name. This parameter is optional, Default value is `root`.
* `security_groups` - (Optional, Set: [`String`]) ID of security group. If both vpc_id and subnet_id are not set, this argument should not be set either.
* `tags` - (Optional, Map) The available tags within this postgresql.
* `voucher_ids` - (Optional, List: [`String`]) Specify Voucher Ids if `auto_voucher` was `1`, only support using 1 vouchers for now.

The `backup_plan` object supports the following:

* `backup_period` - (Optional, List) List of backup period per week, available values: `monday`, `tuesday`, `wednesday`, `thursday`, `friday`, `saturday`, `sunday`. NOTE: At least specify two days.
* `base_backup_retention_period` - (Optional, Int) Specify days of the retention.
* `max_backup_start_time` - (Optional, String) Specify latest backup start time, format `hh:mm:ss`.
* `min_backup_start_time` - (Optional, String) Specify earliest backup start time, format `hh:mm:ss`.

The `db_node_set` object supports the following:

* `zone` - (Required, String) Indicates the node available zone.
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
$ terraform import tencentcloud_postgresql_instance.foo postgres-cda1iex1
```

