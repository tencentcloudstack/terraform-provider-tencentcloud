---
subcategory: "TencentDB for MySQL(cdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_binlog_backup_overview"
sidebar_current: "docs-tencentcloud-datasource-mysql_binlog_backup_overview"
description: |-
  Use this data source to query detailed information of mysql binlog_backup_overview
---

# tencentcloud_mysql_binlog_backup_overview

Use this data source to query detailed information of mysql binlog_backup_overview

## Example Usage

```hcl
data "tencentcloud_mysql_binlog_backup_overview" "binlog_backup_overview" {
  product = "mysql"
}
```

## Argument Reference

The following arguments are supported:

* `product` - (Required, String) The type of cloud database product to be queried, currently only supports `mysql`.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `binlog_archive_count` - The number of archived log backups.
* `binlog_archive_volume` - Archived log backup capacity (in bytes).
* `binlog_backup_count` - The total number of log backups, including remote log backups.
* `binlog_backup_volume` - Total log backup capacity, including off-site log backup (unit is byte).
* `binlog_standby_count` - The number of standard storage log backups.
* `binlog_standby_volume` - Standard storage log backup capacity (in bytes).
* `remote_binlog_count` - The number of remote log backups.
* `remote_binlog_volume` - Remote log backup capacity (in bytes).


