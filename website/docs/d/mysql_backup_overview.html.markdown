---
subcategory: "TencentDB for MySQL(cdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_backup_overview"
sidebar_current: "docs-tencentcloud-datasource-mysql_backup_overview"
description: |-
  Use this data source to query detailed information of mysql backup_overview
---

# tencentcloud_mysql_backup_overview

Use this data source to query detailed information of mysql backup_overview

## Example Usage

```hcl
data "tencentcloud_mysql_backup_overview" "backup_overview" {
  product = "mysql"
}
```

## Argument Reference

The following arguments are supported:

* `product` - (Required, String) The type of cloud database product to be queried, currently only supports `mysql`.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `backup_archive_volume` - Archive backup capacity, including data backup and log backup. Note: This field may return null, indicating that no valid value can be obtained.
* `backup_count` - The total number of user backups in the current region (including data backups and log backups).
* `backup_standby_volume` - Standard storage backup capacity, including data backup and log backup. Note: This field may return null, indicating that no valid value can be obtained.
* `backup_volume` - The total backup capacity of the user in the current region.
* `billing_volume` - The billable capacity of the user&amp;#39;s backup in the current region, that is, the part that exceeds the gifted capacity.
* `free_volume` - The free backup capacity obtained by the user in the current region.
* `remote_backup_volume` - The total capacity of off-site backup of the user in the current region. Note: This field may return null, indicating that no valid value can be obtained.


