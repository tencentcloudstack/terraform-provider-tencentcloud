---
subcategory: "TencentDB for MariaDB(MariaDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mariadb_grant_account_privileges"
sidebar_current: "docs-tencentcloud-resource-mariadb_grant_account_privileges"
description: |-
  Provides a resource to create a mariadb grant_account_privileges
---

# tencentcloud_mariadb_grant_account_privileges

Provides a resource to create a mariadb grant_account_privileges

## Example Usage

```hcl
resource "tencentcloud_mariadb_grant_account_privileges" "grant_account_privileges" {
  instance_id = "tdsql-9vqvls95"
  user_name   = "keep-modify-privileges"
  host        = "127.0.0.1"
  db_name     = "*"
  privileges  = ["SELECT", "INSERT"]
}
```

## Argument Reference

The following arguments are supported:

* `db_name` - (Required, String, ForceNew) Database name. `*` indicates that global permissions will be set (i.e., `*.*`), in which case the `Type` and `Object ` parameters will be ignored. If `DbName` is not `*`, the input parameter `Type` is required.
* `host` - (Required, String, ForceNew) Access host allowed for user. An account is uniquely identified by username and host.
* `instance_id` - (Required, String, ForceNew) Instance ID, which is in the format of `tdsql-ow728lmc` and can be obtained through the `DescribeDBInstances` API.
* `privileges` - (Required, Set: [`String`], ForceNew) Global permissions: SELECT, Insert, UPDATE, DELETE, CREATE, DROP, REFERENCES, INDEX, ALT, CREATE TEMPORARY TABLES, LOCK TABLES, EXECUTE, CREATE VIEW, SHOW VIEW, CREATE ROUTINE, ALT ROUTINE, EVENT, TRIGGER, SHOW DATABASES, REPLICATION CLIENT, REPLICATION SLAVE Library permissions: SELECT, Insert, UPDATE, DELETE, CRAVE EATE, DROP, REFERENCES, INDEX, ALT, CREATE TEMPORARY TABLES, LOCK TABLES, EXECUTE, CREATE VIEW, SHOW VIEW, CREATE ROUTINE, ALT ROUTINE, EVENT, TRIGGER Table/View Permissions: SELECT, Insert, UPDATE, DELETE, CREATE, DROP, REFERENCES, INDEX, ALT, CREATE VIEW, SHOW VIEW, TRIGGER Stored Procedure/Function Permissions: ALT ROUTINE, EXECUTE Field Permissions: Insert, REFERENCES, SELECT, UPDATE.
* `user_name` - (Required, String, ForceNew) Login username.
* `col_name` - (Optional, String, ForceNew) If `Type` is `table` and `ColName` is `*`, the permissions will be granted to the table; if `ColName` is a specific field name, the permissions will be granted to the field.
* `object` - (Optional, String, ForceNew) Type name. For example, if `Type` is `table`, `Object` indicates a specific table name; if both `DbName` and `Type` are specific names, it indicates a specific object name and cannot be `*` or empty.
* `type` - (Optional, String, ForceNew) Type. Valid values: table, view, proc, func, *. If `DbName` is a specific database name and `Type` is `*`, the permissions of the database will be set (i.e., `db.*`), in which case the `Object` parameter will be ignored.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



