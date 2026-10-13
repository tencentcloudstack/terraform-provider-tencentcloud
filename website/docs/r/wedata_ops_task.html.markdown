---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_ops_task"
sidebar_current: "docs-tencentcloud-resource-wedata_ops_task"
description: |-
  Provides a resource to create a wedata ops task
---

# tencentcloud_wedata_ops_task

Provides a resource to create a wedata ops task

## Example Usage

```hcl
resource "tencentcloud_wedata_ops_task" "wedata_ops_task" {
  project_id = "1859317240494305280"
  task_id    = "20251013154418424"
  action     = "START"
}
```

## Argument Reference

The following arguments are supported:

* `action` - (Required, String) Action. Valid values: `START`, `PAUSE`.
* `project_id` - (Required, String, ForceNew) Project Id.
* `task_id` - (Required, String, ForceNew) Task id.
* `enable_data_backfill` - (Optional, Bool) Whether to re-record the intermediate instance from the last pause to the current one when starting. The default value is false, which means no re-recording.
* `kill_instance` - (Optional, Bool) Whether required to terminate the generated instance.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `status` - Task status.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [GetOpsAsyncJob](https://cloud.tencent.com/document/api/1267/123587) - 查询运维中心异步操作详情
* Create: [PauseOpsTasksAsync](https://cloud.tencent.com/document/api/1267/123569) - 异步批量暂停任务
* Create: [StartOpsTasks](https://cloud.tencent.com/document/api/1267/124031) - 异步批量启动任务
* Read: [GetOpsTask](https://cloud.tencent.com/document/api/1267/123586) - 获取任务详情
* Update: [GetOpsAsyncJob](https://cloud.tencent.com/document/api/1267/123587) - 查询运维中心异步操作详情
* Update: [PauseOpsTasksAsync](https://cloud.tencent.com/document/api/1267/123569) - 异步批量暂停任务
* Update: [StartOpsTasks](https://cloud.tencent.com/document/api/1267/124031) - 异步批量启动任务

