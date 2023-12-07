Provides a resource to create a postgresql isolate_db_instance_operation

Example Usage

```hcl
resource "tencentcloud_postgresql_isolate_db_instance_operation" "isolate_db_instance_operation" {
  db_instance_id_set = [local.pgsql_id]
}
```