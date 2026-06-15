## Context

CynosDB (TDSQL-C) is a cloud-native database service. Clusters can have read-only instance groups, but access to these groups must be explicitly opened via the `OpenClusterReadOnlyInstanceGroupAccess` API. This API is asynchronous—it returns a `FlowId` that must be polled via `DescribeFlow` until the operation completes.

The existing codebase already has:
- A `CynosdbService.DescribeFlow()` helper in `tencentcloud/services/cynosdb/service_tencentcloud_cynosdb.go` that polls flow status.
- Multiple cynosdb resources that use this pattern (e.g., `cynosdb_cluster_password_complexity`, `cynosdb_proxy_end_point`).
- The `cynosdb/v20190107` SDK package is already vendored.

## Goals / Non-Goals

**Goals:**
- Provide a Terraform `RESOURCE_KIND_OPERATION` resource that opens read-only instance group access for a CynosDB cluster.
- Support async flow polling via `DescribeFlow` to wait for the operation to complete.
- Follow existing provider patterns for operation resources (Create-only, no Read/Update/Delete logic).
- Include unit tests using gomonkey mocks.

**Non-Goals:**
- No Read/Update/Delete lifecycle management (this is a one-shot operation).
- No import support (operation resources are not importable).
- No acceptance tests requiring real cloud credentials.

## Decisions

### 1. Resource structure: Create-only with empty Read/Delete

**Decision**: Implement Create with the API call + flow polling. Read returns nil (no state to read). Delete returns nil (nothing to destroy).

**Rationale**: This follows the established `RESOURCE_KIND_OPERATION` pattern used by `tencentcloud_cls_open_service_operation` and similar resources. The resource ID will be set to a generated token via `helper.BuildToken()` since there's no meaningful ID to track.

**Alternative considered**: Using `cluster_id` as the resource ID. Rejected because the operation is not idempotent and doesn't represent a persistent resource.

### 2. Async flow polling via existing CynosdbService.DescribeFlow

**Decision**: After calling `OpenClusterReadOnlyInstanceGroupAccess`, use `resource.Retry` with `tccommon.WriteRetryTimeout` to poll `CynosdbService.DescribeFlow(flowId)` until it returns success.

**Rationale**: This reuses the existing, well-tested flow polling infrastructure. The `DescribeFlow` helper already handles status codes (0 = success, 2 = failed).

### 3. Schema with Timeouts block

**Decision**: Include a `Timeouts` block with a `Create` timeout to support configurable wait times for the async operation.

**Rationale**: Required by provider conventions for async operations. Users may need longer timeouts for large clusters.

### 4. Unit tests with gomonkey mocks

**Decision**: Use gomonkey to mock the SDK client methods (`OpenClusterReadOnlyInstanceGroupAccessWithContext`, `DescribeFlow`) and test the Create function logic.

**Rationale**: Per project requirements, new terraform resources use gomonkey-based unit tests rather than acceptance tests.

## Risks / Trade-offs

- **[Risk] Flow polling timeout**: The async operation may take longer than the default timeout. → **Mitigation**: Expose configurable `Timeouts.Create` in the schema.
- **[Risk] API returns nil FlowId**: The API might succeed but return a nil FlowId. → **Mitigation**: Check for nil FlowId after the API call and return `NonRetryableError` if nil.
- **[Trade-off] No state tracking**: Since this is an operation resource, there's no way to detect if the operation was already performed. Re-applying will attempt to open access again. → **Accepted**: This is inherent to the RESOURCE_KIND_OPERATION pattern.
