## Why

The `tencentcloud_scf_function` resource currently lacks the `qualifier` parameter, which is required to specify a function version or alias when performing operations like reading function details, deleting functions, and managing triggers. Without this parameter, users cannot target specific function versions/aliases, limiting the usability of the resource in versioned deployment scenarios.

## What Changes

- Add a new **Optional** + **Computed** `qualifier` parameter (TypeString) to the `tencentcloud_scf_function` resource schema
- Pass `qualifier` to the `GetFunction` API request when reading function state
- Read `qualifier` from the `GetFunction` API response (both top-level `response.Qualifier` and `response.Triggers[].Qualifier`)
- Pass `qualifier` to the `DeleteFunction` API request when deleting the function
- Pass `qualifier` to the `CreateTrigger` and `DeleteTrigger` API requests when managing triggers

## Capabilities

### New Capabilities
- `scf-function-qualifier-param`: Support specifying a function version or alias (`qualifier`) for the `tencentcloud_scf_function` resource, enabling version-aware function operations (read, delete, trigger management).

### Modified Capabilities
<!-- None - this is a new optional parameter addition, no existing capability requirements change -->

## Impact

- **Affected Code**:
  - `tencentcloud/services/scf/resource_tc_scf_function.go` — schema definition and CRUD methods
  - `tencentcloud/services/scf/service_tencentcloud_scf.go` — service layer methods (DescribeFunction, DeleteFunction, CreateTriggers, DeleteTriggers)
  - `tencentcloud/services/scf/resource_tc_scf_function.md` — documentation
  - `tencentcloud/services/scf/resource_tc_scf_function_test.go` — test cases
- **API Dependencies**: `GetFunction`, `DeleteFunction`, `CreateTrigger`, `DeleteTrigger` (all from `scf/v20180416`)
- **Backward Compatibility**: Fully backward compatible — the new parameter is Optional + Computed, existing configurations continue to work without changes
- **Note**: `CreateFunction`, `UpdateFunctionCode`, and `UpdateFunctionConfiguration` APIs do NOT support the `Qualifier` parameter in the current SDK, so the qualifier is not passed during create or update operations