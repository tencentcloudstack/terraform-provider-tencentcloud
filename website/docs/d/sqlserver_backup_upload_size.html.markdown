---
subcategory: "SQLServer"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_backup_upload_size"
sidebar_current: "docs-tencentcloud-datasource-sqlserver_backup_upload_size"
description: |-
  Use this data source to query detailed information of sqlserver datasource_backup_upload_size
---

# tencentcloud_sqlserver_backup_upload_size

Use this data source to query detailed information of sqlserver datasource_backup_upload_size

## Example Usage

```hcl
data "tencentcloud_sqlserver_backup_upload_size" "backup_upload_size" {
  instance_id         = "mssql-4gmc5805"
  backup_migration_id = "mssql-backup-migration-9tj0sxnz"
}
```

## Argument Reference

The following arguments are supported:

* `backup_migration_id` - (Required, String) Backup import task ID, which is returned through the API CreateBackupMigration.
* `instance_id` - (Required, String) ID of imported target instance.
* `incremental_migration_id` - (Optional, String) Incremental import task ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `cos_upload_backup_file_set` - Information of uploaded backups.
  * `file_name` - Backup name.
  * `size` - Backup size.


