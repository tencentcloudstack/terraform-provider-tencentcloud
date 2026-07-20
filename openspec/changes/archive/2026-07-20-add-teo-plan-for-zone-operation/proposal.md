## Why

EdgeOne (TEO) zones that have not yet purchased a plan cannot be fully managed via Terraform. Currently, users must purchase a plan for a zone manually in the console before they can manage the zone's configuration through Terraform. The `CreatePlanForZone` API exists to programmatically purchase a plan for an unbound zone, but there is no Terraform resource exposing this action. Adding a one-time operation resource lets users purchase a plan for a zone entirely from infrastructure-as-code.

## What Changes

- Add a new **operation-type** (one-time action) resource `tencentcloud_teo_plan_for_zone` that purchases a plan for a zone that has no bound plan, via the `CreatePlanForZone` API.
- Create lifecycle calls `CreatePlanForZone` with `ZoneId` and `PlanType` inputs and exposes `ResourceNames` and `DealNames` from the response as computed attributes.
- The resource unique ID is auto-generated (via `helper.BuildToken()`).
- Read, Update, and Delete lifecycles are no-ops (one-time operation; no state to query or revert).
- Provider registration, resource example doc (`.md`), website documentation, and unit test are added.

## Capabilities

### New Capabilities
- `teo-plan-for-zone-operation`: A one-time operation resource that purchases a TEO plan for a zone via the `CreatePlanForZone` API.

### Modified Capabilities
<!-- None: brand-new resource, no existing requirements change. -->

## Impact

- New file: `tencentcloud/services/teo/resource_tc_teo_plan_for_zone_operation.go`
- New file: `tencentcloud/services/teo/resource_tc_teo_plan_for_zone_operation.md`
- New file: `tencentcloud/services/teo/resource_tc_teo_plan_for_zone_operation_test.go`
- New file: `website/docs/r/teo_plan_for_zone.html.markdown` (generated via `make doc`)
- Modified: `tencentcloud/provider.go` (register resource), `tencentcloud/provider.md`, `website/tencentcloud.erb`
- SDK: uses existing `teo/v20220901` API `CreatePlanForZone` (already present in vendored SDK — no SDK changes required). The API is synchronous; no async polling is needed.
- No breaking changes; purely additive.
