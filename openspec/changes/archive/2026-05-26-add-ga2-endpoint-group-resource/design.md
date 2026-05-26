## Context

GA2 endpoint groups are a child object of a listener under a global accelerator. The SDK at `tencentcloud/ga2/v20250115` provides:

- `CreateEndpointGroup(GlobalAcceleratorId, ListenerId, EndpointGroupType, EndpointGroupConfiguration)` → `{TaskId, EndpointGroupId}`
- `DescribeEndpointGroups(GlobalAcceleratorId, Offset, Limit, Filters)` → `{EndpointGroupConfigurationSet[], TotalCount}`
- `ModifyEndpointGroup(GlobalAcceleratorId, ListenerId, EndpointGroupId, ...flat fields)` → `{TaskId}`
- `DeleteEndpointGroups(GlobalAcceleratorId, ListenerId, EndpointGroupIds)` → `{TaskId}`
- `DescribeTaskResult(TaskId)` → `{Status}` (terminal `SUCCESS`)

There is no existing `ga2` service package, no `UseGa2V20250115Client` connectivity entry, and no Terraform resource exposing any GA2 functionality. This change introduces all three.

The styling baseline is `tencentcloud_igtm_monitor` (see `tencentcloud/services/igtm/resource_tc_igtm_monitor.go`): `defer tccommon.LogElapsed`, `defer tccommon.InconsistentCheck`, `tccommon.NewResourceLifeCycleHandleFuncContext`, all SDK calls wrapped in `resource.Retry(tccommon.WriteRetryTimeout, …)`, nil-safe response unpacking via `if result == nil || result.Response == nil { resource.NonRetryableError(…) }`.

## Goals / Non-Goals

**Goals:**
- Resource `tencentcloud_ga2_endpoint_group` with CRUD lifecycle, importable.
- Schema mirrors `CreateEndpointGroup` request 1:1 (top-level + nested block).
- Composite ID `<global_accelerator_id>#<listener_id>#<endpoint_group_id>`, parsed in Read/Update/Delete and re-joined for `d.SetId`.
- Wrap every SDK call in `resource.Retry`. After Create/Update/Delete, additionally poll `DescribeTaskResult` until `Status == SUCCESS`.
- Service helper `DescribeGa2EndpointGroupById` performs paginated lookup using `endpoint-group-id` filter, with `Limit = 100` (API max).
- Defensive nil checks on every SDK return value before pointer dereference.

**Non-Goals:**
- Listing/data-source resources for GA2.
- Other GA2 entity types (accelerator, listener, etc.) — this PR only covers endpoint groups.
- Internal-only fields (`HealthCheckStatus` is read-only inside `EndpointConfigurations` and may be optionally exposed as Computed in nested elements; not modeled separately).

## Decisions

### D1: Composite ID format `<gaId>#<listenerId>#<endpointGroupId>`
`CreateEndpointGroup` requires all three IDs to look up an endpoint group, and the response only returns `EndpointGroupId`. We embed the three required IDs into the resource ID using `#` as separator (consistent with project-wide convention and `tccommon.FILED_SP`). Read/Update/Delete split on `#` and validate `len(parts) == 3`.

### D2: Schema 1:1 with `CreateEndpointGroupRequestParams`
Top-level fields:

| Schema field | SDK field | Type | Optional/Required |
|---|---|---|---|
| `global_accelerator_id` | `GlobalAcceleratorId` | TypeString | Required + ForceNew |
| `listener_id` | `ListenerId` | TypeString | Required + ForceNew |
| `endpoint_group_type` | `EndpointGroupType` | TypeString | Required + ForceNew (`VIRTUAL`/`DEFAULT`) |
| `endpoint_group_configuration` | `EndpointGroupConfiguration` | TypeList(MaxItems=1, Elem=Resource) | Required |

The nested block matches every field of `EndpointGroupConfiguration`, including the inner `endpoint_configurations` list (TypeList of Resource matching `EndpointConfigurations`) and `port_overrides` list (TypeList of Resource matching `PortOverride`).

Computed-only attribute: `endpoint_group_id` exposes the SDK-assigned ID for convenience.

### D3: Update path uses flat fields, not the nested block
`ModifyEndpointGroupRequestParams` is shaped differently from `CreateEndpointGroupRequestParams` — it expects most config attributes flat at top level (no `EndpointGroupConfiguration` wrapper). The Update implementation reads the nested `endpoint_group_configuration[0]` block and copies fields onto the flat `ModifyEndpointGroupRequest`. The user-facing schema remains a single nested block in both Create and Update for consistency. `endpoint_group_type` is `ForceNew` because Modify does not accept it.

### D4: Async polling via `DescribeTaskResult`
After every `Create`/`Modify`/`Delete` SDK call, we get `TaskId`. We then call:

```
resource.Retry(tccommon.WriteRetryTimeout*2, func() *resource.RetryError {
    result, err := client.DescribeTaskResultWithContext(ctx, &DescribeTaskResultRequest{TaskId: &taskId})
    if err != nil { return tccommon.RetryError(err) }
    if result == nil || result.Response == nil || result.Response.Status == nil {
        return resource.NonRetryableError(...)
    }
    if *result.Response.Status == "SUCCESS" { return nil }
    return resource.RetryableError(fmt.Errorf("task %s status=%s", taskId, *result.Response.Status))
})
```

Wrapped into `service.WaitForGa2TaskFinish(ctx, taskId)` for reuse. Terminal: `SUCCESS`. Anything else is a retryable poll until timeout.

### D5: Read uses pagination with API max `Limit=100`
`DescribeEndpointGroups` is a list API, not a get-by-id API. Read passes `GlobalAcceleratorId` plus a `Filters: [{Name: "endpoint-group-id", Values: [<id>]}]` filter; we paginate with `Offset/Limit=100` until either we find the matching record or `len(results) < limit`. Service helper returns `(*EndpointGroupConfigurationSet, error)` where a nil pointer means not found (Read sets `d.SetId("")`).

### D6: Defensive nil checks throughout
- All SDK responses are checked: `if result == nil || result.Response == nil { return resource.NonRetryableError(...) }` before any pointer dereference.
- All response field reads in Read are guarded with `if respData.Field != nil { _ = d.Set(...) }`.
- `EndpointGroupId`, `TaskId`, `Status` are explicitly checked for nil before deref.

### D7: Connectivity client lazy-init
`UseGa2V20250115Client` follows the exact pattern of `UseIgtmV20231024Client`: cached `*ga2v20250115.Client`, lazy-initialized with `me.NewClientProfile(300)` and `WithHttpTransport(&LogRoundTripper{})`. Imports added at the top of `connectivity/client.go`.

## Risks / Trade-offs

- [Risk] Async task may exceed `WriteRetryTimeout*2` for slow regions. Mitigation: choose `WriteRetryTimeout*2` (10 min) which is the project convention for similarly async resources; if it proves insufficient, extend to `*3` per regional metrics.
- [Risk] `DescribeEndpointGroups` filter `endpoint-group-id` is documented but if the API rejects unknown filter names, Read may scan all groups. Mitigation: helper falls back to filtering client-side on the response set if SDK returns an empty result with `TotalCount > 0`.
- [Risk] `EndpointConfigurations.HealthCheckStatus` is a server-side computed field; we mark it Computed in the nested schema so plans don't drift.
- [Trade-off] Schema duplicates flat health-check fields between `endpoint_group_configuration` block and would-be top-level, but we keep them only in the nested block for symmetry with Create. Update internally re-flattens.
