## Why

The `tencentcloud_ga2_forwarding_policy` resource's Read function does not handle `ResourceNotFound` errors returned by the `DescribeForwardingPolicy` API. When the forwarding policy is deleted outside of Terraform (e.g., via console or API), the Read function returns an error instead of gracefully removing the resource from state. This causes Terraform to report errors during plan/apply instead of detecting the resource was deleted and proposing to recreate it.

## What Changes

- Modify `resourceTencentCloudGa2ForwardingPolicyRead` to check if the error returned by `DescribeGa2ForwardingPolicyById` is a `TencentCloudSDKError` with code `"ResourceNotFound"`. When this occurs, log a warning, clear the resource ID, and return nil.

## Capabilities

### New Capabilities
<!-- No new capabilities - this is a bug fix -->

### Modified Capabilities
<!-- No spec-level behavior changes - this is a defensive error handling improvement -->

## Impact

- **Affected code**: `tencentcloud/services/ga2/resource_tc_ga2_forwarding_policy.go` — the `resourceTencentCloudGa2ForwardingPolicyRead` function (lines 153-156)
- **API**: Uses existing `DescribeForwardingPolicy` API (no API changes)
- **Dependencies**: None (uses already-imported `sdkErrors` package)
- **Backward compatibility**: Fully compatible — this change improves error handling without changing schema or normal behavior