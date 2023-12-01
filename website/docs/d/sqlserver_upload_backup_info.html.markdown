---
subcategory: "SQLServer"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_upload_backup_info"
sidebar_current: "docs-tencentcloud-datasource-sqlserver_upload_backup_info"
description: |-
  Use this data source to query detailed information of sqlserver upload_backup_info
---

# tencentcloud_sqlserver_upload_backup_info

Use this data source to query detailed information of sqlserver upload_backup_info

## Example Usage

```hcl
data "tencentcloud_sqlserver_upload_backup_info" "example" {
  instance_id         = "mssql-qelbzgwf"
  backup_migration_id = "mssql-backup-migration-8a0f3eht"
}
```

## Argument Reference

The following arguments are supported:

* `backup_migration_id` - (Required, String) Backup import task ID, which is returned through the API CreateBackupMigration.
* `instance_id` - (Required, String) Instance ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `bucket_name` - Bucket name.
* `expired_time` - Temporary key expiration time.
* `path` - Storage path.
* `region` - Bucket location information.
* `start_time` - Temporary key start time.
* `tmp_secret_id` - Temporary key ID.
* `tmp_secret_key` - Temporary key (Key).
* `x_cos_security_token` - Temporary key (Token).


