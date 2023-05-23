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
resource "tencentcloud_sqlserver_full_backup_migration" "my_migration" {
  instance_id    = "mssql-qelbzgwf"
  recovery_type  = "FULL"
  upload_type    = "COS_URL"
  migration_name = "migration_test"
  backup_files   = []
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



## Import

sqlserver full_backup_migration can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_full_backup_migration.full_backup_migration full_backup_migration_id
```

