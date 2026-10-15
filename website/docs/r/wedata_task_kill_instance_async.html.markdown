---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_task_kill_instance_async"
sidebar_current: "docs-tencentcloud-resource-wedata_task_kill_instance_async"
description: |-
  Provides a resource to create a wedata task kill instance
---

# tencentcloud_wedata_task_kill_instance_async

Provides a resource to create a wedata task kill instance

## Example Usage

```hcl
resource "tencentcloud_wedata_task_kill_instance_async" "wedata_task_kill_instance_async" {
  project_id        = "1859317240494305280"
  instance_key_list = ["20250324192240178_2025-10-13 17:00:00"]
}
```

## Argument Reference

The following arguments are supported:

* `instance_key_list` - (Required, Set: [`String`], ForceNew) Instance id list, which can be obtained from ListInstances.
* `project_id` - (Required, String, ForceNew) Project Id.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [GetOpsAsyncJob](https://cloud.tencent.com/document/api/1267/123587) - 查询运维中心异步操作详情
* Create: [KillTaskInstancesAsync](https://cloud.tencent.com/document/api/1267/123581) - 实例批量终止操作(异步)

