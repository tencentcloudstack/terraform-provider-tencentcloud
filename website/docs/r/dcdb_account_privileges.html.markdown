---
subcategory: "TDSQL for MySQL(DCDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dcdb_account_privileges"
sidebar_current: "docs-tencentcloud-resource-dcdb_account_privileges"
description: |-
  Provides a resource to create a dcdb account_privileges
---

# tencentcloud_dcdb_account_privileges

Provides a resource to create a dcdb account_privileges

## Example Usage

```hcl
resource "tencentcloud_dcdb_account_privileges" "account_privileges" {
  instance_id = "%s"
  account {
    user = "tf_test"
    host = "%s"
  }
  global_privileges = ["SHOW DATABASES", "SHOW VIEW"]
  database_privileges {
    privileges = ["SELECT", "INSERT", "UPDATE", "DELETE", "CREATE"]
    database   = "tf_test_db"
  }

  table_privileges {
    database   = "tf_test_db"
    table      = "tf_test_table"
    privileges = ["SELECT", "INSERT", "UPDATE", "DELETE", "CREATE"]

  }
```

## Argument Reference

The following arguments are supported:

* `account` - (Required, List) The account of the database, including username and host.
* `instance_id` - (Required, String) The ID of instance.
* `column_privileges` - (Optional, List) Permissions for columns in database tables. Optional values for the Privileges permission are:  SELECT, INSERT, UPDATE, REFERENCES.  Note that if this parameter is not passed, the existing privileges are reserved. If you need to clear them, please pass an empty array in the complex type Privileges field.
* `database_privileges` - (Optional, List) Database permissions. Optional values for the Privileges permission are: SELECT, INSERT, UPDATE, DELETE, CREATE, DROP, REFERENCES, INDEX, ALTER, CREATE TEMPORARY TABLES,  LOCK TABLES, EXECUTE, CREATE VIEW, SHOW VIEW, CREATE ROUTINE, ALTER ROUTINE, EVENT, TRIGGER.  Note that if this parameter is not passed, the existing privileges are reserved. If you need to clear them, please pass an empty array in the complex type Privileges field.
* `global_privileges` - (Optional, Set: [`String`]) Global permissions. Among them, the optional value of the permission in GlobalPrivileges is: SELECT, INSERT, UPDATE, DELETE, CREATE, PROCESS, DROP, REFERENCES, INDEX, ALTER, SHOW DATABASES,  CREATE TEMPORARY TABLES, LOCK TABLES, EXECUTE, CREATE VIEW, SHOW VIEW, CREATE ROUTINE, ALTER ROUTINE, EVENT, TRIGGER.  Note that if this parameter is not passed, it means that the existing permissions are reserved. If it needs to be cleared, pass an empty array in this field.
* `table_privileges` - (Optional, List) Permissions for tables in the database. Optional values for the Privileges permission are: SELECT, INSERT, UPDATE, DELETE, CREATE, DROP, REFERENCES, INDEX, ALTER, CREATE VIEW, SHOW VIEW, TRIGGER. Note that if this parameter is not passed, the existing privileges are reserved. If you need to clear them, please pass an empty array in the complex type Privileges field.
* `view_privileges` - (Optional, List) Permissions for database views. Optional values for the Privileges permission are:  SELECT, INSERT, UPDATE, DELETE, CREATE, DROP, REFERENCES, INDEX, ALTER, CREATE VIEW, SHOW VIEW, TRIGGER.  Note that if this parameter is not passed, the existing privileges are reserved. If you need to clear them, please pass an empty array in the complex type Privileges field.

The `account` object supports the following:

* `host` - (Required, String) account host.
* `user` - (Required, String) account name.

The `column_privileges` object supports the following:

* `column` - (Required, String) Database column name.
* `database` - (Required, String) The name of database.
* `privileges` - (Required, Set) Permission information.
* `table` - (Required, String) Database table name.

The `database_privileges` object supports the following:

* `database` - (Required, String) The name of database.
* `privileges` - (Required, Set) Permission information.

The `table_privileges` object supports the following:

* `database` - (Required, String) The name of database.
* `privileges` - (Required, Set) Permission information.
* `table` - (Required, String) Database table name.

The `view_privileges` object supports the following:

* `database` - (Required, String) The name of database.
* `privileges` - (Required, Set) Permission information.
* `view` - (Required, String) Database view name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [DescribeFlow](https://cloud.tencent.com/document/api/557/56485) - This API is used to query task status.
* Create: [ModifyAccountPrivileges](https://cloud.tencent.com/document/api/557/83475) - This API is used to modify the permissions of a TencentDB instance account. \n\n**Note**\n-Only the SELECT permission (that is, set the permission parameter to `["SELECT"]`) of the system database `mysql` can be granted.
* Read: [DescribeAccountPrivileges](https://cloud.tencent.com/document/api/557/19994) - This API is used to query the permissions of a TencentDB account.
* Read: [DescribeDatabaseObjects](https://cloud.tencent.com/document/api/557/19989) - This API is used to query the list of database objects in a TencentDB instance, including tables, stored procedures, views, and functions.
* Read: [DescribeDatabaseTable](https://cloud.tencent.com/document/api/557/19988) - This API is used to query the table information of a TencentDB instance.
* Read: [DescribeDatabases](https://cloud.tencent.com/document/api/557/19987) - This API is used to query the database list of a TencentDB instance.
* Update: [DescribeFlow](https://cloud.tencent.com/document/api/557/56485) - This API is used to query task status.
* Update: [ModifyAccountPrivileges](https://cloud.tencent.com/document/api/557/83475) - This API is used to modify the permissions of a TencentDB instance account. \n\n**Note**\n-Only the SELECT permission (that is, set the permission parameter to `["SELECT"]`) of the system database `mysql` can be granted.

## Import

dcdb account_privileges can be imported using the id, e.g.

```
terraform import tencentcloud_dcdb_account_privileges.account_privileges instanceId#userName#host#dbName#tabName#viewName#colName
```

