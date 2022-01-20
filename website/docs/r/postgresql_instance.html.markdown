---
subcategory: "PostgreSQL"
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

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Required, ForceNew) Availability zone. NOTE: If value modified but included in `db_node_set`, the diff will be suppressed.
* `memory` - (Required) Memory size(in GB). Allowed value must be larger than `memory` that data source `tencentcloud_postgresql_specinfos` provides.
* `name` - (Required) Name of the postgresql instance.
* `root_password` - (Required) Password of root account. This parameter can be specified when you purchase master instances, but it should be ignored when you purchase read-only instances or disaster recovery instances.
* `storage` - (Required) Volume size(in GB). Allowed value must be a multiple of 10. The storage must be set with the limit of `storage_min` and `storage_max` which data source `tencentcloud_postgresql_specinfos` provides.
* `charge_type` - (Optional, ForceNew) Pay type of the postgresql instance. For now, only `POSTPAID_BY_HOUR` is valid.
* `charset` - (Optional, ForceNew) Charset of the root account. Valid values are `UTF8`,`LATIN1`.
* `db_node_set` - (Optional) Specify instance node info for disaster migration.
* `engine_version` - (Optional, ForceNew) Version of the postgresql database engine. Valid values: `10.4`, `11.8`, `12.4`.
* `max_standby_archive_delay` - (Optional) max_standby_archive_delay applies when WAL data is being read from WAL archive (and is therefore not current). Units are milliseconds if not specified.
* `max_standby_streaming_delay` - (Optional) max_standby_streaming_delay applies when WAL data is being received via streaming replication. Units are milliseconds if not specified.
* `project_id` - (Optional) Project id, default value is `0`.
* `public_access_switch` - (Optional) Indicates whether to enable the access to an instance from public network or not.
* `root_user` - (Optional, ForceNew) Instance root account name. This parameter is optional, Default value is `root`.
* `security_groups` - (Optional) ID of security group. If both vpc_id and subnet_id are not set, this argument should not be set either.
* `subnet_id` - (Optional, ForceNew) ID of subnet.
* `tags` - (Optional) The available tags within this postgresql.
* `vpc_id` - (Optional, ForceNew) ID of VPC.

The `db_node_set` object supports the following:

* `zone` - (Required) Indicates the node available zone.
* `role` - (Optional) Indicates node type, available values:`Primary`, `Standby`. Default: `Standby`.

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

