## Context

TencentCloud TEO (EdgeOne) provides the `BindSharedCNAMEWithContext` API to bind/unbind acceleration domains to shared CNAMEs. The API uses a `BindType` field ("bind"/"unbind") to control the operation direction. The `DescribeSharedCNAMEWithContext` API returns shared CNAME info including bound acceleration domains, which can be used to read the current binding state.

This resource follows the RESOURCE_KIND_ATTACHMENT pattern: Create calls the bind API, Delete calls the unbind API, Read reads the current domain list from the describe API, and Update computes the diff and calls bind/unbind for changed domains.

## Goals / Non-Goals

**Goals:**
- Provide a Terraform resource `tencentcloud_teo_domain_shared_cname_attachment` to manage domain-to-shared-CNAME bindings
- Support Create (bind), Read (read domain list), Update (diff-based bind/unbind), and Delete (unbind all) operations
- Use composite ID (`zone_id` + `shared_cname`) with `tccommon.FILED_SP` separator
- Follow existing provider patterns (WithContext APIs, retry, error handling, logging)

**Non-Goals:**
- Managing the shared CNAME resource itself (creation/deletion of shared CNAMEs)
- Exposing `BindType` to users (it is derived from the operation)

## Decisions

### 1. Resource Schema Design

The resource exposes top-level fields (not nested under `bind_shared_cname_maps`):
- `zone_id` (Required, ForceNew, String): The zone ID that the acceleration domain belongs to
- `shared_cname` (Required, ForceNew, String): The shared CNAME to bind to
- `domain_names` (Required, List of String): The acceleration domain names to bind (supports in-place update)

**Rationale**: A single resource instance represents one `shared_cname` binding. Flattening the schema makes it simpler and more idiomatic. `zone_id` and `shared_cname` are ForceNew because they form the resource identity. `domain_names` is mutable and supports in-place update.

### 2. Composite ID Format

ID format: `{zone_id}#{shared_cname}`

Using `tccommon.FILED_SP` (`#`) as separator between zone_id and shared_cname.

**Rationale**: The binding is uniquely identified by the combination of zone and shared CNAME. Domain names are not part of the ID because they can change via Update. This simplifies import.

### 3. Read Implementation

The Read function calls `DescribeSharedCNAMEWithContext` with a filter on the specific shared CNAME, then populates `domain_names` from the `AccelerationDomains` field of the response.

**Rationale**: There is no dedicated "describe binding" API. The `DescribeSharedCNAMEWithContext` API returns `SharedCNAMEInfo` which includes `AccelerationDomains` - we read the current bound domain list from there.

### 4. Update Implementation

Update uses `d.GetChange("domain_names")` to compute the diff between old and new domain lists, then:
1. Calls `BindSharedCNAMEWithContext` with `BindType = "unbind"` for domains removed from the list
2. Calls `BindSharedCNAMEWithContext` with `BindType = "bind"` for domains added to the list

**Rationale**: The same API handles both bind and unbind via the `BindType` parameter. Diff-based update avoids unnecessary API calls for unchanged domains.

### 5. Delete Implementation

Delete reads `domain_names` from state and calls `BindSharedCNAMEWithContext` with `BindType = "unbind"` for all of them.

**Rationale**: The same API handles both bind and unbind via the `BindType` parameter, which is the standard TEO pattern for attachment resources.

### 6. WithContext API Usage

All API calls use the `WithContext` variant (e.g., `BindSharedCNAMEWithContext`, `DescribeSharedCNAMEWithContext`) with a context created via `tccommon.NewResourceLifeCycleHandleFuncContext`.

**Rationale**: Consistent with the existing provider pattern used in other TEO resources.

## Risks / Trade-offs

- [Risk] The `DescribeSharedCNAMEWithContext` API may return a limited number of acceleration domains per shared CNAME. → Mitigation: Set `Limit = 200` in the request.
- [Risk] The `BindSharedCNAMEWithContext` API is not marked as async, but eventual consistency may cause Read to not immediately reflect the binding. → Mitigation: Use retry with `tccommon.ReadRetryTimeout` in the Read function.
