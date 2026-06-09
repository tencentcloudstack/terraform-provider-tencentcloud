## Why

Terraform Provider for TencentCloud currently lacks support for managing TEO (EdgeOne) prefetch tasks. Users need the ability to submit URL prefetch tasks through Terraform to pre-warm cache on EdgeOne nodes, which reduces latency for end-users accessing new or updated content. This is a common operational need when deploying new content or migrating to EdgeOne.

## What Changes

- Add new Terraform resource `tencentcloud_teo_prefetch_task` of type `RESOURCE_KIND_OPERATION` to support creating TEO prefetch cache tasks
- The resource calls `CreatePrefetchTask` API to submit prefetch tasks and polls `DescribePrefetchTasks` API until the task completes (not in `processing` status)
- Since this is an OPERATION resource, RUD (Read/Update/Delete) methods are empty — the resource represents a one-time operation with no persistent state to manage
- Register the new resource in `tencentcloud/provider.go` and `tencentcloud/provider.md`
- Add unit tests using gomonkey mock approach
- Add resource documentation `.md` file

## Capabilities

### New Capabilities
- `teo-prefetch-task-operation`: Supports creating TEO prefetch cache tasks with configurable zone ID, target URLs, prefetch mode, headers, and media segment prefetch control. Polls for task completion status via DescribePrefetchTasks API.

### Modified Capabilities
<!-- No existing capabilities are being modified -->

## Impact

- New files in `tencentcloud/services/teo/`: `resource_tc_teo_prefetch_task_operation.go`, `resource_tc_teo_prefetch_task_operation_test.go`
- Modified files: `tencentcloud/provider.go`, `tencentcloud/provider.md`
- New doc file: `tencentcloud/services/teo/resource_tc_teo_prefetch_task_operation.md`
- Depends on cloud API: `CreatePrefetchTask` and `DescribePrefetchTasks` from `teo/v20220901` SDK package
