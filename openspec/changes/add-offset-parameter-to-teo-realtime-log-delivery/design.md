## Context

The Terraform Provider for TencentCloud currently provides a resource `tencentcloud_teo_realtime_log_delivery` for managing individual realtime log delivery tasks. However, there is no data source to query and list multiple tasks, which is a common use case for users who need to:

- Discover all realtime log delivery tasks across zones
- Filter tasks by specific criteria (zone, task type, status, etc.)
- Paginate through large task lists efficiently
- Export task information for analysis or reporting

The underlying TencentCloud API `DescribeRealtimeLogDeliveryTasks` already supports these capabilities, including pagination with `Offset` and `Limit` parameters, filtering via `Filters` parameter, and sorting with `Order` and `Direction` parameters.

This design leverages existing patterns in the TEO service provider, such as the `tencentcloud_teo_zones` data source, which implements similar filtering and pagination features.

### Current State

- **Resource**: `tencentcloud_teo_realtime_log_delivery` - manages individual tasks via CRD operations
- **Service Layer**: `DescribeTeoRealtimeLogDeliveryById()` - fetches a single task by ID using filter
- **No Data Source**: Missing ability to query multiple tasks with pagination and filtering

### Constraints

- Must maintain backward compatibility (cannot modify existing resource schema)
- Must follow Terraform Plugin SDK v2 patterns and provider conventions
- Must use `helper.Retry()` for final consistency retry logic
- Must include proper error handling with `tccommon.LogElapsed()` and `tccommon.InconsistentCheck()`
- Must include acceptance tests with TF_ACC=1
- Must include documentation in `website/docs/` directory

## Goals / Non-Goals

**Goals:**

1. Add a new data source `tencentcloud_teo_realtime_log_delivery_tasks` to query and list realtime log delivery tasks
2. Support pagination with `Offset` and `Limit` parameters for efficient large dataset handling
3. Support filtering by zone_id, task_id, task_name, task_type, and other relevant fields
4. Support result ordering with `Order` and `Direction` parameters
5. Provide complete task information in the result set matching the API response model
6. Include comprehensive acceptance tests covering pagination, filtering, and edge cases
7. Include documentation with usage examples

**Non-Goals:**

1. Modifying the existing `tencentcloud_teo_realtime_log_delivery` resource
2. Adding write/update/delete operations (this is a read-only data source)
3. Implementing custom pagination logic beyond what the API provides
4. Supporting deprecated or experimental API parameters

## Decisions

### 1. Data Source Naming Convention

**Decision**: Use `tencentcloud_teo_realtime_log_delivery_tasks` (plural) for the data source name.

**Rationale**: Terraform naming conventions use plural forms for list data sources (e.g., `tencentcloud_vpc_instances`, `tencentcloud_cvm_instances`). This clearly distinguishes it from the singular resource `tencentcloud_teo_realtime_log_delivery`. The file naming will follow the pattern `data_source_tc_teo_realtime_log_delivery_tasks.go`.

**Alternatives Considered**:
- `tencentcloud_teo_realtime_log_delivery` (same as resource) - Confusing and violates naming conventions
- `tencentcloud_teo_realtime_log_delivery_list` - Less standard, "list" is implied by data source nature

### 2. Parameter Mapping to API

**Decision**: Map Terraform data source parameters directly to TencentCloud API parameters with minimal transformation.

**Rationale**: The DescribeRealtimeLogDeliveryTasks API already has well-designed parameters for pagination (`Offset`, `Limit`), filtering (`Filters`), and sorting (`Order`, `Direction`). Direct mapping reduces complexity, makes the code more maintainable, and aligns with user expectations if they're familiar with the API.

**Mapping**:
- `offset` → `Offset` (uint64)
- `limit` → `Limit` (uint64)
- `filters` → `Filters` ([]*AdvancedFilter)
- `order` → `Order` (string)
- `direction` → `Direction` (string)

### 3. Service Layer Function Design

**Decision**: Add a new function `DescribeTeoRealtimeLogDeliveryTasks()` in the service layer that returns a slice of task objects.

**Rationale**: The existing `DescribeTeoRealtimeLogDeliveryById()` function is specialized for single-task lookup with hardcoded filters. A new function is needed to support flexible querying with user-provided filters and pagination parameters. This follows the separation of concerns pattern where service layer handles API communication and retry logic.

**Function Signature**:
```go
func (me *TeoService) DescribeTeoRealtimeLogDeliveryTasks(
    ctx context.Context,
    filters []*teo.AdvancedFilter,
    offset uint64,
    limit uint64,
    order string,
    direction string,
) ([]*teo.RealtimeLogDeliveryTask, int64, error)
```

The function returns:
- Task slice
- Total count (for pagination awareness)
- Error

### 4. Filter Implementation Strategy

**Decision**: Implement filters as a list of filter objects following the pattern used in other TEO data sources (e.g., `tencentcloud_teo_zones`).

**Rationale**: This pattern is:
- Consistent with existing data sources in the TEO service
- Flexible enough to support multiple filter criteria
- Type-safe with proper schema validation
- Well-understood by users familiar with the provider

**Filter Schema**:
```go
"filters": {
    Type: schema.TypeList,
    Optional: true,
    Description: "Filter criteria...",
    Elem: &schema.Resource{
        Schema: map[string]*schema.Schema{
            "name": {
                Type: schema.TypeString,
                Required: true,
            },
            "values": {
                Type: schema.TypeSet,
                Required: true,
                Elem: &schema.Schema{
                    Type: schema.TypeString,
                },
            },
            "fuzzy": {
                Type: schema.TypeBool,
                Optional: true,
            },
        },
    },
}
```

### 5. Result Set Structure

**Decision**: Return task list as a computed schema field with nested resource schema matching the API response.

**Rationale**: Nested resource schema provides:
- Clear structure matching the API's RealtimeLogDeliveryTask model
- Type safety for all fields
- Proper Terraform state management
- Good documentation generation

The schema will include all relevant fields from the API response, such as:
- task_id
- zone_id
- task_name
- task_type
- delivery_status
- log_type
- area
- entity_list
- fields
- custom_fields
- And other task metadata

### 6. Error Handling Strategy

**Decision**: Use standard provider error handling patterns with proper retry logic and clear error messages.

**Rationale**:
- `helper.Retry()` handles transient failures and eventual consistency
- `tccommon.LogElapsed()` tracks performance metrics
- `tccommon.InconsistentCheck()` validates state consistency
- Error messages should be actionable and clearly indicate the issue

Error types to handle:
- API authentication failures (credentials)
- Rate limiting (retry with backoff)
- Invalid parameters (user error, don't retry)
- Network issues (transient, retry)
- Zone/task not found (return empty list, not an error)

### 7. Testing Strategy

**Decision**: Implement comprehensive acceptance tests covering all major use cases.

**Rationale**: Acceptance tests (TF_ACC=1) are essential for:
- Validating integration with the actual TencentCloud API
- Ensuring pagination works correctly
- Verifying filter behavior
- Testing edge cases (empty results, large offsets, etc.)
- Maintaining backward compatibility

**Test Scenarios**:
1. Basic query without filters
2. Pagination with offset and limit
3. Filter by zone_id
4. Filter by task_id
5. Filter by task_type
6. Filter by task_name with fuzzy matching
7. Ordering by different fields
8. Empty result handling
9. Large offset exceeding result count
10. Multiple filters combined

## Risks / Trade-offs

### Risk 1: API Behavior Changes

**Risk**: TencentCloud API may change parameter names, types, or behavior in future versions.

**Mitigation**: Use the vendored SDK version (`tencentcloud-sdk-go`) which pins the API contract. Regular updates to follow API changes. Document supported API version.

### Risk 2: Pagination Complexity

**Risk**: Users may misuse offset/limit parameters leading to inefficient queries or confusion.

**Mitigation**: Provide clear documentation with examples. Include reasonable defaults (limit=20, offset=0). Validate parameter values (non-negative integers). Document that results may change between queries if tasks are added/removed.

### Risk 3: Filter Performance

**Risk**: Complex filter queries with many conditions may be slow or hit API rate limits.

**Mitigation**: The API handles filtering on the server side, which is efficient. Document that complex queries may take longer. Implement rate limit handling with appropriate retry logic.

### Risk 4: Large Result Sets

**Risk**: Querying tasks across many zones with no filters could return thousands of results, impacting performance.

**Mitigation**: Encourage users to use filters via documentation. Consider implementing a default limit (e.g., 20) to prevent unbounded queries. Document that users should always use appropriate filtering and pagination.

### Risk 5: Backward Compatibility

**Risk**: This is a new feature, so no backward compatibility concerns. However, future changes could break existing configurations.

**Mitigation**: Follow semantic versioning for provider releases. Test thoroughly before releases. Provide migration guides if breaking changes are necessary in the future.

### Trade-off 1: Flexibility vs. Complexity

**Trade-off**: The filters design is very flexible (supports multiple filters with values and fuzzy matching) but adds schema complexity.

**Decision**: Flexibility is worth the complexity because it provides powerful querying capabilities that match the API's capabilities, and the pattern is consistent with other data sources in the provider.

### Trade-off 2: Total Count vs. Simpler API

**Trade-off**: Including total count in the response is useful for pagination awareness but requires additional parsing and state management.

**Decision**: Include total count because it's valuable for users to know the total number of available results and implement client-side pagination if needed. The API provides this information, so leveraging it is appropriate.

## Migration Plan

This is a new feature addition, so no migration is required for existing users. The implementation plan is:

1. **Create the data source implementation** in `data_source_tc_teo_realtime_log_delivery_tasks.go`
2. **Add service layer function** in `service_tencentcloud_teo.go`
3. **Register the data source** in the TEO service provider (`teo.go`)
4. **Create acceptance tests** in `data_source_tc_teo_realtime_log_delivery_tasks_test.go`
5. **Create documentation** in `data_source_tc_teo_realtime_log_delivery_tasks.md`
6. **Run acceptance tests** to validate the implementation
7. **Generate documentation** (if using docsgen tool)
8. **Review and merge** the changes

**Rollback Strategy**: Since this is a new data source (not modifying existing functionality), rolling back simply involves not releasing or disabling the data source in a subsequent release. Existing users' configurations will not be affected.

## Open Questions

1. **Default limit value**: Should we implement a default limit (e.g., 20) to prevent unbounded queries, or require users to always specify it?
   - **Recommendation**: Default to 20 if not specified, with a maximum of 100 (matching API constraints)

2. **Maximum offset value**: Should we validate or limit the maximum offset value to prevent users from specifying extremely large offsets?
   - **Recommendation**: Allow any non-negative integer value (validated as uint64), but document that specifying an offset larger than the result count returns an empty list

3. **Filter name validation**: Should we validate filter names against a whitelist of supported filters, or pass any filter name to the API and let the API handle errors?
   - **Recommendation**: Pass any filter name to the API to allow future API support for new filters without requiring provider updates. Document the currently supported filters.

4. **Task output fields**: Should we include all fields from the RealtimeLogDeliveryTask API response, or only a subset?
   - **Recommendation**: Include all relevant fields to provide complete information. This matches the pattern of other data sources and provides maximum value to users.
