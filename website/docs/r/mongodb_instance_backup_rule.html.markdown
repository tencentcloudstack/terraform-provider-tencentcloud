---
subcategory: "TencentDB for MongoDB(mongodb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mongodb_instance_backup_rule"
sidebar_current: "docs-tencentcloud-resource-mongodb_instance_backup_rule"
description: |-
  Provides a resource to create mongodb instance backup rule
---

# tencentcloud_mongodb_instance_backup_rule

Provides a resource to create mongodb instance backup rule

## Example Usage

```hcl
resource "tencentcloud_mongodb_instance_backup_rule" "example" {
  instance_id             = "cmgo-rnht8d3d"
  backup_method           = 0
  backup_time             = 10
  backup_retention_period = 7
  backup_version          = 1
}
```

## Argument Reference

The following arguments are supported:

* `backup_method` - (Required, Int) Set automatic backup method. Valid values:
- 0: Logical backup;
- 1: Physical backup;
- 3: Snapshot backup (supported only in cloud disk version).
* `backup_time` - (Required, Int) Set the start time for automatic backup. The value range is: [0,23]. For example, setting this parameter to 2 means that backup starts at 02:00.
* `instance_id` - (Required, String, ForceNew) Instance ID.
* `active_weekdays` - (Optional, String) Specify the specific dates for automatic backups to be performed each week. Format: Enter a number between 0 and 6 to represent Sunday through Saturday (e.g., 1 represents Monday). Separate multiple dates with commas (,). Example: Entering 1,3,5 means the system will perform backups on Mondays, Wednesdays, and Fridays every week. Default: If not set, the default is a full cycle (0,1,2,3,4,5,6), meaning backups will be performed daily.
* `alarm_water_level` - (Optional, Int) Sets the alarm threshold for backup dataset storage space usage. Unit: %. Default value: 100. Value range: [50, 300].
* `backup_frequency` - (Optional, Int) Specify the daily automatic backup frequency. 12: Back up twice a day, approximately 12 hours apart; 24: Back up once a day (default), approximately 24 hours apart.
* `backup_retention_period` - (Optional, Int) Specifies the retention period for backup data. Unit: days, default is 7 days. Value range: [7, 365].
* `backup_version` - (Optional, Int) Backup version. Old version backup is 0, advanced backup is 1. Set this value to 1 when enabling advanced backup.
* `long_term_active_days` - (Optional, String) Specify the specific backup dates to be retained long-term. This setting only takes effect when LongTermUnit is set to weekly or monthly. Weekly Retention: Enter a number between 0 and 6 to represent Sunday through Saturday. Separate multiple dates with commas. Monthly Retention: Enter a number between 1 and 31 to represent specific dates within the month. Separate multiple dates with commas.
* `long_term_expired_days` - (Optional, Int) Long-term backup retention period. Value range [30, 1075].
* `long_term_unit` - (Optional, String) Long-term retention period. Supports selecting specific dates for backups on a weekly or monthly basis (e.g., backup data for the 1st and 15th of each month) to retain for a longer period. Disabled (default): Long-term retention is disabled. Weekly retention: Specify `weekly`. Monthly retention: Specify `monthly`.
* `notify` - (Optional, Bool) Set whether to send failure alerts when automatic backup errors occur.
- true: Send.
- false: Do not send.
* `oplog_expired_days` - (Optional, Int) Incremental backup retention period. Unit: days. Default value: 7 days. Value range: [7,365].

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

mongodb instance backup rule can be imported using the id, e.g.

```
terraform import tencentcloud_mongodb_instance_backup_rule.example cmgo-rnht8d3d
```

