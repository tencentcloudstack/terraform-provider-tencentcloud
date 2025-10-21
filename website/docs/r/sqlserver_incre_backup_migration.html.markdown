---
subcategory: "SQLServer"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_incre_backup_migration"
sidebar_current: "docs-tencentcloud-resource-sqlserver_incre_backup_migration"
description: |-
  Provides a resource to create a sqlserver incre_backup_migration
---

# tencentcloud_sqlserver_incre_backup_migration

Provides a resource to create a sqlserver incre_backup_migration

## Example Usage

```hcl
resource "tencentcloud_sqlserver_incre_backup_migration" "example" {
  instance_id         = "mssql-4gmc5805"
  backup_migration_id = "mssql-backup-migration-9tj0sxnz"
  backup_files        = []
  is_recovery         = "YES"
}
```

## Argument Reference

The following arguments are supported:

* `backup_migration_id` - (Required, String) Backup import task ID, which is returned through the API CreateBackupMigration.
* `instance_id` - (Required, String) ID of imported target instance.
* `backup_files` - (Optional, List: [`String`]) Incremental backup file. If the UploadType of a full backup file is COS_URL, fill in URL here. If the UploadType is COS_UPLOAD, fill in the name of the backup file here. Only 1 backup file is supported, but a backup file can involve multiple databases.
* `is_recovery` - (Optional, String) Whether restoration is required. No: not required. Yes: required. Not required by default.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `incremental_migration_id` - Incremental import task ID.


## Import

sqlserver incre_backup_migration can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_incre_backup_migration.incre_backup_migration incre_backup_migration_id
```

