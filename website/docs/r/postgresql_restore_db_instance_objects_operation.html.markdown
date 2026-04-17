---
subcategory: "TencentDB for PostgreSQL(PostgreSQL)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_postgresql_restore_db_instance_objects_operation"
sidebar_current: "docs-tencentcloud-resource-postgresql_restore_db_instance_objects_operation"
description: |-
  Use this resource to restore database objects (databases, tables) on a PostgreSQL instance from a backup set or a point-in-time target.
---

# tencentcloud_postgresql_restore_db_instance_objects_operation

Use this resource to restore database objects (databases, tables) on a PostgreSQL instance from a backup set or a point-in-time target.

~> **NOTE:** This is a one-time operation resource. Destroying this resource does nothing. To re-execute the operation, use `terraform taint` or re-create the resource.

## Example Usage

### Restore by backup set

```hcl
resource "tencentcloud_postgresql_restore_db_instance_objects_operation" "example" {
  db_instance_id  = "postgres-6bwgamo3"
  restore_objects = ["user"]
  backup_set_id   = "your-backup-set-id"
}
```

### Restore by point-in-time

```hcl
resource "tencentcloud_postgresql_restore_db_instance_objects_operation" "example" {
  db_instance_id      = "postgres-6bwgamo3"
  restore_objects     = ["user", "orders"]
  restore_target_time = "2024-04-30 00:20:27"
}
```

## Argument Reference

The following arguments are supported:

* `db_instance_id` - (Required, String, ForceNew) PostgreSQL instance ID, e.g. `postgres-6bwgamo3`.
* `restore_objects` - (Required, List: [`String`], ForceNew) List of database objects to restore. The restored object name format will be `${original}_bak_${timestamp}`.
* `backup_set_id` - (Optional, String, ForceNew) Backup set ID used for restoration. Exactly one of `backup_set_id` or `restore_target_time` must be specified.
* `restore_target_time` - (Optional, String, ForceNew) Point-in-time target for restoration (Beijing time), e.g. `2024-04-30 00:20:27`. Exactly one of `backup_set_id` or `restore_target_time` must be specified.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



