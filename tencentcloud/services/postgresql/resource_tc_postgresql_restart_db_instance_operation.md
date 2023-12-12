Provides a resource to create a postgresql restart_db_instance_operation

Example Usage

```hcl
resource "tencentcloud_postgresql_restart_db_instance_operation" "restart_db_instance_operation" {
  db_instance_id = local.pgsql_id
}
```