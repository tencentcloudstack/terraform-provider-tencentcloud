## Context

The Tencent Cloud GA2 (Global Accelerator V2) product exposes an ACL model:

```
GlobalAccelerator
└── GlobalAcceleratorAclPolicy  (default action + status)
    └── GlobalAcceleratorAclRule  (protocol / port / source_cidr_block / policy / description)
```

This provider already ships `tencentcloud_ga2_global_accelerator_acl_rule` — a **single-rule** resource whose schema mirrors one `GlobalAcceleratorAclRuleSet` element and whose CRUD maps 1:1 onto `CreateGlobalAcceleratorAclRule` (one `AclEntries` element) / `ModifyGlobalAcceleratorAclRule` / `DeleteGlobalAcceleratorAclRule`.

The cloud API, however, is **batch-oriented**:
- `CreateGlobalAcceleratorAclRule` accepts `AclEntries []*AclEntries` (a list) and returns `GlobalAcceleratorAclRuleIds []*string` (a list of created IDs).
- `DeleteGlobalAcceleratorAclRule` accepts `GlobalAcceleratorAclRuleIds []*string` (a list to delete in one call).
- `DescribeGlobalAcceleratorAclRules` returns the **entire** `GlobalAcceleratorAclRuleSet` for a given `GlobalAcceleratorAclPolicyId`.

So a set-style resource that owns the full rule collection under one ACL policy is the natural fit for the API shape, and it lets users declare the whole rule set in one block instead of one resource per rule.

The vendored SDK at `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ga2/v20250115/` already exposes all four calls. The provider already has `Ga2Service.WaitForGa2TaskFinish(ctx, taskId, timeout)` (polls `DescribeTaskResult` until `Status == SUCCESS`) and `Ga2Service.DescribeGa2GlobalAcceleratorAclRuleById(ctx, policyId, ruleId)` (single-rule lookup). The connectivity binding `UseGa2V20250115Client` is already wired.

## Goals / Non-Goals

**Goals:**
- Provide lifecycle management (create / read / update / delete / import) of the **full set** of ACL rules under one GA2 ACL policy as a single Terraform resource `tencentcloud_ga2_global_accelerator_acl_rule_set`.
- Reuse the batch-friendly API calls: one `CreateGlobalAcceleratorAclRule` call for all new entries, one `DeleteGlobalAcceleratorAclRule` call for all removed entries, and per-rule `ModifyGlobalAcceleratorAclRule` for changed entries (the modify API only accepts a single rule at a time).
- All async writes wait for `Status == SUCCESS` on the returned `TaskId` via the existing `Ga2Service.WaitForGa2TaskFinish`, so the cloud state is consistent before Terraform proceeds.
- Code style matches the existing ga2 resources (single-file layout, `resource.Retry` on every SDK call, defensive nil checks, `defer tccommon.LogElapsed/InconsistentCheck`).
- Coexist with the existing single-rule resource — no state or schema changes to it.

**Non-Goals:**
- Modifying or deprecating the existing `tencentcloud_ga2_global_accelerator_acl_rule` single-rule resource.
- Managing the ACL policy itself (`tencentcloud_ga2_global_accelerator_acl_policy` already exists).
- Datasource implementation (`data_source_tc_ga2_global_accelerator_acl_rule_set`); this change is resource-only.
- Tag management — ACL rules do not carry tags in the GA2 API.

## Decisions

### D1. Composite resource ID = `GlobalAcceleratorId#GlobalAcceleratorAclPolicyId`
Why: the rule set is uniquely identified by the (GA, policy) pair — exactly one rule-set resource exists per ACL policy. Using the composite ID lets the Read path recover both identifiers from `d.Id()` and keeps import to a single token.
Alternative considered: bare `GlobalAcceleratorAclPolicyId`. Rejected — the GA2 create/delete/modify APIs all require `GlobalAcceleratorId` as well as the policy ID, so storing only the policy ID would force the Read path to look up the GA ID through an extra describe call. The composite ID avoids that.

### D2. `acl_entries` is a `TypeList` whose elements are flattened rule fields
Each element exposes: `protocol`, `port`, `source_cidr_block`, `policy`, `description` (user-editable) plus `global_accelerator_acl_rule_id` (computed, server-assigned). Per the user's rule against nesting list data one level deeper, the rule fields are placed directly on the element rather than wrapping them in a `rule { ... }` sub-block.
- `acl_entries` is `Required` with `MinItems: 0` (an empty set is valid and means "delete all rules under the policy").
- The computed `global_accelerator_acl_rule_id` is populated on Create/Read and used as the join key during Update diffing.

### D3. Create = single batched `CreateGlobalAcceleratorAclRule` call
Assemble every `acl_entries` element into one `[]*ga2v20250115.AclEntries` slice and call `CreateGlobalAcceleratorAclRule` once. The response returns `GlobalAcceleratorAclRuleIds` in the same order as the input `AclEntries`, so we map `response.Response.GlobalAcceleratorAclRuleIds[i]` back onto `acl_entries[i].global_accelerator_acl_rule_id`. Then poll `TaskId` via `WaitForGa2TaskFinish`.
- Defensive checks: if `response.Response` is nil, or `GlobalAcceleratorAclRuleIds` is nil/empty, or `len(ids) != len(input)`, return `NonRetryableError` (per the user's rule to never write an empty/missing ID into state).
- `d.SetId(gaId + FILED_SP + policyId)` is set **after** the task succeeds (outside the retry block).

### D4. Read = `DescribeGlobalAcceleratorAclRules` paginated by `Limit=200`
Add a new service helper `Ga2Service.DescribeGa2GlobalAcceleratorAclRulesByPolicyId(ctx, policyId) ([]*ga2v20250115.GlobalAcceleratorAclRuleSet, error)` that loops with `Offset`/`Limit=200` (the documented maximum) and returns the full set. The resource Read calls this, then flattens the set into `acl_entries` (sorting by `global_accelerator_acl_rule_id` for deterministic state).
- If the API returns an empty set, the resource is **not** removed from state — an empty rule set is a valid desired state. Instead, `acl_entries` is set to an empty list.
- If `response.Response` is nil, the helper returns an error (not a silent empty result), to surface real API failures.

### D5. Update = diff-driven create / modify / delete
On `d.HasChange("acl_entries")`, compare old vs. new lists keyed by `global_accelerator_acl_rule_id`:
- **New entries** (present in new, absent in old): batch-create via one `CreateGlobalAcceleratorAclRule` call with all new `AclEntries`; poll `TaskId`; map returned IDs back.
- **Removed entries** (present in old, absent in new): batch-delete via one `DeleteGlobalAcceleratorAclRule` call with all removed `GlobalAcceleratorAclRuleIds`; poll `TaskId`.
- **Changed entries** (same ID, different fields): call `ModifyGlobalAcceleratorAclRule` per entry; poll each `TaskId`.
Because `global_accelerator_acl_policy_id` and `global_accelerator_id` are `ForceNew`, changing them recreates the resource rather than entering Update.
Alternative: always delete-all then create-all on any change. Rejected — wasteful and disrupts rule IDs users may reference elsewhere.

### D6. Delete = single batched `DeleteGlobalAcceleratorAclRule` call
Collect every `global_accelerator_acl_rule_id` from `acl_entries` into one `GlobalAcceleratorAclRuleIds` slice and call `DeleteGlobalAcceleratorAclRule` once; poll `TaskId`.
- If `acl_entries` is empty at delete time, skip the API call (nothing to delete) and return nil.

### D7. Async retry topology
- Read path: `tccommon.ReadRetryTimeout` (wraps the describe call).
- Write paths (Create / batch-Modify / batch-Delete): `tccommon.WriteRetryTimeout` (wraps each SDK call).
- Task polling (after each write succeeds): `d.Timeout(schema.TimeoutCreate|Update|Delete)` passed into `WaitForGa2TaskFinish`.
- Resource-level `Timeouts` block defaults to **5 minutes** for Create/Update/Delete (matches the existing ga2 resources).

### D8. ForceNew choices
- `global_accelerator_id` → ForceNew (the create/delete APIs require it and it cannot be moved between GA instances).
- `global_accelerator_acl_policy_id` → ForceNew (rules belong to a specific policy).
- `acl_entries` → updatable in place (the whole point of the set resource is in-place diffing).
`task_id` is Computed and never user-set.

### D9. Importer = `schema.ImportStatePassthrough`
Import via the composite ID `GlobalAcceleratorId#GlobalAcceleratorAclPolicyId`. The Read path re-derives both IDs from `d.Id()`, lists all rules under the policy, and populates `acl_entries`. Matches the existing single-rule resource's import idiom.

### D10. File layout (single-file resource)
Per the user's standing rule ("不要拆分 _crud.go / _helpers.go"), the entire resource lives in one file:
- `resource_tc_ga2_global_accelerator_acl_rule_set.go`: schema + Create/Read/Update/Delete + `parseGa2GlobalAcceleratorAclRuleSetId` + build/flatten/diff helpers, in that order.
The service-layer helper `DescribeGa2GlobalAcceleratorAclRulesByPolicyId` lives in the existing `service_tencentcloud_ga2.go`.

## Risks / Trade-offs

- **[Risk]** `CreateGlobalAcceleratorAclRule` is documented as returning `GlobalAcceleratorAclRuleIds` in input order, but if the server ever reorders them the back-mapping onto `acl_entries[i]` would assign wrong IDs. → **Mitigation**: the `AclEntries` struct has no client-generated correlation ID, so input-order is the only available join key; we rely on it and document the assumption. If server behavior diverges, Read will self-correct on the next refresh because it matches by rule fields, not by list index.
- **[Risk]** Update performs potentially many sequential `ModifyGlobalAcceleratorAclRule` calls (one per changed rule), each followed by task polling. Large rule sets with many simultaneous changes could be slow. → **Mitigation**: the modify API has no batch form, so per-rule calls are unavoidable; the `Timeouts.Update` budget (default 5 min) bounds the total. Document that for large sets, wholesale replace via `acl_entries` list changes is faster (batch create/delete).
- **[Risk]** Importing a rule set whose policy has many rules produces a large `acl_entries` list in state. → **Mitigation**: this is inherent to the set resource model; the provider already handles large lists for other set-style resources.
- **[Trade-off]** An empty `acl_entries` list is a valid state (delete all rules under the policy) rather than triggering resource deletion. This is intentional — the resource owns the rule set, and "zero rules" is a valid set. The resource itself is only destroyed on `terraform destroy`.

## Migration Plan

This is purely additive. No state migration required:
1. Land the new resource + service helper + provider registration.
2. After release, users opt in by adding `resource "tencentcloud_ga2_global_accelerator_acl_rule_set" "x" { ... }` to their config.
3. Existing `tencentcloud_ga2_global_accelerator_acl_rule` single-rule resources continue to work unchanged — nothing is deprecated.

Rollback: pure revert of the new files + the `provider.go` registration line; no state mutations to undo.

## Open Questions

- None requiring user input. The SDK exposes all needed APIs; the proposal-level decisions (composite ID, batched create/delete, per-rule modify, empty-set validity) are fully determined by the API shape.
