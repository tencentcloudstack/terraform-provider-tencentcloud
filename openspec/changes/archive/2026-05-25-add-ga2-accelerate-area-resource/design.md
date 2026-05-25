## Context

TencentCloud GA2 (Global Accelerator v2) provides APIs for managing accelerate areas under a global accelerator instance. The GA2 SDK (`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ga2/v20250115`) is already vendored. All CRUD operations for accelerate areas are asynchronous—they return a `TaskId` and require polling `DescribeTaskResult` until the task completes.

Currently there is no `tencentcloud/services/ga2/` directory, so this is a new service module.

## Goals / Non-Goals

**Goals:**
- Implement `tencentcloud_ga2_accelerate_area` resource with full CRUD lifecycle
- Handle async operations by polling `DescribeTaskResult` after Create/Modify/Delete
- Support Terraform import via `global_accelerator_id`
- Provide unit tests using gomonkey mocks (not acceptance tests)
- Register the resource in provider.go and provider.md

**Non-Goals:**
- Data source for accelerate areas (separate change)
- Managing the global accelerator instance itself (separate resource)
- Acceptance tests requiring real cloud credentials

## Decisions

### 1. Resource ID Strategy
**Decision**: Use `global_accelerator_id` as the Terraform resource ID.
**Rationale**: Accelerate areas are a sub-resource of a global accelerator instance. The DescribeAccelerateAreas API takes only `GlobalAcceleratorId` to list all areas. There is no single "accelerate area ID" that uniquely identifies the collection—the collection is identified by its parent instance.

### 2. Schema Design
**Decision**: Use `global_accelerator_id` (Required, ForceNew) as the primary identifier, `accelerator_areas` (Required, List) for Create/Update input, and `accelerate_area_set` (Computed, List) for Read output.
**Rationale**: The Create/Modify APIs accept `AcceleratorAreas` with fields (AccelerateRegion, Bandwidth, IspType, IpVersion). The Read API returns `AccelerateAreaSet` with additional computed fields (AcceleratorAreaId, IpAddress, IpAddressInfoSet). Separating input and output schemas avoids drift issues.

### 3. Async Task Polling
**Decision**: After each async operation (Create/Modify/Delete), poll `DescribeTaskResult` with retry until `Status` indicates completion.
**Rationale**: All three mutation APIs return a `TaskId`. The operation is not effective until the task completes. Polling ensures Terraform state is consistent.

### 4. Delete Implementation
**Decision**: In the Delete function, first call DescribeAccelerateAreas to get all `AcceleratorAreaId` values, then pass them to `DeleteAccelerateAreas`.
**Rationale**: The Delete API requires `AcceleratorAreaIds` which are only available from the Read response. This ensures all areas under the instance are deleted.

### 5. Service Layer
**Decision**: Create `service_tencentcloud_ga2.go` with helper functions for API calls and task polling.
**Rationale**: Follows the project pattern of separating service-layer logic from resource CRUD functions.

### 6. Pagination in Read
**Decision**: Implement pagination in DescribeAccelerateAreas using Offset/Limit to fetch all results.
**Rationale**: The API supports pagination. We must fetch all pages to get the complete list of accelerate areas.

## Risks / Trade-offs

- [Risk] Task polling timeout → Mitigation: Use Terraform Timeouts block (Create/Update/Delete) with reasonable defaults (e.g., 10 minutes). Polling loop respects context deadline.
- [Risk] Eventual consistency after task completion → Mitigation: After task completes, perform a Read to verify state before returning.
- [Trade-off] Separating input (`accelerator_areas`) and output (`accelerate_area_set`) schemas adds complexity but avoids plan diffs on computed fields like IpAddress.
