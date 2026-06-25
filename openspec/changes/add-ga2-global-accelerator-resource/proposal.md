## Why

Tencent Cloud Global Accelerator V2 (`ga2`) currently exposes only the `tencentcloud_ga2_endpoint_group` resource in this provider. Users cannot manage the top-level **Global Accelerator instance** (the parent object that owns listeners and endpoint groups) through Terraform, which forces them to create the accelerator instance manually in the console before any Terraform-managed listener/endpoint-group resource can be wired to it. Adding `tencentcloud_ga2_global_accelerator` closes this gap and makes the `ga2` workflow fully Terraform-native.

## What Changes

- Add a new resource `tencentcloud_ga2_global_accelerator` backed by the `ga2` v20250115 SDK.
- Map all `CreateGlobalAccelerator` request parameters to schema fields, in particular: `name`, `instance_charge_type`, `description`, `cross_border_type`, `cross_border_promise_flag`, `tags`.
- Implement async-aware CRUD: `CreateGlobalAccelerator`, `ModifyGlobalAccelerator`, `DeleteGlobalAccelerator` all return a `TaskId` that must be polled via `DescribeTaskResult` until `Status == "SUCCESS"`. Reuse the existing `Ga2Service.WaitForGa2TaskFinish(ctx, taskId, timeout)` helper.
- Add a new service helper `Ga2Service.DescribeGa2GlobalAcceleratorById(ctx, gaId) (*GlobalAcceleratorSet, error)` that wraps `DescribeGlobalAccelerators` with a `global-accelerator-id` filter and pagination capped at the documented maximum (`Limit=100`).
- Surface read-only computed fields from the describe response: `state`, `cname`, `ddos_id`, `create_time`, `listener_counts`, `accelerator_area_counts`, `status`.
- Wire the new resource into `tencentcloud/provider.go` under the `ga2` namespace.
- Author resource markdown documentation `resource_tc_ga2_global_accelerator.md` (example HCL snippet + `terraform import` syntax).
- Author acceptance-test scaffolding `resource_tc_ga2_global_accelerator_test.go`.
- Resource ID is the bare `GlobalAcceleratorId` returned by `CreateGlobalAccelerator` (no composite ID).
- All SDK calls are wrapped with `resource.Retry` (write paths use `tccommon.WriteRetryTimeout`, read paths use `tccommon.ReadRetryTimeout`).
- Tags are forwarded directly via `CreateGlobalAcceleratorRequest.Tags` on Create. On Update, because `ModifyGlobalAcceleratorRequest` does not accept tags, tag changes are reconciled through the standard `svctag.NewTagService(...).ModifyTags(...)` helper using resource name `qcs::ga2:<region>:uin/:globalAccelerator/<gaId>`.

## Capabilities

### New Capabilities
- `ga2-global-accelerator-resource`: Lifecycle management (create / read / update / delete / import) of a Tencent Cloud Global Accelerator V2 instance, including async task polling, full schema parity with `CreateGlobalAccelerator`, and tag management.

### Modified Capabilities
<!-- None: this change only introduces a new resource; it does not alter requirement-level behavior of any existing capability. -->

## Impact

- **New code**:
  - `tencentcloud/services/ga2/resource_tc_ga2_global_accelerator.go` (CRUD + schema + build/flatten helpers, single file, mirroring `tencentcloud_igtm_monitor` style).
  - `tencentcloud/services/ga2/resource_tc_ga2_global_accelerator.md` (resource doc + import syntax, mirroring `resource_tc_config_compliance_pack.md` filename convention).
  - `tencentcloud/services/ga2/resource_tc_ga2_global_accelerator_test.go` (acceptance test skeleton, mirroring `resource_tc_config_compliance_pack_test.go` filename convention).
- **Modified code**:
  - `tencentcloud/services/ga2/service_tencentcloud_ga2.go`: add `DescribeGa2GlobalAcceleratorById`.
  - `tencentcloud/provider.go`: register `tencentcloud_ga2_global_accelerator` in `ResourcesMap`.
- **APIs consumed**: `CreateGlobalAccelerator`, `DescribeGlobalAccelerators`, `ModifyGlobalAccelerator`, `DeleteGlobalAccelerator`, `DescribeTaskResult` (already vendored in `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ga2/v20250115/`).
- **No breaking change**: purely additive; no existing schema or state is modified.
- **No SDK upgrade required**: all required APIs are already present in the vendored SDK.
