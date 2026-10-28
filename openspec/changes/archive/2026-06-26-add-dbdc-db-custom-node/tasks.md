## 1. Service Layer

- [x] 1.1 Append `DescribeDBCustomNodeById(ctx, nodeId)` to `service_tencentcloud_dbdc.go` — wraps `DescribeDBCustomNodes` with `NodeIds=[nodeId]`, `Limit=100`, `resource.Retry` + `ratelimit.Check`; returns single `*dbdcv20201029.DBCustomNode` (nil if not found), nil/length-safe
- [x] 1.2 Reuse existing `DescribeDBCustomTaskStatusById` and `waitDBCustomTaskSucceeded` (already in the `dbdc` package)

## 2. Resource Implementation

- [x] 2.1 Create `resource_tc_dbdc_db_custom_node.go` with full schema (arguments aligned to `CreateDBCustomNodesRequest` minus `client_token`; `login_settings` nested block; `tags` map; computed read-only fields incl. `system_disk`/`data_disks`)
- [x] 2.2 Implement Create: build request (`NodeCount` default 1), call `CreateDBCustomNodes` (retry), poll task until `Succeeded`, `SetId(NodeIds[0])`
- [x] 2.3 Implement Read: call `DescribeDBCustomNodeById`, populate all fields with nil guards, clear ID when not found
- [x] 2.4 Implement Update: on `tags` change call `ModifyDBCustomNodeTags` (AddTags/DeleteTagKeys); on `auto_renew` change call `RenewDBCustomNode` (retry)
- [x] 2.5 Implement Delete: `IsolateDBCustomNode` (retry) -> poll node status until `Isolated` -> `DestroyDBCustomNode` (retry) -> poll task until `Succeeded`
- [x] 2.6 Add `waitDBCustomNodeStatus` helper for the isolate stage

## 3. Provider Registration

- [x] 3.1 Register `tencentcloud_dbdc_db_custom_node` in `provider.go` ResourcesMap

## 4. Documentation & Tests

- [x] 4.1 Create `resource_tc_dbdc_db_custom_node.md` (Example Usage + Import, per `resource_tc_config_compliance_pack.md` convention)
- [x] 4.2 Create `resource_tc_dbdc_db_custom_node_test.go` (per `resource_tc_config_compliance_pack_test.go` convention)
- [x] 4.3 Add website doc `website/docs/r/dbdc_db_custom_node.html.markdown` (+ erb nav entry)

## 5. Verification

- [x] 5.1 `go build ./tencentcloud/...` compiles cleanly
- [x] 5.2 `go vet ./tencentcloud/services/dbdc/...` passes
