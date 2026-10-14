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
  time        = "2023-05-23 01:00:00"
  rename_restore {
    old_name = "keep_pubsub_db2"
    new_name = "rollback_pubsub_db3"
  }
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Instance ID.
* `rename_restore` - (Required, List) Rename the databases listed in ReNameRestoreDatabase.
* `time` - (Required, String) Target time point for rollback.

The `rename_restore` object supports the following:

* `new_name` - (Required, String) New database name.
* `old_name` - (Required, String) Database name. If the OldName database does not exist, a failure will be returned. It can be left empty in offline migration tasks.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `encryption` - TDE encryption, `enable` encrypted, `disable` unencrypted.
  * `db_name` - Database name.
  * `status` - encryption, `enable` encrypted, `disable` unencrypted.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [DescribeFlowStatus](https://www.tencentcloud.com/document/product/238/32113) - This API is used to query flow status.
* Create: [RollbackInstance](https://www.tencentcloud.com/document/product/238/32095) - This API is used to roll back the instance by time point.
* Read: [DescribeDBs](https://www.tencentcloud.com/document/product/238/32114) - This API is used to query the list of databases
* Update: [DescribeFlowStatus](https://www.tencentcloud.com/document/product/238/32113) - This API is used to query flow status.
* Update: [RollbackInstance](https://www.tencentcloud.com/document/product/238/32095) - This API is used to roll back the instance by time point.
* Delete: [DeleteDB](https://www.tencentcloud.com/document/product/238/32117) - This API is used to drop a database.
* Delete: [DescribeFlowStatus](https://www.tencentcloud.com/document/product/238/32113) - This API is used to query flow status.

## Import

sqlserver rollback_instance can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_rollback_instance.rollback_instance mssql-qelbzgwf#2023-05-23 01:00:00#keep_pubsub_db2#rollback_pubsub_db3
```

