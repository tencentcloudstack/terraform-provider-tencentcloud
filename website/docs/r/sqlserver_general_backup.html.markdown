---
subcategory: "SQLServer"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_general_backup"
sidebar_current: "docs-tencentcloud-resource-sqlserver_general_backup"
description: |-
  Provides a resource to create a sqlserver general_backup
---

# tencentcloud_sqlserver_general_backup

Provides a resource to create a sqlserver general_backup

## Example Usage

```hcl
data "tencentcloud_availability_zones_by_product" "zones" {
  product = "sqlserver"
}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  name              = "subnet-example"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.0.0/16"
  is_multicast      = false
}

resource "tencentcloud_security_group" "security_group" {
  name        = "sg-example"
  description = "desc."
}

resource "tencentcloud_sqlserver_basic_instance" "example" {
  name                   = "tf-example"
  availability_zone      = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  charge_type            = "POSTPAID_BY_HOUR"
  vpc_id                 = tencentcloud_vpc.vpc.id
  subnet_id              = tencentcloud_subnet.subnet.id
  project_id             = 0
  memory                 = 4
  storage                = 100
  cpu                    = 2
  machine_type           = "CLOUD_PREMIUM"
  maintenance_week_set   = [1, 2, 3]
  maintenance_start_time = "09:00"
  maintenance_time_span  = 3
  security_groups        = [tencentcloud_security_group.security_group.id]

  tags = {
    "test" = "test"
  }
}

resource "tencentcloud_sqlserver_general_backup" "example" {
  instance_id = tencentcloud_sqlserver_basic_instance.example.id
  backup_name = "tf_example_backup"
  strategy    = 0
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Instance ID in the format of mssql-i1z41iwd.
* `backup_name` - (Optional, String) Backup name. If this parameter is left empty, a backup name in the format of [Instance ID]_[Backup start timestamp] will be automatically generated.
* `db_names` - (Optional, Set: [`String`]) List of names of databases to be backed up (required only for multi-database backup).
* `strategy` - (Optional, Int) Backup policy (0: instance backup, 1: multi-database backup).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `flow_id` - flow id.


## Import

sqlserver general_backups can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_general_backups.general_backups general_backups_id
```

