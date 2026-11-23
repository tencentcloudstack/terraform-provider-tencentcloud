---
subcategory: "TDSQL Boundless(tdmysql)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tdmysql_dbs_backup"
sidebar_current: "docs-tencentcloud-resource-tdmysql_dbs_backup"
description: |-
  # Resource: tencentcloud_tdmysql_dbs_backup
---

# tencentcloud_tdmysql_dbs_backup

# Resource: tencentcloud_tdmysql_dbs_backup

Provides a resource to create a manual backup set for a TDSQL-C (TDSQL Boundless, tdmysql) instance.

## Example Usage

```hcl
resource "tencentcloud_tdmysql_dbs_backup" "example" {
  instance_id   = "tdsqlc-abc123"
  backup_method = "physical"
  backup_type   = "full"
  backup_name   = "tf-manual-backup"
}
```

## Argument Reference

The following arguments are supported:

* `backup_method` - (Required, String) Backup method. It must be consistent with the `backupMethod` returned by DescribeDBSBackupPolicy. Valid values: `physical`, `snapshot`. Changing this value forces recreation.
* `backup_type` - (Required, String) Backup type. Currently only `full` is supported. Changing this value forces recreation.
* `instance_id` - (Required, String) ID of the TDSQL-C instance to back up. Changing this value forces recreation.
* `backup_name` - (Optional, String) Backup comment name. It can be updated in place through ModifyDBSBackupSetComment.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

  * `backup_duration` - Backup set duration in seconds.
  * `backup_method` - Backup method.
  * `backup_name` - Backup comment name.
  * `backup_progress` - Backup progress percentage.
  * `backup_set_id` - Backup set ID.
  * `backup_status` - Backup set status. Valid values: `pending`, `running`, `success`, `failed`.
  * `backup_type` - Backup set type.
  * `db_version` - Instance version.
  * `end_time` - Backup end time.
  * `end_trx_time` - Transaction commit end time.
  * `error_message` - Error message.
  * `expired_time` - Backup expiration time.
  * `file_name` - Backup file name.
  * `file_size` - Backup file size in bytes.
  * `instance_id` - Instance ID.
  * `manual_backup` - Backup trigger. Valid values: `0` (auto backup), `1` (manual backup).
  * `start_time` - Backup start time.
* `backup_set_id` - ID of the backup set returned by CreateDBSBackup.
* `id` - ID of the backup set in the format `instance_id#backup_set_id`.
* `items` - Backup set list returned by DescribeDBSBackupSets.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create:
  - [CreateDBSBackup](https://www.tencentcloud.com/document/api/1292/78258)
  - [DescribeDBSBackupSets](https://www.tencentcloud.com/document/api/1292/80116)
* Read: [DescribeDBSBackupSets](https://www.tencentcloud.com/document/api/1292/80116)
* Update:
  - [DescribeDBSBackupSets](https://www.tencentcloud.com/document/api/1292/80116)
  - [ModifyDBSBackupSetComment](https://www.tencentcloud.com/document/api/1292/78275)
* Delete: [DeleteDBSBackupSets](https://www.tencentcloud.com/document/api/1292/78277)

## Import

A tdmysql backup set can be imported using the `instance_id#backup_set_id` format, e.g.

```
terraform import tencentcloud_tdmysql_dbs_backup.example tdsqlc-abc123#10086
```

