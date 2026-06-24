## Context

The Tencent Cloud GA2 product line is a multi-tier object model:

```
GlobalAccelerator (parent instance, holds tags + cross-border config)
├── Listener                 ← already shipped (tencentcloud_ga2_listener)
│   └── EndpointGroup        ← already shipped (tencentcloud_ga2_endpoint_group)
│       └── ForwardingRule   ← already shipped (tencentcloud_ga2_forwarding_rule)
└── AccelerateArea           ← THIS CHANGE (per-region acceleration entry)
```

An `AccelerateArea` is the acceleration-region entry attached to a global accelerator instance: it declares the access region (`AccelerateRegion`), ISP type, bandwidth, IP version and (optionally) the bound IPs. The provider already ships the listener/endpoint-group/forwarding-rule tiers; the acceleration-region tier must be added so users can construct the full chain in HCL with no console hand-off.

The vendored SDK at `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ga2/v20250115/` already exposes:
- `CreateAccelerateAreasWithContext` → request `{ GlobalAcceleratorId, AcceleratorAreas []*AcceleratorAreas }`, response `{ TaskId }` (**asynchronous; no AcceleratorAreaId returned**).
- `DescribeAccelerateAreasWithContext` → request `{ GlobalAcceleratorId, Offset, Limit }`, response `{ AccelerateAreaSet []*AcceleratorAreas, TotalCount }`.
- `ModifyAccelerateAreasWithContext` → request `{ GlobalAcceleratorId, AcceleratorAreas []*AcceleratorAreas }`, response `{ TaskId }` (asynchronous).
- `DeleteAccelerateAreasWithContext` → request `{ GlobalAcceleratorId, AcceleratorAreaIds []*string }`, response `{ TaskId }` (asynchronous).
- `DescribeTaskResultWithContext` → returns `{ Status }`, used as the polling oracle.

The shared `AcceleratorAreas` SDK struct carries: `AccelerateRegion`, `Bandwidth`, `IspType`, `IpVersion`, `AcceleratorAreaId`, `IpAddress []*string`, `IpAddressInfoSet []*IpAddressInfoSet`. The `IpAddressInfoSet` struct carries `{ IpAddress, IspType }`.

The existing `Ga2Service` already provides `WaitForGa2TaskFinish(ctx, taskId, timeout)`; we will reuse it verbatim.

## Goals / Non-Goals

**Goals:**
- Provide full lifecycle management of a single GA2 acceleration region through Terraform.
- Schema input fields exactly mirror the `CreateAccelerateAreas` payload (the `AcceleratorAreas` element flattened onto the resource), per the user's explicit rule — no field renaming, no synthetic flags.
- All async writes wait for `Status == SUCCESS` on the returned `TaskId` before returning to Terraform.
- Resolve the area ID after Create via `DescribeAccelerateAreas` keyed on `(GlobalAcceleratorId, AccelerateRegion)`, because `CreateAccelerateAreas` does not return `AcceleratorAreaId`.
- Code style matches `tencentcloud_igtm_monitor` and the previously-shipped GA2 resources (single-file resource layout, retry on every SDK call, defensive nil checks on response payloads).
- Filename conventions for the markdown doc and `_test.go` mirror `resource_tc_config_compliance_pack.md` / `_test.go`.

**Non-Goals:**
- A `tencentcloud_ga2_accelerate_areas` datasource — resource-only here.
- Managing more than one acceleration region per resource block (one resource == one `AcceleratorAreaId`). Batch creation via the API's `AcceleratorAreas` list is intentionally narrowed to a single element.
- Any field that exists only in the describe response but not in `CreateAccelerateAreas` is exposed Computed-only, never input.

## Decisions

### D1. Composite resource ID = `<GlobalAcceleratorId>#<AcceleratorAreaId>`
Why: `ModifyAccelerateAreas` and `DeleteAccelerateAreas` are keyed by `AcceleratorAreaId`, while `DescribeAccelerateAreas` is keyed by `GlobalAcceleratorId`. Persisting both in `d.Id()` avoids re-lookups and makes `terraform import` self-contained.
Implementation: use the project-standard separator `tccommon.FILED_SP` (already used by the other GA2 resources). Importer is `schema.ImportStatePassthrough`.
Alternative considered: ID = bare `AcceleratorAreaId` with `global_accelerator_id` sourced from schema. Rejected because passthrough import then cannot populate `global_accelerator_id`.

### D2. Post-create area-ID resolution by region (core business rule)
`CreateAccelerateAreasResponseParams` exposes only `TaskId` — never `AcceleratorAreaId`. The natural key the user controls is `(GlobalAcceleratorId, AccelerateRegion)`. Therefore Create:
1. Calls `CreateAccelerateAreas` with a single-element `AcceleratorAreas` list and captures `TaskId`.
2. Waits for the task via `WaitForGa2TaskFinish(ctx, taskId, d.Timeout(schema.TimeoutCreate))`.
3. Calls `DescribeGa2AccelerateAreaByRegion(ctx, gaId, region)` and reads back the server-generated `AcceleratorAreaId`.
4. `d.SetId(strings.Join([]string{gaId, areaId}, tccommon.FILED_SP))`.
Consequence: `accelerate_region` must be **ForceNew** (it is the natural key; changing it identifies a different area entirely).
Edge case: if the region lookup returns `(nil, nil)` after a successful task, Create returns an explicit error (the area was expected to exist).

### D3. Reuse `WaitForGa2TaskFinish` as-is
It already accepts a caller-supplied `timeout time.Duration`, polls `DescribeTaskResult`, and treats `SUCCESS` as terminal. No accelerate-area-specific behavior is needed.

### D4. Two service helpers on the existing `Ga2Service`
Both follow the established pagination/retry pattern from `DescribeGa2GlobalAcceleratorById`:
- `DescribeGa2AccelerateAreaById(ctx, gaId, areaId) (*ga2v20250115.AcceleratorAreas, error)` — used by Read/Update/Delete; strict-equals on `*item.AcceleratorAreaId == areaId`.
- `DescribeGa2AccelerateAreaByRegion(ctx, gaId, region) (*ga2v20250115.AcceleratorAreas, error)` — used by Create; strict-equals on `*item.AccelerateRegion == region`.

Both helpers:
- Build the `DescribeAccelerateAreasRequest` and set `request.GlobalAcceleratorId` **outside** the for-loop; only `Offset` / `Limit` mutate per page.
- Page with `Limit=100` (the documented maximum), passed as a literal — no new package-level constant.
- Wrap each SDK page in `resource.Retry(tccommon.ReadRetryTimeout, ...)`.
- Defend against nil `result` / `result.Response` and skip nil items / nil match-key pointers.
- Return `(nil, nil)` when not found; the resource Read treats this as "deleted out of band" and calls `d.SetId("")`.

Note: `DescribeAccelerateAreasRequest` has **no per-area Filter slot** (only `GlobalAcceleratorId` + `Offset`/`Limit`), so matching is performed strictly client-side, exactly like `DescribeGa2ForwardingRuleById`.

### D5. Schema parity with the `CreateAccelerateAreas` payload
The resource manages a single `AcceleratorAreas` element, flattened onto the resource. Mapping (every Create input field appears, no extras, no renames):

| Schema field | Type | Required? | ForceNew? | Source SDK field | Notes |
|---|---|---|---|---|---|
| `global_accelerator_id` | `TypeString` | Required | **Yes** | `GlobalAcceleratorId` | An area cannot be moved between accelerators |
| `accelerate_region` | `TypeString` | Required | **Yes** | `AcceleratorAreas[].AccelerateRegion` | Natural key used to resolve `AcceleratorAreaId` post-create |
| `bandwidth` | `TypeInt` | Optional+Computed | No | `AcceleratorAreas[].Bandwidth` | `uint64` in SDK; converted via `helper.IntUint64` |
| `isp_type` | `TypeString` | Optional+Computed | No | `AcceleratorAreas[].IspType` | `BGP` / `三网` / `精品`; default `BGP` |
| `ip_version` | `TypeString` | Optional+Computed | No | `AcceleratorAreas[].IpVersion` | Only `IPv4` supported; default `IPv4` |
| `ip_address` | `TypeSet`, Elem=string | Optional+Computed | No | `AcceleratorAreas[].IpAddress []*string` | TypeSet — order is not semantic |

Computed-only fields (surfaced from the `AcceleratorAreas` describe item, not part of the Create payload):
- `accelerator_area_id` (string) — also stored as the second segment of `d.Id()`.
- `ip_address_info_set` (list of nested blocks, computed) — each block: `ip_address` (string), `isp_type` (string). Sourced from `AcceleratorAreas[].IpAddressInfoSet`.

### D6. ForceNew choices justified by the API
- `global_accelerator_id`: an area belongs to exactly one accelerator; moving it requires recreate.
- `accelerate_region`: it is the natural key used to discover the server-generated `AcceleratorAreaId`; changing it points at a different area and must recreate.

All remaining input fields (`bandwidth`, `isp_type`, `ip_version`, `ip_address`) are updatable in place via `ModifyAccelerateAreas`, which accepts the same `AcceleratorAreas` struct.

### D7. Update sends a single fully-populated `AcceleratorAreas` element
`ModifyAccelerateAreas` takes `{ GlobalAcceleratorId, AcceleratorAreas []*AcceleratorAreas }`. The resource builds a single-element list where:
- `AcceleratorAreaId` is set from the parsed composite ID (identifies which area to modify).
- `AccelerateRegion` is set from state (ForceNew, unchanged) so the server can correlate the entry.
- `Bandwidth` / `IspType` / `IpVersion` / `IpAddress` carry the current schema values.

The Update function short-circuits (skips the Modify call) when none of the mutable fields (`bandwidth`, `isp_type`, `ip_version`, `ip_address`) changed.

### D8. `ip_address` uses `TypeSet`
`IpAddress []*string` is an unordered SDK string slice; `TypeSet` avoids spurious diffs on reordering, mirroring the cert-list decision in `tencentcloud_ga2_listener`. Conversion uses `(*schema.Set).List()` then `helper.InterfacesStringsPoint` (or an explicit `[]*string` build) when constructing requests.

### D9. File layout
Single file: `resource_tc_ga2_accelerate_area.go` (schema + Create/Read/Update/Delete + ID parser + build/flatten helpers). Service-level helpers live in the existing `service_tencentcloud_ga2.go`. This matches the user's strict feedback during prior GA2 changes to avoid `_crud.go` / `_helpers.go` splits.

### D10. `make doc` flow + `provider.md` registration
The resource markdown lives at `tencentcloud/services/ga2/resource_tc_ga2_accelerate_area.md`. The website file at `website/docs/r/ga2_accelerate_area.html.markdown` is **never** hand-edited; it is regenerated by `make doc`. For `make doc` to discover the new resource, we must also append `tencentcloud_ga2_accelerate_area` to the `Global Accelerator(GA2)` Resources section in `tencentcloud/provider.md`.

## Risks / Trade-offs

- **[Risk]** Race / eventual-consistency between task `SUCCESS` and the area appearing in `DescribeAccelerateAreas`. → **Mitigation**: the area lookup itself runs inside `resource.Retry(tccommon.ReadRetryTimeout, ...)` via the service helper; if the area is briefly absent, Create surfaces an explicit error rather than persisting a partial ID.
- **[Risk]** `DescribeAccelerateAreas` has no server-side area filter, so all regions under the accelerator are paged and matched client-side. → **Mitigation**: `Limit=100` (documented max) keeps round-trips low; strict-equality matching mirrors the existing `DescribeGa2ForwardingRuleById` approach.
- **[Risk]** Two areas with the same `AccelerateRegion` under one accelerator would make the region-based post-create resolution ambiguous. → **Mitigation**: the API treats `(GlobalAcceleratorId, AccelerateRegion)` as unique for an acceleration entry; the helper returns the first strict match and Create relies on that uniqueness, which the spec documents.
- **[Trade-off]** `accelerate_region` ForceNew means changing the region destroys and recreates the area. This matches the natural-key constraint and is explicitly documented.

## Migration Plan

This is purely additive. No state migration required:
1. Land the new resource + service helpers + provider registration on the `feat/ga2_nr` branch.
2. After release, users opt in by adding `resource "tencentcloud_ga2_accelerate_area" "x" { ... }`.
3. Existing GA2 resources are unaffected.

Rollback: pure revert of the new files + the `provider.go` and `provider.md` registration lines; no state mutations to undo.

## Open Questions

- None blocking. All required SDK surface area is vendored; the spec-level decisions (composite ID, post-create region resolution, ForceNew set, async polling reuse, TypeSet for `ip_address`) are determinate.
