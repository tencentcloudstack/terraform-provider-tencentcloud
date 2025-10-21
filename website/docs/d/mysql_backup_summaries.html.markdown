---
subcategory: "TencentDB for MySQL(cdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_backup_summaries"
sidebar_current: "docs-tencentcloud-datasource-mysql_backup_summaries"
description: |-
  Use this data source to query detailed information of mysql backup_summaries
---

# tencentcloud_mysql_backup_summaries

Use this data source to query detailed information of mysql backup_summaries

## Example Usage

```hcl
data "tencentcloud_mysql_backup_summaries" "backup_summaries" {
  product         = "mysql"
  order_by        = "BackupVolume"
  order_direction = "ASC"
}
```

## Argument Reference

The following arguments are supported:

* `product` - (Required, String) The type of cloud database product to be queried, currently only supports `mysql`.
* `order_by` - (Optional, String) Specify to sort by a certain item, the optional values include: BackupVolume: backup volume, DataBackupVolume: data backup volume, BinlogBackupVolume: log backup volume, AutoBackupVolume: automatic backup volume, ManualBackupVolume: manual backup volume. By default, they are sorted by BackupVolume.
* `order_direction` - (Optional, String) Specify the sorting direction, optional values include: ASC: forward order, DESC: reverse order. The default is ASC.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `items` - Instance backup statistics entries.
  * `auto_backup_count` - The number of automatic data backups for this instance.
  * `auto_backup_volume` - The automatic data backup capacity of this instance.
  * `backup_volume` - The total backup (including data backup and log backup) of the instance occupies capacity.
  * `binlog_backup_count` - The number of log backups for this instance.
  * `binlog_backup_volume` - The capacity of the instance log backup.
  * `data_backup_count` - The total number of data backups (including automatic backups and manual backups) of the instance.
  * `data_backup_volume` - The total data backup capacity of this instance.
  * `instance_id` - Instance ID.
  * `manual_backup_count` - The number of manual data backups for this instance.
  * `manual_backup_volume` - The capacity of manual data backup for this instance.


