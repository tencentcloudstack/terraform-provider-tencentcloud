## Context

TencentCloud GA2 (Global Accelerator v2) provides APIs for managing accelerate areas under a global accelerator instance. The SDK package `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ga2/v20250115` is already vendored and provides four APIs: `CreateAccelerateAreas`, `DescribeAccelerateAreas`, `ModifyAccelerateAreas`, and `DeleteAccelerateAreas`.

All write operations (Create/Modify/Delete) are asynchronous—they return a `TaskId` and the actual changes take effect after some delay. The resource must poll `DescribeAccelerateAreas` after each write to confirm the operation has completed.

There is no existing ga2 service directory in `tencentcloud/services/`, so all files will be new.

## Goals / Non-Goals

**Goals:**
- Implement a fully functional `tencentcloud_ga2_accelerate_area` resource with CRUD operations
- Handle async operations by polling DescribeAccelerateAreas after Create/Modify/Delete
- Support import via `global_accelerator_id`
- Follow existing provider patterns (retry, error handling, timeouts)
- Provide unit tests using gomonkey mocks

**Non-Goals:**
- Data source for listing accelerate areas (separate resource if needed later)
- Managing the global accelerator instance itself
- Handling task status polling via a separate DescribeTask API (not available; use DescribeAccelerateAreas to verify state)

## Decisions

### 1. Resource ID: Use `global_accelerator_id` as the Terraform resource ID

**Rationale**: Accelerate areas are a sub-resource of a global accelerator instance. The `DescribeAccelerateAreas` API takes `GlobalAcceleratorId` as input and returns all accelerate areas for that instance. Since this resource manages all accelerate areas for a given accelerator, the accelerator ID is the natural resource identifier.

**Alternative considered**: Using a composite ID with `AcceleratorAreaId`. Rejected because the API manages areas as a batch (Create/Modify take a list of areas, not individual ones), making the accelerator ID the correct granularity.

### 2. Schema Design: `accelerator_areas` as input, `accelerate_area_set` as computed output

**Rationale**: The Create and Modify APIs accept `AcceleratorAreas` (a list of area configurations). The Read API returns `AccelerateAreaSet` which includes additional computed fields like `AcceleratorAreaId`, `IpAddress`, and `IpAddressInfoSet`. Separating input from output avoids schema conflicts and makes the resource behavior clear.

- `global_accelerator_id` (Required, ForceNew, String): The GA2 instance ID
- `accelerator_areas` (Required, List): Input configuration for areas
  - `accelerate_region` (Required, String)
  - `bandwidth` (Required, Int)
  - `isp_type` (Optional, String, default "BGP")
  - `ip_version` (Optional, String, default "IPv4")
- `accelerate_area_set` (Computed, List): Read-only output from DescribeAccelerateAreas
  - `accelerator_area_id` (String)
  - `accelerate_region` (String)
  - `bandwidth` (Int)
  - `isp_type` (String)
  - `ip_version` (String)
  - `ip_address` (List of String)
  - `ip_address_info_set` (List)
    - `ip_address` (String)
    - `isp_type` (String)

### 3. Async Handling: Poll DescribeAccelerateAreas after write operations

**Rationale**: Since Create/Modify/Delete are async, we need to verify the operation completed. The approach is to call DescribeAccelerateAreas in a retry loop after each write operation until the expected state is observed. Use `resource.Retry` with the configured timeout.

### 4. Delete: Collect all AcceleratorAreaIds from Read, then call DeleteAccelerateAreas

**Rationale**: The Delete API requires `AcceleratorAreaIds`. We first call DescribeAccelerateAreas to get all current area IDs, then pass them to DeleteAccelerateAreas.

### 5. Service Layer: Create `service_tencentcloud_ga2.go`

**Rationale**: Following the provider pattern, API calls are encapsulated in a service layer file. This provides reusable client initialization and API wrappers.

## Risks / Trade-offs

- [Risk] Async operations may take longer than expected → Use configurable Timeouts in schema (default 10 minutes for create/update/delete)
- [Risk] DescribeAccelerateAreas pagination → Implement full pagination in Read to handle cases with many accelerate areas
- [Trade-off] Managing all areas as a single resource vs individual areas → Chose batch management to match API design; users who need individual area management can use lifecycle blocks
