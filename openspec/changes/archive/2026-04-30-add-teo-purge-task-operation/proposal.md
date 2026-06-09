## Why

EdgeOne (TEO) currently lacks a Terraform resource to trigger cache purge operations. Users need the ability to manage cache invalidation through Infrastructure as Code, enabling automated cache clearing when origin resources are updated. This is a common operational need for CDN/EdgeOne users.

## What Changes

- Add a new Terraform OPERATION resource `tencentcloud_teo_purge_task` that triggers a cache purge task via the `CreatePurgeTask` API
- The resource will support multiple purge types: `purge_url`, `purge_prefix`, `purge_host`, `purge_all`, `purge_cache_tag`
- After creating the purge task, poll `DescribePurgeTasks` to verify task completion and return task results as computed attributes
- Register the new resource in `provider.go` and `provider.md`
- Add corresponding unit tests and documentation

## Capabilities

### New Capabilities
- `teo-purge-task-operation`: Triggers and tracks EdgeOne cache purge tasks via Terraform, supporting URL/prefix/host/all/cache-tag purge types with async status polling

### Modified Capabilities
<!-- No existing capabilities are being modified -->

## Impact

- **New files**: `tencentcloud/services/teo/resource_tc_teo_purge_task_operation.go`, test file, and documentation
- **Modified files**: `tencentcloud/provider.go` (resource registration), `tencentcloud/provider.md` (documentation index)
- **API dependencies**: `CreatePurgeTask` and `DescribePurgeTasks` from `teo/v20220901` SDK package
- **No breaking changes**: This is a purely additive change
