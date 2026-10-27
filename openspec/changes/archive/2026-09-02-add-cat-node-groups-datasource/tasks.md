# Implementation Tasks: Add CAT Node Groups Data Source

## 1. Service Layer

- [x] 1.1 Add `DescribeCatNodeGroupsByFilter` method to `CatService`
  - Location: `tencentcloud/services/cat/service_tencentcloud_cat.go`
  - Build a `cat.NewDescribeNodeGroupsRequest()` from the `paramMap`
  - Map filter keys to request fields: `node_type` → `NodeType` (`[]*int64`), `task_category` → `TaskCategory`, `ip_type` → `IPType`, `name` → `Name`, `region_id` → `RegionID`, `district_id` → `DistrictID`, `net_service_id` → `NetServiceID`, `node_group_type` → `NodeGroupType`, `task_type` → `TaskType`, `probe_type` → `ProbeType`
  - Call `ratelimit.Check(request.GetAction())` before the API call
  - Wrap the SDK call `me.client.UseCatClient().DescribeNodeGroups(request)` in a single `resource.Retry(tccommon.ReadRetryTimeout, ...)` block (no pagination — `DescribeNodeGroups` returns the full tree)
  - Inside the retry block, if `result == nil || result.Response == nil || (result.Response.NodeList == nil && result.Response.DistrictList == nil && result.Response.NetServiceList == nil)` return `resource.NonRetryableError(...)`
  - Return `(nodeList []*cat.NodeTree, districtList []*cat.DistinctOrNetServiceInfo, netServiceList []*cat.DistinctOrNetServiceInfo, errRet error)`
  - Add deferred error logging with `log.Printf("[CRITAL]%s api[%s] fail, ...")` matching existing cat service methods

## 2. Data Source Schema & Read

- [x] 2.1 Create `tencentcloud/services/cat/data_source_tc_cat_node_groups.go`
  - Package `cat`, imports: `context`, `log`, `tccommon`, `helper`, `schema`, `resource`, `cat` SDK
  - Implement `DataSourceTencentCloudCatNodeGroups() *schema.Resource` returning Read function + schema

- [x] 2.2 Define filter (input) schema arguments (all Optional)
  - `node_type`: `schema.TypeList`, `Elem: &schema.Schema{Type: schema.TypeInt}` (maps to `NodeType []*int64`)
  - `task_category`: `schema.TypeInt` (`TaskCategory *int64`)
  - `ip_type`: `schema.TypeInt` (`IPType *int64`)
  - `name`: `schema.TypeString` (`Name *string`)
  - `region_id`: `schema.TypeInt` (`RegionID *int64`)
  - `district_id`: `schema.TypeInt` (`DistrictID *int64`)
  - `net_service_id`: `schema.TypeInt` (`NetServiceID *int64`)
  - `node_group_type`: `schema.TypeInt` (`NodeGroupType *int64`)
  - `task_type`: `schema.TypeInt` (`TaskType *int64`)
  - `probe_type`: `schema.TypeInt` (`ProbeType *uint64`)
  - `result_output_file`: `schema.TypeString`, Optional

- [x] 2.3 Define computed output `node_list` (TypeList) nested schema
  - Top-level `NodeTree`: `id` (TypeString), `content` (TypeString), `children` (TypeList)
    - `NodeLeaf`: `id` (TypeString), `content` (TypeString), `children` (TypeList)
      - `NodeInfoBase`: `id` (TypeString), `content` (TypeString)
  - Each nested `children` uses `Elem: &schema.Resource{Schema: ...}`

- [x] 2.4 Define computed output `district_list` and `net_service_list` (TypeList) schemas
  - `district_list`: each entry `id` (TypeString), `name` (TypeString) — `DistinctOrNetServiceInfo`
  - `net_service_list`: each entry `id` (TypeString), `name` (TypeString) — `DistinctOrNetServiceInfo`

- [x] 2.5 Implement `dataSourceTencentCloudCatNodeGroupsRead` function
  - Add `defer tccommon.LogElapsed("data_source.tencentcloud_cat_node_groups.read")()` and `defer tccommon.InconsistentCheck(d, meta)()`
  - Create `logId`/`ctx` with `tccommon.GetLogId(tccommon.ContextNil)` / `context.WithValue(context.TODO(), tccommon.LogIdKey, logId)`
  - Build `paramMap` from provided filter arguments only (skip unset args); convert `node_type` list to `[]*int64` via `helper.IntInt64` per element; scalar ints via `helper.IntInt64`; string via `helper.String`
  - Call `catService.DescribeCatNodeGroupsByFilter(ctx, paramMap)` inside `resource.Retry(tccommon.ReadRetryTimeout, ...)`; on error `return tccommon.RetryError(e)`; assign returned slices
  - On retry failure path, log `log.Printf("[DATASOURCE] read empty, skip SetId")` before returning the error

- [x] 2.6 Flatten response into state with nil-pointer checks
  - Flatten `nodeList` into `nodeListTmp []map[string]interface{}`: for each `NodeTree`, add `id`/`content` only if non-nil; flatten `Children` (`NodeLeaf`) into `children`; for each `NodeLeaf` flatten `Children` (`NodeInfoBase`) into inner `children`
  - Collect top-level `NodeTree.ID` values into `ids []string`
  - Set `d.SetId(helper.DataResourceIdsHash(ids))`
  - Flatten `districtList` → `districtListTmp` (`id`, `name`), `netServiceList` → `netServiceListTmp` (`id`, `name`), with nil checks
  - `_ = d.Set("node_list", nodeListTmp)`, `_ = d.Set("district_list", districtListTmp)`, `_ = d.Set("net_service_list", netServiceListTmp)`

- [x] 2.7 Handle `result_output_file`
  - If `result_output_file` is set and non-empty, write `nodeListTmp` via `tccommon.WriteToFile(output.(string), nodeListTmp)`; return error on failure

## 3. Provider Registration

- [x] 3.1 Register data source in `tencentcloud/provider.go`
  - Add `"tencentcloud_cat_node_groups": cat.DataSourceTencentCloudCatNodeGroups()` to `DataSourcesMap` (alphabetically after `tencentcloud_cat_node` and before `tencentcloud_cat_probe_data`)

- [x] 3.2 Update provider.go doc-comment index
  - Add the new data source entry under the CAT product's `Data Source` subsection in the provider.go comment block so `make doc` updates `website/tencentcloud.erb`

## 4. Documentation Example File

- [x] 4.1 Create `tencentcloud/services/cat/data_source_tc_cat_node_groups.md`
  - One-line description mentioning CAT: `Use this data source to query detailed information of cat node groups`
  - `Example Usage` HCL block demonstrating filter usage (e.g. `node_type = [1]`, `ip_type = 1`)
  - No `Argument Reference` / `Attribute Reference` sections (auto-generated by `make doc`)

## 5. Unit Tests (gomonkey mock)

- [x] 5.1 Create `tencentcloud/services/cat/data_source_tc_cat_node_groups_test.go`
  - Package `cat_test`, imports: `testing`, `gomonkey/v2`, `schema`, `stretchr/testify/assert`, `cat` SDK, `tccommon`, `connectivity`, `services/cat`
  - Implement a mock `ProviderMeta` returning `*connectivity.TencentCloudClient`

- [x] 5.2 Test successful read flattening
  - Mock `UseCatClient` to return a `cat.Client{}`; mock `DescribeNodeGroups` to return a response with one `NodeTree` (one `NodeLeaf` child with one `NodeInfoBase` grandchild), one district, one ISP
  - Call `res.Read(d, meta)`; assert no error, non-empty `d.Id()`, and `node_list`/`district_list`/`net_service_list` values with nested `children` populated

- [x] 5.3 Test nil/empty response path
  - Mock `DescribeNodeGroups` returning nil response → assert `res.Read` returns an error (NonRetryableError surfaces) and `d.SetId("")` is NOT called prematurely

- [x] 5.4 Test schema correctness
  - Assert `res.Schema` contains `node_list`, `district_list`, `net_service_list`, `result_output_file`, and all filter arguments with correct types

## 6. Validation

- [x] 6.1 Verify code compiles (gofmt in finalize phase only)
- [x] 6.2 Verify `make doc` regenerates `website/docs/d/cat_node_groups.html.markdown` and updates `website/tencentcloud.erb`
- [x] 6.3 Confirm no existing CAT data sources are broken (additive change only)
