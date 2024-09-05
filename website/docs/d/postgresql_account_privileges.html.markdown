---
subcategory: "TencentDB for PostgreSQL(PostgreSQL)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_postgresql_account_privileges"
sidebar_current: "docs-tencentcloud-datasource-postgresql_account_privileges"
description: |-
  Use this data source to query detailed information of postgresql account privileges
---

# tencentcloud_postgresql_account_privileges

Use this data source to query detailed information of postgresql account privileges

## Example Usage

```hcl
data "tencentcloud_postgresql_account_privileges" "example" {
  db_instance_id = "postgres-3hk6b6tj"
  user_name      = "tf_example"
  database_object_set {
    object_name = "postgres"
    object_type = "database"
  }
}
```

## Argument Reference

The following arguments are supported:

* `database_object_set` - (Required, List) Instance database object info.
* `db_instance_id` - (Required, String) Instance ID.
* `user_name` - (Required, String) Instance username.
* `result_output_file` - (Optional, String) Used to save results.

The `database_object_set` object supports the following:

* `object_name` - (Required, String) Database object Name.Note: This field may return null, indicating that no valid value can be obtained.
* `object_type` - (Required, String) Supported database object types: account, database, schema, sequence, procedure, type, function, table, view, matview, column. Note: This field may return null, indicating that no valid value can be obtained.
* `database_name` - (Optional, String) Database name to which the database object belongs. This parameter is mandatory when ObjectType is not database.Note: This field may return null, indicating that no valid value can be obtained.
* `schema_name` - (Optional, String) Schema name to which the database object belongs. This parameter is mandatory when ObjectType is not database or schema.Note: This field may return null, indicating that no valid value can be obtained.
* `table_name` - (Optional, String) Table name to which the database object belongs. This parameter is mandatory when ObjectType is column.Note: This field may return null, indicating that no valid value can be obtained.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `privilege_set` - Privilege list.
  * `object` - Database object.If ObjectType is database, DatabaseName/SchemaName/TableName can be null.If ObjectType is schema, SchemaName/TableName can be null.If ObjectType is table, TableName can be null.If ObjectType is column, DatabaseName/SchemaName/TableName can&amp;#39;t be null.In all other cases, DatabaseName/SchemaName/TableName can be null. Note: This field may return null, indicating that no valid value can be obtained.
    * `database_name` - Database name to which the database object belongs. This parameter is mandatory when ObjectType is not database. Note: This field may return null, indicating that no valid value can be obtained.
    * `object_name` - Database object Name. Note: This field may return null, indicating that no valid value can be obtained.
    * `object_type` - Supported database object types: account, database, schema, sequence, procedure, type, function, table, view, matview, column. Note: This field may return null, indicating that no valid value can be obtained.
    * `schema_name` - Schema name to which the database object belongs. This parameter is mandatory when ObjectType is not database or schema. Note: This field may return null, indicating that no valid value can be obtained.
    * `table_name` - Table name to which the database object belongs. This parameter is mandatory when ObjectType is column. Note: This field may return null, indicating that no valid value can be obtained.
  * `privilege_set` - Privileges the specific account has on database object. Note: This field may return null, indicating that no valid value can be obtained.


