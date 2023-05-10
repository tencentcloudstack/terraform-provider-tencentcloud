---
subcategory: "TencentDB for MongoDB(mongodb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mongodb_instance_backups"
sidebar_current: "docs-tencentcloud-datasource-mongodb_instance_backups"
description: |-
  Use this data source to query detailed information of mongodb instance_backups
---

# tencentcloud_mongodb_instance_backups

Use this data source to query detailed information of mongodb instance_backups

## Example Usage

```hcl
data "tencentcloud_mongodb_instance_backups" "instance_backups" {
  instance_id   = "cmgo-9d0p6umb"
  backup_method = 0
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Instance ID, the format is: cmgo-9d0p6umb.Same as the instance ID displayed in the cloud database console page.
* `backup_method` - (Optional, Int) Backup mode, currently supported: 0-logic backup, 1-physical backup, 2-all backups.The default is logical backup.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `backup_list` - backup list.
  * `backup_desc` - Remark of backup.
  * `backup_method` - Backup method.
  * `backup_name` - Backup mode name.
  * `backup_size` - Size of backup(KN).
  * `backup_type` - Backup mode type.
  * `end_time` - end time of backup.
  * `instance_id` - Instance ID.
  * `start_time` - start time of backup.
  * `status` - Backup status.


