## Why

Activating the CLS (Cloud Log Service) for an account is a prerequisite for using CLS resources, but there is currently no Terraform way to trigger this one-time activation. Users must open the service manually in the console before managing CLS resources.

## What Changes

- Add a new **operation-type** (one-time action) resource `tencentcloud_cls_open_service_operation` that activates CLS for the current account.
- Create lifecycle calls `OpenClsService` (no input parameters) to open the service.
- The resource unique ID is auto-generated (via `helper.BuildToken()`).
- Read lifecycle calls `GetClsService` to expose the current activation status as a computed attribute.
- Delete lifecycle is a no-op (a one-time activation cannot be reverted via API), only removing the resource from state.
- Provider registration, resource example doc (`.md`), website documentation, and unit test are added.

## Capabilities

### New Capabilities
- `cls-open-service-operation`: A one-time operation resource that activates CLS for the account and surfaces the activation status.

### Modified Capabilities
<!-- None: brand-new resource, no existing requirements change. -->

## Impact

- New file: `tencentcloud/services/cls/resource_tc_cls_open_service_operation.go`
- New file: `tencentcloud/services/cls/resource_tc_cls_open_service_operation.md`
- New file: `tencentcloud/services/cls/resource_tc_cls_open_service_operation_test.go`
- New file: `website/docs/r/cls_open_service_operation.html.markdown` (generated via `make doc`)
- Modified: `tencentcloud/provider.go` (register resource), `tencentcloud/provider.md`, `website/tencentcloud.erb`
- SDK: uses existing `cls/v20201016` APIs `OpenClsService` and `GetClsService` (already present in vendored SDK — no SDK changes required).
- No breaking changes; purely additive.
