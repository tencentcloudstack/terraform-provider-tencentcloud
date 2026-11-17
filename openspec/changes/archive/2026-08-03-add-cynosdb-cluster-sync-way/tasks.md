## 1. Schema Definition

- [x] 1.1 Add `SyncWay` schema field (TypeString, Optional, Computed, validate `async`/`semisync`/`sync` via `tccommon.ValidateAllowedStringValue`) to `TencentCynosdbClusterBaseInfo()` in `tencentcloud/services/cynosdb/extension_cynosdb.go`
- [x] 1.2 Add `SemiSyncTimeout` schema field (TypeInt, Optional, Computed, validate `[1000, 4294967295]` via `tccommon.ValidateIntegerInRange`) to `TencentCynosdbClusterBaseInfo()`
- [x] 1.3 Add `BinlogSyncWay` schema field (TypeString, Computed) to `TencentCynosdbClusterBaseInfo()`

## 2. Create Function Changes

- [x] 2.1 Read `SyncWay` and `SemiSyncTimeout` from schema data in `resourceTencentCloudCynosdbClusterCreate`
- [x] 2.2 Set `request.SyncWay` and `request.SemiSyncTimeout` on the `CreateClustersRequest` when the corresponding schema values are present

## 3. Read Function Changes

- [x] 3.1 Read `BinlogSyncWay` and `SemiSyncTimeout` from `DescribeClusterDetail` response (`cluster.SlaveZoneAttr[0]`) in `resourceTencentCloudCynosdbClusterRead`, guarded by a nil/empty check on `SlaveZoneAttr`
- [x] 3.2 Set `BinlogSyncWay` and `SemiSyncTimeout` in state via `d.Set`

## 4. Update Function Changes

- [x] 4.1 Add `SyncWay` and `SemiSyncTimeout` to the `immutableArgs` array in `resourceTencentCloudCynosdbClusterUpdate` so changes after creation return an error

## 5. Unit Tests

- [x] 5.1 Add unit test cases for the `SyncWay`/`SemiSyncTimeout` create request wiring using gomonkey mocks of the cloud API (business-logic-only tests, per resource kind convention)
- [x] 5.2 Add unit test cases for reading `BinlogSyncWay`/`SemiSyncTimeout` from `DescribeClusterDetail` response in the Read function

## 6. Documentation

- [x] 6.1 Update `tencentcloud/services/cynosdb/resource_tc_cynosdb_cluster.md` with usage examples for `SyncWay`, `SemiSyncTimeout`, and `BinlogSyncWay`

## 7. Validation

- [x] 7.1 Verify the code compiles successfully
- [x] 7.2 Verify no lint errors
