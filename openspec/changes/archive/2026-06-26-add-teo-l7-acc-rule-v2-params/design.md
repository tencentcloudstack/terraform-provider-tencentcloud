## Context

The `tencentcloud_teo_l7_acc_rule_v2` resource manages TEO L7 acceleration rules through the cloud API. The resource currently supports `zone_id`, `rule_id`, `status`, `rule_name`, `description`, `branches`, and `rule_priority` parameters. The CRUD operations (Create, Read, Update, Delete) are already implemented using `CreateL7AccRules`, `DescribeL7AccRules`, `ModifyL7AccRule`, and `DeleteL7AccRules` API interfaces from the `teo/v20220901` SDK package.

The core data structure `RuleEngineItem` in the SDK contains: `Status`, `RuleId`, `RuleName`, `Description`, `Branches`, and `RulePriority`.

The current resource uses a composite ID (`zoneId + FILED_SP + ruleId`) and supports import via `ImportStatePassthrough`.

## Goals / Non-Goals

**Goals:**
- Ensure all CRUD API parameters are correctly mapped in the `tencentcloud_teo_l7_acc_rule_v2` resource
- Verify parameter consistency across CreateL7AccRules, ModifyL7AccRule, DeleteL7AccRules, and DescribeL7AccRules interfaces
- Maintain backward compatibility with existing Terraform configurations and state

**Non-Goals:**
- Adding new action types to the `branches.actions` schema (e.g., new operation names)
- Modifying the `RuleEngineItem` data structure or adding new SDK fields
- Changing the composite ID format or import behavior

## Decisions

1. **Parameter mapping approach**: All parameters listed in the requirement (`ZoneId`, `Rules`, `zone_id`, `rule_id`, `status`, `rule_name`, `description`, `branches`, `RuleIds`, `Values`) already exist in the current resource implementation. The Create function maps `zone_id` → `request.ZoneId` and individual fields → `RuleEngineItem` within `request.Rules`. The Modify function maps `zone_id` → `request.ZoneId` and fields → `request.Rule`. The Delete function maps `zone_id` → `request.ZoneId` and `rule_id` → `request.RuleIds`. The Describe (read) uses `Filters` with `Name: "rule-id"` and `Values: [ruleId]`.

2. **DescribeL7AccRules parameter mapping**: The requirement specifies `request.Values` → `rule_id`. In the SDK, `DescribeL7AccRulesRequest` uses `Filters []*Filter` where each `Filter` has `Name` and `Values` fields. The current implementation correctly uses `Filters` with `Name: "rule-id"` and `Values: []string{ruleId}`. This is the correct mapping since the SDK doesn't have a direct `Values` field at the request level.

3. **Backward compatibility**: All parameter additions are Optional or Computed, ensuring existing configurations remain valid.

## Risks / Trade-offs

- [Risk] Parameter naming inconsistency between Create (`ZoneId` → `ZoneId` in requirement vs `zone_id` in schema) → Mitigation: The schema uses snake_case as per Terraform convention; the mapping is correctly handled in the Create/Update functions
- [Risk] The DescribeL7AccRules API uses `Filters` not direct `Values` → Mitigation: Current implementation correctly uses `Filters` with `Name: "rule-id"` which maps to the `rule_id` parameter
