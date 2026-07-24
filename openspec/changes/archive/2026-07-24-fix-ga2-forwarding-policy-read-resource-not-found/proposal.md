## Why

The GA2 resources' Read functions do not handle `ResourceNotFound` errors returned by their describe APIs. When a resource is deleted outside of Terraform (e.g., via console or API), the Read function returns an error instead of gracefully removing the resource from state. This causes Terraform to report errors during plan/apply instead of detecting the resource was deleted and proposing to recreate it.

## What Changes

- Add a unified common helper function `HandleGa2ReadNotFound` in `resource_tc_ga2_common.go` that handles both cases:
  1. SDK `ResourceNotFound` error — delegates to the existing `HandleGa2ResourceNotFoundError` helper (which checks `!d.IsNewResource()`).
  2. Nil/empty response from the service layer — when the resource is not new, logs a warning, clears the resource ID, and returns handled=true; when `d.IsNewResource()` is true, returns an error so the Create → Read cycle can propagate the failure.
- Apply this helper to all 6 GA2 resource Read methods:
  - `resource_tc_ga2_forwarding_policy.go` — `DescribeForwardingPolicy`
  - `resource_tc_ga2_global_accelerator.go` — `DescribeGlobalAccelerators`
  - `resource_tc_ga2_listener.go` — `DescribeListeners`
  - `resource_tc_ga2_endpoint_group.go` — `DescribeEndpointGroups`
  - `resource_tc_ga2_accelerate_area.go` — `DescribeAccelerateAreas`
  - `resource_tc_ga2_forwarding_rule.go` — `DescribeForwardingRule`

## Capabilities

### New Capabilities
<!-- No new capabilities - this is a bug fix -->

### Modified Capabilities
<!-- No spec-level behavior changes - this is a defensive error handling improvement -->

## Impact

- **Affected code**:
  - `tencentcloud/services/ga2/resource_tc_ga2_common.go` — unified helper `HandleGa2ReadNotFound` and legacy `HandleGa2ResourceNotFoundError`
  - `tencentcloud/services/ga2/resource_tc_ga2_forwarding_policy.go` — Read function uses `HandleGa2ReadNotFound`
  - `tencentcloud/services/ga2/resource_tc_ga2_global_accelerator.go` — Read function uses `HandleGa2ReadNotFound`
  - `tencentcloud/services/ga2/resource_tc_ga2_listener.go` — Read function uses `HandleGa2ReadNotFound`
  - `tencentcloud/services/ga2/resource_tc_ga2_endpoint_group.go` — Read function uses `HandleGa2ReadNotFound`
  - `tencentcloud/services/ga2/resource_tc_ga2_accelerate_area.go` — Read function uses `HandleGa2ReadNotFound`
  - `tencentcloud/services/ga2/resource_tc_ga2_forwarding_rule.go` — Read function uses `HandleGa2ReadNotFound`
- **API**: Uses existing describe APIs (no API changes)
- **Dependencies**: None (uses already-imported `sdkErrors` package)
- **Backward compatibility**: Fully compatible — this change improves error handling without changing schema or normal behavior