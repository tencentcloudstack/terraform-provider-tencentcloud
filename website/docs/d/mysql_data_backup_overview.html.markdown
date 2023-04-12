---
subcategory: "TencentDB for MySQL(cdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_data_backup_overview"
sidebar_current: "docs-tencentcloud-datasource-mysql_data_backup_overview"
description: |-
  Use this data source to query detailed information of mysql data_backup_overview
---

# tencentcloud_mysql_data_backup_overview

Use this data source to query detailed information of mysql data_backup_overview

## Example Usage

```hcl
data "tencentcloud_mysql_data_backup_overview" "data_backup_overview" {
  product = "mysql"
}
```

## Argument Reference

The following arguments are supported:

* `product` - (Required, String) The type of cloud database product to be queried, currently only supports `mysql`.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `auto_backup_count` - The total number of automatic backups in the current region.
* `auto_backup_volume` - The total automatic backup capacity of the current region.
* `data_backup_archive_count` - The total number of archive backups in the current region.
* `data_backup_archive_volume` - The total capacity of the current regional archive backup.
* `data_backup_count` - The total number of data backups in the current region.
* `data_backup_standby_count` - The total number of standard storage backups in the current region.
* `data_backup_standby_volume` - The total backup capacity of the current regional standard storage.
* `data_backup_volume` - Total data backup capacity of the current region (including automatic backup and manual backup, in bytes).
* `manual_backup_count` - The total number of manual backups in the current region.
* `manual_backup_volume` - The total manual backup capacity of the current region.
* `remote_backup_count` - The total number of remote backups.
* `remote_backup_volume` - The total capacity of remote backup.


