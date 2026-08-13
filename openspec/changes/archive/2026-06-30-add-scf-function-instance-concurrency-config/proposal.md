## Why

The `tencentcloud_scf_function` resource currently lacks the `instance_concurrency_config` parameter, which is essential for configuring multi-concurrency behavior for Web functions in SCF. This parameter enables users to control single-instance concurrency settings (static/dynamic), max concurrency limits, instance isolation, and session configuration — capabilities that are already supported by the TencentCloud SCF API. Without this parameter, users cannot leverage these concurrency features through Terraform, requiring manual API calls or console operations.

## What Changes

- Add new optional parameter `instance_concurrency_config` (TypeList) to `tencentcloud_scf_function` resource, containing the following sub-fields:
  - `dynamic_enabled` (TypeString, Optional) - Whether to enable intelligent dynamic concurrency
  - `max_concurrency` (TypeInt, Optional) - Maximum single-instance concurrency (range: 1-100)
  - `instance_isolation_enabled` (TypeString, Optional) - Security isolation switch
  - `type` (TypeString, Optional) - Session-Based or Request-Based concurrency mode
  - `mix_node_config` (TypeList, Optional) - Dynamic concurrency parameters
  - `session_config` (TypeList, Optional) - Session configuration parameters
- Wire the new parameter through:
  - `CreateFunction` API call (as input)
  - `UpdateFunctionConfiguration` API call (as input)
  - `GetFunction` API call (as read output, to sync state)

## Capabilities

### New Capabilities

- `scf-function-instance-concurrency-config`: Add `instance_concurrency_config` parameter to `tencentcloud_scf_function` resource for configuring single-instance multi-concurrency behavior for Web functions.

### Modified Capabilities

<!-- No existing capabilities are modified. This is a pure addition of an Optional parameter. -->

## Impact

- **Affected Files**:
  - `tencentcloud/services/scf/resource_tc_scf_function.go` - Schema definition, Create/Read/Update logic
  - `tencentcloud/services/scf/resource_tc_scf_function_test.go` - Unit test cases
  - `tencentcloud/services/scf/resource_tc_scf_function.md` - Documentation
- **Cloud API Dependencies**: Uses existing `GetFunction`, `CreateFunction`, `UpdateFunctionConfiguration` APIs (no new API calls needed)
- **Backward Compatibility**: Fully backward compatible — the new parameter is Optional with no default value; existing configurations will continue to work unchanged
- **SDK**: Uses existing `InstanceConcurrencyConfig` type from `tencentcloud-sdk-go/scf/v20180416`
