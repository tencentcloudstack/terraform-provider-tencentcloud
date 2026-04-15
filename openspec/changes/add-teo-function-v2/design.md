## Context

This design document outlines the implementation of `tencentcloud_teo_function_v2` resource for managing TencentCloud EdgeOne (TEO) functions. The current codebase already has a similar resource `tencentcloud_teo_function` which serves as a reference for implementation patterns and best practices.

TEO functions are serverless JavaScript code that runs at the edge, allowing users to customize content delivery, security rules, and other edge behaviors. The V2 variant follows the same API contract as the existing function resource but is implemented as a separate resource to maintain backward compatibility and allow future API evolution.

Key constraints:
- Must follow existing code patterns in `tencentcloud/services/teo/` directory
- Must use the teo v20220901 SDK package from vendor
- Must support standard Terraform CRUD operations
- Must handle asynchronous operations with state refresh
- Must implement retry logic for transient failures
- Must maintain backward compatibility with existing resources

## Goals / Non-Goals

**Goals:**
- Create a new Terraform resource `tencentcloud_teo_function_v2` for managing TEO functions
- Implement complete CRUD operations using CreateFunction, DescribeFunctions, ModifyFunction, and DeleteFunction APIs
- Support all function attributes: zone_id, function_id, name, remark, content, domain, create_time, update_time
- Implement proper resource ID format: `zone_id#function_id`
- Handle asynchronous operations with state polling
- Provide comprehensive error handling and logging
- Follow existing code patterns and conventions from the codebase

**Non-Goals:**
- Modify the existing `tencentcloud_teo_function` resource
- Implement function rule management (handled by separate resources)
- Implement function runtime environment configuration (handled by separate resources)
- Support function execution or debugging features
- Implement function versioning or rollback capabilities

## Decisions

### 1. Resource Implementation Pattern
**Decision:** Follow the exact pattern of `resource_tc_teo_function.go` implementation

**Rationale:**
- The existing function resource is well-tested and follows all best practices
- Consistency across similar resources reduces maintenance burden
- The pattern includes proper retry logic, error handling, and state refresh
- Using a proven pattern reduces risk of bugs

**Alternatives considered:**
- Create a completely new implementation: More flexibility but higher risk of inconsistencies
- Refactor to share code between function and function_v2: Would reduce duplication but increase complexity and potential for breaking changes

### 2. Resource ID Format
**Decision:** Use composite ID format `zone_id#function_id` with `tccommon.FILED_SP` separator

**Rationale:**
- Matches existing pattern used by `tencentcloud_teo_function`
- Allows importing existing functions
- Enables unique identification of functions across zones
- Follows the project's established conventions

**Alternatives considered:**
- Use only `function_id`: Simpler but loses context and prevents multi-zone support
- Use a JSON-encoded ID: More flexible but not following project conventions

### 3. Immutable vs Mutable Fields
**Decision:** Mark `name` and `zone_id` as immutable (`ForceNew: true`), allow `remark` and `content` to be mutable

**Rationale:**
- `zone_id` as ForceNew is required because it's part of the resource ID
- `name` as ForceNew follows the existing function resource pattern and API limitations
- `remark` and `content` are supported by ModifyFunction API, making them safe to update
- Prevents accidental recreation of resources when updating mutable fields

**Alternatives considered:**
- Allow all fields to be mutable: Would require recreating resources when name changes, which is not supported by the API

### 4. State Refresh Logic
**Decision:** Poll DescribeFunctions API until the `domain` field is populated

**Rationale:**
- Matches the existing function resource implementation
- Domain assignment indicates the function is fully provisioned and ready to use
- Provides a clear, observable state for the refresh mechanism
- Uses the same timeout configuration (10s delay, 3s min timeout, 600s max timeout)

**Alternatives considered:**
- Poll until function_id is returned: Too early, function may not be fully provisioned
- Use a simple sleep: Less reliable and may fail or wait too long
- Check multiple fields: More complex without clear benefit

### 5. Error Handling Strategy
**Decision:** Use resource.Retry with tccommon.WriteRetryTimeout for write operations and tccommon.ReadRetryTimeout for read operations

**Rationale:**
- Handles transient network issues and rate limiting automatically
- Follows the project's standard retry pattern
- Provides consistent timeout behavior across all operations
- Includes proper logging of retries and failures

**Alternatives considered:**
- Immediate return on first error: More fragile and less user-friendly
- Custom retry logic: More control but adds complexity and inconsistency

### 6. Service Layer Integration
**Decision:** Add `DescribeTeoFunctionV2ById` method to TeoService if not already present, otherwise create a new service method

**Rationale:**
- Maintains separation of concerns between resource and service layers
- Allows code reuse if other resources need similar functionality
- Follows the established service layer pattern
- Makes testing easier with clear service boundaries

**Alternatives considered:**
- Call API directly from resource: Simpler but violates layered architecture
- Reuse existing DescribeTeoFunctionById: Would work but naming is confusing

## Risks / Trade-offs

### Risk 1: API Changes
**Description:** The teo v20220901 API may change, breaking the implementation

**Mitigation:**
- Use vendor mode for SDK dependency management
- Pin to specific API version (v20220901)
- Monitor API changelog for breaking changes
- Implement proper error handling for API failures

### Risk 2: Asynchronous Operation Timeout
**Description:** Function creation may take longer than the 600-second timeout, causing false failures

**Mitigation:**
- Set reasonable default timeout (600s) that covers most scenarios
- Document timeout in resource documentation
- Users can retry the operation if timeout occurs
- Consider increasing timeout if issues are reported

### Risk 3: Content Size Limitations
**Description:** Function content is limited to 5MB by the API, but Terraform schemas don't enforce size limits

**Mitigation:**
- Document the 5MB limit in resource description
- Rely on API validation to catch oversized content
- API will return clear error if content is too large

### Risk 4: State Inconsistency
**Description:** State refresh may fail or return stale data, leading to inconsistencies

**Mitigation:**
- Implement tccommon.InconsistentCheck to detect state drift
- Use retry logic for read operations
- Log all state inconsistencies for debugging
- Clear state on deletion or not-found errors

### Risk 5: Rate Limiting
**Description:** Frequent operations may hit API rate limits, causing failures

**Mitigation:**
- Resource.Retry with exponential backoff handles rate limits
- tccommon.WriteRetryTimeout and tccommon.ReadRetryTimeout provide appropriate delays
- Users can reduce operation frequency if issues occur

### Trade-off 1: Code Duplication
**Trade-off:** Duplicating code from `resource_tc_teo_function.go` increases maintenance burden

**Rationale:** Code duplication is acceptable because:
- Both resources may evolve independently
- Sharing code would increase complexity and potential for breaking changes
- The code is well-tested and follows a stable pattern
- Duplicated code is minimal and focused

### Trade-off 2: Polling vs Event-based
**Trade-off:** Using polling for state refresh is less efficient than event-based mechanisms

**Rationale:** Polling is acceptable because:
- The API doesn't provide event-based mechanisms
- Polling is simple, reliable, and matches existing patterns
- 10-second polling interval is reasonable for edge function provisioning
- Timeout of 600s prevents indefinite waiting

## Migration Plan

No migration is required because this is a new resource. Existing `tencentcloud_teo_function` resources will continue to work as before.

For users who want to adopt the new resource:
1. Create new `tencentcloud_teo_function_v2` resources alongside existing resources
2. Import existing functions using `terraform import` with composite ID format
3. Gradually migrate state management to the new resource
4. Eventually remove old resources (optional, not required)

No automated migration tooling is provided. Manual migration is straightforward due to identical APIs.

## Open Questions

None identified at this time. The design is straightforward and follows established patterns.
