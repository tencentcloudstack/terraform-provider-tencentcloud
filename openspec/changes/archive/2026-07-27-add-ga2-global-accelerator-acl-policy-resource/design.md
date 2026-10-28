## Context

The Tencent Cloud GA2 (Global Accelerator V2) product line is a multi-tier object model:

```
GlobalAccelerator (parent instance)
├── Listener
│   └── EndpointGroup / ForwardingPolicy / ForwardingRule
├── AccelerateArea
└── AclPolicy (访问控制策略)   ← this change
    └── AclRule (out of scope here)
```

The provider already ships Terraform resources for `GlobalAccelerator`, `Listener`, `EndpointGroup`, `ForwardingPolicy`, `ForwardingRule` and `AccelerateArea`. The **AclPolicy** layer — which controls whether source traffic to an accelerator instance is accepted or dropped by default, and can be toggled OPEN/CLOSE — is currently console/API-only. This change adds `tencentcloud_ga2_global_accelerator_acl_policy` so the security posture of a GA2 instance is fully declarative.

The vendored SDK at `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ga2/v20250115/` already exposes all four required APIs plus the generic task oracle:

- `CreateGlobalAcceleratorAclPolicyWithContext` → request `{ GlobalAcceleratorId, DefaultAction }`, response `{ TaskId, GlobalAcceleratorAclPolicyId }` (asynchronous).
- `DescribeGlobalAcceleratorAclPoliciesWithContext` → request `{ GlobalAcceleratorId, Offset (*uint64), Limit (*string, max "200") }`, response `{ GlobalAcceleratorAclPolicySet []*GlobalAcceleratorAclPolicies, TotalCount }` (synchronous). `GlobalAcceleratorAclPolicies` = `{ GlobalAcceleratorAclPolicyId, DefaultAction, Status }`.
- `ModifyGlobalAcceleratorAclPolicyWithContext` → request `{ GlobalAcceleratorId, GlobalAcceleratorAclPolicyId, Status }`, response `{ TaskId }` (asynchronous). Only `Status` (OPEN/CLOSE) is mutable.
- `DeleteGlobalAcceleratorAclPolicyWithContext` → request `{ GlobalAcceleratorId, GlobalAcceleratorAclPolicyId }`, response `{ TaskId }` (asynchronous).
- `DescribeTaskResultWithContext` → returns `{ Status }`, used as the polling oracle.

The provider already has a connectivity binding (`UseGa2V20250115Client`), a service helper file (`service_tencentcloud_ga2.go`) containing `WaitForGa2TaskFinish(ctx, taskId, timeout)` (polls `DescribeTaskResult` until `Status == "SUCCESS"`), and the composite-ID convention (`tccommon.FILED_SP == "#"`) used by sibling resources such as `tencentcloud_ga2_forwarding_policy` and `tencentcloud_ga2_endpoint_group`.

## Goals / Non-Goals

**Goals:**
- Provide full lifecycle management of one GA2 ACL policy (create / read / update / delete / import) through Terraform.
- Schema fields exactly mirror the cloud APIs (no field renaming, no synthetic flags).
- All async writes (`Create` / `Modify` / `Delete`) must wait for `Status == SUCCESS` on the returned `TaskId` before returning to Terraform, so dependent resources and `terraform plan` reflect reality.
- Code style matches `tencentcloud_ga2_forwarding_policy` / `tencentcloud_ga2_endpoint_group` and the `tencentcloud_igtm_strategy` reference (single-file resource layout, retry on every SDK call, defensive nil checks on response payloads, composite-ID split helpers).
- Import works via the composite ID `<global_accelerator_id>#<global_accelerator_acl_policy_id>`.

**Non-Goals:**
- Managing `AclRule` (the individual allow/drop entries under a policy). That is a separate child resource (`CreateGlobalAcceleratorAclRule` / `DescribeGlobalAcceleratorAclRules` / `ModifyGlobalAcceleratorAclRule` / `DeleteGlobalAcceleratorAclRule` exist in the SDK) and is out of scope for this change.
- A datasource (`data_source_tc_ga2_global_accelerator_acl_policy` / `_acl_policies`); this change is resource-only.
- Batch creation of multiple policies per accelerator in a single resource block (one resource == one policy).
- Modifying `DefaultAction` after creation — the cloud API does not support it, so it is ForceNew.

## Decisions

### D1. Composite resource ID = `<global_accelerator_id>#<global_accelerator_acl_policy_id>`
Why: every Read/Update/Delete SDK call requires both `GlobalAcceleratorId` and `GlobalAcceleratorAclPolicyId`. `DescribeGlobalAcceleratorAclPolicies` has no per-policy filter slot — only `GlobalAcceleratorId` (+ paging) — so to re-read a single policy after creation we must retain both IDs in state. The sibling `tencentcloud_ga2_forwarding_policy` and `tencentcloud_ga2_endpoint_group` resources already use `tccommon.FILED_SP` ("#") composite IDs, so this is the established idiom.
Alternative considered: store only `GlobalAcceleratorAclPolicyId` and require the user to always pass `global_accelerator_id` separately. Rejected — `global_accelerator_id` is ForceNew and known at create time; embedding it in the ID avoids re-deriving it and keeps Read self-sufficient (critical for `terraform import` and refresh after external changes).
Alternative considered: a different separator. Rejected — `tccommon.FILED_SP` is the provider-wide constant for exactly this purpose.

### D2. Reuse `WaitForGa2TaskFinish` as-is
Why: It already accepts a caller-supplied `timeout time.Duration`, polls `DescribeTaskResult`, and treats `SUCCESS` as terminal. No ACL-policy-specific behavior is needed. This matches the decision recorded for the sibling `tencentcloud_ga2_global_accelerator` change.
Alternative: a separate `WaitForGa2AclPolicyTask` wrapper. Rejected — duplicate code with no extra value.

### D3. Add `DescribeGa2GlobalAcceleratorAclPolicyById` to the existing `Ga2Service`
Pattern matches `DescribeGa2GlobalAcceleratorById` / `DescribeGa2ForwardingPolicyById`:
- Build the request *outside* the for-loop (only `Offset`/`Limit` mutate per page).
- `Limit` set to `"200"` (the SDK-documented maximum for `DescribeGlobalAcceleratorAclPoliciesRequest.Limit`, which is a `*string`), passed as a literal — no new package-level constant, per prior feedback to avoid constant proliferation.
- `Offset` is `*uint64`, starting at `0`.
- Page through `GlobalAcceleratorAclPolicySet`; strict-equal on `*item.GlobalAcceleratorAclPolicyId == policyId` before returning, since the API has no per-policy filter.
- Wrap each SDK page in `resource.Retry(tccommon.ReadRetryTimeout, ...)`; return `resource.NonRetryableError` on nil `Response`.
- Return `(nil, nil)` when not found; the resource layer treats this as "deleted out of band" and calls `d.SetId("")` (after logging the id per the project rule).
Alternative: a single-shot non-paginated call relying on default paging. Rejected — an accelerator could legitimately hold more than the default page size of policies over time, and we follow the same idiom as the sibling helpers for symmetry and correctness.

### D4. `default_action` is ForceNew; only `status` is mutable
Why: `CreateGlobalAcceleratorAclPolicyRequest` accepts `DefaultAction` but `ModifyGlobalAcceleratorAclPolicyRequest` does **not** — it only accepts `Status` (OPEN/CLOSE). Modeling `default_action` as updatable would silently mislead users into thinking a `terraform apply` changed it. ForceNew makes the limitation honest: changing `default_action` destroys and recreates the policy.
`status` is Optional+Computed: it is returned by `DescribeGlobalAcceleratorAclPolicies` as `Status` and is the single field `ModifyGlobalAcceleratorAclPolicy` can change, so it is updatable in place.
Alternative considered: make `status` Computed-only and expose a separate `modify_status` operation. Rejected — Terraform idioms favor a directly-settable Optional+Computed attribute over imperative operations.

### D5. Async retry topology
Every SDK call is wrapped in `resource.Retry(timeoutScope, func() *resource.RetryError { ... })`:
- Read paths: `tccommon.ReadRetryTimeout`.
- Write paths (Create / Modify / Delete): `tccommon.WriteRetryTimeout`.
- Async polling (after the write SDK call succeeds): `d.Timeout(schema.TimeoutCreate|Update|Delete)` passed into `WaitForGa2TaskFinish`.
- Resource-level `Timeouts` block defaults to **5 minutes** for Create/Update/Delete (matches sibling `ga2` resources).
- `d.SetId(...)` and any state mutation happen **outside** and **after** the retry block + task-wait, never inside the retry closure (per project rule).
- Nil-response defenses (`result == nil || result.Response == nil || critical field == nil`) return `resource.NonRetryableError` with a descriptive message rather than dereferencing a nil pointer.
Why: the provider's two-tier retry model — SDK-call-level retries handle transient TencentCloudSDK retryable errors; task-level polling handles eventual consistency. This is the same shape used by every sibling `ga2` resource.

### D6. Schema parity with the four cloud APIs
Mapping (all SDK input fields appear, plus computed outputs; no extras):

| Schema field | Type | SDK source | Mutability |
|---|---|---|---|
| `global_accelerator_id` | `TypeString`, Required, ForceNew | `GlobalAcceleratorId` (Create/Describe/Modify/Delete) | ForceNew |
| `default_action` | `TypeString`, Required, ForceNew | `DefaultAction` (Create only) | ForceNew |
| `status` | `TypeString`, Optional, Computed | `Status` (Modify input + Describe output) | updatable (OPEN/CLOSE) |
| `global_accelerator_acl_policy_id` | `TypeString`, Computed | `GlobalAcceleratorAclPolicyId` (Create output + Describe output) | computed |
| `task_id` | `TypeString`, Computed | `TaskId` (Create/Modify/Delete output) | computed |

Alternative considered: omit `task_id` from schema. Rejected — sibling resources surface async task IDs for operator observability, and the user-supplied mapping explicitly lists `task_id` as a Terraform parameter path.

### D7. Import via composite ID
The resource uses `schema.ImportStatePassthrough`. On import, `Read` splits `d.Id()` on `tccommon.FILED_SP` into `(gaId, policyId)`, queries `DescribeGa2GlobalAcceleratorAclPolicyById(ctx, gaId, policyId)`, and hydrates state. The resource markdown `Import` section documents the composite-ID syntax.
Alternative: a custom `StateContext` importer that parses two CLI args. Rejected — passthrough + a documented composite ID is the sibling-resource idiom and needs no extra code.

### D8. Read "deleted out of band" handling
When `DescribeGa2GlobalAcceleratorAclPolicyById` returns `(nil, nil)`, the resource SHALL **first** log `log.Printf("[CRUD] ga2_global_accelerator_acl_policy id=%s", d.Id())` to preserve the id in logs, **then** call `d.SetId("")` and return nil. This follows the project rule that the id must not be lost before the log line is emitted, so operators can locate which invocation cleared state.
If the resource is `IsNewResource()` (Create-then-immediate-Read), return an explicit error instead of silently clearing.

### D9. File layout (single-file resource)
Per the established `ga2` convention, the entire resource lives in **one** file:
- `resource_tc_ga2_global_accelerator_acl_policy.go`: schema → Create/Read/Update/Delete → composite-ID split helper.
The service-layer helper `DescribeGa2GlobalAcceleratorAclPolicyById` lives in the existing `service_tencentcloud_ga2.go`.
No `_extension.go` file is generated.

### D10. Test strategy — gomonkey mocks, not Terraform acceptance suite
Per the project rule for **newly-added** resources, the test file uses `gomonkey` to mock the cloud-API client methods and exercises only the resource business logic (schema wiring, composite-ID split, nil defenses, retry/task-wait ordering). Tests run via `go test -gcflags=all=-l ./tencentcloud/services/ga2/...`. No `TF_ACC` acceptance suite is authored for this new resource.

## Risks / Trade-offs

- **[Risk]** `DescribeGlobalAcceleratorAclPolicies` has no per-policy filter; we paginate and match client-side. If an accelerator accumulates many policies, Read does multiple round-trips. → **Mitigation**: page size is the documented maximum (`Limit="200"`); the common case (≤ a handful of policies per accelerator) is a single page. No correctness risk because the strict-equal match guarantees we never hydrate the wrong policy.
- **[Risk]** `default_action` is ForceNew, so changing it recreates the policy (and drops any child `AclRule`s until they too are recreated — out of scope here). → **Mitigation**: document the ForceNew behavior in the resource markdown; the cloud API genuinely cannot mutate `DefaultAction`, so ForceNew is the only honest option.
- **[Risk]** TencentCloud occasionally surfaces an async failure as a non-`SUCCESS` terminal status (e.g. `FAIL`) without rich error text. → **Mitigation**: `WaitForGa2TaskFinish` already returns `RetryableError(... current status: <S>)`; the `resource.Retry` budget eventually surfaces a clear timeout to the user. Treating known terminal-failure statuses as non-retryable is out of scope.
- **[Trade-off]** Surfacing `task_id` as a Computed field means it flips on every Update/Delete. → **Mitigation**: Computed fields are not diffed for user-visible drift, so this is benign and matches the sibling resources.
- **[Trade-off]** `status` is Optional+Computed. If a user omits it, Terraform treats it as unknown until the first Read populates it from the cloud (the API returns the actual status). This is the standard Optional+Computed pattern and requires no `CustomizeDiff`.

## Migration Plan

This is purely additive. No state migration required:
1. Land the new resource + service helper + provider registration.
2. After release, users opt in by adding `resource "tencentcloud_ga2_global_accelerator_acl_policy" "x" { ... }` to their config, referencing an existing `tencentcloud_ga2_global_accelerator` via `global_accelerator_id`.
3. Existing `ga2` resources continue to behave exactly as before — nothing changes for them.

Rollback: pure revert of the new files + the `provider.go` registration line; no state mutations to undo.

## Open Questions

- None requiring user input. The SDK exposes all needed APIs; the proposal-level decisions (composite ID, schema parity, async polling reuse, gomonkey tests) are fully determined by the existing sibling-resource conventions.
