---
subcategory: "TencentDB for PostgreSQL(PostgreSQL)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_postgresql_account_privileges_operation"
sidebar_current: "docs-tencentcloud-resource-postgresql_account_privileges_operation"
description: |-
  Provides a resource to modify the privileges of a PostgreSQL account (one-time operation)
---

# tencentcloud_postgresql_account_privileges_operation

Provides a resource to modify the privileges of a PostgreSQL account (one-time operation)

## Example Usage

```hcl
resource "tencentcloud_postgresql_account_privileges_operation" "example" {
  db_instance_id = "postgres-4wdeb0zv"
  user_name      = "tf_example"

  modify_privilege_set {
    modify_type = "grantObject"

    database_privilege {
      object {
        object_type   = "table"
        object_name   = "tf_table"
        database_name = "tf_db"
        schema_name   = "public"
      }

      privilege_set = ["SELECT", "INSERT"]
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `db_instance_id` - (Required, String, ForceNew) Instance ID in the format of postgres-4wdeb0zv.
* `modify_privilege_set` - (Required, List) Privileges to modify. Batch modification supported, up to 50 entries at a time.
* `user_name` - (Required, String, ForceNew) Instance username.

The `database_privilege` object of `modify_privilege_set` supports the following:

* `object` - (Optional, List) Database object.If ObjectType is database, DatabaseName/SchemaName/TableName can be null.If ObjectType is schema, SchemaName/TableName can be null.If ObjectType is table, TableName can be null.If ObjectType is column, DatabaseName/SchemaName/TableName can&amp;#39;t be null.In all other cases, DatabaseName/SchemaName/TableName can be null. Note: This field may return null, indicating that no valid value can be obtained.
* `privilege_set` - (Optional, Set) Privileges the specific account has on database object. Note: This field may return null, indicating that no valid value can be obtained.

The `modify_privilege_set` object supports the following:

* `database_privilege` - (Optional, List) Database objects and the user permissions on these objects. Note: This field may return null, indicating that no valid value can be obtained.
* `is_cascade` - (Optional, Bool) Required only when ModifyType is revokeObject. When the parameter is true, revoking permissions will cascade. The default value is false.
* `modify_type` - (Optional, String) Supported modification method: grantObject, revokeObject, alterRole. grantObject represents granting permissions on object, revokeObject represents revoking permissions on object, and alterRole represents modifying the account type.

The `object` object of `database_privilege` supports the following:

* `object_name` - (Required, String) Database object Name. Note: This field may return null, indicating that no valid value can be obtained.
* `object_type` - (Required, String) Supported database object types: account, database, schema, sequence, procedure, type, function, table, view, matview, column. Note: This field may return null, indicating that no valid value can be obtained.
* `database_name` - (Optional, String) Database name to which the database object belongs. This parameter is mandatory when ObjectType is not database. Note: This field may return null, indicating that no valid value can be obtained.
* `schema_name` - (Optional, String) Schema name to which the database object belongs. This parameter is mandatory when ObjectType is not database or schema. Note: This field may return null, indicating that no valid value can be obtained.
* `table_name` - (Optional, String) Table name to which the database object belongs. This parameter is mandatory when ObjectType is column. Note: This field may return null, indicating that no valid value can be obtained.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [ModifyAccountPrivileges](https://www.tencentcloud.com/document/api/409/78816)

