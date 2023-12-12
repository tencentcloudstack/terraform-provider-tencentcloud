Provides a resource to create a postgresql delete_log_backup_operation

Example Usage

```hcl
resource "tencentcloud_postgresql_delete_log_backup_operation" "delete_log_backup_operation" {
  db_instance_id = "local.pg_id"
  log_backup_id = "local.pg_log_backup_id"
}
```