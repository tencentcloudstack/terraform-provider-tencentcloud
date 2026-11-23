Provides a resource to modify the privileges of a PostgreSQL account (one-time operation)

Example Usage

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
