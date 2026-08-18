## Context

The Terraform Provider for TencentCloud manages CDB (MySQL) instances through the `tencentcloud_mysql_instance`, `tencentcloud_mysql_readonly_instance`, and `tencentcloud_mysql_dr_instance` resources. However, there is no resource to manage **clone** (回档/克隆) instances — i.e., creating a new CDB instance by cloning an existing source instance, optionally rolling back to a specified point-in-time or a specified backup set.

The cloud API provides four interfaces for this capability in the `cdb/v20170320` SDK package:
- `CreateCloneInstance` — creates a clone instance; **async**, returns `AsyncRequestId`
- `DescribeDBInstances` — queries instance details (used as the Read interface, returns `InstanceInfo` items including memory/volume/cpu/zone/status/etc.)
- `UpgradeDBInstance` — adjusts instance configuration (memory/volume/cpu); **async**, returns `AsyncRequestId` and `DealIds`
- `OfflineIsolatedInstances` — offlines (destroys) an isolated instance

Since `CreateCloneInstance` is asynchronous and only returns an `AsyncRequestId` (not the cloned instance ID directly), additional handling is required to obtain the actual cloned instance ID. The SDK also provides:
- `DescribeAsyncRequestInfo` — polls async task status (`INITIAL`/`RUNNING`/`SUCCESS`/`FAILED`/`KILLED`/...)
- `DescribeCloneList` — queries clone task list by source `InstanceId`, returning `CloneItem` items that contain `DstInstanceId` (the newly created cloned instance ID)

This is a RESOURCE_KIND_GENERAL resource — it manages the full CRUD lifecycle of a cloned CDB instance.

## Goals / Non-Goals

**Goals:**
- Add a new Terraform resource `tencentcloud_cdb_clone_instance` to manage the full lifecycle of cloned CDB instances
- Support creating a clone from a source instance with optional rollback time or backup set
- Support async operation polling after Create and Update using `DescribeAsyncRequestInfo`
- Support reading instance configuration from `DescribeDBInstances`
- Support updating instance configuration via `UpgradeDBInstance` (memory, volume, cpu, protect_mode, deploy_mode, slave_zone, backup_zone, device_type, cluster_topology, fourth_zone)
- Support deletion via `OfflineIsolatedInstances`
- Support Import for existing cloned instances
- Implement unit tests using gomonkey mock approach (no terraform test suite)

**Non-Goals:**
- Managing the source instance (that remains under `tencentcloud_mysql_instance`)
- Modifying existing CDB resources or their schemas
- Adding a datasource for clone instances (the `DescribeCloneList`-based datasource `tencentcloud_mysql_clone_list` already exists)
- Supporting engine version upgrade through this resource (use `UpgradeDBInstanceEngineVersion` separately if needed — not in scope)

## Decisions

1. **Resource ID Strategy**: The resource ID is the **cloned (destination) instance ID** (`DstInstanceId`), obtained from `DescribeCloneList` after the async `CreateCloneInstance` task succeeds. This is the ID of the newly created instance, used for subsequent Read/Update/Delete operations.

2. **Async Create Flow**: `CreateCloneInstance` returns `AsyncRequestId` but NOT the cloned instance ID. The Create operation SHALL:
   1. Call `CreateCloneInstance` with retry using `tccommon.WriteRetryTimeout`; validate response and `AsyncRequestId` are non-empty (return `NonRetryableError` if empty)
   2. Poll `DescribeAsyncRequestInfo` with the `AsyncRequestId` until status is `SUCCESS` (return `RetryableError` for `INITIAL`/`RUNNING`, `NonRetryableError` for `FAILED`/`KILLED`)
   3. After async success, call `DescribeCloneList` with the source `InstanceId` to find the `CloneItem` whose `CloneJobId` or most recent `DstInstanceId` corresponds to this clone; extract `DstInstanceId`
   4. Set the resource ID to `DstInstanceId` and call Read to populate state
   The polling uses `d.Timeout(schema.TimeoutCreate)`.

3. **Async Update Flow**: `UpgradeDBInstance` returns `AsyncRequestId`. The Update operation SHALL:
   1. Call `UpgradeDBInstance` with retry using `tccommon.WriteRetryTimeout`; validate response and `AsyncRequestId` are non-empty
   2. Poll `DescribeAsyncRequestInfo` until status is `SUCCESS`
   3. Call Read to refresh state
   The polling uses `d.Timeout(schema.TimeoutUpdate)`.

4. **Read Operation**: The Read operation calls `DescribeDBInstances` with `InstanceIds` set to `d.Id()` (the cloned instance ID). Following the existing `DescribeDBInstanceById` pattern, `QueryClusterInfo` is set to `true`. If the response items are empty, log `[CRUD] cdb_clone_instance id=<id>` before calling `d.SetId("")`. Fields are set only when the response field is non-nil.

5. **Delete Operation**: The delete operation follows the existing pattern in `resource_tc_mysql_dr_instance.go`: first isolate the instance via `IsolateDBInstance`, poll `DescribeDBInstanceById` until status reaches `MYSQL_STATUS_ISOLATED`, then call `OfflineIsolatedInstances`. However, the user-specified delete API is `OfflineIsolatedInstances` only. To handle a freshly cloned instance (which may already be in a deletable state or require isolation first), the Delete operation SHALL: isolate via `IsolateDBInstance` (retry), poll until isolated, then call `OfflineIsolatedInstances`, then poll `DescribeIsolatedDBInstanceById` until the instance disappears. This reuses the proven delete flow from the existing CDB instance resources to ensure robustness. If isolation is unnecessary, the retry will still converge to the isolated state.

6. **Schema Design**: Top-level fields map to the Create/Read/Update API parameters. Fields only used at Create (and not updatable) are NOT marked ForceNew unless the cloud API lacks an Update path. Since this is a GENERAL resource with a real Update operation, only `instance_id` (source) is effectively create-only — but the resource ID is the cloned instance ID, not the source. The source `instance_id`, `specified_rollback_time`, `specified_backup_id` are create-only parameters (ForceNew) because `UpgradeDBInstance` does not accept them. Fields accepted by `UpgradeDBInstance` (memory, volume, cpu, protect_mode, deploy_mode, slave_zone, backup_zone, device_type, cluster_topology, fourth_zone) are updatable. Create-only fields like `uniq_vpc_id`, `uniq_subnet_id`, `instance_name`, `security_group`, `resource_tags`, `cage_id`, `project_id`, `pay_type`, `period`, `zone`, `master_zone`, `slave_zone`(create), `instance_nodes`, `deploy_group_id`, `dry_run`, `src_region`, `specified_sub_backup_id` — these are Create-only and marked ForceNew since `UpgradeDBInstance` does not accept them.

7. **Field Type Mapping**: 
   - `memory`, `volume`, `cpu`, `protect_mode`, `deploy_mode`, `instance_nodes`, `period`, `specified_backup_id`, `specified_sub_backup_id` → TypeInt
   - `instance_id`, `specified_rollback_time`, `uniq_vpc_id`, `uniq_subnet_id`, `instance_name`, `slave_zone`, `backup_zone`, `device_type`, `deploy_group_id`, `cage_id`, `pay_type`, `src_region`, `master_zone`, `zone`, `fourth_zone`, `async_request_id` → TypeString
   - `security_group` → TypeList of TypeString
   - `resource_tags` → TypeList of tag block (`key`, `value`)
   - `cluster_topology` → TypeList (MaxItems:1) nested block (per `ClusterTopology` struct)
   - `dry_run` → TypeBool
   - `project_id` → TypeInt
   - `async_request_id` → Computed TypeString

8. **Field Name Conflicts**: The cloned instance's own ID is the resource ID (not exposed as a schema field to avoid confusion with the source `instance_id`). The source instance ID is the schema field `instance_id`. The `async_request_id` is a Computed field populated from the Create/Update async response.

9. **immutableArgs for Update**: Since this is a GENERAL resource with a real Update, the Update method checks which fields changed and only sends updatable parameters to `UpgradeDBInstance`. Create-only fields that changed trigger re-creation (ForceNew), so they won't appear in Update. However, fields that are neither accepted by `UpgradeDBInstance` nor marked ForceNew must be added to `immutableArgs` and rejected with an error if changed.

## Risks / Trade-offs

- [Risk] `CreateCloneInstance` returns only `AsyncRequestId`, requiring a two-step lookup (poll async status, then `DescribeCloneList`) to obtain the cloned instance ID → Mitigation: Use `DescribeAsyncRequestInfo` to wait for `SUCCESS`, then query `DescribeCloneList` with the source `InstanceId` and pick the latest `DstInstanceId`; if multiple clones exist for the same source, match by `CloneJobId` proximity to the async request, or pick the most recent successful clone. The polling uses the Create timeout.
- [Risk] Async clone task may fail or be killed → Mitigation: Return `NonRetryableError` when `DescribeAsyncRequestInfo` status is `FAILED`/`KILLED`/`REMOVED`, surfacing the task `Info` message.
- [Risk] `OfflineIsolatedInstances` requires the instance to be in isolated state first → Mitigation: Follow the proven isolate-then-offline pattern from existing CDB instance resources (`IsolateDBInstance` → poll until isolated → `OfflineIsolatedInstances` → poll until gone).
- [Risk] Create-only fields changed by the user would be rejected → Mitigation: Mark truly create-only fields as `ForceNew` so Terraform handles re-creation; only fields that are neither updatable nor ForceNew go into `immutableArgs`.
- [Trade-off] The resource does not expose the source instance's region parameter `src_region` for update (it is Create-only) → acceptable since cross-region clone target is fixed at creation.
