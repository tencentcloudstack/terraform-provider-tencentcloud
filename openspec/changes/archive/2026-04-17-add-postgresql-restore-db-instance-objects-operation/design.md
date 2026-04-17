# Design: PostgreSQL Restore DB Instance Objects Operation Resource

## File Layout

| File | Action |
|---|---|
| `tencentcloud/services/postgresql/resource_tc_postgresql_restore_db_instance_objects_operation.go` | New |
| `tencentcloud/services/postgresql/resource_tc_postgresql_restore_db_instance_objects_operation.md` | New |
| `tencentcloud/services/postgresql/resource_tc_postgresql_restore_db_instance_objects_operation_test.go` | New |
| `tencentcloud/provider.go` | Modified (register resource) |

## Code Style Reference

Strictly follow `resource_tc_igtm_strategy.go` style:
- `var (logId, ctx, request, response)` block at top of each handler
- `resource.Retry(tccommon.WriteRetryTimeout, ...)` wrapping API call
- nil check on response before proceeding
- `log.Printf("[DEBUG]...")` after success
- `log.Printf("[CRITAL]...")` on error

## Schema

```
db_instance_id        Required, ForceNew, String   — Instance ID; SetId to this value
restore_objects       Required, ForceNew, List(String) — Objects to restore
backup_set_id         Optional, ForceNew, String   — Backup set ID
restore_target_time   Optional, ForceNew, String   — Point-in-time target
task_id               Computed, Int                — Task ID from API response
```

## Create Handler

1. Build `RestoreDBInstanceObjectsRequest` from schema fields
2. Call `RestoreDBInstanceObjectsWithContext` wrapped in `resource.Retry(WriteRetryTimeout)`
3. Check `response.Response.TaskId != nil`
4. Call `DescribeTasksWithContext` wrapped in `resource.Retry(WriteRetryTimeout)` polling until `TaskSet[0].Status == "Success"`; treat `Fail` / `Pause` as `NonRetryableError`
5. `d.SetId(dbInstanceId)`
6. Call Read handler

## Read Handler

One-time resource: just set `db_instance_id` from `d.Id()` and `task_id` if known. No remote read needed.

## Update Handler

Not applicable (all fields are ForceNew).

## Delete Handler

No-op.

## Async Polling Pattern

```go
flowRequest := postgresqlv20170312.NewDescribeTasksRequest()
flowRequest.TaskId = helper.Int64Uint64(taskId)
err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
    result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePostgresqlV20170312Client().DescribeTasksWithContext(ctx, flowRequest)
    if e != nil {
        return tccommon.RetryError(e)
    }
    if result == nil || result.Response == nil || result.Response.TaskSet == nil {
        return resource.NonRetryableError(fmt.Errorf("DescribeTasks response is nil"))
    }
    if len(result.Response.TaskSet) == 0 {
        return resource.RetryableError(fmt.Errorf("waiting for task to initialize"))
    }
    status := *result.Response.TaskSet[0].Status
    if status == "Success" {
        return nil
    }
    if status == "Fail" || status == "Pause" {
        return resource.NonRetryableError(fmt.Errorf("task failed with status: %s", status))
    }
    return resource.RetryableError(fmt.Errorf("task still running, status: %s", status))
})
```
