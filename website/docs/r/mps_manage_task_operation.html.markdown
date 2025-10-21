---
subcategory: "Media Processing Service(MPS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mps_manage_task_operation"
sidebar_current: "docs-tencentcloud-resource-mps_manage_task_operation"
description: |-
  Provides a resource to create a mps manage_task_operation
---

# tencentcloud_mps_manage_task_operation

Provides a resource to create a mps manage_task_operation

## Example Usage

```hcl
resource "tencentcloud_mps_manage_task_operation" "operation" {
  operation_type = "Abort"
  task_id        = "2600010949-LiveScheduleTask-xxxx"
}
```

## Argument Reference

The following arguments are supported:

* `operation_type` - (Required, String, ForceNew) Operation type. Valid values:`Abort`: task termination. Notice: If the task type is live stream processing (LiveStreamProcessTask), tasks whose task status is `WAITING` or `PROCESSING` can be terminated.For other task types, only tasks whose task status is `WAITING` can be terminated.
* `task_id` - (Required, String, ForceNew) Video processing task ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



