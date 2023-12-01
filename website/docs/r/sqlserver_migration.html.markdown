---
subcategory: "SQLServer"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_migration"
sidebar_current: "docs-tencentcloud-resource-sqlserver_migration"
description: |-
  Provides a resource to create a sqlserver migration
---

# tencentcloud_sqlserver_migration

Provides a resource to create a sqlserver migration

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

resource "tencentcloud_sqlserver_basic_instance" "src_example" {
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

resource "tencentcloud_sqlserver_basic_instance" "dst_example" {
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

resource "tencentcloud_sqlserver_db" "src" {
  instance_id = tencentcloud_sqlserver_basic_instance.src_example.id
  name        = "tf_example_db_src"
  charset     = "Chinese_PRC_BIN"
  remark      = "testACC-remark"
}

resource "tencentcloud_sqlserver_db" "dst" {
  instance_id = tencentcloud_sqlserver_basic_instance.dst_example.id
  name        = "tf_example_db_dst"
  charset     = "Chinese_PRC_BIN"
  remark      = "testACC-remark"
}

resource "tencentcloud_sqlserver_account" "src" {
  instance_id = tencentcloud_sqlserver_basic_instance.src_example.id
  name        = "tf_example_src_account"
  password    = "Qwer@234"
  is_admin    = true
}

resource "tencentcloud_sqlserver_account" "dst" {
  instance_id = tencentcloud_sqlserver_basic_instance.dst_example.id
  name        = "tf_example_dst_account"
  password    = "Qwer@234"
  is_admin    = true
}

resource "tencentcloud_sqlserver_account_db_attachment" "src" {
  instance_id  = tencentcloud_sqlserver_basic_instance.src_example.id
  account_name = tencentcloud_sqlserver_account.src.name
  db_name      = tencentcloud_sqlserver_db.src.name
  privilege    = "ReadWrite"
}

resource "tencentcloud_sqlserver_account_db_attachment" "dst" {
  instance_id  = tencentcloud_sqlserver_basic_instance.dst_example.id
  account_name = tencentcloud_sqlserver_account.dst.name
  db_name      = tencentcloud_sqlserver_db.dst.name
  privilege    = "ReadWrite"
}

resource "tencentcloud_sqlserver_migration" "migration" {
  migrate_name = "tf_test_migration"
  migrate_type = 1
  source_type  = 1
  source {
    instance_id = tencentcloud_sqlserver_basic_instance.src_example.id
    user_name   = tencentcloud_sqlserver_account.src.name
    password    = tencentcloud_sqlserver_account.src.password
  }
  target {
    instance_id = tencentcloud_sqlserver_basic_instance.dst_example.id
    user_name   = tencentcloud_sqlserver_account.dst.name
    password    = tencentcloud_sqlserver_account.dst.password
  }

  migrate_db_set {
    db_name = tencentcloud_sqlserver_db.src.name
  }
}
```

## Argument Reference

The following arguments are supported:

* `migrate_name` - (Required, String) Name of the migration task.
* `migrate_type` - (Required, Int) Migration type (1 structure migration 2 data migration 3 incremental synchronization).
* `source_type` - (Required, Int) Type of migration source 1 TencentDB for SQLServer 2 Cloud server self-built SQLServer database 4 SQLServer backup and restore 5 SQLServer backup and restore (COS mode).
* `source` - (Required, List) Migration source.
* `target` - (Required, List) Migration target.
* `migrate_db_set` - (Optional, List) Migrate DB objects. Offline migration is not used (SourceType=4 or SourceType=5).
* `rename_restore` - (Optional, List) Restore and rename the database in ReNameRestoreDatabase. If it is not filled in, the restored database will be named by default and all databases will be restored. Valid if SourceType=5.

The `migrate_db_set` object supports the following:

* `db_name` - (Optional, String) Name of the migration database.

The `rename_restore` object supports the following:

* `new_name` - (Optional, String) When the new name of the library is used for offline migration, if it is not filled in, it will be named according to OldName. OldName and NewName cannot be filled in at the same time. OldName and NewName must be filled in and cannot be duplicate when used for cloning database.
* `old_name` - (Optional, String) The name of the library. If oldName does not exist, a failure is returned.It can be left blank when used for offline migration tasks.

The `source` object supports the following:

* `cvm_id` - (Optional, String) ID of the migration source Cvm, used when MigrateType=2 (cloud server self-built SQL Server database).
* `instance_id` - (Optional, String) The ID of the migration source instance, which is used when MigrateType=1 (TencentDB for SQLServers). The format is mssql-si2823jyl.
* `ip` - (Optional, String) Migrate the intranet IP of the self-built database of the source Cvm, and use it when MigrateType=2 (self-built SQL Server database of the cloud server).
* `password` - (Optional, String) Password, MigrateType=1 or MigrateType=2.
* `port` - (Optional, Int) The port number of the self-built database of the migration source Cvm, which is used when MigrateType=2 (self-built SQL Server database of the cloud server).
* `subnet_id` - (Optional, String) The subnet ID under the Vpc of the source Cvm is used when MigrateType=2 (ECS self-built SQL Server database). The format is as follows subnet-h9extioi.
* `url_password` - (Optional, String) The source backup password for offline migration, MigrateType=4 or MigrateType=5.
* `url` - (Optional, Set) The source backup address for offline migration. MigrateType=4 or MigrateType=5.
* `user_name` - (Optional, String) User name, MigrateType=1 or MigrateType=2.
* `vpc_id` - (Optional, String) The Vpc network ID of the migration source Cvm is used when MigrateType=2 (cloud server self-built SQL Server database). The format is as follows vpc-6ys9ont9.

The `target` object supports the following:

* `instance_id` - (Optional, String) The ID of the migration target instance, in the format mssql-si2823jyl.
* `password` - (Optional, String) Password of the migration target instance.
* `user_name` - (Optional, String) User name of the migration target instance.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

sqlserver migration can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_migration.migration migration_id
```

