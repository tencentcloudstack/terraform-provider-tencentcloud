## Why

TencentCloud Global Accelerator (GA2) currently lacks Terraform support for managing accelerate areas. Users need to manage accelerate areas (加速地域) programmatically through Infrastructure as Code to automate their global acceleration network topology.

## What Changes

- Add a new Terraform resource `tencentcloud_ga2_accelerate_area` that manages the full lifecycle (CRUD) of accelerate areas under a GA2 global accelerator instance.
- The resource uses asynchronous cloud APIs (CreateAccelerateAreas, DescribeAccelerateAreas, ModifyAccelerateAreas, DeleteAccelerateAreas) with task polling via DescribeTaskResult.
- Register the new resource in `tencentcloud/provider.go` and `tencentcloud/provider.md`.

## Capabilities

### New Capabilities
- `ga2-accelerate-area`: Terraform resource for creating, reading, updating, and deleting accelerate areas under a GA2 global accelerator instance, with async task polling support.

### Modified Capabilities

## Impact

- New files: `tencentcloud/services/ga2/resource_tc_ga2_accelerate_area.go`, `tencentcloud/services/ga2/resource_tc_ga2_accelerate_area_test.go`, `tencentcloud/services/ga2/resource_tc_ga2_accelerate_area.md`, `tencentcloud/services/ga2/service_tencentcloud_ga2.go`
- Modified files: `tencentcloud/provider.go`, `tencentcloud/provider.md`
- Dependencies: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ga2/v20250115` (already vendored)
