# Design: Add TotalCount Parameter to tencentcloud_teo_l7_acc_rule

## Context

The Terraform Provider for TencentCloud currently exposes the `tencentcloud_teo_l7_acc_rule` data source, which queries the `DescribeL7AccRules` API from the Tencent Cloud EdgeOne (TEO) service. The API response includes a `TotalCount` field that indicates the total number of records matching the query criteria. This field is already available in the API response but is not exposed in the Terraform provider's data source schema.

Current state:
- Data source exists at `tencentcloud/services/teo/data_source_tc_teo_l7_acc_rule.go`
- API returns `TotalCount` in the response
- Provider does not expose this field to users

Constraints:
- Must maintain backward compatibility with existing Terraform configurations
- Should not modify existing resource schemas (only add new Optional/Computed fields)
- Must follow Terraform Plugin SDK v2 patterns
- Requires proper documentation updates

## Goals / Non-Goals

**Goals:**
- Expose the `TotalCount` field from the `DescribeL7AccRules` API response
- Enable users to access the total record count for pagination and data validation scenarios
- Maintain full backward compatibility with existing configurations
- Follow established patterns in the codebase for adding computed fields

**Non-Goals:**
- Modifying the `DescribeL7AccRules` API call parameters
- Changing the behavior of existing data source fields
- Adding any new API calls or external dependencies
- Implementing pagination logic in the provider (the user can implement this with the TotalCount information)

## Decisions

### Schema Definition Decision
- **Decision**: Add `TotalCount` as a computed (read-only) field in the data source schema
- **Rationale**: The field is already returned by the API and should be read-only since it's determined by the API response, not user input. This follows the Terraform pattern for exposing API metadata.

### Type Mapping Decision
- **Decision**: Map the API's `TotalCount` (int64) to Terraform's `schema.TypeInt`
- **Rationale**: The field represents a count value, which is an integer. Terraform's `schema.TypeInt` is the appropriate type for this data.

### Implementation Location Decision
- **Decision**: Update the `data_source_tc_teo_l7_acc_rule.go` file in the teo service directory
- **Rationale**: This follows the established file organization pattern where each data source has its own file under `tencentcloud/services/<service>/`.

### No State Migration Needed Decision
- **Decision**: No state migration required
- **Rationale**: We're only adding a new computed field to the data source, not modifying resources or existing schema fields. Data sources don't maintain state, so backward compatibility is guaranteed.

## Risks / Trade-offs

### Risk: API Field Name Change
- **Risk**: Tencent Cloud API might change the `TotalCount` field name or structure in future versions
- **Mitigation**: The provider will need to be updated when such API changes occur. This is a standard maintenance requirement for any provider that depends on external APIs.

### Risk: Large TotalCount Values
- **Risk**: The TotalCount field could potentially exceed Terraform's int64 range for very large datasets
- **Mitigation**: This is unlikely for typical TEO use cases. If encountered, we could document this limitation or consider using a string type in future updates.

### Trade-off: Documentation Maintenance
- **Trade-off**: Adding this field increases documentation maintenance burden
- **Mitigation**: The documentation update is straightforward and minimal. The benefit of providing this information to users outweighs the maintenance cost.

## Migration Plan

Since this change only adds a new computed field to a data source and does not modify resources or existing schema fields:

1. **No migration required** - Data sources don't maintain state, and we're only adding a new field
2. **Backward compatible** - Existing Terraform configurations will continue to work without modification
3. **Rollback strategy** - Simply revert the code changes if issues arise; no state cleanup needed

## Open Questions

None - The requirements are clear and the implementation is straightforward following established patterns in the codebase.
