---
subcategory: "TencentDB for MySQL(cdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_bin_log"
sidebar_current: "docs-tencentcloud-datasource-mysql_bin_log"
description: |-
  Use this data source to query detailed information of mysql bin_log
---

# tencentcloud_mysql_bin_log

Use this data source to query detailed information of mysql bin_log

## Example Usage

```hcl
data "tencentcloud_mysql_bin_log" "bin_log" {
  instance_id = "cdb-fitq5t9h"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Instance ID, in the format: cdb-c1nl9rpv. Same instance ID as displayed in the ApsaraDB for Console page.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `items` - Details of binary log files that meet the query conditions.
  * `binlog_finish_time` - binlog file deadline.
  * `binlog_start_time` - Binlog file start time.
  * `cos_storage_type` - Storage method, 0-regular storage, 1-archive storage, the default is 0.
  * `date` - File storage time, time format: 2016-03-17 02:10:37.
  * `instance_id` - Instance ID, in the format: cdb-c1nl9rpv. Same instance ID as displayed in the ApsaraDB for Console page.
  * `internet_url` - download link.
  * `intranet_url` - download link.
  * `name` - binlog log backup file name.
  * `region` - The region where the local binlog file is located.
  * `remote_info` - Binlog remote backup details.
    * `finish_time` - End time of remote backup task.
    * `region` - The region where remote backup is located.
    * `start_time` - Start time of remote backup task.
    * `status` - Backup task status. Possible values are `SUCCESS`: backup succeeded, `FAILED`: backup failed, `RUNNING`: backup in progress.
    * `sub_backup_id` - The ID of the remote backup subtask.
    * `url` - download link.
  * `size` - Backup file size, unit: Byte.
  * `status` - Backup task status. Possible values are `SUCCESS`: backup succeeded, `FAILED`: backup failed, `RUNNING`: backup in progress.
  * `type` - Specific log type, possible values are: binlog - binary log.


