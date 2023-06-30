---
subcategory: "SQLServer"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_datasource_ins_attribute"
sidebar_current: "docs-tencentcloud-datasource-sqlserver_datasource_ins_attribute"
description: |-
  Use this data source to query detailed information of sqlserver datasource_ins_attribute
---

# tencentcloud_sqlserver_datasource_ins_attribute

Use this data source to query detailed information of sqlserver datasource_ins_attribute

## Example Usage

```hcl
data "tencentcloud_sqlserver_datasource_ins_attribute" "datasource_ins_attribute" {
  instance_id = "mssql-gyg9xycl"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Instance ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `blocked_threshold` - Block process threshold in milliseconds.
* `event_save_days` - Retention period for the files of slow SQL, blocking, deadlock, and extended events.
* `regular_backup_counts` - The number of retained archive backups.
* `regular_backup_enable` - Archive backup status. Valid values: enable (enabled), disable (disabled).
* `regular_backup_save_days` - Archive backup retention period: [90-3650] days.
* `regular_backup_start_time` - Archive backup start date in YYYY-MM-DD format, which is the current time by default.
* `regular_backup_strategy` - Archive backup policy. Valid values: years (yearly); quarters (quarterly);months` (monthly).
* `tde_config` - TDE Transparent Data Encryption Configuration.
  * `certificate_attribution` - Certificate ownership. Self - indicates using the account's own certificate, others - indicates referencing certificates from other accounts, and none - indicates no certificate.
  * `encryption` - TDE encryption, 'enable' - enabled, 'disable' - not enabled.
  * `quote_uin` - Other primary account IDs referenced when activating TDE encryption
Note: This field may return null, indicating that a valid value cannot be obtained.


