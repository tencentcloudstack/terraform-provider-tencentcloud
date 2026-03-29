## Proposal for teo-l7-acc-rule-totalcount Integration

### Summary

Integrate the `TotalCount` parameter from the `DescribeL7AccRules` API response into the `tencentcloud_teo_l7_acc_rule` data source schema. This will expose the total count of rules that match the query criteria to end users.

### Background

The `DescribeL7AccRules` API returns a `TotalCount` field that indicates the total number of L7 access control rules matching the query conditions. Currently, this field is not exposed in the Terraform provider's data source schema, which limits users' ability to understand the full scope of their query results without making additional API calls or manual calculations.

### Motivation

Exposing the `TotalCount` parameter provides several benefits:

1. **Improved User Experience**: Users can see the total count of rules without additional API calls
2. **Better Query Understanding**: Users can better understand the scope and completeness of their query results
3. **Pagination Support**: Enables users to implement pagination logic based on total count
4. **API Alignment**: Brings the Terraform provider data source in alignment with the underlying API's capabilities

### Goals

- Add `total_count` field (Computed type) to the `tencentcloud_teo_l7_acc_rule` data source schema
- Parse the `TotalCount` value from the API response and populate it in the schema
- Update unit tests to verify the `total_count` field is correctly populated
- Ensure backward compatibility - no breaking changes to existing configurations

### Non-Goals

- Modify the API request parameters or query logic
- Change how other fields are processed
- Modify state management or schema structure beyond adding the computed field
- Add pagination functionality (only expose the count, not implement pagination logic)

### Change Type

This is a **feature addition** (non-breaking change). It adds a computed field to an existing data source without modifying any existing behavior or requiring user action.

### Implementation Approach

1. **Schema Update**: Add `total_count` as a computed integer field to the data source schema
2. **Response Parsing**: In the `Read` method, extract the `TotalCount` value from the API response and set it in the schema
3. **Testing**: Add unit test cases to verify the field is correctly populated from API responses

### Impact Assessment

**Positive Impact:**
- Users gain visibility into total rule count
- Aligns with API capabilities
- No breaking changes

**Potential Concerns:**
- None identified - this is a purely additive change

### Success Criteria

- The `total_count` field appears in the data source schema and is populated correctly from API responses
- All existing tests continue to pass
- New tests validate the `total_count` field behavior
- Documentation is updated to reflect the new field

### Dependencies

None - this change is self-contained and does not depend on other changes.

### Risks and Mitigations

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| API response structure changes | Low | Medium | Add error handling for missing or null values |
| Field name conflicts with existing schema | Very Low | High | Verify field name uniqueness before implementation |
| Performance impact from additional parsing | Very Low | Low | Minimal impact - just reading one additional field from existing response |

### Rollout Plan

1. Implement the schema change
2. Add response parsing logic
3. Update tests
4. Update documentation
5. Release with provider version update

### Rollback Plan

If issues arise, the change can be reverted by removing the `total_count` field from the schema and response parsing logic. No state migration is needed since this is a computed field.

### Documentation Updates

- Update data source documentation to include the new `total_count` field
- Add examples showing how to use the field in queries
- Update changelog to document the new feature

### Testing Strategy

**Unit Tests:**
- Test normal case with valid TotalCount value
- Test edge case with TotalCount = 0
- Test with missing TotalCount field (ensure graceful handling)

**Integration Tests:**
- Verify the field is correctly populated from actual API responses
- Test that existing functionality remains unchanged

### Open Questions

None at this time.
