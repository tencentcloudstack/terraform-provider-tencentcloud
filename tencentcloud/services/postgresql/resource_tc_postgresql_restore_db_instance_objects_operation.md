Use this resource to restore database objects (databases, tables) on a PostgreSQL instance from a backup set or a point-in-time target.

~> **NOTE:** This is a one-time operation resource. Destroying this resource does nothing. To re-execute the operation, use `terraform taint` or re-create the resource.

Example Usage

Restore by backup set

```hcl
resource "tencentcloud_postgresql_restore_db_instance_objects_operation" "example" {
  db_instance_id  = "postgres-6bwgamo3"
  restore_objects = ["user"]
  backup_set_id   = "your-backup-set-id"
}
```

Restore by point-in-time

```hcl
resource "tencentcloud_postgresql_restore_db_instance_objects_operation" "example" {
  db_instance_id      = "postgres-6bwgamo3"
  restore_objects     = ["user", "orders"]
  restore_target_time = "2024-04-30 00:20:27"
}
```
