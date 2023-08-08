---
subcategory: "SQLServer"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_full_backup_migration"
sidebar_current: "docs-tencentcloud-resource-sqlserver_full_backup_migration"
description: |-
  Provides a resource to create a sqlserver full_backup_migration
---

# tencentcloud_sqlserver_full_backup_migration

Provides a resource to create a sqlserver full_backup_migration

## Example Usage

```hcl
data "tencentcloud_availability_zones_by_product" "zones" {
  product = "sqlserver"
}

data "tencentcloud_sqlserver_backups" "example" {
  instance_id = tencentcloud_sqlserver_db.example.instance_id
  backup_name = tencentcloud_sqlserver_general_backup.example.backup_name
  start_time  = "2023-07-25 00:00:00"
  end_time    = "2023-08-04 00:00:00"
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

resource "tencentcloud_sqlserver_db" "example" {
  instance_id = tencentcloud_sqlserver_basic_instance.example.id
  name        = "tf_example_db"
  charset     = "Chinese_PRC_BIN"
  remark      = "test-remark"
}

resource "tencentcloud_sqlserver_general_backup" "example" {
  instance_id = tencentcloud_sqlserver_db.example.instance_id
  backup_name = "tf_example_backup"
  strategy    = 0
}

resource "tencentcloud_sqlserver_full_backup_migration" "example" {
  instance_id    = tencentcloud_sqlserver_basic_instance.example.id
  recovery_type  = "FULL"
  upload_type    = "COS_URL"
  migration_name = "migration_test"
  backup_files   = [data.tencentcloud_sqlserver_backups.example.list.0.internet_url]
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) ID of imported target instance.
* `migration_name` - (Required, String) Task name.
* `recovery_type` - (Required, String) Migration task restoration type. FULL: full backup restoration, FULL_LOG: full backup and transaction log restoration, FULL_DIFF: full backup and differential backup restoration.
* `upload_type` - (Required, String) Backup upload type. COS_URL: the backup is stored in users Cloud Object Storage, with URL provided. COS_UPLOAD: the backup is stored in the applications Cloud Object Storage and needs to be uploaded by the user.
* `backup_files` - (Optional, List: [`String`]) If the UploadType is COS_URL, fill in the URL here. If the UploadType is COS_UPLOAD, fill in the name of the backup file here. Only 1 backup file is supported, but a backup file can involve multiple databases.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `backup_migration_id` - Backup import task ID.


## Import

sqlserver full_backup_migration can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_full_backup_migration.full_backup_migration full_backup_migration_id
```

