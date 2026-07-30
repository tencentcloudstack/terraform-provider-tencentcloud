## Context

The Tencent Cloud GA2 product line is a multi-tier object model. With this change, the full hierarchy becomes Terraform-native:

```
GlobalAccelerator                          ← tencentcloud_ga2_global_accelerator (shipped)
└── Listener                               ← tencentcloud_ga2_listener (shipped)
    ├── EndpointGroup                      ← tencentcloud_ga2_endpoint_group (shipped)
    └── ForwardingPolicy (HTTP/HTTPS only) ← THIS CHANGE
        └── ForwardingRule                 ← tencentcloud_ga2_forwarding_rule (shipped)
```

The `ForwardingPolicy` models the L7 host-based routing domain for an HTTP/HTTPS listener. Each HTTP/HTTPS listener gets a default forwarding policy created automatically by the cloud, and additional policies can be created via `CreateForwardingPolicy`. Users reference the `ForwardingPolicyId` when creating `ForwardingRule` resources.

The vendored SDK at `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ga2/v20250115/` already exposes:
- `CreateForwardingPolicyWithContext` → returns `{ TaskId, ForwardingPolicyId }` (asynchronous)
- `DescribeForwardingPolicyWithContext` → paged list keyed by `(GlobalAcceleratorId, ListenerId)`, returning `[]*ForwardingPolicySet`
- `ModifyForwardingPolicyWithContext` → returns `{ TaskId }` (asynchronous)
- `DeleteForwardingPolicyWithContext` → returns `{ TaskId }` (asynchronous)
- `DescribeTaskResultWithContext` → returns `{ Status }`, used as the polling oracle

The existing `Ga2Service` already provides `WaitForGa2TaskFinish(ctx, taskId, timeout)`; we will reuse it verbatim.

## Goals / Non-Goals

**Goals:**
- Provide full lifecycle management of a GA2 ForwardingPolicy through Terraform.
- Schema fields exactly mirror `CreateForwardingPolicyRequest` (no field renaming, no synthetic flags).
- All async writes wait for `Status == SUCCESS` on the returned `TaskId` before returning to Terraform.
- Code style matches the previously-shipped GA2 resources (single-file resource layout, retry on every SDK call, defensive nil checks on response payloads).
- Support `terraform import` via the 3-segment composite ID.

**Non-Goals:**
- Adding a datasource for forwarding policies (resource-only here).
- Any field that exists only in the `ForwardingPolicySet` describe response but not in `CreateForwardingPolicyRequest` is exposed Computed-only, never as input (e.g., `default_host_flag`).

## Decisions

### D1. 3-segment composite resource ID
Why: `Modify` / `Delete` and the lookup helper all need the 3-tuple `(gaId, listenerId, policyId)`. Persisting only the policy ID would force a re-discovery on every apply, which is brittle.

Format: `<gaId>#<listenerId>#<policyId>` using `tccommon.FILED_SP` — same separator already used by existing GA2 resources.

Alternative considered: bare `ForwardingPolicyId`. Rejected — would force every CRUD function to re-fetch parent IDs from state, which is fragile across imports.

### D2. Reuse `WaitForGa2TaskFinish` as-is
Already accepts a caller-supplied `timeout time.Duration`, polls `DescribeTaskResult`, treats `SUCCESS` as terminal. No policy-specific behavior is needed.

### D3. Add `DescribeGa2ForwardingPolicyById` to the existing `Ga2Service`
Pattern matches the previously-added GA2 helpers (`DescribeGa2ForwardingRuleById`):
- Build the request **outside** the for-loop. Set `request.GlobalAcceleratorId`, `request.ListenerId` once.
- Page size = `100` (documented maximum), passed as a literal — no new package-level constant.
- Strict-equals on `*item.ForwardingPolicyId == policyId` (and double-check the parent IDs for paranoia).
- Returns `(nil, nil)` when not found; the resource layer treats this as "deleted out of band" and calls `d.SetId("")`.

Note: `DescribeForwardingPolicyRequest` lacks a `Filters` field; pagination is the only filter mechanism.

### D4. Async retry topology
Every SDK call wrapped in `resource.Retry(timeoutScope, func() *resource.RetryError { ... })`:
- Read paths: `tccommon.ReadRetryTimeout`.
- Write paths (Create / Modify / Delete): `tccommon.WriteRetryTimeout`.
- Async polling (after the write succeeds): `d.Timeout(schema.TimeoutCreate|Update|Delete)` passed into `WaitForGa2TaskFinish`.
- Resource-level `Timeouts` block defaults to **5 minutes** for Create/Update/Delete (matches the other GA2 resources).

### D5. Schema parity with `CreateForwardingPolicyRequest`
Mapping (every CreateForwardingPolicy input field appears, no extras, no renames):

| Schema field | Type | Required? | ForceNew? | Source SDK field |
|---|---|---|---|---|
| `global_accelerator_id` | `TypeString` | Required | **Yes** | `GlobalAcceleratorId` |
| `listener_id` | `TypeString` | Required | **Yes** | `ListenerId` |
| `host` | `TypeString` | Required | No | `Host` |

Computed-only fields (not in `CreateForwardingPolicyRequest`, surfaced from `ForwardingPolicySet`):
- `forwarding_policy_id` (string) — also stored as the 3rd segment of `d.Id()`.
- `default_host_flag` (bool) — whether this policy is the default host for the listener.

### D6. ForceNew choices justified by the API
- `global_accelerator_id`, `listener_id`: `ModifyForwardingPolicy` carries them but only as *identifiers* of which policy to modify; it cannot move a policy across these boundaries. ForceNew prevents users from accidentally requesting a relocation that the API silently rejects.
- `host` is the only field that can be updated via `ModifyForwardingPolicy`.

### D7. Update path semantics
Per `ModifyForwardingPolicyRequest`, only `Host` is an updatable body field. The Update function:
- Skips the SDK call entirely if `host` hasn't changed (the ID fields are ForceNew, so they cannot trigger Update directly).
- Always populates the mandatory ID fields (`GlobalAcceleratorId`, `ListenerId`, `ForwardingPolicyId`) on the request.
- Sets `request.Host` from `d.Get("host")`.
- Awaits the returned `TaskId` via `WaitForGa2TaskFinish(ctx, taskId, d.Timeout(schema.TimeoutUpdate))`.

### D8. File layout
Single file: `resource_tc_ga2_forwarding_policy.go` — schema + Create/Read/Update/Delete + ID parser + build/flatten helpers in that order. Service-level helper lives in the existing `service_tencentcloud_ga2.go`. Matches the pattern used by other GA2 resources.

### D9. `make doc` flow + `provider.md` registration
Per the established workflow, the resource markdown lives at `tencentcloud/services/ga2/resource_tc_ga2_forwarding_policy.md`. The website file at `website/docs/r/ga2_forwarding_policy.html.markdown` is **never** hand-edited; it is regenerated by `make doc`. For `make doc` to discover the new resource, we must also append `tencentcloud_ga2_forwarding_policy` to the `Global Accelerator(GA2)` Resources section in `tencentcloud/provider.md`.

## Risks / Trade-offs

- **[Risk]** `DescribeForwardingPolicy` is keyed by `(gaId, listenerId)` and lacks a per-policy filter. → **Mitigation**: paginate with `Limit=100` and strict-equal on `ForwardingPolicyId` in the helper. For listeners with very large policy populations the page count grows linearly; this is acceptable because the API's documented maximum is the same `100` for all callers.
- **[Risk]** `UnsupportedOperation.DefaultForwardingPolicyOperate` error from Modify/Delete on the default policy. → **Mitigation**: documented in the resource markdown; this is a cloud-side constraint, not a Terraform concern. Users who need to modify/delete the default policy must handle it at the cloud level.
- **[Trade-off]** `global_accelerator_id` and `listener_id` are ForceNew. This means a policy cannot be moved between listeners. → This matches the API's hard constraint and is explicitly documented.

## Migration Plan

This is purely additive. No state migration required:
1. Land the new resource + service helper + provider registration.
2. After release, users opt in by adding `resource "tencentcloud_ga2_forwarding_policy" "x" { ... }` to their config.

Rollback: pure revert of the new files + the `provider.go` and `provider.md` registration lines; no state mutations to undo.

## Open Questions

- None blocking. All required SDK surface area is vendored; the spec-level decisions (3-segment composite ID, ForceNew set, async polling reuse) are determinate.