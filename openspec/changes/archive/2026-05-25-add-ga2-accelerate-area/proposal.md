## Why

TencentCloud Global Accelerator 2.0 (GA2) provides acceleration area management capabilities through cloud APIs, but there is currently no Terraform resource to manage GA2 accelerate areas. Users need a Terraform resource to create, read, update, and delete accelerate areas for their GA2 instances in an infrastructure-as-code workflow.

## What Changes

- Add a new Terraform resource `tencentcloud_ga2_accelerate_area` that manages the full lifecycle (CRUD) of GA2 accelerate areas.
- The resource uses four cloud APIs: `CreateAccelerateAreas`, `DescribeAccelerateAreas`, `ModifyAccelerateAreas`, and `DeleteAccelerateAreas`.
- All write operations (Create/Modify/Delete) are asynchronous (return TaskId), so the resource must poll via `DescribeAccelerateAreas` until the operation takes effect.
- Register the new resource in `provider.go` and `provider.md`.
- Add corresponding documentation (.md example file).

## Capabilities

### New Capabilities
- `ga2-accelerate-area-resource`: Terraform resource for managing GA2 accelerate areas with full CRUD lifecycle, including async operation polling.

### Modified Capabilities

(none)

## Impact

- New files: `tencentcloud/services/ga2/resource_tc_ga2_accelerate_area.go`, `tencentcloud/services/ga2/resource_tc_ga2_accelerate_area_test.go`, `tencentcloud/services/ga2/resource_tc_ga2_accelerate_area.md`, `tencentcloud/services/ga2/service_tencentcloud_ga2.go`
- Modified files: `tencentcloud/provider.go`, `tencentcloud/provider.md`
- Dependencies: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ga2/v20250115` (already in vendor)
