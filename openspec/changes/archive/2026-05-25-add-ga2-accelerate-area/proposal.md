## Why

TencentCloud Global Accelerator (GA2) currently lacks Terraform support for managing accelerate areas. Users need to manage accelerate areas (加速地域) programmatically through Infrastructure as Code to automate their global acceleration configurations.

## What Changes

- Add a new Terraform resource `tencentcloud_ga2_accelerate_area` that manages the full lifecycle (CRUD) of accelerate areas under a GA2 global accelerator instance.
- The resource uses four cloud APIs: `CreateAccelerateAreas`, `DescribeAccelerateAreas`, `ModifyAccelerateAreas`, and `DeleteAccelerateAreas` from the `ga2/v20250115` SDK package.
- All write operations (Create/Modify/Delete) are asynchronous, returning a TaskId. After each write, the resource polls `DescribeAccelerateAreas` to confirm the operation has taken effect.
- Register the new resource in `provider.go` and `provider.md`.

## Capabilities

### New Capabilities
- `ga2-accelerate-area`: CRUD resource for managing accelerate areas under a GA2 global accelerator instance, including region, bandwidth, ISP type, and IP version configuration.

### Modified Capabilities

(none)

## Impact

- New files:
  - `tencentcloud/services/ga2/resource_tc_ga2_accelerate_area.go`
  - `tencentcloud/services/ga2/resource_tc_ga2_accelerate_area_test.go`
  - `tencentcloud/services/ga2/service_tencentcloud_ga2.go` (if not existing)
  - `tencentcloud/services/ga2/resource_tc_ga2_accelerate_area.md`
- Modified files:
  - `tencentcloud/provider.go` (register resource)
  - `tencentcloud/provider.md` (document resource)
- Dependencies: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ga2/v20250115` (already in vendor)
