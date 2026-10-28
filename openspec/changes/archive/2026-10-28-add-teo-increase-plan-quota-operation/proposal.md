## Why

TEO (EdgeOne) users need the ability to increase plan quotas (e.g., site count, precise access control rules, rate limiting rules) through Terraform. Currently, the terraform provider lacks this capability, forcing users to manually use the console or API to purchase additional quotas. This change adds a Terraform framework Action for the `IncreasePlanQuota` API, enabling Infrastructure-as-Code management of TEO plan quota upgrades.

## What Changes

- Add a new framework Action `tencentcloud_teo_increase_plan_quota` (terraform-plugin-framework `action.Action`) that calls the TEO `IncreasePlanQuota` API
- The action is a one-time operation: it invokes the API on `terraform apply` and reports the returned `deal_name` (order number) back to the user via a warning diagnostic. No state is persisted.
- Input parameters: `plan_id`, `quota_type`, `quota_number` (all required)
- Output parameter: `deal_name` (computed) - the order number returned by the API
- Register the new action in `tencentcloud/framework/registry.go` (actionFactories) and list it under the TEO `Action` section in `tencentcloud/provider.md`

## Capabilities

### New Capabilities
- `teo-increase-plan-quota`: Terraform framework Action for the TEO `IncreasePlanQuota` API, enabling users to increase plan quotas for TEO EdgeOne plans via Infrastructure-as-Code

### Modified Capabilities
<!-- None -->

## Impact

- **New file**: `tencentcloud/services/teo/resource_tc_teo_increase_plan_quota_operation.go` - framework action implementation (`action.Action`)
- **New file**: `tencentcloud/services/teo/resource_tc_teo_increase_plan_quota_operation_test.go` - unit tests using gomonkey to mock the TEO client
- **New file**: `tencentcloud/services/teo/action_tc_teo_increase_plan_quota.md` - action documentation source for `make doc`
- **Removed file**: `tencentcloud/services/teo/resource_tc_teo_increase_plan_quota_operation.md` (no longer an SDKv2 resource)
- **Modified**: `tencentcloud/framework/registry.go` - register the new action in `actionFactories`
- **Modified**: `tencentcloud/provider.md` - list the action under the TEO `Action` section
- **Not modified**: `tencentcloud/provider.go` - framework actions are registered via `framework/registry.go`, not the SDKv2 `ResourcesMap`
- **Dependency**: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901` (already in vendor)