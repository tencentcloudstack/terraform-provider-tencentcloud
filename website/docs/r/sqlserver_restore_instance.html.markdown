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
  backup_id   = 3482091273
  rename_restore {
    old_name = "keep_pubsub_db2"
    new_name = "restore_keep_pubsub_db2"
  }
}
```

## Argument Reference

The following arguments are supported:

* `backup_id` - (Required, Int) Backup file ID, which can be obtained through the Id field in the returned value of the DescribeBackups API.
* `instance_id` - (Required, String) Instance ID.
* `rename_restore` - (Required, List) Restore the databases listed in ReNameRestoreDatabase and rename them after restoration. If this parameter is left empty, all databases will be restored and renamed in the default format.

The `rename_restore` object supports the following:

* `new_name` - (Required, String) New database name. In offline migration, OldName will be used if NewName is left empty (OldName and NewName cannot be both empty). In database cloning, OldName and NewName must be both specified and cannot have the same value.
* `old_name` - (Required, String) Database name. If the OldName database does not exist, a failure will be returned.It can be left empty in offline migration tasks.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `encryption` - TDE encryption, `enable` encrypted, `disable` unencrypted.
  * `db_name` - Database name.
  * `status` - encryption, `enable` encrypted, `disable` unencrypted.


## Import

sqlserver restore_instance can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_restore_instance.restore_instance mssql-qelbzgwf#3482091273#keep_pubsub_db2#restore_keep_pubsub_db2
```

