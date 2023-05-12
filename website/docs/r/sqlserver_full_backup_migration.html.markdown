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
resource "tencentcloud_sqlserver_full_backup_migration" "full_backup_migration" {
  instance_id    = "mssql-i1z41iwd"
  recovery_type  = "FULL"
  upload_type    = "COS_URL"
  migration_name = "test_migration"
  backup_files   = []
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) ID of imported target instance.
* `migration_name` - (Required, String) Task name.
* `recovery_type` - (Required, String) Migration task restoration type. FULL: full backup restoration, FULL_LOG: full backup and transaction log restoration, FULL_DIFF: full backup and differential backup restoration.
* `upload_type` - (Required, String) Backup upload type. COS_URL: the backup is stored in users Cloud Object Storage, with URL provided. COS_UPLOAD: the backup is stored in the application Cloud Object Storage and needs to be uploaded by the user.
* `backup_files` - (Optional, List: [`String`]) If the UploadType is COS_URL, fill in the URL here. If the UploadType is COS_UPLOAD, fill in the name of the backup file here. Only 1 backup file is supported, but a backup file can involve multiple databases.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `backup_migration_id` - migration id.
* `backup_migration_set` - backup migration set.
  * `app_id` - app id.
  * `backup_files` - backup files list.
  * `create_time` - create time.
  * `end_time` - end time.
  * `instance_id` - instance id.
  * `is_recovery` - Whether it is the final recovery, the field of the full import task is empty.
  * `message` - msg.
  * `migration_id` - migration id.
  * `migration_name` - migration name.
  * `recovery_type` - recovery type.
  * `region` - region.
  * `start_time` - start time.
  * `status` - Migration task status, 2-created, 7-full import, 8-waiting for increment, 9-import successful, 10-import failed, 12-incremental import.
  * `upload_type` - Backup user upload type, COS_URL - the backup is placed on the user's object storage, and the URL is provided. COS_UPLOAD-The backup is placed on the object storage of the business, and the user uploads it.


## Import

sqlserver full_backup_migration can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_full_backup_migration.full_backup_migration full_backup_migration_id
```

