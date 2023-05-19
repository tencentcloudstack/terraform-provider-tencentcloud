---
subcategory: "SQLServer"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_backup_by_flow_id"
sidebar_current: "docs-tencentcloud-datasource-sqlserver_backup_by_flow_id"
description: |-
  Use this data source to query detailed information of sqlserver datasource_backup_by_flow_id
---

# tencentcloud_sqlserver_backup_by_flow_id

Use this data source to query detailed information of sqlserver datasource_backup_by_flow_id

## Example Usage

```hcl
resource "tencentcloud_sqlserver_general_backup" "general_backup" {
  strategy    = 0
  instance_id = "mssql-qelbzgwf"
  backup_name = "create_sqlserver_backup_name"
}

data "tencentcloud_sqlserver_backup_by_flow_id" "backup_by_flow_id" {
  instance_id = tencentcloud_sqlserver_general_backup.general_backup.instance_id
  flow_id     = tencentcloud_sqlserver_general_backup.general_backup.flow_id
}
```

## Argument Reference

The following arguments are supported:

* `flow_id` - (Required, String) Create a backup process ID, which can be obtained through the [CreateBackup](https://cloud.tencent.com/document/product/238/19946) interface.
* `instance_id` - (Required, String) Instance ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `backup_name` - Backup task name, customizable.
* `backup_way` - Backup method, 0-scheduled backup; 1-manual temporary backup; instance status is 0-creating, this field is the default value 0, meaningless.
* `dbs` - For the DB list, only the library name contained in the first record is returned for a single-database backup file; for a single-database backup file, the library names of all records need to be obtained through the DescribeBackupFiles interface.
* `end_time` - backup end time.
* `external_addr` - External network download address, for a single database backup file, only the external network download address of the first record is returned; single database backup files need to obtain the download addresses of all records through the DescribeBackupFiles interface.
* `file_name` - File name. For a single-database backup file, only the file name of the first record is returned; for a single-database backup file, the file names of all records need to be obtained through the DescribeBackupFiles interface.
* `group_id` - Aggregate Id, this value is not returned for packaged backup files. Use this value to call the DescribeBackupFiles interface to obtain the detailed information of a single database backup file.
* `internal_addr` - Intranet download address, for a single database backup file, only the intranet download address of the first record is returned; single database backup files need to obtain the download addresses of all records through the DescribeBackupFiles interface.
* `start_time` - backup start time.
* `status` - Backup file status, 0-creating; 1-success; 2-failure.
* `strategy` - Backup strategy, 0-instance backup; 1-multi-database backup; when the instance status is 0-creating, this field is the default value 0, meaningless.


