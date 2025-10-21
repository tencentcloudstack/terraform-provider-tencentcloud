---
subcategory: "TencentDB for Redis(crs)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_redis_backup"
sidebar_current: "docs-tencentcloud-datasource-redis_backup"
description: |-
  Use this data source to query detailed information of redis backup
---

# tencentcloud_redis_backup

Use this data source to query detailed information of redis backup

## Example Usage

```hcl
data "tencentcloud_redis_backup" "backup" {
  instance_id   = "crs-c1nl9rpv"
  begin_time    = "2023-04-07 03:57:30"
  end_time      = "2023-04-07 03:57:56"
  status        = [2]
  instance_name = "Keep-terraform"
}
```

## Argument Reference

The following arguments are supported:

* `begin_time` - (Optional, String) start time, such as 2017-02-08 19:09:26.Query the list of backups that the instance started backing up during the [beginTime, endTime] time period.
* `end_time` - (Optional, String) End time, such as 2017-02-08 19:09:26.Query the list of backups that the instance started backing up during the [beginTime, endTime] time period.
* `instance_id` - (Optional, String) The ID of instance.
* `instance_name` - (Optional, String) Instance name, which supports fuzzy search based on instance name.
* `result_output_file` - (Optional, String) Used to save results.
* `status` - (Optional, Set: [`Int`]) Status of the backup task:1: Backup is in the process.2: The backup is normal.3: Backup to RDB file processing.4: RDB conversion completed.-1: The backup has expired.-2: Backup deleted.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `backup_set` - An array of backups for the instance.
  * `backup_id` - Backup ID.
  * `backup_size` - Internal fields, which can be ignored by the user.
  * `backup_type` - Backup type.1: User-initiated manual backup.0: System-initiated backup in the early morning.
  * `end_time` - Backup end time.
  * `expire_time` - Backup file expiration time.
  * `file_type` - Back up file types.
  * `full_backup` - Internal fields, which can be ignored by the user.
  * `instance_id` - The ID of instance.
  * `instance_name` - The name of instance.
  * `instance_type` - Internal fields, which can be ignored by the user.
  * `locked` - Whether the backup is locked.0: Not locked.1: Has been locked.
  * `region` - The region where the backup is located.
  * `remark` - Notes information for the backup.
  * `start_time` - Backup start time.
  * `status` - Backup status.1: The backup is locked by another process.2: The backup is normal and not locked by any process.-1: The backup has expired.3: The backup is being exported.4: The backup export is successful.


