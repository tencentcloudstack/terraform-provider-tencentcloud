---
subcategory: "SQLServer"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_rollback_instance"
sidebar_current: "docs-tencentcloud-resource-sqlserver_rollback_instance"
description: |-
  Provides a resource to create a sqlserver rollback_instance
---

# tencentcloud_sqlserver_rollback_instance

Provides a resource to create a sqlserver rollback_instance

## Example Usage

```hcl
resource "tencentcloud_sqlserver_rollback_instance" "rollback_instance" {
  instance_id = "mssql-qelbzgwf"
  type        = 1
  time        = "2023-05-25 19:14:30"
  dbs         = ["keep_pubsub_db2"]
  rename_restore {
    old_name = "keep_pubsub_db2"
    new_name = "rollback_pubsub_db2"
  }
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Instance ID.
* `time` - (Required, String) Target time point for rollback.
* `type` - (Required, Int) Rollback type. 0: the database rolled back overwrites the original database; 1: the database rolled back is renamed and does not overwrite the original database.
* `dbs` - (Optional, Set: [`String`]) Database to be rolled back.
* `rename_restore` - (Optional, List) Rename the databases listed in ReNameRestoreDatabase. This parameter takes effect only when Type = 1 which indicates that backup rollback supports renaming databases. If it is left empty, databases will be renamed in the default format and the DBs parameter specifies the databases to be restored.

The `rename_restore` object supports the following:

* `new_name` - (Optional, String) New database name. In offline migration, OldName will be used if NewName is left empty (OldName and NewName cannot be both empty). In database cloning, OldName and NewName must be both specified and cannot have the same value.
* `old_name` - (Optional, String) Database name. If the OldName database does not exist, a failure will be returned. It can be left empty in offline migration tasks.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

sqlserver rollback_instance can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_rollback_instance.rollback_instance rollback_instance_id
```

