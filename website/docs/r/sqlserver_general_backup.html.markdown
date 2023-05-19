---
subcategory: "SQLServer"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_general_backup"
sidebar_current: "docs-tencentcloud-resource-sqlserver_general_backup"
description: |-
  Provides a resource to create a sqlserver general_backup
---

# tencentcloud_sqlserver_general_backup

Provides a resource to create a sqlserver general_backup

## Example Usage

```hcl
resource "tencentcloud_sqlserver_general_backup" "general_backup" {
  strategy    = 0
  instance_id = "mssql-qelbzgwf"
  backup_name = "create_sqlserver_backup_name"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Instance ID in the format of mssql-i1z41iwd.
* `backup_name` - (Optional, String) Backup name. If this parameter is left empty, a backup name in the format of [Instance ID]_[Backup start timestamp] will be automatically generated.
* `db_names` - (Optional, Set: [`String`]) List of names of databases to be backed up (required only for multi-database backup).
* `strategy` - (Optional, Int) Backup policy (0: instance backup, 1: multi-database backup).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `flow_id` - flow id.


## Import

sqlserver general_backups can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_general_backups.general_backups general_backups_id
```

