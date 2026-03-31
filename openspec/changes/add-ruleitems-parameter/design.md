## Context

The `tencentcloud_teo_rule_engine` resource currently uses the `DescribeRules` API to fetch rule information. However, the API returns a `RuleItems` array (a list of rule items), but the current implementation only processes and returns the first matching rule item based on `rule_id`. This creates a data access limitation where users cannot see all available rule items returned by the API.

**Current State:**
- Service method `DescribeTeoRuleEngineById` returns `*teo.RuleItem` (single item)
- The method filters through `response.Response.RuleItems` to find a specific rule by ID
- Only one rule item is exposed to Terraform state

**Constraints:**
- Must maintain backward compatibility (cannot break existing TF configurations)
- Cannot modify existing schema except by adding Optional/Computed fields
- Must follow Terraform Plugin SDK v2 patterns
- Must handle potential API rate limits and retries

## Goals / Non-Goals

**Goals:**
- Expose the complete `RuleItems` array from the `DescribeRules` API response as a computed field in the resource
- Maintain backward compatibility with existing Terraform configurations
- Ensure the existing single rule item behavior continues to work when querying by `rule_id`
- Follow the established code patterns in the Teo service implementation

**Non-Goals:**
- Modifying the behavior of create/update/delete operations
- Changing the resource ID format or structure
- Adding new API calls or external dependencies
- Modifying other TEO resources or data sources

## Decisions

### 1. Add Computed Field to Schema
**Decision:** Add a new computed field `rule_items` to the resource schema to expose the complete list of rule items.

**Rationale:** Since the API already returns the complete `RuleItems` array, adding it as a computed field provides users access to all available data without requiring any changes to their existing configurations. Computed fields are read-only and don't affect the resource's create/update behavior.

**Alternatives Considered:**
- **Create a new data source:** This would duplicate code and complicate the user experience. Since the data is already being fetched, exposing it directly is more efficient.
- **Modify existing schema fields:** This would break backward compatibility and is not allowed by the constraints.

### 2. Update Service Method Signature
**Decision:** Keep `DescribeTeoRuleEngineById` returning a single `*teo.RuleItem` for backward compatibility, but add a new service method `DescribeTeoRuleEngineItems` that returns the complete array.

**Rationale:** This approach ensures existing code continues to work without modification while providing a new method to fetch the complete array. It follows the single responsibility principle and maintains clear separation of concerns.

**Alternatives Considered:**
- **Modify the existing method to return both:** This would require changing all call sites and could introduce bugs.
- **Always return the array:** This would break backward compatibility with existing consumers expecting a single item.

### 3. Schema Structure for `rule_items`
**Decision:** Define `rule_items` as a `TypeList` with nested schema matching the existing `rules` field structure in the resource.

**Rationale:** This aligns with the existing schema patterns in the resource and provides a consistent user experience. The structure mirrors the Terraform resource's `rules` field, making it familiar to users.

## Risks / Trade-offs

### Risk: Increased Memory Usage
**Risk:** Fetching and storing all rule items in Terraform state could increase memory usage for zones with many rules.

**Mitigation:** The `rule_items` field is computed and read-only. It doesn't affect the resource's create/update operations. Terraform state will only store what the API returns, which is bounded by the API's own limits.

### Risk: Schema Complexity
**Risk:** Adding a complex nested list structure to the schema could make the resource harder to understand.

**Mitigation:** The structure mirrors the existing `rules` field, which users are already familiar with. Clear documentation will be provided in the resource's markdown file.

### Trade-off: Single vs. Multiple API Calls
**Trade-off:** We're adding another service method, but it uses the same API call with different processing logic.

**Rationale:** This is acceptable because the alternative (modifying existing behavior) would break backward compatibility. The overhead of maintaining an additional method is minimal compared to the benefit of maintaining stability.

## Migration Plan

### Deployment Steps
1. Add the new `DescribeTeoRuleEngineItems` service method to `service_tencentcloud_teo.go`
2. Update the resource schema to include the `rule_items` computed field
3. Modify `resourceTencentCloudTeoRuleEngineRead` to populate the new field
4. Update resource documentation to describe the new field
5. Run existing tests to ensure backward compatibility
6. Add acceptance tests for the new field

### Rollback Strategy
If issues arise:
1. Remove the `rule_items` field from the schema (optional removal, won't affect existing state)
2. Remove or comment out the code that populates the field
3. Revert the service method changes
4. No data migration is needed since the field is computed and read-only

### Backward Compatibility
- Existing Terraform configurations will continue to work without modification
- The resource's create/update/delete operations are not affected
- Existing state files will be compatible (no state migration required)
- Users who don't use the new field will see no impact

## Open Questions

None at this time. The requirements are clear and the implementation approach follows established patterns in the codebase.
