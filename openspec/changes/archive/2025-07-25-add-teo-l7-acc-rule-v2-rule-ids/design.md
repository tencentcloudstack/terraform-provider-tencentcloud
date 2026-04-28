## Context

The `tencentcloud_teo_l7_acc_rule_v2` resource manages TEO L7 acceleration rules via the `CreateL7AccRules` API. Currently, the resource only exposes a single `rule_id` (TypeString, Computed) by taking the first element from the `RuleIds` array in the API response. However, the `CreateL7AccRules` API returns `RuleIds []*string`, which is a list of all created rule IDs. The full list should be exposed as a computed `rule_ids` parameter so users can reference all created rules.

The existing resource schema includes: `zone_id` (Required, ForceNew), `status`, `rule_name`, `description`, `branches` (Optional), and computed fields `rule_id` and `rule_priority`.

## Goals / Non-Goals

**Goals:**
- Add a computed `rule_ids` parameter (TypeList of TypeString) to expose the full list of rule IDs returned by `CreateL7AccRules` API
- Populate `rule_ids` in both Create and Read functions
- Maintain backward compatibility - existing `rule_id` (singular) remains unchanged
- Add unit tests for the new parameter

**Non-Goals:**
- Changing or removing the existing `rule_id` parameter
- Modifying the Update or Delete logic (rule_ids is computed, not user-settable)
- Changing the resource ID format (remains `zone_id#rule_id`)

## Decisions

1. **Schema type for `rule_ids`**: Use `TypeList` with `Elem: &schema.Schema{Type: schema.TypeString}` and `Computed: true`. This matches the API response type `[]*string` and follows existing patterns in the codebase for computed list outputs.

2. **Setting `rule_ids` in Create**: After the `CreateL7AccRules` call succeeds, extract all rule IDs from `result.Response.RuleIds` and flatten to `[]string` for `d.Set("rule_ids", ...)`. This is done inside the retry block since it depends on a successful API call.

3. **Setting `rule_ids` in Read**: The `DescribeL7AccRules` API returns `RuleEngineItem` objects which each contain a `RuleId` field. After reading the rules, collect all `RuleId` values from the response and set `rule_ids`. Note: the Describe API returns the full rule objects, so we extract RuleId from each.

4. **Keep existing `rule_id`**: The `rule_id` (singular) parameter remains for backward compatibility. It continues to store the first rule ID from the list, which is used in the composite resource ID.

## Risks / Trade-offs

- [Risk: Duplicate data] The `rule_id` and `rule_ids[0]` will contain the same value → This is acceptable for backward compatibility; `rule_id` is the primary identifier used in the composite ID, while `rule_ids` provides the complete list.
- [Risk: DescribeL7AccRules returns all rules for a zone, not just the resource's rule] → The current Read implementation already filters by rule ID using the composite ID. We only set `rule_ids` from the matched rule, ensuring consistency with the resource's scope.
