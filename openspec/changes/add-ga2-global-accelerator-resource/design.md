## Context

The Tencent Cloud GA2 (Global Accelerator V2) product line is a multi-tier object model:

```
GlobalAccelerator (parent instance, holds tags + cross-border config)
└── Listener
    └── EndpointGroup
        └── EndpointConfigurations
```

This provider already ships `tencentcloud_ga2_endpoint_group` (committed in the prior change). The parent `GlobalAccelerator` instance, however, can only be created from the console — Terraform users today must hand-create it and pass the resulting ID into `tencentcloud_ga2_endpoint_group.global_accelerator_id`. This change closes that gap.

The vendored SDK at `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ga2/v20250115/` already exposes:
- `CreateGlobalAcceleratorWithContext` → returns `{ TaskId, GlobalAcceleratorId }` (asynchronous)
- `DescribeGlobalAcceleratorsWithContext` → paged list with `Filters` (synchronous)
- `ModifyGlobalAcceleratorWithContext` → returns `{ TaskId }` (asynchronous)
- `DeleteGlobalAcceleratorWithContext` → returns `{ TaskId }` (asynchronous)
- `DescribeTaskResultWithContext` → returns `{ Status }`, used as the polling oracle

The provider already has a connectivity binding (`UseGa2V20250115Client`) and a service helper file (`service_tencentcloud_ga2.go`) containing `WaitForGa2TaskFinish(ctx, taskId, timeout)` which we will reuse verbatim.

## Goals / Non-Goals

**Goals:**
- Provide full lifecycle management of a GA2 Global Accelerator instance (create / read / update / delete / import) through Terraform.
- Schema fields exactly mirror `CreateGlobalAcceleratorRequest` (no field renaming, no synthetic flags), per the user's explicit rule.
- All async writes must wait for `Status == SUCCESS` on the returned `TaskId` before returning to Terraform, so subsequent dependent resources (`tencentcloud_ga2_endpoint_group`, listeners) can immediately reference the newly created accelerator.
- Code style matches `tencentcloud_igtm_monitor` and the recently added `tencentcloud_ga2_endpoint_group` (single-file resource layout, retry on every SDK call, defensive nil checks on response payloads).
- Filename conventions for the markdown doc and `_test.go` mirror `resource_tc_config_compliance_pack.md` / `_test.go`.

**Non-Goals:**
- Managing `Listener`, `EndpointConfigurations`, or `Forwarding Rules` (out of scope; `tencentcloud_ga2_endpoint_group` already exists; listeners will be a future change).
- Bandwidth-package or accelerator-area binding APIs (separate change if needed).
- Datasource implementation (`data_source_tc_ga2_global_accelerator(s)`); this change is resource-only.
- Custom IP / accelerator IP allocation API surface — not part of `CreateGlobalAccelerator`.

## Decisions

### D1. Single resource ID = `GlobalAcceleratorId`
Why: `CreateGlobalAcceleratorResponse` returns `GlobalAcceleratorId` as the unique identifier; `DescribeGlobalAccelerators` / `Modify*` / `Delete*` all key off the same ID. Using a composite ID would be ceremony without value.
Alternative considered: composite `<region>#<gaId>`. Rejected — region is already on the provider config; multi-region state isolation is not needed and would break import UX.

### D2. Reuse `WaitForGa2TaskFinish` as-is
Why: It already accepts a caller-supplied `timeout time.Duration`, polls `DescribeTaskResult`, and treats `SUCCESS` as terminal. No GA-specific behavior is needed.
Alternative: a separate `WaitForGa2GlobalAcceleratorTask` wrapper. Rejected — duplicate code with no extra value.

### D3. Add `DescribeGa2GlobalAcceleratorById` to the existing `Ga2Service`
Pattern matches `DescribeGa2EndpointGroupById`:
- Build the request *outside* the for-loop (only `Offset`/`Limit` mutate per page).
- `Filters` set to `[{Name: "global-accelerator-id", Values: [gaId]}]`.
- Page size = `100` (the SDK-documented maximum), passed as a literal — no new package-level constant, per the user's previous feedback to avoid constant proliferation.
- Final defensive equality check on `*item.GlobalAcceleratorId == gaId`, since filter semantics are advisory server-side.
- Return `(nil, nil)` when not found; the resource layer treats this as "deleted out of band" and calls `d.SetId("")`.
Alternative: a single-shot non-paginated call. Rejected — even if the filter is exact, the API still returns a `Set`; we follow the same idiom as the existing endpoint-group helper for symmetry.

### D4. Tag handling: split between Create and Update
- **Create**: forward `Tags` directly via `CreateGlobalAcceleratorRequest.Tags`. The SDK accepts `[]*Tag`, so we marshal `d.Get("tags").(map[string]interface{})` into `[]*ga2v20250115.Tag` inline. This is a single round-trip and satisfies the user's note about avoiding the deprecated post-create `ModifyTags` pattern.
- **Update**: `ModifyGlobalAcceleratorRequest` does **not** accept tags. Tag drift is reconciled with `svctag.NewTagService(...).ModifyTags(ctx, "qcs::ga2:<region>:uin/:globalAccelerator/<gaId>", replaceTags, deleteTags)` — the standard provider tag pipeline.
- **Read**: hydrate `tags` from `GlobalAcceleratorSet.TagSet`.
Why: avoids an unnecessary `ModifyTags` round-trip on Create while still supporting tag updates without forcing recreation.

### D5. Async retry topology
Every SDK call is wrapped in `resource.Retry(timeoutScope, func() *resource.RetryError { ... })`:
- Read paths: `tccommon.ReadRetryTimeout`.
- Write paths (Create / Modify / Delete): `tccommon.WriteRetryTimeout`.
- Async polling (after the write succeeds): `d.Timeout(schema.TimeoutCreate|Update|Delete)` passed into `WaitForGa2TaskFinish`.
- Resource-level `Timeouts` block defaults to **5 minutes** for Create/Update/Delete (matches `tencentcloud_ga2_endpoint_group`).
Why: the provider's two-tier retry model — SDK-call-level retries handle transient TencentCloudSDK retryable errors; task-level polling handles eventual consistency. This is the same shape used by `tencentcloud_ga2_endpoint_group`, so operators get a uniform mental model.

### D6. Schema parity with `CreateGlobalAcceleratorRequest`
Mapping (all SDK fields appear, no extras):

| Schema field | Type | Source | Notes |
|---|---|---|---|
| `name` | `TypeString`, Optional+Computed | `Name` | ≤60 bytes (validated server-side) |
| `instance_charge_type` | `TypeString`, Optional+Computed, ForceNew | `InstanceChargeType` | `POSTPAID` only today; `Modify` cannot change it → ForceNew |
| `description` | `TypeString`, Optional+Computed | `Description` | ≤100 bytes |
| `cross_border_type` | `TypeString`, Optional+Computed | `CrossBorderType` | `HighQuality` / `Unicom` |
| `cross_border_promise_flag` | `TypeBool`, Optional+Computed | `CrossBorderPromiseFlag` | required when cross-border is used |
| `tags` | `TypeMap`, Optional, Elem=String | `Tags []*Tag` | Create→inline, Update→tagService |
| `state` | `TypeString`, Computed | `State` | from `GlobalAcceleratorSet` |
| `status` | `TypeString`, Computed | `Status` | from `GlobalAcceleratorSet` (distinct from `State`) |
| `cname` | `TypeString`, Computed | `Cname` | |
| `ddos_id` | `TypeString`, Computed | `DdosId` | |
| `create_time` | `TypeString`, Computed | `CreateTime` | |
| `listener_counts` | `TypeInt`, Computed | `ListenerCounts` | |
| `accelerator_area_counts` | `TypeInt`, Computed | `AcceleratorAreaCounts` | |
| `global_accelerator_id` | `TypeString`, Computed | response field | also stored as `d.Id()` |

Alternative considered: collapse `state` and `status` since they sound similar. Rejected — the SDK exposes both as distinct fields; faithfully surfacing both lets users observe whatever semantic difference Tencent Cloud documents (we don't second-guess the API).

### D7. ForceNew choices
- `instance_charge_type` → ForceNew (the modify API has no slot for it).
- `name`, `description`, `cross_border_type`, `cross_border_promise_flag`, `tags` → updatable in place.
Rationale: minimize unnecessary destroys.

### D8. File layout (single-file resource)
Per user's prior strict feedback on the `endpoint_group` change ("不要拆分 _crud.go / _helpers.go"), the entire resource lives in **one** file:
- `resource_tc_ga2_global_accelerator.go`: schema + Create/Read/Update/Delete + `buildXxxRequest` helpers + `flatten` helpers, in that order.
The service-layer helper `DescribeGa2GlobalAcceleratorById` lives in the existing `service_tencentcloud_ga2.go`.

## Risks / Trade-offs

- **[Risk]** TencentCloud occasionally surfaces an async failure as a non-`SUCCESS` terminal status (e.g. `FAIL`) without rich error text. → **Mitigation**: `WaitForGa2TaskFinish` already returns `RetryableError(... current status: <S>)`; the `resource.Retry` budget eventually surfaces a clear timeout to the user. Future enhancement (out of scope here): treat known terminal-failure statuses as non-retryable.
- **[Risk]** `cross_border_promise_flag` is required by the API only when cross-border is used; modeling it as Optional+Computed means Terraform cannot validate the conditional pre-flight. → **Mitigation**: rely on server-side rejection; document the rule in the resource markdown. Adding a `CustomizeDiff` would couple us to undocumented business rules and was deemed over-engineering during the endpoint-group change.
- **[Risk]** Tag updates and accelerator-instance updates are now two round-trips (Modify call + tagService call). → **Mitigation**: this is the standard provider idiom (~50 other resources use it). The two-step is wrapped in the same Update function; on tag-only change we skip the Modify call by guarding on `d.HasChangeExcept("tags")`.
- **[Trade-off]** `instance_charge_type = ForceNew`. Today the API only accepts `POSTPAID`, so this is a no-op constraint; if PREPAID later becomes supported and is mutable, the field can be relaxed in a future minor change. Documented in the spec.

## Migration Plan

This is purely additive. No state migration required:
1. Land the new resource + service helper + provider registration on the `feat/ga2_nr` branch.
2. After release, users opt in by adding `resource "tencentcloud_ga2_global_accelerator" "x" { ... }` to their config.
3. Existing `tencentcloud_ga2_endpoint_group` resources continue to reference `global_accelerator_id` exactly as before — nothing changes for them.

Rollback: pure revert of the new files + the `provider.go` registration line; no state mutations to undo.

## Open Questions

- None requiring user input. The SDK exposes all needed APIs; the proposal-level decisions (single ID, schema parity, async polling reuse) are fully determined.
