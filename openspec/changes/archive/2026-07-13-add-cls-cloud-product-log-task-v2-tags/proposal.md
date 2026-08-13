## Why

The `tencentcloud_cls_cloud_product_log_task_v2` resource currently does not support the `tags` parameter, even though the cloud API (`CreateCloudProductLogCollection` and `ModifyCloudProductLogCollection`) already accepts a `Tags` field to bind tags to the associated logset and topic. Users cannot manage tags on cloud product log collection tasks through Terraform, forcing them to tag resources manually via the console. Adding the `tags` parameter brings the resource to feature parity with the cloud API.

## What Changes

- Add a new optional `tags` parameter (TypeMap, string elements) to the `tencentcloud_cls_cloud_product_log_task_v2` resource schema. Tags are bound to the logset and topic created/associated by the cloud product log collection task.
- Pass `tags` to the `CreateCloudProductLogCollection` API request on resource creation so tags are applied at creation time.
- Pass `tags` to the `ModifyCloudProductLogCollection` API request on resource update when `tags` changes, so tags can be modified in-place without recreating the resource.
- Read back tags from the `DescribeCloudProductLogTasks` response (`TopicTags` / `LogsetTags`) and/or from the `DescribeClsLogset` / `DescribeClsTopicById` responses during the Read operation to populate state.
- Update the resource documentation (`resource_tc_cls_cloud_product_log_task_v2.md`) with a tags usage example.
- Add unit test cases (using gomonkey mocks) covering the tags parameter in create and update flows.

## Capabilities

### New Capabilities
- `cls-cloud-product-log-task-v2-tags`: Support for managing tags on the CLS cloud product log collection task v2 resource, including create, update, and read of the `tags` parameter.

### Modified Capabilities
<!-- No existing capability requirements are being modified. -->

## Impact

- **Resource code**: `tencentcloud/services/cls/resource_tc_cls_cloud_product_log_task_v2.go` — schema definition, create/update/read logic for `tags`.
- **Resource documentation**: `tencentcloud/services/cls/resource_tc_cls_cloud_product_log_task_v2.md` — add tags example.
- **Resource tests**: `tencentcloud/services/cls/resource_tc_cls_cloud_product_log_task_v2_test.go` — add unit test cases for tags (using gomonkey mocks, not the terraform test suite).
- **Cloud API**: `CreateCloudProductLogCollection` (Tags input), `ModifyCloudProductLogCollection` (Tags input), `DescribeCloudProductLogTasks` (TopicTags/LogsetTags output), `DescribeLogsets`/`DescribeTopics` (Tags output). All fields already exist in the vendored SDK; no SDK upgrade is required.
- **Backward compatibility**: The new `tags` parameter is optional; existing configurations and state are unaffected.
