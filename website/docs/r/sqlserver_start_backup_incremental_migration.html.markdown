---
subcategory: "SQLServer"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_start_backup_incremental_migration"
sidebar_current: "docs-tencentcloud-resource-sqlserver_start_backup_incremental_migration"
description: |-
  Provides a resource to create a sqlserver start_backup_incremental_migration
---

# tencentcloud_sqlserver_start_backup_incremental_migration

Provides a resource to create a sqlserver start_backup_incremental_migration

## Example Usage

```hcl
resource "tencentcloud_sqlserver_start_backup_incremental_migration" "start_backup_incremental_migration" {
  instance_id              = "mssql-i1z41iwd"
  backup_migration_id      = "mssql-backup-migration-cg0ffgqt"
  incremental_migration_id = "mssql-incremental-migration-kp7bgv8p"
}
```

## Argument Reference

The following arguments are supported:

* `backup_migration_id` - (Required, String, ForceNew) Backup import task ID, returned by the CreateBackupMigration interface.
* `incremental_migration_id` - (Required, String, ForceNew) Incremental backup import task ID.
* `instance_id` - (Required, String, ForceNew) ID of imported target instance.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



