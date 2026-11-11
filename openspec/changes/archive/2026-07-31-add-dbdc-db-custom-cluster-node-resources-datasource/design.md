## Context

The `dbdc` service (TencentDB Dedicated Cluster, "DB Custom") already has
several Terraform data sources in
`tencentcloud/services/dbdc/data_source_tc_dbdc_db_custom_*.go`, each following
the same pattern: a schema with query arguments, a `Read` function wrapped in
`resource.Retry(tccommon.ReadRetryTimeout, ...)`, and a service-layer helper in
`service_tencentcloud_dbdc.go` that calls the SDK.

The cloud API `DescribeDBCustomClusterNodeResources` (vendored SDK
`dbdc/v20201029`) exposes per-node resource metrics:
`Capacity`, `Allocatable`, `Requests`, `Limits`, `Available` — each a
`MetaResource{Cpu, Memory, Pods}`. Unlike `DescribeDBCustomClusterNodes`, this
API has **no pagination** (`Offset`/`Limit`); it filters by `NodeIds` (max 50
per request). There is currently no Terraform data source for it.

## Goals / Non-Goals

**Goals:**
- Add `tencentcloud_dbdc_db_custom_cluster_node_resources` data source that
  queries node resource info by `cluster_id` (required) and optional `node_ids`.
- Reuse the established `dbdc` data-source code style (mirror
  `data_source_tc_dbdc_db_custom_cluster_nodes.go`).
- Correctly model the nested `MetaResource` blocks in the schema.
- Provide a gomonkey-mock unit test (no real cloud calls).

**Non-Goals:**
- No CRUD resource (this is read-only / RESOURCE_KIND_DATASOURCE).
- No pagination logic (the API does not support it; it filters via `NodeIds`).
- No website/docs generation here — that is handled by `make doc` in the
  finalize stage. Only the `.md` example doc under `tencentcloud/services/dbdc/`
  is authored here.

## Decisions

### Decision 1: Synchronous API, no pagination loop
`DescribeDBCustomClusterNodeResources` is synchronous (returns directly, no
`TaskId`) and has no `Offset`/`Limit` params. It is queried by `ClusterId`
(required) and an optional `NodeIds` list (max 50). Therefore the service helper
performs a single call inside `resource.Retry` rather than a pagination loop.

**Rationale**: The other `dbdc` helpers paginate because their APIs support
it; this one does not. Adding a fake loop would be misleading and could hide
truncation (the API caps `NodeIds` at 50, which is enforced server-side).

### Decision 2: Schema mirrors `DBCustomClusterNodeResource` + `MetaResource`
The output `node_set` is a `TypeList` of `TypeResource` where each element has:
- `node_id` (string, computed)
- Five nested `MetaResource` blocks, each a single-block `TypeList`
  (`MaxItems: 1`) with `cpu` (float), `memory` (float), `pods` (int):
  `capacity`, `allocatable`, `requests`, `limits`, `available`.

`MetaResource` fields use Go pointer types (`*float64`, `*uint64`); the Read
must nil-check each pointer before setting. `cpu`/`memory` are surfaced as
`TypeFloat` and `pods` as `TypeInt` to match the SDK semantics (Cpu/Memory are
cores/GiB, Pods is a count).

### Decision 3: Resource ID is a synthetic token
Per the reference data source `tencentcloud_dbdc_db_custom_cluster_nodes`, the
data source `d.SetId(helper.BuildToken())` — a non-persistent token — because a
query has no natural identity. This matches the existing pattern.

### Decision 4: Service helper signature
Add `DescribeDBCustomClusterNodeResourcesByFilter(ctx, param map[string]interface{})`
returning `([]*dbdcv20201029.DBCustomClusterNodeResource, error)`, mirroring the
return shape of sibling helpers but without `totalCount` (the API returns no
`TotalCount`). The helper:
- Builds `request.ClusterId` (`*string`) and `request.NodeIds` (`[]*string`)
  from the param map.
- Calls `me.client.UseDbdcV20201029Client().DescribeDBCustomClusterNodeResources(request)`
  inside `resource.Retry(tccommon.ReadRetryTimeout, ...)` with `ratelimit.Check`.
- Inside the retry block, returns `resource.NonRetryableError` if
  `result == nil || result.Response == nil || result.Response.NodeSet == nil`.
  On the failure path logs `log.Printf("[DATASOURCE] read empty, skip SetId")`.
- Returns `response.Response.NodeSet` on success.

**Rationale**: Keeps the helper shape consistent with
`DescribeDBCustomClusterNodesByFilter` while honestly reflecting that this API
has no count and no paging.

### Decision 5: Read function structure (per project rules #8, #14)
```
Read:
  defer LogElapsed / InconsistentCheck
  build paramMap: ClusterId (required), NodeIds (optional)
  var respData []*DBCustomClusterNodeResource
  resource.Retry(ReadRetryTimeout): respData, _ = service.DescribeDBCustomClusterNodeResourcesByFilter(...)
  if reqErr != nil: return reqErr
  // On empty result, the retry block already returned NonRetryableError;
  // do NOT d.SetId("") here.
  build nodeSetList with nil-guarded pointer access (Capacity/Allocatable/.../Available each nil-checked; each MetaResource field nil-checked)
  _ = d.Set("node_set", nodeSetList)
  d.SetId(helper.BuildToken())
  if result_output_file: WriteToFile
  return nil
```

### Decision 6: Unit tests via gomonkey (no TF acceptance suite)
Per project rule, new resources use gomonkey mock unit tests (not the TF
acceptance test suite). Tests:
- `TestDbdcDbCustomClusterNodeResourcesDS_ReadBasic` — mock the SDK
  `DescribeDBCustomClusterNodeResources` to return a 2-element `NodeSet` with
  populated `MetaResource` fields; assert schema mapping incl. nested blocks.
- `TestDbdcDbCustomClusterNodeResourcesDS_Schema` — assert schema fields/types.
- `TestDbdcDbCustomClusterNodeResourcesDS_ReadWithEmptyResponse` — mock returns
  empty `NodeSet`; assert `Read` errors (NonRetryableError from service layer).

## Risks / Trade-offs

- **[Risk] `NodeIds` max 50 enforced server-side** → The data source passes
  `node_ids` straight through; if a user supplies >50 the cloud API returns an
  `InvalidParameterValue` error. We document the limit in the schema
  description but do not client-side chunk (the API is the source of truth and
  the limit rarely hit for a single cluster). Trade-off: simpler implementation.
- **[Risk] Nested `MetaResource` could be nil** → All five blocks are optional
  in practice ("此字段可能返回 null"). Each is nil-checked before set; a nil
  block is simply omitted from the state map.
- **[Risk] `pods` is `*uint64` in SDK but schema uses `TypeInt`** → Terraform
  SDK `TypeInt` is used; the value is read via `*node.Capacity.Pods` (uint64 →
  int conversion). Nil-checked first. Acceptable since Pod counts are small.
