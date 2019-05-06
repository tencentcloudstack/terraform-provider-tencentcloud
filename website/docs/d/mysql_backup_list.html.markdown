---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_backup_list"
sidebar_current: "docs-tencentcloud-datasource-mysql_backup_list"
description: |-
  Use this data source to query the list of backup databases.
---

# tencentcloud_mysql_backup_list

Use this data source to query the list of backup databases.

## Example Usage

```hcl
resource "tencentcloud_mysql_backup_list" "default" {
  mysql_id = "my-test-database"
  max_number = 10
  result_output_file = "mytestpath"
}

## Argument Reference

The following arguments are supported:

* `mysql_id` - (Required, ForceNew) Instance ID, such as cdb-c1nl9rpv. It is identical to the instance ID displayed in the database console page.
* `max_number` - (Optional, ForceNew) The latest files to list, rang from 1 to 10000. And the default value is 10.
* `result_output_file` - (Optional, ForceNew) Used to store results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - A list of MySQL backup. Each element contains the following attributes:
  * `backup_id` - ID of Backup task.
  * `backup_model` - Backup method. Supported values include: physical - physical backup, and logical - logical backup.
  * `creator` - The owner of the backup files.
  * `finish_time` - The time at which the backup finishes.
  * `internet_url` - URL for downloads externally.
  * `intranet_url` - URL for downloads internally.
  * `size` - the size of backup file.
  * `time` - The earliest time at which the backup starts. For example, 2 indicates 2:00 am.

