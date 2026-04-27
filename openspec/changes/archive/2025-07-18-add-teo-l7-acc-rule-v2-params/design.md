## Context

The `tencentcloud_teo_l7_acc_rule_v2` resource manages TEO (TencentCloud EdgeOne) L7 acceleration rules. Currently, the resource schema exposes a `rule_id` computed field that stores only the first rule ID from the `CreateL7AccRules` API response. However, the API's response (`response.Response.RuleIds`) returns a list of all rule IDs created, which is not fully exposed to Terraform users.

The existing resource already handles:
- `zone_id` (required, ForceNew) - the site ID
- `status`, `rule_name`, `description`, `branches` - rule configuration fields
- `rule_id` (computed) - single rule ID
- `rule_priority` (computed) - rule priority

The `CreateL7AccRules` API response contains `RuleIds []*string` which represents all rule IDs created in the batch. This needs to be exposed as a new `rule_ids` computed schema field.

## Goals / Non-Goals

**Goals:**
- Add a `rule_ids` computed attribute (TypeList of TypeString) to the `tencentcloud_teo_l7_acc_rule_v2` resource schema
- Populate `rule_ids` from the `CreateL7AccRules` API response in the Create function
- Populate `rule_ids` in the Read function using data from the `DescribeL7AccRules` API response
- Update the resource documentation (.md file)
- Add unit tests for the new `rule_ids` field

**Non-Goals:**
- Modifying the existing `rule_id` computed field behavior
- Changing the resource's composite ID format (zone_id + FILED_SP + rule_id)
- Adding new input parameters beyond `rule_ids`
- Modifying the CreateL7AccRules request structure

## Decisions

1. **`rule_ids` as TypeList of TypeString**: The API returns `RuleIds []*string`, so `rule_ids` SHALL be defined as `schema.TypeList` with `Elem: &schema.Schema{Type: schema.TypeString}` and `Computed: true`. This aligns with the API's list-based response.

2. **Populating `rule_ids` in Create**: After the `CreateL7AccRules` API call succeeds, extract all values from `result.Response.RuleIds` and flatten them for `d.Set("rule_ids", ...)`. The existing logic already reads `result.Response.RuleIds[0]` for the composite ID.

3. **Populating `rule_ids` in Read**: In the Read function, after calling `DescribeTeoL7AccRuleById`, iterate over `respData.Rules` and collect the `RuleId` from each `RuleEngineItem` to build the `rule_ids` list. Since the current query filters by a single rule-id, the list will contain one element matching the current rule, consistent with the Create response.

4. **No change to Update/Delete functions**: The `rule_ids` field is computed-only and does not need to be handled in the Update or Delete functions.

## Risks / Trade-offs

- [Risk] The `DescribeL7AccRules` API returns rules filtered by rule-id, so `rule_ids` in Read will only contain the current rule's ID, not necessarily matching the full Create response if the API ever returns multiple IDs → Mitigation: This is acceptable since the current implementation only creates one rule per resource, and the behavior is consistent with the existing `rule_id` field.
