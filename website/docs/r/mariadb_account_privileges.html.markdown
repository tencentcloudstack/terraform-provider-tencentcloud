---
subcategory: "TencentDB for MariaDB(MariaDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mariadb_account_privileges"
sidebar_current: "docs-tencentcloud-resource-mariadb_account_privileges"
description: |-
  Provides a resource to create a mariadb account_privileges
---

# tencentcloud_mariadb_account_privileges

Provides a resource to create a mariadb account_privileges

## Example Usage

```hcl
resource "tencentcloud_mariadb_account_privileges" "account_privileges" {
  instance_id = "tdsql-9vqvls95"
  accounts {
    user = "keep-modify-privileges"
    host = "127.0.0.1"
  }
  global_privileges = ["ALTER", "CREATE", "DELETE", "SELECT", "UPDATE", "DROP"]
}
```

## Argument Reference

The following arguments are supported:

* `accounts` - (Required, List, ForceNew) account information.
* `instance_id` - (Required, String, ForceNew) instance id.
* `column_privileges` - (Optional, List) Column permission. Valid values of `Privileges`: `SELECT`, `INSERT`, `UPDATE`, `REFERENCES`.Note: if the parameter is left empty, no change will be made to the granted column permissions. To clear the granted column permissions, set `Privileges` to an empty array.
* `database_privileges` - (Optional, List) Database permission. Valid values of `Privileges`: `SELECT`, `INSERT`, `UPDATE`, `DELETE`, `CREATE`, `DROP`, `REFERENCES`, `INDEX`, `ALTER`, `CREATE TEMPORARY TABLES`, `LOCK TABLES`, `EXECUTE`, `CREATE VIEW`, `SHOW VIEW`, `CREATE ROUTINE`, `ALTER ROUTINE`, `EVENT`, `TRIGGER`.Note: if the parameter is left empty, no change will be made to the granted database permissions. To clear the granted database permissions, set `Privileges` to an empty array.
* `function_privileges` - (Optional, List) Database function permissions. Valid values of `Privileges`: `ALTER ROUTINE`, `EXECUTE`.Note: if the parameter is not passed in, no change will be made to the granted function permissions. To clear the granted function permissions, set `Privileges` to an empty array.
* `global_privileges` - (Optional, Set: [`String`]) Global permission. Valid values of `GlobalPrivileges`: `SELECT`, `INSERT`, `UPDATE`, `DELETE`, `CREATE`, `PROCESS`, `DROP`, `REFERENCES`, `INDEX`, `ALTER`, `SHOW DATABASES`, `CREATE TEMPORARY TABLES`, `LOCK TABLES`, `EXECUTE`, `CREATE VIEW`, `SHOW VIEW`, `CREATE ROUTINE`, `ALTER ROUTINE`, `EVENT`, `TRIGGER`.Note: if the parameter is left empty, no change will be made to the granted global permissions. To clear the granted global permissions, set the parameter to an empty array.
* `procedure_privileges` - (Optional, List) Database stored procedure permission. Valid values of `Privileges`: `ALTER ROUTINE`, `EXECUTE`.Note: if the parameter is not passed in, no change will be made to the granted stored procedure permissions. To clear the granted stored procedure permissions, set `Privileges` to an empty array.
* `table_privileges` - (Optional, List) `SELECT`, `INSERT`, `UPDATE`, `DELETE`, `CREATE`, `DROP`, `REFERENCES`, `INDEX`, `ALTER`, `CREATE VIEW`, `SHOW VIEW`, `TRIGGER`.Note: if the parameter is not passed in, no change will be made to the granted table permissions. To clear the granted table permissions, set `Privileges` to an empty array.
* `view_privileges` - (Optional, List) Database view permission. Valid values of `Privileges`: `SELECT`, `INSERT`, `UPDATE`, `DELETE`, `CREATE`, `DROP`, `REFERENCES`, `INDEX`, `ALTER`, `CREATE VIEW`, `SHOW VIEW`, `TRIGGER`.Note: if the parameter is not passed in, no change will be made to the granted view permissions. To clear the granted view permissions, set `Privileges` to an empty array.

The `accounts` object supports the following:

* `host` - (Required, String) user host.
* `user` - (Required, String) user name.

The `column_privileges` object supports the following:

* `column` - (Required, String) Column name.
* `database` - (Required, String) Database name.
* `privileges` - (Required, Set) Permission information.
* `table` - (Required, String) Table name.

The `database_privileges` object supports the following:

* `database` - (Required, String) Database name.
* `privileges` - (Required, Set) Permission information.

The `function_privileges` object supports the following:

* `database` - (Required, String) Database name.
* `function_name` - (Required, String) Function name.
* `privileges` - (Required, Set) Permission information.

The `procedure_privileges` object supports the following:

* `database` - (Required, String) Database name.
* `privileges` - (Required, Set) Permission information.
* `procedure` - (Required, String) Procedure name.

The `table_privileges` object supports the following:

* `database` - (Required, String) Database name.
* `privileges` - (Required, Set) Permission information.
* `table` - (Required, String) Table name.

The `view_privileges` object supports the following:

* `database` - (Required, String) Database name.
* `privileges` - (Required, Set) Permission information.
* `view` - (Required, String) View name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [DescribeFlow](https://cloud.tencent.com/document/api/237/16177)
* Create: [ModifyAccountPrivileges](https://cloud.tencent.com/document/api/237/56486) - This API is used to modify the permissions of a TencentDB instance account. \n\n**Note**\n-Only the SELECT permission (that is, set the permission parameter to `["SELECT"]`) of the system database `mysql` can be granted.An error will be reported if read-write permissions are granted to a read-only account. If the parameter is not passed in, no change will be made to the granted table permissions. To clear the granted view permissions, set `Privileges` to an empty array.
* Read: [DescribeAccountPrivileges](https://cloud.tencent.com/document/api/237/16164) - This API is used to query the permissions of a TencentDB account.
* Update: [DescribeFlow](https://cloud.tencent.com/document/api/237/16177)
* Update: [ModifyAccountPrivileges](https://cloud.tencent.com/document/api/237/56486) - This API is used to modify the permissions of a TencentDB instance account. \n\n**Note**\n-Only the SELECT permission (that is, set the permission parameter to `["SELECT"]`) of the system database `mysql` can be granted.An error will be reported if read-write permissions are granted to a read-only account. If the parameter is not passed in, no change will be made to the granted table permissions. To clear the granted view permissions, set `Privileges` to an empty array.

## Import

mariadb account_privileges can be imported using the id, e.g.

```
terraform import tencentcloud_mariadb_account_privileges.account_privileges account_privileges_id
```

