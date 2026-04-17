## 1. Resource Implementation

- [x] 1.1 Create `resource_tc_postgresql_restore_db_instance_objects_operation.go`:
  - Schema: `db_instance_id` (Required/ForceNew), `restore_objects` (Required/ForceNew, List of String), `backup_set_id` (Optional/ForceNew), `restore_target_time` (Optional/ForceNew), `task_id` (Computed)
  - Create handler: call `RestoreDBInstanceObjects`, poll `DescribeTasks` until `Status == "Success"`
  - Read handler: no-op (set `db_instance_id` from ID)
  - Update handler: not needed (all ForceNew)
  - Delete handler: no-op

## 2. Provider Registration

- [x] 2.1 Register `tencentcloud_postgresql_restore_db_instance_objects_operation` in `provider.go`

## 3. Documentation & Tests

- [x] 3.1 Create `resource_tc_postgresql_restore_db_instance_objects_operation.md`
- [x] 3.2 Create `resource_tc_postgresql_restore_db_instance_objects_operation_test.go`

## 4. Refinements

- [x] 4.1 Remove `task_id` field from schema and create handler
- [x] 4.2 Simplify Read handler to empty body (no-op, no `d.Set` calls)
