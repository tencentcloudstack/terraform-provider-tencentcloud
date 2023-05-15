---
subcategory: "SQLServer"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_backup_commands"
sidebar_current: "docs-tencentcloud-datasource-sqlserver_backup_commands"
description: |-
  Use this data source to query detailed information of sqlserver datasource_backup_command
---

# tencentcloud_sqlserver_backup_commands

Use this data source to query detailed information of sqlserver datasource_backup_command

## Example Usage

```hcl
data "tencentcloud_sqlserver_backup_commands" "backup_command" {
  backup_file_type = "FULL"
  data_base_name   = "db_name"
  is_recovery      = "No"
  local_path       = ""
}
```

## Argument Reference

The following arguments are supported:

* `backup_file_type` - (Required, String) Backup file type. Full: full backup. FULL_LOG: full backup which needs log increments. FULL_DIFF: full backup which needs differential increments. LOG: log backup. DIFF: differential backup.
* `data_base_name` - (Required, String) Database name.
* `is_recovery` - (Required, String) Whether restoration is required. No: not required. Yes: required.
* `local_path` - (Optional, String) Storage path of backup files. If this parameter is left empty, the default storage path will be D:.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - Command list.
  * `command` - Create backup command.
  * `request_id` - Request ID.


