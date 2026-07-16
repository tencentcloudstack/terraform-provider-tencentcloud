## Why

TEO (EdgeOne) users need the ability to increase plan quotas (e.g., site count, precise access control rules, rate limiting rules) through Terraform. Currently, the terraform provider lacks this capability, forcing users to manually use the console or API to purchase additional quotas. This change adds a Terraform OPERATION resource for the `IncreasePlanQuota` API, enabling Infrastructure-as-Code management of TEO plan quota upgrades.

## What Changes

- Add a new RESOURCE_KIND_OPERATION resource `tencentcloud_teo_increase_plan_quota` that calls the `IncreasePlanQuota` API
- The resource is a one-time operation: it creates (calls the API), and read/delete are no-ops as expected for OPERATION resources
- Input parameters: `plan_id`, `quota_type`, `quota_number` (all required and ForceNew)
- Output parameter: `deal_name` (computed) - the order number returned by the API
- Register the new resource in `tencentcloud/provider.go` and `tencentcloud/provider.md`

## Capabilities

### New Capabilities
- `teo-increase-plan-quota`: Terraform resource for the TEO `IncreasePlanQuota` API, enabling users to increase plan quotas for TEO EdgeOne plans via Infrastructure-as-Code

### Modified Capabilities
<!-- None -->

## Impact

- **New file**: `tencentcloud/services/teo/resource_tc_teo_increase_plan_quota_operation.go`
- **New file**: `tencentcloud/services/teo/resource_tc_teo_increase_plan_quota_operation_test.go`
- **New file**: `tencentcloud/services/teo/resource_tc_teo_increase_plan_quota_operation.md`
- **Modified**: `tencentcloud/provider.go` - register the new resource
- **Modified**: `tencentcloud/provider.md` - register the new resource in the provider doc
- **Dependency**: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901` (already in vendor)