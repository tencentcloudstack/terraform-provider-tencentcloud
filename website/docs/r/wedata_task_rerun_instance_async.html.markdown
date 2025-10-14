---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_task_rerun_instance_async"
sidebar_current: "docs-tencentcloud-resource-wedata_task_rerun_instance_async"
description: |-
  Provides a resource to create a wedata task re-run instance
---

# tencentcloud_wedata_task_rerun_instance_async

Provides a resource to create a wedata task re-run instance

## Example Usage

```hcl
resource "tencentcloud_wedata_task_rerun_instance_async" "wedata_task_rerun_instance_async" {
  project_id        = "1859317240494305280"
  instance_key_list = ["20250324192240178_2025-10-13 16:20:00"]
}
```

## Argument Reference

The following arguments are supported:

* `instance_key_list` - (Required, Set: [`String`], ForceNew) Instance id list, which can be obtained from ListInstances.
* `project_id` - (Required, String, ForceNew) Project Id.
* `check_parent_type` - (Optional, String, ForceNew) Whether to check upstream tasks: ALL, MAKE_SCOPE (selected), NONE (do not check anything), default is NONE.
* `redefine_parallel_num` - (Optional, Int, ForceNew) Customize the instance running concurrency. If not configured, the original self-dependency of the task will be used.
* `redefine_param_list` - (Optional, List, ForceNew) Re-run instance custom parameters.
* `redefine_self_workflow_dependency` - (Optional, String, ForceNew) Customized workflow self-dependence: yes to enable, no to disable. If not configured, the original workflow self-dependence will be used.
* `rerun_type` - (Optional, String, ForceNew) Rerun type, 1: self; 3: children; 2: self and children, default 1.
* `skip_event_listening` - (Optional, Bool, ForceNew) Whether to ignore event monitoring when rerunning.
* `son_range_type` - (Optional, String, ForceNew) Downstream instance scope WORKFLOW: workflow PROJECT: project ALL: all cross-workflow dependent projects, default WORKFLOW.

The `redefine_param_list` object supports the following:

* `k` - (Optional, String) Key.
* `v` - (Optional, String) Value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



