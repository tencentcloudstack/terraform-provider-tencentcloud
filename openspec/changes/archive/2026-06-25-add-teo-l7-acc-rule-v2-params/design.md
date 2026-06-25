## Context

The `tencentcloud_teo_l7_acc_rule_v2` resource already exists at `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_v2.go` with the following schema fields: `zone_id`, `status`, `rule_name`, `description`, `branches`, `rule_id` (Computed), `rule_priority` (Computed).

The resource uses the TEO v20220901 SDK APIs:
- **Create**: `CreateL7AccRules` - takes `ZoneId` and `Rules` ([]*RuleEngineItem)
- **Read**: `DescribeL7AccRules` - takes `ZoneId` and `Filters` (with rule-id filter)
- **Update**: `ModifyL7AccRule` - takes `ZoneId` and `Rule` (*RuleEngineItem)
- **Delete**: `DeleteL7AccRules` - takes `ZoneId` and `RuleIds`

The `RuleEngineItem` struct contains: `Status`, `RuleId`, `RuleName`, `Description`, `Branches`, `RulePriority`.

The resource uses a composite ID format: `{zone_id}#{rule_id}` with `tccommon.FILED_SP` as separator.

## Goals / Non-Goals

**Goals:**
- Ensure the `tencentcloud_teo_l7_acc_rule_v2` resource has complete parameter coverage matching the cloud API
- Maintain backward compatibility with existing Terraform configurations
- Follow established patterns in the codebase (retry logic, error handling, composite IDs)

**Non-Goals:**
- Changing the existing schema structure or field types
- Adding support for batch rule operations (the resource manages a single rule)
- Modifying the `tencentcloud_teo_l7_acc_rule` (v1) resource

## Decisions

1. **Schema Design**: Use flat schema with individual fields (`zone_id`, `status`, `rule_name`, `description`, `branches`) rather than a nested `rules` block. This matches the existing implementation pattern and is more user-friendly for a single-rule resource.
   - Alternative: Nested `rules` block wrapping all fields - rejected because it adds unnecessary nesting for a single-rule resource and would break backward compatibility.

2. **Composite ID**: Use `{zone_id}#FILED_SP#{rule_id}` format for resource ID, enabling import support.
   - This is the established pattern in the codebase for resources with multiple identifying fields.

3. **Create API Mapping**: Construct a single `RuleEngineItem` from individual schema fields and wrap in a slice for the `CreateL7AccRules` API's `Rules` parameter.
   - The API accepts a list but we only create one rule per resource instance.

4. **Read API Mapping**: Use `DescribeL7AccRules` with `Filters` containing `rule-id` filter to fetch the specific rule by ID.

5. **Error Handling**: Use `tccommon.RetryError()` for all API calls with appropriate retry timeouts (`WriteRetryTimeout` for CUD, `ReadRetryTimeout` for R).

## Risks / Trade-offs

- [Risk] API returns empty response during transient failures → Mitigation: Check for nil response and empty rules list before setting ID to empty; log the resource ID before clearing it.
- [Risk] Create API returns empty RuleIds → Mitigation: Check response, RuleIds length, and nil pointer before using the returned ID; return NonRetryableError if empty.
