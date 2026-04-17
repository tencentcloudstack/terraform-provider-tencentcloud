# Add PostgreSQL Restore DB Instance Objects Operation Resource

## What

Add a new one-time Terraform operation resource `tencentcloud_postgresql_restore_db_instance_objects_operation` that calls the `RestoreDBInstanceObjects` API to restore database objects (databases, tables) on a PostgreSQL instance from a backup set or a point-in-time target.

## Why

Users need to restore specific database objects (e.g. individual databases or tables) from a backup without replacing the whole instance. This is a destructive, one-time operation that is naturally represented as a Terraform operation resource.

## APIs Used

| Step | API | Notes |
|---|---|---|
| Create | `RestoreDBInstanceObjects` | Async — returns `TaskId` |
| Poll | `DescribeTasks` | Poll until `TaskSet[0].Status == "Success"` |

## Inputs

| Field | Required | Description |
|---|---|---|
| `db_instance_id` | Yes | PostgreSQL instance ID; used as resource ID |
| `restore_objects` | Yes | List of object names to restore |
| `backup_set_id` | No | Backup set ID (mutually exclusive with `restore_target_time`) |
| `restore_target_time` | No | Point-in-time target (mutually exclusive with `backup_set_id`) |

## Output

- `task_id` (Computed) — Task ID returned by `RestoreDBInstanceObjects`
