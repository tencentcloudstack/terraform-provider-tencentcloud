## Context

The TencentCloud Terraform provider serves the CAT (Cloud Automated Testing / 拨测) product via the SDK package `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cat/v20180409`. Three CAT data sources already exist (`tencentcloud_cat_node`, `tencentcloud_cat_probe_data`, `tencentcloud_cat_metric_data`), implemented in `tencentcloud/services/cat/`.

This change adds a fourth data source, `tencentcloud_cat_node_groups`, backed by the cloud API `DescribeNodeGroups` (获取拨测点组). `DescribeNodeGroups` returns a two-level tree of probe node groups (`NodeList []*NodeTree`), plus reference lists for districts (`DistrictList []*DistinctOrNetServiceInfo`) and ISPs (`NetServiceList []*DistinctOrNetServiceInfo`). The request accepts optional filters (node type, task category, IP type, name, region id, district id, net service id, node group type, task type, probe type).

Key SDK types (already vendored, verified):
- `DescribeNodeGroupsRequest`: scalar/pointer fields — `NodeType []*int64`, `TaskCategory *int64`, `IPType *int64`, `Name *string`, `RegionID *int64`, `DistrictID *int64`, `NetServiceID *int64`, `NodeGroupType *int64`, `TaskType *int64`, `ProbeType *uint64`.
- `DescribeNodeGroupsResponse.Response`: `NodeList []*NodeTree`, `DistrictList []*DistinctOrNetServiceInfo`, `NetServiceList []*DistinctOrNetServiceInfo`.
- `NodeTree`: `ID *string`, `Content *string`, `Children []*NodeLeaf`.
- `NodeLeaf`: `ID *string`, `Content *string`, `Children []*NodeInfoBase`.
- `NodeInfoBase`: `ID *string`, `Content *string`.
- `DistinctOrNetServiceInfo`: `ID *string`, `Name *string`.

`DescribeNodeGroups` is a synchronous API with no pagination parameters (no `Limit`/`Offset`). The response is a complete tree, so no pagination loop is required.

## Goals / Non-Goals

**Goals:**
- Provide a read-only Terraform data source that exposes all filter inputs of `DescribeNodeGroups` and flattens the complete response (`NodeList` tree, `DistrictList`, `NetServiceList`) into Terraform state.
- Follow the established CAT data source conventions (`data_source_tc_cat_node.go`) for imports, helper usage, retry wrapping, and `result_output_file` handling.
- Follow the established list-datasource conventions (`data_source_tc_igtm_instance_list.go`) for the retry-wrapped service call and `NonRetryableError` on a nil response inside the retry block (per project rule: a data source Read MUST NOT clear the id on a transient empty API response).
- Provide gomonkey-mocked unit tests (no acceptance suite), per the project rule that new data sources use mock-based unit tests for business logic.
- Remain fully backward compatible (purely additive).

**Non-Goals:**
- No CRUD resource is created (this is a data source only).
- No pagination handling (the API does not paginate).
- No modification of the existing `tencentcloud_cat_node` data source.
- No changes to the `vendor/` SDK (already vendored via go.mod).

## Decisions

### Decision 1: No pagination loop in the service layer
`DescribeNodeGroups` has no `Limit`/`Offset` parameters and returns the complete node-group tree in a single call. We therefore do NOT implement an internal pagination loop (unlike `DescribeIgtmInstanceListByFilter`). The service method makes a single SDK call inside one `resource.Retry(tccommon.ReadRetryTimeout, ...)` block.

**Alternative considered:** Add a defensive pagination loop for future-proofing. **Rejected** because there are no pagination parameters on the request/response and inventing them would cause SDK errors.

### Decision 2: Nil-response handling returns `NonRetryableError` inside the retry block
Per the project rule for `RESOURCE_KIND_DATASOURCE`, the Read retry block MUST check whether the API returned an empty/nil response and return `resource.NonRetryableError` (NOT clear `d.SetId("")`). This prevents a transient API hiccup from wiping the local state id. The check covers `result == nil`, `result.Response == nil`, and `result.Response.NodeList == nil && result.Response.DistrictList == nil && result.Response.NetServiceList == nil`. An empty-but-non-nil response (all three lists empty) is treated as a valid result and does NOT error — it simply yields empty computed lists.

### Decision 3: Schema flattening matches the JsonPath mapping exactly
The output schema follows the documented JsonPath → SchemaName mapping:
- `node_list` (`TypeList`) — top-level `NodeTree` items, each with `id`, `content`, and nested `children`.
  - `children` (`TypeList`) — `NodeLeaf` items, each with `id`, `content`, and nested `children`.
    - inner `children` (`TypeList`) — `NodeInfoBase` items, each with `id`, `content` (no further nesting).
- `district_list` (`TypeList`) — `DistinctOrNetServiceInfo` items, each with `id`, `name`.
- `net_service_list` (`TypeList`) — `DistinctOrNetServiceInfo` items, each with `id`, `name`.

The same base name `children` is reused at two nesting levels (per the mapping), and `id`/`content` are reused at each level. Each `id` field is a `TypeString` (the SDK types are `*string`).

### Decision 4: `node_type` is a `TypeList` of `TypeInt`
The cloud API `NodeType` is `[]*int64` (a slice), while every other input filter is a scalar pointer. To faithfully represent the API, `node_type` is declared as `schema.TypeList` with `Elem: &schema.Schema{Type: schema.TypeInt}`. All other inputs are scalar `TypeInt`/`TypeString` matching their SDK pointer types. When building the request, `node_type` is converted from `[]interface{}` (Terraform) to `[]*int64` (SDK) via `helper.IntPtr` per element.

### Decision 5: Data source id uses `helper.DataResourceIdsHash`
Following the existing `tencentcloud_cat_node` data source convention, the id is computed as `helper.DataResourceIdsHash(ids)` where `ids` are the top-level `NodeTree.ID` values. This gives a deterministic, content-derived id (preferred over `helper.BuildToken()`'s random token for CAT data sources, and consistent with the sibling `cat_node` data source).

### Decision 6: `result_output_file` writes the flattened lists
Following `tencentcloud_cat_node`, `result_output_file` (if set) writes the flattened `nodeListTmp` (the node list maps) to the file. This keeps the CAT data sources internally consistent.

### Decision 7: `node_type` request field — single-value convenience
Because `DescribeNodeGroups.NodeType` is a list but the most common usage is a single value, we accept a Terraform list and convert each element to `*int64`. Users may pass one or more node types; the API takes the union.

## Risks / Trade-offs

- **[Risk] Nested schema depth (three levels of `children`) is verbose.** → **Mitigation:** The nesting is fixed by the cloud API response shape (`NodeTree` → `NodeLeaf` → `NodeInfoBase`); we mirror it exactly so flattening is mechanical and each level has only `id`/`content` (+ optional `children`). No deeper nesting exists in the SDK.
- **[Risk] Empty response could be confused with API error.** → **Mitigation:** We distinguish nil/empty: a structurally nil response (`result == nil`, `Response == nil`, or all three lists nil) raises `NonRetryableError` so the outer retry keeps trying; a non-nil response with empty lists is a valid empty result.
- **[Risk] `node_type` being a list diverges from the other scalar filters and may surprise users.** → **Mitigation:** This mirrors the cloud API exactly (`NodeType []*int64`). Documented in the schema `Description` and the example `.md`.
- **[Risk] Forgetting to update `provider.go` doc-comment breaks `make doc`.** → **Mitigation:** tasks.md explicitly lists updating both the `DataSourcesMap` entry and the provider.go doc-comment index.
- **[Trade-off] No deterministic ordering guarantee from the API.** The `node_list`/`district_list`/`net_service_list` order follows whatever `DescribeNodeGroups` returns; we do not re-sort.
