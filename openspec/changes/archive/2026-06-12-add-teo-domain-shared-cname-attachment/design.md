## Context

TencentCloud TEO (EdgeOne) provides the `BindSharedCNAME` API to bind/unbind acceleration domains to shared CNAMEs. The API uses a `BindType` field ("bind"/"unbind") to control the operation direction. The `DescribeSharedCNAME` API returns shared CNAME info including bound acceleration domains, which can be used to verify binding status.

This is a standard RESOURCE_KIND_ATTACHMENT pattern: Create calls the bind API, Delete calls the unbind API, and Read verifies the binding exists via the describe API.

## Goals / Non-Goals

**Goals:**
- Provide a Terraform resource `tencentcloud_teo_domain_shared_cname_attachment` to manage domain-to-shared-CNAME bindings
- Support Create (bind), Read (verify binding), and Delete (unbind) operations
- Use composite ID (`zone_id` + `shared_cname` + `domain_names`) with `tccommon.FILED_SP` separator
- Follow existing provider patterns (retry, error handling, logging)

**Non-Goals:**
- Managing the shared CNAME resource itself (creation/deletion of shared CNAMEs)
- Supporting Update operations (attachment resources are immutable - destroy and recreate)
- Exposing `BindType` to users (it is derived from the operation: create=bind, delete=unbind)

## Decisions

### 1. Resource Schema Design

The resource exposes:
- `zone_id` (Required, ForceNew, String): The zone ID that the acceleration domain belongs to
- `bind_shared_cname_maps` (Required, ForceNew, List): The binding relationships between domains and shared CNAMEs
  - `shared_cname` (Required, String): The shared CNAME to bind to
  - `domain_names` (Required, List of String): The acceleration domain names to bind

**Rationale**: Since this is a CRD-only resource (no Update API), all fields are ForceNew. The `BindType` is not exposed because it's implicitly determined by the lifecycle operation.

### 2. Composite ID Format

ID format: `{zone_id}#{shared_cname}#{domain_name1,domain_name2,...}`

Using `tccommon.FILED_SP` (`#`) as separator between zone_id, shared_cname, and a comma-joined list of domain_names.

**Rationale**: The binding is uniquely identified by the combination of zone, shared CNAME, and the set of domains. This allows proper import support.

### 3. Read Implementation

The Read function calls `DescribeSharedCNAME` with a filter on the specific shared CNAME, then checks if the expected domains are present in the `AccelerationDomains` field of the response.

**Rationale**: There is no dedicated "describe binding" API. The `DescribeSharedCNAME` API returns `SharedCNAMEInfo` which includes `AccelerationDomains` - we can verify the binding exists by checking if our domains appear in that list.

### 4. Delete Implementation

Delete calls `BindSharedCNAME` with `BindType = "unbind"` using the same `BindSharedCNAMEMaps` structure.

**Rationale**: The same API handles both bind and unbind via the `BindType` parameter, which is the standard TEO pattern for attachment resources.

## Risks / Trade-offs

- [Risk] The `DescribeSharedCNAME` API returns at most 100 acceleration domains per shared CNAME. If a shared CNAME has more than 100 bound domains, the Read function may not find all expected domains. → Mitigation: This is an API limitation; document it. In practice, most shared CNAMEs have fewer than 100 bound domains.
- [Risk] The `BindSharedCNAME` API is not marked as async, but eventual consistency may cause Read to not immediately reflect the binding. → Mitigation: Use retry with `tccommon.ReadRetryTimeout` in the Read function.
