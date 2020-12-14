---
subcategory: "SQLServer"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_backups"
sidebar_current: "docs-tencentcloud-datasource-sqlserver_backups"
description: |-
  Use this data source to query the list of SQL Server backups.
---

# tencentcloud_sqlserver_backups

Use this data source to query the list of SQL Server backups.

## Example Usage

```hcl
data "tencentcloud_sqlserver_backups" "foo" {
  instance_id = "mssql-3cdq7kx5"
  start_time  = "2020-06-17 00:00:00"
  end_time    = "2020-06-22 00:00:00"
}
```

## Argument Reference

The following arguments are supported:

* `end_time` - (Required) End time of the instance list, like yyyy-MM-dd HH:mm:ss.
* `instance_id` - (Required) Instance ID.
* `start_time` - (Required) Start time of the instance list, like yyyy-MM-dd HH:mm:ss.
* `result_output_file` - (Optional) Used to store results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - A list of SQL Server backup. Each element contains the following attributes:
  * `db_list` - Database name list of the backup.
  * `end_time` - End time of the backup.
  * `file_name` - File name of the backup.
  * `id` - ID of the backup.
  * `instance_id` - Instance ID.
  * `internet_url` - URL for downloads externally.
  * `intranet_url` - URL for downloads internally.
  * `size` - The size of backup file. Unit is KB.
  * `start_time` - Start time of the backup.
  * `status` - Status of the backup. `1` for creating, `2` for successfully created, 3 for failed.
  * `strategy` - Strategy of the backup. `0` for instance backup, `1` for multi-databases backup.
  * `trigger_model` - The way to trigger backup. `0` for timed trigger, `1` for manual trigger.


