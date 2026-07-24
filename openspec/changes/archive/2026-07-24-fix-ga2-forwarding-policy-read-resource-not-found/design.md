## Context

The GA2 resources use various describe APIs (e.g., `DescribeForwardingPolicy`, `DescribeGlobalAccelerators`, `DescribeListeners`, `DescribeEndpointGroups`, `DescribeAccelerateAreas`, `DescribeForwardingRule`) to read the current state of their resources. The Read functions currently call the service layer and, if an error is returned, propagate it directly to Terraform without checking whether the error means the resource was deleted externally.

The GA2 SDK defines `ResourceNotFound` as a known error code (see `vendor/.../ga2/v20250115/errors.go`). This change brings all GA2 resources in line with the established pattern of handling `ResourceNotFound` gracefully.

## Goals / Non-Goals

**Goals:**
- When any GA2 describe API returns a `ResourceNotFound` SDK error, the Read function should log a warning, clear the resource ID, and return nil (instead of returning an error)
- When any GA2 describe API returns an empty/nil response (resource not found by ID), the Read function should log a warning, clear the resource ID, and return nil
- When the resource is new (`d.IsNewResource()` is true, i.e., during the initial Create → Read cycle), both of the above cases should propagate the error normally
- The implementation should be abstracted into a unified common helper function `HandleGa2ReadNotFound` to avoid code duplication
- The change must be minimal and not affect any other behavior

**Non-Goals:**
- No schema changes
- No changes to Create, Update, or Delete functions
- No changes to the service layer

## Decisions

**Decision 1: Create a unified helper function `HandleGa2ReadNotFound`**

Instead of duplicating the not-found check pattern in every resource, a single exported function in `resource_tc_ga2_common.go` encapsulates the logic. The function handles both cases:
1. SDK `ResourceNotFound` error — delegates to the existing `HandleGa2ResourceNotFoundError` helper.
2. Nil/empty response from the service layer — when the resource is not new, logs a warning and clears the ID; when `d.IsNewResource()` is true, returns an error so the Create → Read cycle can propagate the failure.

The function returns `(handled bool, err error)` — a clean two-value contract for callers.

**Decision 2: Use `!d.IsNewResource()` guard for both SDK error and nil response**

Following the reference pattern from other resources, the not-found check includes `!d.IsNewResource()` to prevent incorrectly clearing the ID during the initial resource creation cycle (when the Read is called from within Create). This applies uniformly to both the SDK `ResourceNotFound` error case and the nil response case.

**Decision 3: Apply to all GA2 resources uniformly**

The same pattern is applied to all 6 GA2 resources:
- `tencentcloud_ga2_forwarding_policy` (DescribeForwardingPolicy)
- `tencentcloud_ga2_global_accelerator` (DescribeGlobalAccelerators)
- `tencentcloud_ga2_listener` (DescribeListeners)
- `tencentcloud_ga2_endpoint_group` (DescribeEndpointGroups)
- `tencentcloud_ga2_accelerate_area` (DescribeAccelerateAreas)
- `tencentcloud_ga2_forwarding_rule` (DescribeForwardingRule)

## Risks / Trade-offs

- **Risk**: If the API returns `ResourceNotFound` for a transient reason (e.g., API inconsistency) during a normal Read, the resource will be removed from state. → **Mitigation**: This is the same behavior as all other Terraform resources. The retry mechanism in the service layer already handles transient errors within the `ReadRetryTimeout` window. If the API returns `ResourceNotFound` consistently, the resource truly does not exist.