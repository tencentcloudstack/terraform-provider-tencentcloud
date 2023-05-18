---
subcategory: "SQLServer"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_upload_incremental_info"
sidebar_current: "docs-tencentcloud-datasource-sqlserver_upload_incremental_info"
description: |-
  Use this data source to query detailed information of sqlserver upload_incremental_info
---

# tencentcloud_sqlserver_upload_incremental_info

Use this data source to query detailed information of sqlserver upload_incremental_info

## Example Usage

```hcl
data "tencentcloud_sqlserver_upload_incremental_info" "upload_incremental_info" {
  instance_id              = "mssql-j8kv137v"
  backup_migration_id      = "migration_id"
  incremental_migration_id = ""
}
```

## Argument Reference

The following arguments are supported:

* `backup_migration_id` - (Required, String) Backup import task ID, which is returned through the API CreateBackupMigration.
* `incremental_migration_id` - (Required, String) ID of the incremental import task.
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


