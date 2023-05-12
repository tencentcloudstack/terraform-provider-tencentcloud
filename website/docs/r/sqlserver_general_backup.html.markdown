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
resource "tencentcloud_sqlserver_general_backup" "general_backup" {
  strategy    = 0
  db_names    = ["db1", "db2"]
  instance_id = "mssql-i1z41iwd"
  backup_name = "bk_name"
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
* `backup_files` - Details of backup file list.
  * `backup_format` - Backup file format (pkg - packaged backup file, single - single library backup file).
  * `backup_name` - backup name.
  * `backup_way` - Backup mode, 0-scheduled backup; 1-manual temporary backup; 2-regular backup.
  * `cross_backup_addr` - Destination domain download link for cross-region backup.
  * `cross_backup_status` - Target region and backup status of cross-region backup.
  * `dbs` - The name of the library for backing up files.
  * `end_time` - end time.
  * `external_addr` - External network download address, this value is not returned for single-database backup files; the download address of single-database backup files is obtained through the DescribeBackupFiles interface.
  * `file_name` - backup file name.
  * `group_id` - Aggregate Id, this value is not returned for packaged backup files. Use this value to call the DescribeBackupFiles interface to obtain the detailed information of a single database backup file.
  * `id` - backup id.
  * `internal_addr` - Intranet download address, this value is not returned for a single database backup file; the download address of a single database backup file is obtained through the DescribeBackupFiles interface.
  * `region` - region.
  * `size` - file size(k).
  * `start_time` - start time.
  * `status` - Backup file status (0-creating; 1-success; 2-failure).
  * `strategy` - Backup strategy (0-instance backup; 1-multi-database backup).
* `backup_id` - Backup ID.


## Import

sqlserver general_backups can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_general_backups.general_backups general_backups_id
```

