## Context

The `tencentcloud_teo_l7_acc_rule_v2` resource manages L7 acceleration rules for TencentCloud EdgeOne (TEO). It currently exposes individual fields from a single `RuleEngineItem` (status, rule_name, description, branches, rule_id, rule_priority) at the top level of the resource schema.

The `DescribeL7AccRules` API returns a `Rules` field of type `[]*RuleEngineItem` in its response. The requirement is to expose this full `Rules` list as a new computed attribute on the resource, allowing users to see the complete rules information returned by the API.

Current resource file: `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_v2.go`
Extension file with helper functions: `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_extension.go`

## Goals / Non-Goals

**Goals:**
- Add a new `Computed` attribute `rules` to the `tencentcloud_teo_l7_acc_rule_v2` resource that exposes the full `Rules` list from `DescribeL7AccRules` response.
- Maintain backward compatibility with existing resource configurations.
- Include proper unit tests for the new attribute.

**Non-Goals:**
- Changing the existing top-level fields (status, rule_name, description, branches, rule_id, rule_priority).
- Making the `rules` attribute writable (it is read-only/Computed).
- Modifying the Create, Update, or Delete logic.

## Decisions

1. **Attribute type**: The `rules` attribute will be a `TypeList` of `schema.Resource` with `Computed: true`. Each element represents a `RuleEngineItem` with fields: `status`, `rule_id`, `rule_name`, `description`, `branches`, `rule_priority`.

   *Rationale*: This mirrors the API response structure directly. Using `TypeList` preserves ordering which is important for rule priority.

2. **Reuse existing schema helpers**: The `branches` sub-schema within each rule element will reuse the existing `TencentTeoL7RuleBranchBasicInfo` helper function from the extension file, but with all fields set to `Computed: true` since this is a read-only attribute.

   *Rationale*: Consistency with existing code patterns and avoiding duplication.

3. **Read function update**: The Read function will populate the `rules` attribute by flattening the entire `respData.Rules` slice from the `DescribeL7AccRules` response.

   *Rationale*: The data is already fetched by the existing Read function; we just need to set the additional attribute.

4. **Flatten helper**: Create a dedicated flatten function to convert `[]*RuleEngineItem` into the format expected by Terraform's `d.Set()`.

   *Rationale*: Keeps the Read function clean and makes the flatten logic testable.

## Risks / Trade-offs

- [Risk] Large nested schema increases state size → Acceptable since this is a computed attribute and the data is already being fetched.
- [Risk] Schema complexity from deeply nested branches → Mitigated by reusing existing helper functions for branch schema definition.
