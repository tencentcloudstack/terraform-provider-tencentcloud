## Context

The Tencent Cloud GA2 product line is a multi-tier object model. With this change, the full hierarchy becomes Terraform-native:

```
GlobalAccelerator                          ← tencentcloud_ga2_global_accelerator (shipped)
└── Listener                               ← tencentcloud_ga2_listener (shipped)
    ├── EndpointGroup                      ← tencentcloud_ga2_endpoint_group (shipped)
    └── ForwardingPolicy (HTTP/HTTPS only)
        └── ForwardingRule                 ← THIS CHANGE
```

The `ForwardingPolicy` itself is provisioned automatically by the cloud when the L7 listener is created (it has no dedicated `Create*Policy` SDK call); users reference its `ForwardingPolicyId` directly when authoring rules. This change therefore models only the leaf `ForwardingRule` object.

The vendored SDK at `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ga2/v20250115/` already exposes:
- `CreateForwardingRuleWithContext` → returns `{ TaskId, ForwardingRuleId }` (asynchronous)
- `DescribeForwardingRuleWithContext` → paged list keyed by `(GlobalAcceleratorId, ListenerId, ForwardingPolicyId)`, returning `[]*ForwardingRuleSet`
- `ModifyForwardingRuleWithContext` → returns `{ TaskId }` (asynchronous)
- `DeleteForwardingRuleWithContext` → returns `{ TaskId }` (asynchronous)
- `DescribeTaskResultWithContext` → returns `{ Status }`, used as the polling oracle

The existing `Ga2Service` already provides `WaitForGa2TaskFinish(ctx, taskId, timeout)`; we will reuse it verbatim.

## Goals / Non-Goals

**Goals:**
- Provide full lifecycle management of a GA2 ForwardingRule through Terraform.
- Schema fields exactly mirror `CreateForwardingRuleRequest` (no field renaming, no synthetic flags), per the user's explicit rule.
- All async writes wait for `Status == SUCCESS` on the returned `TaskId` before returning to Terraform.
- Code style matches the previously-shipped GA2 resources (single-file resource layout, retry on every SDK call, defensive nil checks on response payloads).
- Filename conventions mirror `resource_tc_config_compliance_pack.md` / `_test.go`.

**Non-Goals:**
- A `tencentcloud_ga2_forwarding_policy` resource (the policy is auto-provisioned by the cloud; users obtain `ForwardingPolicyId` from the listener-level describe — modeling it as a separate resource would create an empty pass-through with no API surface).
- A `tencentcloud_ga2_forwarding_rules` datasource — resource-only here.
- Any field that exists only in the `ForwardingRuleSet` describe response but not in `CreateForwardingRule` is exposed Computed-only, never input.

## Decisions

### D1. 4-segment composite resource ID
Why: `Modify` / `Delete` and the lookup helper all need the 4-tuple `(gaId, listenerId, policyId, ruleId)`. Persisting only the rule ID would force a re-discovery on every apply, which is brittle (`DescribeForwardingRule` is keyed by the full `(gaId, listenerId, policyId)` triple, not by rule ID alone).
Format: `<gaId>#<listenerId>#<policyId>#<ruleId>` using `tccommon.FILED_SP` — same separator already used by the existing `tencentcloud_ga2_endpoint_group` (3-segment) and `tencentcloud_ga2_listener` (2-segment).
Alternative considered: bare `ForwardingRuleId`. Rejected — would force every CRUD function to first re-fetch parent IDs from state by walking the listener, which is fragile across imports.

### D2. Reuse `WaitForGa2TaskFinish` as-is
Already accepts a caller-supplied `timeout time.Duration`, polls `DescribeTaskResult`, treats `SUCCESS` as terminal. No rule-specific behavior is needed.

### D3. Add `DescribeGa2ForwardingRuleById` to the existing `Ga2Service`
Pattern matches the previously-added GA2 helpers:
- Build the request **outside** the for-loop. Set `request.GlobalAcceleratorId`, `request.ListenerId`, `request.ForwardingPolicyId` once.
- Page size = `100` (documented maximum), passed as a literal — no new package-level constant.
- Strict-equals on `*item.ForwardingRuleId == ruleId` (and double-check the parent IDs for paranoia).
- Returns `(nil, nil)` when not found; the resource layer treats this as "deleted out of band" and calls `d.SetId("")`.

Note: `DescribeForwardingRuleRequest` lacks a `Filters` field (unlike `DescribeListeners` / `DescribeEndpointGroups`); pagination is the only filter mechanism.

### D4. Async retry topology
Every SDK call wrapped in `resource.Retry(timeoutScope, func() *resource.RetryError { ... })`:
- Read paths: `tccommon.ReadRetryTimeout`.
- Write paths (Create / Modify / Delete): `tccommon.WriteRetryTimeout`.
- Async polling (after the write succeeds): `d.Timeout(schema.TimeoutCreate|Update|Delete)` passed into `WaitForGa2TaskFinish`.
- Resource-level `Timeouts` block defaults to **5 minutes** for Create/Update/Delete (matches the other GA2 resources).

### D5. Schema parity with `CreateForwardingRuleRequest`
Mapping (every CreateForwardingRule input field appears, no extras, no renames):

| Schema field | Type | Required? | ForceNew? | Source SDK field |
|---|---|---|---|---|
| `global_accelerator_id` | `TypeString` | Required | **Yes** | `GlobalAcceleratorId` |
| `listener_id` | `TypeString` | Required | **Yes** | `ListenerId` |
| `forwarding_policy_id` | `TypeString` | Required | **Yes** | `ForwardingPolicyId` |
| `rule_conditions` | `TypeSet`, nested block | Required | No | `RuleConditions []*RuleCondition` |
| └ `rule_condition_type` | `TypeString` | Required | — | `RuleConditionType` |
| └ `rule_condition_value` | `TypeSet` of String | Required | — | `RuleConditionValue []*string` |
| `rule_actions` | `TypeSet`, nested block | Required | No | `RuleActions []*RuleAction` |
| └ `rule_action_type` | `TypeString` | Required | — | `RuleActionType` |
| └ `rule_action_value` | `TypeString` | Required | — | `RuleActionValue` |
| `origin_headers` | `TypeSet`, nested block | Optional+Computed | No | `OriginHeaders []*OriginHeader` |
| └ `key` | `TypeString` | Required | — | `Key` |
| └ `value` | `TypeString` | Required | — | `Value` |
| `enable_origin_sni` | `TypeBool` | Optional+Computed | No | `EnableOriginSni` |
| `origin_sni` | `TypeString` | Optional+Computed | No | `OriginSni` |
| `origin_host` | `TypeString` | Optional+Computed | No | `OriginHost` |

Computed-only fields (not in `CreateForwardingRuleRequest`, surfaced from `ForwardingRuleSet`):
- `forwarding_rule_id` (string) — also stored as the 4th segment of `d.Id()`.

### D6. ForceNew choices justified by the API
- `global_accelerator_id`, `listener_id`, `forwarding_policy_id`: `ModifyForwardingRule` carries them but only as *identifiers* of which rule to modify; it cannot move a rule across these boundaries. ForceNew prevents users from accidentally requesting a relocation that the API silently rejects.
- All other input fields are updatable in place via `ModifyForwardingRule`.

### D7. `rule_conditions` / `rule_actions` / `origin_headers` use `TypeSet`
The cloud API treats these as unordered logical collections (the rule semantics depend on `Type`+`Value`, not on slice order). Using `TypeList` would create spurious diffs on rotation. We model them as `TypeSet`; build helpers convert via `(*schema.Set).List()`. This mirrors the `endpoint_configurations` decision in `tencentcloud_ga2_endpoint_group` and the `server_certificates` / `client_ca_certificates` decision in `tencentcloud_ga2_listener`.

`rule_condition_value` (the list of values for a single condition) is also modeled as `TypeSet`.

### D8. Update path semantics
Per `ModifyForwardingRuleRequest`, every body field (`RuleConditions`, `RuleActions`, `OriginHeaders`, `EnableOriginSni`, `OriginSni`, `OriginHost`) is updatable. The Update function:
- Skips the SDK call entirely if no body field changed (the 4 ID fields are ForceNew, so they cannot trigger Update directly).
- Always populates the 4 mandatory ID fields (`GlobalAcceleratorId`, `ListenerId`, `ForwardingPolicyId`, `ForwardingRuleId`) on the request.
- Forwards each body field whose schema getter returns a non-zero value.
- Awaits the returned `TaskId` via `WaitForGa2TaskFinish(ctx, taskId, d.Timeout(schema.TimeoutUpdate))`.

### D9. Read response shape inconsistency
`ForwardingRuleSet` exposes `RuleCondition` / `RuleAction` (singular field names), while the request structs use `RuleConditions` / `RuleActions` (plural). The flatten helper bridges that delta when populating the schema; this is purely an SDK quirk surfaced to the caller, not a Terraform-level concern.

### D10. File layout
Single file: `resource_tc_ga2_forwarding_rule.go` — schema + Create/Read/Update/Delete + ID parser + build/flatten helpers in that order. Service-level helper lives in the existing `service_tencentcloud_ga2.go`. Matches the user's strict feedback to avoid `_crud.go` / `_helpers.go` splits.

### D11. `make doc` flow + `provider.md` registration
Per the established workflow, the resource markdown lives at `tencentcloud/services/ga2/resource_tc_ga2_forwarding_rule.md`. The website file at `website/docs/r/ga2_forwarding_rule.html.markdown` is **never** hand-edited; it is regenerated by `make doc`. For `make doc` to discover the new resource, we must also append `tencentcloud_ga2_forwarding_rule` to the `Global Accelerator(GA2)` Resources section in `tencentcloud/provider.md`.

## Risks / Trade-offs

- **[Risk]** `DescribeForwardingRule` is keyed by `(gaId, listenerId, policyId)` and lacks a per-rule filter. → **Mitigation**: paginate with `Limit=100` and strict-equal on `ForwardingRuleId` in the helper. For policies with very large rule populations the page count grows linearly; this is acceptable because the API's documented maximum is the same `100` for all callers.
- **[Risk]** `ForwardingPolicyId` is not creatable through Terraform (no SDK API), so users must source it from the cloud (e.g. via the listener-level describe). → **Mitigation**: documented in the resource markdown; users can keep the value as a HCL variable or pull it from the listener resource attributes once we add `forwarding_policy_id` exposure to `tencentcloud_ga2_listener` in a future change. Out of scope for this change.
- **[Risk]** SDK Request uses plural names (`RuleConditions` / `RuleActions`) while the Set struct uses singular (`RuleCondition` / `RuleAction`). → **Mitigation**: flatten helper translates explicitly; covered by a Read scenario in the spec.
- **[Trade-off]** `forwarding_policy_id` ForceNew means users cannot move a rule between policies in place; they must destroy + recreate. This matches the API's hard constraint and is explicitly documented.

## Migration Plan

This is purely additive. No state migration required:
1. Land the new resource + service helper + provider registration on the `feat/ga2_nr` branch.
2. After release, users opt in by adding `resource "tencentcloud_ga2_forwarding_rule" "x" { ... }` to their config.

Rollback: pure revert of the new files + the `provider.go` and `provider.md` registration lines; no state mutations to undo.

## Open Questions

- None blocking. All required SDK surface area is vendored; the spec-level decisions (4-segment composite ID, ForceNew set, async polling reuse, TypeSet for rule collections) are determinate.
