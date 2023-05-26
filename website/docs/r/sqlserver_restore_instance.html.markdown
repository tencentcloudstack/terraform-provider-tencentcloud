---
subcategory: "SQLServer"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_restore_instance"
sidebar_current: "docs-tencentcloud-resource-sqlserver_restore_instance"
description: |-
  Provides a resource to create a sqlserver restore_instance
---

# tencentcloud_sqlserver_restore_instance

Provides a resource to create a sqlserver restore_instance

## Example Usage

```hcl
resource "tencentcloud_sqlserver_restore_instance" "restore_instance" {
  instance_id = "mssql-qelbzgwf"
  backup_id   = 3461718019
  rename_restore {
    old_name = "keep_pubsub_db2"
    new_name = "new_keep_pubsub_db2"
  }
  type    = 1
  db_list = ["keep_pubsub_db2"]
}
```

## Argument Reference

The following arguments are supported:

* `backup_id` - (Required, Int) Backup file ID, which can be obtained through the Id field in the returned value of the DescribeBackups API.
* `instance_id` - (Required, String) Instance ID.
* `db_list` - (Optional, Set: [`String`]) The database that needs to be covered and rolled back is required only when the file is covered and rolled back.
* `rename_restore` - (Optional, List) Restore the databases listed in ReNameRestoreDatabase and rename them after restoration. If this parameter is left empty, all databases will be restored and renamed in the default format.
* `type` - (Optional, Int) Rollback type, 0-overwrite method; 1-rename method, default 1.

The `rename_restore` object supports the following:

* `new_name` - (Optional, String) New database name. In offline migration, OldName will be used if NewName is left empty (OldName and NewName cannot be both empty). In database cloning, OldName and NewName must be both specified and cannot have the same value.
* `old_name` - (Optional, String) Database name. If the OldName database does not exist, a failure will be returned.It can be left empty in offline migration tasks.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

sqlserver restore_instance can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_restore_instance.restore_instance restore_instance_id
```

