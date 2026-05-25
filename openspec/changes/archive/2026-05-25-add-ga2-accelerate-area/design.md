## Context

TencentCloud GA2 (Global Accelerator 2.0) provides cloud APIs for managing accelerate areas attached to a global accelerator instance. The SDK package `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ga2/v20250115` is already vendored. No existing Terraform resource covers GA2 accelerate area management.

The four relevant APIs are:
- `CreateAccelerateAreas` - creates accelerate areas (async, returns TaskId)
- `DescribeAccelerateAreas` - queries accelerate areas (paginated with Offset/Limit)
- `ModifyAccelerateAreas` - modifies accelerate areas (async, returns TaskId)
- `DeleteAccelerateAreas` - deletes accelerate areas by IDs (async, returns TaskId)

## Goals / Non-Goals

**Goals:**
- Implement `tencentcloud_ga2_accelerate_area` resource with full CRUD lifecycle
- Handle async operations by polling `DescribeAccelerateAreas` after Create/Modify/Delete
- Support import via `global_accelerator_id`
- Follow existing provider patterns (retry, error handling, timeouts)

**Non-Goals:**
- Data source for listing accelerate areas (separate future work)
- Managing the global accelerator instance itself
- Exposing TaskId as a user-facing attribute (it's only used internally for polling)

## Decisions

### 1. Resource ID uses `global_accelerator_id`

**Rationale**: The accelerate areas are a sub-resource of a global accelerator instance. The `DescribeAccelerateAreas` API takes `GlobalAcceleratorId` as the primary query parameter and returns all areas for that instance. Since the resource manages all accelerate areas for a given accelerator, the resource ID is the `GlobalAcceleratorId` itself.

**Alternative considered**: Using a composite ID with individual `AcceleratorAreaId` values. Rejected because the Create/Modify APIs operate on the full set of areas for an accelerator, not individual areas.

### 2. Async polling after write operations

**Rationale**: All write APIs (Create/Modify/Delete) return a `TaskId` indicating async execution. The resource must call `DescribeAccelerateAreas` in a polling loop after each write operation to confirm the changes have taken effect, using `resource.Retry` with the configured timeout.

### 3. Schema structure for `accelerator_areas`

**Rationale**: The `accelerator_areas` field maps to `[]*AcceleratorAreas` in the SDK. It is a TypeList of objects with fields: `accelerate_region` (Required), `bandwidth` (Required), `isp_type` (Optional, default BGP), `ip_version` (Optional, default IPv4), `accelerator_area_id` (Computed), `ip_address` (Computed), `ip_address_info_set` (Computed).

For the Modify API, `AcceleratorAreaId` is needed to identify which area to modify, so it must be included in the input. The schema marks it as Optional+Computed.

### 4. Delete requires collecting AcceleratorAreaIds from Read

**Rationale**: The `DeleteAccelerateAreas` API requires `AcceleratorAreaIds`. The Delete function must first call `DescribeAccelerateAreas` to get all current area IDs, then pass them to the delete API.

### 5. File organization

- `tencentcloud/services/ga2/resource_tc_ga2_accelerate_area.go` - resource definition and CRUD
- `tencentcloud/services/ga2/resource_tc_ga2_accelerate_area_test.go` - unit tests with gomonkey mocks
- `tencentcloud/services/ga2/resource_tc_ga2_accelerate_area.md` - example usage documentation
- `tencentcloud/services/ga2/service_tencentcloud_ga2.go` - service layer (if ga2 service file doesn't exist yet)

### 6. Timeouts

**Rationale**: Since all write operations are async, the schema must declare a `Timeouts` block with Create, Update, and Delete timeouts. Default timeout: 10 minutes for each.

## Risks / Trade-offs

- [Risk] Async polling may time out if the backend is slow → Mitigation: Use configurable Timeouts in schema, default 10 minutes, with clear error messages on timeout.
- [Risk] Delete requires reading all area IDs first, adding an extra API call → Mitigation: This is unavoidable given the API design; the extra call is lightweight.
- [Trade-off] Managing all areas as a single resource (vs. one resource per area) means adding/removing a single area requires updating the entire set → This matches the API design where Create/Modify operate on the full set.
