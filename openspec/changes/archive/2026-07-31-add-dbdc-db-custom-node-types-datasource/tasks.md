## 1. Service Layer

- [x] 1.1 Append `DescribeDBCustomNodeTypesByFilter(ctx context.Context, param map[string]interface{}) (ret []*dbdcv20201029.DBCustomNodeTypeInfo, totalCount int64, errRet error)` to `tencentcloud/services/dbdc/service_tencentcloud_dbdc.go` — wraps `DescribeDBCustomNodeTypes`; translates `paramMap["Filters"]` to `[]*dbdcv20201029.Filter`; single SDK call with `ratelimit.Check` (NO pagination loop, NO nested retry — the caller owns the retry); guard nil response (`result == nil || result.Response == nil || result.Response.NodeTypeSet == nil`) → `NonRetryableError` with `log.Printf("[DATASOURCE] read empty, skip SetId")`; return `response.Response.NodeTypeSet`
- [x] 1.2 Verify the helper reuses the existing `DbdcService` struct and `me.client.UseDbdcV20201029Client()` accessor already present in the file

## 2. Data Source Implementation

- [x] 2.1 Create `tencentcloud/services/dbdc/data_source_tc_dbdc_db_custom_node_types.go` with `DataSourceTencentCloudDbdcDbCustomNodeTypes()` returning `*schema.Resource` (Read-only, no Create/Update/Delete)
- [x] 2.2 Define the schema: optional `filters` (TypeList of schema.Resource with `name` Required string + `values` Required TypeList of TypeString), optional `result_output_file` (TypeString), computed `node_type_set` (TypeList of schema.Resource with flattened fields `zone`, `node_type`, `node_family`, `cpu` TypeInt, `memory` TypeInt, `status`, `system_disk_types` TypeList of TypeString, `data_disk_types` TypeList of TypeString)
- [x] 2.3 Implement `dataSourceTencentCloudDbdcDbCustomNodeTypesRead`: `defer tccommon.LogElapsed(...)` + `defer tccommon.InconsistentCheck(...)`, build `paramMap` from `filters` (translate to `[]*dbdcv20201029.Filter` like `data_source_tc_dbdc_db_custom_nodes.go`), wrap service call in `resource.Retry(tccommon.ReadRetryTimeout, ...)` using `tccommon.RetryError(e)`, on failure `log.Printf("[DATASOURCE] read empty, skip SetId")` + return error
- [x] 2.4 In the Read, map each `DBCustomNodeTypeInfo` into a `map[string]interface{}` with nil guards for `Zone`, `NodeType`, `NodeFamily`, `CPU`, `Memory`, `Status`; map `SystemDiskTypes`/`DataDiskTypes` (each `[]*string`) to `[]string` lists only when non-nil; append to `nodeTypeSetList` and `_ = d.Set("node_type_set", nodeTypeSetList)`
- [x] 2.5 Set `d.SetId(helper.BuildToken())` after successful read; handle `result_output_file` via `tccommon.WriteToFile` when provided

## 3. Provider Registration

- [x] 3.1 Register `tencentcloud_dbdc_db_custom_node_types` in `tencentcloud/provider.go` DataSourcesMap (next to the existing `tencentcloud_dbdc_db_custom_images` entry), mapping to `dbdc.DataSourceTencentCloudDbdcDbCustomNodeTypes()`

## 4. Documentation

- [x] 4.1 Create `tencentcloud/services/dbdc/data_source_tc_dbdc_db_custom_node_types.md` — one-sentence description starting with "Use this data source to query ..." mentioning dbdc; Example Usage with a `filters` example (zone/node-family) and an `output` referencing `node_type_set`; NO Import section; NO Argument Reference / Attribute Reference sections (auto-generated)

## 5. Unit Tests

- [x] 5.1 Create `tencentcloud/services/dbdc/data_source_tc_dbdc_db_custom_node_types_test.go` using gomonkey mocks (NOT the Terraform test suite); mock `DescribeDBCustomNodeTypesByFilter` to return a sample `[]*dbdcv20201029.DBCustomNodeTypeInfo` and verify the Read populates `node_type_set` schema correctly (including nil `SystemDiskTypes`/`DataDiskTypes` skip behavior)
- [x] 5.2 Add a second test case covering the empty/nil response path asserting that the Read returns an error and does NOT clear the id (NonRetryableError behavior)

## 6. Verification (separate from code tasks)

- [x] 6.1 Code compiles cleanly (build verification performed by separate downstream flow, not in this phase)
- [x] 6.2 `make doc` generates website docs in the finalize phase (tfpacer-finalize skill)
