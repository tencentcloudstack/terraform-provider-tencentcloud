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
* `column_privileges` - (Optional, List) &amp;quot;Permissions for columns in database tables. Optional values for the Privileges permission are:&amp;quot;&amp;quot;SELECT, INSERT, UPDATE, REFERENCES.&amp;quot;&amp;quot;Note that if this parameter is not passed, the existing privileges are reserved. If you need to clear them, please pass an empty array in the complex type Privileges field.&amp;quot;.
* `database_privileges` - (Optional, List) &amp;quot;Database permissions. Optional values for the Privileges permission are:&amp;quot;&amp;quot;SELECT, INSERT, UPDATE, DELETE, CREATE, DROP, REFERENCES, INDEX, ALTER, CREATE TEMPORARY TABLES,&amp;quot;&amp;quot;LOCK TABLES, EXECUTE, CREATE VIEW, SHOW VIEW, CREATE ROUTINE, ALTER ROUTINE, EVENT, TRIGGER.&amp;quot;&amp;quot;Note that if this parameter is not passed, the existing privileges are reserved. If you need to clear them, please pass an empty array in the complex type Privileges field.&amp;quot;.
* `global_privileges` - (Optional, Set: [`String`]) &amp;quot;Global permissions. Among them, the optional value of the permission in GlobalPrivileges is:&amp;quot;&amp;quot;SELECT, INSERT, UPDATE, DELETE, CREATE, PROCESS, DROP, REFERENCES, INDEX, ALTER, SHOW DATABASES,&amp;quot;&amp;quot;CREATE TEMPORARY TABLES, LOCK TABLES, EXECUTE, CREATE VIEW, SHOW VIEW, CREATE ROUTINE, ALTER ROUTINE, EVENT, TRIGGER.&amp;quot;&amp;quot;Note that if this parameter is not passed, it means that the existing permissions are reserved. If it needs to be cleared, pass an empty array in this field.&amp;quot;.
* `table_privileges` - (Optional, List) &amp;quot;Permissions for tables in the database. Optional values for the Privileges permission are:&amp;quot;&amp;quot;SELECT, INSERT, UPDATE, DELETE, CREATE, DROP, REFERENCES, INDEX, ALTER, CREATE VIEW, SHOW VIEW, TRIGGER.&amp;quot;&amp;quot;Note that if this parameter is not passed, the existing privileges are reserved. If you need to clear them, please pass an empty array in the complex type Privileges field.&amp;quot;.
* `view_privileges` - (Optional, List) &amp;quot;Permissions for database views. Optional values for the Privileges permission are:&amp;quot;&amp;quot;SELECT, INSERT, UPDATE, DELETE, CREATE, DROP, REFERENCES, INDEX, ALTER, CREATE VIEW, SHOW VIEW, TRIGGER.&amp;quot;&amp;quot;Note that if this parameter is not passed, the existing privileges are reserved. If you need to clear them, please pass an empty array in the complex type Privileges field.&amp;quot;.

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



## Import

dcdb account_privileges can be imported using the id, e.g.

```
terraform import tencentcloud_dcdb_account_privileges.account_privileges instanceId#userName#host#dbName#tabName#viewName#colName
```

