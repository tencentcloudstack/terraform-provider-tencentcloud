## Why

The TencentCloud CLS (Cloud Log Service) product supports cross-account/cross-region log delivery tasks, but the Terraform provider currently has no resource to manage this capability. Users must create and manage CLS deliver tasks manually via the console or API, causing infrastructure drift and preventing these tasks from being managed as code. Adding a `tencentcloud_cls_cls_deliver_task` resource closes this gap, letting users manage the full lifecycle (create, read, update, delete) of CLS deliver tasks through Terraform.

## What Changes

- Add a new resource `tencentcloud_cls_cls_deliver_task` of type RESOURCE_KIND_GENERAL under the `cls` service, implementing full CRUD via the cloud APIs `CreateCLSDeliverTask`, `DescribeCLSDeliverTasks`, `ModifyCLSDeliverTask`, and `DeleteCLSDeliverTask`.
- Define schema with top-level fields (`task_name`, `compliance`, `has_services_log`, `enable`) and nested blocks (`source_topic_config`, `target_topic_config`, `deliver_rule`) plus computed outputs (`task_id`, `status`, `progress`, `create_time`, `update_time`, `uin`, `total`).
- Register the resource in `tencentcloud/provider.go` and `tencentcloud/provider.md`.
- Add a service-layer `DescribeClsClsDeliverTaskById` helper in `service_tencentcloud_cls.go` that queries `DescribeCLSDeliverTasks` with a `taskId` filter.
- Add a unit test file using gomonkey mocks.
- Add a resource example `.md` file.

## Capabilities

### New Capabilities
- `cls-cls-deliver-task`: Manage the lifecycle of a CLS cross-account/cross-region deliver task, including source topic selection, target topic configuration, deliver rule, compliance consent and service log delivery toggle.

### Modified Capabilities
<!-- None - this is a new resource with no existing spec-level behavior changes. -->

## Impact

- **New files**:
  - `tencentcloud/services/cls/resource_tc_cls_cls_deliver_task.go` (resource CRUD + schema)
  - `tencentcloud/services/cls/resource_tc_cls_cls_deliver_task_test.go` (gomonkey unit tests)
  - `tencentcloud/services/cls/resource_tc_cls_cls_deliver_task.md` (example doc)
  - `openspec/changes/add-cls-cls-deliver-task/specs/cls-cls-deliver-task/spec.md`
- **Modified files**:
  - `tencentcloud/services/cls/service_tencentcloud_cls.go` (add `DescribeClsClsDeliverTaskById`)
  - `tencentcloud/provider.go` (register `tencentcloud_cls_cls_deliver_task`)
  - `tencentcloud/provider.md` (list `tencentcloud_cls_cls_deliver_task`)
- **APIs**: Uses `cls/v20201016` SDK package - `CreateCLSDeliverTask`, `DescribeCLSDeliverTasks`, `ModifyCLSDeliverTask`, `DeleteCLSDeliverTask`. No SDK version bump needed (already vendored).
- **Dependencies**: None new. Backward compatible (purely additive).