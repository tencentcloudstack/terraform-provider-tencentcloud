---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_backup_list"
sidebar_current: "docs-tencentcloud-tencentcloud_mysql_backup_list"
description: |-
Use this data source to query the list of backup databases.

---
#tencentcloud_mysql_backup_list

Use this data source to query the list of backup databases.

##Example Usage
```
resource "tencentcloud_mysql_account" "default" { 
    mysql_id = "my-test-database" 
    name = "tf_account" 
    password = "…" 
    description = "My test account" 
}
```
 

##Argument Reference

The following arguments are supported:

- `mysql_id` - (Required) Instance ID, such as cdb-c1nl9rpv. It is identical to the instance ID displayed in the database console page.

- `max_number` - (Optional) The latest files to list, rang from 1 to 10000. And the default value is 10.

- `result_output_file` - (Optional) Used to store results.

##Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `time` - The earliest time at which the backup starts. For example, 2 indicates 2:00 am.

- `finish_time` - The time at which the backup finishes.

- `size` - the size of backup file.

- `backup_id` - ID of Backup task.

- `backup_model` - Backup method. Supported values include: physical - physical backup, and logical - logical backup.

- `intranet_url` - URL for downloads internally.

- `internet_url` - URL for downloads externally.

- `creator` - The owner of the backup files.