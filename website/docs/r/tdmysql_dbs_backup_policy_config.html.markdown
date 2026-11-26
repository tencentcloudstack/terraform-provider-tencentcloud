---
subcategory: "TDSQL Boundless(tdmysql)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tdmysql_dbs_backup_policy_config"
sidebar_current: "docs-tencentcloud-resource-tdmysql_dbs_backup_policy_config"
description: |-
  Provides a resource to manage the TDSQL Boundless DBS backup policy configuration.
---

# tencentcloud_tdmysql_dbs_backup_policy_config

Provides a resource to manage the TDSQL Boundless DBS backup policy configuration.

## Example Usage

```hcl
resource "tencentcloud_tdmysql_dbs_backup_policy_config" "example" {
  instance_id = "tdsql-12345"
  backup_policy = {
    backup_method         = "physical"
    backup_start_time     = "02:00"
    backup_end_time       = "04:00"
    enable_full           = 1
    enable_log            = 1
    full_retention_period = 7
    log_retention_period  = 7
    period_time           = "Monday,Tuesday,Wednesday,Thursday,Friday,Saturday,Sunday"
    storage_type          = "COS"
  }
}
```

## Argument Reference

The following arguments are supported:

* `backup_policy` - (Required, Object) Backup policy.
* `instance_id` - (Required, String) Instance ID.

The `backup_policy` object supports the following:

* `backup_end_time` - (Optional, String) Backup end time.
* `backup_method` - (Optional, String) Backup method. Valid values: `physical` (physical backup), `snapshot` (snapshot backup).
* `backup_start_time` - (Optional, String) Backup start time.
* `enable_full` - (Optional, Int) Whether to enable full backup.
* `enable_log` - (Optional, Int) Whether to enable log backup.
* `full_retention_period` - (Optional, Int) Full backup retention period. Currently can only be set to 7 days.
* `log_retention_period` - (Optional, Int) Log retention days. Currently can only be set to 7 days.
* `period_time` - (Optional, String) Days of the week to perform backup.
* `storage_type` - (Optional, String) Storage type. Valid values: `COS`, `SNAPSHOT`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the backup policy configuration, same as the locator attribute.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create:
  - [DescribeDBSBackupPolicy](https://www.tencentcloud.com/document/api/1292/80117)
  - [ModifyDBSBackupPolicy](https://www.tencentcloud.com/document/api/1292/78276)
* Read: [DescribeDBSBackupPolicy](https://www.tencentcloud.com/document/api/1292/80117)
* Update:
  - [DescribeDBSBackupPolicy](https://www.tencentcloud.com/document/api/1292/80117)
  - [ModifyDBSBackupPolicy](https://www.tencentcloud.com/document/api/1292/78276)

## Import

TDSQL Boundless DBS backup policy config can be imported using the id, e.g.

```
terraform import tencentcloud_tdmysql_dbs_backup_policy_config.example tdsql-12345
```

