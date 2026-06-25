## Context

The `tencentcloud_teo_l7_acc_rule_v2` resource currently maps the `RuleEngineItem` structure's sub-fields (`status`, `rule_name`, `description`, `branches`, `rule_id`, `rule_priority`) individually in the Terraform schema. The `CreateL7AccRules` API accepts `request.Rules` (type `[]*RuleEngineItem`) and the `DescribeL7AccRules` API returns `response.Rules` (type `[]*RuleEngineItem`). However, there is no direct Terraform schema field that mirrors the complete `Rules` list from the API response.

This design adds a new computed `rules` field that faithfully represents the `response.Rules` output, enabling users to access the full structured rule data without needing to assemble it from individual sub-fields.

## Goals / Non-Goals

**Goals:**
- Add a new computed `rules` parameter (TypeList of `RuleEngineItem` blocks) to the `tencentcloud_teo_l7_acc_rule_v2` resource
- The `rules` field SHALL be populated in the Read function from `response.Rules`
- The `rules` field SHALL map `request.Rules` on the Create path via the existing sub-field construction logic

**Non-Goals:**
- No changes to existing schema fields (`zone_id`, `status`, `rule_name`, `description`, `branches`, `rule_id`, `rule_priority`)
- No changes to Modify/Delete API call logic
- No new API dependencies or SDK upgrades

## Decisions

### Decision 1: Make `rules` a Computed-only field
The `rules` field will be `TypeList`, `Computed: true`, and **not** `Optional` or `Required`. This avoids conflicts with the existing sub-fields (`status`, `rule_name`, `description`, `branches`) that already construct the `RuleEngineItem` for Create/Modify requests. Users continue to use the existing sub-fields for input, and the `rules` field serves as a unified computed output.

**Alternatives considered:**
- Making `rules` Optional → would create conflicting input paths with existing sub-fields
- Making `rules` both Optional and Conflicted with sub-fields → overly complex for the use case

### Decision 2: Use the existing `RuleEngineItem` structure directly
The `rules` element schema SHALL mirror the fields of `RuleEngineItem`: `rule_id` (Computed), `status` (Computed), `rule_name` (Computed), `description` (Computed), `branches` (Computed), `rule_priority` (Computed). All sub-fields within the list element are Computed since they come from the API response.

### Decision 3: Reuse existing `TencentTeoL7RuleBranchBasicInfo` for branches
The `branches` sub-field within each `rules` element will reuse the existing helper function `TencentTeoL7RuleBranchBasicInfo(1)` to avoid code duplication.

## Risks / Trade-offs

- **Risk**: The `rules` field may contain redundant data alongside existing top-level fields → **Mitigation**: Document clearly that `rules` is computed-only and the top-level fields remain the primary input mechanism
- **Risk**: Large `rules` output may impact state file size → **Mitigation**: The `rules` field mirrors existing API behavior and contains the same data as the individual fields already present in state
