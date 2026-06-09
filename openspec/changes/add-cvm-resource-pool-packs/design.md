## Context

The CVM service currently lacks Terraform resource support for resource pool packs management. Resource pool packs are a CVM feature that allows users to purchase and manage compute resource pools in advance. 

The implementation follows the standard Terraform Provider pattern for TencentCloud:
- Service layer in `service_tencentcloud_cvm.go` for API calls
- Resource layer in `resource_tc_cvm_resource_pool_packs.go` for Terraform lifecycle
- Uses `tencentcloud-sdk-go` CVM client for API interactions

**Constraints:**
- No Update API available - resource is create/read/delete only (ForceNew on all fields)
- Must follow existing CVM service patterns for consistency
- Must maintain backward compatibility (new resource, no breaking changes)

## Goals / Non-Goals

**Goals:**
- Provide Terraform resource for CVM resource pool packs lifecycle management
- Support standard Terraform operations: create, read, delete
- Include retry logic for eventual consistency
- Provide comprehensive test coverage
- Generate proper documentation

**Non-Goals:**
- Update functionality (no Update API available from CVM)
- Data source for querying existing packs (can be added later if needed)
- Management of individual instances within the pack (separate concern)
- Cross-region resource pool management in a single resource

## Decisions

### 1. Resource Naming: `tencentcloud_cvm_resource_pool_packs`

**Decision:** Use plural form `resource_pool_packs` instead of singular `resource_pool_pack`.

**Rationale:** 
- The API name is `PurchaseResourcePoolPacks` (plural)
- Aligns with the API semantics where you can purchase multiple packs in one call
- Consistent with existing resources like `tencentcloud_clb_instances`

**Alternative considered:** Singular form would be more conventional for Terraform resources, but API alignment takes precedence.

### 2. Resource ID Strategy

**Decision:** Use the pack ID returned from `PurchaseResourcePoolPacks` API as the Terraform resource ID.

**Rationale:**
- Simple, single-value identifier
- No composite ID needed (unlike resources with multiple parent/child relationships)
- Directly maps to API query parameter

**Alternative considered:** Composite ID with region - rejected because region is typically handled at provider level.

### 3. ForceNew on All Fields

**Decision:** Mark all resource fields as `ForceNew: true` in the schema.

**Rationale:**
- No Update API available from CVM service
- Any change to pack configuration requires destroy + recreate
- Makes behavior explicit to users

**Alternative considered:** Blocking updates with an error - rejected because ForceNew is the idiomatic Terraform approach.

### 4. Service Layer Methods

**Decision:** Add three service layer methods in `service_tencentcloud_cvm.go`:
- `CreateCvmResourcePoolPacks()` - wraps `PurchaseResourcePoolPacks`
- `DescribeCvmResourcePoolPackById()` - wraps `DescribeResourcePoolPacks` with filtering
- `DeleteCvmResourcePoolPacks()` - wraps `TerminateResourcePoolPacks`

**Rationale:**
- Separation of concerns: service layer handles API interactions, resource layer handles Terraform logic
- Enables retry logic at service layer
- Consistent with existing CVM service patterns

**Alternative considered:** Direct API calls in resource - rejected for maintainability and consistency.

### 5. Retry Strategy

**Decision:** Use `resource.Retry` with `tccommon.ReadRetryTimeout` for query operations.

**Rationale:**
- Handles eventual consistency issues common in cloud APIs
- Standard pattern used throughout the provider
- Prevents spurious failures during resource creation/deletion verification

### 6. Error Handling

**Decision:** Use `defer tccommon.LogElapsed()` and `defer tccommon.InconsistentCheck()` pattern.

**Rationale:**
- Standard error handling pattern in the provider
- Provides consistent logging and debugging experience
- Handles inconsistent state detection

## Risks / Trade-offs

### Risk: No Update Support
**Impact:** Users cannot modify pack configuration without destroy+recreate.  
**Mitigation:** 
- Clearly document in resource `.md` that all fields are ForceNew
- Mark all schema fields with `ForceNew: true` to make behavior explicit
- This is a CVM API limitation, not a provider limitation

### Risk: API Quota Limits
**Impact:** Rapid create/delete cycles during testing could hit API rate limits.  
**Mitigation:** 
- Use retry logic with exponential backoff
- Document rate limit considerations in resource documentation
- Use `ratelimit.Check()` before each API call

### Risk: Resource Deletion Failures
**Impact:** Pack termination might fail if resources are still in use.  
**Mitigation:** 
- Implement proper error messaging from `TerminateResourcePoolPacks` API
- Document prerequisites for deletion (e.g., no active instances)
- Return clear error messages to guide users

### Trade-off: No Data Source
**Decision:** Not implementing a data source in initial version.  
**Rationale:** 
- Keeps scope focused on resource management
- Data source can be added incrementally if user demand exists
- Query API is available and can be easily wrapped later

## Migration Plan

**Deployment:**
1. Merge code changes to main branch
2. No migration needed - this is a new resource with no state impact
3. Update provider documentation website

**Rollback:**
- Simply remove the resource registration from `provider.go`
- No state migration needed
- No impact on existing resources

**Testing:**
- Run acceptance tests with `TF_ACC=1`
- Verify create/read/delete lifecycle
- Test ForceNew behavior on field changes
- Verify error handling for invalid inputs

## Open Questions

None - all implementation details are clear based on the CVM API documentation.
