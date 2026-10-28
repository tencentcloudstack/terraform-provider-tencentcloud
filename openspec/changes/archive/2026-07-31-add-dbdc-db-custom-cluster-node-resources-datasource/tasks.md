## 1. Service Layer

- [x] 1.1 Append `DescribeDBCustomClusterNodeResourcesByFilter(ctx context.Context, param map[string]interface{}) (ret []*dbdcv20201029.DBCustomClusterNodeResource, errRet error)` to `tencentcloud/services/dbdc/service_tencentcloud_dbdc.go` — wraps the synchronous SDK `DescribeDBCustomClusterNodeResources`; maps `ClusterId` (`*string`) and `NodeIds` (`[]*string`) from param map; single call inside `resource.Retry(tccommon.ReadRetryTimeout, ...)` with `ratelimit.Check` (no pagination loop — API has no Offset/Limit); inside retry returns `resource.NonRetryableError(fmt.Errorf("Describe dbdc_db_custom_cluster_node_resources failed, Response is nil."))` when `result == nil || result.Response == nil || result.Response.NodeSet == nil`; on retry failure path logs `log.Printf("[DATASOURCE] read empty, skip SetId")`; returns `response.Response.NodeSet` on success

## 2. Data Source Implementation

- [x] 2.1 Create `tencentcloud/services/dbdc/data_source_tc_dbdc_db_custom_cluster_node_resources.go` with schema: `cluster_id` (Required String → `request.ClusterId`), `node_ids` (Optional List of String → `request.NodeIds`, schema description notes max 50 per request), `result_output_file` (Optional String), `node_set` (Computed List of Resource). The `node_set` element schema: `node_id` (Computed String), and five nested `MetaResource` blocks `capacity`/`allocatable`/`requests`/`limits`/`available` (each TypeList MaxItems 1, Computed) with nested `cpu` (Computed Float), `memory` (Computed Float), `pods` (Computed Int)
- [x] 2.2 Implement `dataSourceTencentCloudDbdcDbCustomClusterNodeResourcesRead`: defer `tccommon.LogElapsed` + `tccommon.InconsistentCheck`; build paramMap from `cluster_id` and `node_ids`; call service helper inside `resource.Retry(tccommon.ReadRetryTimeout)`; on error return it (do NOT `d.SetId("")` on empty — service layer already returns NonRetryableError); build `nodeSetList` with nil guards for every pointer (each `MetaResource` block nil-checked, then its `Cpu`/`Memory`/`Pods` nil-checked before set); `_ = d.Set("node_set", nodeSetList)`; `d.SetId(helper.BuildToken())`; write `result_output_file` via `tccommon.WriteToFile` if set
- [x] 2.3 No `#`-style composite ID or import — data source uses `helper.BuildToken()` for its id, mirroring `tencentcloud_dbdc_db_custom_cluster_nodes`

## 3. Provider Registration

- [x] 3.1 Register `"tencentcloud_dbdc_db_custom_cluster_node_resources": dbdc.DataSourceTencentCloudDbdcDbCustomClusterNodeResources()` in `tencentcloud/provider.go` DataSourcesMap (next to the existing `tencentcloud_dbdc_db_custom_cluster_nodes` entry)
- [x] 3.2 Add `tencentcloud_dbdc_db_custom_cluster_node_resources` entry to `tencentcloud/provider.md` in the DBDC Data Source section

## 4. Documentation & Tests

- [x] 4.1 Create `tencentcloud/services/dbdc/data_source_tc_dbdc_db_custom_cluster_node_resources.md` — one-line description "Use this data source to query detailed information of DB Custom cluster node resources" + Example Usage (query by `cluster_id`, and query by `cluster_id` + `node_ids`); no Argument/Attribute Reference sections
- [x] 4.2 Create `tencentcloud/services/dbdc/data_source_tc_dbdc_db_custom_cluster_node_resources_test.go` with gomonkey mock unit tests (no TF acceptance suite): `TestDbdcDbCustomClusterNodeResourcesDS_ReadBasic` (mock `DescribeDBCustomClusterNodeResources` returning a 2-element `NodeSet` with populated `MetaResource` blocks; assert `node_set` mapping incl. nested capacity/allocatable/requests/limits/available cpu/memory/pods; assert `d.Id()` non-empty), `TestDbdcDbCustomClusterNodeResourcesDS_Schema` (assert schema fields/types incl. nested blocks), `TestDbdcDbCustomClusterNodeResourcesDS_ReadWithEmptyResponse` (mock returns nil/empty `NodeSet`; assert Read errors)

## 5. Verification

- [x] 5.1 Ensure all code compiles cleanly (no `go build`/`go vet` to be run in this stage per project rules; verify by inspection that SDK types/fields referenced exist in vendored `dbdc/v20201029` models.go)
- [x] 5.2 Ensure unit test file is syntactically correct and mocks the correct SDK method `DescribeDBCustomClusterNodeResources` with the right request/response types
