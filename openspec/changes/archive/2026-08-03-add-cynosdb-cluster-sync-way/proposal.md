## Why

The CynosDB `CreateClusters` API supports `SyncWay` (async/semisync/sync) and `SemiSyncTimeout` (ms) for cluster synchronization configuration, and `DescribeClusterDetail` returns `SlaveZoneAttr.BinlogSyncWay` / `SlaveZoneAttr.SemiSyncTimeout` for the current slave-zone synchronization settings. The Terraform resource `tencentcloud_cynosdb_cluster` does not expose these parameters, so users cannot configure or read the cluster synchronization mode through Terraform and must fall back to the console or raw API calls.

## What Changes

- Add `SyncWay` (Optional, string, ForceNew/immutable) parameter to `tencentcloud_cynosdb_cluster` resource, passed to the `CreateClusters` API as `SyncWay`. Valid values: `async`, `semisync`, `sync`.
- Add `SemiSyncTimeout` (Optional, int, ForceNew/immutable) parameter to `tencentcloud_cynosdb_cluster` resource, passed to the `CreateClusters` API as `SemiSyncTimeout`. Valid range: `[1000, 4294967295]` ms, API default `10000`.
- Read `BinlogSyncWay` and `SemiSyncTimeout` from the `DescribeClusterDetail` API response (`Detail.SlaveZoneAttr[0]`) in the Read function, so imported resources and state refreshes populate these values.
- The existing mappings (`ModifyInsGrpSecurityGroups` -> `rw_group_sg`/`ro_group_sg`, `SwitchServerlessCluster` -> `serverless_status_flag`, `ModifyTags` -> `tags`) are already implemented in the resource and remain unchanged by this proposal.

## Capabilities

### New Capabilities
- `cynosdb-cluster-sync-way`: Enable the `SyncWay`, `SemiSyncTimeout` and `BinlogSyncWay` parameters on the `tencentcloud_cynosdb_cluster` resource so users can specify and read the cluster synchronization mode (async/semisync/sync) and the semi-sync timeout.

### Modified Capabilities
<!-- No existing specs require modification -->

## Impact

- **Affected files:**
  - `tencentcloud/services/cynosdb/extension_cynosdb.go` — add `SyncWay`, `SemiSyncTimeout`, `BinlogSyncWay` schema fields to `TencentCynosdbClusterBaseInfo()` (used by `ResourceTencentCloudCynosdbCluster`)
  - `tencentcloud/services/cynosdb/resource_tc_cynosdb_cluster.go` — wire `SyncWay`/`SemiSyncTimeout` into `CreateClustersRequest` in Create; read `BinlogSyncWay`/`SemiSyncTimeout` from `DescribeClusterDetail` (`Detail.SlaveZoneAttr`) in Read; add immutable-args handling in Update since no update API is mapped for these fields
  - `tencentcloud/services/cynosdb/resource_tc_cynosdb_cluster.md` — update documentation with new usage example
- **SDK dependency:** `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107` — already vendored with `CreateClustersRequest.SyncWay`/`SemiSyncTimeout` and `CynosdbClusterDetail.SlaveZoneAttr[].BinlogSyncWay`/`SemiSyncTimeout` fields; no SDK update needed.
- **Backward compatibility:** fully backward compatible — the new parameters are Optional and default to not being set (API defaults apply).
- **API constraints:** `SyncWay` and `SemiSyncTimeout` are only accepted by `CreateClusters` in the mapped APIs; the requirement maps only Create (write) and `DescribeClusterDetail` (read). No update API is mapped, so these fields must be immutable after creation (rejected via immutable-args check in Update).
