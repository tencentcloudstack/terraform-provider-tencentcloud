---
subcategory: "TencentDB for PostgreSQL(PostgreSQL)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_postgresql_base_backups"
sidebar_current: "docs-tencentcloud-datasource-postgresql_base_backups"
description: |-
  Use this data source to query detailed information of postgresql base_backups
---

# tencentcloud_postgresql_base_backups

Use this data source to query detailed information of postgresql base_backups

## Example Usage

```hcl
data "tencentcloud_postgresql_base_backups" "base_backups" {
  min_finish_time = "%s"
  max_finish_time = "%s"

  order_by      = "StartTime"
  order_by_type = "asc"
}

data "tencentcloud_postgresql_base_backups" "base_backups" {
  filters {
    name   = "db-instance-id"
    values = [local.pgsql_id]
  }

  order_by      = "Size"
  order_by_type = "asc"
}
```

## Argument Reference

The following arguments are supported:

* `filters` - (Optional, List) Filter instances using one or more criteria. Valid filter names: `db-instance-id`: Filter by instance ID (in string format). `db-instance-name`: Filter by instance name (in string format). `db-instance-ip`: Filter by instance VPC IP (in string format). `base-backup-id`: Filter by base backup ID (in string format).
* `max_finish_time` - (Optional, String) Maximum end time of a backup in the format of `2018-01-01 00:00:00`. It is the current time by default.
* `min_finish_time` - (Optional, String) Minimum end time of a backup in the format of `2018-01-01 00:00:00`. It is 7 days ago by default.
* `order_by_type` - (Optional, String) Sorting order. Valid values: `asc` (ascending), `desc` (descending).
* `order_by` - (Optional, String) Sorting field. Valid values: `StartTime`, `FinishTime`, `Size`.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `name` - (Optional, String) Filter name.
* `values` - (Optional, Set) One or more filter values.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `base_backup_set` - List of full backup details.
  * `backup_method` - Backup method, including physical and logical.
  * `backup_mode` - Backup mode, including automatic and manual.
  * `db_instance_id` - Instance ID.
  * `expire_time` - Backup expiration time.
  * `finish_time` - Backup end time.
  * `id` - Unique ID of a backup file.
  * `name` - Backup file name.
  * `size` - Backup set size in bytes.
  * `start_time` - Backup start time.
  * `state` - Backup task status.


