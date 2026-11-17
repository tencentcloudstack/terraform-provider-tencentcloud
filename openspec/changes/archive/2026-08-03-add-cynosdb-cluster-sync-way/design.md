## Context

The `tencentcloud_cynosdb_cluster` resource (defined in `tencentcloud/services/cynosdb/resource_tc_cynosdb_cluster.go`, schema in `TencentCynosdbClusterBaseInfo()` from `extension_cynosdb.go`) currently supports creating a cluster with `CreateClusters`, but does not expose the cluster synchronization mode parameters.

**Current state:**
- Resource file: `tencentcloud/services/cynosdb/resource_tc_cynosdb_cluster.go`
- Schema: `TencentCynosdbClusterBaseInfo()` in `tencentcloud/services/cynosdb/extension_cynosdb.go`
- Service layer: `tencentcloud/services/cynosdb/service_tencentcloud_cynosdb.go` (`DescribeClusterById` returns `*cynosdb.CynosdbClusterDetail`)
- SDK: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107` (already vendored, includes the required fields)

**API behavior analysis (verified against vendored SDK):**

| API | SyncWay in Request | SemiSyncTimeout in Request | BinlogSyncWay in Response | SemiSyncTimeout in Response |
|-----|--------------------|---------------------------|--------------------------|-----------------------------|
| `CreateClusters` | Yes (`SyncWay *string`, async/semisync/sync) | Yes (`SemiSyncTimeout *int64`, [1000, 4294967295] ms, default 10000) | No | No |
| `DescribeClusterDetail` | N/A | N/A | Yes (`Detail.SlaveZoneAttr[].BinlogSyncWay *string`) | Yes (`Detail.SlaveZoneAttr[].SemiSyncTimeout *int64`) |

The requirement's other mappings (`ModifyInsGrpSecurityGroups` -> `rw_group_sg`/`ro_group_sg`, `SwitchServerlessCluster` -> `serverless_status_flag`, `ModifyTags` -> `tags`) are already implemented in the resource and are not part of this change's implementation work.

**Key constraint:** The write path for `SyncWay`/`SemiSyncTimeout` is only `CreateClusters` per the requirement; the read path is `DescribeClusterDetail` (`Detail.SlaveZoneAttr[0]`). No update API is mapped for these fields, so they must be immutable after creation.

## Goals / Non-Goals

**Goals:**
- Add `SyncWay` (Optional, string, immutable after creation) parameter to `tencentcloud_cynosdb_cluster` with valid values `async`, `semisync`, `sync`
- Add `SemiSyncTimeout` (Optional, int, immutable after creation) parameter to `tencentcloud_cynosdb_cluster` with valid range `[1000, 4294967295]` ms
- Pass `SyncWay` and `SemiSyncTimeout` to the `CreateClusters` API request when specified by the user
- Read `BinlogSyncWay` and `SemiSyncTimeout` from the `DescribeClusterDetail` API response (`Detail.SlaveZoneAttr[0]`) in the Read function to support state refresh and import
- Implement immutable-args handling in the Update function so changes to these fields return a clear error (consistent with the existing `db_mode`/`min_cpu`/`max_cpu` pattern)
- Maintain full backward compatibility — existing configurations continue to work unchanged

**Non-Goals:**
- Making `SyncWay`/`SemiSyncTimeout`/`BinlogSyncWay` updatable after creation (no update API is mapped in the requirement)
- Passing `BinlogSyncWay`/`SemiSyncTimeout` through `ModifyClusterSlaveZone`/`AddClusterSlaveZone` update paths (not part of the mapped requirement; noted as a future extension)
- Adding these fields to any datasource or to `tencentcloud_cynosdb_cluster_v2` (out of scope)

## Decisions

### Decision 1: Schema field naming follows the requirement contract (`SyncWay`, `SemiSyncTimeout`, `BinlogSyncWay`)

**Rationale:** The requirement explicitly maps the target Terraform schema names as `SyncWay`, `SemiSyncTimeout`, and `BinlogSyncWay`. These map 1:1 to the API field names, keeping the mapping unambiguous. `SyncWay` and `SemiSyncTimeout` are user-configurable (Optional); `BinlogSyncWay` is only present in the Describe response, so it is a Computed field refreshed from the API.

### Decision 2: `SyncWay` and `SemiSyncTimeout` are immutable (not ForceNew)

**Rationale:** The mapped write API is `CreateClusters` only; there is no mapped update API for these fields. Following the existing pattern in this resource (e.g., `db_mode`, `min_cpu`, `max_cpu`, `auto_pause`, `auto_pause_delay`, `storage_pay_mode` use an `immutableArgs` array in the Update function), we add these fields to the immutable-args check so that changing them after creation returns a clear error instead of silently destroying/recreating the resource.

### Decision 3: Read `BinlogSyncWay` and `SemiSyncTimeout` from `DescribeClusterDetail.Detail.SlaveZoneAttr`

**Rationale:** `DescribeClusterById` already returns the full `*cynosdb.CynosdbClusterDetail`, which includes `SlaveZoneAttr []*SlaveZoneAttrItem` with `BinlogSyncWay` and `SemiSyncTimeout`. The Read function reads the first slave-zone attribute entry (`SlaveZoneAttr[0]`, guarded by a nil/empty check) and sets these values, consistent with how the resource already reads `cluster.SlaveZones[0]` for `slave_zone`. Guarding with nil checks prevents panics when the cluster has no slave zone.

### Decision 4: Validation of `SyncWay` and `SemiSyncTimeout`

**Rationale:** `SyncWay` validates against allowed values `async`/`semisync`/`sync` using `tccommon.ValidateAllowedStringValue`; `SemiSyncTimeout` validates the `[1000, 4294967295]` range using `tccommon.ValidateIntegerInRange`. This mirrors the existing schema validation style in `extension_cynosdb.go` and gives users clearer errors than API-side failures.

### Decision 5: No SDK update required

**Rationale:** The vendored SDK (`vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107`) already contains `CreateClustersRequest.SyncWay`/`SemiSyncTimeout` and `SlaveZoneAttrItem.BinlogSyncWay`/`SemiSyncTimeout`, so no dependency bump is needed for this change.

## Risks / Trade-offs

- **[Risk] `SlaveZoneAttr` is empty when the cluster has no slave zone**: The Describe response may return an empty `SlaveZoneAttr` list.
  - **Mitigation:** Read function checks `cluster.SlaveZoneAttr != nil && len(cluster.SlaveZoneAttr) > 0` before reading index 0, and skips setting when empty.

- **[Risk] `SemiSyncTimeout` mismatch between create-time value and read-back value**: If the API normalizes the timeout (e.g., clamps to defaults), the state could drift.
  - **Mitigation:** `SemiSyncTimeout` is read back from `DescribeClusterDetail` and may need `Computed` semantics alongside `Optional` to absorb normalization differences. `SyncWay`/`BinlogSyncWay` are strings with stable values.

- **[Risk] Changing `SyncWay`/`SemiSyncTimeout` after creation errors out**: Users may expect in-place update.
  - **Mitigation:** Clear immutable-args error message (`argument 'X' cannot be modified`), consistent with existing behavior for other immutable fields. A future change can map `ModifyClusterSlaveZone`/`AddClusterSlaveZone` to support updates.

- **[Risk] CamelCase schema keys deviate from provider snake_case convention**: The requirement explicitly specifies `SyncWay`/`SemiSyncTimeout`/`BinlogSyncWay`.
  - **Mitigation:** These names are kept as the requirement contract; terraform SDK v2 accepts alphanumeric schema keys. Documentation generation follows the schema keys.
